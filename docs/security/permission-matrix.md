# Permission Matrix

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input

## Matrix Rules

- Decision values are `Allow` or `Deny`.
- `Scope` is part of the authorization decision and must be enforced by the
  backend.
- Frontend hidden or disabled actions do not satisfy authorization.
- Denied access must not expose or mutate restricted data.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8. New architecture must preserve these backend authorization rules.

## Resource Scope Definitions

| Scope | Definition |
|---|---|
| Governed records | All CRM records in the single committed team workspace. |
| Team records | All CRM records in the single committed sales team workspace. |
| Owned/assigned records | Records owned by or assigned to the Sales user. |
| Related child records | Contacts, quotes, contracts, payments, activities, notes, tasks, and history linked to an owned/assigned parent record. |
| Archived records | Non-deleted records removed from active/default views and available only through explicit archived filters, history, logs, or audit/report context. |
| Global operation logs | Administrator-only cross-record operation and security event query. |

## Permission Matrix

| ID | Priority | Actor | Action | Resource | Condition | Decision | Audit Requirement | Acceptance IDs |
|---|---|---|---|---|---|---|---|---|
| PM-001 | P0 | Unauthenticated | Access | Any protected CRM route, API, record, report, import/export, or log | No authenticated session | Deny | Log protected access failure where applicable without sensitive details. | ACC-001, ACC-002 |
| PM-002 | P0 | Disabled user | Access | Any protected CRM data | Account disabled | Deny | Log authentication or access failure. | ACC-001, ACC-002 |
| PM-003 | P0 | Administrator | Manage | Users and roles | Authenticated active Administrator | Allow | Operation log for create/update/disable/role change. | ACC-001, ACC-002, ACC-022 |
| PM-004 | P0 | Sales Manager | Manage | Users and roles | Any condition | Deny | Log denied attempt when routed to backend. | ACC-002 |
| PM-005 | P0 | Sales | Manage | Users and roles | Any condition | Deny | Log denied attempt when routed to backend. | ACC-002 |
| PM-006 | P0 | Administrator | Change role/status | User account | Target account exists; change does not remove last active Administrator | Allow | Operation log with target user, old/new values, result. | ACC-001, ACC-002, ACC-022 |
| PM-007 | P0 | Administrator | Disable or downgrade | Last active Administrator | Change would leave no active Administrator | Deny | Operation log with blocked result. | ACC-001, ACC-002, ACC-022 |
| PM-008 | P0 | Administrator | Create/view/edit/search/filter | Lead, company/customer, contact, opportunity, quote, contract, payment, activity, note, task | Governed records | Allow | Record-local history for business changes; operation log where required. | ACC-002 to ACC-016 |
| PM-009 | P0 | Sales Manager | Create/view/edit/search/filter | Lead, company/customer, contact, opportunity, quote, contract, payment, activity, note, task | Team records | Allow | Record-local history for business changes; operation log where required. | ACC-002 to ACC-016 |
| PM-010 | P0 | Sales | Create | Lead | Authenticated Sales | Allow | Record-local history once persisted. | ACC-003, ACC-014 |
| PM-011 | P0 | Sales | View/edit/search/filter | Lead | Owned/assigned lead | Allow | Record-local history for mutations. | ACC-003, ACC-014, ACC-015 |
| PM-012 | P0 | Sales | View/edit/search/filter | Lead | Non-owned/non-assigned lead | Deny | No restricted data exposure; log denied detail access where applicable. | ACC-002, ACC-003, ACC-015 |
| PM-013 | P0 | Sales | Qualify/edit/convert | Unassigned lead | Lead has no owner | Deny | No mutation; safe denied feedback. | ACC-003, ACC-004 |
| PM-014 | P0 | Administrator, Sales Manager | Assign/transfer | Lead or team-owned parent record | Governed/team scope | Allow | Record-local owner-change history; operation log where applicable. | ACC-002, ACC-003, ACC-014, ACC-022 |
| PM-015 | P0 | Sales | Assign/transfer | Lead or parent record | Any condition | Deny | Denied mutation creates no owner change. | ACC-002, ACC-003 |
| PM-016 | P0 | Sales | Create/view/edit | Company/customer or contact | Related to owned/assigned record or created by user in allowed workflow | Allow | Record-local history for mutations. | ACC-005, ACC-006, ACC-014 |
| PM-017 | P0 | Sales | View/edit | Company/customer or contact | No owned/assigned relation | Deny | No existence or sensitive value exposure. | ACC-002, ACC-005, ACC-006 |
| PM-018 | P0 | Sales | Create/view/edit/change stage/close | Opportunity | Owned/assigned opportunity and business rules pass | Allow | Record-local history for stage/status/closure. | ACC-007, ACC-008, ACC-013, ACC-014 |
| PM-019 | P0 | Sales | View/edit/change stage/close | Opportunity | Non-owned/non-assigned opportunity | Deny | No mutation or restricted data exposure. | ACC-002, ACC-007, ACC-008 |
| PM-020 | P0 | Sales | Create/view/edit/status change | Quote | Related opportunity is owned/assigned and quote rules pass | Allow | Record-local history and operation log for accepted quote. | ACC-009, ACC-014, ACC-022 |
| PM-021 | P0 | Sales | Create/view/edit/status change | Contract | Related opportunity/contract is owned/assigned and contract rules pass | Allow | Record-local history and operation log for signature, termination, status changes. | ACC-010, ACC-014, ACC-022 |
| PM-022 | P0 | Sales | Create/view/edit/status change | Payment plan or actual payment | Related contract is owned/assigned and payment rules pass | Allow | Record-local history and operation log for actual payment and overdue events. | ACC-011, ACC-014, ACC-022 |
| PM-023 | P0 | Sales | Create/view/edit | Activity, note, task | Related record is owned/assigned; task ownership rule passes | Allow | Record-local history for relevant changes. | ACC-012, ACC-014, ACC-021 |
| PM-024 | P0 | Sales | View | Record-local history | Related record is owned/assigned | Allow | No edit allowed through normal CRM actions. | ACC-014 |
| PM-025 | P0 | Sales | View | Record-local history | Related record is not owned/assigned | Deny | No restricted event detail exposure. | ACC-002, ACC-014 |
| PM-026 | P0/P1 | Administrator | Archive | Eligible CRM record | Governed record; downstream obligations resolved or archived | Allow | Record-local history and operation log. | ACC-002, ACC-014, ACC-022 |
| PM-027 | P0/P1 | Sales Manager | Archive | Eligible CRM record | Team record; downstream obligations resolved or archived | Allow | Record-local history and operation log. | ACC-002, ACC-014, ACC-022 |
| PM-028 | P0 | Sales | Archive | Any core CRM record | Any condition | Deny | Denied mutation creates no archive state. | ACC-002 |
| PM-029 | P0 | Any role | Hard delete | Core CRM record | Any condition | Deny | Attempts are unavailable or rejected; no normal hard-delete log required unless backend receives a denied attempt. | ACC-002, ACC-016 |
| PM-030 | P0/P1 | Administrator | View | Archived records | Explicit archived filter, history, log, audit, or report context | Allow | Query result follows privacy and audit rules. | ACC-014, ACC-015, ACC-023 |
| PM-031 | P0/P1 | Sales Manager | View | Archived records | Explicit archived filter and team scope | Allow | Query result follows privacy and audit rules. | ACC-014, ACC-015, ACC-023 |
| PM-032 | P0/P1 | Sales | View | Archived records | Owned/assigned related record and explicit archived filter or history context | Allow | No global archive browsing. | ACC-014, ACC-015 |
| PM-033 | P0/P1 | Sales | View | Archived records | Non-owned/non-assigned record | Deny | No restricted data exposure. | ACC-002, ACC-014, ACC-015 |
| PM-034 | P1 | Administrator | Import | CSV CRM records | Governed records; row validation and related references pass | Allow | Operation log with object, scope, counts, result. | ACC-020, ACC-022 |
| PM-035 | P1 | Sales Manager | Import | CSV CRM records | Team records; row validation, permissions, and related references pass | Allow | Operation log with object, scope, counts, result. | ACC-020, ACC-022 |
| PM-036 | P1 | Sales | Import | CSV CRM records | Any condition | Deny | Denied access event where applicable. | ACC-020 |
| PM-037 | P1 | Administrator | Export | CSV CRM records | Governed records; explicit export confirmation | Allow | Operation log with filters, archived inclusion, count, result. | ACC-020, ACC-022 |
| PM-038 | P1 | Sales Manager | Export | CSV CRM records | Team records; explicit export confirmation | Allow | Operation log with filters, archived inclusion, count, result. | ACC-020, ACC-022 |
| PM-039 | P1 | Sales | Export | CSV CRM records | Any condition | Deny | Denied access event where applicable. | ACC-020 |
| PM-040 | P1 | Administrator | View/query | Global operation logs | Authenticated active Administrator | Allow | Query is read-only; log detail follows privacy rules. | ACC-022 |
| PM-041 | P1 | Sales Manager | View/query | Global operation logs | Any condition | Deny | Denied access event where applicable. | ACC-022 |
| PM-042 | P1 | Sales | View/query | Global operation logs | Any condition | Deny | Denied access event where applicable. | ACC-022 |
| PM-043 | P1 | Administrator | View/query | Basic reports | Governed records | Allow | Reports use persisted authorized records only. | ACC-023 |
| PM-044 | P1 | Sales Manager | View/query | Basic reports and team overview | Team records | Allow | Reports use persisted authorized records only. | ACC-018, ACC-023 |
| PM-045 | P1 | Sales | View/query | Manager/admin reports or team overview | Any condition | Deny | No aggregate or drill-in leakage. | ACC-018, ACC-023 |
| PM-046 | P1 | Administrator, Sales Manager, Sales | View | Reminders | Authorized due/overdue related records only | Allow | Reminder visibility follows record permissions. | ACC-021 |
| PM-047 | P1 | Any authenticated role | View | Reminders | Unauthorized related records | Deny | Hidden from reminder list; no existence leakage. | ACC-002, ACC-021 |
| PM-048 | P1 | Administrator, Sales Manager, Sales | Receive duplicate warning | Company/contact/lead create or edit | Actor is authorized to create/edit the current record | Allow | Warning is not a merge and does not expose unauthorized matched record details. | ACC-019 |

## Required Permission Test Patterns

- Positive tests for each allowed role/action/resource condition.
- Negative tests for each denied role/action/resource condition.
- IDOR-style tests using direct record identifiers outside actor scope.
- Mutation tests proving denied actions create no business-state change.
- Visibility tests proving unauthorized records are hidden from list, search,
  filter, reminder, report, import/export, and log contexts.
- Audit tests proving sensitive allowed and blocked operations create the
  required history or operation-log events.
