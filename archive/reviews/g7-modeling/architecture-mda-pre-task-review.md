# Architecture MDA Pre-Task Review

## Document Control

- Project: CRM System
- Review Gate: G7/G8 pre-task review
- Reviewer Agent: Architecture
- Reviewed Date: 2026-05-27
- Output Location: `archive/reviews/g7-modeling/architecture-mda-pre-task-review.md`
- Review Scope: MDA Modeling draft against accepted architecture design and product acceptance matrix.

Implementation remains blocked until G8 passes.

## Review Decision

**Blocked**

The MDA package correctly carries the overall architecture direction: Go API,
React/Vite frontend, PostgreSQL 16, REST/OpenAPI, server-side sessions,
backend authorization, transactional audit/history, import/export/report
ownership, and DigitalOcean/Caddy/Cloudflare/Spaces deployment direction are
all visible in the PSM and traceability model.

However, the current PSM is not yet concrete enough for G8 Task Planning and
Frontend/Backend implementation to proceed without guessing. The blockers
below are architecture-owned P0/P1 mapping gaps. They must be repaired in
`modeling/PSM.md` and, where affected, `modeling/traceability-matrix.md` before
Architecture should recommend entering task planning.

No implementation code was reviewed or changed. `modeling/` was not edited by
this review.

## Reviewed Inputs

- `modeling/CIM.md`
- `modeling/PIM.md`
- `modeling/PSM.md`
- `modeling/domain-model.md`
- `modeling/state-machines.md`
- `modeling/domain-events.md`
- `modeling/traceability-matrix.md`
- `modeling/test-model.md`
- `docs/architecture/architecture.md`
- `docs/architecture/module-boundaries.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`
- `docs/product/acceptance-matrix.md`

## Positive Findings

| Area | Finding |
|---|---|
| Platform direction | PSM maps backend to `apps/api/` Go, frontend to `apps/web/` React/Vite, shared contracts to `packages/shared/`, PostgreSQL 16, and production infrastructure to DigitalOcean/Caddy/Cloudflare/Spaces. |
| Lead restore repair | Invalid Lead restore is present as `/leads/{id}/restore-pending`, maps to `PIM-CMD-006`, `PSM-API-005`, `PSM-TX-003`, `SM-LEAD`, and `EVT-STATUS-CHANGED`. |
| Duplicate warning repair | Duplicate warning is represented by `PSM-DUP`, `/duplicate-warnings/check`, `DUPLICATE_WARNING_REQUIRED`, warning token acknowledgement, safe masking, and mutation interception. |
| Core traceability | `ACC-001` through `ACC-023` each map to CIM, PIM, PSM, test IDs, and pending G8/G12 placeholders. |
| No downgrade detected | No P0/P1 acceptance item was deleted, downgraded, or marked done without evidence. |

## P0/P1 Blockers

### AMDA-BLOCK-001: Money Representation And Authoritative Calculation Are Not Finalized In PSM

- Severity: **P0 Blocker**
- Affected Acceptance: `ACC-009`, `ACC-010`, `ACC-011`, `ACC-013`, `ACC-023`
- Architecture Sources:
  - `docs/architecture/frontend-backend-contract.md` requires the money DTO strategy to be finalized in PSM and forbids frontend floating-point authority.
  - `docs/architecture/data-design.md` requires payment and quote/contract amount rules to be enforced through durable constraints and transactions.
- Current PSM Evidence:
  - `modeling/PSM.md` lists `amount`, `due_amount`, `paid_amount`, amount indexes, payment locks, and report amount fields, but does not choose the platform representation for API DTOs, Go domain values, or PostgreSQL columns.

Why this blocks:

Quote acceptance, contract amount mismatch, payment no-overpayment, Won full
payment, and report sums all depend on exact monetary behavior. Without a
PSM-level decision, Backend may choose `numeric`, integer minor units, or an
inconsistent Go/TypeScript representation; Frontend may represent money
differently from backend; QA cannot define precise boundary tests.

Required repair:

- Add a PSM money strategy section that defines:
  - API DTO representation for all money fields.
  - Go domain/application type strategy.
  - PostgreSQL column type and precision/scale or minor-unit integer rule.
  - Rounding and comparison rules for quote, contract, payment, and report sums.
  - Mapping to `PSM-DB-006` through `PSM-DB-010`, `PSM-API-008` through `PSM-API-010`, `PSM-API-016`, `PSM-TX-008`, and `TM-009` through `TM-013` / `TM-023`.

### AMDA-BLOCK-002: PSM Lacks Resource-Level Authorization Policy And Scope Loader Mapping

- Severity: **P0 Blocker**
- Affected Acceptance: `ACC-001`, `ACC-002`, `ACC-014`, `ACC-015`, `ACC-018`, `ACC-020`, `ACC-021`, `ACC-022`, `ACC-023`
- Architecture Sources:
  - `docs/architecture/authz-architecture.md` requires PSM to specify policy function names, inputs for each resource, scope loader queries for Sales/Sales Manager/Administrator, endpoint-to-policy mapping, denied error mapping, stale session/role recheck mechanics, and last-Administrator mechanics.
  - `docs/architecture/module-boundaries.md` requires application services to call authorization before protected reads/mutations and requires import/export/report modules to authorize before selecting source records.
- Current PSM Evidence:
  - `modeling/PSM.md` defines `PSM-AUTHZ` and broad Required Policy text per API group, but it does not define resource-specific policy functions, scope loader mappings, or endpoint-to-policy rows for lead/company/contact/opportunity/quote/contract/payment/activity/note/task/history/log/import/export/report/reminder.

Why this blocks:

The architecture makes backend authorization the primary P0 safety boundary.
If the PSM only says "role/scope" at API-group level, Task Planner and Backend
must infer how each resource loads scope and which policy method each endpoint
calls. That creates direct risk for IDOR, unauthorized report/export
aggregation, hidden record leakage, and inconsistent Sales vs Manager behavior.

Required repair:

- Add PSM tables for:
  - endpoint/action to policy function mapping.
  - resource to Sales/Sales Manager/Administrator scope loader mapping.
  - read/list/export/report/reminder authorization-before-query mapping.
  - denied error mapping for safe `PERMISSION_DENIED`, safe `NOT_FOUND`, and `BUSINESS_RULE_BLOCKED`.
  - stale session/role recheck and last-Administrator transaction mapping.
- Ensure the mapping covers all P0/P1 API groups and aligns with
  `docs/architecture/authz-architecture.md` resource scope rules.

### AMDA-BLOCK-003: Deployment, Backup, Restore, And Retention Are Too Thin In PSM

- Severity: **P0 Blocker**
- Affected Acceptance: `ACC-016`, `ACC-017`, `ACC-022`
- Architecture Sources:
  - `docs/architecture/deployment-notes.md` defines the production topology, required environment variables, deployment units, migration strategy, backup schedule, restore rehearsal, runbooks, and production verification inputs.
  - `docs/architecture/integration-design.md` requires encrypted off-host backups, restore into isolated environment, verification of CRM data/history/logs/reports, and retained operational evidence.
  - `docs/architecture/data-design.md` defines retention expectations for CRM records, contracts/payments, history, operation logs, imports, exports, and backups.
- Current PSM Evidence:
  - `modeling/PSM.md` includes `PSM-INFRA`, `PSM-JOBS`, and `ARCH-ACC-007`, but does not map required environment variables, migration ordering, backup retention/encryption/checksum metadata, restore rehearsal evidence, or application retention obligations to PSM elements.

Why this blocks:

`ACC-017` is P0 and requires real operation with persistent data, backup, and
restore evidence. If PSM only names the infrastructure target, Task Planner and
Integration Owner must infer the concrete tasks and tests for deploy,
migration, backup, restore, and retention. That weakens architecture's
accepted production constraints.

Required repair:

- Add PSM infra/operability mapping for:
  - deployment units and required environment variables.
  - migration sequencing and pre-migration backup requirement.
  - backup schedule, encryption, checksum, storage, and retention.
  - restore rehearsal flow and evidence requirements.
  - retention mapping for core records, contracts/payments, history, operation logs, import/export metadata, and backups.
  - trace links to `TM-016`, `TM-017`, and `TM-022`.

## P2 Improvements

| ID | Improvement | Rationale |
|---|---|---|
| AMDA-P2-001 | Replace broad references such as `PSM-API-*` and `PSM-API entity groups` in traceability with concrete PSM IDs where practical. | Improves Task Planner precision and avoids over-broad task mapping. |
| AMDA-P2-002 | Add PRD/NFR/security references to the `Architecture Acceptance Inside PSM` table. | The acceptance standard recommends architecture acceptance map to PRD IDs, product acceptance, NFRs/security where relevant, PSM elements, tasks, and tests. |
| AMDA-P2-003 | Clarify the traceability row text `PSM rule API-007` for `ACC-004`. | `API-007` is an architecture API rule, not a PSM ID; use an explicit source reference to avoid ID ambiguity. |
| AMDA-P2-004 | Split `PSM-AUDIT` into explicit history and operation-log PSM rows or add sub-rows. | Architecture has separate `history` and `operationlog` modules; explicit PSM rows would help task decomposition. |
| AMDA-P2-005 | Add explicit OpenAPI artifact names for DTO groups beyond enums. | Helps Frontend/Backend/QA plan generated client, request/response schema, examples, and contract tests. |

## No-Downgrade Check

| Check | Result |
|---|---|
| P0/P1 acceptance rows preserved | Passed |
| P0/P1 status not weakened | Passed |
| Mock/static/non-persistent behavior introduced | Not found |
| Implementation started before G8 | Not found in reviewed modeling artifacts |
| Product acceptance overridden by architecture/modeling | Not found |

## Recommendation

**Do not enter G8 Task Planning yet.**

Return to Domain Modeling to repair the PSM and traceability gaps listed above.
After repair, Architecture should perform a focused re-review of:

1. Money representation and exact calculation mapping.
2. Resource-level authz policy/scope loader mapping.
3. Deployment/backup/restore/retention PSM mapping.

If those focused items pass and no new P0/P1 gap is introduced, Architecture
can recommend G8 Task Planning from its domain.

## Files Modified By This Review

- `archive/reviews/g7-modeling/architecture-mda-pre-task-review.md`
