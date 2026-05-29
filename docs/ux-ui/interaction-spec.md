# Interaction Spec

## Document Control

- Project: CRM System
- Phase: G4 UX Design
- Owner Agent: UX Designer
- Source: `docs/ux-ui/ux-flows.md`, `docs/business/business-rules.md`
- Status: Accepted as Architecture Input

## Interaction Principles

- Interactions must make allowed actions obvious and blocked actions explainable.
- Validation should prevent invalid state transitions without hiding business
  rules.
- Permission denial must be explicit and must not expose restricted data.
- Terminal and archive actions require confirmation.
- UX feedback must not imply persistence unless save succeeds.

## Interaction Matrix

| ID | Priority | Interaction | Trigger | Feedback | Failure / Recovery | Acceptance IDs |
|---|---|---|---|---|---|---|
| IX-001 | P0 | Sign in | Submit credentials | Loading, then role workspace | Invalid credentials keep user on sign-in with error | ACC-001 |
| IX-002 | P0 | Open protected record | Select list item or link | Detail loads if authorized | Permission denied with safe return; no restricted data shown | ACC-002, ACC-015 |
| IX-003 | P0 | Save form | Create/edit submit | Field validation, save loading, success confirmation | Inline errors and form summary; stay in form | ACC-003, ACC-005, ACC-006, ACC-007, ACC-009, ACC-010, ACC-011, ACC-012 |
| IX-004 | P0 | Qualify lead | Select Valid or Invalid | Status confirmation and history update | Invalid requires reason; Unassigned Sales action denied | ACC-004 |
| IX-005 | P0 | Convert lead | Select convert action | Customer/contact/opportunity flow opens | Invalid or converted lead blocked with reason | ACC-004, ACC-007 |
| IX-006 | P0 | Change opportunity stage | Select stage action | Stage change confirmation and history update | Forbidden transition shows reason and required data | ACC-008, ACC-014 |
| IX-007 | P0 | Close opportunity Won | Select Won | Confirmation, then terminal Won state | Blocked until full payment recorded | ACC-013 |
| IX-008 | P0 | Close opportunity Lost | Select Lost | Lost reason prompt, confirmation, terminal Lost state | Lost reason required | ACC-013 |
| IX-009 | P0 | Accept quote | Select Accept | Accepted quote indicator | Prevents multiple Accepted quotes for same opportunity | ACC-009 |
| IX-010 | P0 | Create contract | Select Accepted quote | Contract form opens with quote/customer/opportunity context | Expired quote blocked; amount mismatch needs reason | ACC-010 |
| IX-011 | P0 | Sign contract | Change to Signed | Signed/effective date prompt and confirmation | Missing signed/effective date blocks transition | ACC-010 |
| IX-012 | P0 | Record payment | Submit payment | Payment status updates | Zero, negative, or overpayment blocked | ACC-011 |
| IX-013 | P0/P1 | Create task | Submit task | Task appears in related record and reminders if due | Missing title/owner/due date/status blocked | ACC-012, ACC-021 |
| IX-014 | P1 | Open reminder | Select reminder | Related record detail opens | Unauthorized related record hidden or denied | ACC-021 |
| IX-015 | P1 | Resolve reminder | Complete task, sign contract, or record payment | Reminder list updates | Inactive reminder removed from active list | ACC-021 |
| IX-016 | P1 | Duplicate warning | Enter matching company/contact/lead data | Non-blocking warning | User can continue save; no merge or overwrite | ACC-019 |
| IX-017 | P1 | CSV import | Submit CSV | Progress, then success and row-level errors | Unsupported format or invalid rows reported | ACC-020 |
| IX-018 | P1 | CSV export | Select export | Export starts and completes | Unauthorized export denied | ACC-020 |
| IX-019 | P0/P1 | Archive record | Select archive | Confirmation and success | Blocked by active obligations with related links and retry | ACC-002, ACC-014, ACC-015, ACC-021 |
| IX-020 | P0/P1 | View history/logs | Open history or admin logs | Timeline or log results | Log edit unavailable; non-admin global log denied | ACC-014, ACC-022 |
| IX-021 | P1 | View reports | Open reports | Report groups load | Empty state; unauthorized records excluded | ACC-018, ACC-023 |
| IX-022 | P0 | Owner transfer | Manager transfers owner | Transfer confirmation and affected task note | Unauthorized transfer denied | ACC-002, ACC-014 |
| IX-023 | P0 | Search/filter entity list | Enter query or filter | Authorized filtered results update | Invalid filter keeps input and shows correction message | ACC-015 |
| IX-024 | P0 | Open entity detail | Select authorized list row | Entity detail opens | Unauthorized entity hidden or denied without restricted data | ACC-015, ACC-002 |
| IX-025 | P0 | Manage user account | Administrator edits user status or role | Save confirmation and operation-log event | Non-admin denied; failed save preserves input | ACC-001, ACC-002, ACC-022 |
| IX-026 | P0 | View role capability summary | Administrator opens role detail | Role capabilities shown read-only unless later configured | Non-admin denied | ACC-001, ACC-002 |

## Validation Behavior

- Required fields validate on submit and after field blur where helpful.
- Status transitions validate before mutation.
- Amount validation shows remaining amount where overpayment is blocked.
- Contract date validation depends on contract status.
- Duplicate warning appears before final save but does not block save.
- Entity list filters validate before applying when the filter value has a
  constrained format.
- User/role management validates required identity and role fields before save.

## Confirmation Behavior

Confirmation is required for:
- Won closure.
- Lost closure.
- Archive.
- Contract termination.
- Import run after file validation.

Confirmation must show:
- Target record.
- Business effect.
- Whether history/log events will be created.
- Any irreversible or terminal consequence.

## Recovery Behavior

- Validation failure keeps user input intact.
- Permission denial provides return to previous safe screen.
- Partial CSV import failure shows failed rows and successful row count.
- Blocked archive lists active downstream obligations and links to them.
- Failed save must not display success.

## Accessibility Requirements

- All form controls require accessible labels.
- Error summaries must be keyboard focusable and link to fields.
- Modal confirmations must trap focus and return focus after close.
- Lists and tables must support keyboard navigation and meaningful row labels.
- Status messages must be announced to assistive technology where applicable.
