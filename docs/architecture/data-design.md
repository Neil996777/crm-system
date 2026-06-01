# Data Design

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Data Strategy

The CRM uses one self-hosted PostgreSQL instance in the committed scope, with service-level data
ownership enforced by database or schema boundaries and independent database
users.

The architecture may use either independent databases per service or independent
schemas per service inside the same PostgreSQL instance. The implementation
choice must preserve the same rule: a service can access only its owned
database/schema through its own credentials.

## Data Ownership Map

| Service | Owned Data | Database / Schema Name | Service DB User | Forbidden Direct Access |
|---|---|---|---|---|
| identity-authz-service | users, roles, sessions, permission policy metadata | `identity_authz` | `crm_identity_authz_user` | All CRM business record tables. |
| lead-service | leads, lead owner/status, lead qualification, lead conversion marker | `lead` | `crm_lead_user` | account, contact, opportunity, commercial, work, audit, reporting tables. |
| account-service | accounts/customers, contacts, account ownership/status | `account` | `crm_account_user` | lead, opportunity, commercial, work, audit, reporting tables. |
| opportunity-service | opportunities, stages, stage history, closure data | `opportunity` | `crm_opportunity_user` | lead, account, commercial payment/contract tables, work, audit, reporting tables. |
| commercial-service | quotes, contracts, payment plans, actual payments, payment status | `commercial` | `crm_commercial_user` | opportunity source tables, lead, account, work, audit, reporting tables. |
| work-service | activities, notes, tasks, reminder projection data | `work` | `crm_work_user` | source CRM record tables, commercial source tables, audit, reporting tables. |
| audit-history-service | record-local history, operation logs, security/audit events | `audit_history` | `crm_audit_history_user` | source domain tables. |
| reporting-service | read model projections, report snapshots, overview metrics | `reporting` | `crm_reporting_user` | source domain tables. |
| import-export-service | import/export runs, row results, export metadata | `import_export` | `crm_import_export_user` | target domain tables. |

## Data Access Rules

- Service credentials must be scoped to the owning database/schema only.
- Cross-service database grants are prohibited for P0/P1 paths.
- Direct foreign keys across service-owned schemas/databases are prohibited.
- Cross-service references use stable IDs, public summaries, events, or read
  models.
- Database migrations are owned and executed per service.
- Backup jobs may access the PostgreSQL instance for backup purposes only; they
  are not business services and must not expose data through product APIs.
- History/audit storage credentials must not expose normal update/delete paths
  for history or operation log records.
- Backup job credentials must be separated from service business credentials
  and restricted to backup/restore operations.

## Aggregate Ownership

| Aggregate | Owning Service | Notes |
|---|---|---|
| User | identity-authz-service | Includes active/disabled state and role assignment. |
| Lead | lead-service | Conversion marker prevents duplicate conversion. |
| Account / Customer | account-service | No hard delete; archive or lifecycle state only. |
| Contact | account-service | Contact details require restricted handling. |
| Opportunity | opportunity-service | Won/Lost are terminal in the committed scope. |
| Quote | commercial-service | Exactly one quote per opportunity (DEC-018); it is the Accepted quote linked to the contract. |
| Contract | commercial-service | Pending Signature requires expected signed date and note. |
| Payment Plan / Actual Payment | commercial-service | Blocks zero, negative, and overpayment amounts. |
| Activity / Note / Task | work-service | Related record reference must be validated through public contract. |
| History Event | audit-history-service | Append-only through normal CRM workflows. |
| Operation Log | audit-history-service | Administrator-only query. |
| Report Projection | reporting-service | Derived read model, not source truth. |
| Import / Export Run | import-export-service | Long-running operation state and row results. |

## Consistency Strategy

Local transactions:

- A service must use a local transaction for changes to its owned aggregate and
  outbox/audit event candidate where applicable.

Cross-service flows:

- Cross-service database transactions are not allowed.
- Cross-service writes use target service Command APIs.
- Multi-service operations use idempotency keys, correlation IDs, reliable
  event publication, and compensating state where needed.
- G6 PSM must represent which flows are synchronous, asynchronous, or mixed.

## History And Audit Data

Source services must publish or submit history/audit events for accepted
business mutations. The audit-history-service owns storage and query of:

- record-local business history
- admin/global operation logs
- access denials where required
- import/export run summaries
- sensitive lifecycle events

History/log data must be append-only through normal CRM workflows. No service
may provide normal product APIs to edit or delete history/log records.

Append-only architecture requirements:

- audit-history-service provides append and query APIs only. It provides no
  normal update/delete API for history events or operation logs.
- Administrative UI may query operation logs but cannot edit or delete them.
- Database permissions for the audit-history write path must be modeled so
  normal application roles cannot update or delete accepted history/log rows.
- History and operation log rows include `eventId`, `eventVersion`,
  `producerService`, `aggregateType`, `aggregateId`, `actorId`, `occurredAt`,
  `correlationId`, `causationId`, `safeSummary`, `prevHash`, and `eventHash`
  or an equivalent tamper-evidence field.
- Hash/tamper evidence must not include unrestricted sensitive field values in
  plaintext payloads.

## Reporting Data

The reporting-service owns read models and projections. It may build them from:

- domain events
- approved source service Query APIs
- explicit rebuild operations through public contracts

It must not directly query source service databases. Reports must apply
authorization before returning rows, summaries, or aggregates.

## Import / Export Data

The import-export-service owns run and row-result data only. It must not write
target domain tables directly.

Imports mutate domain data by calling target service Command APIs, preserving:

- role/scope authorization
- required-field validation
- domain state rules
- duplicate warning behavior
- history/audit event creation
- row-level safe error summaries

Exports use target service Query APIs or approved read models and must include
authorized records only.

## Backup And Restore Data

PostgreSQL backup uses encrypted timestamped backup files on the ECS host and
must copy encrypted backup evidence to an off-server target before production
release. Backups retain the whole PostgreSQL instance for the committed release, including
service-owned databases/schemas, history, operation logs, reporting
projections, and import/export run metadata.

Backup access controls:

- Local backup directory: `/opt/crm-system/backups/postgres`.
- Directory permissions: readable only by root/ops account and backup job.
- Backup files are encrypted before storage or immediately after creation.
- Encryption keys/secrets are stored outside the repository and outside product
  service configuration files where possible.
- Backup checksums are recorded for restore evidence.
- Backup logs must not print restricted data, full connection strings, or
  secrets.
- Restore procedures must run in a controlled target or controlled recovery
  procedure and must record operator, timestamp, source backup, checksum,
  decryption step, result, and cleanup action.

Restore validation must prove:

- core CRM data survives restore
- record-local history and admin operation logs survive restore
- service credentials remain correctly scoped after restore
- no backup secrets are exposed in product UI or public repository

## Concurrency Tokens

Editable P0/P1 aggregates must expose a numeric `version` and `updatedAt`.
Commands that mutate editable records must include `expectedVersion`.

Required aggregates:

- Lead
- Account / Customer
- Contact
- Opportunity
- Quote
- Contract
- Payment Plan where editable
- Activity / Note / Task where editable

Stale writes must fail with `VERSION_CONFLICT` and include a safe latest record
summary. Services must not silently overwrite stale user edits.

## Lifecycle Protection

- Won/Lost opportunity states are terminal in the committed scope.
- Closed opportunities reject ordinary update and stage transition commands.
- Post-close notes and follow-up tasks are stored in work-service and reference
  the closed opportunity without editing it.
- Archived records remain queryable only through explicit archived filters and
  cannot be modified by ordinary active-record commands.
