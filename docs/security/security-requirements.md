# Security Requirements

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input
- Source:
  - `docs/product/prd.md`
  - `docs/product/acceptance-matrix.md`
  - `docs/business/business-rules.md`
  - `docs/business/role-permission-scenarios.md`
  - `docs/business/edge-cases.md`
  - `docs/ux-ui/ui-spec.md`
  - `docs/ux-ui/component-spec.md`

## Scope

This document defines security requirements for authentication, authorization,
data visibility, sensitive operations, privacy-sensitive handling, auditability,
and abuse-case coverage.

Security may clarify or strengthen P0/P1 product behavior. It must not
downgrade, delete, merge away, weaken, or accept partial P0/P1 behavior.
Frontend visibility is UX guidance only; protected data and protected actions
require backend authorization enforcement. Implementation remains blocked until
G8 passes.

## Security Requirement Index

| ID | Priority | Requirement | Primary Acceptance IDs | Downstream Owner |
|---|---|---|---|---|
| SEC-001 | P0 | Authenticate every CRM user before protected route, API, data, report, import, export, or log access. | ACC-001, ACC-002 | Architecture, Backend, QA |
| SEC-002 | P0 | Enforce the three-role model: Administrator, Sales Manager, and Sales. | ACC-001, ACC-002 | Architecture, Backend, QA |
| SEC-003 | P0 | Enforce record-level visibility for Sales owned/assigned scope and related child records. | ACC-002, ACC-003 to ACC-015 | Architecture, Backend, QA |
| SEC-004 | P0 | Deny unauthorized access without exposing restricted record names, existence, sensitive values, or before/after values. | ACC-002, ACC-015 | Frontend, Backend, QA |
| SEC-005 | P0 | Apply backend authorization to every protected create, read, update, assign, transfer, archive, close, import, export, report, and log query action. | ACC-002 | Architecture, Backend, QA |
| SEC-006 | P0 | Preserve no-hard-delete behavior for core CRM records in normal product workflows. | ACC-002, ACC-005, ACC-016 | Architecture, Backend, Audit |
| SEC-007 | P0 | Require explicit authorization, confirmation, validation, and audit events for sensitive operations. | ACC-002, ACC-014, ACC-022 | Architecture, Frontend, Backend, QA |
| SEC-008 | P0 | Protect Administrator user and role management, including last-Administrator protection. | ACC-001, ACC-002, ACC-022 | Architecture, Backend, QA, Audit |
| SEC-009 | P0 | Record-local history must be append-only through normal CRM workflows and visible only by record permission. | ACC-014 | Architecture, Backend, QA, Audit |
| SEC-010 | P1 | Admin/global operation logs must be Administrator-only, queryable, append-only through normal CRM workflows, and testable. | ACC-022 | Architecture, Backend, QA, Audit |
| SEC-011 | P1 | CSV import must validate authorization, required fields, related references, dangerous content, and row-level failure behavior before mutation. | ACC-020 | Architecture, Backend, QA |
| SEC-012 | P1 | CSV export must include only authorized records, require explicit confirmation, and create an export operation-log event. | ACC-020, ACC-022 | Architecture, Backend, QA, Audit |
| SEC-013 | P1 | Basic reports must use authorized persisted records only and must not leak unauthorized aggregates. | ACC-018, ACC-023 | Architecture, Backend, QA |
| SEC-014 | P0/P1 | Sensitive data display, errors, logs, import results, and reports must follow data classification and masking rules. | ACC-002, ACC-014, ACC-020, ACC-022, ACC-023 | Frontend, Backend, QA, Audit |
| SEC-015 | P0/P1 | Retention, archive, and deletion boundaries must follow the v1 privacy requirements. | ACC-002, ACC-014, ACC-016, ACC-022 | Architecture, Backend, Audit |
| SEC-016 | P0/P1 | Authentication and authorization failures must be logged at the appropriate level without exposing sensitive details to unauthorized users. | ACC-001, ACC-002, ACC-022 | Architecture, Backend, QA |
| SEC-017 | P0/P1 | Security controls must be verifiable with positive, negative, abuse-case, and audit-log tests. | ACC-001 to ACC-023 | QA, Integration, Audit |
| SEC-018 | P0 | No P0/P1 core CRM path may rely on mock, static-only, in-memory-only, or non-persistent behavior. | ACC-016, ACC-017 | Architecture, QA, Audit |

## Authentication Requirements

| ID | Requirement | Acceptance IDs |
|---|---|---|
| AUTH-001 | Protected CRM routes and APIs require an authenticated session. | ACC-001, ACC-002 |
| AUTH-002 | Invalid credentials, disabled accounts, unavailable accounts, and unavailable authentication states return a unified sign-in failure message to unauthenticated users. | ACC-001 |
| AUTH-003 | Disabled users cannot access protected CRM data, reports, import/export, user management, or logs. | ACC-001, ACC-002 |
| AUTH-004 | Authentication state must be rechecked on protected API requests, not only at frontend route entry. | ACC-001, ACC-002 |
| AUTH-005 | Session invalidation, logout, and expired-session behavior must prevent further protected data access. | ACC-001, ACC-002 |
| AUTH-006 | Login failures and protected access failures create security-relevant operation-log events without leaking sensitive credential details. | ACC-001, ACC-022 |

## Authorization Requirements

| ID | Requirement | Acceptance IDs |
|---|---|---|
| AUTHZ-001 | Administrator can manage users, roles, governed CRM records, archived records, reports, import/export, and global operation logs. | ACC-001, ACC-002, ACC-020, ACC-022, ACC-023 |
| AUTHZ-002 | Sales Manager can view and manage all team CRM records, team assignments, eligible team archives, team import/export, team overview, and team reports. | ACC-002, ACC-018, ACC-020, ACC-023 |
| AUTHZ-003 | Sales can create and manage owned/assigned CRM records and related child records only. | ACC-002 to ACC-015 |
| AUTHZ-004 | Sales cannot view non-owned/non-assigned CRM records, team overview, manager/admin reports, import/export, global operation logs, user management, or archive actions. | ACC-002, ACC-018, ACC-020, ACC-022, ACC-023 |
| AUTHZ-005 | Sales Manager and Sales cannot manage users or roles. | ACC-001, ACC-002 |
| AUTHZ-006 | Sales Manager and Sales cannot view admin/global operation logs. | ACC-022 |
| AUTHZ-007 | Authorization must evaluate actor, action, resource, record owner/assignee, related parent scope, archived state, terminal state, and business rule conditions. | ACC-002 to ACC-023 |
| AUTHZ-008 | Denied mutation attempts must not change data, write partial business state, or create misleading success feedback. | ACC-002, ACC-016 |
| AUTHZ-009 | Administrator and Sales Manager actions are recorded as their own actor identity; the system must not silently act on behalf of Sales users. | ACC-014, ACC-022 |

## Sensitive Operations

Sensitive operations require backend authorization, UI confirmation where
applicable, validation, durable mutation only after validation succeeds, and
record-local history or admin/global operation-log events as specified.

| Operation | Priority | Confirmation | Required Audit Behavior | Acceptance IDs |
|---|---|---|---|---|
| User role/status change | P0 | Required before save | Admin/global operation log; include actor, target user, old/new role or status, result, and reason when available. | ACC-001, ACC-002, ACC-022 |
| Grant or remove Administrator | P0 | Required before save | Admin/global operation log; blocked if it would leave no active Administrator. | ACC-001, ACC-002, ACC-022 |
| Owner assignment or transfer | P0 | Required for transfer | Record-local history and operation log where applicable. | ACC-003, ACC-014, ACC-022 |
| Opportunity terminal Won/Lost closure | P0 | Required | Record-local history; operation log where applicable. | ACC-013, ACC-014 |
| Quote acceptance | P0 | Required when changing accepted quote | Record-local history and operation log. | ACC-009, ACC-014, ACC-022 |
| Contract signature, activation, completion, termination | P0 | Required for termination and high-impact status change | Record-local history and operation log. | ACC-010, ACC-014, ACC-022 |
| Actual payment record creation or status update | P0 | Required for high-impact payment change | Record-local history and operation log. | ACC-011, ACC-014, ACC-022 |
| Archive eligible record | P0/P1 | Required | Record-local history and operation log; no hard delete. | ACC-002, ACC-014, ACC-022 |
| CSV import run | P1 | Required before mutation | Operation log with object type, actor, scope, result counts, and failure counts. | ACC-020, ACC-022 |
| CSV export run | P1 | Required before export | Operation log with object type, filters, archived inclusion, estimated/exported count, scope, and result. | ACC-020, ACC-022 |
| Admin/global operation-log query | P1 | No destructive confirmation | Access is Administrator-only; query failures logged as security-relevant when suspicious. | ACC-022 |

## Secure Error And Display Requirements

- Permission denied states must not reveal restricted record names, existence,
  field values, owner names, contact details, payment values, or sensitive
  before/after values.
- Form validation may identify the invalid field and validation rule, but it
  must not expose unauthorized related-record data.
- Import row errors must prefer row number, field name, validation rule, and
  safe summaries. Full customer, contact, contract, payment, or log-sensitive
  raw values are not displayed by default.
- Export confirmation and export result states must not show unnecessary
  sensitive sample data.
- History and operation-log details may show before/after values only when the
  actor is authorized for the record or is an Administrator viewing global logs,
  and only according to the privacy data classification rules.
- Reports must exclude unauthorized records from both rows and aggregates.

## Downstream Verification Requirements

- Architecture must define the authn/authz architecture, enforcement points,
  data isolation conditions, audit storage, import/export security, retention,
  and backup handling without weakening these requirements.
- QA must create positive and negative permission tests for every P0/P1
  permission class, plus abuse-case tests for the scenarios in
  `docs/security/abuse-cases.md`.
- Integration must prove permissions, persistence, history, operation logs,
  import/export, and reports across frontend, backend, and persistent services.
- Audit must verify that P0/P1 acceptance items are not marked Done without
  implementation evidence, QA evidence, integration evidence, audit pass, and
  no open P0/P1 blocker.
