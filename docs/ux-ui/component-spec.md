# Component Spec

## Document Control

- Project: CRM System
- Phase: G4 UI Design
- Owner Agent: UI Designer
- Source: `docs/ux-ui/ui-spec.md`, `docs/ux-ui/interaction-spec.md`,
  `docs/ux-ui/screen-state-spec.md`
- Status: Accepted as Architecture Input

## Component Principles

- Components must support dense operational CRM workflows.
- Each interactive component must expose loading, disabled, focus, hover,
  error, and success states where applicable.
- Icon-only controls require accessible names and tooltips.
- Components must not imply permissions that backend/security enforcement does
  not grant.
- Visual variants must map to UX states and acceptance behavior.
- Generic feedback components must use safe summaries and avoid echoing
  restricted or sensitive raw values unless the role is authorized to view those
  values.

## Component Index

| ID | Component | Purpose | Variants | Required States | Acceptance IDs |
|---|---|---|---|---|---|
| CMP-001 | App Shell | Global navigation and page frame | Desktop sidebar, tablet collapsed, mobile top/drawer | Loading role, permission-filtered nav, active section | ACC-001, ACC-002 |
| CMP-002 | Page Header | Entity or workspace title and primary actions | Workspace, list, detail, admin | Loading, action disabled, permission filtered | ACC-015 |
| CMP-003 | Toolbar | Search, filter, sort, refresh, density, create | Entity list, report, log, import/export | Loading, invalid filter, no selection | ACC-015 |
| CMP-004 | Data Table | Dense operational lists | Entity, report, log, import error | Loading, empty, error, selected, sorted, permission filtered | ACC-015, ACC-020, ACC-022, ACC-023 |
| CMP-005 | Detail Header | Record identity, status, owner, actions | Lead, opportunity, quote, contract, payment | Loading, archived, terminal, permission denied | ACC-003 to ACC-014 |
| CMP-006 | Form Section | Structured create/edit inputs | Basic, grouped, required, readonly | Error, disabled, success, conflict | ACC-003 to ACC-012 |
| CMP-007 | Status Badge | Entity lifecycle or business state | Neutral, info, success, warning, danger | Active, overdue, terminal, archived | ACC-004, ACC-008, ACC-010, ACC-011, ACC-013 |
| CMP-008 | Stage Stepper | Opportunity stage visibility and actions | Horizontal desktop, compact mobile | Current, complete, blocked, terminal | ACC-008, ACC-013 |
| CMP-009 | Confirmation Modal | Confirm terminal or high-impact actions | Won, Lost, archive, terminate, import, export, role/status change | Warning, danger, blocked, loading | ACC-013, ACC-020, ACC-021 |
| CMP-010 | Alert | Inline or page-level feedback | Info, success, warning, error, permission denied | Dismissible where safe, persistent for blockers | ACC-002, ACC-019, ACC-021 |
| CMP-011 | Reminder Row | Due/overdue task, contract, or payment item | Task, contract, payment | Due, overdue, inactive, permission hidden | ACC-021 |
| CMP-012 | History Timeline | Record-local business events | Compact, expanded | Empty, loading, event detail | ACC-014 |
| CMP-013 | Operation Log Table | Administrator global log query | Table with event detail | Empty, loading, denied, query error | ACC-022 |
| CMP-014 | Metric Tile | Basic report summary | Count, sum, status group | Loading, zero, filtered | ACC-018, ACC-023 |
| CMP-015 | Import File Panel | CSV import workflow | Upload, validating, progress, result | Unsupported format, partial success, failure | ACC-020 |
| CMP-016 | Row-Level Error Table | CSV failed row feedback | Compact table | Empty, errors, copy/export if later allowed | ACC-020 |
| CMP-017 | Permission Denied Panel | Safe denial feedback | Page, section, inline | Return action, no restricted details | ACC-002 |
| CMP-018 | Empty State | No authorized data or no results | List, detail section, report, reminder | Role-aware next action | ACC-015, ACC-018, ACC-021, ACC-023 |
| CMP-019 | Action Menu | Secondary record actions | Table row, detail header | Permission-filtered, disabled, loading | ACC-002 |
| CMP-020 | User Table | Administrator user management | User list | Loading, empty, error, selected | ACC-001, ACC-002 |
| CMP-021 | Role Selector | Assign user role | Single role select | Required, disabled, saved, invalid | ACC-001, ACC-002 |
| CMP-022 | Role Summary Table | Show role capabilities | Read-only summary | Loading, denied | ACC-001, ACC-002 |

## Component Details

### CMP-001: App Shell

Structure:
- Left navigation on desktop.
- Collapsible navigation on tablet.
- Top navigation or drawer on mobile.
- Role/account utility area.

States:
- Loading role.
- Permission-filtered navigation.
- Active section.
- Access denied fallback.

Accessibility:
- Keyboard reachable navigation.
- Current section announced through active state.

### CMP-002: Page Header

Structure:
- Title.
- Optional subtitle/context.
- Status badge.
- Primary action.
- Secondary action menu.

Rules:
- Primary action must be the most common allowed action for that screen.
- High-risk actions stay in confirmation flow.
- Header actions must reflow on mobile into stacked or menu layout.

### CMP-003: Toolbar

Structure:
- Search input.
- Filter controls.
- Sort control.
- Refresh.
- Optional density.
- Optional archived filter.
- Optional create action.

States:
- Search loading.
- Invalid filter.
- No results.
- Permission-filtered action visibility.

### CMP-004: Data Table

Structure:
- Column headers.
- Sortable columns where useful.
- Row status badges.
- Row action menu.
- Pagination or load-more area.

Rows must show:
- Primary name.
- Owner or role-relevant owner field.
- Status/stage.
- Key date or amount where entity needs it.
- Warning or overdue indicator where relevant.

Mobile:
- Convert table row to stacked summary row.
- Preserve primary name, status, owner, key date/amount, and row action.

### CMP-005: Detail Header

Structure:
- Entity name.
- Status/stage badge.
- Owner.
- Key amount/date.
- Primary actions.
- More actions.

Variants:
- Normal.
- Archived.
- Terminal Won/Lost.
- Permission denied.

### CMP-006: Form Section

Rules:
- Labels are always visible.
- Required fields are indicated.
- Field-level errors appear near fields.
- Form-level summary appears when multiple errors block save.
- Save action shows loading and then success or error.
- Stale edit conflict shows refresh/discard/retry options.

### CMP-007: Status Badge

Variants:
- Neutral: Draft, Unassigned, Open.
- Info: Sent, Pending Signature, Active.
- Success: Accepted, Paid, Won, Completed.
- Warning: Overdue, Partially Paid, Expired, duplicate warning.
- Danger: Lost, Terminated, Invalid, permission denied.

Rules:
- Badge text must be explicit.
- Color alone is insufficient.

### CMP-008: Stage Stepper

Purpose:
- Show opportunity stage progression.

States:
- Current.
- Completed.
- Blocked.
- Terminal Won.
- Terminal Lost.

Rules:
- Blocked stage actions show missing requirement reason.
- Terminal states disable stage changes.

### CMP-009: Confirmation Modal

Used for:
- Won.
- Lost.
- Archive.
- Contract termination.
- Import run.
- Export run.
- User role/status change.

Required content:
- Target record.
- Action effect.
- Required reason field where applicable.
- History/log notice where applicable.
- Confirm and cancel actions.

Export content:
- Object.
- Filter conditions.
- Authorization scope.
- Estimated record count when available.
- Archived inclusion or exclusion.
- Audit notice.
- No unnecessary sensitive sample data.

Role/status change content:
- Target user.
- Old role and new role.
- Old status and new status.
- Access impact summary.
- Audit/log notice.
- Explicit confirm and cancel actions.

Blocked state:
- Governance-risk operations such as disabling or downgrading the last
  Administrator show a blocked state with no confirm action.
- Confirmation UI does not replace backend authorization or server-side
  governance enforcement.

### CMP-010: Alert

Variants:
- Info.
- Success.
- Warning.
- Error.
- Permission denied.

Rules:
- Blocking alerts remain visible until issue is resolved or user leaves flow.
- Non-blocking warnings must not look like hard errors.
- Error, warning, success, and permission-denied alerts use safe summaries.
  They must not echo restricted record names, restricted field values, full
  contact details, payment values, or sensitive raw before/after values unless
  role authorization allows it.

### CMP-011: Reminder Row

Content:
- Type icon.
- Entity name.
- Due date or expected signed date.
- Status.
- Owner/context.
- Open action.

Rules:
- Unauthorized reminders are hidden.
- Inactive reminders are not shown in active reminder list.

### CMP-012: History Timeline

Content:
- Event type.
- Actor.
- Timestamp.
- Resource.
- Before/after values when available.

Rules:
- Timeline is read-only.
- Empty state names the record context.
- Summary rows avoid sensitive raw values by default.
- Event detail may show before/after values only according to Security data
  classification and the user's role authorization.

### CMP-013: Operation Log Table

Content:
- Event ID.
- Actor.
- Action.
- Resource.
- Timestamp.
- Result.
- Before/after where relevant.

Rules:
- Administrator only.
- Read-only.
- Table summaries avoid sensitive raw values by default.
- Event detail may show before/after values only according to Security data
  classification and Administrator authorization.
- Query errors use safe summaries and do not leak restricted values.

### CMP-014: Metric Tile

Content:
- Metric label.
- Count or amount.
- Filter context.
- Drill-in where authorized.

Rules:
- Empty data shows zero state.
- Do not imply advanced forecasting.

### CMP-015: Import File Panel

Structure:
- Object selector.
- CSV file input.
- Validation summary.
- Run import button.
- Progress.
- Result summary.

States:
- No file.
- Unsupported format.
- Validating.
- Importing.
- Partial success.
- Failed.

### CMP-016: Row-Level Error Table

Content:
- Row number.
- Field.
- Error message.
- Suggested correction when available.

Rules:
- Keep successful row count visible.
- Do not imply failed rows were imported.
- Show row number, field, and validation rule before any row value.
- Do not default to full contact, customer, amount, or payment value echoing.
- Any displayed row value must follow Security data classification and role
  authorization.

### CMP-017: Permission Denied Panel

Content:
- Safe denial message.
- Return action.
- Optional request-access placeholder only if later scoped.

Rules:
- Do not reveal restricted record name or existence unless already authorized.
- Do not provide bypass actions.
- Use a safe denial summary for page, section, and inline variants.

### CMP-018: Empty State

Content:
- Entity-specific empty message.
- Allowed next action if role can create.
- No unauthorized action prompt.

Variants:
- No data.
- No filtered results.
- No authorized records.

### CMP-019: Action Menu

Rules:
- Permission-filter unavailable actions.
- Use disabled state with reason for business-blocked actions where safe.
- High-risk actions require confirmation.

### CMP-020: User Table

Content:
- User name.
- Status.
- Assigned role.
- Last known access status where later available.
- Row action.

Rules:
- Administrator only.
- Non-admin sees permission denied, not a partial user list.

### CMP-021: Role Selector

Rules:
- Required when editing a user.
- Shows Administrator, Sales Manager, Sales.
- Role/status changes require pre-save confirmation before saving.
- Confirmation must show target user, old role, new role, old status, new
  status, access impact, and audit/log notice.
- Granting Administrator, removing Administrator, status changes, and account
  disablement all use the confirmation flow.
- Disabling or downgrading the last Administrator is blocked with no save
  action.
- Save creates visible success feedback and operation log event only after the
  confirmed action succeeds.
- UI confirmation does not replace backend authorization or server-side
  governance enforcement.

### CMP-022: Role Summary Table

Purpose:
- Present role capabilities for Administrator review.

Rules:
- Read-only unless later scope promotes role customization.
- Non-admin denied.

## Visual State Coverage

Every component used in P0/P1 flows must define:
- Default.
- Hover/focus.
- Disabled.
- Loading.
- Empty where applicable.
- Error.
- Success.
- Permission denied where applicable.
- Conflict or blocked where applicable.

## Accessibility Notes

- Focus indicator must be visible.
- Touch/click targets should be at least 44px where practical.
- Error messages must be programmatically associated with fields.
- Tables need meaningful headers and row labels.
- Modal focus must be contained and restored.
- Color contrast must meet WCAG AA targets.
