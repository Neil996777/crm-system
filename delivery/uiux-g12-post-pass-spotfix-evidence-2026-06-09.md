# UI/UX G12 Post-Pass Spot-Fix Evidence

Date: 2026-06-09
From: Codex (execution)
To: Claude (G12 re-confirm audit)
Status: Spot-fix returned. Codex does not self-pass G12.

## Scope

Release-owner live-use feedback:

- `BLK-UIUX-G12-004` — dashboard card content clipped by fixed grid row height
  and `overflow:hidden`.
- `BLK-UIUX-G12-005` — Card→Focus should activate from the whole card, not only
  the corner expand button.

Changed files:

- `frontend/src/pages/WorkOverview.tsx`
- `frontend/src/styles/design-system.css`
- `frontend/e2e/overview.spec.ts`
- `planning/blockers.md`
- `planning/gate-status.md`

## Implementation Summary

Dashboard clipping:

- `.managerDashboardGrid` now uses `grid-auto-rows: minmax(268px, auto)`.
- `.salesDashboardGrid` now uses `grid-auto-rows: minmax(308px, auto)`.
- Dashboard panels and their overview content containers no longer crop card
  contents: `.dashboardPanel`, `.pipelineViz`, `.dashboardList`,
  `.paymentRows`, and `.timeline` no longer use the clipping path.
- Visual rhythm is preserved through the same minimum row heights and existing
  spacing/tokens.

Card→Focus trigger:

- `DashboardPanel` now makes the whole `.dashboardPanel` the Focus trigger with
  `role="button"`, `tabIndex=0`, click, Enter, and Space handling.
- The expand glyph remains as a visual affordance only (`aria-hidden="true"`),
  avoiding nested interactive controls.
- The dashboard cards remain read-only summaries; activation only changes the
  dashboard presentation into per-card Focus.

E2E coverage strengthened:

- Manager focus test now asserts the dashboard card itself is a keyboard-reachable
  button and opens Focus by clicking the card body.
- Manager "最近活动" and Sales "我的销售漏斗" Focus checks click the whole card.
- Reduced-motion Focus test now uses keyboard Space on the whole card.
- A6 stability test now checks all dashboard cards at 1440px and 1680px for
  zero scroll/client overflow, then keeps the prior 900px no-horizontal-overflow
  assertion.

## Constraint Check

- Frontend-only: yes.
- Backend/shared/root-api diff: 0. `git diff --name-only -- services shared api backend`
  and `git status --short -- services shared api backend` returned no paths.
- Existing endpoints only: yes. No API/backend/shared/root-api files changed.
- New colors: none. CSS diff adds no hex/rgb/hsl/color-mix values; Focus outline
  uses existing `--primary`.
- zh-CN: yes. No user-visible English copy added.
- Enum/role comparison values unchanged: yes.
- No e2e downgrade: no assertions removed or skipped; skip/only scan returned no
  matches.
- Backend-needed aggregation/blocker: none.

## Verification

Commands run from `frontend/` unless noted:

- `npx tsc --noEmit` — PASS.
- `npm run build` — PASS; produced `dist/assets/index-rHR7gSrr.css` and
  `dist/assets/index-CsDtRBPy.js`.
- `npx playwright test e2e/overview.spec.ts` — PASS, 6/6.
- `npm run test:e2e` — PASS, 49/49, 0 skipped.

Commands run from repo root:

- `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e`
  — no matches.
- `git diff --check` — PASS.
- `git diff --name-only -- services shared api backend` — no output.
- `git status --short -- services shared api backend` — no output.
- `git diff -- frontend/src/styles/design-system.css | rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}|rgba?\\(|hsla?\\(|color-mix\\()"`
  — no matches.

## Handoff

Returned to Claude for dashboard screenshot re-verification:

- 8 dashboard cards no longer clip content at 1440px and wider.
- Whole-card click and keyboard activation enter Card→Focus.
- Full e2e remains green with 0 skipped tests.

