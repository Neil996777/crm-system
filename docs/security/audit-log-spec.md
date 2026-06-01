# Audit Log Spec

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input

## Audit Model

The CRM has two audit-related surfaces:

| Surface | Priority | Purpose | Visibility | Acceptance IDs |
|---|---|---|---|---|
| Record-local business history | P0 | Shows collaboration and business timeline for a related CRM record. | Users authorized for the related record. | ACC-014 |
| Admin/global operation log | P1 | Gives Administrator cross-record query for access-sensitive and operation-sensitive events. | Administrator only. | ACC-022 |

Both surfaces must be append-only through normal CRM actions. Normal CRM users
must not be able to edit or hard-delete history or operation-log events.
Architecture reset on 2026-05-29: implementation is blocked until the restarted
delivery flow passes G8. New architecture must preserve these append-only audit
requirements.

## Common Event Schema

| Field | Required | Description |
|---|---|---|
| event_id | Yes | Stable event type such as `EVT-STAGE-CHANGED`. |
| event_uid | Yes | Unique immutable event instance identifier. |
| occurred_at | Yes | Server-side timestamp. |
| actor_user_id | Yes when authenticated | User who performed the action. |
| actor_role | Yes when authenticated | Role at the time of action. |
| actor_display | Yes | Safe display name or system actor label. |
| action | Yes | Human-readable action name. |
| resource_type | Yes | Lead, Customer, Contact, Opportunity, Quote, Contract, Payment, Task, User, Import, Export, Auth, or Log. |
| resource_id | Yes when resource exists | Identifier of affected resource. |
| parent_resource_type | When applicable | Parent object type for child records. |
| parent_resource_id | When applicable | Parent object identifier for permission checks and traceability. |
| result | Yes | Success, denied, blocked, failed, or system-applied. |
| reason_code | When applicable | Validation, permission, business-rule, or security reason code. |
| before_summary | When applicable | Safe or authorized before state. |
| after_summary | When applicable | Safe or authorized after state. |
| diff_classification | When applicable | Internal, Confidential, Restricted, or Security Critical. |
| scope_summary | When applicable | Team, governed, owned/assigned, filters, archived inclusion, or import/export scope. |
| request_context | Architecture-defined | Non-secret request trace fields for audit correlation. |
| acceptance_ids | Yes | Related ACC IDs. |

Secrets, credentials, session tokens, raw password data, and raw authentication
factors must never be stored in audit events.

## Immutability Requirements

| ID | Requirement | Acceptance IDs |
|---|---|---|
| AUD-IMM-001 | Record-local history and admin/global operation logs cannot be edited through normal CRM actions. | ACC-014, ACC-022 |
| AUD-IMM-002 | Audit event creation must occur in the same durable workflow as the sensitive business change where applicable, so success is not shown without the required event. | ACC-014, ACC-016, ACC-022 |
| AUD-IMM-003 | Denied or blocked sensitive operations must not mutate business data, but may create operation-log events. | ACC-002, ACC-022 |
| AUD-IMM-004 | Architecture must define append-only or tamper-evident storage behavior for operation logs. | ACC-022 |
| AUD-IMM-005 | Audit-log access itself is Administrator-only for global logs and must not expose logs to Sales Manager or Sales. | ACC-022 |

## Event Catalog

| Event ID | Priority | Surface | Trigger | Actor | Resource | Required Fields / Notes | Acceptance IDs |
|---|---|---|---|---|---|---|---|
| EVT-AUTH-LOGIN-SUCCEEDED | P0/P1 | Admin/global operation log | User signs in successfully | User | Auth | actor, role, occurred_at, result | ACC-001, ACC-022 |
| EVT-AUTH-LOGIN-FAILED | P0/P1 | Admin/global operation log | Sign-in fails | Unknown or user if known safely | Auth | safe reason code; no credential details | ACC-001, ACC-022 |
| EVT-AUTH-ACCESS-DENIED | P0/P1 | Admin/global operation log where applicable | Protected access denied | User or unauthenticated | Resource or route | safe resource type; avoid restricted names for non-admin display | ACC-002, ACC-022 |
| EVT-USER-ROLE-CHANGED | P0/P1 | Admin/global operation log | Administrator changes user role | Administrator | User | target user, old/new role, result | ACC-001, ACC-002, ACC-022 |
| EVT-USER-STATUS-CHANGED | P0/P1 | Admin/global operation log | Administrator enables/disables user | Administrator | User | target user, old/new status, result | ACC-001, ACC-002, ACC-022 |
| EVT-LAST-ADMIN-BLOCKED | P0/P1 | Admin/global operation log | Last Administrator downgrade/disable blocked | Administrator | User | target user, blocked reason | ACC-001, ACC-002, ACC-022 |
| EVT-OWNER-CHANGED | P0 | Record-local history and operation log where applicable | Owner assignment or transfer | Administrator or Sales Manager | Lead or parent CRM record | old/new owner, affected open task transfer summary | ACC-003, ACC-014, ACC-022 |
| EVT-LEAD-QUALIFIED | P0 | Record-local history | Lead marked Valid | Sales or Sales Manager | Lead | qualification result | ACC-004, ACC-014 |
| EVT-LEAD-DISQUALIFIED | P0 | Record-local history | Lead marked Invalid | Sales or Sales Manager | Lead | invalid reason | ACC-004, ACC-014 |
| EVT-LEAD-CONVERTED | P0 | Record-local history | Lead converted to customer/opportunity context | Sales or Sales Manager | Lead | related customer/opportunity ids | ACC-004, ACC-005, ACC-007, ACC-014 |
| EVT-STAGE-CHANGED | P0 | Record-local history and operation log where applicable | Opportunity stage changes | Sales or Sales Manager | Opportunity | old/new stage, required data summary | ACC-008, ACC-014, ACC-022 |
| EVT-STATUS-CHANGED | P0 | Record-local history and operation log where applicable | Status changes on quote, contract, payment, task, or other CRM record | Authorized user or system | CRM record | old/new status, result | ACC-009 to ACC-014, ACC-022 |
| EVT-QUOTE-ACCEPTED | P0/P1 | Record-local history and operation log | Quote accepted | Sales or Sales Manager | Quote | opportunity, quote id | ACC-009, ACC-014, ACC-022 |
| EVT-CONTRACT-SIGNED | P0/P1 | Record-local history and operation log | Contract becomes Signed | Sales or Sales Manager | Contract | signed/effective date, amount summary | ACC-010, ACC-014, ACC-022 |
| EVT-CONTRACT-TERMINATED | P0/P1 | Record-local history and operation log | Contract terminated | Sales Manager or Administrator | Contract | termination reason | ACC-010, ACC-014, ACC-022 |
| EVT-PAYMENT-RECORDED | P0/P1 | Record-local history and operation log | Actual payment recorded | Sales or Sales Manager | Payment | amount summary, contract id, status result | ACC-011, ACC-014, ACC-022 |
| EVT-PAYMENT-OVERDUE | P0/P1 | Record-local history and operation log where applicable | Payment becomes overdue | System or authorized user | Payment | due date, unpaid summary | ACC-011, ACC-021, ACC-022 |
| EVT-TASK-COMPLETED | P0 | Record-local history | Task completed | Task owner, Sales Manager, Administrator | Task | completion timestamp | ACC-012, ACC-014 |
| EVT-TASK-CANCELLED | P0 | Record-local history | Task cancelled | Task owner, Sales Manager, Administrator | Task | cancellation reason | ACC-012, ACC-014 |
| EVT-OPPORTUNITY-WON | P0 | Record-local history and operation log where applicable | Opportunity closed Won | Sales or Sales Manager | Opportunity | closure evidence summary (related contract Signed) | ACC-013, ACC-014, ACC-022 |
| EVT-OPPORTUNITY-LOST | P0 | Record-local history and operation log where applicable | Opportunity closed Lost | Sales or Sales Manager | Opportunity | lost reason | ACC-013, ACC-014, ACC-022 |
| EVT-RECORD-ARCHIVED | P0/P1 | Record-local history and operation log | Eligible record archived | Administrator or Sales Manager | CRM record | archive target, downstream obligation result | ACC-002, ACC-014, ACC-022 |
| EVT-IMPORT-RUN | P1 | Admin/global operation log | CSV import run completes or fails | Administrator or Sales Manager | Import | object type, scope, total rows, success count, failure count, result | ACC-020, ACC-022 |
| EVT-EXPORT-RUN | P1 | Admin/global operation log | CSV export run completes or fails | Administrator or Sales Manager | Export | object type, filters, archived inclusion, exported count, result | ACC-020, ACC-022 |
| EVT-REPORT-ACCESS-DENIED | P1 | Admin/global operation log where applicable | Unauthorized report access denied | User | Report | report type, denied result | ACC-018, ACC-023 |

## Query Requirements

| Surface | Query Requirement |
|---|---|
| Record-local history | Query by related record; visible according to record permission; ordered by occurred_at; supports event type filtering where useful. |
| Admin/global operation log | Administrator-only query by event ID, actor, resource type, date range, result, import/export run, access failure, owner/stage/status/payment/archive categories. |
| Reports and exports | Must not query unauthorized records first and filter only in the UI; authorization applies at the data source. |

## Testability Requirements

- Every required event in the catalog must have at least one positive test for
  event creation.
- Permission-denied cases must prove no restricted data exposure and no
  business mutation.
- Sensitive successful mutations must prove both data change and required audit
  event.
- Operation logs must be inaccessible to Sales Manager and Sales.
- Record-local history must be inaccessible for unauthorized records.
- Import/export tests must verify counts, result status, and operation-log
  event creation without storing sensitive sample rows in generic output.
