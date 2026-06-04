# G12 Rework Package #6 — CRM System (Claude → Codex, consolidated final tail)

## Document Control

- Source: G12 full whole-repo systematic RE-SWEEP, 2026-06-04
  (`archive/reviews/g12-audit/g12-resweep-2026-06-04.md`)
- Decision: **REWORK #6** — G11/G12 `Gate Blocked`. Do NOT release.
- Status going in: BLK-G12-015..025 all Resolved/verified; no over-reach; no downgrade;
  build/vet/frontend green. These five (BLK-G12-026..030) are the **entire** remaining tail —
  one MAJOR + three MINOR + one LOW-batch. This is the consolidated final pass.
- Executor: Codex; then return to Claude for a minimal spot re-audit of exactly these five items
  + a build/test confirm → **G12 PASS** if clean. Do NOT self-pass G12.

## ⚖️ TWO HARD CONSTRAINTS (unchanged — read first)

**1. NO DOWNGRADE.** P0/P1 items may only become `Done`, `Blocked`, or `Formal Scope Change by
User`. No mock/stub/TODO/static/non-persistent to satisfy any item; do not delete or weaken any
test. The 3 retired tests stay absent.

**2. NO OVER-REACH (不得越权).** Fix EXACTLY these five, nothing more:
- No committed-decision (DEC-001..022), domain-model, architecture, service-boundary, or scope
  change.
- **Do NOT introduce a shared S2S/error/outbox library** and do NOT refactor the copy-pasted
  per-service helpers. The work-verifier divergence in BLK-G12-030 is to be aligned IN PLACE
  within the work service only — not by extracting a shared package. If you believe a shared lib
  is warranted, STOP and raise a blocker; do not do it unilaterally.
- Do NOT delete `shared/contracts` (already raised as a separate decision).
- If any item appears to need a model/scope/decision change, STOP and record a blocker for
  Claude/user instead of deciding it.

Everything is TDD: add a failing test that reproduces the gap, then fix to green.

---

## 🟠 BLK-G12-026 (MAJOR, P0 audit) — make identity-authz denial-path security audit durable

**Where:** `services/identity-authz/internal/handler/auth.go:244-252` (`appendAccessDenied`,
covering EVT-AUTH-LOGIN-FAILED and the unauthenticated / invalid-session / inactive-user /
authz-version-stale / user-admin-denied denials) and `permission.go:55-65`
(`UserAccessDenied` → EVT-AUTH-ACCESS-DENIED).

**What's wrong:** these write via the non-transactional `h.outbox.Append` and **swallow the
error with `log.Printf`** (`auth.go:250`, `permission.go:64`). BLK-G12-017 made the *mutation*
(change) audit events durable but left these *denial* events best-effort, so a failed append on a
failed-login or access-denied event is silently lost — and these are exactly the events most
relevant to intrusion / brute-force detection.

**Required fix:** make these denial-event appends durable and non-discarded, using the SAME
mechanism BLK-G12-017 already established in this service — append the outbox row inside a DB
transaction (or, where there is no surrounding business mutation to bind to, write the outbox row
in its own committed transaction and **propagate** the append error rather than swallowing it).
Do NOT add a new shared library; reuse the identity-authz dispatcher/outbox already wired in
`internal/event`. Where a denial path genuinely cannot fail the user-facing response (e.g. an
auth check that must still return 401/403 to the client), the outbox write must still be durable
and retryable (row persisted, dispatcher delivers), and any persistence failure must be surfaced
in logs AND not silently dropped from the audit trail (no `log.Printf`-and-continue that loses
the event).

**Acceptance (TDD):** a fail-first test proves that when the outbox append for a failed-login /
access-denied event fails, the event is NOT silently lost (the write is retried/durable or the
operation surfaces the failure) — mirroring `TestUserAdminAuditOutboxFailureRollsBackMutation` /
`TestOutboxDispatcherDeliversAuditHistoryAndRetries`. No `log.Printf`-swallowed audit append
remains on any identity-authz denial path.

---

## 🟡 BLK-G12-027 (MINOR) — remove dead identity-authz audit client carrying a stale collapsed id

**Where:** `services/identity-authz/internal/authz/audit_client.go` (whole file:
`NewAuditClient`/`AuditClient`/`AppendOperationLog`, incl. the hard-coded `EVT-USER-ADMIN-CHANGED`
at ~line 53).

**What's wrong:** after BLK-G12-017 moved identity audit onto the dispatcher, this file has zero
call sites, and it still hard-codes the collapsed id `EVT-USER-ADMIN-CHANGED` that BLK-G12-021
replaced with the distinct catalog ids. Dead, and a latent re-collapse risk.

**Required fix:** re-verify zero call sites (`rg`), then delete the file (same evidence discipline
as BLK-G12-025). If anything still references it, STOP — that means 017's migration is incomplete;
raise it. Update `scripts/test_cleanup_contract.py` (or equivalent) to assert the file is gone.

**Acceptance:** file removed; `go build ./...` + `go test ./... -count=1` green in
`services/identity-authz`; cleanup guard asserts absence; `rg 'EVT-USER-ADMIN-CHANGED'` returns no
production matches.

---

## 🟡 BLK-G12-028 (MINOR) — pin the catalog eventId for the 8 untested-mapping events

**Where:** the emission/dispatch tests for EVT-LEAD-CONVERTED, EVT-QUOTE-ACCEPTED,
EVT-PAYMENT-RECORDED, EVT-OPPORTUNITY-WON, EVT-OPPORTUNITY-LOST, EVT-RECORD-ARCHIVED,
EVT-IMPORT-RUN, EVT-EXPORT-RUN.

**What's wrong:** each has an emission path and a handler/outbox test, but no test asserts the
**derived catalog `eventId`** produced by `auditEventContract` (contrast EVT-STAGE-CHANGED, pinned
at `opportunity/internal/event/dispatcher_test.go:135`). A mapping regression would go undetected;
audit-log-spec Testability is met only loosely.

**Required fix (test-only):** add assertions (extend the existing tests; do not weaken any) that
pin each of the 8 to its catalog `EVT-*` id at the point the event is mapped/dispatched. No
production-code change expected — if you find a mapping is actually wrong, that is a separate
finding; raise it rather than silently changing semantics.

**Acceptance:** each of the 8 has a test asserting its catalog id; `go test ./... -count=1` green
in the affected services.

---

## 🟡 BLK-G12-029 (MINOR) — close the contacts sub-resource existence leak

**Where:** `services/account/internal/handler/account_command.go` `listContactsForAccount`
(returns `403 PERMISSION_DENIED` for an existing-but-unreadable account).

**What's wrong:** `getAccount` for the same id returns `404 NOT_FOUND` (per BLK-G12-024's Denial
Contract), but the contacts sub-resource returns 403 — leaking which account ids exist via the
sub-resource. Same existence-leak class BLK-G12-024 closed for the primary reads.

**Required fix:** apply the committed Denial Contract here too — a non-owned/unreadable account's
contacts listing returns the same safe `404 NOT_FOUND` with no data, consistent with `getAccount`.
Do not change what an authorized caller receives.

**Acceptance (TDD):** a test proves a non-owner Sales user listing contacts for an existing
non-owned account gets `404 NOT_FOUND` (no existence leak, no contact data); owner/Manager/Admin
unchanged.

---

## 🔵 BLK-G12-030 (LOW) — consistency/robustness batch (in-place, no shared lib)

All three are in-place fixes; **do not** extract a shared library for any of them.

- **(a) work S2S verifier shape.** `services/work/internal/authz/service_token.go` uses
  `ServiceClaims` + `errors.New("invalid claims")` instead of the `ServiceTokenClaims` /
  `ErrServiceAuthFailed` shape the other six verifiers use. Align the type/error names **within
  the work service** so the contract matches (behavior is already equivalent — this is naming/shape
  consistency only). If aligning would require touching shared code, STOP and raise it instead.
- **(b) concurrent lead-conversion double-submit mapping.** A racing second create that hits the
  unique-index violation (after the pre-check) currently returns generic `400 VALIDATION_FAILED`
  instead of the idempotent `200`-with-existing-row. No duplicate is created (the index holds);
  fix the error mapping so the unique-violation on the lead-conversion idempotency key returns the
  existing record (200), matching the sequential-retry path. Sites:
  `account/internal/handler/account_command.go` and the opportunity equivalent.
- **(c) account duplicate-warning two-transaction window.** `account/internal/handler/duplicate_check.go:43`
  persists the duplicate token and appends the `DuplicateWarningRaised` outbox row in two separate
  transactions. The error is already propagated (not swallowed), and this is not a P0/P1 mutation,
  so this is the lowest priority — if it can be made single-transaction without disturbing the
  duplicate-check flow, do so; otherwise document why it is acceptably two-phase and leave it. Your
  call, but record the decision.

**Acceptance:** (a) work verifier uses the aligned type/error names, tests green; (b) a test proves
a concurrent/duplicate-key lead-conversion create returns 200-with-existing-row, no duplicate;
(c) either single-transaction with a test, or a recorded rationale.

---

## Definition of Done

- BLK-G12-026 fixed with the named fail-first test green; BLK-G12-027/028/029 fixed with named
  tests; BLK-G12-030 (a)/(b) fixed with tests, (c) fixed-or-documented.
- No new fakes; no weakened/deleted tests; no scope/model/decision change; **no shared-lib
  refactor**; `shared/contracts` untouched; any item needing a decision raised as a blocker.
- `delivery/tasks.md`, `modeling/traceability-matrix.md`, `planning/blockers.md` updated to real
  artifacts; commits made.
- Return to Claude for a minimal spot re-audit of exactly these five + a build/test confirm. If
  clean, **Claude passes G12.** Do NOT self-pass G12.
