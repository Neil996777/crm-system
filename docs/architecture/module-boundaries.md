# Module Boundaries

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Boundary Rule

Service boundaries are business capability boundaries. They are not page,
controller, table, repository, or technical-layer boundaries.

Each service may contain its own internal modules such as API handlers,
application commands, domain model, repository, event publisher, and adapters.
Other services may depend only on public contracts, not internal modules.

## Service Boundaries

| Service | Responsibilities | Non-Responsibilities | Owns | Depends On |
|---|---|---|---|---|
| gateway-bff | External API entry, frontend aggregation, request routing, correlation ID propagation, response shaping for web app. | Business rule authority, direct database access, bypassing domain services. | API edge routes, frontend-oriented DTO composition. | All backend service public APIs. |
| identity-authz-service | Authentication, sessions, users, roles, active/disabled state, permission decisions, last-admin protection. | Owning CRM records, replacing domain service business rules. | User, session, role, permission policy data. | audit-history-service for operation/security events. |
| lead-service | Lead create/edit/query, owner assignment state, qualification state, duplicate warning input for leads, conversion-once guard. | Owning customer/contact or opportunity data after conversion. | Lead aggregate, lead status, lead source, lead owner, lead conversion record. | identity-authz-service, account-service, opportunity-service, audit-history-service. |
| account-service | Company/customer and contact lifecycle, ownership, duplicate warning input for account/contact records. | Owning lead, opportunity, quote, contract, payment, task, or report data. | Account/customer aggregate, contact aggregate, account owner/scope. | identity-authz-service, audit-history-service. |
| opportunity-service | Opportunity lifecycle, stage transitions, Won/Lost terminal states, closure rules, expected amount/date. | Owning quote, contract, or payment records. | Opportunity aggregate, stage history pointer, closure data. | identity-authz-service, account-service, commercial-service, audit-history-service. |
| commercial-service | Quote lifecycle, contract lifecycle, payment plans, actual payments, payment status, amount integrity. | Owning opportunity state, lead/account state, task/reminder ownership. | Quote, contract, payment plan, actual payment aggregates. | identity-authz-service, opportunity-service, audit-history-service. |
| work-service | Activities, notes, follow-up tasks, task status, due/overdue reminder eligibility. | Owning target CRM records or commercial states. | Activity, note, task, reminder projection data. | identity-authz-service, record-owning services, audit-history-service. |
| audit-history-service | Record-local history, admin/global operation logs, append-only event storage, log query. | Replacing source service state, mutating business aggregates. | History event, operation log, security/audit event data. | identity-authz-service for admin log query authz context. |
| reporting-service | Team overview, basic sales reports, read model/projections for authorized metrics. | Directly querying source service databases, owning source truth. | Report read models, aggregate metrics snapshots. | identity-authz-service, domain events from source services. |
| import-export-service | CSV import/export run tracking, row-level results, export result metadata, long-running operation state. | Directly writing domain tables or bypassing domain validation. | Import/export run, row result, generated export metadata. | identity-authz-service, target domain service APIs, audit-history-service. |

## Dependency Direction

Allowed dependency types:

- Public HTTP API or future RPC contract.
- Published domain event contract.
- Shared contract/DTO package generated from architecture-approved contracts.
- Operational health check contract.

Forbidden dependency types:

- Direct import of another service's internal code.
- Direct read/write of another service's database tables.
- Shared repository or shared business-rule package across services.
- Generic database CRUD API used as a substitute for business contracts.

## Gateway Boundary

The gateway-bff is not a data owner for CRM domain records. It may:

- authenticate request context with identity-authz-service
- route commands to the owning service
- aggregate UI-oriented read views from public Query APIs
- attach and propagate correlation IDs
- normalize safe error responses

It may not:

- decide business state transitions
- write domain databases
- bypass service authorization
- return fields that target services did not authorize

## Shared Package Boundary

Shared packages may contain:

- API DTO schemas
- event schemas
- error codes
- permission action constants
- generated clients
- tracing/correlation helpers

Shared packages may not contain:

- domain aggregate methods
- repositories
- service-specific database models
- business-rule implementations
- permission bypass helpers

## Archive And Deletion Boundary

No service may hard-delete core CRM records in P0/P1 scope. Eligible records are
archived or transitioned through domain states. Archive behavior remains owned
by the record-owning service and must publish history/audit events.
