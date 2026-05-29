# Responsive Spec

## Document Control

- Project: CRM System
- Phase: G4 UI Design
- Owner Agent: UI Designer
- Source: `docs/ux-ui/ui-spec.md`, `docs/ux-ui/component-spec.md`
- Status: Accepted as Architecture Input

## Responsive Principles

- Desktop is primary for high-volume CRM operations.
- Tablet must preserve list/detail productivity with adaptive navigation.
- Mobile must support inspection, reminders, light edits, and urgent actions
  without hiding core state.
- Responsive behavior must not remove P0/P1 capabilities; if a workflow is
  awkward on mobile, the UI must still provide an accessible path or state that
  explains the limitation.
- Architecture reset on 2026-05-29: implementation is blocked until the restarted delivery flow passes G8.

## Breakpoints

| Name | Width | Layout Rule |
|---|---|---|
| Mobile | < 768px | Single-column layout, drawer navigation, stacked rows, full-width forms. |
| Tablet | 768px to 1199px | Collapsed navigation, two-column where space allows, responsive tables. |
| Desktop | >= 1200px | Sidebar navigation, full data tables, split list/detail where useful. |

## Global Responsive Layout

| Area | Desktop | Tablet | Mobile |
|---|---|---|---|
| Navigation | Persistent left sidebar | Collapsed sidebar or top trigger | Drawer or top navigation |
| Page header | Title, context, status, actions in one row | Title/context with actions wrapping | Stacked title, status, primary action, more menu |
| Toolbar | Inline search/filter/sort/actions | Wrap filters into second row or panel | Search first, filters in drawer/panel, actions in menu |
| Tables | Full table with key columns | Reduced columns plus row expansion | Stacked summary rows |
| Detail | Main content plus side summary where useful | Main content with collapsible side sections | Single-column sections with sticky primary action when safe |
| Modals | Centered modal | Centered or full-width panel | Full-screen sheet for complex forms |

## Screen Responsive Matrix

| Screen | Desktop | Tablet | Mobile | Notes |
|---|---|---|---|---|
| Sign In | Centered panel | Centered panel | Full-width panel with safe margins | Keyboard flow must remain simple |
| Role Workspace / Work Overview | Summary strip plus work/reminder columns | Summary wraps; reminders below or side panel | Stacked summary, reminders, assigned work | Preserve reminder visibility |
| Entity List Pattern | Full data table | Reduced table with expandable row | Stacked row cards with key fields | Must preserve search/filter and open detail |
| Entity Detail Pattern | Header, field groups, related sections, side summary | Related sections become tabs/accordions | Single-column accordions | Primary status and actions remain visible |
| Lead Detail | Form and conversion section side by side if space allows | Conversion below form | Stacked form, qualification, conversion | Required errors remain near fields |
| Customer/Contact Detail | Customer header, contacts table, related tabs | Contacts become compact table/list | Contacts as stacked rows | Contact create remains accessible |
| Opportunity Detail | Stage stepper full width, related sections below | Compact stepper, tabs | Vertical stage list or compact current stage | Terminal states must be obvious |
| Quote Detail | Form plus related contract context | Form then related context | Stacked fields and status actions | Accepted quote indicator remains visible |
| Contract Detail | Status/date/amount sections plus payment links | Sections stack with summary header | Single-column, key dates near top | Pending Signature reminder state remains visible |
| Payment Detail | Payment plan and actual payment side by side | Sections stack | Single-column, remaining amount near action | Overpayment context remains visible |
| Reminder Center | Grouped list with filters | Filters collapse | Type tabs or filter drawer | Related record action remains one tap/click away |
| Manager Overview | Summary metrics, team table, risk panels | Summary wraps; table reduces | Summary and risk list; detail opens separately | Avoid hiding overdue and pending signature risks |
| Import/Export | Object selector, file panel, result table | Result table reduced | Full-screen flow with row errors stacked | Long-running progress visible |
| History And Admin Logs | Timeline/table with filters | Filters collapse | Timeline/list with event detail sheet | Logs remain read-only |
| Reports | Metric grid plus grouped tables | Metric grid wraps | Stacked metrics and grouped lists | Basic reports only |
| Archive Confirmation | Modal with obligation list | Modal or side panel | Full-screen confirmation sheet | Related obligations remain navigable |
| Admin User/Role Management | User table plus detail panel | User list and detail tabs | User list then full-screen user detail | Non-admin denial must not leak data |

## Entity List Responsive Rules

Desktop table columns:
- Primary name.
- Status/stage.
- Owner.
- Key amount or key date.
- Updated time or due date where relevant.
- Row actions.

Tablet columns:
- Primary name.
- Status/stage.
- Owner.
- Key date or amount.
- Row actions.

Mobile row summary:
- Primary name.
- Status/stage badge.
- Owner or role context.
- Key date/amount.
- Warning/overdue indicator.
- More action.

Rules:
- Search remains visible above the list.
- Filters can collapse into a panel or drawer.
- Permission-filtered records do not appear.
- Archived filter must be explicit.

## Detail Responsive Rules

Desktop:
- Header remains at top of content.
- Related sections may use tabs or side panels.
- Long forms are grouped into sections.

Tablet:
- Header actions wrap into a compact action row.
- Related sections use tabs or accordions.
- Tables reduce columns.

Mobile:
- Header stacks title, status, owner, primary action.
- Use accordions for long details.
- Use full-screen sheets for complex create/edit forms.
- Keep destructive actions away from primary save actions.

## Form Responsive Rules

Desktop:
- Two-column forms are allowed for related short fields.
- Long text and notes span full width.

Tablet:
- Two-column only for short paired fields.
- Validation summaries remain above fields.

Mobile:
- Single-column fields.
- Submit actions fixed at bottom only when they do not cover content.
- Error summaries appear before form fields and link to failed fields.

## Data Table Responsive Rules

- Tables must not require horizontal scrolling for critical fields on mobile.
- If horizontal scroll is unavoidable for admin/log/report tables, freeze or
  repeat the primary identifier visually in each row summary.
- Row actions move into a menu on small screens.
- Empty, loading, error, and permission states must occupy the table area.

## Modal And Drawer Responsive Rules

| Component | Desktop | Tablet | Mobile |
|---|---|---|---|
| Confirmation modal | Center modal | Center modal or side panel | Full-screen sheet |
| Filter panel | Inline or popover | Side panel | Drawer |
| User detail | Side panel or detail page | Detail tab/page | Full-screen page |
| Import result | Table area | Reduced table | Stacked row errors |

Rules:
- Modals must preserve focus.
- Full-screen sheets must have clear close/cancel and primary action.
- Blocking reasons must remain visible without scrolling past the confirm
  action.

## Accessibility Responsive Rules

- Focus order follows visual order after reflow.
- Target size should remain at least 44px where practical.
- Text must not be truncated when it changes business meaning.
- Status badges require text, not color alone.
- Error and success messages remain reachable by keyboard and screen reader.

## Responsive Verification Notes

QA and Frontend should verify at minimum:
- Desktop >= 1200px.
- Tablet around 768px to 1199px.
- Mobile around 390px width.
- Entity list/search/filter/open detail.
- Opportunity stage, quote, contract, payment, reminder, archive, admin user
  management, import result, reports, and permission-denied states.
