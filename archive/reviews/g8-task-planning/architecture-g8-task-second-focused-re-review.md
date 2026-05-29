# Architecture G8 Task Planning Second Focused Re-Review

## Decision

Passed from Architecture second focused re-review perspective.

The remaining Architecture blocker `ARCH-G8-BLOCKER-004` is now resolved. The
PSM platform module runtime paths have been repaired to align with the
`app/domain/repository/workflow` boundaries already defined in `delivery/tasks.md`.

Implementation remains blocked until the full G8 gate is approved by the Task
Planner owner and all required reviewers.

## Blocker Resolution Result

| Blocker | Previous Status | Second Focused Re-Review Result | Evidence | Remaining Action |
|---|---|---|---|---|
| ARCH-G8-BLOCKER-004: Backend paths did not explicitly preserve accepted layer boundaries. | Partially resolved; still blocking. | Resolved. | `delivery/tasks.md` requires HTTP handlers under `apps/api/internal/app/http/`, application use cases under `apps/api/internal/app/`, domain models/policies/value objects under `apps/api/internal/domain/`, workflow orchestration under `apps/api/internal/workflow/`, repository interfaces under `apps/api/internal/repository/`, and PostgreSQL implementations under `apps/api/internal/repository/postgres/`. `modeling/PSM.md` Platform Modules now map PSM-HTTP, PSM-AUTH, PSM-AUTHZ, PSM-USERS, PSM-CRM, PSM-WORKFLOW, PSM-AUDIT, PSM-DUP, PSM-IMPORTEXPORT, PSM-REPORT, PSM-REMINDER, PSM-JOBS, and PSM-REPO to those same path conventions. | None for Architecture G8. Keep the execution watch item for boundary drift during implementation. |

## New P0/P1 Architecture Blockers

No new P0/P1 Architecture blocker was introduced by this PSM path repair.

## Recommendation

Architecture recommends G8 pass from Architecture review perspective.

The Task Planner gate owner should continue using the repaired `delivery/`
artifacts and PSM as the implementation boundary source, with no code work
starting until G8 is formally approved.
