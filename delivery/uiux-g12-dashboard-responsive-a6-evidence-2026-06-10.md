# UI/UX G12 Dashboard Responsive A6 Evidence — BLK-UIUX-G12-010

Date: 2026-06-10  
Owner: Codex (frontend execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-010`
- Acceptance basis: `docs/ux-ui/requirements/uiux-implementation.requirements.md` A6 (`桌面优先 + 响应式不破版`).
- Surface: dashboard/workbench KPI strip, manager dashboard grid, sales dashboard grid, donut legend, dashboard row text, and KPI metric tile layout.
- Out of scope: backend/API/shared changes, role/enum semantics, Card -> Focus hero motion (`BLK-UIUX-G12-006`), persistent focus selector rail (`BLK-UIUX-G12-007`), single-return focus header (`BLK-UIUX-G12-008`), and B3 live polling (`BLK-UIUX-G12-009`).

## Spec / Decision Note

- No new DEC was required for this patch. This is an implementation correction under the locked A6 responsive/no-break requirement and the release-owner blocker `BLK-UIUX-G12-010`.
- The correction keeps the 1440px desktop-first composition intact while adding graceful degradation for narrower desktop widths that were previously untested.

## Implementation

- `frontend/src/styles/design-system.css`
  - Adds dashboard-local responsive column reduction:
    - `<=1439px`: KPI strip `4 -> 2`; manager dashboard `4 -> 2`; sales dashboard `3 -> 2`.
    - `<=1079px`: KPI strip and dashboard card grids become single-column for graceful degradation near 1024px and below.
  - Keeps dashboard grid definitions scoped to `.dashboardKpis`, `.managerDashboardGrid`, and `.salesDashboardGrid`.
  - Prevents donut legend CJK character wrapping by applying single-line ellipsis to legend labels and nowrap to percentages.
  - Prevents KPI value/icon collision by adding `min-width:0`, block value rendering, `overflow-wrap:anywhere`, and fixed icon sizing.
  - Fixes dashboard list row overflow by giving row grid children `min-width:0` and making row meta text a block ellipsis target instead of relying on inline overflow behavior.
- `frontend/e2e/overview.spec.ts`
  - Strengthens `TEST-UIUX-A6-001` to cover `1680`, `1440`, `1280`, `1180`, `1024`, and `900` widths.
  - Extends card overflow assertions to include both dashboard data cards and top KPI metric tiles.
  - Adds geometry checks for KPI value/icon overlap and donut legend wrapping.
  - Keeps existing dashboard count and live-report visibility checks.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Manager grid degrades from 4 columns | PASS — `<=1439px` uses 2 columns; `<=1079px` uses 1 column. |
| KPI strip degrades from 4 columns | PASS — `<=1439px` uses 2 columns; `<=1079px` uses 1 column. |
| Sales grid degrades from 3 columns | PASS — `<=1439px` uses 2 columns; `<=1079px` uses 1 column. |
| 1200-1400px band clean | PASS — e2e now covers 1280px and 1180px with no card scroll/client overflow, no legend wrapping, and no KPI overlap. |
| 1024px graceful degradation | PASS — e2e covers 1024px with no card scroll/client overflow, no legend wrapping, and no KPI overlap. |
| 1440px primary look preserved | PASS — e2e still covers 1440px; grid remains 4-column manager layout at 1440px with text overflow protections only. |
| Donut legend no char-wrap | PASS — legend labels are nowrap/ellipsis; e2e checks legend label height and `white-space`. |
| KPI amount does not overlap icon | PASS — KPI value can wrap and icon keeps fixed size; e2e checks value/icon rectangles do not intersect. |
| Hero/focus/live behavior unaffected | PASS — overview e2e still passes focus, A7 reduced-motion, and B3 live polling tests. |
| Pure frontend | PASS — only CSS, overview e2e, delivery evidence, and planning records changed. |
| 0 backend/shared/root-api diff | PASS — `git diff --name-only -- services shared api backend` and `git status --short -- services shared api backend` returned no files. |
| No new color | PASS — CSS/e2e diff color-literal scan returned no matches. |
| zh-CN preserved | PASS — no new visible English UI text added. |
| Enum/role comparison values unchanged | PASS — no enum/role comparison logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 50/50 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npx playwright test e2e/overview.spec.ts`  
  Result: PASS, exit 0; 7 passed, 10.6s.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built `index-adHxixKd.css` and `index-CC1VeCcU.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 50 passed, 0 skipped, 16.4s.
- `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --check`  
  Result: PASS, no whitespace errors.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- frontend/src/styles/design-system.css frontend/e2e/overview.spec.ts | rg "^\\+.*(#|rgba?\\(|hsla?\\(|color-mix\\(|linear-gradient\\()"`  
  Result: PASS, no matches.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `110789608a503f079d7ce3061ff5e26dad3d52613475fa319dbbb90f9442118e` |
| `frontend/dist/assets/index-adHxixKd.css` | `a3eb96cdda01fe24f0b4e55ae6fffcb153ddbcbf13c88143413af1d52c40c938` |
| `frontend/dist/assets/index-CC1VeCcU.js` | `b3ccacb3bef39a10f84a687a0f83f7bdcf362c359228859a31db58a644a3c7b6` |

## Handoff

BLK-UIUX-G12-010 is returned to Claude for independent re-verification. Codex does not self-resolve or self-pass G12.
