# User Scenarios

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Scenario Index

| ID | Priority | Role | Scenario | Goal | Acceptance IDs | Status |
|---|---|---|---|---|---|---|
| SCN-001 | P0 | Sales | Create and qualify a new lead | Convert valid sales need into CRM context | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007 | Accepted as Architecture Input |
| SCN-002 | P0 | Sales | Manage opportunity through quote and contract | Preserve pipeline, quote, and contract history | ACC-007, ACC-008, ACC-009, ACC-010, ACC-014 | Accepted as Architecture Input |
| SCN-003 | P0 | Sales | Track payment and close opportunity | Close only with valid payment or lost reason | ACC-011, ACC-013, ACC-014 | Accepted as Architecture Input |
| SCN-004 | P0/P1 | Sales | Manage follow-up work | Keep tasks, notes, and reminders actionable | ACC-012, ACC-021 | Accepted as Architecture Input |
| SCN-005 | P0/P1 | Sales Manager | Review and coordinate team pipeline | Manage team records, risks, and assignments | ACC-002, ACC-018, ACC-023 | Accepted as Architecture Input |
| SCN-006 | P0/P1 | Administrator | Govern users, access, and logs | Maintain operational control and auditability | ACC-001, ACC-002, ACC-022 | Accepted as Architecture Input |
| SCN-007 | P1 | Administrator, Sales Manager | Import and export CRM data | Bulk load or extract authorized records safely | ACC-020 | Accepted as Architecture Input |
| SCN-008 | P1 | Sales Manager, Administrator | Review basic sales reports | Understand persisted team or governed sales state | ACC-018, ACC-023 | Accepted as Architecture Input |

## Scenario Details

### SCN-001: Sales Creates And Qualifies A Lead

Preconditions:
- Sales is authenticated.
- Sales has permission to create a lead.

Main path:
1. Sales creates a lead with source, status, and lead name or company name.
2. Sales records company/contact details and need summary where available.
3. Sales qualifies the lead as Valid.
4. Sales creates or links customer and contact records.
5. Sales creates an opportunity from the qualified need.

Expected result:
- Lead, customer/contact, opportunity, and lead conversion history are
  persisted.

Failure paths:
- Missing required fields block save.
- Sales cannot qualify Unassigned leads.
- Invalid lead cannot convert to opportunity.

### SCN-002: Sales Manages Opportunity Through Quote And Contract

Preconditions:
- Sales owns or is assigned to the opportunity.
- Related customer exists.

Main path:
1. Sales creates opportunity with expected amount and expected close date.
2. Sales moves through allowed pipeline stages.
3. Sales creates a quote and sends it.
4. Sales marks one quote Accepted.
5. Sales creates Pending Signature contract from the Accepted quote.
6. Sales records expected signed date and required contract note.
7. Sales signs or activates contract when signed/effective date is available.

Expected result:
- Opportunity, quote, and contract records are persisted with history.

Failure paths:
- Forbidden stage transitions are rejected.
- Expired quote cannot be linked to new contract.
- Contract amount mismatch requires difference reason.

### SCN-003: Sales Tracks Payment And Closes Opportunity

Preconditions:
- Contract exists and user is authorized.

Main path:
1. Sales creates payment plan.
2. Sales records actual payment.
3. Partial payment sets payment status to Partially Paid.
4. Full payment sets payment status to Paid.
5. Sales closes opportunity as Won after full payment.

Alternative path:
- Sales closes opportunity as Lost with lost reason before terminal Won.

Failure paths:
- Zero, negative, or overpayment amounts are rejected.
- Won is rejected without full payment.
- Won/Lost cannot be reopened in v1.

### SCN-004: Sales Manages Follow-Up Work

Preconditions:
- Related CRM record exists and Sales is authorized.

Main path:
1. Sales creates activity, note, or task on a related CRM record.
2. Sales sets task due date and owner.
3. Sales receives in-app reminders for due/overdue tasks.
4. Sales completes or cancels task.

Alternative reminders:
- Pending Signature contracts past expected signed date create reminders.
- Due/overdue payments create reminders.

Failure paths:
- Completed/cancelled tasks do not create active reminders.
- Unauthorized related records are hidden.

### SCN-005: Sales Manager Reviews And Coordinates Team Pipeline

Preconditions:
- Sales Manager is authenticated.
- Team records exist.

Main path:
1. Sales Manager opens team overview.
2. Sales Manager reviews leads, opportunities, quotes, contracts, payments,
   tasks, and pipeline status.
3. Sales Manager opens details to inspect history and risks.
4. Sales Manager assigns or transfers team work.

Expected result:
- Team work is visible and actionable according to manager permissions.

Failure paths:
- Sales users cannot access manager overview.
- Owner transfer preserves or transfers open tasks according to business rule.

### SCN-006: Administrator Governs Users, Access, And Logs

Preconditions:
- Administrator is authenticated.

Main path:
1. Administrator manages users and roles.
2. Administrator reviews governed CRM records.
3. Administrator opens global operation logs.
4. Administrator reviews access failures and key CRM operation events.

Expected result:
- Governance and audit-sensitive activity are visible to Administrator.

Failure paths:
- Sales Manager and Sales cannot access global operation logs.
- Logs cannot be edited through normal CRM actions.

### SCN-007: Administrator Or Sales Manager Imports And Exports CRM Data

Preconditions:
- User is Administrator or Sales Manager.
- CSV data or export target is available.

Main path:
1. User imports CSV.
2. System validates required fields row by row.
3. Valid rows are imported.
4. Invalid rows are reported.
5. User exports authorized records.

Failure paths:
- Unsupported file formats are rejected.
- Invalid rows do not corrupt existing records.
- Sales users cannot import/export in v1.

### SCN-008: Sales Manager Or Administrator Reviews Basic Sales Reports

Preconditions:
- User is Administrator or Sales Manager.
- Persisted CRM records exist or report state is empty.

Main path:
1. User opens reports.
2. Reports show counts and sums by committed groupings.
3. Report data respects authorization.

Failure paths:
- Sales users cannot access manager/admin reports.
- Empty data returns zero or empty report state.
