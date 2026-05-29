# Security MDA Pre-Task Review

## Document Control

- Project: CRM System
- Review Gate: G7/G8 Pre-Task Security Review
- Reviewer Agent: Security Compliance
- Date: 2026-05-27
- Status: Passed
- Output Location: `archive/reviews/g7-modeling/security-mda-pre-task-review.md`

Implementation remains blocked until G8 passes.

## Review Scope

This review evaluates whether the drafted MDA Modeling package can be safely
handed to G8 Task Planning without causing task planners, engineers, QA, or
audit to guess security behavior.

Reviewed inputs:

- `modeling/CIM.md`
- `modeling/PIM.md`
- `modeling/PSM.md`
- `modeling/domain-model.md`
- `modeling/state-machines.md`
- `modeling/domain-events.md`
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`
- `docs/security/security-requirements.md`
- `docs/security/permission-matrix.md`
- `docs/security/privacy-requirements.md`
- `docs/security/audit-log-spec.md`
- `docs/security/abuse-cases.md`
- `docs/security/compliance-risks.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/product/acceptance-matrix.md`

Rules applied:

- P0/P1 items cannot be downgraded, removed, merged away, weakened, or accepted
  as partial work.
- Frontend hiding never satisfies authorization.
- Any open P0/P1 security blocker would block G8 Task Planning.
- This review does not edit `modeling/` and does not create implementation
  code.

## Gate Decision

Decision: Passed.

Security Compliance finds no P0/P1 security blocker preventing the project from
entering G8 Task Planning, provided the task plan explicitly carries the model,
security, privacy, audit, abuse-case, and test obligations listed below.

## P0/P1 Blockers

None.

## Security Coverage Assessment

| Area | Result | Evidence |
|---|---|---|
| Three-role authorization | Passed | `CIM-002` to `CIM-004`, `PIM-POLICY-AUTHZ`, `DM-AUTH-001` to `DM-AUTH-004`, `PSM-AUTHZ`, `PSM-API-*`, `TM-002`, and `TP-AUTH-001` to `TP-AUTH-007` preserve Administrator, Sales Manager, and Sales behavior. |
| Owned/assigned/team/governed scope | Passed | `PIM-VO-SCOPE`, `PIM-POLICY-AUTHZ`, `DM-AUTH-*`, `DM-REP-001`, `PSM-API-*`, and test patterns require scope checks for protected reads, mutations, lists, reports, reminders, exports, and logs. |
| Archived context | Passed | `CIM-017`, `PIM-VO-ARCHIVE`, `PIM-SVC-ARCHIVE`, `SM-ARCHIVE`, `PSM-TX-009`, and `TM-015` model archived records as explicit authorized context, not deletion. |
| Last Administrator | Passed | `PIM-POLICY-LASTADMIN`, `SM-USER`, `EVT-LAST-ADMIN-BLOCKED`, `PSM-IDX-001`, `PSM-TX-002`, and `INV-002` preserve the server-side guard and concurrent protection expectation. |
| Stale session / role downgrade | Passed | `PIM-SESSION`, `PSM-AUTH`, `PSM-DB-002`, `ABT-008`, and `TP-AUTH-001` require active-user and role-version recheck on protected APIs. |
| Safe denial | Passed | `CIM-INV-001`, `CIM-INV-002`, `PIM-QRY-003`, `PSM-UI-004`, common error DTOs, and `ABT-001` / `ABT-002` require no restricted data exposure and no mutation on denied actions. |
| Record-local history | Passed | `CIM-015`, `PIM-HISTORY`, `DM-014`, `PSM-AUDIT`, `PSM-DB-014`, `EVT-*`, `TM-014`, and `INV-009` require append-only, scoped, transaction-linked record evidence. |
| Global operation log | Passed | `CIM-016`, `PIM-OPLOG`, `DM-015`, `PSM-DB-015`, `PSM-API-012`, `EVT-*`, `TM-022`, and `ABT-007` preserve Administrator-only global operation log behavior. |
| Duplicate warning privacy | Passed | `CIM-018`, `PIM-SVC-DUP`, `PIM-VO-DUPKEY`, `DM-RULE-018`, `PSM-DUP`, `PSM-DTO-009`, `PSM-TX-004`, `TM-019`, and `ABT-003` prevent merge/overwrite and unauthorized match-detail leakage. |
| Import/export | Passed | `PIM-IMPORT`, `PIM-EXPORT`, `SM-IMPORT`, `SM-EXPORT`, `PSM-IMPORTEXPORT`, `PSM-TX-010`, `PSM-TX-011`, `TM-020`, `ABT-004`, and `ABT-006` cover authorization, row validation, safe errors, dangerous cell handling, and operation logs. |
| Reports and manager overview | Passed | `PIM-SVC-REPORT`, `PIM-QRY-007`, `PIM-QRY-008`, `PSM-REPORT`, `PSM-API-016`, `TM-018`, `TM-023`, and `ABT-005` require authorization before aggregation and Sales denial. |
| Reminders | Passed | `PIM-SVC-REMINDER`, `PIM-QRY-006`, `SM-CONTRACT`, `SM-PAYMENT`, `SM-TASK`, `PSM-REMINDER`, `TM-021`, and `ABT-021` require authorized due/overdue visibility only. |
| Abuse-case traceability | Passed | `test-model.md` maps IDOR, direct mutation bypass, duplicate probing, export/report leakage, dangerous CSV cells, log access, and stale role/session concepts to acceptance IDs. |

## P2 Improvements

These are not blockers for G8 entry, but Task Planning should carry them as
explicit task/test coverage to prevent ambiguity during implementation.

| ID | Priority | Improvement | Rationale | Suggested Owner |
|---|---|---|---|---|
| SEC-MDA-P2-001 | P2 | In G8, map tasks not only to `ACC-*` and `TM-*`, but also to `SEC-*`, `PM-*`, `PRIV-*`, `AUD-*`, and `ABUSE-*` IDs for security-sensitive work. | The MDA traceability is sufficient for G7, but implementation task closure will be easier to audit if security source IDs are explicit in the task map. | Task Planner, QA TDD |
| SEC-MDA-P2-002 | P2 | Add explicit implementation tasks for database-level append-only protections for `record_history_events` and `operation_log_events`. | MDA requires append-only behavior through normal workflows; task planning should translate this into migrations, DB role/privilege restrictions, repository API shape, and negative tests. | Task Planner, Backend, QA TDD |
| SEC-MDA-P2-003 | P2 | Add explicit task/test coverage for import raw-file lifecycle and export short-lived delivery cleanup if file storage is implemented. | The privacy requirements define raw import and generated export retention limits. G8 should ensure upload/download mechanics do not accidentally retain restricted files longer than allowed. | Task Planner, Backend, Security Compliance |
| SEC-MDA-P2-004 | P2 | Add abuse-test tasks for duplicate-warning rate/probing controls if repeated probing detection is implemented in v1. | MDA marks suspicious duplicate probing event emission as conditional. If implemented, it needs concrete thresholds, safe payloads, and tests. | Task Planner, QA TDD, Security Compliance |
| SEC-MDA-P2-005 | P2 | Add contract tests for safe error/warning DTOs and generated frontend client types. | The safe denial and duplicate-warning contracts are central security boundaries; contract tests reduce drift between OpenAPI, Go handlers, and React handling. | Task Planner, Architecture, QA TDD |

## Required G8 Carry-Forward Conditions

Task Planning may proceed only if the plan keeps these obligations intact:

- Every P0/P1 `ACC-*` must map to implementation, unit/integration/E2E or
  manual verification, integration evidence, and later audit evidence.
- Every protected endpoint and repository query must have backend authn/authz
  tasks; UI hiding alone is insufficient.
- Every denied mutation must have tests proving no business-state mutation and
  no restricted data exposure.
- Every sensitive successful mutation must have transaction-scoped
  history/log-writing tasks and tests.
- Import/export/report/reminder tasks must enforce authorization before query,
  aggregation, output generation, or displayed reminder rows.
- Duplicate-warning tasks must preserve safe warning payloads, token binding,
  no automatic merge/overwrite/link/ownership-transfer, and unauthorized match
  masking.
- Last Administrator, stale session/role recheck, archived context, IDOR,
  global operation-log denial, and record-local history scope must each be
  represented in QA and abuse-test tasks.
- No task may satisfy P0/P1 behavior with mock, stub, TODO, static-only,
  in-memory-only, or non-persistent core paths.

## Recommendation

Security Compliance recommends entering G8 Task Planning after the remaining
G7 pre-task reviews complete. This recommendation does not approve
implementation by itself; implementation remains blocked until G8 passes.

## Files Modified

- `archive/reviews/g7-modeling/security-mda-pre-task-review.md`
