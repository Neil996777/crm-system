# Architecture MDA Focused Re-Review

## Document Control

- Project: CRM System
- Review Gate: G7 focused re-review before G8 Task Planning
- Reviewer Agent: Architecture
- Reviewed Date: 2026-05-27
- Output Location: `archive/reviews/g7-modeling/architecture-mda-focused-re-review.md`
- Review Scope: Domain Modeling repair for Architecture G7/G8 pre-task P0 findings.

Implementation remains blocked until G8 passes.

## Review Decision

**Passed**

The focused repair closes the three Architecture-owned P0 blockers from
`archive/reviews/g7-modeling/architecture-mda-pre-task-review.md`.

No implementation code was reviewed or changed. No task planning artifacts were
created. `modeling/` was reviewed only and was not edited by this re-review.

## Reviewed Inputs

- `archive/reviews/g7-modeling/architecture-mda-pre-task-review.md`
- `archive/reviews/g7-modeling/domain-modeling-mda-repair-note.md`
- `modeling/PSM.md`
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`
- `PROJECT_CONTEXT.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`

## Focused Findings

| Original Blocker | Result | Re-Review Finding |
|---|---|---|
| AMDA-BLOCK-001: Money representation and authoritative calculation not finalized in PSM. | Closed | `PSM-MONEY-001` through `PSM-MONEY-006` define API Money DTOs, Go `Money` value-object strategy, PostgreSQL `BIGINT` minor-unit storage, decimal parsing and rejection rules, exact quote/contract/payment comparisons, Won full-payment checks, and authorized report sums. |
| AMDA-BLOCK-002: PSM lacks resource-level authorization policy and scope loader mapping. | Closed | `PSM-AUTHZ-001` through `PSM-AUTHZ-014` and `PSM-SCOPE-001` through `PSM-SCOPE-014` define endpoint/action policy functions, policy inputs, resource scope loaders, authorization-before-query, safe denial behavior, stale session/role recheck, and last Administrator mechanics. |
| AMDA-BLOCK-003: Deployment, backup, restore, and retention too thin in PSM. | Closed | `PSM-INFRA-001` through `PSM-INFRA-007` and `PSM-RET-001` through `PSM-RET-006` define deployment units, required environment variables, migration ordering, pre-migration backup, encrypted backup metadata/checksum, retention, restore rehearsal, and operational evidence. |

## Money Review

| Check | Result | Evidence |
|---|---|---|
| API DTO strategy is explicit. | Passed | `PSM-MONEY-001` and `PSM-DTO-012` define `{ amountMinor: string, currencyCode: string }`, with frontend display values non-authoritative. |
| Go type strategy is explicit. | Passed | `PSM-MONEY-002` requires a single Go `Money` value object with `MinorUnits int64` and `CurrencyCode string`; floating-point authority is forbidden. |
| PostgreSQL storage is explicit. | Passed | `PSM-MONEY-003` maps authoritative money columns to `BIGINT` minor units plus `CHAR(3)` currency code across `PSM-DB-006` through `PSM-DB-010`. |
| Parsing, rounding, comparison, report sums, and Won full-payment are executable. | Passed | `PSM-MONEY-004` rejects over-precision instead of silent rounding; `PSM-MONEY-005` defines exact minor-unit comparisons and full-payment checks; `PSM-MONEY-006` defines authorized SQL sums. |
| Affected ACC/TM/PSM IDs are mapped. | Passed | `ACC-009`, `ACC-010`, `ACC-011`, `ACC-013`, `ACC-018`, and `ACC-023` map to `PSM-MONEY-*` and `TM-009` through `TM-013`, `TM-018`, `TM-023` where relevant. |

## Authorization Review

| Check | Result | Evidence |
|---|---|---|
| Endpoint/action to policy mapping is sufficient for G8 planning. | Passed | `PSM-AUTHZ-001` through `PSM-AUTHZ-014` cover users, leads, company/contact, opportunities, quotes, contracts, payments, work items, history/logs, duplicate warnings, import/export, reminders, and reports. |
| Resource scope loaders cover required resources. | Passed | `PSM-SCOPE-001` through `PSM-SCOPE-014` cover lead, company/customer, contact, opportunity, quote, contract, payment, activity/note/task, history, operation log, import/export, reports, reminders, and duplicate warning candidates. |
| Authorization-before-query is clear. | Passed | PSM mandates scoped predicates before rows, aggregates, reminders, imports, exports, reports, and duplicate warning lookup are returned. |
| Safe error mapping is clear. | Passed | PSM distinguishes safe `NOT_FOUND`, `PERMISSION_DENIED`, `BUSINESS_RULE_BLOCKED`, and duplicate warning responses without exposing restricted details. |
| Stale session/role recheck is clear. | Passed | `RequireActiveSession` reloads user status/role version from PostgreSQL before protected API execution and rejects disabled, revoked, expired, or stale-role sessions before repository queries. |
| Last Administrator mechanics are clear. | Passed | Role/status mutation uses `PSM-TX-002`, locks target user and active Administrator set/count, writes `EVT-LAST-ADMIN-BLOCKED` on blocked removal, and commits without changing target role/status. |

## Deployment, Backup, Restore, And Retention Review

| Check | Result | Evidence |
|---|---|---|
| Required environment variables are mapped. | Passed | `PSM-INFRA-002` lists CRM environment, origin, timezone, default currency, database, session, password pepper if selected, backup, log, and CORS variables. |
| Migration order and pre-migration backup are mapped. | Passed | `PSM-INFRA-003` requires encrypted pre-migration backup, checksum, current migration version, staging rehearsal, ordered migrations, and health/persistence smoke checks. |
| Backup encryption, checksum, storage, and retention are mapped. | Passed | `PSM-INFRA-004`, `PSM-INFRA-005`, and `PSM-RET-006` define nightly encrypted PostgreSQL backups to DigitalOcean Spaces, metadata/checksum, and retention windows. |
| Restore rehearsal is mapped. | Passed | `PSM-INFRA-006` requires isolated restore verification covering migrations, active Administrator login, disabled-user denial, core CRM data, history, operation logs, reports, and reminders. |
| Application retention is mapped. | Passed | `PSM-RET-001` through `PSM-RET-006` cover core CRM records, contracts/payments, record-local history, global operation logs, import/export metadata, generated files, and backups. |
| Trace to TM-016/TM-017/TM-022 exists. | Passed | `traceability-matrix.md` and `test-model.md` map persistence, production operation, backup/restore, logs, retention, and operational evidence to `TM-016`, `TM-017`, and `TM-022`. |

## New P0/P1 Blockers

None found.

## P2 Improvements

| ID | Improvement | Owner | Rationale |
|---|---|---|---|
| AMDA-FRR-P2-001 | During G8, split `PSM-AUTHZ-*` and `PSM-SCOPE-*` into explicit backend unit/integration test tasks and frontend denial-state tasks. | Task Planner / QA TDD | The model is sufficient, but task decomposition should preserve policy granularity. |
| AMDA-FRR-P2-002 | During G8, create separate operational tasks for deploy, migration, backup, restore rehearsal, smoke test, and evidence retention under `ACC-017`. | Task Planner / Integration Owner | Prevents production readiness from being collapsed into one broad infrastructure task. |
| AMDA-FRR-P2-003 | During G8, create contract-test tasks for Money DTO examples and invalid scale handling. | QA TDD / Backend / Frontend | Ensures frontend, OpenAPI, Go, and PostgreSQL stay aligned on minor-unit money behavior. |

## No-Downgrade Check

| Check | Result |
|---|---|
| P0/P1 acceptance preserved | Passed |
| P0/P1 status not weakened | Passed |
| Product acceptance overridden by Architecture or MDA | Not found |
| Mock/static/non-persistent behavior introduced | Not found |
| Implementation started before G8 | Not found in reviewed artifacts |
| Task planning artifacts created during re-review | Not found |

## Recommendation

Architecture recommends entering **G8 Task Planning**.

This is not an approval to implement. Implementation remains blocked until G8
Task Planning passes with traceable tasks, tests, integration coverage, and no
open P0/P1 blocker.

## Files Modified By This Review

- `archive/reviews/g7-modeling/architecture-mda-focused-re-review.md`
- `PROJECT_CONTEXT.md`
