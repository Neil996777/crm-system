# Screen State Spec

## Document Control

- Project: CRM System
- Phase: G4 UX Design (G4b)
- Owner Agent: UX Designer
- Source: `docs/ux-ui/screen-flows.md`, `docs/ux-ui/interaction-spec.md`
- Status: Accepted as Architecture Input; canonical interactive-state set added
  2026-06-06 (grounded in the locked dashboard), pending re-acceptance.

## Canonical Interactive States (1:1 map for UI Design)

This document is the authoritative source for the **canonical interactive state
names**. UI Design (`design-system.md`) must render exactly these names, and the
modern interaction layer (`interaction-spec.md` Part B) specifies the behavior
and transition into/out of each. Use these names verbatim everywhere:

`loading · empty · error · disabled · selected · focused · hover ·
permission-denied · optimistic-update · success`

- The original six per-screen states (`loading`, `empty`, `error`, `success`,
  `disabled`, `permission-denied`) remain mandatory for every P0/P1 screen in the
  matrix below (no-downgrade: nothing here is weakened).
- The four interaction-level states (`selected`, `focused`, `hover`,
  `optimistic-update`) are cross-cutting component states defined in
  "Cross-Screen State Requirements" below and in `interaction-spec.md` Part B
  (B4). They apply wherever the corresponding interaction exists (selectable
  rows, focusable controls, hoverable surfaces, safe optimistic edits) rather
  than being enumerated per screen.
- Transition behavior (skeleton→content crossfade, optimistic apply + rollback,
  inline-error vs toast, the card→focus hero transition, live-update highlight)
  is owned by `interaction-spec.md` Part B; this file owns the per-screen state
  coverage and the static state contracts.

## State Principles

- Every P0/P1 screen must define loading, empty, error, success, disabled, and
  permission-denied states where relevant.
- Disabled actions must explain what requirement is missing when that
  explanation is safe to reveal.
- Permission-denied states must not reveal restricted record data.
- State design must support QA verification and manual reproduction.
- `selected`, `focused`, `hover`, and `optimistic-update` must be visually and
  behaviorally distinct from each other and from `disabled`/`permission-denied`
  (see `interaction-spec.md` Part B4 and B6).

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
- Stage/payment updates must refresh visible state (DEC-020).
- History-relevant changes should expose a path to record-local history.

### Disabled Or Blocked

- Disabled actions should explain missing required data or permission where
  safe.
- Blocked business transitions must state the business reason, such as a Signed
  contract required before Won (DEC-017).
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

### Selected

- Persistent selection of a row/card/item; drives bulk-action availability.
- Must be visually distinct from `hover` and `focused`.
- Clearable via toggle, Esc, or an explicit clear control.

### Focused

- Keyboard focus produces a visible focus-visible affordance on every
  interactive element.
- Focus must be restored after drawer/stage/modal close and must never be lost
  to the page body after a transition or live update (see Part B5/B6).

### Hover

- Lightweight affordance only (lift/tint); must never be the only way to reveal
  essential information or a critical action.
- Suppressed on touch and under `prefers-reduced-motion`.

### Optimistic-Update

- For safe inline edits only: apply the change in the UI immediately with a
  pending affordance, then reconcile.
- On success, the optimistic value is confirmed (no contradictory flash); on
  failure, **roll back** to the prior value and show an inline `error` with
  retry — feedback must never imply persistence that did not occur.
- Forbidden for business-gated/terminal actions (qualify, stage change, Won/Lost
  close, accept quote, contract create/sign, payment record, archive, owner
  transfer, user/role change) — these use confirm-then-commit with real
  `loading`/`success`/`error` (see Part B4/B7 and Interaction Spec Part A).

### List Archetype (商机 exemplar, generalizes to all CRM lists)

- The full **List Archetype interaction pattern** (data scope by role, debounced
  live-count search + facets, columns/sort, drawer drill-in, selection + bulk
  actions, pagination decision, per-state list behavior, list live-update + pause,
  and list keyboard/a11y) is specified in `interaction-spec.md` **Part B8**
  ("List Archetype"). The Opportunity List, Lead List, Company/Customer List,
  Contact List, Quote List, Contract List, Payment List, and Activity/Note/Task
  List rows in the Screen State Matrix above all render against that one pattern,
  using the canonical state names verbatim. List visuals are owed by
  `design-system.md` (see Part B8.11 + the Part B "UI Handoff" list).

### Live Update (no full reload)

- Auto-refreshing surfaces (dashboard panels, and inheriting entity lists /
  reminder counts) patch changed rows/values in place without a full-page
  reload, preserving scroll, focus, selection, and in-progress edits.
- A changed row/value shows a transient `arrived` highlight that is redundant
  with the already-correct static value (no motion-only signaling).
- A `暂停` (pause) view-preference buffers updates (never drops them) and offers
  a non-blocking "有 N 条新更新" apply affordance. See `interaction-spec.md` B3.

## UX-To-UI Handoff Requirements

UI Design must define:
- Visual hierarchy for dense CRM list/detail work.
- Component states for validation, warnings, disabled actions, success,
  loading, empty, permission denied, blocked, and conflict.
- Visual tokens for the canonical interactive states `selected`, `focused`
  (focus-visible ring), `hover`, and the `optimistic-update` pending affordance,
  each visually distinct.
- Responsive behavior for lists, detail panels, forms, tables, and modals.
- Accessible focus order, labels, and keyboard behavior for all core states.
- The visual treatments enumerated in `interaction-spec.md` Part B
  ("UI Handoff — Visual Treatments Owed"): motion tokens, `arrived` highlight,
  skeleton/shimmer, live-dot states, scrim, strip-card, rail flyout, toast,
  new-updates pill, and hit-area padding — all consistent with the LOCKED
  palette and the approved dashboard mockup.
