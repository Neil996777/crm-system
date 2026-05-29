# Security Compliance G8 Task Planning Review

## Decision

Passed from Security Compliance review perspective, with P2 improvements.

No P0/P1 blockers were found. The moved `delivery/` artifacts preserve the
P0/P1 security requirements for backend authorization, scope loaders,
authorization-before-query, safe denial, audit logging, privacy masking,
import/export controls, duplicate-warning privacy, report/reminder scoping, and
no mock/static/non-persistent core paths.

Implementation remains blocked until the full G8 gate is approved by the Task
Planner owner and all required reviewers.

## Reviewed Inputs

| Input | Review Purpose |
|---|---|
| `../../company/operating-model.md` | G8 reviewer role, no-downgrade rule, blocker behavior, and P0/P1 release impact. |
| `../../standards/acceptance-matrix-standard.md` | Required acceptance traceability, evidence, and no-downgrade controls. |
| `../../standards/status-and-priority-standard.md` | Priority, status, and blocker severity language. |
| `PROJECT_CONTEXT.md` | Current G8 draft status, active inputs, and implementation block before G8 pass. |
| `delivery/tasks.md` | End-to-end task definitions, production/test file surfaces, security checks, manual verification, TDD guards, and blocker rules. |
| `delivery/acceptance-task-map.md` | ACC to TASK/TM/evidence mapping for P0/P1 acceptance. |
| `delivery/task-dependencies.md` | Task sequencing and dependency risk. |
| `delivery/delivery-plan.md` | End-to-end capability sequencing and integration expectations. |
| `delivery/blockers.md` | Open blocker register and security watch items. |
| `docs/security/security-requirements.md` | Authentication, authorization, audit, import/export, reports, privacy, and verification requirements. |
| `docs/security/permission-matrix.md` | Role/action/resource/scope rules and required permission test patterns. |
| `docs/security/privacy-requirements.md` | Data classification, masking, retention, import/export, report, reminder, and duplicate warning privacy rules. |
| `docs/security/audit-log-spec.md` | Record-local history, global operation log, event schema, immutability, and testability requirements. |
| `docs/security/abuse-cases.md` | Abuse scenarios for IDOR, unauthorized mutation, stale sessions, import/export, duplicate enumeration, reports, reminders, and audit leakage. |
| `docs/security/compliance-risks.md` | Security/compliance risks that must remain visible through planning and later audit. |
| `modeling/PSM.md` | Authz policy mapping, resource scope loaders, safe denial, stale-session behavior, audit/log storage, import/export, reports, reminders, retention, and architecture acceptance. |
| `modeling/test-model.md` | TM, TP-AUTH, INV, and ABT coverage for security-sensitive P0/P1 acceptance. |
| `modeling/traceability-matrix.md` | ACC to CIM/PIM/PSM/TM/TASK reverse traceability. |

## Findings

| ID | Severity | Finding | Evidence | Required Action |
|---|---|---|---|---|
| SEC-G8-001 | Pass | Backend authorization and scope-loader requirements are covered as executable task work, not frontend-only permission hiding. | `delivery/tasks.md` TASK-002 adds `apps/api/internal/authorization/policies.go`, `scope_loaders.go`, `scope_queries.go`, authz middleware, policy tests, IDOR tests, stale-role tests, and requires every protected backend path to call central policy/scope loaders. `modeling/PSM.md` maps `PSM-AUTHZ-001` to `PSM-AUTHZ-014` and `PSM-SCOPE-001` to `PSM-SCOPE-014`. | None for G8 pass. |
| SEC-G8-002 | Pass | Safe denial and privacy masking are preserved for unauthorized records, duplicate warnings, import row errors, reports, reminders, and logs. | `delivery/tasks.md` TASK-002, TASK-015, TASK-019, TASK-020, TASK-021, TASK-022, and TASK-023 require safe denied/missing behavior, hidden unauthorized rows, masked duplicate details, safe row errors, authorization-before-output, and safe before/after values. `docs/security/privacy-requirements.md` and `docs/security/security-requirements.md` define the same masking rules. | None for G8 pass. |
| SEC-G8-003 | Pass | Record-local history and Administrator global operation logs are planned with append-only, scoped, transaction-aware behavior. | TASK-014 requires transaction-linked history writes, scoped reads, no edit/delete route, and rollback behavior where testable. TASK-022 requires Administrator-only global log query, event writes for sensitive operations, immutability tests, safe before/after values, and restore-preserved log evidence. `docs/security/audit-log-spec.md` and `modeling/PSM.md` align on event schema and storage expectations. | None for G8 pass. |
| SEC-G8-004 | Pass | Import/export, duplicate warning, report, and reminder abuse cases are covered with authorization-before-query/output tests. | TASK-019 includes unauthorized duplicate masking and `ABT-003`. TASK-020 includes scoped import/export, Sales denial, dangerous cell handling, retained run metadata, `ABT-004`, and `ABT-006`. TASK-021 requires unauthorized reminders excluded before output. TASK-018 and TASK-023 require unauthorized rows excluded before aggregation with `ABT-005`. `modeling/test-model.md` includes `TP-AUTH-007`, `INV-012`, and the relevant ABT rows. | None for G8 pass. |
| SEC-G8-005 | Pass | No P0/P1 security requirement is downgraded, deleted, merged away, or accepted as mock/static/non-persistent behavior. | `delivery/tasks.md` global planning rules and per-task no-downgrade rows reject mock, stub, TODO, static-only, in-memory-only, localStorage-only, frontend-only authorization, client-side filtering of unscoped data, and unscoped report aggregation. `delivery/blockers.md` WATCH-004, WATCH-005, and WATCH-006 promote violations to P0 blockers. | None for G8 pass. |
| SEC-G8-006 | P2 Issue | Direct SEC/PM/PRIV/ABUSE ID traceability is behaviorally covered but not always written directly into each task row. | Tasks primarily cite ACC, TM, PSM, INV, TP-AUTH, and ABT IDs. This is sufficient for G8 pass because security behavior maps through PSM/test model, but later reviewers would benefit from direct references to `SEC-*`, `PM-*`, `PRIV-*`, and `ABUSE-*` for high-risk tasks such as TASK-002 and TASK-018 to TASK-023. | Before implementation task kickoff or G10 QA hardening, add direct security requirement/abuse-case references to the relevant delivery task rows or QA test cases without changing scope or priority. |
| SEC-G8-007 | P2 Issue | Import/export operation-log evidence is split across TASK-020 and TASK-022, which is acceptable but should be made explicit during execution. | TASK-020 depends on TASK-022 and covers scoped import/export, dangerous cells, safe row errors, retention, and run states. TASK-022 covers import/export operation-log events. Because the security requirement spans ACC-020 and ACC-022, closure evidence must show both task outputs together. | During TASK-020/TASK-022 execution, ensure import/export operation-log event tests are referenced from both closure records or a shared evidence note so ACC-020 cannot close without the required operation-log evidence. |

## P0/P1 Blockers

No P0/P1 blockers were found in Security Compliance review.

Specifically:

- No P0/P1 authentication, authorization, privacy, audit, import/export,
  report, reminder, duplicate-warning, retention, or abuse-case requirement is
  missing from the G8 delivery plan.
- No task relies on frontend-only authorization, client-side filtering of
  unauthorized records, mock data, static UI, in-memory-only storage, or
  non-persistent behavior to satisfy a P0/P1 path.
- Authorization-before-query/output is explicitly carried into lists,
  reports, reminders, import/export, duplicate warnings, history, and logs.
- Safe denial and masking requirements remain testable through TM, TP-AUTH,
  INV, and ABT coverage.

## P2 Improvements

| ID | Improvement | Owner |
|---|---|---|
| SEC-G8-006 | Add direct `SEC-*`, `PM-*`, `PRIV-*`, and `ABUSE-*` references to high-risk delivery task rows or QA cases for easier security reverse traceability. | Security Compliance + QA TDD + Task Planner |
| SEC-G8-007 | Tie TASK-020 and TASK-022 closure evidence together for CSV import/export operation-log proof. | crm-backend-operations-reporting + crm-qa-tdd-owner + crm-audit-traceability-closure |

## Recommendation

Security Compliance recommends G8 pass from the security, privacy, permission,
audit, and abuse-case planning perspective.

The P2 improvements should be treated as traceability and execution-hardening
work. They do not block G8 because the required security behaviors are already
covered through `delivery/`, `modeling/PSM.md`, and `modeling/test-model.md`.
