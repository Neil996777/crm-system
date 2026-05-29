# Domain Modeling G8 Task Review

## Decision

**Gate decision: Blocked for Domain Modeling.**

The moved `delivery/` artifacts cover `ACC-001` to `ACC-023` at a high level,
but G8 should not pass yet because several P0 task-planning gaps prevent a
clean `ACC -> CIM/PIM/PSM -> TM -> TASK -> tests -> reproducible validation`
chain.

## Reviewed Inputs

- `delivery/tasks.md`
- `delivery/task-dependencies.md`
- `delivery/delivery-plan.md`
- `delivery/acceptance-task-map.md`
- `delivery/blockers.md`
- `modeling/CIM.md`
- `modeling/PIM.md`
- `modeling/PSM.md`
- `modeling/domain-model.md`
- `modeling/state-machines.md`
- `modeling/domain-events.md`
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`
- `docs/product/acceptance-matrix.md`

## Findings

| ID | Severity | Finding | Evidence | Required correction |
|---|---|---|---|---|
| DM-G8-001 | P0 Blocker | Record-local history is planned after the mutation tasks that already require history evidence, creating an impossible closure dependency. | `TASK-003` requires owner-change history, while `TASK-014` later modifies mutation services from `TASK-003` to `TASK-013` and depends on `TASK-003 to TASK-013` being complete. `delivery/task-dependencies.md` also says a dependency is satisfied only with production code, tests, and evidence. This conflicts with `ACC-008` and `ACC-014`, which require history events and record-local history. | Re-sequence history as an enabling slice before or inside the first mutation tasks, or explicitly mark each mutation task as not closable until history instrumentation/tests are included. Update dependencies so no task can be "complete" before required transactional history evidence exists. |
| DM-G8-002 | P0 Blocker | `TASK-012` does not cover all P0 related-record contexts for activities, notes, and tasks. | `ACC-012` requires scenarios for related `lead/customer/contact/opportunity/quote/contract/payment` records. `TASK-012` files/tests/manual verification cover leads, companies/customers, opportunities, contracts, and payments, but omit contact and quote contexts. `TM-012` has the same omission. | Expand `TASK-012` and `TM-012` to include contact and quote related-record scope, UI wiring, automated tests, and manual verification. |
| DM-G8-003 | P0 Blocker | User role/status lifecycle and last-Administrator behavior are accepted model elements but lack explicit task-level code/test ownership. | `PIM-CMD-002` maps user role/status change and last-admin guard to `ACC-001`, `ACC-002`, and `ACC-022`; `PSM-API-003` maps `/users`, `/roles` to the same acceptance IDs; `SM-USER` includes Administrator create/disable/enable/change-role transitions. `TASK-001` and `TASK-002` cover login and authorization policy, but their planned files and traces do not include `PSM-API-003`, user management handlers/service, or role/status lifecycle verification. | Add explicit coverage in `TASK-001`/`TASK-002` or a dedicated end-to-end task tied to the same ACC/TM chain for Administrator user creation, role/status changes, stale role/session rejection, and last-active-Administrator blocking. |
| DM-G8-004 | P2 Issue | Some MDA traces use compressed or slash notation that is readable by humans but weak for audit automation. | Examples: `PSM-AUTHZ-001..014`, `PSM-DB-011..013`, `CIM-CAP-004/005`, `PIM-CMD-005..007`. | Prefer explicit IDs in task traces or add a short expansion rule so Audit can mechanically verify every referenced model element. |
| DM-G8-005 | P2 Issue | Modeling metadata status is stale relative to the project context. | Several modeling files still say `Draft for G6/G7 Review` or `Awaiting Architecture Focused Re-Review`, while `PROJECT_CONTEXT.md` says G7 reviews passed and G8 is ready for review. | After blocker repair, update modeling document-control statuses or add a review note explaining that these are accepted G7 inputs with G8 task review underway. |
| DM-G8-006 | P2 Issue | `TASK-007` references `packages/shared/src/money.ts` as if `TASK-009` may have created it earlier, but `TASK-009` depends on `TASK-007`/`TASK-008`. | `TASK-007` files-to-modify says "if created by TASK-009 earlier in execution order"; dependency matrix places `TASK-009` after `TASK-007` and `TASK-008`. | Move shared Money DTO setup to the first task that needs opportunity amount semantics or make it an explicit planned file in `TASK-007`/foundation scope. |

## P0/P1 Blockers

| Blocker ID | Severity | Affected Acceptance | Affected Tasks | Required action |
|---|---|---|---|---|
| DM-G8-001 | P0 | ACC-003, ACC-004, ACC-008, ACC-014, ACC-016 | TASK-003 to TASK-014, TASK-016 | Repair the history/event dependency plan so transactional history is available before dependent mutation tasks claim completion. |
| DM-G8-002 | P0 | ACC-012 | TASK-012, TM-012 | Add contact and quote contexts to planned implementation, automated tests, and manual verification. |
| DM-G8-003 | P0 | ACC-001, ACC-002, ACC-022 | TASK-001, TASK-002, TASK-022 or new task | Add explicit user role/status lifecycle and last-Administrator implementation/test trace. |

No additional P1-only blockers were found. P1 items remain covered at the
acceptance-task map level, subject to the P0 repairs above and later QA,
integration, and audit evidence.

## P2 Improvements

- Expand shorthand MDA references into explicit model IDs to improve audit
  traceability.
- Align modeling document-control statuses with the current G8 review state.
- Correct the `TASK-007` / `TASK-009` Money file sequencing note.

## Recommendation

Do not pass G8 from the Domain Modeling review lane yet. Repair the P0 blockers
in `delivery/` and, where needed, `modeling/test-model.md` or
`modeling/traceability-matrix.md`, then request Domain Modeling re-review. The
plan is close structurally, but the current blocker set would allow task
execution to start with missing or impossible MDA-backed evidence paths.
