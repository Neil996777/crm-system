# Product Manager G8 Task Planning Review

## Decision

Passed from Product Manager review perspective, with P2 improvements.

No P0/P1 blockers were found. Product acceptance scope is preserved, P0/P1
items are not downgraded, and the moved `delivery/` artifacts describe
end-to-end user capability delivery rather than a technical chores list.

Implementation remains blocked until the full G8 gate is approved by the Task
Planner owner and all required reviewers.

## Reviewed Inputs

| Input | Review Purpose |
|---|---|
| `../../company/operating-model.md` | G8 reviewer role, no-downgrade rule, and blocker behavior. |
| `../../standards/acceptance-matrix-standard.md` | Product acceptance source-of-truth and required trace fields. |
| `../../standards/status-and-priority-standard.md` | Priority, status, and blocker severity language. |
| `PROJECT_CONTEXT.md` | Current phase, active artifacts, and implementation block before G8 pass. |
| `docs/product/acceptance-matrix.md` | Source of truth for ACC-001 to ACC-023. |
| `delivery/tasks.md` | End-to-end task definitions, user capability slices, test requirements, manual verification, and no-downgrade controls. |
| `delivery/acceptance-task-map.md` | ACC to TASK/TM/evidence mapping. |
| `delivery/delivery-plan.md` | Milestone sequencing and release-blocking coverage. |
| `delivery/task-dependencies.md` | Capability dependency order and blocker triggers. |
| `delivery/blockers.md` | Current blocker register and carry-forward watch items. |
| `modeling/traceability-matrix.md` | ACC to CIM/PIM/PSM/TM/TASK trace. |

## Findings

| ID | Severity | Finding | Evidence | Required Action |
|---|---|---|---|---|
| PM-G8-001 | Pass | P0/P1 product acceptance coverage is complete in moved delivery artifacts. | `delivery/acceptance-task-map.md` maps ACC-001 to ACC-023 one-to-one to TASK-001 to TASK-023 and TM-001 to TM-023. `delivery/tasks.md` task index also lists all 23 tasks with priority, user capability, ACC, TM, owner, dependencies, and evidence status. | None for Product Manager G8 pass. |
| PM-G8-002 | Pass | Tasks are framed as end-to-end user capabilities, not standalone technical chores. | `delivery/tasks.md` defines each task with task goal, business capability, acceptance item, production/test file changes, automated tests, manual verification, MDA trace, TDD guard, no-downgrade rules, and blocker record. | None for Product Manager G8 pass. |
| PM-G8-003 | Pass | P0/P1 no-downgrade controls are explicit and preserved. | `delivery/tasks.md`, `delivery/delivery-plan.md`, and `delivery/blockers.md` prohibit mock, stub, TODO, static-only, localStorage-only, in-memory-only, or non-persistent behavior for core scope. | None for Product Manager G8 pass. |
| PM-G8-004 | Pass | Delivery artifacts were moved out of `docs/planning`; active design docs are not polluted by task documents. | `delivery/` contains `tasks.md`, `task-dependencies.md`, `delivery-plan.md`, `acceptance-task-map.md`, and `blockers.md`; `docs/planning/` only contains `agent-registry.md`. | None. |
| PM-G8-005 | P2 Issue | Product acceptance matrix back-links are stale now that G8 task planning exists. | `docs/product/acceptance-matrix.md` still shows `Related Tasks` as `N/A until planning` and `Related Tests` as `N/A until QA planning` for ACC-001 to ACC-023, while `delivery/acceptance-task-map.md` and `modeling/traceability-matrix.md` already provide TASK/TM mappings. | After G8 review consensus, update acceptance matrix back-links to TASK-001 to TASK-023 and TM-001 to TM-023, without changing acceptance wording or priority. |
| PM-G8-006 | P2 Issue | Some task `Files to modify` rows use broad labels instead of concrete file paths. This does not reduce product scope, but it weakens execution handoff precision. | Examples include cross-task phrases such as "protected handlers from later tasks", "Entity detail pages", "List/detail handlers and pages", "All repositories from TASK-001 to TASK-015", and service groups for history/log events. | Before implementation starts, Task Planner or owning implementation agents should replace broad labels with concrete files as those files exist, while keeping the same ACC/TM/task scope. |

## P0/P1 Blockers

No P0/P1 blockers were found in Product Manager review.

Specifically:

- No P0/P1 acceptance item from `docs/product/acceptance-matrix.md` is missing
  from `delivery/acceptance-task-map.md`.
- No P0/P1 item is downgraded, deleted, merged away, or marked optional.
- No task claims completion by mock, static UI, TODO, local-only state, or
  non-persistent behavior.
- Every task has a planned user capability, acceptance item, real production
  code surface, automated test requirement, and reproducible manual
  verification path.

## P2 Improvements

| ID | Improvement | Owner |
|---|---|---|
| PM-G8-005 | Refresh `docs/product/acceptance-matrix.md` related task/test back-links after G8 consensus. | Product Manager + Task Planner + QA TDD |
| PM-G8-006 | Tighten broad `Files to modify` entries into concrete file paths before implementation task kickoff. | Task Planner + owning CRM implementation agents |

## Recommendation

Product Manager recommends G8 pass from product acceptance perspective.

The Task Planner gate owner should continue collecting the remaining required
reviewer decisions. The P2 items above should be handled as cleanup and handoff
hardening, not as Product Manager blockers to G8.
