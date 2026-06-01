# CRM System — G7 Delivery Plan (Acceptance-Driven, End-to-End)

## Document Control

- Project: CRM System
- Gate: G7 (Task Planning) — planning artifact only, no implementation code
- Role: Task Planner
- Date: 2026-06-01
- Authority (single source; nothing invented): `docs/product/acceptance-matrix.md`
  (ACC-001..023), `docs/product/business-capability-map.md` (CAP-001..012),
  `docs/product/decision-log.md` (DEC-001..022), `modeling/CIM.md`,
  `modeling/PIM.md`, `modeling/PSM.md`, `modeling/test-model.md`,
  `modeling/traceability-matrix.md`, and the architecture/security docs under
  `docs/architecture/` and `docs/security/`.

This package is the directly-executable delivery plan for the implementation
agent (Codex) at G9–G11. Codex does not plan, design, or decide: every task is
self-contained with concrete file paths, contracts, tests, and a definition of
done. If something is genuinely undecided, it is recorded as a Blocker — never
guessed.

The four artifacts:

| File | Purpose |
|---|---|
| `delivery/delivery-plan.md` | This file — how to read the plan, global conventions, stack, repo layout, test strategy, no-downgrade and traceability rules, phased sequence. |
| `delivery/tasks.md` | The task plan. TASK-001..TASK-040, each with all 17 schema fields. |
| `delivery/acceptance-task-map.md` | ACC-001..023 → TASK-IDs. Proves 23/23 coverage. |
| `delivery/task-dependencies.md` | Dependency DAG and build order. |

## How To Read A Task

Each task in `tasks.md` carries all 17 fields. Codex executes a task by:

1. Reading the **Reference docs** IDs (CIM/PIM/PSM/CONTRACT/api-spec/PM/SEC/ABUSE/PRIV/TEST).
2. Writing the **Automated tests** FIRST (TDD) so they FAIL, then implementing.
3. Creating/editing exactly the **File changes** listed.
4. Satisfying the **Definition of Done** and **No-downgrade items** (real DB,
   real permission, real history event — never mock/stub/TODO/static/in-memory).
5. Marking **Status** done only when the mapped **ACC** passes by the stated
   **Acceptance method** and tests are green.

## Committed Stack (DEC-021, DEC-022, ADR-ARCH-001)

- **Backend:** Go microservices, the accepted SVC-001..010 (PSM Service Mapping),
  one service per Go module under `services/<service-name>/`. Owner agent for all
  backend services is `backend-engineer` (PSM Service Mapping; SVC-ACC-002).
- **Database:** PostgreSQL, schema-per-service with a dedicated DB user per
  service (`data-design.md` Data Ownership Map). Migrations via `golang-migrate`
  under `services/<svc>/migrations/`. No service touches another service's schema
  (SVC-ACC-007; `module-boundaries.md` Forbidden dependency types).
- **Frontend:** React + TypeScript under `frontend/`, talking ONLY to
  gateway-bff (`frontend-backend-contract.md`; `module-boundaries.md` Gateway
  Boundary). Screens per `ui-spec.md` UI-001..017. Owner agent `frontend-engineer`.
- **Deployment / operations:** owner agent `infrastructure-ops` (recognized project
  role — `AGENTS.md`, `process/process-gap-register.md`; required reviewer at
  G5/G8/G11/G12) for runtime environment, deployment, and release-evidence tasks
  (TASK-039/040: deploy, HTTPS/TLS, security group, monitoring, operator access,
  off-server backup + restore), with `backend-engineer` doing the in-repo
  implementation. Release-evidence items are proven at G11 and audited at G12.
- **Tests:** Go `testing` + `testcontainers` (unit/integration against REAL
  PostgreSQL — no mocks for P0/P1 persistence, SVC-ACC-011, DEC-008); Playwright
  (E2E) under `e2e/`. Test types per `test-model.md` Test-type convention.
- **Orchestration:** Docker Compose (`docker-compose.yml`), reverse-proxy +
  HTTPS per `deployment-notes.md`.

## Committed Repo Layout (monorepo, DEC-022)

```
crm-system/
  docker-compose.yml                  # all services + PostgreSQL + reverse proxy (TASK-001/039)
  shared/
    contracts/                        # API DTO schemas, event schemas, error codes,
                                      # permission action constants, correlation helpers
                                      # (module-boundaries.md Shared Package Boundary —
                                      # NO domain methods/repos/business rules)
  services/
    gateway-bff/                      # SVC-001  (Go module; no DB)
    identity-authz/                   # SVC-002  (Go module; schema identity_authz)
    lead/                             # SVC-003  (schema lead)
    account/                          # SVC-004  (schema account)
    opportunity/                      # SVC-005  (schema opportunity)
    commercial/                       # SVC-006  (schema commercial)
    work/                             # SVC-007  (schema work)
    audit-history/                    # SVC-008  (schema audit_history)
    reporting/                        # SVC-009  (schema reporting)
    import-export/                    # SVC-010  (schema import_export)
  frontend/                           # React + TS; talks only to gateway-bff
  e2e/                                # Playwright E2E specs
```

Per-service internal layout (each `services/<svc>/`):

```
services/<svc>/
  go.mod                              # independent module
  cmd/server/main.go                  # process entry + health endpoint
  internal/
    handler/                          # HTTP API handlers (Command/Query)
    domain/                           # aggregates, state machines, invariants (PIM)
    repo/                             # PostgreSQL repository (owned schema only)
    event/                            # outbox writer + dispatcher (FLOW-006)
    authz/                            # permission-check client (SVC-002) + S2S token
  migrations/                         # golang-migrate SQL (NNNN_*.up.sql / .down.sql)
  internal/.../*_test.go              # unit + integration (testcontainers)
```

Service↔schema↔DB-user↔module mapping (PSM Data Ownership / `data-design.md`):

| SVC | Module dir | Schema | DB user |
|---|---|---|---|
| SVC-001 gateway-bff | `services/gateway-bff/` | (none) | (none) |
| SVC-002 identity-authz | `services/identity-authz/` | `identity_authz` | `crm_identity_authz_user` |
| SVC-003 lead | `services/lead/` | `lead` | `crm_lead_user` |
| SVC-004 account | `services/account/` | `account` | `crm_account_user` |
| SVC-005 opportunity | `services/opportunity/` | `opportunity` | `crm_opportunity_user` |
| SVC-006 commercial | `services/commercial/` | `commercial` | `crm_commercial_user` |
| SVC-007 work | `services/work/` | `work` | `crm_work_user` |
| SVC-008 audit-history | `services/audit-history/` | `audit_history` | `crm_audit_history_user` |
| SVC-009 reporting | `services/reporting/` | `reporting` | `crm_reporting_user` |
| SVC-010 import-export | `services/import-export/` | `import_export` | `crm_import_export_user` |

## Global Conventions (apply to every task)

- **Business APIs only**, not generic CRUD (`api-spec.md` API Strategy).
- **Common request fields** on every protected call: `correlationId`, actor id,
  actor role, actor active/disabled status, authz/session version, caller service
  id for internal calls, target action, target resource type/id, idempotency key
  for non-idempotent writes, `expectedVersion` for editable P0/P1 records
  (`api-spec.md` Common Request Requirements).
- **Common response envelope** + **Error Codes** taxonomy from `api-spec.md`
  (`PERMISSION_DENIED`, `SCOPE_DENIED`, `VERSION_CONFLICT`,
  `INVALID_TRANSITION`, `OVERPAYMENT_BLOCKED`, `EARLY_WON_BLOCKED`,
  `ARCHIVE_BLOCKED_ACTIVE_OBLIGATION`, `TERMINAL_RECORD_READ_ONLY`,
  `LOST_REASON_REQUIRED`, `SERVICE_AUTH_FAILED`, `DUPLICATE_WARNING`, …). Errors
  never expose restricted names/amounts/contact details/existence of unauthorized
  records (`ui-spec.md` Data Display Safety; ABUSE-014).
- **Server-side authorization** via SVC-002 `POST /internal/permissions/check`
  on every protected action (CONTRACT-001, ARCH-ACC-002). Frontend hiding/disabling
  NEVER satisfies authorization (`permission-matrix.md` Matrix Rules).
- **Service-to-service** calls use a signed service token:
  `Authorization: Bearer <service-token>`, `X-Service-Id`, `X-Correlation-Id`,
  `X-Intent`, audience, max 5-min lifetime; missing/expired/wrong-audience/
  disallowed-intent rejected with `SERVICE_AUTH_FAILED` (CONTRACT-019,
  `authz-architecture.md` S2S, ARCH-ACC-009, SVC-ACC-008). No cross-service DB
  access and no cross-service internal imports (SVC-ACC-006/007).
- **Optimistic concurrency:** every editable P0/P1 record DTO carries numeric
  `version` + `updatedAt`; mutations require `expectedVersion`; stale edits
  return `VERSION_CONFLICT` (CONTRACT-020, ARCH-ACC-010).
- **History / operation-log on sensitive mutation:** every sensitive write emits
  a record-local history event and/or operation-log event through the SVC-008
  trusted-internal append contract, delivered via the producing service's
  database **outbox table + background dispatcher** in the same durable workflow
  (FLOW-006, CONTRACT-013/014, ARCH-ACC-003, AUD-IMM-002). Event schema fields
  from `data-design.md` History And Audit Data (`eventId`, `eventVersion`,
  `producerService`, `aggregateType`, `aggregateId`, `actorId`, `occurredAt`,
  `correlationId`, `causationId`, `safeSummary`, `prevHash`, `eventHash`); event
  names from `audit-log-spec.md` EVT-* catalog.
- **No hard delete** of any core CRM record by any role (PIM-INV-040, PM-029,
  DEC-011, ABUSE-017). Eligible records are archived (PIM-SM-010).
- **Overdue / reminder evaluation is on-read** against a supplied `businessDate`
  in `Asia/Shanghai` (PSM Resolved Mechanisms; `api-spec.md` Reminder Query).
  This makes overdue tests deterministic.
- **Single currency (CNY)**, no tax/discount/multi-currency (PIM-INV-024, DEC-013).

## Test Strategy & TDD (DEC-008, SVC-ACC-011, test-model.md)

- TDD is mandatory: for each task the listed `TEST-*` cases are written first and
  must FAIL before implementation, then pass. No test is deleted or weakened to go
  green (No-Downgrade; `test-model.md`).
- **Unit:** single-aggregate guard/invariant logic inside one service.
- **Integration:** service API + REAL PostgreSQL (testcontainers) + permission
  check + emitted history/oplog event + cross-service flow. Every P0 abuse and
  permission negative case has a backend/API negative test, not UI-only
  (`abuse-cases.md`, SEC-017, ARCH-ACC-002).
- **E2E:** Playwright across gateway-bff to persisted result and back.
- **Manual:** deployment smoke and operational evidence (ACC-017).
- The `TEST-*` IDs are exactly those defined in `modeling/test-model.md`; tasks
  cite the family/cases each must implement.

## No-Downgrade Rule (operating-model.md, DEC-008)

P0/P1 acceptance items are only `Done`, `Blocked`, or `Formal Scope Change by
User`. No P0/P1 item may be satisfied by mock, stub, TODO, static UI, or
non-persistent behavior. A "fake core" must not pass (SVC-ACC-011). Persistence is
proven against real PostgreSQL and verified after refresh, re-login, and service
restart (ACC-016, ARCH-ACC-015).

## Traceability Rule

Every task carries an explicit MDA chain
`CIM- → PIM- → PSM- → CONTRACT- → ACC- → TEST-` so each code change is traceable
to the model. `delivery/acceptance-task-map.md` proves all 23 ACC items map to
≥1 task; the PSM `PSM Traceability` `Task ID = pending (G8)` placeholders resolve
to the TASK-IDs in this package.

## Build-To-Current-Model Reminders (post-2026-06-01)

- Won = related contract **Signed** (NOT full payment) — DEC-017, PIM-INV-007/035.
- Exactly **one quote per opportunity** — DEC-018, PIM-INV-012.
- Payment tracking retained but **decoupled from Won** (post-sale follow-up) —
  DEC-019; overpayment still blocked (PIM-INV-023).
- Opportunity has **no Status field**; Pipeline Stage is the sole lifecycle
  dimension — DEC-020. Stages: New Opportunity, Needs Confirmed, Quote, Contract
  Negotiation, Won, Lost (Won/Lost terminal). No `Payment In Progress` /
  `Contract Signed` stage.
- Do NOT plan any task around a removed stage, a second quote, an opportunity
  status field, or a full-payment-gates-Won rule.

## Phased Delivery Sequence

**Phase 0 — Foundation / platform (prerequisites for everything):**
TASK-001 (monorepo + Docker Compose + PostgreSQL + migrate scaffold + shared
contracts), TASK-002 (identity-authz auth/session — ACC-001), TASK-003
(identity-authz permission + three-role scope + S2S token — ACC-002), TASK-004
(audit-history append-only history + oplog + outbox — ACC-014/022 core),
TASK-005 (gateway-bff routing/correlation/safe-error — ACC-015 spine),
TASK-006 (frontend app shell + auth + sign-in screen UI-001/UI-002).

**Phase 1 — Capability vertical slices (one+ task per remaining ACC):**
Leads (TASK-007..009), Account/Contact (TASK-010..012), Opportunity
(TASK-013..016), Commercial: quote/contract/payment (TASK-017..023), Work:
activity/note/task (TASK-024..025), Reminders (TASK-026), Record history view
(TASK-027), Admin operation log view + user/role admin (TASK-028..029), Core
retrieval list/detail/search/filter (TASK-030), Duplicate warning (TASK-031),
Archive lifecycle (TASK-032), Team overview (TASK-033), Basic reports (TASK-034),
Import (TASK-035), Export (TASK-036), Persistence verification (TASK-037),
Data classification & retention (TASK-038).

**Phase 2 — Deployment + release evidence (G11/G12 evidence, ACC-017):**
TASK-039 (deploy on runtime host via Docker Compose, reverse proxy, HTTPS/TLS,
security group, health/monitoring), TASK-040 (encrypted off-server backup +
restore rehearsal). These are explicitly flagged release-evidence tasks; their
ARCH-ACC items are `Release-evidence pending` and are proven at G11 and audited
at G12.

Within a phase, follow the DAG in `delivery/task-dependencies.md`: a slice's
backend command/persistence task precedes its frontend task, and both depend on
the foundation (auth + permission + the owning service scaffold + audit-history).
