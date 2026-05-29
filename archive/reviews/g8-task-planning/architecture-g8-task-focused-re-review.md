# Architecture G8 Task Planning Focused Re-Review

## Decision

Blocked from Architecture focused re-review perspective.

ARCH-G8-BLOCKER-001, ARCH-G8-BLOCKER-002, and ARCH-G8-BLOCKER-003 are resolved
for G8 planning. ARCH-G8-BLOCKER-004 is partially repaired in
`delivery/tasks.md`, but not fully resolved because `modeling/PSM.md` still
maps the backend runtime modules to the old package locations while the repaired
delivery plan requires the `app/domain/repository/workflow` structure. That
remaining PSM-to-task mismatch keeps the backend architecture boundary
non-authoritative for implementation.

Implementation remains blocked until G8 passes.

## Blocker Resolution Table

| Previous Blocker | Severity | Focused Re-Review Result | Evidence | Remaining Action |
|---|---|---|---|---|
| ARCH-G8-BLOCKER-001: Shared OpenAPI/generated client/error/enum contract assets missing from planned delivery paths. | P1 | Resolved. | `delivery/tasks.md:48` to `delivery/tasks.md:51` makes `packages/shared/openapi/crm.v1.yaml`, `packages/shared/generated/`, shared enum/error catalogs, and `tests/contract/` authoritative. `TASK-001` adds OpenAPI, generated client, enum/error assets, and contract tests. `delivery/test-plan.md` adds a Contract test group for OpenAPI/generated client/enum/error/Money DTO evidence. | None for Architecture G8. Keep WATCH-007 during execution. |
| ARCH-G8-BLOCKER-002: Money convention appeared after first P0 opportunity amount task. | P0 | Resolved. | `TASK-007` now adds `apps/api/internal/domain/money/money.go`, `packages/shared/src/money.ts`, `tests/contract/money_contract.test.ts`, and updates OpenAPI/generated client. Its completion standard requires expected amount as integer minor-unit Money, and its MDA trace includes `PSM-MONEY-001` to `PSM-MONEY-004`. | None for Architecture G8. Minor trace cleanup may later add ACC-007 to PSM Money DTO rows, but the executable sequencing blocker is closed. |
| ARCH-G8-BLOCKER-003: TASK-020 and TASK-022 hidden cycle for import/export operation-log evidence. | P1 | Resolved. | `delivery/tasks.md:56` to `delivery/tasks.md:60` defines incremental global operation-log delivery. `delivery/task-dependencies.md` keeps TASK-020 dependent on TASK-022 for log infrastructure only, while TASK-020 owns import/export event evidence. `TASK-022` explicitly says import/export log verification is performed in TASK-020 and linked back to ACC-022. | None for Architecture G8. Keep WATCH-009 during execution. |
| ARCH-G8-BLOCKER-004: Backend paths did not explicitly preserve accepted layer boundaries. | P1 | Partially resolved; still blocking. | `delivery/tasks.md:35` to `delivery/tasks.md:47` now defines `apps/api/internal/app/http/`, `apps/api/internal/app/`, `apps/api/internal/domain/`, `apps/api/internal/workflow/`, `apps/api/internal/repository/`, and `apps/api/internal/repository/postgres/`, and requires every backend feature task to expand to the applicable layer files before closure. However, `modeling/PSM.md:35` to `modeling/PSM.md:47` still maps PSM modules to old runtime locations such as `apps/api/internal/http/`, `apps/api/internal/auth/`, `apps/api/internal/authorization/`, `apps/api/internal/crm/`, `apps/api/internal/audit/`, `apps/api/internal/importexport/`, `apps/api/internal/reports/`, `apps/api/internal/reminders/`, and `apps/api/internal/jobs/`. | Update `modeling/PSM.md` platform module runtime locations to match the repaired delivery architecture, or explicitly document the equivalence between PSM module names and the new `app/domain/repository/workflow` paths. Until then, implementation agents still have conflicting source-of-truth paths. |

## New P0/P1 Architecture Blockers

No new P0/P1 Architecture blocker was introduced by the repair package beyond
the unresolved portion of ARCH-G8-BLOCKER-004.

## Recommendation

Architecture recommends G8 remain blocked until `modeling/PSM.md` is aligned
with the repaired backend path conventions in `delivery/tasks.md`.

After that targeted PSM repair, a narrow Architecture re-check can focus only on
ARCH-G8-BLOCKER-004. Do not start implementation before G8 approval.
