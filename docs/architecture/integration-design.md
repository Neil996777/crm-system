# Integration Design

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Integration Strategy

The CRM uses internal service APIs for synchronous commands and queries, plus
domain events for history, reporting, reminders, and integration evidence.

All cross-service interactions must include:

- correlation ID
- caller service identity
- actor context where user initiated
- timeout
- retry policy when safe
- idempotency key for writes
- safe error handling
- audit/history behavior where applicable

## Cross-Service Flow Matrix

| Flow ID | Flow | Primary Flow Owner Agent | Services | Sync Calls | Events | Failure Recovery | Correlation ID |
|---|---|---|---|---|---|---|---|
| INT-FLOW-001 | Sign in and protected work | backend-engineer | gateway-bff, identity-authz-service, target service, audit-history-service | auth/session, permission check | UserSignedIn, UserAccessDenied | Deny safely; no target mutation. | Required |
| INT-FLOW-002 | Lead to opportunity | backend-engineer | gateway-bff, identity-authz-service, lead-service, account-service, opportunity-service, audit-history-service, reporting-service | permission, create/link account/contact, create opportunity | LeadConverted, OpportunityCreated | Idempotent conversion; converted lead cannot convert again; failed downstream call leaves no false success. | Required |
| INT-FLOW-003 | Opportunity to quote/contract | backend-engineer | opportunity-service, commercial-service, audit-history-service, reporting-service | opportunity summary, commercial commands | QuoteAccepted, ContractStatusChanged | Reject invalid links; maintain history on successful mutation. | Required |
| INT-FLOW-004 | Payment to Won | backend-engineer | commercial-service, opportunity-service, audit-history-service, reporting-service | payment status verification | PaymentRecorded, OpportunityClosed | Overpayment/early Won rejected; retry safe by idempotency key. | Required |
| INT-FLOW-005 | Work reminders | backend-engineer | work-service, commercial-service, opportunity-service, account-service, audit-history-service | related record summary/eligibility | TaskStatusChanged, PaymentOverdue, ContractStatusChanged | Hide unauthorized/inactive records; stale reminder refresh. | Required |
| INT-FLOW-006 | Record history and operation logs | backend-engineer | source services, audit-history-service | append/query log APIs where needed | OperationLogAppended, HistoryEventAppended | Source operation must define behavior if audit append fails. P0 sensitive mutations require reliable history path. | Required |
| INT-FLOW-007 | Reports and overview | backend-engineer | source services, reporting-service, identity-authz-service | report query, optional projection rebuild | source domain events, ReportProjectionUpdated | Rebuild projection from approved contracts if projection stale. | Required |
| INT-FLOW-008 | CSV import/export | backend-engineer | import-export-service, target services, audit-history-service, reporting-service | target service commands/queries | ImportRunCompleted, ExportRunCompleted | Row-level failure isolation; valid rows do not corrupt existing records. | Required |
| INT-FLOW-009 | Backup and restore | infrastructure-ops | PostgreSQL, backup job, runtime services | health and restore checks | operational log event candidate | Restore rehearsal before launch; same-host backup risk recorded. | Required for evidence |
| INT-FLOW-010 | Archive eligibility | backend-engineer | record-owning services, work-service, commercial-service, audit-history-service | obligation checks, archive command | RecordArchived, ArchiveBlocked candidate | Block archive when active obligations exist; return obligation DTO and retry after refresh. | Required |
| INT-FLOW-011 | Owner transfer and open work transfer | backend-engineer | record-owning services, work-service, audit-history-service | owner transfer command, open work transfer command/query | OwnerChanged, OpenWorkTransferred | Idempotent transfer; retry pending transfer; manual exception requires privileged reason. | Required |
| INT-FLOW-012 | Duplicate warning | backend-engineer | lead-service, account-service, identity-authz-service | safe duplicate lookup, proceed-after-warning command | DuplicateWarningRaised | No merge/overwrite; warning token is single-use and idempotent with command key. | Required |
| INT-FLOW-013 | Close Lost terminal lifecycle | backend-engineer | opportunity-service, work-service, audit-history-service, reporting-service | close lost command, post-close work command | OpportunityClosedLost, WorkItemCreated | Lost reason required; later edits rejected; notes/tasks remain allowed through work-service. | Required |

## Lead Conversion Sequence

```mermaid
sequenceDiagram
  participant U as Sales
  participant G as gateway-bff
  participant P as identity-authz-service
  participant L as lead-service
  participant A as account-service
  participant O as opportunity-service
  participant H as audit-history-service
  participant R as reporting-service

  U->>G: Convert lead
  G->>P: Check permission
  P-->>G: Allowed
  G->>L: ConvertLead(command, idempotencyKey)
  L->>L: Validate state and conversion-once guard
  L->>A: CreateOrLinkAccountContact
  A-->>L: Account/contact references
  L->>O: CreateOpportunity
  O-->>L: Opportunity reference
  L->>L: Persist converted state
  L-->>H: LeadConverted event
  O-->>H: OpportunityCreated event
  L-->>R: LeadConverted event
  O-->>R: OpportunityCreated event
  L-->>G: Conversion result
  G-->>U: Converted lead and opportunity link
```

## Payment To Won Sequence

```mermaid
sequenceDiagram
  participant U as Sales
  participant G as gateway-bff
  participant P as identity-authz-service
  participant C as commercial-service
  participant O as opportunity-service
  participant H as audit-history-service
  participant R as reporting-service

  U->>G: Record payment
  G->>P: Check commercial permission
  P-->>G: Allowed
  G->>C: RecordPayment(command, idempotencyKey)
  C->>C: Validate amount and overpayment
  C-->>H: PaymentRecorded event
  C-->>R: PaymentRecorded event
  C-->>G: Payment status

  U->>G: Close opportunity Won
  G->>P: Check opportunity close permission
  P-->>G: Allowed
  G->>O: CloseWon(command, idempotencyKey)
  O->>C: GetPaymentStatusSummary
  C-->>O: Paid / not paid
  O->>O: Persist terminal Won if fully paid
  O-->>H: OpportunityClosed event
  O-->>R: OpportunityClosed event
  O-->>G: Won result
  G-->>U: Won state
```

## Import Sequence

```mermaid
sequenceDiagram
  participant U as Admin/Manager
  participant G as gateway-bff
  participant P as identity-authz-service
  participant I as import-export-service
  participant T as target domain service
  participant H as audit-history-service

  U->>G: Start CSV import
  G->>P: Check import permission
  P-->>G: Allowed
  G->>I: StartImport(file metadata)
  I->>I: Parse and validate rows
  loop each valid row
    I->>T: Domain command with idempotency key
    T-->>I: Success or safe row error
  end
  I-->>H: ImportRunCompleted event
  I-->>G: Run result summary
  G-->>U: Row-level import result
```

## Archive Eligibility Sequence

```mermaid
sequenceDiagram
  participant G as gateway-bff
  participant R as record-owning service
  participant W as work-service
  participant C as commercial-service
  participant H as audit-history-service

  G->>R: GetArchiveEligibility(recordId)
  R->>W: Query active obligations
  R->>C: Query commercial obligations
  W-->>R: Open task/follow-up obligations
  C-->>R: Pending signature/payment/quote obligations
  alt obligations exist
    R-->>G: ARCHIVE_BLOCKED with activeObligations[]
    R-->>H: ArchiveBlocked event when required
  else clear
    G->>R: ArchiveRecord(expectedVersion, reason)
    R-->>H: RecordArchived event
    R-->>G: Archived result
  end
```

## Owner Transfer Reliability

Owner transfer uses an idempotent owner transfer command plus an
`OwnerChanged` event. The record-owning service owns the record owner state.
work-service owns task and follow-up assignment state.

Required recovery behavior:

- The owner change command returns `workTransferStatus`.
- `Completed` means all open task/follow-up ownership was transferred or a
  privileged manual exception was recorded.
- `PendingRetry` means owner change is saved and work transfer is queued for
  retry by event ID/idempotency key.
- `Failed` means retry budget is exhausted and the record is blocked for
  operator review before release evidence can pass.
- Every status change emits history/operation log evidence.

## Import / Export Integration Scope

V1 import/export target routing:

| Object Type | Target Service | Mutation / Query Rule |
|---|---|---|
| Lead | lead-service | Import through lead command; export through authorized query. |
| Account / Customer | account-service | Import through account command; export through authorized query. |
| Contact | account-service | Import through contact command; export through authorized query. |
| Opportunity | opportunity-service | Import through opportunity command; export through authorized query. |
| Quote / Contract / Payment | commercial-service | Import only through supported commercial commands; export through authorized query. |
| Activity / Note / Task | work-service | Import through work commands; export through authorized query. |

Unsupported object types are rejected before mutation. Import/export service
may store run metadata and row results only; it may not mutate target service
tables directly.

## Reliability Rules

- Default internal Query API timeout is 3 seconds. Default internal Command API
  timeout is 5 seconds. Longer-running operations must use operation-status
  contracts instead of holding synchronous requests open.
- Internal service calls must present valid service authentication as described
  in `authz-architecture.md`.
- Retry only idempotent operations or commands with idempotency keys.
- Cross-service command retries must not duplicate business records.
- Long-running import/export operations must expose status and result query
  contracts.
- Event consumers must handle duplicate events by event ID.
- Event consumers must tolerate out-of-order events unless the contract requires
  ordered processing for a specific aggregate.
- Correlation ID is mandatory across gateway, services, events, logs, tests, and
  integration evidence.

## Event Delivery Strategy

For v1, Architecture requires an outbox-equivalent reliable publication pattern
for P0/P1 events. The exact implementation may be:

- database outbox table per producing service plus background dispatcher, or
- a transactionally persisted event record plus explicit replay mechanism.

G6 PSM must choose and model the concrete pattern. G8 tasks must include tests
for duplicate, retry, and failure cases.

## Integration Evidence Requirements

Integration Owner must later prove:

- service chain and correlation ID for every P0/P1 cross-service flow
- persisted data evidence after service restart
- role/scope enforcement across service boundaries
- history/log creation for sensitive flows
- import/export row result and audit evidence
- backup and restore rehearsal evidence before release
- off-server backup evidence before production release
