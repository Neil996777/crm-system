# G12 Rework Package #5 — CRM System (Claude → Codex)

## Document Control

- Source: G12 whole-repo systematic sweep, 2026-06-04
  (`archive/reviews/g12-audit/g12-systematic-sweep-2026-06-04.md`)
- Decision: **REWORK #5** — G11/G12 `Gate Blocked`. Do NOT release.
- These findings were in code never touched by rounds 1–4, so never previously audited.
- Executor: Codex; then return to Claude for a full systematic re-sweep (not a narrow check).

## ⚖️ TWO HARD CONSTRAINTS (read first)

**1. NO DOWNGRADE.** P0/P1 items may only become `Done`, `Blocked`, or `Formal Scope Change by
User`. Do not satisfy any P0/P1 item with mock/stub/TODO/static/non-persistent behavior, and do
not delete or weaken any test to make it pass.

**2. NO OVER-REACH (不得越权).** Fix EXACTLY the listed items, nothing more:
- Do NOT change committed decisions (DEC-001..022), the domain model, the architecture, the
  service boundaries, or scope.
- Fix the IDOR holes by ADDING the missing ownership check using the SAME pattern the sibling
  services already use (`CanReadAccount`/`CanReadLead`/`CanReadOpportunity`). Do not invent a
  new authz scheme and do not change what data the endpoint returns beyond enforcing scope.
- Do NOT refactor/rename/"improve" unrelated code. In particular, the root cause is copy-pasted
  per-service helpers — a shared-library refactor is tempting but is itself a SCOPE CHANGE:
  **do not introduce a shared S2S/error/outbox library.** If you believe one is warranted, STOP
  and raise a blocker for Claude/user — do not do it unilaterally.
- For dead-code deletion, only remove items you have re-verified are zero-call-site; never
  remove anything live.
- If any fix appears to require a model/scope/decision change, STOP and record a blocker in
  `planning/blockers.md` for Claude/user. Deciding it yourself is itself over-reach.

Everything is TDD: add a failing test that reproduces the gap, then fix to green.

---

## 🔴 BLK-G12-015 (BLOCKER, P0) — commercial single-record read IDOR

`services/commercial/internal/handler/contract_query.go:25` (`getContract`) and
`quote_query.go:24` (`getQuote`) return a record fetched by id with no ownership check.

**Fix:** add an ownership/CanRead check before returning, matching the commercial mutation-side
rule already present (`contract_status.go:47`: Sales may act only on owner_id == actor.ID) and
the by-id read pattern in account/lead/opportunity. Administrators/Managers per their committed
scope; Sales only own/assigned. On deny, return the committed safe denial (see BLK-G12-024 for
the existence-leak contract — apply the no-existence-leak form).

**Acceptance (TDD):** tests prove a Sales user gets denied (no record data leaked) on a
non-owned contract and quote read by id, and an owner/Manager/Admin succeeds; no data exposed
on denial.

## 🔴 BLK-G12-016 (BLOCKER, P0) — import-export getImportRun IDOR

`services/import-export/internal/handler/import_run.go:169` returns an import run by id with no
actor/ownership check (exposes filename + per-row PII).

**Fix:** enforce actor scope (the run stores ActorID/ActorRole/TeamID) consistent with the
write-side restriction and the committed permission model. **Acceptance:** non-owner/Sales
denied with no leak; owner/Manager/Admin per scope succeeds.

---

## 🟠 BLK-G12-017 (MAJOR, P0 audit) — identity-authz durable/transactional audit

`services/identity-authz`: audit/op-log for login, role/status change, and last-admin-block is
emitted post-commit, non-transactionally, with the error discarded (`user_admin.go:211`); and
`identity_authz.outbox_events` is written but never dispatched (dead table).

**Fix (choose the minimal consistent option, no new shared lib):** make the audit emission for
these security-critical events durable in the same workflow as the mutation — either (a) add the
same transactional outbox + dispatcher the other services use (`inTransaction`/`NewOutboxTx` +
a `DispatchOnce` loop delivering to audit-history), or (b) if the synchronous audit client is
kept, make it part of the transaction boundary and stop discarding its error. Remove the dead
orphaned outbox path if you do not wire it. Do not lose the security events.

**Acceptance (TDD):** a fail-first test proves a failed audit delivery does not silently succeed
(mutation+audit are atomic or the event is durably retried), for a role/status change.

## 🟠 BLK-G12-018 (MAJOR, P0 / ACC-012) — work changeTaskStatus optimistic concurrency

`services/work/internal/handler/work_command.go:240` doesn't read `expectedVersion` and funnels
`ErrVersionConflict` into `writeOutboxFailure` instead of `409 VERSION_CONFLICT`.

**Fix:** adopt the same optimistic-concurrency pattern as the other write services
(read `expectedVersion`, reject 0, return `VERSION_CONFLICT` on stale). **Acceptance:** a
stale-version task-status change returns `VERSION_CONFLICT` and does not lose the concurrent edit.

## 🟠 BLK-G12-019 (MAJOR) — lead conversion downstream idempotency

`services/lead/internal/client/conversion_client.go` + the account/opportunity create handlers:
a retry after partial success re-creates the account (orphaned duplicate).

**Fix:** make the downstream creates idempotent (honor an idempotency key derived from the lead
conversion) so a retry is safe. **Acceptance:** a test simulates account-created-then-
opportunity-failed, retries, and asserts no duplicate account.

## 🟠 BLK-G12-020 (MAJOR) — add X-Correlation-Id to the 7 synchronous internal clients

Sites: `opportunity/internal/authz/commercial_client.go:44`,
`lead/internal/client/conversion_client.go:114`, `account/internal/handler/archive.go:148`,
`work/internal/handler/work_query.go:112`, `import-export/internal/client/audit_client.go:92,145`,
`lead/internal/client/audit_client.go:84`, `identity-authz/internal/authz/audit_client.go:85`.

**Fix:** set `X-Correlation-Id` on each (same as the async outbox callers). **Acceptance:** a
test asserts the header is present on each client call.

## 🟠 BLK-G12-021 (MAJOR) — emit the 8 missing EVT catalog events under their catalog ids

Map to the distinct ids in `docs/security/audit-log-spec.md` (do NOT invent semantics):
EVT-AUTH-LOGIN-SUCCEEDED, EVT-AUTH-LOGIN-FAILED, EVT-AUTH-ACCESS-DENIED, EVT-USER-ROLE-CHANGED,
EVT-USER-STATUS-CHANGED, EVT-LAST-ADMIN-BLOCKED, EVT-CONTRACT-TERMINATED, EVT-REPORT-ACCESS-DENIED.
For EVT-STATUS-CHANGED / EVT-TASK-COMPLETED / EVT-TASK-CANCELLED / EVT-PAYMENT-OVERDUE: these may
be intentional by-design collapses — reconcile with the spec; if the spec and PSM already say
on-read/generic, leave them and note it, do NOT force-emit (that could be over-reach).

**Acceptance:** each of the 8 has an emission path and a positive creation test asserting the
catalog `EVT-*` id.

## 🟠 BLK-G12-022 (MAJOR) — account + commercial outbox→reporting delivery tests

Add the reporting-target dispatcher test (set `ReportingServiceURL`, assert delivery with S2S
headers + X-Correlation-Id + failure-retry) to `services/account/internal/event/dispatcher_test.go`
and `services/commercial/internal/event/dispatcher_test.go`, matching the lead test.

## 🟠 BLK-G12-023 (MAJOR, doc) — api-spec Close-Won schema

`docs/architecture/api-spec.md:378-392`: add the required `contractId` to the Close-Won request
schema (and remove/clarify the unused `idempotencyKey` on close-won/close-lost). Doc-only; the
code is correct — do not change the handler.

---

## 🟡 BLK-G12-024 (MINOR) — error/denial uniformity batch

- Reconcile the denial contract so an existing-but-unreadable single-record read does NOT leak
  existence (committed Denial Contract: must not reveal whether an unauthorized record exists).
  Make opportunity/lead/account record reads consistent (no 403-vs-404 split).
- identity-authz: give auth sign-in/out failures a stable `code`/`category` envelope like every
  other service (`auth.go:253`).
- audit-history verifier: fail-closed on empty audience/intent (`service_token.go:38`).
- identity-authz verifier: use `time.Now().UTC()` (`service_token.go:64`).
- `/internal/sessions/check` cookie-only (`auth.go:67`): confirm intentional and document it, or
  gate it — your call, but document the decision; do not silently leave the divergence.

## 🟡 BLK-G12-025 (LOW) — cleanup batch

- Delete verified dead code: lead `internal/client/audit_client.go`; zero-call helpers
  (`IsServiceAuthFailed`, `HasSurface`, `IsForeignKeyError`, `ProjectionTableName`,
  `ExportScope`, `EventEnvelope`). For the unconsumed `shared/contracts` module: do NOT delete or
  refactor it unilaterally (it may be intended published vocabulary) — raise it for a decision.
- opportunity dispatcher test: add reporting failure-path + UID-dedup assertions.
- Fix the 2 test-model id citations (TEST-NAV-RETRIEVE-006, TEST-TASK-LIFECYCLE-004).
- Add `.secrets` to the project `.gitignore` and confirm no secret file is tracked.

---

## Definition of Done

- BLK-G12-015..023 fixed with the named fail-first tests green; BLK-G12-024/025 addressed.
- No new fakes; no weakened/deleted tests; no scope/model change; no unilateral shared-lib
  refactor; any item needing a decision raised as a blocker instead of decided.
- `delivery/tasks.md`, `modeling/traceability-matrix.md`, `planning/blockers.md` updated to real
  artifacts; commits made.
- Return to Claude for a FULL whole-repo systematic re-sweep. Do not self-pass G12.
