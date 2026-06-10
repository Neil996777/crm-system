# UI/UX G12 Hero Timing Evidence — BLK-UIUX-G12-012

Date: 2026-06-10  
Owner: Codex (UX design + frontend execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-012`
- Release-owner direction: slow the overview -> focus hero entrance from 320ms to ~450ms, scale reverse to ~310ms, keep focus-rail switching at 220ms.
- Surfaces: manager/sales dashboard Card -> Focus stage transition only.
- Out of scope: backend/API/shared changes, role/enum semantics, focus selector composition (`BLK-UIUX-G12-007`), single-return header (`BLK-UIUX-G12-008`), B3 polling (`BLK-UIUX-G12-009`), and A6 responsive fixes (`BLK-UIUX-G12-010/011`).

## Spec / Decision Record

- `docs/ux-ui/interaction-spec.md`
  - B1 now records a narrow exception under `DEC-UX-HEROTIME-01`: the global B1 scale remains `80/140/220/320ms`, but the Card -> Focus hero uses dedicated timing values of ~450ms enter and ~310ms reverse.
  - B2 choreography now states the hero enter window as ~450ms, with the card FLIP/content reveal running 40-450ms and the selector-rail collapse stagger scaled to ~96-450ms.
  - Reverse is now ~310ms and still quicker than enter.
  - Focus-rail switching remains `motion-base` 220ms and is explicitly not part of the slower hero timing.
  - Decisions includes `DEC-UX-HEROTIME-01`, citing `planning/blockers.md` `BLK-UIUX-G12-012`.
- `docs/ux-ui/design-system.md`
  - Focus layout notes now point to `DEC-UX-HEROTIME-01` for the hero timing and state that selector switching is unaffected.

## Implementation

- `frontend/src/styles/motion.css`
  - Added dedicated hero variables: `--motion-hero: 450ms` and `--motion-hero-exit: 310ms`.
  - Existing B1 variables are unchanged: `--motion-instant: 80ms`, `--motion-fast: 140ms`, `--motion-base: 220ms`, `--motion-slow: 320ms`.
- `frontend/src/pages/WorkOverview.tsx`
  - Updated JS phase constants to `focusEnterMs = 450` and `focusExitMs = 310`.
  - Kept `focusSwitchMs = 220` and `focusReducedMs = 80`.
  - Writes `--motion-hero` and `--motion-hero-exit` to the focus root from the same JS constants so the CSS timeline and React phase timeout stay aligned.
- `frontend/src/styles/design-system.css`
  - Enter: `dashboardStageEnter` and `dashboardStageContentIn` run for `calc(var(--motion-hero) - 40ms)` with the existing 40ms delay, ending at 450ms.
  - Enter selector rail: `dashboardStripEnter` uses the hero window with 96ms + 30ms/index stagger and a 144ms duration (`calc(var(--motion-hero) - 306ms)`), so the last manager rail item ends at 450ms.
  - Exit: `dashboardStageExit`, `dashboardStageContentOut`, and `dashboardScrimOut` use `--motion-hero-exit` (310ms).
  - Exit selector rail: `dashboardStripExit` keeps the 16ms/index reverse stagger and uses a 198ms duration (`calc(var(--motion-hero-exit) - 112ms)`), so the last manager rail item ends at 310ms.
  - `dashboardStageSwitch` remains on `--motion-base` 220ms.
  - Reduced-motion rules remain on `--motion-instant` opacity-only animations with `transform:none`.
- `frontend/src/components/ui/index.tsx`
  - Scaled `--strip-enter-delay` from `80 + index * 24` to `96 + index * 30`.
  - Kept `--strip-exit-delay` at `index * 16` to fit the 310ms exit window with the new CSS duration.
- `frontend/e2e/overview.spec.ts`
  - `TEST-UIUX-A7-001` still records named `animationstart` events.
  - The recorder now also captures CSS animation duration + delay and asserts totals:
    - `dashboardStageEnter` ~= 450ms.
    - `dashboardStageExit` ~= 310ms.
    - `dashboardStageSwitch` ~= 220ms.
  - Existing assertions for enter/strip/switch/exit/reduced-motion names, stable 8-item selector rail, single return control, and focus restoration remain.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Hero enter is ~450ms | PASS — CSS total is 40ms delay + 410ms stage animation; e2e asserts `dashboardStageEnter` total ~= 450ms. |
| Full B2 enter choreography scaled | PASS — stage/content use the 450ms hero window; selector rail stagger starts at 96ms, steps by 30ms, and the last manager rail item ends at 450ms. |
| Reverse is ~310ms and quicker than enter | PASS — stage/content/scrim exit use 310ms; strip exit max delay 112ms + 198ms duration ends at 310ms; e2e asserts `dashboardStageExit` total ~= 310ms. |
| Focus-rail switch remains 220ms | PASS — `focusSwitchMs` and `dashboardStageSwitch` remain `--motion-base`; e2e asserts ~= 220ms and switch still does not trigger `dashboardStageExit`. |
| CSS timing and JS phase constants stay in sync | PASS — JS constants are `450/310/220/80`; focus root exposes the same hero values as CSS variables; e2e records real CSS totals. |
| Reduced-motion remains snap/no travel | PASS — `focusReducedMs` stays 80ms; reduced CSS still uses `dashboardReducedFocusAppear/Exit` with `transform:none`; A7 reduced assertions still pass. |
| Prior focus rail/header/live/responsive fixes unaffected | PASS — overview e2e passes all 7 dashboard tests, including A6, B3, A7, rail selection, and single-return header assertions. |
| Pure frontend | PASS — changed docs, frontend CSS/TSX/e2e, delivery/planning records only. |
| 0 backend/shared/root-api diff | PASS — backend/shared/root-api diff and status checks returned no files. |
| No new color | PASS — exact CSS color-literal diff scan returned no matches. |
| zh-CN preserved | PASS — no new visible English UI text added. |
| Enum/role comparison values unchanged | PASS — no enum/role comparison logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 50/50 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npx playwright test e2e/overview.spec.ts -g TEST-UIUX-A7-001`  
  Result: PASS, exit 0; 1 passed, 2.9s.
- `cd frontend && npx playwright test e2e/overview.spec.ts`  
  Result: PASS, exit 0; 7 passed, 12.0s.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built `index-hY67B0Hc.css` and `index-qO7se2PA.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 50 passed, 0 skipped, 21.6s.
- `rg -n "\\b(test|describe)\\.(only|skip)|test\\.skip|describe\\.skip" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- frontend/src/styles/design-system.css frontend/src/styles/motion.css frontend/src/components/ui/index.tsx frontend/src/pages/WorkOverview.tsx frontend/e2e/overview.spec.ts docs/ux-ui/interaction-spec.md docs/ux-ui/design-system.md | rg -n "^\\+.*(#[0-9A-Fa-f]{3,8}\\b|rgba?\\(|hsla?\\(|color-mix\\(|linear-gradient\\()"`
  Result: PASS, no matches.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `4db202cc67509aa58dbaeb1c697287dce0c46f732281fd4ec2ee155ec9c6fd9a` |
| `frontend/dist/assets/index-hY67B0Hc.css` | `9aff5abb39ab579a1fcda27ecb04230695773ba9d780f53f760b7611c624f209` |
| `frontend/dist/assets/index-qO7se2PA.js` | `3e147c5d22e7ee7c91455d17bf8ef2ba93541bc1ae59973587107974c39b7d2d` |

## Handoff

BLK-UIUX-G12-012 is returned to Claude for independent re-verification. Codex does not self-resolve or self-pass G12.
