# Screen Flows

## Document Control

- Project: CRM System
- Phase: G4 UX Design
- Owner Agent: UX Designer
- Source: `docs/ux-ui/user-journeys.md`, `docs/ux-ui/ux-flows.md`
- Status: Accepted as Architecture Input

## Screen Flow Index

| ID | Flow | Entry Screen | Exit Screen | Acceptance IDs | Status |
|---|---|---|---|---|---|
| SF-001 | Sign in to role workspace | Sign In | Role Workspace | ACC-001, ACC-002 | Accepted as Architecture Input |
| SF-002 | Lead create and qualification | Lead List / Work Overview | Opportunity Detail | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007 | Accepted as Architecture Input |
| SF-003 | Opportunity pipeline | Opportunity List | Opportunity Detail / Terminal State | ACC-007, ACC-008, ACC-013 | Accepted as Architecture Input |
| SF-004 | Quote and contract | Opportunity Detail | Contract Detail | ACC-009, ACC-010 | Accepted as Architecture Input |
| SF-005 | Payment and closure | Contract Detail | Opportunity Terminal State | ACC-011, ACC-013 | Accepted as Architecture Input |
| SF-006 | Tasks and reminders | Work Overview / Reminder Center | Related Record Detail | ACC-012, ACC-021 | Accepted as Architecture Input |
| SF-007 | Team overview | Manager Overview | Team Record Detail | ACC-018, ACC-023 | Accepted as Architecture Input |
| SF-008 | Import/export | Import/Export | Import Result / Export Result | ACC-020 | Accepted as Architecture Input |
| SF-009 | History and logs | Record Detail / Admin Logs | Related Record Detail | ACC-014, ACC-022 | Accepted as Architecture Input |
| SF-010 | Archive | Record Detail | Record Detail / Archived Filter | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | Accepted as Architecture Input |
| SF-011 | Entity list/detail/search/filter pattern | Entity List | Entity Detail | ACC-015 | Accepted as Architecture Input |
| SF-012 | Administrator user and role management | Admin User/Role Management | User Detail / Role Detail | ACC-001, ACC-002 | Accepted as Architecture Input |

## Common Screens

| Screen | Purpose | Primary Roles |
|---|---|---|
| Sign In | Authenticate user and start role-scoped session | All |
| Role Workspace | Entry dashboard for allowed work areas | All |
| Work Overview | Assigned active work and reminders | Sales |
| Lead List | Search, filter, create, open, and assign leads | Sales, Sales Manager |
| Lead Detail | View/edit lead and qualification status | Sales, Sales Manager |
| Company/Customer List | Search, filter, create, and open company/customer records | Sales, Sales Manager |
| Customer Detail | View/edit company/customer and related contacts | Sales, Sales Manager |
| Contact List | Search, filter, create, and open contact records | Sales, Sales Manager |
| Contact Detail | View/edit contact and related company/customer context | Sales, Sales Manager |
| Opportunity List | Search, filter, create, and open opportunity records | Sales, Sales Manager |
| Opportunity Detail | Pipeline, quote, contract, payment, tasks, history | Sales, Sales Manager |
| Quote List | Search, filter, create, and open quote records | Sales, Sales Manager |
| Quote Detail | Quote status, amount, validity, accepted state | Sales, Sales Manager |
| Contract List | Search, filter, create, and open contract records | Sales, Sales Manager |
| Contract Detail | Contract status, dates, notes, amount, payment links | Sales, Sales Manager |
| Payment List | Search, filter, create, and open payment records | Sales, Sales Manager |
| Payment Detail | Payment plan and actual payment status | Sales, Sales Manager |
| Activity/Note/Task List | Search, filter, create, and open activity, note, or task records or related-record sections | Sales, Sales Manager |
| Activity/Note/Task Detail | View/edit activity, note, or task context | Sales, Sales Manager |
| Reminder Center | Due/overdue tasks, contracts, and payments | Sales, Sales Manager |
| Manager Overview | Team pipeline and operational summary | Sales Manager |
| Admin User/Role Management | Manage users and assigned roles | Administrator |
| User Detail | View/edit user access state and role assignment | Administrator |
| Role Detail | View role capability summary | Administrator |
| Admin Logs | Global operation log query | Administrator |
| Reports | Basic report summaries | Administrator, Sales Manager |
| Import/Export | CSV import and export | Administrator, Sales Manager |

## SF-001: Sign In To Role Workspace

Screen path:
1. Sign In
2. Loading authentication
3. Role Workspace
4. Permission-filtered navigation

Error path:
- Invalid credentials -> Sign In with error.
- Disabled user -> Sign In with access denied.

Exit:
- Role Workspace.

## SF-002: Lead Create And Qualification

Screen path:
1. Work Overview or Lead List
2. Lead Create
3. Lead Detail
4. Qualification action
5. Customer/Contact Create or Link
6. Opportunity Create
7. Opportunity Detail

Error path:
- Missing required lead fields -> Lead Create with inline errors.
- Unassigned lead opened by Sales -> permission or disabled-action state.
- Invalid lead conversion -> blocked state with return path.

Exit:
- Opportunity Detail or Lead Detail.

## SF-003: Opportunity Pipeline

Screen path:
1. Opportunity List
2. Opportunity Detail
3. Stage action
4. Required data prompt if needed
5. Stage success feedback and history update

Error path:
- Forbidden transition -> blocked stage message.
- Won without a Signed contract -> blocked close with Signed-contract requirement message (DEC-017).
- Lost without reason -> close reason required.

Exit:
- Opportunity Detail with updated stage or terminal state.

## SF-004: Quote And Contract

Screen path:
1. Opportunity Detail
2. Quote Create
3. Quote Detail
4. Quote status action
5. Contract Create from Accepted Quote
6. Contract Detail

Error path:
- Expired quote selected -> blocked contract action.
- Each opportunity has exactly one quote (DEC-018).
- Missing expected signed date -> contract form error.
- Amount mismatch -> difference reason required.

Exit:
- Contract Detail.

## SF-005: Payment And Closure

Screen path:
1. Contract Detail
2. Payment Plan Create
3. Actual Payment Record
4. Payment status update
5. Opportunity Close action
6. Terminal state confirmation

Error path:
- Invalid payment amount -> payment form error.
- Overpayment -> blocked payment with remaining amount context.
- Won without a Signed contract -> blocked close (DEC-017).

Exit:
- Opportunity Detail in Won or Lost state.

## SF-006: Tasks And Reminders

Screen path:
1. Work Overview or Reminder Center
2. Reminder List
3. Related Record Detail
4. Task/contract/payment resolution action
5. Reminder list refresh

Error path:
- Related record no longer authorized -> permission denied and safe return.
- Completed/cancelled task -> no active reminder.

Exit:
- Related Record Detail or Reminder Center.

## SF-007: Team Overview

Screen path:
1. Manager Overview
2. Team list or summary panel
3. Team Record Detail
4. Assign/transfer action
5. Transfer confirmation

Error path:
- Sales opens Manager Overview -> permission denied.
- Empty team records -> empty state.

Exit:
- Manager Overview or Team Record Detail.

## SF-008: Import/Export

Screen path:
1. Import/Export
2. Object selection
3. CSV selection or export scope
4. Progress
5. Result summary
6. Row-level errors or export result

Error path:
- Unsupported format -> file error.
- Partial import failure -> row-level error table.
- Sales access -> permission denied.

Exit:
- Import Result or Export Result.

## SF-009: History And Logs

Screen path:
1. Record Detail
2. Record-local History
3. Related event detail or related record

Admin path:
1. Admin Logs
2. Filter/search
3. Event result
4. Related record detail where allowed

Error path:
- Non-admin global log access -> permission denied.
- Log edit attempt -> action unavailable.

Exit:
- Record Detail or Admin Logs.

## SF-010: Archive

Screen path:
1. Record Detail
2. Archive action
3. Confirmation
4. Success feedback
5. Active list without archived record

Blocked path:
1. Archive action
2. Blocked by active obligations
3. Related obligation list
4. User opens related item
5. User resolves item and retries archive

Exit:
- Record Detail, archived filter, or active list.

## SF-011: Entity List/Detail/Search/Filter Pattern

Covered entities:
- Leads.
- Companies/customers.
- Contacts.
- Opportunities.
- Quotes.
- Contracts.
- Payments.
- Activities.
- Notes.
- Tasks.

Screen path:
1. Entity List.
2. Search and filter controls.
3. Authorized result list.
4. Entity Detail.
5. Create or edit form where role permits.
6. Save confirmation and return to Entity Detail or Entity List.

Entity-specific entry points:
- Lead List opens Lead Detail.
- Company/Customer List opens Customer Detail.
- Contact List opens Contact Detail.
- Opportunity List opens Opportunity Detail.
- Quote List opens Quote Detail.
- Contract List opens Contract Detail.
- Payment List opens Payment Detail.
- Activity/Note/Task List opens Activity/Note/Task Detail or related-record
  section.

Error path:
- Empty result -> entity-specific empty state.
- Invalid filter -> filter error with preserved filter inputs.
- Unauthorized record -> record hidden or permission denied without restricted
  data.
- Save failure -> form error and preserved input.

Exit:
- Entity Detail, Entity List, or related parent record.

## SF-012: Administrator User And Role Management

Screen path:
1. Role Workspace.
2. Admin User/Role Management.
3. User List.
4. User Detail.
5. Edit user status or assigned role.
6. Save confirmation and operation-log event.
7. Role Detail for role capability summary when needed.

Error path:
- Sales Manager or Sales opens user/role management -> permission denied.
- Missing required user fields -> validation error.
- Disabled user cannot access CRM after change.
- Failed save keeps entered data and shows error.

Exit:
- Admin User/Role Management, User Detail, or Role Workspace.
