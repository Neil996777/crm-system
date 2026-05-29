# G8 Task Planning Repair Note

## Document Control

- Project: CRM System
- Repair Owner: task-planner
- Date: 2026-05-27
- Status: Ready for focused G8 re-review

Implementation remains blocked until G8 passes.

## Reviewed Blocked Lanes

- `archive/reviews/g8-task-planning/integration-g8-task-review.md`
- `archive/reviews/g8-task-planning/qa-tdd-g8-task-review.md`
- `archive/reviews/g8-task-planning/domain-modeling-g8-task-review.md`
- `archive/reviews/g8-task-planning/architecture-g8-task-review.md`

## Repairs Applied

| Review Finding | Severity | Repair |
|---|---|---|
| Record-local history was sequenced after mutation tasks that required history evidence. | P0 | `delivery/tasks.md` now defines incremental history evidence: each mutation task owns its own transaction-linked history append tests, and `TASK-014` owns cross-record query/UI review. `delivery/task-dependencies.md` explicitly records this closure rule. |
| `TASK-017` restore/smoke omitted report/reminder prerequisites. | P0 | `TASK-017` dependencies now include `TASK-018`, `TASK-021`, `TASK-022`, and `TASK-023`; ACC-017 evidence now includes reports and reminders. |
| Import/export and operation logs had a hidden closure cycle. | P1 | `TASK-022` now owns log infrastructure/query and current event coverage; `TASK-020` owns import/export events against that infrastructure and contributes ACC-022 evidence. |
| Integration docs contained placeholder blockers/status rows. | P1 | `docs/integration/*` now explicitly state that G11 execution has not started and point to `delivery/` for current G8 planning evidence. Fake open blocker rows were removed. |
| `TASK-012` omitted contact and quote related-record coverage. | P0 | `TASK-012`, `delivery/test-plan.md`, `delivery/test-cases.md`, and `modeling/test-model.md` now include lead, customer, contact, opportunity, quote, contract, and payment contexts. |
| QA test-plan/test-cases were placeholder-only. | P1 | Concrete G8 QA baselines were added under `delivery/test-plan.md` and `delivery/test-cases.md`; `docs/qa/*` now point to delivery baselines and reserve docs for later G10 reports. |
| User role/status lifecycle, stale-role recheck, and last-Administrator behavior lacked explicit task ownership. | P0 | `TASK-002` now includes `PSM-API-003`, user handlers/service, role/status lifecycle, stale-session tests, and last-Administrator tests/manual verification. |
| Shared OpenAPI/generated client/error/enum contract was missing. | P1 | `TASK-001` establishes shared OpenAPI/generated client/error/enum assets and contract tests; planning rules require every API-changing task to update those contract assets. |
| Money convention appeared after opportunity amount. | P0 | `TASK-007` now creates Money DTO/value-object/minor-unit convention and tests before quote/contract/payment/report tasks depend on it. |
| Backend file paths did not clearly preserve architecture boundaries. | P1 | `delivery/tasks.md` now defines required `app/domain/repository/workflow` path conventions and backend feature task expansion rules; planned paths were mechanically aligned away from old `internal/crm`, `internal/postgres`, and `internal/http` locations. |

## Files Updated

- `delivery/tasks.md`
- `delivery/task-dependencies.md`
- `delivery/delivery-plan.md`
- `delivery/acceptance-task-map.md`
- `delivery/blockers.md`
- `delivery/test-plan.md`
- `delivery/test-cases.md`
- `docs/qa/test-plan.md`
- `docs/qa/test-cases.md`
- `docs/qa/qa-report.md`
- `docs/qa/regression-report.md`
- `docs/integration/integration-report.md`
- `docs/integration/acceptance-status.md`
- `docs/integration/blocker-list.md`
- `modeling/test-model.md`
- `modeling/traceability-matrix.md`
- `PROJECT_CONTEXT.md`
- `docs/planning/agent-registry.md`

## Recommendation

Request focused re-review from Integration Owner, QA/TDD, Domain Modeling, and
Architecture. Product Manager and Security Compliance had no P0/P1 blockers in
the first pass, but may review the repaired package if desired.
