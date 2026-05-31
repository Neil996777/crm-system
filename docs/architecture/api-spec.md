# API Specification

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## API Strategy

The CRM uses contract-first service APIs. G5 defines required contract families
and representative endpoints. G6 PSM must refine them into platform-specific
models, and G7 must create contract and acceptance tests before G8 task
planning.

All APIs are business APIs, not generic database CRUD APIs.

## Common Request Requirements

Every protected request must include or derive:

- `correlationId`
- authenticated actor ID
- actor role
- actor active/disabled status
- actor authorization version or session version
- caller service identity for internal calls
- target action
- target resource type and resource ID where applicable
- idempotency key for non-idempotent writes
- `expectedVersion` for editable P0/P1 records

## Common Response Envelope

```text
Success:
{
  "correlationId": "...",
  "data": { ... }
}

Failure:
{
  "correlationId": "...",
  "error": {
    "code": "PERMISSION_DENIED",
    "category": "permission | validation | business_rule | conflict | dependency | system",
    "safeMessage": "...",
    "fieldErrors": [
      { "field": "amount", "code": "MUST_BE_POSITIVE", "safeMessage": "..." }
    ]
  }
}
```

Errors must not expose restricted record names, amounts, contact details, or
existence of unauthorized records.

## Contract Families

| Contract Family | Purpose | Required Before |
|---|---|---|
| Command API | Create, update, archive, transition, import, export, and other writes. | G8 |
| Query API | List, detail, search, filter, history, logs, reports, and safe lookups. | G8 |
| Event Contract | Domain event publishing for history, reporting, reminders, and integration evidence. | G8 |
| Error Contract | Safe validation, permission, blocked transition, conflict, and dependency failures. | G8 |
| Permission Contract | Actor/action/resource/scope checks and denial categories. | G8 |
| Idempotency Contract | Prevent duplicate writes for commands and long-running jobs. | G8 |
| Operation Status Contract | Import/export and slow command progress/result retrieval. | G8 |

## Service API Summary

| Service | Command APIs | Query APIs | Event Contracts |
|---|---|---|---|
| identity-authz-service | sign in, sign out, create/update user, change role/status, revoke sessions | current user, session check, permission check, role list | UserSignedIn, UserSignedOut, UserAccessDenied, UserRoleStatusChanged, SessionRevoked |
| lead-service | create/update/assign lead, qualify lead, restore invalid lead, convert lead, proceed after duplicate warning | lead list/detail/search, lead duplicate check signal, archive eligibility | LeadCreated, LeadOwnerChanged, LeadQualified, LeadConverted, DuplicateWarningRaised |
| account-service | create/update/archive account, create/update contact, create/link account/contact, proceed after duplicate warning | account/contact list/detail/search, account/contact summary, duplicate warning, archive eligibility | AccountCreated, ContactCreated, AccountArchived, DuplicateWarningRaised, OwnerChanged |
| opportunity-service | create/update opportunity, change stage, close won/lost, archive opportunity, change owner | opportunity list/detail/search, opportunity summary, closure eligibility, archive eligibility | OpportunityCreated, OpportunityStageChanged, OpportunityClosedWon, OpportunityClosedLost, OwnerChanged |
| commercial-service | create/update quote, change quote status, create/update contract, change contract status, create payment plan, record payment | quote/contract/payment list/detail, payment status summary, contract reminder eligibility | QuoteAccepted, ContractStatusChanged, PaymentRecorded, PaymentOverdue |
| work-service | create activity/note/task, update task status, transfer open work on owner change | activity/task list/detail, reminder list, active obligation list | WorkItemCreated, TaskStatusChanged, ReminderStateChanged, OpenWorkTransferred |
| audit-history-service | append history/log event through trusted internal contract | record history, admin operation log | HistoryEventAppended, OperationLogAppended |
| reporting-service | rebuild projection, refresh metric snapshot | team overview, basic sales reports | ReportProjectionUpdated |
| import-export-service | start import, start export, cancel run where safe | run status, row results, export metadata | ImportRunCompleted, ExportRunCompleted |

## Representative API Contracts

### Permission Check

```text
POST /internal/permissions/check

Request:
{
  "actorId": "user-id",
  "action": "opportunity.close_won",
  "resource": {
    "type": "opportunity",
    "id": "opportunity-id"
  },
  "context": {
    "ownerId": "user-id",
    "teamId": "single-team"
  },
  "correlationId": "..."
}

Response:
{
  "allowed": true,
  "scope": "owned | team | all",
  "denialCategory": null
}
```

### Lead Conversion

```text
POST /leads/{leadId}/convert

Request:
{
  "idempotencyKey": "...",
  "target": {
    "createOrLinkAccount": true,
    "contactInput": { ... },
    "opportunityInput": { ... }
  }
}

Success:
{
  "leadId": "...",
  "accountId": "...",
  "contactIds": ["..."],
  "opportunityId": "...",
  "status": "ConvertedToOpportunity"
}
```

Failure cases:

- invalid lead not restored
- lead already converted
- missing account/contact/opportunity required fields
- permission denied
- downstream account or opportunity service unavailable

### Editable Record Concurrency

All editable P0/P1 record detail DTOs must include a numeric `version` and an
`updatedAt` timestamp. Mutating commands must include `expectedVersion`.

```text
PATCH /opportunities/{opportunityId}

Request:
{
  "expectedVersion": 12,
  "changes": { ... }
}

Conflict:
{
  "correlationId": "...",
  "error": {
    "code": "VERSION_CONFLICT",
    "category": "conflict",
    "safeMessage": "The record changed after it was opened.",
    "latest": {
      "id": "opportunity-id",
      "version": 13,
      "updatedAt": "timestamp",
      "updatedByDisplay": "authorized-safe-display"
    }
  }
}
```

The frontend must reload or merge from the latest authoritative DTO. Services
must not silently overwrite stale edits.

### Archive Eligibility

```text
GET /{resourceType}/{resourceId}/archive-eligibility

Success:
{
  "resourceType": "account | opportunity | contract",
  "resourceId": "...",
  "canArchive": false,
  "recordVersion": 7,
  "obligations": [
    {
      "type": "open_task | pending_signature_contract | unpaid_payment | active_quote",
      "id": "...",
      "service": "work-service | commercial-service",
      "status": "Open | PendingSignature | Unpaid",
      "dueDate": "YYYY-MM-DD or null",
      "ownerDisplay": "authorized-safe-display",
      "blocking": true,
      "safeMessage": "..."
    }
  ]
}
```

```text
POST /{resourceType}/{resourceId}/archive

Request:
{
  "expectedVersion": 7,
  "reason": "required-safe-text"
}
```

Failure cases:

- `ARCHIVE_BLOCKED_ACTIVE_OBLIGATION`
- `VERSION_CONFLICT`
- `PERMISSION_DENIED`
- `TERMINAL_RECORD_READ_ONLY` where applicable

Archive attempts and successful archive commands must create history or
operation log events.

### Owner Transfer

```text
POST /{resourceType}/{resourceId}/owner-transfer

Request:
{
  "expectedVersion": 8,
  "newOwnerId": "user-id",
  "reason": "required-safe-text",
  "manualWorkReassignment": {
    "enabled": false,
    "reason": null
  }
}

Success:
{
  "resourceId": "...",
  "ownerId": "new-owner-id",
  "version": 9,
  "workTransferStatus": "Completed | PendingRetry | Failed",
  "transferredOpenWorkCount": 4
}
```

Owner transfer emits `OwnerChanged`. Work-service must transfer open tasks and
follow-ups by event or explicit internal command using an idempotency key.
Manual work reassignment exceptions require Administrator or Sales Manager
permission and a required reason.

### Close Opportunity Lost

```text
POST /opportunities/{opportunityId}/close-lost

Request:
{
  "idempotencyKey": "...",
  "expectedVersion": 10,
  "closeDate": "YYYY-MM-DD",
  "lostReason": {
    "code": "PRICE | COMPETITOR | NO_BUDGET | NO_DECISION | OTHER",
    "detail": "required when OTHER or policy requires detail"
  }
}

Success:
{
  "opportunityId": "...",
  "status": "Lost",
  "closedAt": "timestamp",
  "version": 11
}
```

Failure cases:

- missing `lostReason`
- already Won/Lost
- stale `expectedVersion`
- permission denied

After Won/Lost, opportunity-service must reject ordinary edit/stage commands
with `TERMINAL_RECORD_READ_ONLY`. Post-close notes/tasks are created through
work-service with a related opportunity reference.

### Duplicate Warning

Duplicate detection must normalize company names, contact names, phone, email,
province/city, tax or business identifiers where available, and whitespace/case
variants. The check may query lead-service and account-service through safe
lookup contracts only.

```text
POST /duplicate-checks

Request:
{
  "targetType": "lead | account | contact",
  "candidate": { ... }
}

Warning:
{
  "result": "PossibleDuplicate",
  "warningToken": "single-use-token",
  "normalizedFields": ["companyName", "phone", "province"],
  "matches": [
    {
      "type": "lead | account | contact",
      "id": "...",
      "matchStrength": "High | Medium | Low",
      "safeSummary": "authorized-safe-summary",
      "visible": true
    }
  ],
  "rules": ["COMPANY_PHONE_MATCH"]
}
```

Proceed-after-warning requires the warning token:

```text
POST /leads
{
  "idempotencyKey": "...",
  "proceedWarningToken": "...",
  "input": { ... }
}
```

Proceeding creates a new record only. It must not merge, overwrite, or relink
existing records automatically.

### Record Payment

```text
POST /contracts/{contractId}/payments

Request:
{
  "idempotencyKey": "...",
  "amount": "decimal-string",
  "paymentDate": "YYYY-MM-DD",
  "note": "optional"
}

Success:
{
  "paymentId": "...",
  "contractId": "...",
  "paymentStatus": "PartiallyPaid | Paid",
  "remainingAmount": "decimal-string"
}
```

Failure cases:

- zero or negative amount
- overpayment
- unauthorized actor
- contract not in payable state
- duplicate idempotency key mismatch

### Close Opportunity Won

```text
POST /opportunities/{opportunityId}/close-won

Request:
{
  "idempotencyKey": "...",
  "closeDate": "YYYY-MM-DD"
}

Success:
{
  "opportunityId": "...",
  "status": "Won",
  "closedAt": "timestamp"
}
```

The opportunity-service must verify full payment with commercial-service or an
accepted event-backed payment status projection before persisting Won.

### Reminder Query

```text
GET /reminders?businessDate=YYYY-MM-DD

Success:
{
  "timezone": "Asia/Shanghai",
  "businessDate": "YYYY-MM-DD",
  "rows": [
    {
      "id": "...",
      "sourceService": "work-service | commercial-service",
      "type": "task_due | task_overdue | contract_pending_signature | payment_due | payment_overdue",
      "relatedRecord": {
        "type": "lead | account | opportunity | contract",
        "id": "...",
        "display": "authorized-safe-display"
      },
      "ownerDisplay": "authorized-safe-display",
      "dueDate": "YYYY-MM-DD",
      "status": "DueToday | Overdue | PendingSignature",
      "priority": "P0 | P1 | normal",
      "version": 1
    }
  ]
}
```

Due/overdue calculation uses workspace timezone `Asia/Shanghai` and the
requested `businessDate`. Completed, cancelled, signed, terminated, fully paid,
archived, and unauthorized items are excluded.

### Report Metrics

```text
GET /reports/sales-overview?from=YYYY-MM-DD&to=YYYY-MM-DD&groupBy=owner|stage|province

Success:
{
  "scope": "team | all",
  "filters": {
    "archived": "active_default",
    "from": "YYYY-MM-DD",
    "to": "YYYY-MM-DD",
    "groupBy": "owner"
  },
  "currency": "CNY",
  "metrics": {
    "leadCount": 0,
    "opportunityCount": 0,
    "quoteAmount": "0.00",
    "contractAmount": "0.00",
    "paidAmount": "0.00",
    "receivableAmount": "0.00",
    "wonCount": 0,
    "lostCount": 0
  },
  "breakdowns": {
    "leadsByStatus": [],
    "opportunitiesByStage": [],
    "quotesByStatus": [],
    "contractsByStatus": [],
    "paymentsByStatus": []
  },
  "groups": []
}
```

`metrics` carries the flat KPI tiles. `breakdowns` carries the mandatory
BR-014 / BR-017 status-and-stage groupings. Each `breakdowns.*` array is a list
of group rows using the following DTO:

```text
GroupRow:
{
  "key": "string",        // canonical status/stage code, e.g. "Qualified", "Negotiation"
  "label": "string",      // display label for the key
  "count": 0,             // number of records in this group
  "amount": "0.00"        // summed amount for the dimension, "0.00" when count-only
}
```

Required `breakdowns` dimensions and amount semantics (faithful to BR-017):

- `leadsByStatus`: Lead records grouped by lead status. `amount` is `"0.00"`
  (count-only dimension).
- `opportunitiesByStage`: Opportunity records grouped by current stage. `amount`
  is the summed expected amount.
- `quotesByStatus`: Quote records grouped by status. `amount` is the summed quote
  amount.
- `contractsByStatus`: Contract records grouped by status. `amount` is the summed
  contract amount.
- `paymentsByStatus`: Payment records grouped by status. Each row additionally
  carries `dueAmount` and `paidAmount` so Payment-plan due amount and Actual
  payment paid amount are both representable per BR-017; `amount` mirrors the
  dimension-primary value (due amount).

The optional `groupBy=owner|stage|province` pivot populates the top-level
`groups[]` array (same `GroupRow` DTO) and does not replace the mandatory
`breakdowns`.

reporting-service must apply authorization before aggregation. Empty authorized
results return zero `metrics`, empty `breakdowns.*` arrays, and an empty
`groups[]`, not permission denial. Archived records are excluded unless an
authorized user applies an explicit archived filter.

### Import And Export Runs

Committed import/export scope covers leads, accounts, contacts, opportunities, quotes,
contracts, payments, tasks/activities, and notes only when the target service
has a supported domain command/query contract. Unsupported object types must be
rejected before row mutation.

CSV import row validation must include required fields, enum values, date and
amount formats, duplicate warning behavior, permission scope, and target service
state rules. Dangerous CSV formula cells beginning with `=`, `+`, `-`, `@`, tab,
or carriage return must be rejected or safely escaped for export.

Temporary upload/export files are retained for 24 hours after run completion or
failed run termination, then deleted by a scheduled cleanup job. The
import-export-service must expose deletion status in run metadata and must log
cleanup success/failure. G6 PSM must model the concrete cleanup job and storage
path before G8.

## Event Contract Requirements

Each event must include:

- event ID
- event type
- event version
- producer service
- aggregate type
- aggregate ID
- actor ID or system actor
- occurred at
- correlation ID
- causation ID where applicable
- safe summary payload
- acceptance ID references where applicable

Event payloads must minimize sensitive data. Amounts, contact details, and
before/after values require explicit security review.

## Error Codes

Required error categories:

| Category | Example Codes |
|---|---|
| authentication | UNAUTHENTICATED, SESSION_EXPIRED, SESSION_REVOKED, USER_DISABLED, AUTHZ_VERSION_STALE |
| permission | PERMISSION_DENIED, SCOPE_DENIED, ADMIN_ONLY |
| validation | REQUIRED_FIELD_MISSING, INVALID_DATE, INVALID_AMOUNT |
| business_rule | INVALID_TRANSITION, LEAD_ALREADY_CONVERTED, QUOTE_EXPIRED, OVERPAYMENT_BLOCKED, EARLY_WON_BLOCKED, ARCHIVE_BLOCKED_ACTIVE_OBLIGATION, TERMINAL_RECORD_READ_ONLY, LOST_REASON_REQUIRED |
| conflict | VERSION_CONFLICT, DUPLICATE_IDEMPOTENCY_KEY, DUPLICATE_WARNING |
| dependency | SERVICE_UNAVAILABLE, DEPENDENCY_TIMEOUT, SERVICE_AUTH_FAILED |
| operation | IMPORT_PARTIAL_FAILURE, EXPORT_FAILED |

## Acceptance Mapping

Each P0/P1 acceptance item maps to at least one service API or event contract.
The detailed contract-to-acceptance trace must be represented in MDA
traceability before G8.
