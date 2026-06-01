# UX Flows

## Document Control

- Project: CRM System
- Phase: G4 UX Design
- Owner Agent: UX Designer
- Source: `docs/product/acceptance-matrix.md`, `docs/business/*`
- Status: Accepted as Architecture Input

## Flow Index

| ID | Priority | User Goal | Primary Role | Acceptance IDs | Status |
|---|---|---|---|---|---|
| UX-001 | P0 | Sign in and enter role-scoped CRM | All roles | ACC-001, ACC-002 | Accepted as Architecture Input |
| UX-002 | P0 | Create and qualify lead | Sales | ACC-003, ACC-004 | Accepted as Architecture Input |
| UX-003 | P0 | Create customer/contact from valid lead | Sales | ACC-005, ACC-006 | Accepted as Architecture Input |
| UX-004 | P0 | Manage opportunity pipeline | Sales | ACC-007, ACC-008, ACC-013 | Accepted as Architecture Input |
| UX-005 | P0 | Create quote and contract | Sales | ACC-009, ACC-010 | Accepted as Architecture Input |
| UX-006 | P0 | Record payment and close deal | Sales | ACC-011, ACC-013 | Accepted as Architecture Input |
| UX-007 | P0/P1 | Create follow-up and handle reminders | Sales, Sales Manager | ACC-012, ACC-021 | Accepted as Architecture Input |
| UX-008 | P1 | Review and manage team work | Sales Manager | ACC-018, ACC-023 | Accepted as Architecture Input |
| UX-009 | P1 | Import/export CSV | Administrator, Sales Manager | ACC-020 | Accepted as Architecture Input |
| UX-010 | P0/P1 | Review history, logs, and reports | Administrator, Sales Manager, Sales | ACC-014, ACC-022, ACC-023 | Accepted as Architecture Input |
| UX-011 | P0/P1 | Archive eligible records | Administrator, Sales Manager | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | Accepted as Architecture Input |

## Common Flow Requirements

- Every create/edit flow uses field-level validation and form-level summary for
  blocking errors.
- Every permission denial states that access is unavailable and provides a safe
  return path.
- Every successful mutation gives visible confirmation and keeps the user in a
  useful next context.
- Every destructive or terminal business action requires confirmation.
- Every list/detail flow supports loading, empty, error, permission denied, and
  success states where applicable.

## UX-001: Sign In And Enter Role-Scoped CRM

Entry point:
- Sign-in screen.

Main flow:
1. User enters credentials.
2. UX shows loading state.
3. On success, user enters role-scoped CRM area.
4. Navigation exposes only allowed sections.

Success feedback:
- User sees role-appropriate workspace and active work.

Failure and recovery:
- Invalid credentials show error and keep user on sign-in.
- Disabled or unauthenticated access returns to sign-in.

Required screen states:
- Loading, error, success, permission denied.

## UX-002: Create And Qualify Lead

Entry point:
- Lead list, assigned work, or quick create.

Main flow:
1. User opens lead create.
2. User enters lead name or company name, source, and status.
3. User saves.
4. User qualifies as Valid or Invalid.
5. UX routes Valid path to customer/contact/opportunity creation.

Success feedback:
- Save and qualification confirmation.

Failure and recovery:
- Missing fields display inline errors.
- Unassigned lead qualification by Sales is disabled or denied.
- Invalid lead conversion is blocked with reason and allowed recovery route.

Required screen states:
- Loading, empty lead list, validation error, permission denied, success.

## UX-003: Create Customer/Contact From Valid Lead

Entry point:
- Valid lead detail or conversion flow.

Main flow:
1. User creates or links company/customer.
2. User creates contact under company/customer.
3. UX confirms created/linked records.
4. User proceeds to opportunity creation.

Failure and recovery:
- Contact without company/customer link is blocked.
- Missing contact method or role note is shown inline.

Required screen states:
- Search empty, duplicate warning, validation error, success.

## UX-004: Manage Opportunity Pipeline

Entry point:
- Opportunity list, customer detail, or converted lead.

Main flow:
1. User opens opportunity detail.
2. UX shows stage, required next data, related quote/contract/payment status,
   and history.
3. User changes stage when allowed.
4. UX creates visible history feedback.

Failure and recovery:
- Forbidden transition shows blocked reason.
- Missing required data links user to the relevant form section.
- Won is reached when the related contract is Signed; closing Won without a
  Signed contract is blocked (DEC-017).
- Lost without reason is blocked.

Required screen states:
- Loading, empty opportunities, validation error, blocked transition, success,
  terminal Won/Lost.

## UX-005: Create Quote And Contract

Entry point:
- Opportunity detail.

Main flow:
1. User creates quote.
2. User sends quote.
3. User marks quote Accepted.
4. User creates Pending Signature contract.
5. User enters expected signed date and contract note.
6. User later enters signed/effective date for signed lifecycle states.

Failure and recovery:
- Expired quote cannot start contract.
- Each opportunity has exactly one quote (DEC-018).
- Contract amount mismatch requires reason.
- Missing expected signed date blocks Pending Signature save.

Required screen states:
- Quote empty state, quote status feedback, contract validation errors,
  accepted quote indicator, success.

## UX-006: Record Payment And Close Deal

Entry point:
- Contract detail or opportunity payment section.

Main flow:
1. User creates payment plan.
2. User records actual payment.
3. UX updates payment status.
4. User closes opportunity as Won when the related contract is Signed, or Lost
   with lost reason (DEC-017).

Failure and recovery:
- Zero, negative, and overpayment values are blocked.
- Terminal close requires confirmation.
- Closed opportunity disables stage edits and allows notes/tasks (DEC-020).

Required screen states:
- Payment empty state, validation errors, overdue status, terminal close
  confirmation, success.

## UX-007: Create Follow-Up And Handle Reminders

Entry point:
- Record detail, active work, or reminder area.

Main flow:
1. User creates activity, note, or task.
2. UX shows task due state.
3. Reminder area shows due/overdue tasks, pending-signature contracts past
   expected signed date, and due/overdue payments.
4. User opens related record and resolves work.

Failure and recovery:
- Unauthorized reminders are hidden.
- Completed/cancelled tasks disappear from active reminders.

Required screen states:
- Empty reminders, overdue indicators, permission-filtered records, success.

## UX-008: Review And Manage Team Work

Entry point:
- Sales Manager team overview.

Main flow:
1. Sales Manager opens team overview.
2. UX shows pipeline, active tasks, reminders, and basic report summaries.
3. Manager opens detail, assigns/transfers owner, or archives eligible record.

Failure and recovery:
- Sales denied from manager overview.
- Archive blocked by active downstream obligations with links to related items.

Required screen states:
- Empty team data, loading, permission denied, blocked archive, transfer
  success.

## UX-009: Import/Export CSV

Entry point:
- Import/export section for Administrator or Sales Manager.

Main flow:
1. User selects object type and CSV.
2. UX validates file type.
3. UX shows import progress.
4. UX shows success count and row-level errors.
5. User exports authorized records.

Failure and recovery:
- Unsupported format rejected.
- Partial failures are summarized with row details.
- Sales denied from import/export.

Required screen states:
- File validation error, long-running progress, partial success, failure,
  export success.

## UX-010: Review History, Logs, And Reports

Entry point:
- Record detail, admin operation log, or reports.

Main flow:
1. User opens record-local history from a record.
2. Administrator opens global operation logs.
3. Administrator or Sales Manager opens reports.
4. UX shows filters, empty results, and related record navigation.

Failure and recovery:
- Log edit is unavailable.
- Non-admin global log access is denied.
- Unauthorized report records are excluded.

Required screen states:
- Empty history/log/report, permission denied, filter error, success.

## UX-011: Archive Eligible Records

Entry point:
- Authorized record detail action.

Main flow:
1. Administrator or Sales Manager selects archive.
2. UX shows confirmation and effect summary.
3. If active downstream obligations exist, UX blocks archive and lists related
   tasks/contracts/payments.
4. User resolves or archives related active obligations.
5. User retries archive.

Failure and recovery:
- Sales archive action unavailable or denied.
- Archived record leaves active lists and active reminders.

Required screen states:
- Confirmation, blocked archive, related obligation list, retry, success.
