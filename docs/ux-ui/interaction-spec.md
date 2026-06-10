# Interaction Spec

## Document Control

- Project: CRM System
- Phase: G4 UX Design (G4b)
- Owner Agent: UX Designer
- Source: `docs/ux-ui/ux-flows.md`, `docs/business/business-rules.md`,
  locked dashboard mockups (`docs/ux-ui/mockups/dashboard-v7-*.png` and their
  `_src/*.html`)
- Status: Accepted as Architecture Input; modern interaction layer (Part B)
  added 2026-06-06, grounded in the locked dashboard, pending re-acceptance.
  List Archetype (Part B8, 商机 exemplar) instantiated 2026-06-06 on the accepted
  Part B + accepted decisions (DEC-UX-NAV-01 / DEC-UX-MOTION-02 /
  DEC-UX-LIVE-03 / DEC-UX-LIVE-04 / DEC-UX-FOCUSRAIL-01 /
  DEC-UX-FOCUSEXIT-01 / DEC-UX-HEROTIME-01).

## Reading Order

- **Part A — Interaction Matrix** (below, unchanged): the business-rule-driven
  per-action contracts. No-downgrade: these requirements stand as-is.
- **Part B — Modern Interaction Layer** (new section near the end of this file):
  the reusable motion, state-transition, navigation, and cross-archetype
  patterns that make the product feel current ("科技感") without trading away
  usability or accessibility. Part B extends Part A; it never weakens it. Part B8
  ("List Archetype") instantiates the LIST page archetype concretely for the
  商机 list as the reusable pattern for all CRM list pages.
- **Canonical interactive states** live in `screen-state-spec.md` (the 1:1 map
  UI Design renders against). Part B references those state names; it does not
  redefine them.

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
| IX-007 | P0 | Close opportunity Won | Select Won | Confirmation, then terminal Won state | Blocked until the related contract is Signed (DEC-017) | ACC-013 |
| IX-008 | P0 | Close opportunity Lost | Select Lost | Lost reason prompt, confirmation, terminal Lost state | Lost reason required | ACC-013 |
| IX-009 | P0 | Accept quote | Select Accept | Accepted quote indicator | Each opportunity has exactly one quote (DEC-018) | ACC-009 |
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

---

# Part B — Modern Interaction Layer (Motion, States, Navigation, Archetypes)

This part formalizes the interaction model the locked dashboard converged on and
extends it into a coherent, reusable system. It owns **behavior, flow, and
motion**. It does NOT redefine visuals: every place a behavior needs a visual
treatment (a color, a shadow, a shimmer gradient, a focus-ring style) names the
treatment by intent and **defers the rendered values to UI Design's
`design-system.md`**. Those deferrals are collected in "UI Handoff — Visual
Treatments Owed" at the end of Part B so UI can reconcile 1:1.

## B0. Design Principles For Motion ("科技感", restrained)

- **Every animation must earn its place.** An animation is allowed only when it
  does one of: (1) clarifies a state change, (2) shows spatial/object
  continuity, or (3) gives direct feedback to a user action. Decorative motion
  with no informational job is **rejected**.
- **No motion may delay task completion.** Motion is feedback layered on top of
  an already-committed state change, never a gate in front of it. The user can
  act again before an animation finishes; in-flight animations interrupt
  cleanly (no "wait for the slide to end").
- **Continuity over teleport.** When an object changes role on screen (a grid
  card becoming the focused stage, a row collapsing into a strip card), it moves
  rather than disappears-and-reappears, so the user keeps object identity.
- **Performance guardrail.** Animate `transform` and `opacity` only. No
  animation of `width`, `height`, `top`, `left`, `margin`, or anything that
  triggers layout/reflow. Size/position changes are expressed as transforms of a
  pre-laid-out target. Target 60fps; never block the main thread on motion.
- **Reduced motion is a first-class path, not a fallback bolt-on.** See B6.
- **Restraint budget.** At most one "hero" transition on screen at a time (the
  card→focus stage transition is the hero). Everything else is micro-feedback
  ≤ 220ms. Simultaneous competing animations are rejected.

## B1. Motion Token Scale

UX defines the timing/easing scale below as **behavioral tokens** (names +
durations + curves + usage rules). UI Design must mirror these exact names and
values in `design-system.md` so motion is consistent across every surface. If UI
needs to adjust a value for rendering reasons, it changes the token here first
(no per-component one-offs).

### Durations

| Token | Duration | Use for |
|---|---|---|
| `motion-instant` | 80ms | State toggles with no spatial travel: button press, checkbox/toggle thumb, badge swap, focus-ring appear. |
| `motion-fast` | 140ms | Small local moves and feedback: hover-lift, chevron rotate, tooltip, inline value flash start, toast slide-in. |
| `motion-base` | 220ms | Standard transitions: drawer/side-panel open, strip-card collapse, nav rail expand/collapse, skeleton→content crossfade. |
| `motion-slow` | 320ms | The hero card→focus stage transition and its reverse only. Reserved; do not use elsewhere. |

Rule: if a motion would need longer than `motion-slow`, it is too big — split it
or remove it. Stagger (sequencing several elements) uses a 24–40ms step, capped
so the whole staggered group still completes within its token budget.

Exception: per DEC-UX-HEROTIME-01, the Card→Focus hero transition keeps the B1
scale intact but uses dedicated hero timing values for this one signature motion:
enter ~450ms and reverse ~310ms. `motion-base` remains 220ms for focus-rail
switching and all standard/micro transitions; `motion-slow` remains 320ms and is
not globally redefined.

### Easing Curves

| Token | Curve | Use for |
|---|---|---|
| `ease-standard` | `cubic-bezier(0.2, 0, 0, 1)` | Default for moves that both start and end on screen (the stage transition, drawer open). |
| `ease-decelerate` | `cubic-bezier(0, 0, 0, 1)` | Elements entering the screen (toast in, panel appear, content revealing after load). Fast start, soft landing. |
| `ease-accelerate` | `cubic-bezier(0.4, 0, 1, 1)` | Elements leaving the screen (toast out, dismissed panel). Eases in, exits quickly. |
| `ease-emphasis` | `cubic-bezier(0.2, 0, 0, 1)` w/ slight overshoot allowed only on `motion-instant` press feedback | Tactile press only. Never on layout-scale motion. |

Default if unspecified: `motion-base` + `ease-standard`.

### What Animates vs. What Does Not

- **Animate:** focus/stage transition, nav rail width change (as transform),
  strip-card collapse, drawer/side-panel, skeleton→content, row highlight flash
  on live update, toast in/out, hover-lift, press, chevron, toggle thumb,
  live-dot pulse, optional value count-up (B3).
- **Never animate:** text reflow, table column reflow, page scroll position on
  data refresh, anything that moves content the user is reading, validation that
  must be read immediately (errors appear instantly, no fade that delays
  reading), permission-denied messaging (appears instantly).

## B2. Hero Transition — Card Grid → Focus Stage

This is the locked dashboard behavior, formalized. It is the product's signature
"科技感" moment and the template for any overview→detail drill-in that stays on
one screen (see B7 "inline drill-in").

### Layout end-states (from the locked mockup)

- **Overview state:** left sidebar expanded `248px` (icon + text); content is an
  equal-card grid (`roleGrid`, 3 columns of equal-height panels). No scrim.
- **Focus state:** left sidebar collapsed to a `72px` icon-only rail; content
  becomes a two-column `1fr / 300px` layout — the clicked panel expands to the
  left **stage** (full detail: large funnel, data table, tools), and the right
  column is a persistent **focus selector rail** of compact strip cards (`92px`
  tall, title + single key value + live dot). Per DEC-UX-FOCUSRAIL-01, the rail
  lists the full dashboard card set for the current workspace and never removes
  or reorders items while focus changes. For the locked manager dashboard this is
  all **8** cards including the currently-focused one; role-scoped variants list
  their full authorized card set and must not introduce unauthorized manager-only
  cards. The focused card has a visible selected state and `aria-current="true"`.
  A subtle scrim
  (`rgba(15,23,42,.06)`) overlays the workspace area below the topbar to push the
  stage forward. Topbar persists unchanged.

### Trigger

- Click/Enter/Space on a grid panel's expand affordance OR anywhere on the panel
  body (the whole card is the target; the corner `expand` glyph is the visible
  affordance). Each panel is a real button/`role="button"`, focusable, with an
  accessible name like `展开「我的销售漏斗」`.

### Motion choreography (total ~450ms per DEC-UX-HEROTIME-01, `ease-standard`)

Expressed as transform/opacity only; the focus layout is pre-computed and
elements are transformed from their overview positions into it.

1. **0–80ms:** scrim fades in (opacity 0→1, `ease-decelerate`); the clicked card
   lifts (shadow elevates, `motion-instant`) to signal "this one is chosen".
2. **40–450ms (overlaps):** the clicked card translates/scales toward the stage
   slot and its inner detail content crossfades from compact to full
   (`ease-standard`). The card is the same DOM object moving — object continuity,
   not a new element.
3. **~96–450ms (staggered ~30ms step):** the dashboard card set translates toward
   the right selector rail and scales down to strip-card size, fading dense
   internals to the single key value. The rail includes the active card and
   preserves the original dashboard order. Once focus mode is entered, the rail
   remains stable: no item is added, removed, or reordered during in-focus
   switching.
4. **~96–450ms:** the sidebar collapses 248→72px, expressed as a transform of the
   label column to opacity 0 + a rail width change driven by the grid track (the
   labels fade `motion-fast`; icons stay put — they do not move horizontally, so
   the eye keeps anchor).
5. **Settle:** focus-visible ring moves to the stage's primary heading / first
   actionable control (see B5 focus order). Live regions announce the change
   (B6).

### Reverse (返回 / Esc)

- Symmetric reverse on the dedicated hero-exit timing (~310ms, still faster than
  entry and preserving the standard "exit is quicker than enter" feel),
  `ease-standard`. Scrim fades out `ease-accelerate`. Strip cards expand back to
  grid panels in original positions; sidebar re-expands; focus returns to the
  grid panel that was opened (focus restoration, B5).

### Entry/exit + keyboard

- **Enter focus:** click/Enter/Space on a panel.
- **Switch focus:** click/Enter/Space on any right-rail selector item. The left
  stage swaps to that panel via a quick content crossfade on `motion-base`; the
  rail item set and order do not change. The newly focused rail item receives the
  selected treatment and `aria-current="true"`; the previous one clears it.
- **Exit focus:** a single `返回` button (always present, top-right of the stage,
  with the back-chevron icon) or the `Esc` key (global while in focus state).
  Switching focus is not exit; it is never modeled as a full reverse+forward.
- Per DEC-UX-FOCUSEXIT-01, the visible `Esc 退出` / `Esc 返回` hint chip from the
  mockup is removed. The keyboard shortcut remains functional, but no screen
  text advertises it in the stage header.
- Only one focus stage exists at a time; focus is single-select.

### States this transition must honor

- If the clicked panel's data is still **loading**, the stage opens immediately
  with the skeleton state (B4 `loading`) — never block the transition on data.
- If the panel hits **error** after opening, the stage shows the inline error
  state (B4 `error`) with retry; the transition itself is not reverted.
- If a role is **not permitted** to expand a panel (rare; most panels are
  role-scoped already), the panel is rendered in `disabled` state and is not a
  focus trigger; it never opens to a `permission-denied` stage. Permission
  scoping happens before render, per existing Part A IX-002.

## B3. Live / Real-Time Updates (no jarring reload)

Formalizes the `实时更新 / 暂停` segmented toggle + live-dot model. This is how
ALL auto-refreshing surfaces behave (dashboard panels today; entity lists and
reminder counts inherit it).

### Refresh model

- **No full-page reload, ever, on a live update.** Updates patch in place at the
  row/value granularity. Scroll position, focus, selection, expanded panels, and
  in-progress edits are preserved across a refresh (an update must never move
  content the user is reading or steal focus).
- **Changed-row highlight ("arrived" flash):** when a row/value changes, it gets
  a transient highlight — the locked `arrived` treatment (left accent bar + tint
  background) holds for ~1.6s then fades out over `motion-base`. This is the only
  signal that says "this just changed." It is additive: the underlying state
  (badge/number) is already correct and readable without the flash (no
  motion-only signaling — see B6).
- **Layout stability:** new rows enter at the top with a height-stable insert
  (the inserted row animates opacity/translateY only; surrounding rows do not
  reflow-jump — space is reserved, then the row fades/slides into it). If many
  rows arrive at once, they batch (see coalescing) and the list does not thrash.
- **Value changes** (e.g. a KPI metric, a payment total) may use an optional
  count-up tween (B-micro) capped at `motion-base`; if reduced-motion is on, the
  value snaps.

### Debounce / coalesce

- Incoming updates are **coalesced** within a short window (recommended 800ms–1s
  debounce window — flagged as a tuning decision, see Decisions) so a burst of
  events produces one visual update, not a flicker storm.
- At most one highlight flash per row per coalescing window. A row that changes
  twice quickly flashes once with the final value.
- The "更新于 HH:MM" / "数据更新时间" meta in the header reflects the last applied
  batch, not each event.

### Pause behavior (`暂停`)

- Toggling to `暂停`: stops applying live updates to the view. The live-dot stops
  pulsing and goes to a neutral/paused visual (UI owes the paused dot
  treatment). No new highlights occur.
- While paused, incoming changes are **buffered, not dropped**. A non-blocking
  affordance appears: e.g. a pill `有 N 条新更新 · 点击刷新` (zh-CN). Activating
  it applies the buffered batch with the normal highlight flash and returns the
  meta timestamp to live.
- Resuming `实时更新` applies the buffer (with highlights) and continues live.
- Pause is per-user, per-surface, and is a view preference only — it never
  changes server state or stops data from being recorded; it only controls
  whether the view auto-applies.

### SSE-pushed events (parked future feature: "已回款" event)

- The interaction contract is defined now so the future push path has a landing,
  but the feature itself stays parked (do not implement; this is UX intent for
  Architecture/Domain to note, not a new P0/P1 scope item).
- Per DEC-UX-LIVE-04, dashboard live behavior is implemented in the current
  frontend via client-side polling of existing GET endpoints. This is the
  transition mechanism until a future server-push/SSE path is formally scoped.
- A server-pushed event (e.g. a payment "已回款") surfaces **exactly like a poll
  update**: it patches the relevant row/value in place and triggers the single
  `arrived` highlight. It does NOT pop a modal, does NOT reorder beyond the
  defined top-insert, and does NOT steal focus.
- If the relevant surface is paused, the pushed event goes into the same buffer
  and increments the `有 N 条新更新` pill.
- If the user is not on the relevant surface, escalation to the Reminder Center /
  notification path is the existing notification flow (UX-007), NOT a new
  interrupt. SSE is a delivery mechanism, not a new UX surface.
- Connection loss on the SSE/live channel degrades gracefully to the existing
  `offline/reconnect` handling: the live-dot shows a reconnecting state, the view
  keeps last-known data (no blank-out), and a quiet "连接已恢复，已更新" applies
  the catch-up batch on reconnect. (Reconnect visual owed to UI.)

## B4. Canonical Interactive States + Transitions

The canonical state **names** and their per-screen presence are owned by
`screen-state-spec.md` (Part A there). Here we specify, once, the **behavior and
the transition into/out of each state** as reusable rules, using EXACTLY these
names so UI maps 1:1:

`loading · empty · error · disabled · selected · focused · hover ·
permission-denied · optimistic-update · success`

| State | Trigger | Behavior | Transition IN | Transition OUT / Recovery |
|---|---|---|---|---|
| `loading` | A surface/region awaits data or a mutation result. | Show a **skeleton** matched to the real content's shape (not a spinner for content regions; spinner only for in-button action). Skeleton uses a subtle shimmer. Preserve layout (no jump when data lands). Do not show stale success underneath. | Skeleton appears instantly (no fade-in delay). Shimmer loops `motion-base`-paced. | Skeleton → content **crossfade** `motion-fast` `ease-decelerate`; content does not slide the page. On reduced-motion, snap. |
| `empty` | Authorized query returns zero records / first-run. | Name the missing data type; offer the allowed next action for the role; never suggest unauthorized actions (per Part A). Calm, not an error. | Replaces skeleton with a single crossfade `motion-fast`. | When data arrives (e.g. via create or live update), empty → list with the new row using the `arrived` insert. |
| `error` | Load/mutation/filter fails. | **Inline** error in the failing region (not a toast) when the user is looking at that region — keep their input/context, offer 重试. A toast is used only for background/async failures the user isn't currently watching. Errors appear **instantly** (no animated delay before they're readable). | Error block replaces skeleton/content in place; for a failed mutation it replaces the `optimistic-update` (rollback, below). | 重试 → back to `loading`. Errors are dismissible only by retry or by the user navigating; never auto-dismiss a blocking error. |
| `disabled` | Action not currently allowed (missing prerequisite data or business precondition), but the reason is safe to reveal. | Control is non-interactive AND explains the missing requirement (tooltip/inline helper) per Part A. Must not look identical to a permitted control. Not focus-trappable as an action, but discoverable/inspectable (the explanation is reachable). | No motion; state is static. | Becomes enabled with an `motion-instant` affordance change (e.g. ring/opacity) the moment its precondition is met. |
| `selected` | User selects a row/card/item (single or in a bulk-select set). | Persistent visual selection (UI owes the selected token: tint + accent). Drives bulk-action bar (B7). Distinct from `hover` and `focused`. | `motion-instant` tint apply. | Deselect on toggle/Esc/clear; `motion-instant` tint remove. |
| `focused` | Keyboard focus lands on an interactive element. | A clearly visible **focus-visible ring** (keyboard-driven; not shown for mouse unless needed). Drives focus order (B5). The hero stage transition moves focus deliberately (B2). | Ring appears `motion-instant`; may travel `motion-fast` between siblings (subtle), snaps on reduced-motion. | Ring moves with focus; restored on drawer/stage close. |
| `hover` | Pointer over an interactive surface. | Lightweight affordance only: card **hover-lift** (translateY -2px + shadow step, `motion-fast`), row tint (the `hovered` treatment from the mockup), button background step. Hover is never the only way to reveal essential info or actions (no hover-only critical actions). | `motion-fast` `ease-standard`. | Reverse `motion-fast`. Disabled on touch/reduced-motion (no lift). |
| `permission-denied` | User reaches a resource the role may not see. | Explicit denial, **no restricted data**, safe return path (per Part A IX-002). Appears instantly; no teasing animation. Never reached via the hero transition (B2). | Replaces content with the denial state directly (no reveal animation). | Return-to-safe-context control; focus moves to it. |
| `optimistic-update` | User commits an inline edit / toggle / quick status change where latency would otherwise stall them. | **Apply the change in the UI immediately** (so the product feels instant/"科技感") while the request is in flight; show a subtle pending affordance (e.g. a quiet inline spinner or reduced-opacity badge) so the user knows it isn't yet confirmed. Optimistic apply is allowed ONLY where Part A does not require a server precondition/confirmation. **Terminal/irreversible/business-gated actions (Won/Lost close, archive, contract sign, payment record, qualify) are NEVER optimistic** — they follow Part A confirm-then-commit with a real `loading`/`success` cycle. | Value updates `motion-instant`; pending affordance appears. | **Success:** pending affordance resolves to `success`; the optimistic value is now confirmed (no second flash). **Failure:** **roll back** the optimistic value to its prior state with a `motion-fast` revert + an inline `error` explaining the failure and offering 重试. The user must never be left believing a failed change persisted (operating-model: UX feedback must not imply persistence unless save succeeds). |
| `success` | A mutation confirms (server acknowledged). | Identify the saved/changed record (per Part A). Lightweight, non-blocking confirmation: inline check/state flip for in-context saves; a **toast** for actions whose result isn't visible in the current viewport. Auto-dismiss success toasts (~3–4s) — success is not a blocker. Refresh the affected visible state (DEC-020). | `motion-fast` `ease-decelerate` (toast slides in / inline check appears). | Toast auto-dismisses `ease-accelerate`; inline success persists as the new resting state. |

### State precedence (when several could apply)

`permission-denied` > `error` > `disabled` > `loading` > `empty` >
`optimistic-update` > `selected` > `focused` > `hover` > `success` resting.
(Security/permission and failure always win over progress/affordance states.)

## B5. Navigation & Focus Interaction

### Either-collapsed-OR-expanded rule (locked)

The left navigation is in exactly one of two modes, never a hybrid:

- **Expanded** `248px` icon + text — the overview/default mode.
- **Collapsed** `72px` icon-only rail — the focus-stage mode.

The collapse/expand is driven by entering/leaving the focus stage (B2), animated
as part of the hero transition. The two modes are mutually exclusive; there is no
half-width state and no simultaneous "icons + floating labels" persistent state.

### Hover-expand decision (collapsed rail) — RECOMMENDATION, flagged for user

**Recommendation: YES — the collapsed `72px` rail hover-expands to a temporary
icon+text flyout overlay, and this does NOT violate the either/or rule.**

Rationale:

- Icon-only rails have a real legibility/recognition cost; 14 nav items
  (工作台, 线索, 公司客户, 联系人, 商机, 报价, 合同, 回款, 任务, 提醒中心,
  报表, 导入导出, 用户与角色, 操作日志) cannot all be unambiguously recognized
  from a glyph, especially the 管理 items. A temporary on-demand label flyout
  removes guesswork without permanently widening the rail.
- It does not break the either/or rule because the rule governs the **persistent
  layout state** (what occupies grid space). A hover/focus flyout is a transient
  **overlay** that does not reflow the stage, does not change the `72px` track,
  and disappears on mouse-out / blur. The persistent layout is still strictly
  collapsed.
- Behavior spec for the flyout:
  - Trigger: pointer hover over the rail OR keyboard focus entering a rail icon
    (so keyboard users get the label too — not hover-only, per B6).
  - Appearance: labels slide out as an overlay flyout `motion-fast`
    `ease-decelerate`, anchored to each icon; the active item is indicated. It
    floats above the scrim/stage; it never pushes the stage.
  - Dismissal: mouse-out (with a small ~150ms close delay to avoid flicker) or
    focus leaving the rail; `ease-accelerate` out.
  - Reduced motion: flyout appears/disappears with no slide (snap), still on
    hover/focus.
  - It is purely a label aid — clicking an item still navigates exactly as the
    icon does.
- **Alternative (rejected as default): stay strictly icon-only** with only
  per-icon tooltips. Lower implementation cost and zero overlay, but weaker
  legibility for a 14-item zh-CN nav and worse for new users. Acceptable only if
  the user prefers minimalism over discoverability.

This is **DEC-UX-NAV-01 — ACCEPTED by user 2026-06-06: YES** (hover/focus flyout).
See Decisions. Now binding for all archetypes.

### Keyboard focus order

- **Global landmark order:** skip-link → primary nav (rail/sidebar) → topbar
  (search, quick links, user menu) → main content → (stage tools when in focus
  mode). Standard landmark roles so screen-reader users can jump.
- **Overview grid:** Tab order follows visual reading order (KPIs row, then
  `roleGrid` left-to-right, top-to-bottom). Each panel is one Tab stop (its
  expand action); inner interactive elements are reachable after entering.
- **On entering focus stage:** focus moves to the stage heading region / first
  actionable control (not back to top); the rail and topbar remain reachable by
  Shift-Tab.
- **Esc / 返回:** Esc always exits the focus stage from anywhere within it;
  focus **restores** to the grid panel that was opened (focus restoration is
  mandatory, not optional).
- **Drawers/side-panels (B7):** open → focus moves into the drawer and is
  **trapped** there; Esc closes; focus returns to the trigger row/control.
- **focus-visible affordance:** every interactive element shows the
  `focused`-state ring on keyboard focus; the ring must meet contrast against
  both light card and tinted/selected backgrounds (contrast value owed to UI).

## B6. Accessibility (binding, not advisory)

- **`prefers-reduced-motion: reduce` is mandatory.** When set: disable the hero
  transition's travel/scale (state still changes — the stage just appears,
  cross-faded or snapped); disable hover-lift, count-up, live-dot pulse animation
  (the dot stays visible, just not pulsing), strip-card collapse travel, toast
  slide (toast still appears/disappears, just no slide), and skeleton shimmer
  (skeleton stays as a static placeholder). **No information is lost in
  reduced-motion mode** — only the motion is removed. This is the primary
  fallback and must be tested as its own path.
- **No motion-only signaling.** The `arrived` highlight, live-dot pulse, and any
  flash are always redundant with a static, readable signal (the correct
  badge/number/timestamp is already present). A user who never sees the animation
  still gets the same information.
- **Keyboard operability.** Every action reachable by pointer is reachable by
  keyboard: focus stage enter/exit, rail navigation + flyout labels, live
  pause/resume toggle, "有 N 条新更新" apply pill, bulk-select, drawer open/close,
  inline edit. No hover-only or pointer-only action.
- **Focus visibility & restoration.** Visible focus ring on all interactive
  elements; focus is never lost to `body` after a transition, drawer close, or
  live update; focus is restored to the logical origin on Esc/close.
- **Live region announcements.** Entering the focus stage announces the focused
  panel ("已展开 我的销售漏斗"); `success`/`error` resolutions announce via
  polite live regions; the `arrived`/live update announces a concise summary
  ("回款已更新：深圳蓝海制造 ¥120,000") via a **polite** (not assertive) region,
  and announcements are coalesced (one per batch) so screen readers are not
  flooded. Pause stops live announcements; the new-updates pill is announced once.
- **Target sizes.** Interactive targets meet a minimum touch/click size
  (recommended ≥ 40×40px hit area — note the mockup's rail icons, expand glyphs,
  toggle, and badges-as-buttons must satisfy this even when the visual glyph is
  smaller; UI owes the hit-area padding spec). Strip cards and rail icons must
  remain comfortably clickable in collapsed mode.
- **Reduced-motion + live updates.** Live updates still apply in reduced-motion
  mode; they just snap (no flash fade, no count-up, no slide-insert).

## B7. Cross-Archetype Interaction Patterns

So that list / detail / form / reports inherit the modern, continuous model
proven on the dashboard — and do NOT regress to "click row → full-page reload →
static table." Each archetype states its governing principle; all of them honor
Part A's business contracts and the canonical states (B4).

### Lists (leads, companies, contacts, opportunities, quotes, contracts, payments, activities/notes/tasks)

- **Principle: live, in-place, no full reload.** Filtering, searching, sorting,
  and paging patch the result region in place (skeleton → results crossfade);
  the page chrome and scroll context are preserved.
- **Search/filter:** **debounced** (~250–300ms keystroke debounce — tuning
  decision) with a **live result count** ("共 N 条") that updates as the query
  narrows. Invalid/constrained filters validate before applying and preserve
  input (Part A). Empty result → `empty` state naming the entity (Part A).
- **Selection & bulk actions:** rows support `selected`; a bulk-action bar
  appears (slide-in `motion-fast`) when ≥1 row is selected, showing only
  role-permitted bulk actions; bulk terminal/irreversible actions still route
  through confirmation (Part A). Clear-selection via Esc.
- **Drill-in:** see "Inline drill-in vs full-page nav" below.
- **Loading:** skeleton rows matched to column shape; never a blank table.

### Detail surfaces

- **Principle: drawer/side-panel for in-context detail; full route only for a
  deep, standalone workspace.**
  - **Inline drill-in (default):** opening a record from a list, or expanding a
    dashboard panel, uses an **in-place focus stage (B2) or a right-side drawer**
    — the user keeps their list/overview context behind a scrim and returns with
    Esc/返回 without a navigation/history reload. This is the modern continuity
    model and the default for "open this row to see/triage it."
  - **Full-page navigation (deliberate):** reserved for entering a record's own
    rich workspace (e.g. Opportunity Detail with its pipeline + quote + contract
    + payment + history sub-sections) where the user will stay and do extended
    work, and for first-load deep links / shareable URLs. Full-page nav still
    transitions without a white-flash reload (content region swaps via
    skeleton→content; nav/topbar persist).
  - **Rule of thumb:** triage/glance/quick-edit → drawer/stage (stay in context);
    sustained multi-step work or a shareable destination → full route. Drill-in
    must never feel like a 90s full-page postback.

### Forms (create/edit)

- **Principle: inline, optimistic where safe, never optimistic where business
  rules gate it.**
  - Field validation on blur + on submit; form-level summary for blocking errors
    (Part A). Errors appear instantly and keep input.
  - **Optimistic inline edit** (B4 `optimistic-update`) for low-risk field edits
    and quick toggles (e.g. editing a note, toggling a simple flag) so they feel
    instant, with rollback on failure.
  - **Confirm-then-commit (NOT optimistic)** for every Part A gated action:
    qualify lead, change opportunity stage, Won/Lost close, accept quote, create
    contract, sign contract, record payment, archive, transfer owner, user/role
    changes. These show real `loading` → `success`/`error`, with confirmation
    dialogs (focus-trapped, Part A) where Part A requires them.
  - Drawers used for create/edit keep the parent context behind a scrim and trap
    focus.

### Reports

- **Principle: progressive, drillable, never a static dump.**
  - Report sections load with skeletons; charts animate in once (a single
    `motion-base` reveal — bars grow / line draws — purely as a load affordance,
    suppressed under reduced-motion). No looping chart animation.
  - Drill-in from a report figure to its underlying records uses the inline
    drawer/stage pattern, preserving the report context.
  - `empty`/zero-report state per Part A; unauthorized records excluded silently
    (no count leakage), per Part A and security.

### Shared modern affordances (all archetypes)

- **Toasts:** success (auto-dismiss) and background-error toasts; slide-in
  `ease-decelerate`, out `ease-accelerate`; never used for a blocking error the
  user is actively looking at (that's inline). Toast stack is capped; older
  toasts collapse.
- **Skeleton loading:** the default for any content region awaiting data
  (replaces spinners for content).
- **Empty / first-run:** calm, instructive, role-aware (Part A).
- **Optimistic-then-reconcile:** the default feel for safe edits; confirm-commit
  for gated business actions.

## B8. List Archetype — Interaction Pattern (worked exemplar: 商机 list)

This section instantiates the LIST page archetype's interaction design concretely
for the **商机 (Opportunity) list**, and states it as the reusable pattern for
**all CRM list pages** (线索 / 公司客户 / 联系人 / 商机 / 报价 / 合同 / 回款 /
任务). It is grounded in Part A (business contracts), Part B0–B7 (motion, states,
navigation, archetype principles), and the accepted decisions DEC-UX-NAV-01 /
DEC-UX-MOTION-02 / DEC-UX-LIVE-03. It does not re-litigate those; it uses them.

Nav state on a list page: per B5's either-collapsed-OR-expanded rule, **a list
page renders with the sidebar EXPANDED (`248px`)** — a list is an overview/triage
surface, not the focus view. The collapsed `72px` rail is reserved for the focus
stage (B2). The collapsed-rail hover/focus flyout (DEC-UX-NAV-01) is not relevant
on a list page because the rail is expanded there.

Visual treatments for everything below are owed by `design-system.md` (see the
"UI Handoff" list at the end of Part B and the new owed items in B8.11). Behavior
is owned here; UI renders by state name.

### B8.1 Data scope by role (drives default filter + result count)

Data scope is enforced server-side (per `security/permission-matrix.md`) and the
list UI only ever shows what the role is authorized to see — there is no
client-side "see more" affordance for out-of-scope records, and the result count
("共 N 条") reflects the authorized scope only (no count leakage, per Part A and
security).

| Role | Default authorized scope (商机) | Default filter shown | "共 N 条" counts | Extracted from |
|---|---|---|---|---|
| Sales (销售) | Owned/assigned opportunities only (PM-018/PM-019). | Implicit owner scope = self; a "负责人" facet is present but cannot widen beyond self (selecting another owner returns the authorized subset, typically empty). | Only owned/assigned. | PM-011, PM-018, PM-019 |
| Sales Manager (销售经理) | Team records (the single committed team workspace) (PM-009). | "负责人" facet defaults to "全部（团队）" and can filter down to any team member. | Whole team scope. | PM-009 |
| Administrator (管理员) | Governed records (all CRM records in the workspace) (PM-008). | "负责人" facet defaults to "全部" across the governed workspace. | All governed records. | PM-008 |

Pattern rule (all lists): the **default filter is the role's own scope**, never a
wider scope the role would then be denied on. The result count is computed within
scope. A Sales user never sees a count that hints at records they cannot open.
Out-of-scope rows are simply absent (hidden), not shown as `permission-denied`
rows — `permission-denied` is reserved for an attempt to open a specific
out-of-scope record reached by a direct path (B8.7 / Part A IX-002).

Archived records are excluded by default and only appear behind an explicit
"已归档" filter, scoped per PM-030/031/032/033 (Sales only within owned/assigned).

### B8.2 Filter + search (debounced live count)

- **Search box** (顶部工具栏, CMP-003 Toolbar): free-text over the entity's key
  display fields (商机名 / 客户名). Keystroke-**debounced ~250–300ms**
  (DEC-UX-LIVE-03); each settled query patches the result region in place
  (skeleton → results crossfade, B7 Lists), never a full-page reload, and keeps
  page chrome + scroll context.
- **Live result count:** a persistent `共 N 条` updates as the query/filters
  narrow (count reflects authorized scope, B8.1). While a debounced query is
  in flight the count shows a brief `loading` affordance, then settles; the count
  never flickers per keystroke (it updates on the settled, coalesced query).
- **Filter chips / facets** (商机 exemplar):
  - 阶段 (stage): New Opportunity / Needs Confirmed / Quote / Contract
    Negotiation / Won / Lost (multi-select; extracted from PRD opportunity states).
  - 负责人 (owner): scoped per B8.1.
  - 金额 (expected amount range): min/max numeric range over 预计金额.
  - 预计签约日期 (expected close date): date range.
  - 状态/归档 (archived): 进行中 (default) / 已归档.
  - (Generalization: other lists swap the entity-specific facets — e.g. 线索 uses
    来源/状态; 报价 uses 报价状态/有效期; 合同 uses 合同状态/签约日期; 回款 uses
    回款状态/到期日期; 任务 uses 状态/到期日期/负责人. The facet *mechanics* are
    identical.)
- **Constrained-filter validation:** facets with a constrained format (e.g. the
  金额 range, the date range) validate **before applying** (Part A "filters
  validate before applying"); an invalid range (e.g. min > max, malformed date)
  keeps the input and shows an inline correction message — it does not clear the
  user's other applied filters.
- **Applied-filter display:** active facets render as removable chips in an
  "已筛选:" row beneath the toolbar, each with an `×` to remove that one facet;
  a **`清除全部`** control clears all facets + search and returns to the role's
  default scope (B8.1). Removing/clearing re-runs the debounced query in place.
- **Empty distinctions (two different copies, never merged):**
  - **No-data-yet** (`empty`, authorized scope genuinely has zero records):
    e.g. `还没有商机。从合格线索转化或新建一个商机开始。` + the role-permitted
    next action (新建商机) where the role may create. Calm, instructive (B4
    `empty`). For a role that cannot create here, no create CTA is shown (Part A:
    empty must not suggest unauthorized actions).
  - **No-results-after-filter** (`empty`, data exists but the current
    query/filters match nothing): e.g. `没有符合当前筛选条件的商机。` + a
    **`清除全部筛选`** affordance. This copy must NOT imply the user has no data;
    it points at the filter, not at emptiness.

### B8.3 Columns & sort (商机 exemplar)

Exemplar column set (left→right), with alignment and sort:

| Column | Content | Align | Sortable | Notes |
|---|---|---|---|---|
| (select) | Row checkbox | center | no | Selection (B8.5); header = select-all-in-page. |
| 商机名 | Opportunity name (primary link target) | left | yes (A–Z) | Click = drill-in (B8.4); truncates with title tooltip. |
| 客户 | Related company/customer | left | yes | Links to customer in detail, not from list row. |
| 阶段 | Current pipeline stage badge | left | yes (pipeline order) | Badge, not raw text; stage order, not alphabetical. |
| 负责人 | Owner | left | yes | Hidden/!widened per scope (B8.1). |
| 金额 | 预计金额 (expected amount) | **right** | yes (numeric) | **Tabular numerals**, currency, right-aligned. |
| 预计签约 | Expected close date | right | yes (date) | Tabular; overdue-vs-future styling owed to UI. |
| 更新时间 | Last updated | right | yes (date) | **Default sort: 更新时间 降序** (most recently touched first). |

- **Default sort:** `更新时间` descending — the most operationally relevant order
  for a triage list. (Generalization: every list defaults to `更新时间` desc
  unless the entity has a stronger natural order — e.g. 回款 may default to
  到期日期 asc to surface due-soonest; note such per-entity overrides where they
  exist; otherwise 更新时间 desc.)
- **Sort interaction:** single-column sort; clicking a sortable header toggles
  asc/desc and shows a direction indicator; sorting patches the result region in
  place (skeleton → results crossfade), preserves selection and scroll context,
  and never triggers a full reload (B7 Lists). Sort is announced to AT (polite).
- **Right-aligned tabular numerics:** 金额 and all money/count/date figures use
  tabular numerals and right alignment so columns align digit-to-digit
  (reconciles with `design-system.md` numeric/tabular rule). UI owes the exact
  type token; behavior owes the alignment + tabular requirement.
- **Stage as a badge:** 阶段 renders as a stage badge (one of the six states),
  not free text; Won/Lost read as terminal. Badge color/treatment owed to UI.

### B8.4 Row interaction (drill-in + row quick actions)

- **Row click → DETAIL via the Part B drill-in pattern (NOT a full-page
  postback).** Per B7 "Detail surfaces", opening a 商机 row from the list is a
  **triage/glance** action, so it uses the **right-side drawer** (in-context
  drill-in): the list stays behind a scrim, the drawer slides in (`motion-base`
  `ease-standard`), focus moves into and is **trapped** in the drawer, and
  **Esc / drawer 关闭 returns to the exact list position** (scroll, selection, and
  applied filters preserved — continuity, no reload). This is the default
  "open this row to see/triage it" behavior.
  - The drawer surfaces the record summary + the high-frequency actions (推进阶段,
    编辑, open history). A **"打开完整商机"** link in the drawer escalates to the
    **full Opportunity Detail route** (the deep, standalone workspace with
    pipeline + quote + contract + payment + history) for sustained multi-step work
    or a shareable URL — that is the deliberate full-page case (B7 rule of thumb).
  - Deep links / first-load to a specific opportunity open the full route
    directly (B7), still without a white-flash reload.
- **Row-level quick actions** (shown on hover/focus and in the row's Action Menu,
  CMP-019; never hover-only — keyboard-reachable per B6):
  - **编辑 (edit)** of low-risk inline fields → may be **optimistic**
    (`optimistic-update`, B4) where Part A does not gate it (e.g. editing a free
    note field): apply immediately with a pending affordance, reconcile, roll back
    + inline `error` on failure.
  - **推进阶段 (advance stage)** → **confirm-then-commit, NOT optimistic.** Stage
    change is a Part A gated business transition (IX-006, PM-018, PRD transition
    table): it shows real `loading` → `success`/`error`, validates the transition
    server-side, and surfaces the required-data / forbidden-transition reason on
    failure (e.g. `New Opportunity → Needs Confirmed` requires related
    customer/contact; `Contract Negotiation → Won` requires a Signed contract,
    DEC-017). Terminal **Won / Lost** closures are always confirm-then-commit with
    the Part A confirmation dialog (Won shows the Signed-contract gate; Lost
    prompts the required lost reason).
- **Quick-action gating rule (all lists):** a quick action is **optimistic only
  if Part A does not require a server precondition/confirmation**; every gated /
  terminal / irreversible action (qualify, stage change, Won/Lost, accept quote,
  contract create/sign, payment record, archive, owner transfer) is
  confirm-then-commit (B4 `optimistic-update` forbidden set, B7 Forms). An action
  the role/record-state does not permit renders `disabled` with a safe
  explanation (e.g. 推进阶段 disabled on a terminal Won/Lost opportunity:
  `商机已结束，无法推进阶段。`), not hidden, when the reason is safe to reveal.

### B8.5 Selection + bulk actions

- **Multi-select:** per-row checkbox sets the row to `selected` (B4); the header
  checkbox is **select-all-in-page** (selects the currently loaded/visible page of
  rows only). Selecting ≥1 row slides in the **bulk-action bar** (`motion-fast`).
- **Select-all-in-page vs all-matching:** when select-all-in-page is active and
  more rows match the filter than are on the page, an inline affordance offers
  **"选择全部匹配的 N 条"** (all-matching across the filtered, in-scope result
  set). This distinction is explicit so a bulk action's target count is never
  ambiguous — the bulk bar always states the exact target ("已选择 12 条" vs
  "已选择全部匹配的 247 条"). All-matching is still bounded by the role's
  authorized scope (B8.1); it can never include out-of-scope records.
- **Bulk-action bar:** shows only **role-permitted** bulk actions (gated per
  `permission-matrix.md`); shows the live selected count; has a **清除选择**
  (also via **Esc**, B6). Exemplar 商机 bulk actions and their gating:
  - **批量转移负责人 (bulk owner transfer):** Sales Manager / Administrator only
    (PM-014); **denied for Sales** (PM-015) so it does not appear in a Sales
    user's bulk bar. Confirm-then-commit (IX-022); confirmation states the target
    count + that open tasks/follow-ups transfer with the owner (BR rule).
  - **批量归档 (bulk archive):** Sales Manager (team) / Administrator (governed)
    only (PM-026/027); **denied for Sales** (PM-028). Confirm-then-commit;
    per-record obligation check still applies (archive blocked by active
    obligations lists them per Part A) — a bulk archive reports per-row success
    vs blocked (partial-result, like CSV import), never silently skips.
  - **批量导出 (bulk export selected):** Sales Manager / Administrator only
    (PM-037/038); **denied for Sales** (PM-039). Explicit export confirmation;
    logged.
  - (No bulk 推进阶段 in committed scope — stage transitions are per-record gated
    with per-record required data; bulk stage advance is intentionally absent to
    avoid bypassing per-transition validation. Note as a deliberate omission, not
    a gap.)
- **Bulk terminal/irreversible actions stay confirm-then-commit** (Part A): the
  confirmation dialog (focus-trapped, B6) states target record count, business
  effect, and irreversibility; bulk results that partially fail show a per-row
  outcome summary (succeeded / blocked-with-reason) and never imply full success.

### B8.6 Pagination strategy — DECISION: paginated, not infinite-scroll

**Recommendation (adopted for all CRM lists): classic pagination (page-size +
page navigation), NOT infinite scroll.**

Rationale:

- **Count legibility:** CRM lists are operational/reporting surfaces where
  "共 N 条" and "第 X / Y 页" are meaningful and must stay legible; infinite scroll
  makes total counts and "where am I" ambiguous.
- **Bulk-select correctness:** the select-all-in-page vs select-all-matching
  distinction (B8.5) is only crisp with a defined page; infinite scroll blurs
  "this page" and makes bulk-target counts ambiguous and error-prone for
  destructive actions (transfer/archive). Correctness of destructive bulk actions
  outranks the marginal scroll convenience.
- **Layout stability + live updates:** pagination bounds the live-patch surface
  (B3) to one page, keeping top-insert behavior and the `arrived` highlight
  predictable; infinite scroll fights live top-inserts and reflow stability.
- **Determinism for QA/keyboard:** a bounded page is easier to keyboard-navigate
  and to verify against acceptance.

Behavior:

- Default page size (e.g. 25; dense manager view may offer 50) with a page-size
  selector; page navigation (上一页 / 下一页 / 第 X 页) and the `共 N 条` /
  `第 X / Y 页` readout.
- Changing page patches the result region in place (skeleton rows → results
  crossfade), preserves applied filters + sort, and resets scroll to the top of
  the result region only (not the page chrome). Page change is **not** a full
  reload.
- Selection scope is explicit across pages: changing page does **not** silently
  extend an in-page selection; if the user chose "全部匹配的 N 条", that
  all-matching selection persists across page changes and the bulk bar keeps
  showing the all-matching count.
- **"加载更多" is not used** for these lists; if a future high-volume surface
  needs it, it must preserve count legibility and the bulk-select distinction
  (raise back to UX, do not silently switch).

### B8.7 States (exact canonical names)

Uses the canonical state names verbatim (B4 / `screen-state-spec.md`):

- **`loading`:** **skeleton rows** matched to the column shape (never a blank
  table, never a content spinner) (B4, B7 Lists). Skeleton appears instantly,
  shimmer loops `motion-base`-paced; skeleton → results **crossfade**
  `motion-fast` `ease-decelerate`. Layout is preserved so rows do not jump when
  data lands. The `共 N 条` count shows its own brief `loading` while a query is
  in flight.
- **`empty`:** two distinct copies per B8.2 — **no-data-yet** (name the entity +
  role-permitted next action) vs **no-results-after-filter** (point at the filter
  + offer 清除全部筛选). Single crossfade in; never an error tone.
- **`error`:** a **load/filter failure** shows an **inline** error in the result
  region (not a toast, since the user is looking at it), keeps the applied
  filters/search input, and offers **重试** → returns to `loading`
  (Part A: search/filter errors preserve filter inputs). A background/async
  failure (e.g. a bulk action kicked off then failing out of view) uses a toast.
  Errors appear instantly (no fade delay before readable).
- **`permission-denied`:** rows the role may not see are **absent** (hidden,
  B8.1), not rendered as denied rows. `permission-denied` is reached only when the
  user follows a **direct path to a specific out-of-scope record** (e.g. a stale
  deep link / row id outside scope, Part A IX-002): the detail surface shows the
  explicit denial with **no restricted data** and a **return-to-list** control;
  focus moves to that control (B6). Actions the role can't perform are `disabled`
  with a safe reason (B8.4) — not `permission-denied`.
- **`disabled`:** row-level quick actions or bulk actions not currently allowed
  (terminal-state stage advance; a Sales user's transfer/archive/export bulk
  actions are simply **not present**, since they are role-denied, rather than
  shown disabled — disabled is for *missing-prerequisite* on an otherwise
  available action, hidden is for *role-denied*). Disabled controls explain the
  missing requirement where safe (B4 `disabled`).
- **`optimistic-update` + rollback:** only the safe inline edits (B8.4) apply
  optimistically — value updates `motion-instant` with a pending affordance; on
  success it confirms with no second flash; on failure it **rolls back** to the
  prior value with a `motion-fast` revert + inline `error` + 重试 (feedback must
  never imply a failed change persisted — operating-model).
- **`selected` / `focused` / `hover`:** per B4 — `selected` (persistent,
  checkbox-driven, drives the bulk bar), `focused` (visible focus-visible ring,
  keyboard), `hover` (row tint, no movement; reveals quick actions but is never
  the *only* way to reach them). Visually distinct from each other and from
  `disabled` / `permission-denied` (UI owes the distinct tokens).
- **`success`:** an in-context row mutation (e.g. an inline edit confirmed, a
  stage advance succeeded and the row's 阶段 badge flips) shows lightweight inline
  success + refreshes the affected visible state (DEC-020); a bulk action whose
  result isn't fully visible uses an auto-dismiss toast (~3–4s) summarizing the
  outcome.

### B8.8 Live updates on a list (实时更新 / 暂停)

The Part B3 live model applies to the list result region (lists inherit it per
B3 / B7 Lists):

- **In-place row patch, no reload, no scroll jump:** when a visible row's data
  changes (e.g. another user advances a stage, a 金额 changes), the cell patches
  in place and the row shows the single transient **`arrived` highlight**
  (~1.6s hold, fade `motion-base`) — redundant with the already-correct static
  value (no motion-only signaling, B6). Scroll, focus, selection, applied filters,
  and any in-progress inline edit are preserved (B3).
- **New-row arrival affordance:** a newly-matching row (within the current page's
  sort window — e.g. a just-created opportunity at the top of `更新时间` desc)
  enters at the top with a **height-stable insert** (opacity/translateY only;
  surrounding rows do not reflow-jump). If many arrive at once they **coalesce**
  (DEC-UX-LIVE-03 ~800ms–1s window) into one batch — no flicker storm — and the
  `共 N 条` / `更新于 HH:MM` meta reflects the applied batch, not each event.
  A row that changes twice within a window flashes once with the final value.
- **No content the user is reading moves:** live updates never reorder beyond the
  defined top-insert and never steal focus or move a row the user is mid-edit on.
- **Pause (`暂停`):** the per-user, per-surface live toggle stops auto-applying
  updates to this list; the live-dot goes static (paused). Incoming changes are
  **buffered, not dropped**, and a non-blocking pill **`有 N 条新更新 · 点击刷新`**
  appears; activating it (or resuming `实时更新`) applies the buffered batch with
  the normal `arrived` highlights and returns the meta timestamp to live. Pause
  is a **view preference only** — it never changes server state (B3).
- **Value count-up (DEC-UX-MOTION-02):** a numeric cell/aggregate that changes
  (e.g. a list-footer total, or a 金额 on a live update) may use the conservative
  count-up (≤ `motion-base`, only on meaningful deltas); snaps under
  reduced-motion (B6).
- **Reduced-motion + reconnect:** under `prefers-reduced-motion` live updates
  still apply but snap (no flash fade, no count-up, no slide-insert) (B6).
  SSE/live channel loss degrades to the existing `offline/reconnect` handling —
  keep last-known rows (no blank-out), live-dot shows reconnecting, and a quiet
  catch-up batch applies on reconnect (B3).

### B8.9 Accessibility & keyboard (binding, per B6)

- **Landmark / focus order:** skip-link → primary nav (expanded sidebar) → topbar
  (search, quick links, user menu) → list toolbar (search, filters, sort,
  density, 实时更新/暂停 toggle, 新建) → result table → pagination (B5 landmark
  order; the table is within main content).
- **Row focus order:** the table is keyboard-navigable (Part A: lists/tables
  support keyboard navigation with meaningful row labels). Tab/arrow order follows
  visual reading order (header controls → rows top-to-bottom → row's interactive
  cells: select checkbox → primary link → quick actions). Each row has a
  meaningful accessible label (e.g. `商机：深圳蓝海制造 · 阶段 报价 · 负责人 张伟`).
- **Keyboard select:** the row checkbox is focusable and toggles `selected` with
  Space; the header select-all and the "选择全部匹配" affordance are keyboard
  operable; **Esc clears the selection** (and closes the bulk bar) (B6, B7 Lists).
- **Drill-in drawer:** opening a row moves focus **into** the drawer and **traps**
  it; **Esc closes the drawer** and **restores focus to the originating row**
  (focus restoration mandatory, B5/B6); the list position is preserved.
- **focus-visible:** every interactive element (row link, checkbox, header sort,
  filter chips, quick actions, pagination, live toggle, "有 N 条新更新" pill)
  shows the `focused` ring on keyboard focus, contrast-compliant on light and
  tinted/selected row backgrounds (ring token owed to UI).
- **Live-region announcements (polite, coalesced):** sort changes, a settled
  filter/search result ("共 12 条"), and live `arrived` updates announce concisely
  via a **polite** region, one announcement per coalesced batch (not flooded);
  pause stops live announcements and the new-updates pill is announced once
  (B6). `success`/`error` resolutions announce politely.
- **Reduced-motion:** all list motion (skeleton shimmer, crossfade,
  `arrived` flash, count-up, bulk-bar slide, drawer travel) snaps/static under
  `prefers-reduced-motion` with **no information lost** (B6).
- **Target sizes:** checkboxes, sort headers, quick-action glyphs, the live
  toggle, and the new-updates pill meet the B6 ≥40×40px hit-area minimum even
  where the visual glyph is smaller (hit-area padding owed to UI).

### B8.10 List-archetype acceptance hooks (for QA / acceptance-matrix)

These map the list behavior to existing acceptance IDs so QA can verify (no new
scope; references existing matrix items):

- Authorized scope + count by role: ACC-002, ACC-015 (PM-008/009/011/018/019).
- Search/filter (debounced, live count, invalid-filter preserve-input): ACC-015
  (IX-023).
- Open detail (drawer drill-in) + permission-denied on direct out-of-scope path:
  ACC-002, ACC-015 (IX-002, IX-024).
- Row stage advance (confirm-then-commit, forbidden-transition reason): ACC-008,
  ACC-014 (IX-006).
- Won/Lost from list quick action (terminal, gated): ACC-013 (IX-007, IX-008).
- Bulk owner transfer (Manager/Admin allow, Sales deny): ACC-002, ACC-014
  (IX-022, PM-014/015).
- Bulk archive (blocked-by-obligation per row, Sales deny): ACC-002, ACC-014
  (IX-019, PM-026/027/028).
- Bulk export (Manager/Admin allow, Sales deny): ACC-020 (PM-037/038/039).
- Empty (no-data vs no-results) + role-aware next action: ACC-015 (and ACC-018/
  ACC-023 for report-style empties).
- Live update in place / pause-buffer: inherits B3; verify no scroll jump, no
  focus steal, buffered-not-dropped on pause.

### B8.11 UI Handoff additions for the list archetype

Adds to the Part B "UI Handoff — Visual Treatments Owed" list (reconcile in
`design-system.md`; these are list-specific and must agree with CMP-003 Toolbar,
CMP-004 Data Table, CMP-018 Empty State, CMP-019 Action Menu in `component-spec.md`):

- **Skeleton rows** matched to the 商机 column shape (and the generalized list
  column shapes).
- **Stage badge** treatment for the six opportunity states (and per-entity status
  badges), with terminal Won/Lost reading.
- **Right-aligned tabular-numeric** column styling for 金额 / dates (reconcile
  with the design-system numeric/tabular rule).
- **Applied-filter chips** + `清除全部` styling; **constrained-filter inline
  error** styling.
- **Bulk-action bar** styling and the **select-all-in-page vs 选择全部匹配** count
  affordance.
- **Pagination control** styling (`共 N 条` / `第 X / Y 页` / page-size selector).
- **Row `selected` / `hover` / `focused` distinction** on a dense table row, plus
  the row-level **`arrived`** highlight on live update (must agree with the
  dashboard `arrived` source of truth).
- **Empty-state copies** distinction (no-data vs no-results) visual tone.
- Per B6, list checkboxes / sort headers / quick-action glyphs **hit-area
  padding** to the ≥40×40px minimum.

UI must NOT change these behaviors; a visual constraint forcing a behavior change
is raised back to UX (G4b), not diverged silently.

## B-micro. Micro-Interactions (inventory + timing)

Only these micro-interactions exist; anything not listed is out of scope (avoid
gratuitous motion). All snap/disable under reduced-motion.

| Micro-interaction | Where | Timing / curve | Notes |
|---|---|---|---|
| Button press | All buttons | `motion-instant` `ease-emphasis` | Scale ~0.98 + shadow drop; tactile, not bouncy. |
| Card hover-lift | Grid panels, KPI cards, strip cards | `motion-fast` `ease-standard` | translateY -2px + shadow step; off on touch/reduced-motion. |
| Row hover tint | List rows, focus-stage table | `motion-fast` | The `hovered` treatment; no movement, tint only. |
| Toggle (实时更新/暂停 segment, switches) | Live toggle, settings | `motion-instant` | Thumb/segment slide; state is also text-labeled (not color-only). |
| Expand/collapse chevron | Panels, sections, accordions | `motion-fast` | Rotate 90–180°; pairs with content reveal. |
| Live-dot pulse | Live surfaces | gentle ~2s loop, `ease-standard` | Pulsing = "live"; paused = static dot (UI owes paused token). Redundant with text status. Off on reduced-motion (stays visible, static). |
| `arrived` highlight flash | Changed rows/values on live update | hold ~1.6s, fade `motion-base` | The only "just changed" signal; redundant with static value (B3/B6). |
| Value count-up | KPI/total on refresh | ≤ `motion-base`, `ease-decelerate` | ACCEPTED (DEC-UX-MOTION-02, 2026-06-06): ON but conservative — only on meaningful metric changes, not every tick; snaps on reduced-motion. |
| Skeleton shimmer | Loading regions | loop, `motion-base`-paced | Static placeholder under reduced-motion. |
| Toast in/out | Success/bg-error | in `motion-fast` `ease-decelerate`, out `ease-accelerate` | Auto-dismiss success ~3–4s. |
| Stage transition | Card→focus | Dedicated hero timing per DEC-UX-HEROTIME-01: ~450ms in / ~310ms out, `ease-standard` | The one hero motion (B2). In-focus selector switching still uses `motion-base` (220ms). |

## Decisions — Resolved / Flagged

- **DEC-UX-NAV-01 — Collapsed-rail hover/focus-expand. ACCEPTED by user 2026-06-06: YES.**
  The collapsed `72px` rail hover/focus-expands to a temporary icon+text flyout
  overlay (does not reflow the stage or change the 72px track, so it does not break
  the either-collapsed-OR-expanded rule). Rejected alternative: strict icon-only +
  tooltips. This is now binding for all archetypes; the flyout's visual treatment is
  owed by `design-system.md` (Rail flyout).
- **DEC-UX-MOTION-02 — KPI value count-up on refresh. ACCEPTED by user 2026-06-06: ON (conservative).**
  KPI/total values count-up only on meaningful metric deltas, ≤ `motion-base`
  (~220ms), `ease-decelerate`, and snap under `prefers-reduced-motion`. Binding.
- **DEC-UX-LIVE-03 — Coalescing/debounce windows.** Defaulted 2026-06-06 to the
  recommended ~800ms–1s live-update coalescing and ~250–300ms search-keystroke
  debounce; final values remain a tuning decision to validate against real latency
  during execution. Owner: Architecture (tuning only; the behavior is fixed).
- **DEC-UX-LIVE-04 — Dashboard live transition mechanism. ACCEPTED by release owner 2026-06-10.**
  Source: release-owner direction recorded in `planning/blockers.md`
  `BLK-UIUX-G12-009`. Because true server-push/SSE remains a parked future path
  and C1 forbids backend/API/shared changes in this spot-fix, the current
  dashboard B3 live layer is implemented via client-side polling of the same
  existing GET endpoints already used by dashboard refresh. Polling uses the B3
  patch-in-place model, the DEC-UX-LIVE-03 coalescing window, change-driven
  `arrived` highlights, and the functional `实时更新` / `暂停` pause buffer. This
  replaces static/dead live UI and does not alter the future SSE interaction
  contract.
- **DEC-UX-FOCUSRAIL-01 — Focus rail as persistent selector. ACCEPTED by release owner 2026-06-10: Option A.**
  Source: release-owner direction recorded in `planning/blockers.md`
  `BLK-UIUX-G12-007`. The previously locked B2 model used a "7 non-active cards"
  rail, so switching focus caused the old active card to pop into the rail and the
  clicked card to leave it. That behavior was spec-faithful but fatiguing in live
  use. The binding revision is Option A: the focus rail is a persistent selector
  that lists the full dashboard card set for the current workspace in fixed order
  (locked manager dashboard = all 8 cards, including the current one; role-scoped
  variants list their full authorized set). The current item shows the selected
  treatment using existing locked tokens only and carries `aria-current="true"`.
  Clicking/Enter/Space on any rail item changes only the left stage via the
  existing calm content crossfade; the rail does not add/remove/reorder items.
  BLK-UIUX-G12-006 hero enter/exit remains unchanged.
- **DEC-UX-FOCUSEXIT-01 — Focus header has one visible exit control. ACCEPTED by release owner 2026-06-10.**
  Source: release-owner direction recorded in `planning/blockers.md`
  `BLK-UIUX-G12-008`. The locked focus mockup showed a `返回` button plus a
  visible `Esc 退出`/`Esc 返回` hint chip; implementation further duplicated the
  word "返回". The binding revision supersedes that mockup chip: the focus stage
  header shows exactly one visible exit control, the `返回` button with its
  back-chevron icon. The `Esc` key remains a global focus-exit shortcut for B6
  keyboard operability, but the stage header no longer displays an `Esc` hint
  chip or other visible `Esc` text. Data-scope badges, BLK-UIUX-G12-006 hero
  enter/exit motion, and BLK-UIUX-G12-007 selector rail behavior are unchanged.
- **DEC-UX-HEROTIME-01 — Card→Focus hero timing is slower than the B1 slow token. ACCEPTED by release owner 2026-06-10.**
  Source: release-owner direction recorded in `planning/blockers.md`
  `BLK-UIUX-G12-012`. The B2 hero transition behavior remains approved, but its
  original `motion-slow` 320ms entry felt too fast in review. This binding
  revision raises only the Card→Focus hero timing to a dedicated enter duration
  of ~450ms and a proportional reverse duration of ~310ms. The B1 base scale is
  not globally changed: `motion-base` remains 220ms for in-focus selector
  switching and other standard transitions; reduced-motion remains the
  `motion-instant` opacity-only snap; no other motion inherits the hero timing.
- **SSE "已回款" push** remains a **parked future feature** (not new P0/P1 scope).
  Part B only defines how it would surface (B3) so it has a non-jarring landing;
  it must not be implemented as part of current scope without a formal scope
  change. Owner: PM/user.

## UI Handoff — Visual Treatments Owed (reconcile in `design-system.md`)

Part B specifies behavior and references these treatments by intent; UI Design
must supply the concrete visual values/tokens (consistent with the LOCKED
palette) so the approved mockup and this spec agree:

- **Motion tokens** (`motion-instant/fast/base/slow`, `ease-standard/decelerate/
  accelerate/emphasis`): mirror exact names + values from B1 into the design
  system's motion section.
- **`arrived` / changed-row highlight**: accent-bar + tint values + fade
  (the mockup `arrived` class is the source of truth).
- **Skeleton + shimmer**: skeleton fill color and shimmer gradient.
- **Live-dot states**: live (pulsing) vs paused (static/neutral) vs reconnecting
  dot treatments.
- **Hover / selected / focused tokens**: `hovered` row tint, `selected`
  tint+accent, focus-visible ring (color + width + offset, contrast-compliant on
  light and tinted backgrounds).
- **Scrim**: the focus-stage workspace scrim (`rgba(15,23,42,.06)` in mockup) —
  confirm token.
- **Strip-card spec**: `92px` compact card, single-key-value layout (from mockup
  `sideCard`).
- **Rail flyout**: the hover/focus label-flyout overlay styling (if DEC-UX-NAV-01
  = yes).
- **Toast**: success vs error toast styling and stack behavior.
- **Disabled vs permission-denied vs enabled** visual distinction (must be
  visually unambiguous and contrast-compliant).
- **New-updates pill** (`有 N 条新更新 · 点击刷新`) styling.
- **Hit-area padding** so small glyphs (rail icons, expand, toggle,
  badge-buttons) meet the B6 target-size minimum.

UI must NOT change the behaviors specified in Part B; where a visual constraint
forces a behavior change, raise it back to UX (G4b) rather than diverging
silently (per G4c→G4b dependency).
