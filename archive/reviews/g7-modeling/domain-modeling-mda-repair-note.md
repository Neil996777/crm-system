# Domain Modeling MDA Repair Note

## Document Control

- Project: CRM System
- Gate Context: G7 MDA Pre-Task Review
- Owner Agent: Domain Modeling
- Date: 2026-05-27
- Output Location: `archive/reviews/g7-modeling/domain-modeling-mda-repair-note.md`
- Status: Repair complete; awaiting Architecture focused re-review

Implementation remains blocked until G8 passes.

## Scope

This note records the Domain Modeling repair for the Architecture G7/G8
pre-task P0 findings. It is process evidence and is intentionally stored under
`archive/`, not `docs/`.

No implementation code was written. No task plan, delivery plan, or G8 task
artifact was created.

## Modified Active Modeling Artifacts

- `modeling/PSM.md`
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`
- `PROJECT_CONTEXT.md`

## Repair Summary

| Architecture Finding | Repair |
|---|---|
| Money representation and authoritative calculation were not finalized in PSM. | Added `PSM-MONEY-001` through `PSM-MONEY-006`, selected integer minor-unit representation, defined Money DTO, Go value-object strategy, PostgreSQL `BIGINT` minor-unit columns, parsing/rounding rules, exact quote/contract/payment comparisons, Won full-payment checks, and authorized SQL report sums. |
| Resource-level authorization policy and scope loader mapping were missing. | Added `PSM-AUTHZ-001` through `PSM-AUTHZ-014` and `PSM-SCOPE-001` through `PSM-SCOPE-014`, covering endpoint/action policy functions, resource-specific scope loaders for Administrator, Sales Manager, and Sales, authorization-before-query, safe denial mapping, stale session/role recheck, and last-Administrator transaction mapping. |
| Deployment, backup, restore, and retention mapping were too thin. | Added `PSM-INFRA-001` through `PSM-INFRA-007` and `PSM-RET-001` through `PSM-RET-006`, covering deployment units, required environment variables, migration ordering, pre-migration backup, encrypted backup schedule, checksum metadata, retention, restore rehearsal, and operational evidence. |

## Traceability Repairs

- `modeling/traceability-matrix.md` now maps affected acceptance rows to the
  new PSM IDs, including `ACC-002`, `ACC-004`, `ACC-009` through `ACC-011`,
  `ACC-013`, `ACC-015` through `ACC-018`, and `ACC-020` through `ACC-023`.
- The `ACC-004` Architecture Source wording was clarified to distinguish
  architecture API rule `API-007` from explicit PSM mappings.
- `modeling/test-model.md` now maps affected test rows and invariant concepts
  to the new money, authorization, infrastructure, restore, and retention PSM
  IDs.

## Gate Position

No known Domain Modeling-owned P0/P1 blocker remains after this repair.

G8 Task Planning must not begin until Architecture completes a focused
re-review and passes these repaired items.
