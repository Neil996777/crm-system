# Business Rules

## Document Control

- Project: CRM System
- Phase: G4 Business Design
- Owner Agent: Business Analyst
- Source: `docs/product/prd.md`, `docs/product/acceptance-matrix.md`
- Status: Accepted as Architecture Input

## Rule Principles

- Rules clarify committed PRD and acceptance behavior.
- P0/P1 rules cannot be downgraded, deleted, weakened, merged away, or accepted
  as partial work.
- Any unresolved P0/P1-affecting business rule must become a blocker.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8.

## Rule Index

| ID | Priority | Rule | Acceptance IDs | Status |
|---|---|---|---|---|
| BR-001 | P0 | Role and visibility model | ACC-001, ACC-002 | Accepted as Architecture Input |
| BR-002 | P0 | No hard delete for core CRM records | ACC-002, ACC-005, ACC-016 | Accepted as Architecture Input |
| BR-003 | P0 | Required fields for valid save | ACC-003, ACC-005, ACC-006, ACC-007, ACC-009, ACC-010, ACC-011, ACC-012 | Accepted as Architecture Input |
| BR-004 | P0 | Lead ownership and qualification | ACC-003, ACC-004 | Accepted as Architecture Input |
| BR-005 | P0 | Opportunity lifecycle and closure | ACC-007, ACC-008, ACC-013 | Accepted as Architecture Input |
| BR-006 | P0 | Quote lifecycle and accepted quote constraint | ACC-009 | Accepted as Architecture Input |
| BR-007 | P0 | Contract lifecycle and expected signed date | ACC-010, ACC-021 | Accepted as Architecture Input |
| BR-008 | P0 | Payment amount and status rules | ACC-011, ACC-013, ACC-021 | Accepted as Architecture Input |
| BR-009 | P0 | Activity, note, and task rules | ACC-012, ACC-021 | Accepted as Architecture Input |
| BR-010 | P0/P1 | History and operation log rules | ACC-014, ACC-022 | Accepted as Architecture Input |
| BR-011 | P1 | Duplicate warning rules | ACC-019 | Accepted as Architecture Input |
| BR-012 | P1 | CSV import/export rules | ACC-020 | Accepted as Architecture Input |
| BR-013 | P1 | Reminder rules | ACC-021 | Accepted as Architecture Input |
| BR-014 | P1 | Basic report rules | ACC-018, ACC-023 | Accepted as Architecture Input |
| BR-015 | P0 | Persistence rule for core CRM paths | ACC-016, ACC-017 | Accepted as Architecture Input |
| BR-016 | P0/P1 | Archive eligibility and archived-record behavior | ACC-002, ACC-014, ACC-015, ACC-018, ACC-021, ACC-023 | Accepted as Architecture Input |
| BR-017 | P1 | Basic report metric definitions | ACC-018, ACC-023 | Accepted as Architecture Input |
| BR-018 | P1 | CSV import/export object scope | ACC-020 | Accepted as Architecture Input |
| BR-019 | P1 | Duplicate exact-match normalization | ACC-019 | Accepted as Architecture Input |
| BR-020 | P0 | Post-close editability | ACC-008, ACC-013, ACC-014 | Accepted as Architecture Input |
| BR-021 | P1 | Reminder date basis | ACC-021 | Accepted as Architecture Input |

## Rule Details

### BR-001: Role And Visibility Model

Actors:
- Administrator, Sales Manager, Sales

Rule:
- Administrator can govern all CRM records and operation logs.
- Sales Manager can view and manage all team records.
- Sales can view and manage only owned/assigned records and related child
  records.
- Unauthenticated users cannot access core CRM data.

Exception behavior:
- Unauthorized access is denied without exposing or mutating data.

Verification:
- Permission scenarios must cover create, view, edit, assign/transfer,
  close/archive, import/export, report, and audit actions.

### BR-002: No Hard Delete For Core CRM Records

Actors:
- Administrator, Sales Manager, Sales

Rule:
- The system does not allow hard deletion of core CRM records.
- Administrator and Sales Manager may archive eligible records.
- Sales cannot archive records.

Exception behavior:
- Delete attempts are rejected or unavailable.

Verification:
- Core CRM records remain recoverable for history, reports, and audit-sensitive
  review.

### BR-003: Required Fields For Valid Save

Rule:
- Lead requires lead name or company name, source, and status; owner is
  required before Pending Qualification or later states.
- Company/customer requires company name, customer status, and owner.
- Contact requires contact name, related company/customer, and at least one
  contact method or role note.
- Opportunity requires related company/customer, owner, stage, expected
  amount, and expected close date.
- Quote requires related opportunity, related company/customer, quote amount,
  status, validity end date, and owner.
- Contract requires related customer, related opportunity, accepted quote,
  contract amount, status, and contract note.
- Pending Signature contract requires expected signed date.
- Signed, Active, Completed, and post-signature Terminated contract requires
  signed/effective date.
- Payment plan requires related contract, due amount, due date, and status.
- Actual payment requires related contract, paid amount, payment date, and
  status.
- Activity/note requires related CRM record, type, content, actor, and
  timestamp.
- Task requires related CRM record, owner, due date, status, and title.

Exception behavior:
- Missing required fields block save.

### BR-004: Lead Ownership And Qualification

Rule:
- Unassigned leads may exist only before assignment.
- Unassigned leads cannot be qualified, edited by Sales, or converted.
- Invalid leads cannot convert unless first restored to Pending Qualification
  by Administrator or Sales Manager.
- Converted leads cannot be converted again.

History:
- Owner assignment, qualification, disqualification, and conversion create
  record-local history events.

### BR-005: Opportunity Lifecycle And Closure

Rule:
- Opportunity transitions must follow the PRD transition table.
- Won occurs when the related contract is Signed (DEC-017); full payment is not a
  Won precondition.
- Lost requires lost reason.
- Won and Lost are terminal in the committed scope.
- Reopen is not allowed in the committed scope (post-signing breach is handled at
  the contract level via Terminated; the opportunity stays Won).

Exception behavior:
- Forbidden stage transitions are rejected without data mutation.

Amended 2026-06-01 (DEC-017): Won is reached at contract signing, not full payment;
the `Payment In Progress` stage is removed from the pipeline.

### BR-006: Quote Lifecycle And Accepted Quote Constraint

Rule:
- Each opportunity has exactly one quote (DEC-018); the system records the quote
  result, not multi-round negotiation.
- The quote moves through Draft, Sent, Accepted, Rejected, Expired.
- A contract can reference only the opportunity's Accepted quote.
- An Expired quote cannot be linked to a new contract.

Exception behavior:
- Forbidden quote transitions are rejected without data mutation.

Amended 2026-06-01 (DEC-018): one quote per opportunity; the previous
multiple-quotes / one-Accepted-at-a-time constraint is removed.

### BR-007: Contract Lifecycle And Expected Signed Date

Rule:
- Pending Signature contracts require expected signed date and contract note.
- Pending Signature contracts do not require signed/effective date.
- Expected signed date is the planned signature deadline used by contract
  reminders.
- Signed, Active, Completed, and post-signature Terminated contracts require
  signed/effective date.
- Contract amount may differ from accepted quote amount only with a difference
  reason.

Out of committed P0/P1 scope:
- Contract approval, electronic signature, and template generation.

### BR-008: Payment Amount And Status Rules

Rule:
- The system uses one currency for quote, contract, and payment amounts.
- Payment amount must be greater than zero.
- Overpayment is blocked.
- Partial payment is supported.
- Overdue means due date passed and unpaid amount remains.
- Payment tracking is post-sale collection follow-up and is decoupled from Won
  (DEC-019); it is not a closing precondition.

Out of committed P0/P1 scope:
- Multi-currency, tax, and discount automation.

Amended 2026-06-01 (DEC-019): removed "full payment required before Won"; payment
plans, actual payments, status, overdue reminders, and reports are retained as
collection/visibility, decoupled from the Won decision.

### BR-009: Activity, Note, And Task Rules

Rule:
- Activity, note, and task must link to a related CRM record.
- Completed and Cancelled tasks are not active reminders.
- On owner transfer of a parent record, open tasks and follow-ups transfer to
  the new owner unless manually reassigned during transfer.

Exception behavior:
- Unauthorized users cannot view or create restricted related records.

### BR-010: History And Operation Log Rules

Rule:
- Record-local business history is visible according to record permissions.
- Admin/global operation logs are Administrator-only.
- History and operation logs cannot be modified through normal CRM actions.
- Administrator and Sales Manager act as themselves and do not silently act on
  behalf of Sales users.

Required events:
- Owner changes, stage/status changes, quote acceptance, contract signature or
  termination, payment records, overdue payment, opportunity won/lost, archive,
  import, export, and login/access failures.

### BR-011: Duplicate Warning Rules

Rule:
- Warn on exact company name match.
- Warn on contact phone or email match.
- Warn on lead company/contact match.
- Warning does not block save.
- No automatic merge or overwrite occurs in the committed scope.

### BR-012: CSV Import/Export Rules

Rule:
- CSV is the only required import/export format for the committed release.
- Import validates required fields row by row.
- Invalid rows are reported with row-level errors.
- Valid existing records must not be corrupted by failed rows.
- Export includes authorized records only.

### BR-013: Reminder Rules

Rule:
- In-app reminders cover due/overdue tasks.
- In-app reminders cover Pending Signature contracts past expected signed date.
- In-app reminders cover due/overdue payments.
- Completed/cancelled tasks, signed contracts, terminated contracts, and fully
  paid contracts do not create active reminders.
- Unauthorized reminder records are hidden.

### BR-014: Basic Report Rules

Rule:
- Sales Manager can view team overview and team reports.
- Administrator can view reports across governed CRM records.
- Sales users cannot access manager/admin reports.
- Reports use persisted authorized records only.
- Required report groupings are leads by status, opportunities by stage,
  quotes by status/amount, contracts by status/amount, and payments by
  status/amount.

### BR-015: Persistence Rule For Core CRM Paths

Rule:
- Core CRM data must persist across refresh, logout/login, and service restart.
- No P0 core CRM path may rely on mock, static-only, TODO, in-memory-only, or
  non-persistent behavior.

Verification:
- Persistence must be proven before any related P0/P1 item can be considered
  complete.

### BR-016: Archive Eligibility And Archived-Record Behavior

Rule:
- Eligible archive targets in the committed scope are Lead, Company/customer, Contact,
  Opportunity, Quote, Contract, Payment plan, Actual payment, Activity, Note,
  and Task.
- Administrator can archive eligible records.
- Sales Manager can archive eligible team records.
- Sales cannot archive records.
- Won and Lost opportunities may be archived but remain terminal.
- Archived records are excluded from active work lists, active reminders, and
  default operational views.
- Archived records remain available to authorized users through explicit
  archived filters, record-local history, operation logs, and audit/report
  evidence where applicable.
- Archived records are not hard deleted and must not break related business
  history.

Exception behavior:
- Archiving a record with active downstream obligations may require the user to
  resolve or archive related active tasks, pending-signature contracts, or
  unpaid payment items first.

### BR-017: Basic Report Metric Definitions

Rule:
- Leads by status count Lead records grouped by lead status.
- Opportunities by stage count Opportunity records grouped by current stage.
- Opportunity amount uses expected amount.
- Quotes by status/amount group Quote records by status and sum quote amount.
- Contracts by status/amount group Contract records by status and sum contract
  amount.
- Payments by status/amount group Payment plan or Actual payment records by
  status and sum due amount or paid amount respectively.
- Sales Manager reports use team records.
- Administrator reports use governed records.
- Sales users cannot access manager/admin reports.
- Archived records are excluded from default active reports and may appear only
  when an explicit archived filter is applied by an authorized user.

Exception behavior:
- Empty data returns zero or empty report state.

### BR-018: CSV Import/Export Object Scope

Rule:
- The committed CSV import/export minimum objects are Lead, Company/customer, Contact,
  Opportunity, Quote, Contract, Payment plan, Actual payment, Activity, Note,
  and Task.
- Administrator may import/export governed records.
- Sales Manager may import/export team records.
- Sales cannot import/export in the committed scope.
- Import validates each row against required fields, permissions, and related
  record references.
- Invalid rows are reported with row-level errors and do not corrupt existing
  records.
- Export includes authorized records only and excludes archived records by
  default unless an explicit archived filter is selected by an authorized user.

Exception behavior:
- Unsupported formats are rejected.

### BR-019: Duplicate Exact-Match Normalization

Rule:
- Company exact-name duplicate matching trims leading/trailing spaces and
  compares case-insensitively.
- Contact phone duplicate matching compares normalized digits after removing
  spaces, hyphens, and parentheses.
- Contact email duplicate matching trims spaces and compares case-insensitively.
- Lead company/contact duplicate matching uses the same company, phone, and
  email normalization rules.
- Duplicate warnings do not block save, do not merge records, and do not
  overwrite existing records.

Exception behavior:
- If data is too incomplete to evaluate a duplicate rule, no duplicate warning
  is required for that rule.

### BR-020: Post-Close Editability

Rule:
- Won and Lost opportunities are terminal in the committed scope.
- After Won/Lost closure, opportunity stage and close data are not
  editable through normal Sales workflow.
- Post-close notes and tasks may still be added through normal permissions.
- Related quote, contract, and payment records remain governed by their own
  lifecycle rules and permissions.
- Any allowed post-close activity creates normal history events.

Exception behavior:
- Reopen, stage rollback, early Won, or Lost without reason are rejected.

### BR-021: Reminder Date Basis

Rule:
- Reminder due/overdue evaluation uses the CRM workspace business date.
- Until Architecture defines timezone handling, the business date is the
  deployment-configured local date.
- A task or payment is due on its due date and overdue after the due date has
  passed while the item remains incomplete or unpaid.
- A Pending Signature contract is overdue for reminder purposes after expected
  signed date has passed while the contract remains Pending Signature.
- Reminder visibility follows record permission rules.

Exception behavior:
- Completed/cancelled tasks, signed contracts, terminated contracts, and fully
  paid contracts do not create active reminders.
