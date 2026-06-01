# PSM

Platform Specific Model — the platform-specific engineering design that reflects
the accepted G5 architecture and maps the PIM domain objects, behaviors,
state machines, and invariants onto the committed service decomposition, service
contracts, data ownership, cross-service reliability, and deployment boundary. It
reflects and traces the accepted architecture; it does not re-decide, renumber,
or invent it.

## Document Control

- Project: CRM System
- Date: 2026-06-01
- Role: Domain Modeling
- Gate: G6 (MDA Modeling)
- Scope note: This PSM reflects the accepted G5 architecture set (service
  decomposition SVC-001..010, flows INT-FLOW-001..013, service-architecture
  acceptance SVC-ACC-001..012, contract families, data design, authz, integration
  design, deployment notes, audit-log spec, privacy requirements). It traces the
  PIM (its direct upstream: PIM-001..028, PIM-SM-001..011, PIM-INV-001..052,
  PIM-BEH-001..034) and the CIM vocabulary. It invents no new service, contract,
  flow, event, or ID beyond what the architecture committed. Where a contract or
  schema is defined in an architecture doc, this PSM REFERENCES that doc's
  location rather than copying the full schema (drift avoidance). Open
  product/BA decisions BLK-001/002/003 were RESOLVED 2026-06-01 by a Formal Scope
  Change by User (decision-log.md DEC-017..020) and are no longer open. Amended
  2026-06-01: Won = related contract Signed (not full payment); exactly one quote
  per opportunity; payment tracking retained but decoupled from Won; Opportunity
  Status removed (Pipeline Stage is the sole lifecycle dimension); `Payment In
  Progress` / `Contract Signed` pipeline stages removed. Architecture-deferred
  mechanisms that the architecture actually committed are resolved here (overdue
  trigger, sync/async per flow, outbox-vs-replay, import/export temp-file cleanup);
  any genuinely architecture-pending input is recorded under `## Open / Blocked`.

## Tier Altitude Statement

This document stays at platform-specific engineering altitude and is in-tier for:
service mapping, API/event/error/permission contract mapping, data ownership +
forbidden access + consistency strategy, cross-service flow reliability
(idempotency/retry/timeout/compensation/correlation ID), architecture acceptance,
deployment boundary, and PSM traceability. It does NOT re-derive CIM vocabulary,
re-state PIM state machines, or re-decide architecture; it maps FROM PIM IDs and
REFERENCES architecture-doc locations for full schemas. Task IDs, test IDs, and
integration evidence are forward placeholders produced at G7/G8/G11 and are marked
pending, not invented here.

## Technical Models

Per service/aggregate, the platform-specific element: owning service, persistence
(database/schema per `data-design.md` Data Ownership Map), API surface (per
`api-spec.md` Service API Summary), and the key DTO/read-model (per
`api-spec.md` / `frontend-backend-contract.md`). Maps From PIM IDs; Acceptance ID.

| ID | Type | Name | Service | Persistence (schema) | API Surface (api-spec.md) | Key DTO / Read Model | Maps From PIM | Acceptance ID |
|---|---|---|---|---|---|---|---|---|
| PSM-001 | Service / DB / API / Auth | User & Session & Permission model | SVC-002 | `identity_authz` (users, roles, sessions, permission policy) | sign in/out, create/update user, change role/status, revoke sessions; current user, session/permission check | Current-user DTO, permission-check response (api-spec.md "Permission Check"), opaque secure session cookie (authz-architecture.md) | PIM-001, PIM-002, PIM-003 | ACC-001, ACC-002, ACC-022 |
| PSM-002 | Service / DB / API | Lead aggregate model | SVC-003 | `lead` (leads, owner/status, qualification, conversion marker) | create/update/assign/qualify/restore/convert lead, proceed-after-warning; list/detail/search, duplicate signal, archive eligibility | Lead detail DTO (`version`,`updatedAt`), Lead Conversion DTO (api-spec.md "Lead Conversion") | PIM-004, PIM-003, PIM-021 | ACC-003, ACC-004, ACC-019 |
| PSM-003 | Service / DB / API | Account & Contact aggregate model | SVC-004 | `account` (accounts/customers, contacts, ownership/status) | create/update/archive account, create/update/link contact, proceed-after-warning; list/detail/search/summary, duplicate, archive eligibility | Account/Contact detail DTO (`version`,`updatedAt`), account/contact safe summary | PIM-005, PIM-006, PIM-003, PIM-021 | ACC-005, ACC-006, ACC-019 |
| PSM-004 | Service / DB / API | Opportunity aggregate model | SVC-005 | `opportunity` (opportunities, Pipeline Stage, closure data — no separate Status dimension, DEC-020) | create/update opportunity, change stage, close won (related contract Signed, DEC-017)/lost, archive, change owner; list/detail/search, summary, closure/archive eligibility | Opportunity detail DTO (`version`,`updatedAt`), Close-Lost/Close-Won DTO (api-spec.md), opportunity summary | PIM-007, PIM-012, PIM-003 | ACC-007, ACC-008, ACC-013 |
| PSM-005 | Service / DB / API | Quote aggregate model | SVC-006 | `commercial` (quotes — exactly one per opportunity, DEC-018) | create/update quote, change quote status; quote list/detail | Quote detail DTO (`version`,`updatedAt`), quote status summary | PIM-008, PIM-012 | ACC-009 |
| PSM-006 | Service / DB / API | Contract aggregate model | SVC-006 | `commercial` (contracts) | create/update contract, change contract status (Signed drives Opportunity Won, DEC-017); contract detail, reminder eligibility | Contract detail DTO (`version`,`updatedAt`), accepted-quote constraint | PIM-009, PIM-012 | ACC-010 |
| PSM-007 | Service / DB / API | Payment Plan & Actual Payment model | SVC-006 | `commercial` (payment plans, actual payments, payment status) | create payment plan, record payment; payment list/detail, payment status summary. Payment tracking retained as post-sale follow-up, decoupled from Won (DEC-019) | Record-Payment DTO (api-spec.md "Record Payment": `paymentStatus`,`remainingAmount`) | PIM-010, PIM-011, PIM-012 | ACC-011, ACC-013 |
| PSM-008 | Service / DB / API | Activity / Note / Task / Reminder model | SVC-007 | `work` (activities, notes, tasks, reminder projection) | create activity/note/task, update task status, transfer open work; activity/task list, reminder list, active obligation list | Reminder row DTO (api-spec.md "Reminder Query"), task detail DTO (`version`,`updatedAt`) | PIM-013, PIM-014, PIM-015, PIM-016, PIM-017 | ACC-012, ACC-021 |
| PSM-009 | Service / DB / API | History & Operation-Log event model | SVC-008 | `audit_history` (record-local history, operation logs, security/audit events) | append history/log (trusted internal), record history query, admin operation-log query | History/Log event row (data-design.md fields: `eventId`,`eventVersion`,`producerService`,`aggregateType`,`aggregateId`,`actorId`,`occurredAt`,`correlationId`,`causationId`,`safeSummary`,`prevHash`,`eventHash`); event catalog EVT-* (audit-log-spec.md) | PIM-018, PIM-019, PIM-027 | ACC-014, ACC-022 |
| PSM-010 | Service / DB / API | Report & Team-Overview read model | SVC-009 | `reporting` (read-model projections, report snapshots, overview metrics) | rebuild projection, refresh snapshot; team overview, basic sales reports | Sales-overview DTO + `GroupRow`/`breakdowns` (api-spec.md "Report Metrics") | PIM-024, PIM-025, PIM-020 | ACC-018, ACC-023 |
| PSM-011 | Service / DB / API | Core Retrieval / BFF aggregation model | SVC-001 | none (not a data owner; aggregates public Query APIs) | list/detail/search/filter aggregation, correlation-ID propagation, safe error normalization | Authorized list/detail DTO, invalid-filter error (frontend-backend-contract.md) | PIM-026, PIM-020 | ACC-015 |
| PSM-012 | Service / DB / API | Import / Export run model | SVC-010 | `import_export` (import/export runs, row results, export metadata) | start import/export, cancel safe; run status, row results, export metadata | Import/Export run DTO, row-error DTO, export metadata (frontend-backend-contract.md "Import/Export UX Contract") | PIM-022, PIM-023 | ACC-020 |
| PSM-013 | Cross-cutting model | Data Classification & Retention model | SVC-002..010 (per owned data) | per-service owned schema + `audit_history`; backups per `data-design.md` Backup | classification carried on record/log DTOs; retention enforced at storage, no hard-delete API | Classification tag (`diff_classification`), retention metadata; durations per privacy-requirements.md Retention Policy | PIM-027, PIM-028, PIM-INV-049..052 | ACC-014, ACC-016, ACC-022 |
| PSM-014 | Deployment model | Runtime deployment boundary | SVC-001..010 + PostgreSQL + reverse proxy + backup job | Docker Compose on `srv-volcengine-sh-01`; off-server backup `srv-aliyun-bj-01` | HTTPS ingress via reverse proxy; per-service health endpoint; backup/restore job (deployment-notes.md) | Release evidence record (endpoint, TLS, security group, health, backup/restore) | (PIM-OPEN-004 / deployment) | ACC-016, ACC-017 |

Cross-cutting note: PIM-020 (Archive State) applies to every eligible record
aggregate (Lead PSM-002, Account/Contact PSM-003, Opportunity PSM-004, Contract
PSM-006). It is homed for view-exclusion in PSM-010 (reports) and PSM-011
(retrieval) and is exercised on those record aggregates via FLOW-010 (archive
eligibility); it is not re-listed in each record row's Maps-From to avoid
duplication.

## Service Mapping

Reuses the accepted SVC-001..010 from `service-architecture-adr.md` /
`module-boundaries.md` / `service-acceptance-map.md`. Owner agent =
`backend-engineer` per the service-acceptance-map (SVC-001 gateway-bff also
`backend-engineer`). Bounded context = `module-boundaries.md` Responsibilities.
Deployment boundary: all services are independent Docker containers co-located on
one runtime host via Docker Compose (ADR-ARCH-001).

| Service ID | Service | Service Owner Agent | Business Capability / Bounded Context | PIM IDs | PSM IDs | Acceptance IDs | Deployment Boundary |
|---|---|---|---|---|---|---|---|
| SVC-001 | gateway-bff | backend-engineer | External API entry, frontend aggregation, routing, correlation-ID propagation, safe error normalization; not a data owner | PIM-026, PIM-020 | PSM-011 | ACC-015 (+ edge for ACC-001..023) | Independent container; HTTPS ingress via reverse proxy (ADR-ARCH-001/005) |
| SVC-002 | identity-authz-service | backend-engineer | Authentication, sessions, users, roles, active/disabled state, permission decisions, last-admin protection | PIM-001, PIM-002, PIM-003 | PSM-001, PSM-013 | ACC-001, ACC-002, ACC-022 | Independent container; owns `identity_authz` |
| SVC-003 | lead-service | backend-engineer | Lead lifecycle, owner assignment, qualification, conversion-once guard, lead duplicate input | PIM-004, PIM-003, PIM-021 | PSM-002 | ACC-003, ACC-004, ACC-019 | Independent container; owns `lead` |
| SVC-004 | account-service | backend-engineer | Company/customer + contact lifecycle, ownership, account/contact duplicate input | PIM-005, PIM-006, PIM-003, PIM-021 | PSM-003 | ACC-005, ACC-006, ACC-019 | Independent container; owns `account` |
| SVC-005 | opportunity-service | backend-engineer | Opportunity lifecycle, stage transitions, Won/Lost terminal closure | PIM-007, PIM-012, PIM-003 | PSM-004 | ACC-007, ACC-008, ACC-013 | Independent container; owns `opportunity` |
| SVC-006 | commercial-service | backend-engineer | Quote, contract, payment plan, actual payment, payment status, amount integrity | PIM-008, PIM-009, PIM-010, PIM-011, PIM-012 | PSM-005, PSM-006, PSM-007 | ACC-009, ACC-010, ACC-011, ACC-013 | Independent container; owns `commercial` |
| SVC-007 | work-service | backend-engineer | Activities, notes, follow-up tasks, task status, due/overdue reminder eligibility | PIM-013, PIM-014, PIM-015, PIM-016, PIM-017 | PSM-008 | ACC-012, ACC-021 | Independent container; owns `work` |
| SVC-008 | audit-history-service | backend-engineer | Record-local history, admin operation logs, append-only event storage and query | PIM-018, PIM-019, PIM-027 | PSM-009 | ACC-014, ACC-022 | Independent container; owns `audit_history` |
| SVC-009 | reporting-service | backend-engineer | Team overview, basic sales reports, authorized read models/projections | PIM-024, PIM-025, PIM-020 | PSM-010 | ACC-018, ACC-023 | Independent container; owns `reporting` |
| SVC-010 | import-export-service | backend-engineer | CSV import/export run tracking, row results, export metadata, long-running state | PIM-022, PIM-023 | PSM-012 | ACC-020 | Independent container; owns `import_export` |

## Service Contract Mapping

Contract IDs (CONTRACT-*) group the API/Event/Error/Permission contract families
per service. Schema/Contract Location is a POINTER to the authoritative
architecture doc (api-spec.md, integration-design.md, audit-log-spec.md,
authz-architecture.md, frontend-backend-contract.md) — full payloads are NOT
copied here. Permission rules cite `permission-matrix.md` (PM-*). Event names
reuse `api-spec.md` Service API Summary and the `audit-log-spec.md` EVT-* catalog.

| Contract ID | Service | Type | Producer | Consumers | Schema / Contract Location | Permission Rule | Acceptance IDs |
|---|---|---|---|---|---|---|---|
| CONTRACT-001 | SVC-002 | API/Permission/Error | SVC-002 | SVC-001, all domain services | api-spec.md "Permission Check"; authz-architecture.md Permission/Denial/Session contracts | PM-001, PM-002, PM-003..PM-007 | ACC-001, ACC-002 |
| CONTRACT-002 | SVC-002 | Event | SVC-002 | SVC-008 | api-spec.md (UserSignedIn/Out, UserAccessDenied, UserRoleStatusChanged, SessionRevoked); audit-log-spec.md EVT-AUTH-LOGIN-SUCCEEDED/FAILED, EVT-AUTH-ACCESS-DENIED, EVT-USER-ROLE-CHANGED, EVT-USER-STATUS-CHANGED, EVT-LAST-ADMIN-BLOCKED | PM-003, PM-006, PM-007, PM-040 (admin-only log) | ACC-001, ACC-002, ACC-022 |
| CONTRACT-003 | SVC-003 | API/Permission/Error | SVC-003 | SVC-001, SVC-002, SVC-010 | api-spec.md Service API Summary (lead) + "Lead Conversion","Duplicate Warning" | PM-010..PM-015, PM-048 | ACC-003, ACC-004, ACC-019 |
| CONTRACT-004 | SVC-003 | Event | SVC-003 | SVC-004, SVC-005, SVC-008, SVC-009 | api-spec.md (LeadCreated, LeadOwnerChanged, LeadQualified, LeadConverted, DuplicateWarningRaised); audit-log-spec.md EVT-OWNER-CHANGED, EVT-LEAD-QUALIFIED, EVT-LEAD-DISQUALIFIED, EVT-LEAD-CONVERTED | PM-011, PM-014, PM-024/PM-025 (history visibility) | ACC-003, ACC-004, ACC-014 |
| CONTRACT-005 | SVC-004 | API/Permission/Error | SVC-004 | SVC-001, SVC-002, SVC-003, SVC-005, SVC-006, SVC-007, SVC-010 | api-spec.md Service API Summary (account/contact) + "Archive Eligibility","Owner Transfer","Duplicate Warning" | PM-008, PM-009, PM-016, PM-017, PM-026..PM-028, PM-048 | ACC-005, ACC-006, ACC-019 |
| CONTRACT-006 | SVC-004 | Event | SVC-004 | SVC-008, SVC-009 | api-spec.md (AccountCreated, ContactCreated, AccountArchived, DuplicateWarningRaised, OwnerChanged); audit-log-spec.md EVT-OWNER-CHANGED, EVT-RECORD-ARCHIVED, EVT-STATUS-CHANGED | PM-016, PM-026/PM-027 | ACC-005, ACC-006, ACC-014, ACC-022 |
| CONTRACT-007 | SVC-005 | API/Permission/Error | SVC-005 | SVC-001, SVC-002, SVC-004, SVC-006, SVC-007 | api-spec.md Service API Summary (opportunity) + "Close Opportunity Won/Lost","Archive Eligibility","Owner Transfer","Editable Record Concurrency" | PM-018, PM-019 | ACC-007, ACC-008, ACC-013 |
| CONTRACT-008 | SVC-005 | Event | SVC-005 | SVC-006, SVC-008, SVC-009 | api-spec.md (OpportunityCreated, OpportunityStageChanged, OpportunityClosedWon, OpportunityClosedLost, OwnerChanged); audit-log-spec.md EVT-STAGE-CHANGED, EVT-OPPORTUNITY-WON, EVT-OPPORTUNITY-LOST, EVT-LEAD-CONVERTED (downstream), EVT-RECORD-ARCHIVED | PM-018 | ACC-007, ACC-008, ACC-013, ACC-014, ACC-022 |
| CONTRACT-009 | SVC-006 | API/Permission/Error | SVC-006 | SVC-001, SVC-002, SVC-005, SVC-007, SVC-010 | api-spec.md Service API Summary (commercial) + "Record Payment","Editable Record Concurrency" | PM-020, PM-021, PM-022 | ACC-009, ACC-010, ACC-011, ACC-013 |
| CONTRACT-010 | SVC-006 | Event | SVC-006 | SVC-005, SVC-008, SVC-009 | api-spec.md (QuoteAccepted, ContractStatusChanged, PaymentRecorded, PaymentOverdue); audit-log-spec.md EVT-QUOTE-ACCEPTED, EVT-CONTRACT-SIGNED, EVT-CONTRACT-TERMINATED, EVT-PAYMENT-RECORDED, EVT-PAYMENT-OVERDUE, EVT-STATUS-CHANGED | PM-020, PM-021, PM-022 | ACC-009, ACC-010, ACC-011, ACC-014, ACC-021, ACC-022 |
| CONTRACT-011 | SVC-007 | API/Permission/Error | SVC-007 | SVC-001, SVC-002, record-owning services, SVC-006 | api-spec.md Service API Summary (work) + "Reminder Query","Owner Transfer" | PM-023, PM-046, PM-047 | ACC-012, ACC-021 |
| CONTRACT-012 | SVC-007 | Event | SVC-007 | SVC-008, SVC-009 | api-spec.md (WorkItemCreated, TaskStatusChanged, ReminderStateChanged, OpenWorkTransferred); audit-log-spec.md EVT-TASK-COMPLETED, EVT-TASK-CANCELLED, EVT-STATUS-CHANGED | PM-023 | ACC-012, ACC-014, ACC-021 |
| CONTRACT-013 | SVC-008 | API/Permission/Error | SVC-008 | All mutation-producing services (write); SVC-001/SVC-002 (query) | api-spec.md Service API Summary (audit-history) + data-design.md "History And Audit Data"; audit-log-spec.md Common Event Schema + Query Requirements | PM-024, PM-025 (record history), PM-040..PM-042 (admin-only operation log), AUD-IMM-001..005 | ACC-014, ACC-022 |
| CONTRACT-014 | SVC-008 | Event | SVC-008 | (terminal sink; append-only store) | api-spec.md (HistoryEventAppended, OperationLogAppended); audit-log-spec.md Event Catalog EVT-* (full); data-design.md tamper-evidence fields | PM-029 (no hard delete), AUD-IMM-002 | ACC-014, ACC-016, ACC-022 |
| CONTRACT-015 | SVC-009 | API/Permission/Error | SVC-009 | SVC-001, SVC-002 | api-spec.md "Report Metrics" (`GroupRow`/`breakdowns`); data-design.md "Reporting Data" | PM-043, PM-044, PM-045; authz-before-aggregate | ACC-018, ACC-023 |
| CONTRACT-016 | SVC-009 | Event | source services → SVC-009 | SVC-009 | api-spec.md (ReportProjectionUpdated; consumes source domain events); audit-log-spec.md EVT-REPORT-ACCESS-DENIED | PM-045 | ACC-018, ACC-023 |
| CONTRACT-017 | SVC-010 | API/Permission/Error | SVC-010 | SVC-001, SVC-002, target domain services | api-spec.md "Import And Export Runs" + Operation Status Contract; integration-design.md "Import / Export Integration Scope" | PM-034..PM-039 (Sales denied) | ACC-020 |
| CONTRACT-018 | SVC-010 | Event | SVC-010 | SVC-008, SVC-009 | api-spec.md (ImportRunCompleted, ExportRunCompleted); audit-log-spec.md EVT-IMPORT-RUN, EVT-EXPORT-RUN | PM-034..PM-039 | ACC-020, ACC-022 |
| CONTRACT-019 | All services | Permission/Service-Auth | each service | each caller | authz-architecture.md "Service-To-Service Authorization" (serviceId, Bearer service-token, X-Intent, audience, 5-min lifetime, `SERVICE_AUTH_FAILED`); service-boundary-security.md Service-To-Service Permission Input Matrix | STB-003, SVC-ACC-008 | ACC-002, ACC-014, ACC-022 |
| CONTRACT-020 | All editable services | Error/Concurrency | each editable service | callers | api-spec.md "Editable Record Concurrency" (`expectedVersion`,`VERSION_CONFLICT`); api-spec.md "Error Codes" | n/a (integrity) | ACC-005..ACC-013 (editable records) |

## Service Data Ownership Mapping

Per-service owned data/storage from `data-design.md` Data Ownership Map; read
models published; forbidden direct access from `module-boundaries.md` /
`data-design.md`; consistency strategy from `data-design.md` Consistency Strategy.
Data classification and committed retention (PIM-027/028, PIM-INV-049..052,
privacy-requirements.md) are carried into the owned-data/retention notes;
concrete durations are from privacy-requirements.md Retention Policy and the
storage/TTL/backup are anchored here at PSM.

| Service | Owned Data / Storage | Read Models Published | Forbidden Direct Access | Consistency Strategy | Classification & Retention (durations: privacy-requirements.md) |
|---|---|---|---|---|---|
| SVC-002 | `identity_authz`: users, roles, sessions, permission policy (user `crm_identity_authz_user`) | current-user, session-check, permission-check responses | All CRM business record tables | Local txn for user/role/session; authz-version increment + session revoke on role/status change | Security Critical (PRIV-001); user identity/role/status retained while account exists + 7 years in operation logs after deactivation; no hard delete |
| SVC-003 | `lead`: leads, owner/status, qualification, conversion marker (`crm_lead_user`) | lead list/detail/search, duplicate-check signal, archive-eligibility | account, contact, opportunity, commercial, work, audit, reporting tables | Local txn for Lead aggregate + outbox event candidate; conversion-once guard inside aggregate | Confidential (PRIV-002); active retained; archived 7 years after archive or final linked opportunity closure; no hard delete |
| SVC-004 | `account`: accounts/customers, contacts, ownership/status (`crm_account_user`) | account/contact summary, list/detail/search, duplicate, archive-eligibility | lead, opportunity, commercial, work, audit, reporting tables | Local txn for Account+Contact aggregate + outbox; `expectedVersion` on editable | Confidential (PRIV-003/004); contact details restricted handling; archived 7 years after archive or final related closure; no hard delete |
| SVC-005 | `opportunity`: opportunities, Pipeline Stage (sole lifecycle dimension, DEC-020), closure data (`crm_opportunity_user`) | opportunity summary, closure/archive eligibility | lead, account, commercial payment/contract tables, work, audit, reporting | Local txn + outbox; Won (on related contract Signed, DEC-017)/Lost terminal lock; reads contract-signed/payment status via SVC-006 summary, not DB | Confidential (PRIV-005); retained 7 years after Won/Lost or archive; terminal non-reopen; no hard delete |
| SVC-006 | `commercial`: quotes, contracts, payment plans, actual payments, payment status (`crm_commercial_user`) | payment-status summary, contract reminder eligibility, quote/contract/payment detail | opportunity source tables, lead, account, work, audit, reporting | Local txn for Contract aggregate (plans+payments) so overpayment/full-payment invariants are atomic; idempotency key on payment commands; outbox | Restricted (PRIV-006/007/008); quote/contract/payment retained 7 years after closure/completion/full-payment/archive; no hard delete |
| SVC-007 | `work`: activities, notes, tasks, reminder projection (`crm_work_user`) | reminder list, active-obligation list, activity/task list/detail | source CRM record tables, commercial source tables, audit, reporting | Local txn + outbox; reminder projection derived from owned tasks + SVC-006 eligibility query; owner-transfer idempotent | Confidential (PRIV-009); retained with related record, 7 years after related closure/archive; no hard delete |
| SVC-008 | `audit_history`: record-local history, operation logs, security/audit events (`crm_audit_history_user`) | record-history query, admin operation-log query | source domain tables | Append-only; same durable workflow as sensitive mutation (AUD-IMM-002); tamper-evidence hash chain (`prevHash`/`eventHash`); no update/delete path | Restricted/Security Critical (PRIV-010/011/016); history ≥ related record; operation logs 7 years business / ≥3 years access-failure; append-only, no hard delete |
| SVC-009 | `reporting`: read-model projections, report snapshots, overview metrics (`crm_reporting_user`) | team overview, basic sales reports | source domain tables | Derived read model from source events / approved Query APIs; rebuild from contracts if stale; authz-before-aggregate | Confidential/Restricted-when-amounts (PRIV-014); snapshots not retained unless controlled artifact defined; source follows own retention |
| SVC-010 | `import_export`: import/export runs, row results, export metadata (`crm_import_export_user`) | run status, row results, export metadata | target domain tables | Mutates only via target Command APIs; row-level failure isolation; idempotency key per row; operation-status contract | Restricted (PRIV-012/013); raw import file processing-duration only; result metadata 1 year; export file not retained server-side beyond 24h temp window then deleted; imported records follow own class |
| Backup (job, not a service) | encrypted timestamped PostgreSQL backup on `srv-volcengine-sh-01` `/opt/crm-system/backups/postgres`; off-server copy to `srv-aliyun-bj-01` | none (no product API) | product service business credentials; UI exposure | Whole-instance backup; 7-day local retention; off-server copy + restore rehearsal before production release | Backups must not shorten application retention below privacy-requirements.md Retention Policy; keys outside repo/service config (data-design.md Backup) |

## Cross-Service Flow Mapping

Reuses INT-FLOW-001..013 from `integration-design.md` (referenced as FLOW-001..013
below; same intent/IDs). Primary flow owner agent, services, trigger,
contracts/events, idempotency, retry/timeout, compensation, and correlation ID are
reflected from `integration-design.md` Reliability Rules and Event Delivery
Strategy. Architecture-committed mechanism resolutions are recorded in the
Resolved Mechanisms note below.

| Flow ID | Flow | Primary Flow Owner Agent | Services | Trigger | Contracts / Events | Idempotency | Retry / Timeout | Compensation | Correlation ID |
|---|---|---|---|---|---|---|---|---|---|
| FLOW-001 | Sign in & protected work | backend-engineer | SVC-001, SVC-002, target service, SVC-008 | User authn / protected request | CONTRACT-001, CONTRACT-002; UserSignedIn, UserAccessDenied | Idempotent reads; session create idempotent by credential check | Query 3s / Command 5s; deny safely on dependency failure | Deny with no target mutation | Required |
| FLOW-002 | Lead to opportunity | backend-engineer | SVC-001, SVC-002, SVC-003, SVC-004, SVC-005, SVC-008, SVC-009 | Convert lead command | CONTRACT-003/004/005/007/008; LeadConverted, OpportunityCreated | `idempotencyKey` on ConvertLead; conversion-once guard | Command 5s; retry safe by idempotency key | Failed downstream leaves no false success; converted lead cannot reconvert | Required |
| FLOW-003 | Opportunity to quote/contract | backend-engineer | SVC-005, SVC-006, SVC-008, SVC-009 | Quote/contract command | CONTRACT-007/009/010; QuoteAccepted, ContractStatusChanged | `expectedVersion`; one-Accepted-quote guard in SVC-006 | Command 5s | Reject invalid links; maintain history on success only | Required |
| FLOW-004 | Contract signing to Won (payment tracked post-sale) | backend-engineer | SVC-006, SVC-005, SVC-008, SVC-009 | Contract signed verification → Close Won | CONTRACT-009/007/010/008; ContractStatusChanged(Signed), OpportunityClosed(Won) | `idempotencyKey` on SignContract & CloseWon | Command 5s; retry safe by idempotency key | Won requires Signed contract (DEC-017); payment decoupled (DEC-019) — overpayment still rejected on payment; early Won (no Signed contract) rejected | Required |
| FLOW-005 | Work reminders | backend-engineer | SVC-007, SVC-006, SVC-005, SVC-004, SVC-008 | Reminder query at business date | CONTRACT-011/009/010; TaskStatusChanged, PaymentOverdue, ContractStatusChanged | Query idempotent; overdue evaluation deterministic by business date | Query 3s | Hide unauthorized/inactive; stale reminder refresh | Required |
| FLOW-006 | Record history & operation logs | backend-engineer | source services, SVC-008 | Sensitive mutation | CONTRACT-013/014; HistoryEventAppended, OperationLogAppended | Event ID dedupe; append-only | Command 5s; reliable append in same durable workflow | Source defines behavior if audit append fails; P0 sensitive mutation requires reliable history path (outbox) | Required |
| FLOW-007 | Reports & overview | backend-engineer | source services, SVC-009, SVC-002 | Report/overview query or projection rebuild | CONTRACT-015/016; source domain events, ReportProjectionUpdated | Projection rebuild idempotent | Query 3s | Rebuild projection from approved contracts if stale | Required |
| FLOW-008 | CSV import/export | backend-engineer | SVC-010, target services, SVC-008, SVC-009 | Start import/export run | CONTRACT-017/018; ImportRunCompleted, ExportRunCompleted | Per-row idempotency key; run idempotent | Operation-status contract (long-running); per-call 5s | Row-level failure isolation; valid rows do not corrupt existing records | Required |
| FLOW-009 | Backup & restore | infrastructure-ops | PostgreSQL, backup job, runtime services | Scheduled backup / restore rehearsal | data-design.md Backup; deployment-notes.md Restore | Backup file per run, no overwrite | Backup job retry per ops; daily ~02:00 | Restore rehearsal before launch; same-host risk recorded | Required for evidence |
| FLOW-010 | Archive eligibility | backend-engineer | record-owning services, SVC-007, SVC-006, SVC-008 | Archive pre-check / archive command | CONTRACT-005/007/009/011; api-spec.md "Archive Eligibility"; RecordArchived, ArchiveBlocked | `expectedVersion` on archive; eligibility check idempotent | Query 3s / Command 5s | Block archive with active obligations; return obligation DTO; retry after refresh | Required |
| FLOW-011 | Owner transfer & open-work transfer | backend-engineer | record-owning services, SVC-007, SVC-008 | Owner transfer command | CONTRACT-005/007/011; api-spec.md "Owner Transfer"; OwnerChanged, OpenWorkTransferred | Idempotent transfer; work transfer by event ID/idempotency key | Command 5s; PendingRetry queued | `Failed` blocks record for operator review before release evidence | Required |
| FLOW-012 | Duplicate warning | backend-engineer | SVC-003, SVC-004, SVC-002 | Create/edit with duplicate match | CONTRACT-003/005; api-spec.md "Duplicate Warning"; DuplicateWarningRaised | Single-use `warningToken` idempotent with command key | Query 3s | No merge/overwrite; proceed creates new record only | Required |
| FLOW-013 | Close Lost terminal lifecycle | backend-engineer | SVC-005, SVC-007, SVC-008, SVC-009 | Close-lost command | CONTRACT-007/008/011; api-spec.md "Close Opportunity Lost"; OpportunityClosedLost, WorkItemCreated | `idempotencyKey` + `expectedVersion` on close-lost | Command 5s | Lost reason required; later edits rejected (`TERMINAL_RECORD_READ_ONLY`); notes/tasks via SVC-007 | Required |

### Resolved Mechanisms (architecture committed these — resolved here, not deferred)

- Overdue-evaluation trigger (PIM-OPEN-002 / BLK-A01): RESOLVED as on-read
  evaluation. `api-spec.md` "Reminder Query" computes due/overdue against the
  supplied `businessDate` in `Asia/Shanghai` at query time; `PaymentOverdue` /
  task-overdue are emitted on the read/evaluation path, not by a mandatory
  scheduled sweep. Deterministic for G7 test design via the supplied business
  date. Refs: api-spec.md "Reminder Query", integration-design.md FLOW-005.
- Sync vs async per flow: RESOLVED per `integration-design.md` Reliability Rules
  — synchronous internal Command/Query calls (5s/3s timeouts) for FLOW-001/002/
  003/004/005/010/011/012/013 user-initiated paths; asynchronous reliable events
  for history/reporting/reminder propagation (FLOW-006/007 and the event legs of
  002/003/004/008/013). Mixed flows use sync command + async event publication.
- Outbox vs replay (Event Delivery Strategy): RESOLVED as the database outbox
  table per producing service plus a background dispatcher (the first of the two
  architecture-permitted options), satisfying the same-durable-workflow guarantee
  (AUD-IMM-002, ARCH-RISK-003). G8 tasks must include duplicate/retry/failure
  tests (deferred to G8 task IDs).
- Import/export temp-file cleanup: RESOLVED per `api-spec.md` "Import And Export
  Runs" — temporary upload/export files retained 24h after run completion/failure
  on a controlled path under `import_export` ownership, then deleted by a
  scheduled cleanup job; deletion status exposed in run metadata; cleanup
  success/failure logged. Satisfies PIM-INV-052, PRIV-012/013.

## Architecture Acceptance

ARCH-ACC-* records the architecture concerns this PSM reflects, including the
carried release blockers (off-server backup + restore rehearsal, HTTPS/TLS,
security group, monitoring). Source cites the architecture/risk doc; PSM Element
is the modeled element; Verification is the planned method; Status is `Reflected`
(modeled in PSM) or `Release-evidence pending` (proven at G11/G12).

| ID | Architecture Concern | Source | PSM Element | Service | Product Acceptance ID | Verification | Status |
|---|---|---|---|---|---|---|---|
| ARCH-ACC-001 | Service boundaries do not degrade into shared DB/logic | ARCH-RISK-001, module-boundaries.md | Data Ownership Mapping, CONTRACT-019 | SVC-001..010 | ACC-002, ACC-016 | Static/audit check | Reflected |
| ARCH-ACC-002 | Authorization enforced server-side, not only frontend/gateway | ARCH-RISK-002, authz-architecture.md | CONTRACT-001, permission rules PM-* | SVC-002 + domain services | ACC-002 | Direct-API denial tests | Reflected |
| ARCH-ACC-003 | Audit/history not lost after sensitive mutation | ARCH-RISK-003, AUD-IMM-002 | Outbox resolution, CONTRACT-013/014, FLOW-006 | SVC-008 + producers | ACC-014, ACC-022 | Reliability/integration test | Reflected |
| ARCH-ACC-004 | Encrypted local backup + off-server copy + restore rehearsal | ARCH-RISK-004, ADR-ARCH-004, deployment-notes.md | PSM-014, FLOW-009, Backup ownership row | Backup job, PostgreSQL | ACC-016, ACC-017 | Off-server copy + restore rehearsal evidence | Release-evidence pending |
| ARCH-ACC-005 | Import/export goes through domain services | ARCH-RISK-005, integration-design.md | CONTRACT-017, FLOW-008, SVC-010 ownership | SVC-010 | ACC-020 | Integration test | Reflected |
| ARCH-ACC-006 | Reporting uses read models / events, not source tables | ARCH-RISK-006, data-design.md Reporting | CONTRACT-015/016, SVC-009 ownership | SVC-009 | ACC-018, ACC-023 | Data-access audit | Reflected |
| ARCH-ACC-007 | Cross-service timeout/retry/idempotency consistent | ARCH-RISK-008, integration-design.md Reliability | FLOW-001..013, CONTRACT-020 | all flow services | ACC-004..ACC-021 (flows) | Integration test | Reflected |
| ARCH-ACC-008 | Production HTTPS endpoint + TLS evidence | ARCH-RISK-009, ADR-ARCH-005, deployment-notes.md | PSM-014, reverse-proxy ingress, secure cookie | SVC-001, reverse proxy | ACC-017 | HTTPS/TLS + security-header evidence | Release-evidence pending |
| ARCH-ACC-009 | Service-to-service calls use signed tokens, not network trust | ARCH-RISK-011, authz-architecture.md S2S | CONTRACT-019 | all services | ACC-002, ACC-022 | Security test (SERVICE_AUTH_FAILED) | Reflected |
| ARCH-ACC-010 | Stale edits cannot overwrite newer updates | ARCH-RISK-012, data-design.md Concurrency | CONTRACT-020 (`expectedVersion`,`VERSION_CONFLICT`) | editable services | ACC-005..ACC-013 | Concurrency test | Reflected |
| ARCH-ACC-011 | Archive blocked on active downstream obligations | ARCH-RISK-013, integration-design.md FLOW-010 | FLOW-010, archive-eligibility contract | record-owning + SVC-006/007 | ACC-002, ACC-014 | Integration test | Reflected |
| ARCH-ACC-012 | CSV formula-injection safety | ARCH-RISK-014, api-spec.md Import/Export | CONTRACT-017 (dangerous-cell handling) | SVC-010 | ACC-020 | Export safety test | Reflected |
| ARCH-ACC-013 | Security-group / network exposure boundary | deployment-notes.md Network Exposure, STB-008 | PSM-014 (no public DB/internal ports/backup dir) | runtime host, reverse proxy | ACC-017 | Security-group inbound-rule evidence | Release-evidence pending |
| ARCH-ACC-014 | Monitoring / health / observability before release | deployment-notes.md Health And Observability | PSM-014 per-service health endpoint, log correlation | SVC-001..010, ops | ACC-017 | Monitoring/alert/health evidence | Release-evidence pending |
| ARCH-ACC-015 | Persistence survives restart for all owned data | ARCH-RISK (persistence), data-design.md | Data Ownership Mapping, PSM-001..012 | SVC-002..010, PostgreSQL | ACC-016 | Restart persistence evidence | Release-evidence pending |

## Service Architecture Acceptance

Reuses SVC-ACC-001..012 from `service-architecture-acceptance.md` verbatim in
intent. Owner-agent and completion standards are reflected; PSM Element points to
the modeled element that carries the governance concern.

| ID | Priority | Service Governance Concern | PSM Element | Completion Standard | Verification | Status |
|---|---|---|---|---|---|---|
| SVC-ACC-001 | P0 | Service-boundary-first governance | Service Mapping (SVC-001..010) | Capability map, service list, ADR exist and are represented | Document review | Reflected |
| SVC-ACC-002 | P0 | Every service has exactly one owner agent | Service Mapping (owner agent = backend-engineer) | One owner agent per service | Service list review | Reflected |
| SVC-ACC-003 | P0 | P0/P1 acceptance maps to services | PSM Traceability (ACC-001..023) | Every P0/P1 maps to ≥1 service | Traceability review | Reflected (this PSM) |
| SVC-ACC-004 | P0 | Service data ownership explicit | Service Data Ownership Mapping | Each service declares owned data + forbidden access | Architecture review/audit | Reflected |
| SVC-ACC-005 | P0 | Contracts exist before implementation | Service Contract Mapping (CONTRACT-001..020) | API/event/error/permission contracts cover P0/P1, referenced by tasks | Contract/task review | Reflected (tasks pending G8) |
| SVC-ACC-006 | P0 | Cross-service internal dependency prohibited | Forbidden Direct Access column; CONTRACT-019 | No cross-service internal imports/shared business code | Static check/audit | Implementation evidence pending (G9/G12) |
| SVC-ACC-007 | P0 | Cross-service DB access prohibited | Data Ownership Mapping (DB users, forbidden tables) | No direct cross-service DB read/write | Code/data-access audit | Implementation evidence pending (G9/G12) |
| SVC-ACC-008 | P0 | Service-to-service calls enforce security boundaries | CONTRACT-019 (signed token, intent, audience) | Authn/authz/audit/rotation/rejection rules exist + tested | Security review/tests | Reflected (tests pending G10) |
| SVC-ACC-009 | P0 | Cross-service reliability designed and tested | Cross-Service Flow Mapping + Resolved Mechanisms | Idempotency/retry/timeout/compensation/correlation/outbox/owner-transfer/archive/terminal | Integration tests | Reflected (evidence pending G11) |
| SVC-ACC-010 | P0 | AI tasks constrained by service boundaries | (forward to G8 tasks) | Tasks include service, owner agent, contract, ACC, forbidden access | Task review | Pending G8 |
| SVC-ACC-011 | P0 | Fake core implementation cannot pass | Technical Models (real persistence per service) | No mock/stub/TODO/static/non-persistent satisfies P0/P1 | QA/audit | Implementation evidence pending (G10/G12) |
| SVC-ACC-012 | P0 | End-to-end evidence for P0/P1 cross-service flows | FLOW-001..013 + PSM Traceability Integration Evidence column | Evidence incl. ACC, env, steps, result, service chain, correlation ID, endpoint, monitoring, backup/restore | Integration report/audit | Pending G11/G12 |

## PSM Traceability

ACC-001..023 → owning Service → Owner Agent (backend-engineer) → Contract → PIM ID
→ PSM ID → ARCH-ACC → Task ID → Test ID → Integration Evidence. Task/Test/
Integration are FORWARD placeholders produced at G7/G8/G11 — marked `pending`,
not invented.

| Product Acceptance ID | Service ID | Service Owner Agent | Contract ID | PIM ID | PSM ID | Architecture Acceptance ID | Task ID | Test ID | Integration Evidence |
|---|---|---|---|---|---|---|---|---|---|
| ACC-001 | SVC-002 | backend-engineer | CONTRACT-001, 002 | PIM-001, PIM-002 | PSM-001 | ARCH-ACC-002, 008 | pending (G8) | pending (G7) | pending (G11) |
| ACC-002 | SVC-002 (+ all domain) | backend-engineer | CONTRACT-001, 019 | PIM-001, PIM-002, PIM-003 | PSM-001 | ARCH-ACC-001, 002, 009 | pending (G8) | pending (G7) | pending (G11) |
| ACC-003 | SVC-003 | backend-engineer | CONTRACT-003, 004 | PIM-004, PIM-003 | PSM-002 | ARCH-ACC-002, 010 | pending (G8) | pending (G7) | pending (G11) |
| ACC-004 | SVC-003 | backend-engineer | CONTRACT-003, 004 | PIM-004 | PSM-002 | ARCH-ACC-003, 007 | pending (G8) | pending (G7) | pending (G11) |
| ACC-005 | SVC-004 | backend-engineer | CONTRACT-005, 006 | PIM-005 | PSM-003 | ARCH-ACC-010, 011 | pending (G8) | pending (G7) | pending (G11) |
| ACC-006 | SVC-004 | backend-engineer | CONTRACT-005, 006 | PIM-006 | PSM-003 | ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G11) |
| ACC-007 | SVC-005 | backend-engineer | CONTRACT-007, 008 | PIM-007, PIM-012 | PSM-004 | ARCH-ACC-007, 010 | pending (G8) | pending (G7) | pending (G11) |
| ACC-008 | SVC-005 | backend-engineer | CONTRACT-007, 008 | PIM-007 | PSM-004 | ARCH-ACC-003, 007 | pending (G8) | pending (G7) | pending (G11) |
| ACC-009 | SVC-006 | backend-engineer | CONTRACT-009, 010 | PIM-008, PIM-012 | PSM-005 | ARCH-ACC-003, 010 | pending (G8) | pending (G7) | pending (G11) |
| ACC-010 | SVC-006 | backend-engineer | CONTRACT-009, 010 | PIM-009, PIM-012 | PSM-006 | ARCH-ACC-010 | pending (G8) | pending (G7) | pending (G11) |
| ACC-011 | SVC-006 | backend-engineer | CONTRACT-009, 010 | PIM-010, PIM-011, PIM-012 | PSM-007 | ARCH-ACC-007, 003 | pending (G8) | pending (G7) | pending (G11) |
| ACC-012 | SVC-007 | backend-engineer | CONTRACT-011, 012 | PIM-013, PIM-014, PIM-015 | PSM-008 | ARCH-ACC-003 | pending (G8) | pending (G7) | pending (G11) |
| ACC-013 | SVC-005 | backend-engineer | CONTRACT-007, 008, 009, 010 | PIM-007, PIM-011 | PSM-004, PSM-007 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G11) |
| ACC-014 | SVC-008 | backend-engineer | CONTRACT-013, 014 | PIM-018, PIM-019 | PSM-009 | ARCH-ACC-003 | pending (G8) | pending (G7) | pending (G11) |
| ACC-015 | SVC-001 | backend-engineer | CONTRACT-001 (authz), record-owning Query APIs | PIM-026, PIM-020 | PSM-011 | ARCH-ACC-002 | pending (G8) | pending (G7) | pending (G11) |
| ACC-016 | SVC-002..010 | backend-engineer | CONTRACT-020, 014; data-design.md | PIM-001..028 (persisted) | PSM-013 (+ PSM-001..012) | ARCH-ACC-004, 015 | pending (G8) | pending (G7) | pending (G11) |
| ACC-017 | runtime host + SVC-001..010 | backend-engineer (+ infrastructure-ops env evidence) | deployment-notes.md (no API contract) | (deployment, PIM-OPEN-004) | PSM-014 | ARCH-ACC-004, 008, 013, 014 | pending (G8) | pending (G7) | pending (G11) |
| ACC-018 | SVC-009 | backend-engineer | CONTRACT-015, 016 | PIM-024 | PSM-010 | ARCH-ACC-006 | pending (G8) | pending (G7) | pending (G11) |
| ACC-019 | SVC-003, SVC-004 | backend-engineer | CONTRACT-003, 005 | PIM-021 | PSM-002, PSM-003 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G11) |
| ACC-020 | SVC-010 | backend-engineer | CONTRACT-017, 018 | PIM-022, PIM-023 | PSM-012 | ARCH-ACC-005, 012 | pending (G8) | pending (G7) | pending (G11) |
| ACC-021 | SVC-007 | backend-engineer | CONTRACT-011, 010 (payment overdue) | PIM-015, PIM-016, PIM-017 | PSM-008 | ARCH-ACC-007 | pending (G8) | pending (G7) | pending (G11) |
| ACC-022 | SVC-008 | backend-engineer | CONTRACT-013, 014, 002 | PIM-019, PIM-001 | PSM-009 | ARCH-ACC-003, 009 | pending (G8) | pending (G7) | pending (G11) |
| ACC-023 | SVC-009 | backend-engineer | CONTRACT-015, 016 | PIM-025, PIM-020 | PSM-010 | ARCH-ACC-006 | pending (G8) | pending (G7) | pending (G11) |

Coverage note: all 23 product acceptance items (ACC-001..023, including the P1
items ACC-018..023) trace through an owning service + contract + PIM + PSM + ARCH-ACC.
ACC-017 deployment is in-tier at PSM via PSM-014 (deployment boundary). No P0/P1
item is downgraded, merged, or satisfied by a non-persistent/mock path.

## Open / Blocked

The three former product/BA blockers BLK-001/002/003 were RESOLVED 2026-06-01 by a
Formal Scope Change by User (decision-log.md DEC-017..020); they are no longer open
pre-G8 decisions.

- BLK-001 (ACC-007) — Opportunity Status enumerated value set — RESOLVED 2026-06-01
  by Formal Scope Change (DEC-020). The separate Status dimension is removed;
  Pipeline Stage is the sole opportunity lifecycle dimension. PSM-004 / CONTRACT-007
  carry no separate Status field. Ref: DEC-020, PIM-OPEN-003 (resolved), CIM-016
  (retired).
- BLK-002 (ACC-011, ACC-013) — Multi-plan contract "fully paid" aggregation —
  RESOLVED 2026-06-01 by Formal Scope Change (DEC-017 + DEC-019). Won no longer
  depends on full payment (Won = contract Signed, DEC-017), so "contract fully
  paid" is no longer a Won precondition. PSM-007 / CONTRACT-009 retain the
  contract-level overpayment ceiling and payment-status summary for post-sale
  collection (payment decoupled, DEC-019). Ref: DEC-017, DEC-019, PIM-OPEN-005
  (resolved).
- BLK-003 (ACC-009) — Second-quote-accept observable outcome — RESOLVED 2026-06-01
  by Formal Scope Change (DEC-018). With exactly one quote per opportunity there is
  no second quote to accept, so the losing-path ambiguity no longer exists. PSM-005
  / CONTRACT-009/010 reflect the one-quote rule. Ref: DEC-018, PIM-OPEN-001
  (resolved), EDGE-012.

Architecture-deferred-but-now-resolved (recorded here for visibility; resolved in
"Resolved Mechanisms" above, not left open): overdue-evaluation trigger (BLK-A01 →
on-read), sync/async per flow, outbox-vs-replay (→ database outbox + dispatcher),
import/export temp-file cleanup (→ 24h temp window + scheduled cleanup job). No
genuinely architecture-pending PSM input remains: every mechanism the architecture
left to "G6 PSM must choose/model" has been chosen above from an
architecture-permitted option.

Carried-forward release blockers (release-time evidence, not gate blockers;
modeled at PSM altitude via ARCH-ACC-004/008/013/014/015): encrypted off-server
backup copy + restore rehearsal, HTTPS/TLS endpoint, security-group rules, and
monitoring/health evidence. These are proven at G11 and audited at G12. Ref:
blockers.md "Carried-forward Release Blockers", deployment-notes.md, RISK-002/004/
009/011.
