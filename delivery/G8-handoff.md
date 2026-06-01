# G8 Execution Handoff — CRM System (Claude → Codex)

## Document Control

- Project: CRM System
- Gate: G8 — Task Planning → Implementation **[HANDOFF: Claude → Codex]**
- Prepared by: Claude (planning platform), 2026-06-01
- Executor: Codex (execution platform) — G9 Implementation, G10 QA, G11 Integration
- Returns to: Claude for **G12 independent audit** after G11 passes
- This is the single self-contained entry document for execution. Read it first.

## 1. Platform handoff (what this is)

Planning (G1–G7) is complete and gate-passed. This package hands execution to
Codex. **Claude writes no implementation code.** Codex implements the committed
plan, runs QA, and produces integration evidence; then Claude resumes for the
independent G12 audit before any release decision.

**Codex MUST NOT:** re-plan, re-design, re-decide scope, introduce out-of-scope
capability, change committed decisions (decision-log DEC-001..022), or downgrade
P0/P1. If something is genuinely undecided or blocked, STOP and record a blocker —
do not guess. (The plan is built to be zero-TBD; you should not hit this.)

## 2. What to read, in order

1. `delivery/delivery-plan.md` — committed stack, repo layout, global conventions, test strategy, phasing.
2. `delivery/tasks.md` — the 40 tasks (TASK-001..040), each self-contained with a 17-field schema.
3. `delivery/task-dependencies.md` — the build order / acyclic DAG. Execute respecting prerequisites.
4. `delivery/acceptance-task-map.md` — ACC-001..023 → tasks (prove full coverage as you complete tasks).

Authority (the model/contracts you build to — do not re-derive, reference by ID):
- MDA: `modeling/CIM.md`, `modeling/PIM.md`, `modeling/PSM.md`, `modeling/traceability-matrix.md`, `modeling/test-model.md`
- Contracts/architecture: `docs/architecture/*` (api-spec, frontend-backend-contract, data-design, module-boundaries, integration-design, authz-architecture, deployment-notes, service-architecture-adr)
- Acceptance/decisions: `docs/product/acceptance-matrix.md` (ACC-001..023), `docs/product/decision-log.md` (DEC-001..022)
- Security: `docs/security/*` (permission-matrix, abuse-cases, security-requirements, privacy-requirements, audit-log-spec)

## 3. Committed stack & repo layout (DEC-021/022, ADR-ARCH-001)

- Backend: Go microservices, SVC-001..010, one Go module per service under `services/<service-name>/`; owner `backend-engineer`.
- Database: PostgreSQL, schema-per-service with a dedicated DB user per service (`data-design.md` Data Ownership Map); migrations via `golang-migrate` under `services/<svc>/migrations/`. No service touches another service's schema.
- Frontend: React + TypeScript under `frontend/`, talking ONLY to gateway-bff; owner `frontend-engineer`. Screens per `docs/ux-ui/ui-spec.md` UI-001..017.
- Tests: Go `testing` + `testcontainers` (unit/integration against REAL PostgreSQL — no mocks for P0/P1 persistence); Playwright (E2E) under `e2e/`.
- Orchestration: Docker Compose; single runtime host `srv-volcengine-sh-01` (Volcengine ECS); off-server backup target `srv-aliyun-bj-01`. Deployment/ops owner `infrastructure-ops`.

## 4. Per-task execution contract

For each task, honor all 17 fields. In particular:
- **TDD fail-first**: write the named `TEST-*` tests first (they must FAIL), then implement until green. Never delete or weaken a test to pass.
- **Real persistence / no fakes**: no mock/stub/TODO/static/non-persistent behavior may satisfy a P0/P1 item (SVC-ACC-011, DEC-008). Integration tests run against real PostgreSQL.
- **Service boundaries**: a service accesses only its own schema; cross-service interaction is via Command/Query APIs or events with the S2S signed token (authz-architecture); never a cross-service DB read/write.
- **MDA traceability**: keep the CIM→PIM→PSM→CONTRACT→ACC→TEST chain intact; cite it in code/PR as the task does.
- **Status**: update each task's Status (Not Started → In Progress → Done / Blocked) as you go. P0/P1 valid states are only `Done`, `Blocked`, or `Formal Scope Change by User`.
- Fill the `traceability-matrix.md` Task/Test/Integration columns (currently `pending`) as tasks/tests/evidence are produced.

## 5. Build to the current model (post-2026-06-01 scope change, DEC-017..020)

- Opportunity **Won = related contract Signed** (NOT full payment). Verify the related contract is Signed (via commercial-service / S2S) before persisting Won.
- **Exactly one quote per opportunity** (DB UNIQUE + domain guard); no multi-quote / second-accept path.
- **Payment tracking retained but decoupled from Won** (plans/actual/status/overdue reminders/reports stay; overpayment blocked; single currency); payment never gates Won.
- Opportunity has **no Status field** — Pipeline Stage is the sole lifecycle dimension (New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost).
- **Do NOT implement** retired tests: `TEST-OPP-STATUS-ENUM-001`, `TEST-QUOTE-ACCEPT-003`, `TEST-PAYMENT-FULLPAID-AGG-001`.

## 6. Gate expectations for the execution platform

- **G9 Implementation**: implement all tasks in dependency order; every task's automated tests green; no-downgrade honored; no P0/P1 faked.
- **G10 QA Execution**: execute the test model; produce QA/defect reports; regression.
- **G11 Integration**: end-to-end evidence for every P0/P1 service-backed flow (ACC-001..023), including the carried release evidence (below). Provide evidence with ACC id, environment, steps, result, service chain, correlation ID.
- Then **return to Claude for G12 independent audit** before any release decision. Do not skip gates or bypass reviewers.

## 7. Standing conditions, carried blockers, real-world dependencies

- **Operator access** (pre-G8 condition): Security Compliance reviewed and **approved with conditions** (`archive/reviews/g8-handoff/security-operator-access-review-2026-06-01.md`). G9 must create the named least-privilege deploy/ops user + document SSH key ownership/sudo boundary (TASK-039); evidenced at G11, re-audited at G12.
- **Carried release blockers** (G11/G12 release evidence, ARCH-ACC-004/008/013/014/015): encrypted off-server backup + restore rehearsal (TASK-040); HTTPS/TLS endpoint; security-group / network exposure; monitoring/health; restart-survival persistence.
- **Real-world dependency**: a valid HTTPS endpoint / TLS certificate (and domain, if used) must be provided for production release (TASK-039 blocker) — do not invent a domain; record if unavailable.
- **No open product decisions**: the three former G6 blockers (BLK-001/002/003) are resolved by DEC-017..020; there are no undecided P0/P1 product rules for Codex to resolve.

## 8. Governance Codex must follow

- No-Downgrade: P0/P1 cannot be downgraded, deleted, weakened, or merged; valid states `Done` / `Blocked` / `Formal Scope Change by User`.
- No mock/stub/TODO/static/non-persistent behavior for P0/P1.
- Service-boundary-first; every cross-capability flow has a named Primary Flow Owner Agent (per `process/process-gap-register.md`).
- Infrastructure Ops is a required reviewer at G8/G11/G12.
- On a kickback, set the affected gate to `Gate Blocked` in `planning/gate-status.md`, record it in `planning/blockers.md`, and return to Claude.

## 9. Handoff package contents (all on disk)

- `delivery/delivery-plan.md`, `delivery/tasks.md`, `delivery/task-dependencies.md`, `delivery/acceptance-task-map.md`, `delivery/G8-handoff.md` (this file)
- The full authority set under `modeling/`, `docs/`, and the registers under `planning/` + `process/`
- Gate evidence under `archive/reviews/` (incl. the operator-access review)
