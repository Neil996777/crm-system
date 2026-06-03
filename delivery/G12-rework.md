# G12 Rework Package — CRM System (Claude → Codex KICKBACK)

## Document Control

- Project: CRM System
- Source: G12 independent audit, 2026-06-03 (`archive/reviews/g12-audit/g12-audit-decision-2026-06-03.md`)
- Decision: **REWORK** — G11 is `Gate Blocked`. Do NOT release.
- Executor: Codex (G9/G10/G11 rework), then return to Claude for focused G12 re-audit.
- This is a self-contained, directly-executable kickback package. Read it first.

## Rules for this rework

- Fix the items below; do **not** re-plan, re-design, or change committed decisions
  (decision-log DEC-001..022) or the model. No-downgrade: a P0/P1 finding may only be
  fixed, not weakened/merged away.
- Every fix is **TDD**: add a failing test that reproduces the gap, then make it pass.
  Never delete or weaken an existing test to go green.
- Real persistence / real evidence only: no mock/stub/TODO, and no hand-typed evidence
  where a tool can capture it. Infra evidence must be a captured command transcript.
- When an item is fixed, update its task Status in `delivery/tasks.md`, fill the
  `modeling/traceability-matrix.md` Task/Test/Integration columns, set the matching
  blocker in `planning/blockers.md` to `Resolved` with the deciding artifact, and commit.
- When all BLOCKER + MAJOR items are done and re-evidenced, return control to Claude.
  Claude re-audits dimensions 4 (boundary/security) and 5 (release evidence) plus an
  independent compile/test re-run. Do not self-pass G12.

## Priority order

Fix BLOCKERs first (release-blocking), then MAJORs, then MINORs. Suggested sequence:
BLK-G12-001 → 004 (same code path) → 002 → 005 (same service) → 003 → 007 → 006 → 008.

---

## 🔴 BLK-G12-001 (BLOCKER, P0) — Deliver audit events to audit-history

**Where:** `services/opportunity`, `services/commercial`, `services/work`, `services/account`
(cmd/server/main.go and the mutation handlers using `outbox.Append`).

**What's wrong:** These services write `EVT-*` events only to their local
`<schema>.outbox_events` table. No relay/dispatcher delivers them to
audit-history-service, so the required audit records never reach the append-only store.

**Contracts violated:** AUD-IMM-002; `docs/security/audit-log-spec.md` Event Catalog
(EVT-STAGE-CHANGED, EVT-OPPORTUNITY-WON/LOST, EVT-QUOTE-ACCEPTED, EVT-CONTRACT-SIGNED,
EVT-PAYMENT-RECORDED, EVT-OWNER-CHANGED, EVT-RECORD-ARCHIVED); ACC-014, ACC-022;
`integration-design.md` Event Delivery Strategy (outbox + dispatcher).

**Required fix:** Implement the outbox relay that was specified at G5
(`integration-design.md`: "database outbox table per producing service plus background
dispatcher"). A background dispatcher in each producing service must read unsent
`outbox_events` rows and deliver them to audit-history `/internal/events/append`
(with the S2S signed token), marking rows delivered, retrying on failure by event UID,
and tolerating duplicates (idempotent by event UID at the consumer). Bring the four
services to the same durable behavior already used by lead/import-export/identity-authz.

**Acceptance / verify (TDD):**
- Integration test (real Postgres testcontainer) per service: perform a sensitive
  mutation (e.g. opportunity stage change, contract sign, payment recorded, owner change,
  archive) and assert the corresponding `EVT-*` event is queryable in audit-history within
  the dispatcher cycle. Tag with the EVT id.
- Test the failure path: audit-history append fails → event stays undelivered and is
  retried (no silent loss), and the duplicate delivery is de-duplicated by event UID.
- Record the event-id → service → consumer mapping in the traceability matrix.

---

## 🔴 BLK-G12-004 (MAJOR, paired with 001) — Transactional outbox + no dropped error

**Where:** `services/opportunity/internal/handler/opportunity_command.go:84-99` and the
same pattern in commercial/work/account mutation handlers.

**What's wrong:** The business mutation (`repo.Create`) and `outbox.Append` use separate
pool connections (not one transaction), and the outbox error is discarded
(`_ = h.outbox.Append(...)`). A 201 can be returned with no event row.

**Required fix:** Write the business row and its outbox event in the **same DB
transaction**, and never discard the append error (fail the request / roll back if the
event cannot be enqueued). This is the "same durable workflow" guarantee of AUD-IMM-002.

**Acceptance / verify:** Test that injecting an outbox-write failure rolls back the
mutation (record absent) and returns an error — i.e. you can never observe a persisted
business change without its enqueued event.

---

## 🔴 BLK-G12-002 (BLOCKER) — Add S2S auth to reporting internal write

**Where:** `services/reporting/internal/handler/overview_query.go:61-72`
(`POST /internal/projections`).

**What's wrong:** The endpoint has no `verifyServiceToken`. Any caller reaching the
service can inject/overwrite report read-models.

**Contracts violated:** `authz-architecture.md` Service-To-Service Authorization;
SEC-SVC-BLK-002.

**Required fix:** Apply the same signed-service-token middleware used on other
`/internal/*` endpoints (audience = reporting, allowed intent for projection ingest,
expiry ≤ 5 min, correlation id). Reject missing/expired/wrong-audience/disallowed-intent
with `SERVICE_AUTH_FAILED`.

**Acceptance / verify (TDD):** Negative tests — no token / expired / wrong audience /
wrong intent all return `SERVICE_AUTH_FAILED` and mutate no projection; positive test with
a valid token succeeds.

---

## 🔴 BLK-G12-005 (MAJOR, paired with 002) — Wire the projection ingest path

**Where:** `services/reporting` (no producer actually calls `/internal/projections`).

**What's wrong:** The projection store is never populated, so report authorization is
enforced over empty/untrusted data; PM-043..045 report correctness is unverifiable.

**Required fix:** Connect the source domain events (lead/account/opportunity/commercial)
to the reporting projection (via the outbox dispatcher from BLK-G12-001 or the agreed
event path) so the read-model is actually built from authorized domain events, per
`architecture.md` (reporting from owned read model / events, not cross-service table reads).

**Acceptance / verify:** Integration test — emit domain events, assert the reporting
read-model updates and that an authorized report query (Admin/Manager scope) returns the
expected aggregate while Sales is denied with no leakage (PM-043..045).

---

## 🔴 BLK-G12-003 (BLOCKER, P0, ACC-017) — Substantiate security-group closure

**Where:** `docs/release/evidence/` (reopens BLK-G11-002).

**What's wrong:** The only raw `DescribeSecurityGroupAttributes` export still shows TCP
8088/8443/3389 open to `0.0.0.0/0`. The "post-cleanup" files asserting removal are
hand-authored and contradict the raw export. The CRM instance also shares the `Default`
security group (all-protocol self-referential allow).

**Contracts violated:** ACC-017; `deployment-notes.md` Network Exposure; ARCH-ACC-013/014.

**Required fix:**
1. Actually remove the non-CRM public rules (8088/8443/3389) and **re-export the raw
   `DescribeSecurityGroupAttributes` API response after cleanup** showing they are gone.
   The raw API export — not a hand-written summary — is the deciding artifact.
2. Confirm CRM `8080` and PostgreSQL `5432` are absent from any `0.0.0.0/0` rule in the
   raw export.
3. Move CRM to a **dedicated least-exposure security group** off the shared `Default`
   (BLK-G12-007), exposing only 22/80/443 publicly.

**Acceptance / verify:** Post-cleanup raw API export on disk with no 8088/8443/3389 and no
public 8080/5432, plus external negative probes (below in BLK-G12-006). Reconcile
`planning/blockers.md` BLK-G11-002 / BLK-G12-003 to the new artifact.

---

## 🟠 BLK-G12-007 (MAJOR) — Dedicated security group

Move the CRM instance to a dedicated security group (not shared `Default`). Public
ingress limited to 22/80/443; internal ports reachable only within the VPC/Docker
network. Capture the raw API export of the new group binding as evidence. (Do together
with BLK-G12-003.)

---

## 🟠 BLK-G12-006 (MAJOR) — Make release evidence independently verifiable

Replace hand-typed assertions with **captured command transcripts** under
`docs/release/evidence/`:

- **TLS / HTTPS:** `curl -vI https://118.196.44.193` (show HTTP/2 200 + HSTS),
  `openssl s_client -connect 118.196.44.193:443` and/or `certbot certificates`
  (issuer, notBefore/notAfter, SAN), and the HTTP→HTTPS redirect transcript.
- **Renewal:** captured `certbot renew --dry-run` output (the IP cert is ~6-day; renewal
  is safety-critical).
- **External negative probes:** from outside the host, evidence that
  `118.196.44.193:8080` and `:5432` are refused/timed out (e.g. `nc -vz` / `curl`).
- **Restore rehearsal (TASK-040):** a real transcript — `sha256sum -c` verify line,
  decrypt step, and post-restore catalog queries (`\du`, schema/row counts) — not just an
  operator-authored counts JSON.
- **Operator-access conditions 2 & 3:** `getent passwd crm-deploy crm-ops`, the sudoers
  boundary excerpt, `sshd_config` showing `PasswordAuthentication no` + root login policy,
  key fingerprints/ownership, and `ls -l` on the secrets + backup directories proving
  non-world-readable (700/600). The named least-privilege user must be the deploy actor,
  not root.

**Acceptance / verify:** Each of the five evidence items above is a captured artifact a
third party can re-check, referenced from TASK-039/040 and the ACC-017 evidence file.

---

## 🟡 BLK-G12-008 (MINOR) — Test traceability + hardening

- Add/re-tag the untagged P0 tests so every test-model family resolves 1:1:
  `TEST-ABUSE-MUTATE-001` (ABUSE-003), `TEST-AUTHZ-SCOPE-005` (denied mutation = no change),
  `TEST-PERM-CRUD-ADMIN-001 / -MGR-001 / -SALES-001` (PM-008/009/010 role-scoped CRUD allow),
  `TEST-OWNER-TRANSFER-002` (transfer) and `TEST-OWNER-TRANSFER-004` (open-work cascade,
  EDGE-024 / PIM-INV-030/033 — currently unverified).
- `services/work/internal/authz/service_token.go` — enforce the 5-minute max token
  lifetime cap like the other services.
- `identity-authz` — either enforce `AUTHZ_VERSION_STALE` on role change (compare authz
  version + revoke affected sessions) or remove the dead mechanism and document that
  per-request live user reload is the committed control.
- Remove the stray `截屏2026-06-03 16.32.12.png` from the repo root.

---

## Definition of Done for this rework

- BLK-G12-001..007 fixed with the named tests green and the named evidence captured.
- BLK-G12-008 closed (traceability 1:1; small hardening done).
- `delivery/tasks.md` statuses, `modeling/traceability-matrix.md` columns, and
  `planning/blockers.md` resolutions updated to the real artifacts; commits made.
- Affected QA (G10) and integration (G11) evidence re-run for the changed paths
  (audit-event delivery, reporting authz, security-group/TLS/restore/operator-access).
- Control returned to Claude for focused G12 re-audit (dimensions 4 + 5 + independent
  compile/test re-run). Do not self-pass G12.
