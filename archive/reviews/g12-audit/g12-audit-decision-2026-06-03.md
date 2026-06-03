# G12 (Independent Audit) — Gate Decision [KICKBACK: Claude → Codex]

## Document Control

- Project: CRM System
- Gate: G12 — Audit → Release/Rework
- Owner: Audit (Claude, audit platform), independent of execution
- Method: Five independent, parallel author≠auditor audits (read-only, evidence-based,
  charged to falsify Codex's "G9–G11 complete" claim), each requiring `file:line`
  evidence. A sixth dimension (independent compile/test re-run) is pending and is not
  required to reach this decision.
- Date: 2026-06-03
- Decision: **REWORK — Gate NOT Passed. Kicked back to Codex.** No release.
- Archive note: Audit evidence only. Not design authority.

## Scope audited

Codex returned at G11 claiming G9 (implementation), G10 (QA), G11 (integration) all
`Gate Passed`, all 40 tasks Done, and the two G11 release blockers (BLK-G11-001/002)
Resolved. G12 independently re-verified seven dimensions across `services/` (10 Go
services), `frontend/`, `deploy/`, `scripts/`, `docs/release/evidence/`, and the
governance registers.

## Per-dimension verdicts

| # | Dimension | Verdict | Summary |
|---|---|---|---|
| 1 | Completeness & traceability | PASS-WITH-CONCERNS | 40/40 tasks substantiated; 23/23 ACC genuinely backed; traceability matrix filled with real refs. Concern: infra evidence self-attested; runtime test re-run not yet done. |
| 2 | No-fakes / no-downgrade | PASS | Zero TODO/stub/mock on P0/P1 paths; real PostgreSQL persistence; 3 retired tests absent; post-scope-change model (Won=Signed, one quote, payment decoupled, no Status) enforced in code + DB constraints. |
| 3 | Test integrity / TDD | PASS-WITH-CONCERNS | Real testcontainers (real Postgres), no skip/weakening, denied actions assert no mutation. Concern: several P0 TEST-* IDs untagged (traceability gap). |
| 4 | Service boundary / security | **FAIL** | 2 BLOCKERs: audit events not delivered to audit-history; reporting internal write endpoint unauthenticated. |
| 5 | Release evidence (G11) | **FAIL** | Security-group "cleanup" unsubstantiated and contradicts the only raw API export; TLS/restore/operator-access evidence not independently verifiable. |
| 6 | Independent compile/test re-run | PENDING | Not required for this decision; recommended before re-audit close. |
| 7 | Model fidelity | PASS (folded into #2) | Implementation matches DEC-017..020. |

## Why REWORK (the decisive findings)

Two dimensions FAILED on genuine P0 contract violations. An independent audit cannot
pass a release on an unsubstantiated, self-contradicting blocker resolution or a broken
audit-trail, regardless of how much else is correct.

### BLOCKERs

1. **Audit events never reach the durable audit log (BLK-G12-001, P0).**
   `services/opportunity|commercial|work|account` write their required `EVT-*` events
   (STAGE-CHANGED, OPPORTUNITY-WON/LOST, QUOTE-ACCEPTED, CONTRACT-SIGNED,
   PAYMENT-RECORDED, OWNER-CHANGED, RECORD-ARCHIVED) only to a local `outbox_events`
   table. **No relay/dispatcher delivers them to audit-history-service.** Those audit
   records therefore never exist in the append-only store. Violates AUD-IMM-002 and the
   audit-log-spec Event Catalog; breaks ACC-014 (record-local history) and ACC-022
   (operation log). (lead/import-export/identity-authz correctly write synchronously —
   so the gap is partial and inconsistent, not total.)

2. **reporting internal write endpoint is unauthenticated (BLK-G12-002).**
   `services/reporting/internal/handler/overview_query.go:61` — `POST /internal/projections`
   performs no `verifyServiceToken`. Any caller reaching the service can inject/overwrite
   report read-models. Violates authz-architecture Service-To-Service rules / SEC-SVC-BLK-002.

3. **Security-group closure is unsubstantiated and self-contradicting (BLK-G12-003, P0, ACC-017).**
   The only genuine Volcengine `DescribeSecurityGroupAttributes` export on disk still lists
   TCP 8088/8443/3389 open to `0.0.0.0/0`. The "post-cleanup" files asserting these were
   removed are hand-authored conclusions with no raw API re-export, and they contradict the
   raw export. BLK-G11-002 "Resolved" is not backed by an artifact and is **reopened**.
   Additionally the CRM instance shares the `Default` security group (all-protocol
   self-referential allow), not a dedicated least-exposure group.

### MAJORs

- **BLK-G12-004** — outbox append is not in the same DB transaction as the business
  mutation and its error is discarded (`_ = h.outbox.Append`); even locally an event row
  can be missing after a successful 201. AUD-IMM-002 "same durable workflow."
- **BLK-G12-005** — reporting's projection ingest path is never invoked; report authz is
  enforced over empty/untrusted data (PM-043..045 effectively unverifiable).
- **BLK-G12-006** — release evidence is not independently verifiable: TLS facts
  (issuer/expiry/SAN/redirect) hand-typed with no captured `openssl s_client`/`curl -vI`/
  `certbot certificates`; restore rehearsal is operator-authored counts with no command
  transcript or checksum-verify output; operator-access conditions 2 & 3 (no routine root
  SSH / key ownership; secrets+backup dir non-world-readable) have zero captured artifacts
  and the deployment was performed as root; no external negative probe shows 8080/5432
  refused from the public internet.
- **BLK-G12-007** — move CRM off the shared `Default` security group to a dedicated
  least-exposure group.

### MINOR

- **BLK-G12-008** — test traceability + hardening: untagged P0 tests
  (TEST-ABUSE-MUTATE-001, TEST-AUTHZ-SCOPE-005, TEST-PERM-CRUD-ADMIN/MGR/SALES-001,
  TEST-OWNER-TRANSFER-002/004); `work` S2S verifier omits the 5-minute max-lifetime cap;
  `AUTHZ_VERSION_STALE` mechanism is dead (mitigated by per-request live user reload);
  stray `截屏2026-06-03 16.32.12.png` at repo root.

## What is genuinely sound (no rework needed)

Per-service DB isolation (own user + schema, no cross-service DB access); real Postgres
persistence with no fakes; audit **storage** is correctly append-only (REVOKE UPDATE/DELETE);
server-side authorization + ownership re-checks; session cookie posture (HttpOnly/Secure/
SameSite + server-side revocation + disabled-user rejection); committed model fidelity;
backup scripts (encrypted, 7-day retention, off-server copy, checksums). The defects are
concentrated in **event-delivery closure** and **evidence verifiability**, not in the
foundation.

## Decision

**Gate NOT Passed — REWORK.** Per the collaboration model, G11 is set back to
`Gate Blocked`; the findings are registered in `planning/blockers.md` as
BLK-G12-001..008 (BLK-G11-002 reopened); a self-contained rework instruction package is
issued at `delivery/G12-rework.md` and handed to Codex. Codex remediates, re-runs the
affected QA/integration evidence, and returns to Claude for a focused G12 re-audit of
dimensions 4 and 5 (plus the pending independent compile/test re-run) before any release
decision. No P0/P1 may be downgraded to clear a finding (no-downgrade rule).
