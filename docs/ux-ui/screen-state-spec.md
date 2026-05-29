# Screen State Spec

## Document Control

- Project: CRM System
- Phase: G4 UX Design
- Owner Agent: UX Designer
- Source: `docs/ux-ui/screen-flows.md`, `docs/ux-ui/interaction-spec.md`
- Status: Accepted as Architecture Input

## State Principles

- Every P0/P1 screen must define loading, empty, error, success, disabled, and
  permission-denied states where relevant.
- Disabled actions must explain what requirement is missing when that
  explanation is safe to reveal.
- Permission-denied states must not reveal restricted record data.
- State design must support QA verification and manual reproduction.

## Screen State Matrix

| Screen | Loading | Empty | Error | Success | Disabled | Permission Denied |
|---|---|---|---|---|---|---|
| Sign In | Authenticating credentials | N/A | Invalid credentials or disabled account | Role workspace opens | Submit disabled until required fields present | Unauthenticated protected access returns here |
| Role Workspace | Loading role and navigation | No assigned active work | Role load failed | Role-scoped sections visible | Unauthorized sections hidden/unavailable | No core data exposed |
| Work Overview | Loading assigned records/reminders | No assigned active work | List failed to load | Assigned work visible | Actions unavailable without selection | Unauthorized records hidden |
| Lead List | Loading leads | No authorized leads | Search/filter error | Authorized leads listed | Bulk actions unavailable when none selected | Non-owned leads hidden for Sales |
| Lead Detail | Loading lead | N/A | Load or save failed | Save/qualification confirmed | Qualify disabled when required data missing | Unauthorized lead denied |
| Company/Customer List | Loading companies/customers | No authorized companies/customers | Search/filter error | Authorized companies/customers listed | Bulk actions unavailable when none selected | Unauthorized records hidden |
| Customer Detail | Loading customer/contact data | No contacts yet | Save/link failed | Customer/contact saved | Contact save disabled until required fields present | Unauthorized customer denied |
| Contact List | Loading contacts | No authorized contacts | Search/filter error | Authorized contacts listed | Bulk actions unavailable when none selected | Unauthorized records hidden |
| Contact Detail | Loading contact | N/A | Load or save failed | Contact saved | Save disabled until required fields present | Unauthorized contact denied |
| Opportunity List | Loading opportunities | No authorized opportunities | Search/filter error | Authorized opportunities listed | Bulk actions unavailable when none selected | Unauthorized records hidden |
| Opportunity Detail | Loading opportunity | N/A | Stage/save failed | Stage/history updated | Forbidden transition disabled or blocked | Unauthorized opportunity denied |
| Quote List | Loading quotes | No authorized quotes | Search/filter error | Authorized quotes listed | Create disabled until opportunity/customer context exists | Unauthorized records hidden |
| Quote Detail | Loading quote | No quotes for opportunity | Status/save failed | Quote saved/status updated | Accept disabled for invalid/expired quote | Unauthorized quote denied |
| Contract List | Loading contracts | No authorized contracts | Search/filter error | Authorized contracts listed | Create disabled until Accepted quote exists | Unauthorized records hidden |
| Contract Detail | Loading contract | No contract yet | Save/status failed | Contract saved/status updated | Status action disabled until required date/reason present | Unauthorized contract denied |
| Payment List | Loading payments | No authorized payment records | Search/filter error | Authorized payments listed | Create disabled until contract exists | Unauthorized records hidden |
| Payment Detail | Loading payments | No payment plan/payment yet | Invalid amount or save failed | Payment recorded/status updated | Record disabled until required amount/date/status present | Unauthorized payment denied |
| Activity/Note/Task List | Loading activity/note/task records | No authorized activity/note/task records | Search/filter error | Authorized records listed | Create disabled until related record exists | Unauthorized records hidden |
| Activity/Note/Task Detail | Loading activity/note/task detail | N/A | Load or save failed | Record saved | Save disabled until required fields present | Unauthorized record denied |
| Reminder Center | Loading reminders | No active reminders | Reminder load failed | Related record opens or reminder resolves | Inactive reminders unavailable | Unauthorized reminders hidden |
| Manager Overview | Loading team data | No team records | Summary failed to load | Team records and summaries visible | Team actions disabled without selection | Sales denied |
| Admin User/Role Management | Loading users and roles | No user records configured | Query or save failed | User/role changes saved | Save disabled until required fields present | Non-admin denied |
| User Detail | Loading user | N/A | Load or save failed | User status/role saved | Invalid status/role changes disabled | Non-admin denied |
| Role Detail | Loading role | N/A | Load failed | Role capabilities shown | Editing unavailable unless allowed by later scope | Non-admin denied |
| Import/Export | Loading form or export | No import file selected | Unsupported format or row errors | Import/export summary shown | Run disabled until object/file/scope valid | Sales denied |
| Record-local History | Loading history | No history events | History load failed | Timeline shown | Edit unavailable | Unauthorized history denied |
| Admin Logs | Loading logs | No log events match filter | Query failed | Log results shown | Edit unavailable | Non-admin denied |
| Reports | Loading report data | Zero/empty report state | Report load/filter error | Report summaries shown | Drill-in unavailable when no records | Sales denied; unauthorized records excluded |
| Archive Confirmation | Checking obligations | No active obligations | Obligation check failed | Archive success | Confirm disabled when active obligations block archive | Sales denied |

## Cross-Screen State Requirements

### Loading

- Preserve layout position where possible.
- Do not show stale success after a new mutation starts.
- Long-running CSV import must show progress or pending state.

### Empty

- Empty states must name the missing data type.
- Empty states must show allowed next action where the role can create data.
- Empty states must not suggest unauthorized actions.

### Error

- Field-level errors appear near fields.
- Form-level summaries list blocking issues.
- Save failures keep entered data.
- Search/filter errors preserve filter inputs for correction.

### Success

- Success feedback must identify the saved or changed record.
- Stage/status/payment updates must refresh visible status.
- History-relevant changes should expose a path to record-local history.

### Disabled Or Blocked

- Disabled actions should explain missing required data or permission where
  safe.
- Blocked business transitions must state the business reason, such as full
  payment required before Won.
- Archive blocked by active obligations must list related obligations and
  provide entry points.

### Permission Denied

- Show denial without restricted data.
- Provide return to previous safe context.
- Do not indicate whether a hidden record exists unless already authorized by
  another path.

### Conflict

- If a record changes while a user is editing, UX must prevent silent overwrite.
- User should be offered refresh, discard, or retry after review.
- Conflict handling must preserve no-downgrade behavior for persistence and
  history.

## UX-To-UI Handoff Requirements

UI Design must define:
- Visual hierarchy for dense CRM list/detail work.
- Component states for validation, warnings, disabled actions, success,
  loading, empty, permission denied, blocked, and conflict.
- Responsive behavior for lists, detail panels, forms, tables, and modals.
- Accessible focus order, labels, and keyboard behavior for all core states.
