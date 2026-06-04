# G12 Whole-Repo Systematic Sweep — Findings [KICKBACK #5]

## Document Control

- Project: CRM System
- Gate: G12 — Audit → Release/Rework
- Owner: Audit (Claude), independent of execution
- Date: 2026-06-04
- Method: 7 parallel author≠auditor whole-repo sweeps by cross-cutting dimension
  (transactional integrity; dead code; dispatcher test coverage; cross-service uniformity;
  test integrity & EVT/permission coverage; security/authz; model fidelity & traceability),
  NOT scoped to recent changes. One dimension (uniformity) was re-run after a crash.
- Decision: **REWORK #5 — Gate NOT Passed.** No release.

## Why this sweep was run

The prior five rounds each audited only "what Codex just changed," so code never touched by
a fix was never audited. A full-repo systematic sweep was run to drain the tail. It found
release-blocking issues in areas no prior round inspected — vindicating the sweep.

## Findings (deduplicated, by severity)

### BLOCKER — broken object-level authorization (IDOR), real data exposure

1. **commercial `getContract` / `getQuote` IDOR.** `contract_query.go:25`, `quote_query.go:24`
   read `repo.Find` (`WHERE id=$1`, no owner scope) and return with **no CanRead/ownership
   check**; the commercial domain has no `CanRead` at all. Any authenticated Sales user can
   read any contract/quote (amounts, customer, owner) by id enumeration. List endpoints ARE
   role-scoped; only the by-id read is open. Violates PM-020/021. Sibling services
   (account/lead/opportunity) all gate by-id reads with `CanRead*`.
2. **import-export `getImportRun` IDOR.** `import_run.go:169` — no actor/ownership check; runs
   store ActorID/Role/TeamID and the DTO exposes filename + per-row error detail (lead PII).
   Any authenticated user can read another user's import run by id. Write side restricts
   non-Sales; the read does not.

### MAJOR

3. **identity-authz audit is not durable/transactional.** No `inTransaction`/`NewOutboxTx`
   (`internal/event/outbox.go:28`); user mutations commit, then audit/op-log is emitted
   post-commit, non-transactionally, with the HTTP error discarded (`user_admin.go:211`
   `log.Printf`). Also `identity_authz.outbox_events` is written but **never dispatched** (no
   `DispatchOnce` anywhere) — a dead orphaned table. The most security-critical audit
   (login, role/status change, last-admin-block) is the least reliable. AUD-IMM-002, ACC-022.
4. **work `changeTaskStatus` has no optimistic concurrency.** `work_command.go:240` (TASK-024,
   **ACC-012 P0**) doesn't read `expectedVersion`; the repo's `ErrVersionConflict` is funneled
   into `writeOutboxFailure` instead of `409 VERSION_CONFLICT`. Stale concurrent edits are
   silently lost. work is the only write service emitting no `VERSION_CONFLICT`.
5. **lead conversion saga lacks downstream idempotency.** `conversion_client.go` calls account
   `POST /internal/accounts` and opportunity `POST /internal/opportunities`; neither create
   honors an idempotency key, and the lead-side key is persisted only after both succeed
   (`lead_convert.go:71`). A retry after partial success (account created, opportunity failed)
   re-creates the account → orphaned duplicate.
6. **7 synchronous internal clients omit `X-Correlation-Id`** (trace break; STB-003):
   `opportunity/internal/authz/commercial_client.go:44`, `lead/internal/client/conversion_client.go:114`,
   `account/internal/handler/archive.go:148`, `work/internal/handler/work_query.go:112`,
   `import-export/internal/client/audit_client.go:92,145`, `lead/internal/client/audit_client.go:84`,
   `identity-authz/internal/authz/audit_client.go:85`. (Async outbox callers all comply.)
7. **EVT catalog: 8 of 24 events never emitted under their catalog id, untested.** All 6
   identity/auth events (EVT-AUTH-LOGIN-SUCCEEDED/FAILED, EVT-AUTH-ACCESS-DENIED,
   EVT-USER-ROLE-CHANGED, EVT-USER-STATUS-CHANGED, EVT-LAST-ADMIN-BLOCKED) collapse to raw
   strings / one `EVT-USER-ADMIN-CHANGED`; `EVT-CONTRACT-TERMINATED` folds into the generic
   status event; `EVT-REPORT-ACCESS-DENIED` is never emitted. Violates audit-log-spec
   Testability ("every required event has a positive creation test"). (EVT-STATUS-CHANGED /
   TASK-COMPLETED / TASK-CANCELLED / PAYMENT-OVERDUE are arguable by-design collapses —
   reconcile with the spec, don't blindly emit.)
8. **account + commercial outbox→reporting delivery untested** (BLK-G12-014 class). Both
   deliver to reporting in prod but their dispatcher tests never set `ReportingServiceURL`, so
   reporting delivery / S2S / correlation-id / failure-retry are untested.
9. **api-spec Close-Won request schema omits required `contractId`** (`api-spec.md:378-392`),
   so a spec-conformant client cannot satisfy DEC-017. Doc fix (code is correct).
10. **identity-authz codeless error envelope.** A 2-arg `writeError` (`auth.go:253`) emits no
    `code`/`category`; the gateway can't map it to a stable machine code like every other
    service. Internally inconsistent (it also has `writeErrorCode`).

### MINOR / LOW (batch cleanup)

- **Denial existence leak / 403-vs-404 inconsistency:** opportunity/lead/account-getAccount
  return 403 PERMISSION_DENIED for an existing-but-unreadable record (leaks existence), while
  account-getContact returns 404. Reconcile to the committed Denial Contract (must not reveal
  whether an unauthorized record exists → no existence leak).
- audit-history verifier doesn't fail-closed on empty audience/intent (`service_token.go:38`).
- identity-authz `/internal/sessions/check` is cookie-gated, not S2S (`auth.go:67`) — confirm
  intentional + document, or gate (do not force without deciding).
- identity-authz verifier uses local `time.Now()` not `.UTC()` (`service_token.go:64`).
- Dead code: lead `internal/client/audit_client.go` (re-wireable non-transactional path);
  the entire `shared/contracts` module is imported by no service; zero-call-site helpers
  (`IsServiceAuthFailed`, `HasSurface`, `IsForeignKeyError`, `ProjectionTableName`,
  `ExportScope`, `EventEnvelope`).
- opportunity dispatcher test depth: no reporting failure-path, no UID-dedup assertion.
- 2 test-model id mis-citations (TEST-NAV-RETRIEVE-006, TEST-TASK-LIFECYCLE-004 — behavior
  covered, id not cited).
- `.secrets` not in the project `.gitignore` (confirm no secret is tracked).

## What is clean

Model fidelity (DEC-017..020) is faithful in code+migrations; stage enum exact; one quote per
opportunity DB-enforced; Won via real S2S; payment decoupled; traceability honest (13/13
spot-checked); retired tests absent; weak-test hygiene clean (0 skips/empty); append-only audit
DB-enforced; gateway-only external exposure; no cross-service DB access; session/cookie
hardening present.

## Root cause

`service_token.go`, the `writeError`/error helpers, and the outbox/dispatcher code are
**copy-pasted per service with no shared library**. That duplication is how the work 5-min-cap,
work version-conflict, identity-authz outbox/error, and correlation-id drifts went unnoticed.
NOTE: a shared-library refactor is itself a scope change — it must NOT be undertaken
unilaterally by execution; if warranted, raise it for a decision.

## Decision

**REWORK #5.** Findings registered as BLK-G12-015..025 in `planning/blockers.md`; rework
package `delivery/G12-rework-5.md` (with explicit no-downgrade AND no-over-reach constraints).
After Codex remediates, Claude re-runs the SAME whole-repo systematic sweep (not a narrow
scope) to verify the tail is genuinely drained before any G12 pass.
