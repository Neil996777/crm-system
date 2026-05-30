# Frontend Backend Contract

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Purpose

This document defines the contract expectations between the web frontend,
gateway-bff, and backend services. It does not define visual design.

## Frontend Entry Pattern

The frontend talks to gateway-bff only. The frontend does not call internal
domain services directly.

```text
Web App -> gateway-bff -> internal services
```

The gateway-bff may aggregate safe response data, but final business authority
belongs to domain services.

## Required UI State Support

Backend contracts must support the UX/UI service-backed states:

| UI State | Backend Support |
|---|---|
| Loading | Request lifecycle and operation status for long-running runs. |
| Empty | Authorized empty response distinct from permission denial where needed. |
| Validation error | Field-level error codes and safe messages. |
| Permission denied | Safe denial category and return path. |
| Disabled action | Optional allowed-action hints, never authoritative alone. |
| Blocked transition | Business-rule error code and safe explanation. |
| Conflict/stale data | Conflict/version signal where concurrency control applies. |
| Partial failure | Row-level import result and safe summaries. |
| Read-only audit/history | Query-only contracts with no normal mutation API. |
| Long-running operation | Run status, progress, completion, and failure result. |
| Sensitive display | Masking/safe summary and authorization-aware DTOs. |
| Archived context | Active/default filter plus explicit archived filter semantics. |

## Common DTO Requirements

All frontend-facing DTOs must:

- include stable IDs
- include `version` and `updatedAt` for editable P0/P1 records
- avoid leaking unauthorized record existence
- avoid exposing restricted fields unless authorized
- distinguish display labels from authoritative state
- expose status/stage values from accepted domain enumerations
- include safe related-record summaries only when authorized

## Editable Record Contract

Editable record detail responses must include:

```text
{
  "id": "...",
  "version": 12,
  "updatedAt": "timestamp",
  "updatedByDisplay": "authorized-safe-display"
}
```

Mutations must send `expectedVersion`. A stale mutation returns
`VERSION_CONFLICT` with the latest safe summary. The UI must present a reload or
retry path and must not treat the stale save as successful.

## Pagination, Search, And Filtering

Entity list APIs must support:

- page or cursor contract
- search term
- basic filters
- archived/default active filter
- permission-filtered rows
- empty state response
- invalid filter response

The backend must apply authorization before returning rows.

## Mutation UX Contract

Write responses must support:

- success state
- created/updated record ID
- authoritative status/stage after save
- validation errors
- blocked business transition errors
- permission denial
- conflict/stale data response
- dependency failure response when downstream service is unavailable

The frontend must not infer success from optimistic UI alone.

## Import / Export UX Contract

Import/export must use long-running operation contracts:

- start run
- get run status
- get row results
- get export metadata
- get safe failure summary
- authorized object type list
- dangerous CSV formula handling result
- temporary file retention/deletion status

Import row errors must include row number, field, rule code, and safe message.
They must not leak unauthorized matched record details.

Export metadata responses must include requested object type, filters, row
count, generated timestamp, expiry/deletion timestamp, file safety treatment,
and operation log reference.

## Archive UX Contract

Archive attempts must support a pre-check and a blocked result:

```text
{
  "canArchive": false,
  "recordVersion": 7,
  "obligations": [
    {
      "type": "open_task",
      "id": "...",
      "service": "work-service",
      "status": "Open",
      "dueDate": "YYYY-MM-DD",
      "ownerDisplay": "authorized-safe-display",
      "blocking": true,
      "safeMessage": "..."
    }
  ]
}
```

The UI must be able to refresh eligibility after obligations are resolved.

## Duplicate Warning UX Contract

Duplicate warning responses must include safe match summaries and a
single-use `warningToken`. Proceeding after warning creates a new record and
does not merge, overwrite, or relink existing records automatically.

## Reminder Row Contract

Reminder rows must include source service, reminder type, related safe record
summary, owner display, due date, status, and version. Due and overdue states
use workspace timezone `Asia/Shanghai` and a supplied business date.

## History And Log UX Contract

History/log responses must include:

- event ID
- action
- actor display summary
- timestamp
- resource type
- safe before/after values where authorized
- correlation ID where useful for support or audit

History and logs are read-only through normal CRM UI.

The normal CRM UI must expose no edit/delete actions for history or operation
log records.

## Gateway Error Normalization

gateway-bff may normalize transport details, but must preserve:

- error code
- category
- field errors
- correlation ID
- safe message

It must not convert permission denial or business-rule blocks into generic
success or partial acceptance.
