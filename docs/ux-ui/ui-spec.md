# UI Spec

## Document Control

- Project: CRM System
- Phase: G4 UI Design
- Owner Agent: UI Designer
- Source: `docs/ux-ui/user-journeys.md`, `docs/ux-ui/ux-flows.md`,
  `docs/ux-ui/screen-flows.md`, `docs/ux-ui/interaction-spec.md`,
  `docs/ux-ui/screen-state-spec.md`
- Status: Accepted as Architecture Input

## UI Direction

The CRM UI is an operational ToB sales workspace. It should be dense, calm,
predictable, and optimized for repeated daily use.

UI principles:
- Prioritize list/detail scanning, fast form completion, and visible workflow
  status.
- Use a restrained neutral base with clear semantic colors for status and
  validation.
- Avoid marketing-style hero layouts, decorative sections, and visual elements
  that do not support work.
- Use icon buttons for common tools when the icon is familiar and always
  provide accessible labels or tooltips.
- Use text labels for high-risk or business-critical actions such as Won,
  Lost, Archive, Import, and user status changes.
- Do not hide missing business behavior behind visual polish.
- P0/P1 UI must not weaken product acceptance, UX behavior, permission rules,
  or persistence expectations.

## Visual System

### Layout

| Token | Rule |
|---|---|
| Page shell | Persistent app shell with left navigation on desktop and compact top navigation on smaller screens. |
| Content width | Use full available width for operational tables and detail work; constrain forms only when readability requires it. |
| Grid | 12-column desktop grid; 8-column tablet grid; single-column mobile flow. |
| Spacing | 4px base scale. Common spacing: 8, 12, 16, 24. |
| Border radius | 6px default; 8px maximum for cards/modals unless implementation design system requires otherwise. |
| Density | Tables and list rows should support compact, standard, and comfortable density if later enabled; default is standard operational density. |

### Typography

| Usage | Rule |
|---|---|
| Page title | Clear entity or workspace name, not marketing copy. |
| Section heading | Short, task-oriented label. |
| Table text | Compact, readable, no negative letter spacing. |
| Form labels | Always visible. Placeholders must not replace labels. |
| Error text | Concise and placed near the failed field plus form summary when multiple errors exist. |

### Color And Status

| Usage | Rule |
|---|---|
| Base | Neutral background and text hierarchy for long operational sessions. |
| Primary action | One consistent primary action color. |
| Success | Used for saved, completed, paid, won states. |
| Warning | Used for duplicate warnings, pending signature risk, overdue soon, amount mismatch requiring reason. |
| Danger | Used for lost, termination, archive confirmation, denied destructive actions, invalid state. |
| Info | Used for guidance, neutral reminders, and history context. |
| Permission denied | Use neutral/danger combination without revealing restricted data. |

Colors must not be the only signal. Pair color with text, icon, or status label.

### Icons

Use icons for:
- Search.
- Filter.
- Sort.
- Refresh.
- Add/create.
- Edit.
- More actions.
- History.
- Reminder.
- Export/import.
- Archive.
- Warning/error/success states.

Every icon-only control requires an accessible name and tooltip.

### Data Display Safety

UI text must not expose sensitive or restricted raw values in generic feedback
states. This applies to permission-denied messages, toast messages, form error
summaries, import row errors, table-level errors, and log summaries.

Rules:
- Use safe summaries by default for errors and notifications. Do not echo full
  contact details, payment values, customer-sensitive text, or restricted record
  identifiers unless the current role is authorized to view that field.
- Permission-denied states must avoid revealing restricted record names,
  restricted field values, or whether a restricted record exists.
- History and operation-log event detail may show before/after values only when
  allowed by Security data classification and the user's role authorization.
- Import row errors prioritize row number, field name, and validation rule.
  They must not default to full contact, amount, payment, or customer value
  echoing.
- Export result UI must summarize completion and counts without unnecessary
  sensitive sample rows or raw examples.

## Navigation

### Desktop Navigation

Structure:
- Left navigation.
- Top utility area for user identity, role, global search if later supported,
  and account actions.
- Main content area.

Primary sections:
- Work Overview.
- Leads.
- Companies/Customers.
- Contacts.
- Opportunities.
- Quotes.
- Contracts.
- Payments.
- Activities/Notes/Tasks.
- Reminders.
- Reports.
- Import/Export.
- Admin: Users/Roles, Operation Logs.

Visibility:
- Sales sees only allowed sections.
- Sales Manager sees team sections, reports, import/export, and team overview.
- Administrator sees governance and global log sections.
- Unauthorized sections are hidden or unavailable according to security design;
  UI visibility does not replace authorization enforcement.

### Mobile And Tablet Navigation

Structure:
- Collapsible navigation.
- Prioritize Work Overview, Search/Filters, and current detail context.
- Avoid table overflow that hides critical actions.

Mobile must support inspection and light updates, but high-volume CRM work is
desktop-first unless later promoted.

## Screen Index

| ID | Screen | Primary Roles | Source Flow | Acceptance IDs | Status |
|---|---|---|---|---|---|
| UI-001 | Sign In | All | SF-001 | ACC-001, ACC-002 | Accepted as Architecture Input |
| UI-002 | Role Workspace / Work Overview | All / Sales | SF-001, SF-006 | ACC-002, ACC-012, ACC-015, ACC-021 | Accepted as Architecture Input |
| UI-003 | Entity List Pattern | Sales, Sales Manager, Administrator where allowed | SF-011 | ACC-015 | Accepted as Architecture Input |
| UI-004 | Entity Detail Pattern | Sales, Sales Manager, Administrator where allowed | SF-011 | ACC-003 to ACC-016 | Accepted as Architecture Input |
| UI-005 | Lead Detail And Qualification | Sales, Sales Manager | SF-002 | ACC-003, ACC-004 | Accepted as Architecture Input |
| UI-006 | Customer/Contact Detail | Sales, Sales Manager | SF-002, SF-011 | ACC-005, ACC-006 | Accepted as Architecture Input |
| UI-007 | Opportunity Detail | Sales, Sales Manager | SF-003 | ACC-007, ACC-008, ACC-013, ACC-014 | Accepted as Architecture Input |
| UI-008 | Quote Detail | Sales, Sales Manager | SF-004 | ACC-009 | Accepted as Architecture Input |
| UI-009 | Contract Detail | Sales, Sales Manager | SF-004 | ACC-010, ACC-021 | Accepted as Architecture Input |
| UI-010 | Payment Detail | Sales, Sales Manager | SF-005 | ACC-011, ACC-013, ACC-021 | Accepted as Architecture Input |
| UI-011 | Reminder Center | Sales, Sales Manager | SF-006 | ACC-012, ACC-021 | Accepted as Architecture Input |
| UI-012 | Manager Overview | Sales Manager | SF-007 | ACC-018, ACC-023 | Accepted as Architecture Input |
| UI-013 | Import/Export | Administrator, Sales Manager | SF-008 | ACC-020 | Accepted as Architecture Input |
| UI-014 | History And Admin Logs | Administrator, Sales Manager, Sales by scope | SF-009 | ACC-014, ACC-022 | Accepted as Architecture Input |
| UI-015 | Reports | Administrator, Sales Manager | SF-007, SF-009 | ACC-018, ACC-023 | Accepted as Architecture Input |
| UI-016 | Archive Confirmation | Administrator, Sales Manager | SF-010 | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | Accepted as Architecture Input |
| UI-017 | Admin User/Role Management | Administrator | SF-012 | ACC-001, ACC-002, ACC-022 | Accepted as Architecture Input |

## Screen Specifications

### UI-001: Sign In

Purpose:
- Authenticate users and route them into role-scoped CRM.

Layout:
- Centered sign-in panel on neutral page background.
- Product name and environment indicator if applicable.
- Credential fields, submit button, and error area.

Components:
- Text inputs.
- Password input.
- Primary button.
- Error alert.

States:
- Loading: submit button shows progress and fields remain stable.
- Error: one generic unauthenticated failure message for invalid credentials,
  disabled accounts, unavailable accounts, or other unavailable sign-in states.
- Disabled: submit disabled until required fields present.
- Success: route to role workspace.

Security:
- The unauthenticated sign-in screen must not distinguish invalid credentials
  from disabled or unavailable accounts.
- Disabled-account detail is visible only inside Admin User/Role Management to
  authorized Administrators.

Accessibility:
- Labels visible.
- Error alert announced.
- Keyboard focus starts on first field.

### UI-002: Role Workspace / Work Overview

Purpose:
- Give users immediate access to allowed active work and reminders.

Layout:
- App shell with navigation.
- Summary strip for assigned or team work, depending on role.
- Main area with active work list and reminder panel.
- Secondary area for recent history or quick actions where role permits.

Components:
- Navigation item.
- Summary metric.
- Entity list.
- Reminder list.
- Quick action menu.
- Status badge.

States:
- Empty assigned work.
- Loading active work.
- Permission-filtered sections.
- Error loading list.

### UI-003: Entity List Pattern

Purpose:
- Provide consistent list, search, filter, and open behavior for P0 entities.

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

Layout:
- Page header with entity name and allowed primary action.
- Toolbar with search, filters, sort, density, refresh, and optional archived
  filter.
- Data table on desktop.
- Compact row list on mobile.
- Empty/error/permission states below toolbar.

Components:
- Toolbar.
- Search input.
- Filter menu or filter panel.
- Sort control.
- Data table.
- Pagination or incremental loading control.
- Row action menu.
- Status badge.

States:
- Loading table rows.
- Empty authorized results.
- Invalid filter.
- Permission-filtered rows.
- Archived filter active.

### UI-004: Entity Detail Pattern

Purpose:
- Provide consistent detail structure for core CRM records.

Layout:
- Header with entity title, status/stage badge, owner, primary actions, and
  more actions.
- Main detail area with record fields.
- Related tabs or sections: activity, notes, tasks, history, linked objects.
- Right-side or lower summary area for reminders, key dates, and risk signals.

Components:
- Detail header.
- Status/stage badge.
- Form section.
- Related-record table.
- History timeline.
- Action menu.
- Confirmation modal.

States:
- Loading detail.
- Save success.
- Save error.
- Permission denied.
- Disabled actions for missing business prerequisites.
- Conflict state for stale edits.

### UI-005: Lead Detail And Qualification

Purpose:
- Capture lead information and support qualification/conversion.

Layout:
- Lead header with owner and status.
- Required fields section.
- Qualification action area.
- Conversion section for customer/contact/opportunity.
- History section.

Components:
- Required text fields.
- Source selector.
- Owner display.
- Qualification segmented control or action buttons.
- Invalid reason field.
- Convert action.

States:
- Unassigned lead state.
- Sales denied/disabled qualification on Unassigned lead.
- Invalid reason required.
- Converted read-only conversion action.

### UI-006: Customer/Contact Detail

Purpose:
- Manage ToB account and contacts.

Layout:
- Customer header with status and owner.
- Customer fields.
- Contacts section as related table/list.
- Related opportunities, contracts, payments, and history.

Components:
- Customer status badge.
- Contact table.
- Add contact action.
- Related-record tabs.

States:
- No contacts.
- Missing contact method or role note.
- Unauthorized customer/contact.

### UI-007: Opportunity Detail

Purpose:
- Manage sales pipeline and closure.

Layout:
- Opportunity header with stage, status, expected amount, expected close date.
- Stage path or stepper.
- Required next data panel.
- Related quote, contract, payment, task, and history sections.

Components:
- Stage stepper.
- Stage action button.
- Blocked transition alert.
- Close Won/Lost confirmation.
- Related object cards or tables.

States:
- Forbidden transition.
- Missing required data.
- Won blocked until full payment.
- Lost reason required.
- Terminal Won/Lost.

### UI-008: Quote Detail

Purpose:
- Manage quote status and contract linkage.

Layout:
- Quote header with status, amount, validity end date.
- Quote form fields.
- Accepted quote indicator.
- Related opportunity and contract link.

Components:
- Status badge.
- Validity date field.
- Accept/reject actions.
- Expired warning.

States:
- Draft, Sent, Accepted, Rejected, Expired.
- Accepted quote conflict prevention.
- Expired quote contract-link blocked.

### UI-009: Contract Detail

Purpose:
- Manage record-based contract lifecycle.

Layout:
- Contract header with status, amount, related opportunity and quote.
- Date section for expected signed date and signed/effective date.
- Required contract note.
- Payment and reminder sections.
- History section.

Components:
- Contract status badge.
- Date inputs.
- Amount field.
- Difference reason field.
- Contract note area.
- Status action controls.

States:
- Pending Signature requires expected signed date.
- Signed/Active/Completed/post-signature Terminated require signed/effective
  date.
- Amount mismatch requires reason.
- Pending Signature reminder warning when past expected signed date.

### UI-010: Payment Detail

Purpose:
- Track payment plan and actual payment.

Layout:
- Payment header with status.
- Contract link and amount context.
- Payment plan fields.
- Actual payment entry area.
- Remaining amount and overdue signal.

Components:
- Amount input.
- Date input.
- Payment status badge.
- Remaining amount display.
- Validation alert.

States:
- Unpaid, Partially Paid, Paid, Overdue.
- Zero/negative amount blocked.
- Overpayment blocked.
- Full payment unlocks Won path.

### UI-011: Reminder Center

Purpose:
- Present actionable due/overdue work.

Layout:
- Reminder list grouped by type: task, contract, payment.
- Filters for owner, due state, type.
- Row links to related record.

Components:
- Reminder row.
- Type icon.
- Due status badge.
- Open related record action.

States:
- No active reminders.
- Unauthorized reminders hidden.
- Inactive reminder removed after resolution.

### UI-012: Manager Overview

Purpose:
- Let Sales Manager inspect team pipeline and risk.

Layout:
- Team summary strip.
- Pipeline table or compact stage summary.
- Reminders and overdue work panel.
- Team record list.

Components:
- Summary metric.
- Team record table.
- Assignment action.
- Transfer confirmation.

States:
- Empty team data.
- Transfer success.
- Sales denied.

### UI-013: Import/Export

Purpose:
- Support CSV import/export for authorized roles.

Layout:
- Object selector.
- Import file area.
- Export scope area.
- Progress and result summary.
- Row-level error table.

Components:
- File input.
- Object selector.
- Progress indicator.
- Result summary.
- Error table.

States:
- Unsupported format.
- Long-running progress.
- Partial success.
- Export denied for Sales.
- Export confirmation required before export runs.

Export confirmation:
- Shows selected object.
- Shows filter conditions and explicit archived inclusion or exclusion.
- Shows authorization scope for the current role.
- Shows estimated record count before export when available.
- Shows audit notice that the export action will be logged.
- Does not display unnecessary sensitive sample rows or raw example data in the
  confirmation, progress, or result states.

### UI-014: History And Admin Logs

Purpose:
- Show record-local history and administrator operation logs.

Layout:
- Record-local history appears in record detail as timeline.
- Admin logs appear as queryable table.

Components:
- Timeline.
- Log table.
- Filter toolbar.
- Event detail panel.

States:
- Empty history/log.
- Log query error.
- Non-admin global log denied.
- Edit unavailable.

Security display rules:
- Timeline and log summaries use safe, role-authorized labels by default.
- Generic errors, toasts, and denied states do not echo restricted or sensitive
  raw values.
- Event detail panels may show before/after values only according to Security
  data classification and the user's role authorization.

### UI-015: Reports

Purpose:
- Show basic persisted sales metrics.

Layout:
- Report summary grid.
- Grouped tables for leads, opportunities, quotes, contracts, payments.
- Filter toolbar.

Components:
- Metric tile.
- Grouped table.
- Drill-in link.
- Empty state.

States:
- Empty/zero data.
- Unauthorized records excluded.
- Sales denied.
- Archived filter explicit when available.

### UI-016: Archive Confirmation

Purpose:
- Confirm archive and handle active downstream obligations.

Layout:
- Confirmation modal or panel.
- Effect summary.
- Active obligation list when blocked.

Components:
- Confirmation modal.
- Warning alert.
- Related obligation list.
- Retry action.

States:
- Confirmation.
- Blocked by active obligations.
- Archive success.
- Sales denied.

### UI-017: Admin User/Role Management

Purpose:
- Let Administrator manage users and assigned roles.

Layout:
- User list with search/filter.
- User detail panel.
- Role capability summary.
- Operation-log link or history hint.

Components:
- User table.
- User status badge.
- Role selector.
- Role summary table.
- Role/status change confirmation.

States:
- Loading users/roles.
- No users configured.
- Validation error.
- Non-admin denied.
- User status/role saved.
- Role/status change confirmation before save.
- Last Administrator downgrade/disable blocked.

Role/status change confirmation:
- Required before saving user status changes, role changes, Administrator
  grants, Administrator removals, or account disablement.
- Shows target user.
- Shows old role and new role.
- Shows old status and new status.
- Shows access impact summary.
- Shows audit/log notice.
- Requires explicit confirm and cancel actions.

Blocked governance state:
- Disabling, deactivating, or downgrading the last Administrator is shown as a
  blocked state with a safe explanation and no confirm action.
- UI confirmation and blocked-state presentation do not replace backend
  authorization or server-side governance enforcement.

## UI Verification Notes

- UI implementation must be checked against all screen states in
  `screen-state-spec.md`.
- UI must not treat hidden navigation as authorization enforcement.
- UI must preserve UX recovery paths for validation, permission denial, blocked
  transitions, partial CSV failure, archive blockers, and stale edit conflicts.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8.
