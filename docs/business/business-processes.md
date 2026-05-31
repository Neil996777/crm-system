# Business Processes

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Process Principles

- The committed business loop starts at lead entry and ends at won/lost opportunity
  closure with preserved history.
- P0/P1 process behavior must stay traceable to product acceptance IDs.
- Business processes may clarify or strengthen P0/P1 behavior, but must not
  downgrade, delete, merge away, or weaken it.
- No process may rely on mock, static-only, TODO, in-memory-only, or
  non-persistent behavior.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8.

## Process Index

| ID | Priority | Process | Primary Actors | Acceptance IDs | Status |
|---|---|---|---|---|---|
| BP-001 | P0 | Login and role entry | Administrator, Sales Manager, Sales | ACC-001, ACC-002 | Accepted as Architecture Input |
| BP-002 | P0 | Lead intake, assignment, and qualification | Sales, Sales Manager | ACC-003, ACC-004, ACC-014, ACC-016 | Accepted as Architecture Input |
| BP-003 | P0 | Customer and contact setup | Sales, Sales Manager | ACC-005, ACC-006, ACC-014, ACC-016 | Accepted as Architecture Input |
| BP-004 | P0 | Opportunity pipeline management | Sales, Sales Manager | ACC-007, ACC-008, ACC-013, ACC-014 | Accepted as Architecture Input |
| BP-005 | P0 | Quote to contract management | Sales, Sales Manager | ACC-009, ACC-010, ACC-014 | Accepted as Architecture Input |
| BP-006 | P0 | Payment tracking and won/lost closure | Sales, Sales Manager | ACC-011, ACC-013, ACC-014 | Accepted as Architecture Input |
| BP-007 | P0/P1 | Activities, notes, tasks, and reminders | Sales, Sales Manager | ACC-012, ACC-021 | Accepted as Architecture Input |
| BP-008 | P0/P1 | Team management and overview | Sales Manager, Administrator | ACC-002, ACC-018, ACC-023 | Accepted as Architecture Input |
| BP-009 | P1 | CSV import and export | Administrator, Sales Manager | ACC-020 | Accepted as Architecture Input |
| BP-010 | P0/P1 | History, operation logs, and reports | Administrator, Sales Manager, Sales | ACC-014, ACC-022, ACC-023 | Accepted as Architecture Input |
| BP-011 | P0/P1 | Archive and active-work filtering | Administrator, Sales Manager | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | Accepted as Architecture Input |

## BP-001: Login And Role Entry

Trigger:
- User signs in to the CRM.

Main flow:
1. User provides credentials.
2. CRM authenticates the user.
3. CRM loads the assigned role: Administrator, Sales Manager, or Sales.
4. CRM applies role and record-visibility rules to all subsequent actions.

Business outcome:
- The user can operate only within the role and record scope assigned to them.

Exception flow:
- Invalid credentials, disabled users, and unauthenticated requests are denied.
- Denied access must not expose CRM data.

Acceptance:
- ACC-001, ACC-002

## BP-002: Lead Intake, Assignment, And Qualification

Trigger:
- A lead is created, imported, assigned, or opened for qualification.

Main flow:
1. Sales, Sales Manager, or Administrator creates a lead with lead name or
   company name, source, and status.
2. Unassigned leads may exist before assignment.
3. Administrator or Sales Manager assigns owner before Pending Qualification or
   later states.
4. Sales or Sales Manager records qualification result.
5. Valid leads can be converted to customer/contact/opportunity context.
6. Invalid leads require invalid reason.
7. Lead changes create business history.

Business outcome:
- Lead origin, owner, status, qualification result, and conversion history are
  preserved.

Exception flow:
- Sales cannot qualify, edit, or convert Unassigned leads.
- Invalid leads cannot convert unless restored to Pending Qualification by
  Administrator or Sales Manager.
- Converted leads cannot be converted again.

Acceptance:
- ACC-003, ACC-004, ACC-014, ACC-016

## BP-003: Customer And Contact Setup

Trigger:
- A valid lead is converted or a user creates customer/contact records directly.

Main flow:
1. User creates or links a company/customer.
2. User records customer status and owner.
3. User creates one or more contacts under the company/customer.
4. Each contact has contact name, related company/customer, and at least one
   contact method or role note.
5. Related records become available in authorized list, detail, search, and
   filter views.

Business outcome:
- ToB account and contact structure is available for opportunity, quote,
  contract, payment, and follow-up work.

Exception flow:
- Missing company/customer link blocks contact save.
- Unauthorized access is denied.
- Hard delete is unavailable for core CRM records.

Acceptance:
- ACC-005, ACC-006, ACC-014, ACC-015, ACC-016

## BP-004: Opportunity Pipeline Management

Trigger:
- A valid business need is ready to track as an opportunity.

Main flow:
1. User creates an opportunity linked to customer, owner, stage, status,
   expected amount, and expected close date.
2. User advances stage through the sales pipeline.
3. Each allowed stage transition creates record-local history.
4. Sales closure as Won requires full payment.
5. Sales closure as Lost requires lost reason.

Business outcome:
- Opportunity progress, status, value, and closure rationale are traceable.

Exception flow:
- Forbidden transitions are rejected.
- Won and Lost are terminal in the committed scope.
- Reopen is unavailable in the committed scope.

Acceptance:
- ACC-007, ACC-008, ACC-013, ACC-014

## BP-005: Quote To Contract Management

Trigger:
- An opportunity reaches quote or contract negotiation work.

Main flow:
1. User creates one or more quotes linked to opportunity and customer.
2. User sends, accepts, rejects, or lets quotes expire according to quote rules.
3. Only one quote can be Accepted for an opportunity at a time.
4. User creates a Pending Signature contract from an Accepted quote.
5. Pending Signature contract requires expected signed date and contract note.
6. Signed/Active/Completed/post-signature Terminated contracts require
   signed/effective date.

Business outcome:
- Quote and contract decisions are traceable and linked to the opportunity.

Exception flow:
- Expired quote cannot be linked to a new contract.
- Contract amount differing from accepted quote amount requires difference
  reason.
- Contract approval, electronic signature, and template generation are out of
  committed P0/P1 scope.

Acceptance:
- ACC-009, ACC-010, ACC-014

## BP-006: Payment Tracking And Won/Lost Closure

Trigger:
- A contract has payment terms or actual payments to record.

Main flow:
1. User creates payment plan records linked to contract.
2. User records actual payments.
3. Partial payment updates status to Partially Paid.
4. Full payment updates status to Paid.
5. Overdue status applies when due date has passed and unpaid amount remains.
6. Opportunity can move to Won only after full payment is recorded.

Business outcome:
- Payment plan, actual payments, overdue state, and opportunity closure are
  auditable.

Exception flow:
- Zero, negative, or overpayment amounts are rejected.
- Single currency applies in the committed scope.
- Tax, discount, and multi-currency automation are out of committed P0/P1 scope.

Acceptance:
- ACC-011, ACC-013, ACC-014

## BP-007: Activities, Notes, Tasks, And Reminders

Trigger:
- User records follow-up, collaboration, or reminder-driven work.

Main flow:
1. User creates activities, notes, or tasks against lead, customer, contact,
   opportunity, quote, contract, payment, or related CRM record.
2. Open tasks have owner, due date, status, and title.
3. In-app reminders show authorized due/overdue tasks.
4. In-app reminders show Pending Signature contracts past expected signed date.
5. In-app reminders show due/overdue payments.

Business outcome:
- Follow-up work remains visible and traceable for authorized team members.

Exception flow:
- Completed/cancelled tasks are not active reminders.
- Signed contracts, terminated contracts, and fully paid contracts do not
  create active reminders.
- Unauthorized records are hidden.

Acceptance:
- ACC-012, ACC-021

## BP-008: Team Management And Overview

Trigger:
- Sales Manager or Administrator reviews team work or transfers ownership.

Main flow:
1. Sales Manager views all team records in the committed scope.
2. Sales Manager assigns or transfers team work.
3. Open tasks and follow-ups transfer with the parent owner unless manually
   reassigned during transfer.
4. Sales Manager reviews team overview for leads, opportunities, quotes,
   contracts, payments, tasks, and pipeline status.
5. Sales Manager may archive eligible team records.

Business outcome:
- Team coordination is possible without broad Sales-user visibility.

Exception flow:
- Sales users cannot access manager overview.
- Sales users cannot archive records.
- Administrator and Sales Manager act as themselves, not silently on behalf of
  Sales users.

Acceptance:
- ACC-002, ACC-018, ACC-023

## BP-009: CSV Import And Export

Trigger:
- Administrator or Sales Manager imports or exports CRM data.

Main flow:
1. Authorized user selects CSV import or export.
2. User imports or exports Lead, Company/customer, Contact, Opportunity, Quote,
   Contract, Payment plan, Actual payment, Activity, Note, or Task records.
3. Import validates required fields row by row.
4. Valid rows are imported.
5. Invalid rows are reported with row-level errors.
6. Export includes authorized records only.

Business outcome:
- Bulk data work supports committed operations without corrupting existing records or
  exposing unauthorized data.

Exception flow:
- Unsupported formats are rejected.
- Invalid rows do not corrupt existing records.
- Sales users cannot import/export in the committed scope.

Acceptance:
- ACC-020

## BP-010: History, Operation Logs, And Reports

Trigger:
- User opens record-local history, Administrator opens global operation logs,
  or Administrator/Sales Manager opens reports.

Main flow:
1. Authorized users view record-local business history from related records.
2. Administrator views global operation log events.
3. Administrator and Sales Manager view basic sales reports.
4. Reports show counts and sums for leads, opportunities, quotes, contracts,
   and payments using persisted authorized records.
5. Default reports exclude archived records unless an authorized user applies an
   explicit archived filter.

Business outcome:
- Collaboration, audit, and management reporting are supported by persisted
  traceable data.

Exception flow:
- Record-local history cannot be edited through normal CRM actions.
- Sales Manager and Sales cannot access global operation logs.
- Sales users cannot access manager/admin reports.

Acceptance:
- ACC-014, ACC-022, ACC-023

## BP-011: Archive And Active-Work Filtering

Trigger:
- Administrator or Sales Manager archives an eligible record.

Main flow:
1. Authorized user opens an eligible record.
2. User archives the record.
3. CRM records archive history and global operation log where applicable.
4. Archived record is removed from active work lists, active reminders, and
   default operational views.
5. Authorized users can retrieve archived records through explicit archived
   filters and audit/history views.

Business outcome:
- Records can leave active operations without hard deletion or history loss.

Exception flow:
- Sales cannot archive records.
- Active downstream obligations may need resolution before archive.

Acceptance:
- ACC-002, ACC-014, ACC-015, ACC-021, ACC-023
