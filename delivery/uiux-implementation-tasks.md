# UI/UX Implementation Tasks

Status: **UI/UX G11 complete — implementation, QA, and integration evidence
ready for Claude G12 audit**.
Claude passed the UI/UX G8 handoff audit and DEC-UIUX-A5-001 token checkpoint on
2026-06-07. UIUX-001..014 are implemented and verified by `npm run build` plus
`npm run test:e2e` (45/45 passed, 0 skips) from `frontend/`.

Read first:

- `delivery/uiux-implementation-delivery-plan.md`
- `docs/ux-ui/requirements/uiux-implementation.requirements.md`
- `docs/ux-ui/design-system.md`
- `docs/ux-ui/mockups/*.png`
- `docs/ux-ui/mockups/_src/*.html`
- `frontend/src/i18n/labels.ts`

Global hard constraints on every task: C1 frontend/design only; C2 no downgrade
of P0/P1 or G12 security fixes; C3 preserve zh-CN and enum/role comparison
values; C4 display real enums via `labels.ts`; C5 e2e green with no skip or
assertion reduction; C6 match locked mockups and design-system.

Every task below binds to ACC-018 and ACC-023 because this follow-on change
completes the committed CAP-009 design-realization layer, and because the shared
UI system must support the team overview/basic reports surfaces without visual
or functional downgrade.

---

## UIUX-001 — React Design-System Foundation

1. **Status:** Done — token checkpoint passed; React CSS variables, motion
   styles, and UI primitives are wired into the frontend; `npm run build` green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** global frontend shell and all 14 nav pages.
4. **Acceptance binding:** ACC-018, ACC-023; A1; C1-C6.
5. **Objective:** Extract `docs/ux-ui/design-system.md` into project-local CSS
   variables and UI primitives, using locked token values plus the approved
   DEC-UIUX-A5-001 text-only `*-ink` contrast tokens and existing frontend
   dependencies.
6. **Implementation scope:** `frontend/src/styles.css`,
   `frontend/src/styles/design-system.css`, `frontend/src/styles/motion.css`,
   `frontend/src/components/ui/*`.
7. **Do not change:** backend, API clients' request/response contracts, data
   models, enum comparison values, roles, authorization decisions.
8. **Design-is-implemented acceptance:** CSS variables include the locked color,
   typography, spacing, radius, shadow, state, and motion token values; only
   DEC-UIUX-A5-001 may add text-only `*-ink` colors, and those may not be used
   for fill/background/button/badge background/border/icon/graph roles;
   primitives exist for card/panel, button, badge, table, toolbar, form field,
   skeleton, empty/error/permission-denied, pagination, drawer, metric card,
   chart shells, and toast/live affordances.
9. **Verification:** static token spot-check against `design-system.md` and
   `delivery/uiux-token-exception-a5-2026-06-07.md`; `npm run build` remains
   green after foundation wiring.

## UIUX-002 — App Shell, Navigation, Topbar, Dashboard / Focus Stage

1. **Status:** Done — expanded shell, role-filtered nav, topbar, dashboard,
   live banner, and focus-stage implementation landed; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** `Shell`, `Nav`, `WorkOverview`, dashboard/focus page
   types.
4. **Acceptance binding:** ACC-018, ACC-023; A1, A2, A3, A5, A6, A7; C1-C6.
5. **Objective:** Replace the old dark/sidebar shell with the locked 248px
   expanded sidebar, 64px topbar, dashboard grids, KPI cards, live banner, and
   focus-stage visual model.
6. **Implementation scope:** `frontend/src/app/Shell.tsx`,
   `frontend/src/app/Nav.tsx`, `frontend/src/pages/WorkOverview.tsx`, reusable
   dashboard/focus primitives.
7. **Required behavior:** nav remains role-filtered as today; focus mode uses the
   72px rail and temporary hover/focus flyout only as specified; topbar keeps
   user identity and logout; no routing or auth behavior change.
8. **Design-is-implemented acceptance:** dashboard manager/sales variants match
   `dashboard-v7-sales.png`, `dashboard-v7-manager.png`, and
   `dashboard-v7-manager-focus.png` in shell geometry, token use, KPI/grid
   structure, card elevation, icon-led panels, topbar, and zh-CN copy.
9. **Verification:** browser screenshots for 工作台 at desktop width; keyboard
   focus through nav/topbar/main; build green.

## UIUX-003 — List Archetype And Opportunity List Variants

1. **Status:** Done — opportunity list archetype landed with search, stage
   filter, result count, selected rows, pagination structure, and Sales-hidden
   bulk affordances; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** all entity lists, with 商机 list as the exemplar.
4. **Acceptance binding:** ACC-018, ACC-023; A1, A2, A3, A4, A5, A6, A7; C1-C6.
5. **Objective:** Implement the reusable list archetype and apply it first to
   `OpportunityList` for manager/admin and sales variants.
6. **Implementation scope:** list primitives, `OpportunityList.tsx`, existing
   `listOpportunities` usage, selectors used by list e2e tests.
7. **Required behavior:** search/filter/pagination structure, live result count,
   selected/hover/focused rows, bulk bar, role-gated bulk actions, no infinite
   scroll, real six-stage labels from `labels.ts`.
8. **Role gates:** Sales must not see bulk owner transfer or bulk archive; manager
   and admin may see permitted bulk affordances; no bulk stage advance is added.
9. **Design-is-implemented acceptance:** `list-opportunities.png` and
   `list-opportunities-sales.png` page types are represented with the expanded
   shell, page header, filters/chips, table/list body, selected/hover/focused
   states, pagination, and safe bulk semantics.
10. **Verification:** existing opportunity e2e assertions remain; added UI checks
    verify sales hidden bulk actions and real stage display labels.

## UIUX-004 — Detail And Form Archetypes, Opportunity Detail/Form

1. **Status:** Done — opportunity detail/form exemplar landed with terminal
   read-only treatment and Sales owner self-lock on create form; build and e2e
   green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** detail/form page types, `OpportunityDetail`,
   opportunity create/edit form, `StageStepper`, close dialogs.
4. **Acceptance binding:** ACC-018, ACC-023; A1, A2, A3, A4, A5, A6; C1-C6.
5. **Objective:** Implement locked detail/form primitives and apply them to the
   opportunity detail/form exemplar.
6. **Implementation scope:** detail header, stage stepper, relationship panels,
   timeline panel styling, form layout, validation state, save/error/success
   states.
7. **Required behavior:** terminal Won/Lost records are read-only; stage changes
   are linear and server-confirmed; Won remains gated by Signed contract; Lost
   requires reason; new/edit forms exclude terminal stages; Sales owner is
   prefilled/locked to self where the existing role context allows.
8. **Design-is-implemented acceptance:** `detail-opportunity.png` and
   `form-opportunity.png` are represented with locked tokens, cards, badges,
   stepper, field validation, disabled explanations, and zh-CN copy; no business
   action becomes optimistic when it is business-gated.
9. **Verification:** existing opportunity stage/close e2e remains green; added
   checks verify terminal detail read-only affordances and form stage options
   exclude Won/Lost.

## UIUX-005 — Apply CRUD Archetypes To 8 Entity Areas

1. **Status:** Done — eight entity list/detail/form surfaces share the new
   archetype classes, selected states, status pills, and table wrapping; build
   and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** 线索, 公司/客户, 联系人, 商机, 报价, 合同, 回款,
   任务.
4. **Acceptance binding:** ACC-018, ACC-023; A1, A2, A3, A4, A5, A6; C1-C6.
5. **Objective:** Apply the reusable list/detail/form visual system to the eight
   CRUD entity areas without changing their existing API calls or business rules.
6. **Implementation scope:** `LeadList/Detail`, `AccountList/Detail`,
   `ContactList`, `QuoteList/Detail`, `ContractList/Detail`,
   `PaymentList/Detail`, `TaskList`, shared dialogs/panels.
7. **Required behavior:** preserve existing create/edit/transition/archive/
   history behavior; use real label maps for statuses; keep authorized data scope
   and safe denial behavior.
8. **Design-is-implemented acceptance:** each entity uses the same list/detail/
   form shell, table/card density, badges, money/date alignment, state treatments,
   and action hierarchy as the locked page types; no page remains in the old
   plain/table-only visual style.
9. **Verification:** all existing entity e2e specs continue to pass with updated
   selectors only where necessary.

## UIUX-006 — Reports And Manager Overview CAP-009 Realization

1. **Status:** Done — manager overview and basic reports use real API metrics
   for KPI cards, pipeline funnel, breakdown bars, and wrapped data tables;
   access-gating e2e assertions are green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** `ManagerOverview`, `BasicReports`, report charts and
   metrics.
4. **Acceptance binding:** ACC-018, ACC-023 primary; A1, A2, A3, A5, A6, A7; C1-C6.
5. **Objective:** Realize the committed team overview and basic reports design
   using KPI cards, funnel/distribution visuals, grouped breakdowns, and empty/
   loading/error states.
6. **Implementation scope:** report primitives and existing
   `frontend/src/api/reports.ts` data fields only.
7. **Required behavior:** metrics align to `OverviewMetrics`; pipeline and
   breakdowns use real enum labels; manager scope remains team and admin scope
   remains governed/all where existing APIs authorize; Sales remains denied from
   reports via existing nav/API behavior.
8. **Design-is-implemented acceptance:** `reports-team.png` is represented:
   overview metric cards, pipeline by six stages, lead/opportunity/quote/
   contract/payment breakdowns, money as tabular `¥`, date/range controls, export
   affordance visual only where existing behavior supports it.
9. **Verification:** `overview.spec.ts` and `reports.spec.ts` remain green; added
   visual/DOM checks verify chart sections are not replaced by a generic table
   only.

## UIUX-007 — Special Page Types: Admin, Reminders, Import/Export, Operation Log

1. **Status:** Done — admin, reminders, import/export, and operation-log
   surfaces use the locked special-page structure; last-admin options are visibly
   disabled, reminders keep source semantics, import/export tables are wrapped,
   and operation logs are read-only; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** `UserManagement`, `ReminderCenter`,
   `ImportExportPage`, `OperationLogs`.
4. **Acceptance binding:** ACC-018, ACC-023; A1, A2, A3, A4, A5, A6; C1-C6.
5. **Objective:** Apply locked special-page designs and preserve each page's
   security and data semantics.
6. **Implementation scope:** admin table, last-admin disabled treatment, reminder
   rows/groups, import/export result cards and row-error table, read-only oplog
   timeline/table.
7. **Required behavior:** user/role and operation-log pages remain admin-only;
   last active administrator guard is visibly disabled; reminders do not invent
   read/unread semantics; import/export requires confirmation for export and uses
   real ImportRun/ExportRun fields; operation logs use safe summaries and remain
   read-only.
8. **Design-is-implemented acceptance:** `admin-users.png`,
   `reminders-center.png`, `import-export.png`, and `operation-log.png` page
   types are represented with locked shell, tables/cards, badges, state coverage,
   and zh-CN copy.
9. **Verification:** `user-admin.spec.ts`, `reminders.spec.ts`,
   `import.spec.ts`, `export.spec.ts`, and `oplog.spec.ts` remain green; added
   checks cover no read/unread reminder switch and no log edit/delete controls.

## UIUX-008 — Canonical State Layer (A3)

1. **Status:** Done — canonical loading/empty/error/disabled/selected/focused/
   permission/read-only treatments are wired through shared CSS and
   representative pages; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** all shared primitives and all 14 nav pages.
4. **Acceptance binding:** ACC-018, ACC-023; A3; C1-C6.
5. **Objective:** Implement reusable state treatments for `loading`, `empty`,
   `error`, `disabled`, `selected`, `focused`, `hover`, `permission-denied`,
   `optimistic-update`, and `success`.
6. **Implementation scope:** UI primitives, page-level state wiring, aria-live
   regions, skeletons, row selection, disabled explanations, permission-denied
   panels, success/error feedback.
7. **Required behavior:** business-gated actions are not optimistic; permission
   denial leaks no restricted data; disabled states explain safe prerequisites;
   loading is skeleton-based and layout-stable.
8. **Design-is-implemented acceptance:** every page type has visible or
   testable coverage for the canonical states relevant to it, using the static
   visual treatments from `design-system.md` §8 and state rules from
   `screen-state-spec.md`.
9. **Verification:** component/page tests or e2e checks cover representative
   loading/empty/error/disabled/permission-denied/selected/focused states.

## UIUX-009 — Role And Permission Gate Review (A4)

1. **Status:** Done — A4 affordances implemented for Sales bulk hiding,
   terminal read-only details, terminal-free create forms, Sales owner self-lock,
   admin-only nav, last-admin disabled options, and report scope assertions;
   e2e assertions are green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** Opportunity list/detail/form, admin pages,
   operation logs, import/export, reports, nav.
4. **Acceptance binding:** ACC-018, ACC-023; A4; C1-C6.
5. **Objective:** Verify and adjust frontend affordances so they mirror the
   already-reviewed authorization rules without becoming the source of security.
6. **Implementation scope:** nav visibility, action visibility/disabled states,
   owner fields, terminal record actions, last-admin affordances.
7. **Required gates:** Sales hides bulk owner transfer/bulk archive; terminal
   opportunity details are read-only; create/edit stage controls exclude Won/Lost;
   Sales owner field is self-locked; user/role and operation log pages are admin
   only; last admin disable/downgrade is visibly disabled; data range copy stays
   Sales self / Manager team / Admin all.
8. **Design-is-implemented acceptance:** UI affordances cannot invite an action
   that server authorization forbids; hidden vs disabled follows the yardstick;
   no frontend change widens data scope or changes role strings.
9. **Verification:** e2e role scenarios remain green and add UI assertions for
   the A4 rules above.

## UIUX-010 — Accessibility Baseline (A5)

1. **Status:** Done — labels, focus-visible ring, keyboard-reachable nav, state
   landmarks, table wrappers, and representative A5 e2e assertion landed; build
   and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** all interactive components and page types.
4. **Acceptance binding:** ACC-018, ACC-023; A5; C1-C6.
5. **Objective:** Bring the redesigned UI to the required accessibility baseline.
6. **Implementation scope:** semantic landmarks, accessible names, labels,
   visible focus ring, keyboard order, drawer/dialog focus management, aria-live
   announcements, table labels, target sizes.
7. **Required behavior:** every pointer action is keyboard reachable; focus is
   restored after drawer/dialog/focus-stage close; forms have labels; status is
   not color-only; contrast meets AA with locked colors.
8. **Design-is-implemented acceptance:** keyboard-only operation works for nav,
   search/filter, row selection, pagination, dialogs, live toggle, and all primary
   CRUD/report/admin flows; focus-visible ring uses the locked token treatment.
9. **Verification:** Playwright keyboard path checks plus manual focused-state
   screenshot evidence for representative pages.

## UIUX-011 — Desktop-First Responsive Stability (A6)

1. **Status:** Done — desktop-first responsive shell, wrapped wide tables,
   stable grid/card dimensions, and narrow-desktop overflow e2e assertion
   landed; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** shell, grids, tables, forms, reports, special pages.
4. **Acceptance binding:** ACC-018, ACC-023; A6; C1-C6.
5. **Objective:** Ensure the locked 1440px desktop composition is sharp and
   narrower desktop/tablet widths degrade without overflow, overlap, or broken
   controls.
6. **Implementation scope:** CSS grid constraints, table overflow strategy,
   fixed-height panels, min-width/min-height rules, text truncation, responsive
   stacking where required.
7. **Required behavior:** no full mobile breakpoint commitment in this scope, but
   pages must not horizontally explode or overlap at reasonable desktop/tablet
   widths; fixed grids use `min-height:0` and stable dimensions from the design
   system.
8. **Design-is-implemented acceptance:** screenshots at 1440px and at least one
   narrower viewport show no incoherent overlap, clipped controls, or unusable
   table/action areas.
9. **Verification:** browser screenshot review and build green.

## UIUX-012 — Conservative Motion And Reduced-Motion Path (A7)

1. **Status:** Done — conservative motion tokens, focus/live pulse, skeleton
   shimmer, rail flyout timing, and reduced-motion snap rules are implemented in
   CSS; build and e2e green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** dashboard focus stage, rail flyout, live/count
   affordances, skeletons, hover/focus/selected transitions, toasts.
4. **Acceptance binding:** ACC-018, ACC-023; A7; C1-C6.
5. **Objective:** Implement the accepted conservative motion layer from
   `interaction-spec.md` Part B without decorative or blocking animation.
6. **Implementation scope:** motion CSS tokens, `prefers-reduced-motion` rules,
   count-up utility for meaningful KPI deltas, live update highlight, rail
   flyout timing, skeleton shimmer, toast transitions.
7. **Required behavior:** NAV-01 flyout is hover/focus accessible; MOTION-02
   count-up is conservative; LIVE-03 debounce/coalesce behavior is represented
   where current data refresh behavior exists; reduced-motion snaps and loses no
   information.
8. **Design-is-implemented acceptance:** no animation exceeds the specified
   tokens; no layout/reflow animation is introduced; all motion has a static
   readable equivalent.
9. **Verification:** reduced-motion browser check and focused screenshots for
   rail flyout/live/count states.

## UIUX-013 — Fold In G8 Observations A8/A9

1. **Status:** Done — A8 reminder type/status/priority display uses real
   backend values through label maps; A9 import/export run fields, audit/cleanup
   status, retainedUntil, file safety, and archivedIncluded display are aligned;
   full e2e is green.
2. **Owner agent:** frontend-engineer.
3. **Affected UI surface:** Reminder Center; Import/Export.
4. **Acceptance binding:** ACC-018, ACC-023; A8, A9; C1-C6.
5. **Objective:** Close the two implementation observations explicitly folded
   into the yardstick.
6. **Implementation scope:** `frontend/src/i18n/labels.ts` display maps only if
   needed; Reminder Center status/priority display; Import/Export sample/result
   wording and `archivedIncluded` display consistency.
7. **Required behavior:** status/priority values are displayed from real backend
   values without changing backend comparisons; ImportRun/ExportRun fields align
   with `frontend/src/api/importexport.ts`; export result must not contradict the
   includeArchived checkbox.
8. **Design-is-implemented acceptance:** Reminder rows show real type/status/
   priority labels; Import/Export result cards and tables match `ImportRun` and
   `ExportRun`, including partial failure row errors, audit status, file safety,
   retention, and archived inclusion.
9. **Verification:** reminder/import/export e2e assertions updated and green;
   static diff confirms labels-only enum display additions.

## UIUX-014 — E2E, Regression, And Handoff Evidence (C5)

1. **Status:** Done — selector updates and added A3/A4/A5/A6/A8/A9 assertions
   landed; `npm run build` green and `npm run test:e2e` passed 45/45 with 0
   skips before G11 handoff.
2. **Owner agent:** qa-execution with frontend-engineer.
3. **Affected UI surface:** full frontend test suite and representative browser
   screenshots.
4. **Acceptance binding:** ACC-018, ACC-023; A1-A9; C1-C6.
5. **Objective:** Update Playwright selectors for the redesigned DOM while
   preserving or increasing assertion strength, then record implementation
   evidence for G11/G12.
6. **Implementation scope:** `frontend/e2e/*.spec.ts`, frontend build/test
   commands, screenshot/manual evidence files if used.
7. **Required behavior:** no `test.skip`, `test.only`, disabled suites, weakened
   assertions, or English-copy regressions. Selector updates should prefer roles,
   labels, and stable accessible names.
8. **Design-is-implemented acceptance:** all existing e2e scenarios remain
   covered; added checks cover A3 state samples, A4 role gates, A5 keyboard
   baseline, A8/A9 observations, and representative visual page-type presence.
9. **Verification:** `npm run build` and `npm run test:e2e` pass from
   `frontend/`; report includes no backend/API/model diff for this scope.
