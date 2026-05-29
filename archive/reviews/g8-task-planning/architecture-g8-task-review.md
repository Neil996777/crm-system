# Architecture G8 Task Planning Review

## Decision

Blocked from Architecture review perspective.

The moved `delivery/` artifacts are materially stronger than a development task
list: they map ACC-001 to ACC-023 to TASK-001 to TASK-023, include production
file surfaces, tests, manual verification, MDA trace, TDD guards, no-downgrade
language, and blocker rules. However, Architecture cannot recommend G8 pass
until the P0/P1 blockers below are repaired because the current plan leaves
shared contracts, Money sequencing, backend module boundaries, and operation-log
dependencies partially non-executable against the accepted architecture and PSM.

Implementation remains blocked until the Task Planner gate owner approves G8
with all required reviewer decisions.

## Reviewed Inputs

| Input | Review Purpose |
|---|---|
| `../../company/operating-model.md` | G8 owner/reviewer rule, no-downgrade rule, blocker behavior. |
| `../../standards/acceptance-matrix-standard.md` | Evidence, traceability, and P0/P1 completion rules. |
| `../../standards/status-and-priority-standard.md` | Severity and gate status language. |
| `../../workflows/software-delivery.md` | Phase 9 task-planning criteria and implementation block before G8 pass. |
| `PROJECT_CONTEXT.md` | Current G8 review status and active artifact locations. |
| `delivery/tasks.md` | Planned end-to-end task scope, code paths, tests, dependencies, manual verification, MDA trace. |
| `delivery/task-dependencies.md` | Executable order, dependency rules, and blocker triggers. |
| `delivery/delivery-plan.md` | Milestone sequencing and release-blocking coverage. |
| `delivery/acceptance-task-map.md` | ACC to TASK/TM/evidence mapping. |
| `delivery/blockers.md` | Open blockers and carry-forward watch items. |
| `docs/architecture/architecture.md` | Accepted system architecture, module map, dependency direction, shared contract and persistence decisions. |
| `docs/architecture/module-boundaries.md` | Backend/frontend/shared/infra boundaries and dependency rules. |
| `docs/architecture/api-spec.md` | REST/OpenAPI source-of-truth, endpoint groups, safe errors, authz and idempotency rules. |
| `docs/architecture/data-design.md` | PostgreSQL persistence, table/constraint/transaction/retention expectations. |
| `docs/architecture/authz-architecture.md` | Backend authn/authz enforcement, role scopes, safe denial, last Administrator rules. |
| `docs/architecture/frontend-backend-contract.md` | OpenAPI/generated client assets, DTO/error/UI state contracts, frontend/backend separation. |
| `docs/architecture/integration-design.md` | Web/API/DB/job/audit/import/export/report/backup integration obligations. |
| `docs/architecture/deployment-notes.md` | Production topology, environment variables, backup/restore, runbook and verification requirements. |
| `docs/architecture/risk-register.md` | Architecture risks that must be mapped into task, QA, integration, or audit evidence. |
| `modeling/PSM.md` | Platform-specific modules, API/data/authz/scope/Money/infra/retention mappings. |
| `modeling/traceability-matrix.md` | ACC to CIM/PIM/PSM/TM/TASK traceability. |
| `modeling/test-model.md` | TM, invariant, abuse, and operational verification expectations. |

## Findings

| ID | Severity | Finding | Evidence | Required Action |
|---|---|---|---|---|
| ARCH-G8-001 | P1 Blocker | Shared API contract work is not executable against the accepted architecture. The task plan lists per-feature `packages/shared/src/*.ts` files, but no task adds or modifies the OpenAPI source, generated TypeScript client, enum/error catalog, examples, or contract tests required by the architecture. This weakens frontend/backend separation because implementation agents would have to guess whether hand-written TypeScript files or OpenAPI are authoritative. | `docs/architecture/api-spec.md:15` to `docs/architecture/api-spec.md:19` define REST under `/api/v1`, OpenAPI at `packages/shared/openapi/crm.v1.yaml`, and generated TypeScript client usage. `docs/architecture/frontend-backend-contract.md:23` to `docs/architecture/frontend-backend-contract.md:31` list required shared contract assets. `delivery/tasks.md:89`, `delivery/tasks.md:137`, and later tasks list `packages/shared/src/*.ts` but not the OpenAPI/generated contract assets. | Update G8 planning so each affected capability task, or a explicitly mapped contract-support slice tied to the relevant ACC/TM tasks, includes `packages/shared/openapi/crm.v1.yaml`, generated client/type outputs, role/status/error enum assets, API examples, and contract tests. Keep backend as the authority for validation/authz/business behavior. |
| ARCH-G8-002 | P0 Blocker | Money sequencing is inconsistent for the P0 opportunity task. TASK-007 requires an opportunity amount and expected close date, but the shared Money DTO/value object is introduced in TASK-009, and TASK-007 says it may modify `packages/shared/src/money.ts` only if TASK-009 created it earlier. TASK-009 is later in both sequence and dependency order, so the plan permits ACC-007 implementation before the Money convention exists. | `delivery/tasks.md:59` places TASK-007 before TASK-009. `delivery/tasks.md:229` to `delivery/tasks.md:244` make amount part of ACC-007 completion and prohibit client-only/optional amount handling. `delivery/tasks.md:234` references `packages/shared/src/money.ts` only if TASK-009 created it earlier. `delivery/tasks.md:281` introduces `packages/shared/src/money.ts` in TASK-009. `modeling/PSM.md:165` to `modeling/PSM.md:170` require integer minor-unit Money representation, and `modeling/PSM.md:183` to `modeling/PSM.md:184` include `PSM-DB-006` opportunity money columns. | Move Money DTO, Go value object, PostgreSQL minor-unit convention, and tests into TASK-007 or an earlier dependency that TASK-007 explicitly depends on. Also ensure TASK-007 traces to the relevant `PSM-MONEY-*` elements wherever opportunity amount is accepted, persisted, displayed, or reported. |
| ARCH-G8-003 | P1 Blocker | Operation-log dependency order is not executable for ACC-020/ACC-022. TASK-020 depends on TASK-022 for operation-log evidence, while TASK-022 claims completion only after import/export events are written and manually verified. Import/export services do not exist until TASK-020, so TASK-022 cannot complete as written before TASK-020, and TASK-020 cannot start as written without TASK-022. | `delivery/task-dependencies.md:42` says TASK-020 depends on TASK-022. `delivery/task-dependencies.md:44` says TASK-022 depends on TASK-001, TASK-002, and TASK-014 only. `delivery/tasks.md:593` to `delivery/tasks.md:605` require TASK-022 to add operation logs and modify auth, owner/stage/status, quote, contract, payment, archive, import, and export services, and to manually trigger import/export log evidence. `delivery/tasks.md:545` to `delivery/tasks.md:557` creates import/export in TASK-020 and requires retained run metadata/evidence. | Split operation-log infrastructure from event-coverage closure, or change dependencies so import/export event coverage is completed where import/export exists. A valid repair is: TASK-022 provides append-only log infrastructure/query/access denial and earlier event coverage; TASK-020 then adds import/export operation events and evidence against the existing log writer, with ACC-022 evidence updated accordingly. |
| ARCH-G8-004 | P1 Blocker | Planned backend file paths do not consistently preserve the accepted layer boundaries. Architecture and PSM name application/domain/repository responsibilities and repository interfaces, while the task plan mostly writes handlers plus `internal/crm` and `internal/postgres` files. Without explicit repository interfaces, domain policy locations, and application-service boundaries, implementers can satisfy file lists while bypassing the required HTTP -> service -> domain/policy -> repository -> PostgreSQL dependency direction. | `docs/architecture/architecture.md:86` to `docs/architecture/architecture.md:104` define `apps/api/internal/app/`, `apps/api/internal/domain/`, `apps/api/internal/repository/`, and dependency direction. `docs/architecture/module-boundaries.md:37` to `docs/architecture/module-boundaries.md:50` define `http`, `auth`, `authorization`, `crm`, `workflow`, `audit`, and `repository` dependencies. `modeling/PSM.md:38` to `modeling/PSM.md:48` includes `apps/api/internal/users/`, `apps/api/internal/workflow/`, and `apps/api/internal/repository/`. `delivery/tasks.md` repeatedly lists `apps/api/internal/postgres/*_repository.go` and feature files, but not the repository interfaces or application/domain layer files needed to enforce the boundary. | Amend task planned files to include the architectural layer contracts: application services/use cases, domain policy/value objects, repository interfaces, PostgreSQL implementations, and workflow services where applicable. If `internal/postgres` is intended as an implementation detail under repository, record that convention explicitly in delivery planning and align it with PSM/module-boundary names before G8 pass. |
| ARCH-G8-005 | Pass | P0/P1 acceptance coverage and high-level MDA trace are present. | `delivery/acceptance-task-map.md` maps ACC-001 to ACC-023 one-to-one to TASK-001 to TASK-023 and TM-001 to TM-023. `modeling/traceability-matrix.md` maps each ACC to CIM/PIM/PSM/TM/TASK. `delivery/tasks.md` includes an MDA trace row for every task. | None for coverage. Keep these mappings intact when repairing blockers. |
| ARCH-G8-006 | Pass | Persistence and production-operation intent is aligned with architecture at the milestone level. | TASK-016 requires PostgreSQL-backed restart persistence, no mock/static/in-memory core path, and representative refresh/relogin/restart evidence. TASK-017 requires Compose/Caddy deployment, migrations, encrypted backup, checksum, restore rehearsal, smoke, no-secrets verification, and deployment-equivalent manual evidence. | None for the milestone intent. See P2 improvements for production runbook hardening. |

## P0/P1 Blockers

| Blocker ID | Severity | Affected Acceptance / Tasks | Blocking Condition | Required Repair Before Architecture G8 Pass |
|---|---|---|---|---|
| ARCH-G8-BLOCKER-001 | P1 | ACC-001 to ACC-023 / TASK-001 to TASK-023 | Shared OpenAPI/generated client/error/enum contract assets are missing from planned delivery file paths. | Add executable shared contract file and test coverage planning aligned to `packages/shared/openapi/crm.v1.yaml`, `packages/shared/generated/`, shared enums/errors, and API examples. |
| ARCH-G8-BLOCKER-002 | P0 | ACC-007, ACC-009, ACC-010, ACC-011, ACC-013, ACC-018, ACC-023 / TASK-007, TASK-009 to TASK-011, TASK-013, TASK-018, TASK-023 | Money convention appears after the first P0 task that requires amount persistence. | Introduce Money DTO/value object/DB convention before or inside TASK-007 and make downstream tasks depend on the same asset. |
| ARCH-G8-BLOCKER-003 | P1 | ACC-020, ACC-022 / TASK-020, TASK-022 | TASK-020 and TASK-022 contain a hidden circular dependency for import/export operation-log evidence. | Repair sequencing by separating log infrastructure from import/export event coverage or by reordering dependencies so required services exist before claimed evidence. |
| ARCH-G8-BLOCKER-004 | P1 | ACC-001 to ACC-023, especially ACC-016 / all backend tasks | Planned backend paths do not explicitly include accepted application/domain/repository boundaries and can be implemented in a way that bypasses architecture dependency direction. | Align task file plans with `app`/`domain`/`repository`/`workflow`/PostgreSQL implementation boundaries or document an equivalent path convention in delivery planning and PSM. |

## P2 Improvements

| ID | Improvement | Rationale | Suggested Owner |
|---|---|---|---|
| ARCH-G8-P2-001 | Add explicit production runbook file planning for rollback, user bootstrap, and secret rotation in TASK-017. | `docs/architecture/deployment-notes.md` requires these runbooks before release. TASK-017 currently plans deploy and backup/restore runbooks but not all production-readiness runbooks. | `crm-integration-operations-delivery` + `crm-backend-foundation-platform` |
| ARCH-G8-P2-002 | Replace broad `Files to modify` phrases with concrete paths once blocker repairs introduce layer files. | Phrases such as "protected handlers from later tasks", "Entity detail pages", and grouped service names are understandable during draft planning but weaken implementation handoff precision. | Task Planner + owning implementation agents |
| ARCH-G8-P2-003 | Update `delivery/blockers.md` after repairs to record the resolved Architecture G8 blockers or their watch-item equivalents. | The current blocker register says no open P0/P1 blocker is known for G8 review. After this review, the architecture blockers should be reflected until fixed or formally superseded by revised delivery artifacts. | Task Planner |

## Recommendation

Architecture recommends G8 remain blocked until ARCH-G8-BLOCKER-001 through
ARCH-G8-BLOCKER-004 are repaired in the `delivery/` planning artifacts and
rechecked.

Do not start implementation. The repair should keep the existing ACC/TM/MDA
coverage, preserve P0/P1 no-downgrade rules, and avoid moving delivery task
artifacts back into `docs/`.
