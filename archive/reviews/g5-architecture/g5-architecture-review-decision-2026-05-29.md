# G5 Architecture Review Decision

## Document Control

- Project: CRM System
- Review Date: 2026-05-29
- Gate: G5 Architecture Design
- Gate Owner: Architecture
- Required Reviewers: Product Manager, Business Analyst, UX Designer, UI Designer, Security Compliance
- Additional Reviewer: Infrastructure Ops
- Decision: Block
- Archive Note: This file records G5 review evidence. It is not design authority and does not replace active architecture documents.

## Review Scope

Reviewed current active architecture draft:

- `docs/architecture/architecture.md`
- `docs/architecture/module-boundaries.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`
- `docs/architecture/service-architecture-adr.md`
- `docs/architecture/service-architecture-acceptance.md`
- `docs/architecture/service-acceptance-map.md`

Reviewers also checked their upstream source documents under:

- `docs/product/`
- `docs/business/`
- `docs/ux-ui/`
- `docs/security/`
- `../../company/infrastructure/`

No implementation code was written. No P0/P1 item was downgraded, deleted,
merged away, weakened, or accepted as partial work.

## Gate Decision

G5 does not pass.

All participating reviewers returned `Block`:

| Reviewer | Decision | Summary |
|---|---|---|
| Product Manager | Block | ACC-017/OQ-001 endpoint and ownership gaps; ACC-019 duplicate warning contract gap. |
| Business Analyst | Block | Owner transfer, Lost closure, archive obligations, duplicate normalization, report metrics, reminder date rules, import/export object scope gaps. |
| UX Designer | Block | Archive eligibility UX, duplicate warning UX, and stale/conflict UX contracts incomplete. |
| UI Designer | Block | Archive obligation UI, stale edit conflict, report DTO, duplicate proceed-after-warning, export metadata, and reminder row DTO gaps. |
| Security Compliance | Block | Service-to-service authentication, session invalidation, backup security, HTTPS/TLS, audit tamper resistance, and CSV security gaps. |
| Infrastructure Ops | Block | Same-host-only backup conflicts with company production backup requirement; security group record, monitoring/alerting, capacity, runtime paths, and root-ops plan gaps. |

## Consolidated P0 Blockers

| ID | Owner | Blocker | Required Fix |
|---|---|---|---|
| G5-BLK-001 | Architecture, Infrastructure Ops | ACC-017/OQ-001 is not fully closed for G5 because v1 access endpoint strategy, environment ownership, security group evidence, monitoring target, and deployment evidence ownership are incomplete. | Define IP/domain strategy, endpoint evidence fields, owner split across Architecture / Infrastructure Ops / Backend, security group record requirement, monitoring target, and release evidence path. |
| G5-BLK-002 | Architecture, Infrastructure Ops | Same-host-only local PostgreSQL backup conflicts with company infrastructure P0 expectation that production backups must not exist only on the production server. | Either define an external backup target or explicitly mark same-host-only backup as a release-blocking gap that must close before production release. |
| G5-BLK-003 | Architecture, Security Compliance | Service-to-service authentication mechanism is not concrete. Docker internal network is not a trust boundary. | Define internal service authentication, credential storage, rotation, rejection behavior, and caller verification for internal APIs. |
| G5-BLK-004 | Architecture, Security Compliance | Session/token strategy and role/status change invalidation are not concrete enough for ACC-001/ACC-002. | Define session/token storage, expiration, logout, disabled-user handling, and role/status re-evaluation or invalidation behavior. |
| G5-BLK-005 | Architecture, Security Compliance | Backup files contain restricted/security-critical data but backup access control, encryption decision, permissions, and restore data protection are not defined. | Define backup directory permissions, encryption decision, allowed access subject, secret handling, restore logging, and restore privacy controls. |
| G5-BLK-006 | Architecture, Security Compliance | Production public entry allows HTTP/HTTPS wording and does not force HTTPS/TLS for login/session traffic. | Define HTTPS-only production access, HTTP-to-HTTPS redirect, TLS certificate handling, secure session transport, and reverse proxy security headers. |
| G5-BLK-007 | Architecture, Business Analyst, UX Designer, UI Designer | Archive active downstream obligation behavior is not architected. | Define archive eligibility API, active obligation DTO, blocked response, retry/refresh behavior, and history/operation-log events. |
| G5-BLK-008 | Architecture, UI Designer, UX Designer | P0 editable-record stale/conflict handling lacks a concrete concurrency token contract. | Define version/revision/updatedAt/etag or equivalent for P0 editable records, expected-version command fields, and conflict recovery response. |
| G5-BLK-009 | Architecture, Business Analyst | Owner transfer rule for open tasks/follow-ups is not fully represented. | Define OwnerChanged event/command, work-service transfer behavior, manual reassignment exception, failure recovery, and history/log events. |
| G5-BLK-010 | Architecture, Business Analyst | Lost closure and post-close editability are incomplete. | Define CloseLost contract, required lostReason, terminal edit protection, post-close note/task allowed path, and audit events. |

## Consolidated P1 Blockers / Required Fixes

| ID | Owner | Issue | Required Fix |
|---|---|---|---|
| G5-ISS-001 | Architecture, Product Manager, BA, UX/UI | Duplicate warning is too abstract for ACC-019. | Define duplicate normalization, cross lead/account/contact lookup, safe match summary, warning response, proceed-after-warning or equivalent mechanism, and no merge/no overwrite behavior. |
| G5-ISS-002 | Architecture, BA, UI | Report metric contract is not concrete. | Define report response DTO, groupings, amount fields, zero state, active/default archive filter, and authorization-before-aggregation behavior. |
| G5-ISS-003 | Architecture, BA, UI | Reminder eligibility and row DTO are incomplete. | Define workspace timezone/business date, due/overdue rules for tasks, pending-signature contracts, and payments, plus reminder row DTO. |
| G5-ISS-004 | Architecture, BA, Security, UI | Import/export contracts are incomplete. | Define v1 object scope, target service routing, row validation, export preview/metadata DTO, dangerous CSV formula handling, temporary file retention/deletion, authorization, and operation log. |
| G5-ISS-005 | Architecture, Security | Audit/history append-only or tamper-evident design is insufficient. | Define service/API no-update/no-delete rules, database permission constraints, admin UI prohibition, and optional hash/chaining or operation-level tamper evidence. |
| G5-ISS-006 | Architecture, Infrastructure Ops | ECS capacity and runtime path planning are incomplete. | Define minimum resource expectations, memory/disk thresholds, upgrade triggers, `/opt/crm-system`, volume, log, backup, and deployment paths. |
| G5-ISS-007 | Infrastructure Ops, Security Compliance | Long-term root SSH operation remains unresolved. | Define deploy/ops user plan, root usage boundary, and security review timing. |

## Positive Coverage Confirmed

Reviewers agreed the draft already covers these directions:

- Physical multi-service Go backend with Docker isolation.
- Service-boundary-first design and DDD-aligned capability boundaries.
- No unified database CRUD service.
- One PostgreSQL instance with service-owned database/schema and independent service credentials.
- No direct cross-service database access.
- Gateway/BFF is not the business owner.
- Frontend hiding/disabled actions are not authorization.
- Reporting must not directly query source service tables.
- Import/export must not bypass domain services.
- History/audit must be durable and read-only through normal product workflows.
- G5 draft is explicitly not implementation approval.

## Required Repair Sequence

1. Architecture repairs P0 blockers in architecture documents.
2. Infrastructure Ops supplies or requests required server/security-group/monitoring/backup evidence inputs.
3. Security Compliance reviews service authentication, session/token, HTTPS, backup security, audit tamper resistance, and CSV security fixes.
4. Product Manager reviews ACC-017/OQ-001 and ACC-019 fixes.
5. Business Analyst reviews owner transfer, Lost closure, archive obligations, duplicate normalization, report/reminder/import-export rule coverage.
6. UX Designer and UI Designer review archive, duplicate, conflict, report, export, reminder, and UI-facing DTO fixes.
7. Architecture produces a G5 repair note.
8. Required reviewers perform G5 re-review.

## Final Decision

G5 Architecture Design is `Gate Blocked`.

The project may not proceed to G6 MDA Modeling, G7 traceability/test model, G8
task planning, or implementation until G5 blockers are repaired and required
reviewers approve the architecture without open P0/P1 blockers.
