# UI/UX G12 A7 Motion Evidence

Date: 2026-06-10
From: Codex (execution)
To: Claude (G12 A7 re-verification audit)
Status: `BLK-UIUX-G12-006` returned for review. Codex does not self-pass G12.

## Scope

Blocker: `BLK-UIUX-G12-006` - Card->Focus hero transition (B2) was previously
instant/blinking instead of the locked A7 motion behavior.

Implemented scope already landed before this evidence return:

- `frontend/src/styles/motion.css`
- `frontend/src/styles/design-system.css`
- `frontend/src/pages/WorkOverview.tsx`
- `frontend/src/components/ui/index.tsx`
- `frontend/e2e/overview.spec.ts`

No implementation, CSS, or e2e files were edited as part of this evidence-only
return.

## Locked Part B Alignment

Source: `docs/ux-ui/interaction-spec.md` Part B.

| Locked item | Evidence |
|---|---|
| B1 motion token scale | `motion.css` mirrors the locked durations: `--motion-instant: 80ms`, `--motion-fast: 140ms`, `--motion-base: 220ms`, `--motion-slow: 320ms`; `--motion-medium` remains as `var(--motion-base)` for existing consumers. Easing tokens match B1, including `ease-standard`, `ease-decelerate`, `ease-accelerate`, and `ease-emphasis`. |
| B2 layout end-states | Overview remains the 8-card dashboard grid; Focus renders the stage plus 7 right-side strip cards. `.shell.focusMode` keeps the collapsed nav rail path, and `.dashboardFocusPage::before` renders the locked scrim value `rgba(15, 23, 42, .06)`. |
| B2 trigger | Whole dashboard cards remain focusable card targets with click, Enter, and Space activation; the visual expand glyph stays non-interactive. |
| B2 0-80ms scrim/chosen-card cue | `dashboardScrimIn`/`dashboardScrimOut` run on the focus page scrim. The chosen card's source DOM rectangle is captured before transition and used as the stage FLIP origin. |
| B2 clicked card -> stage FLIP | `WorkOverview.tsx` captures source/target rectangles with `getBoundingClientRect()` and writes `--hero-start-x`, `--hero-start-y`, `--hero-start-scale-x`, and `--hero-start-scale-y`; `dashboardStageEnter` and `dashboardStageExit` consume those variables with transform/opacity only. |
| B2 compact -> full content crossfade | The stage body is wrapped in `.dashboardStageContent`; `dashboardStageContentIn`, `dashboardStageContentOut`, and `dashboardStageSwitch` animate the content layer independently of layout. |
| B2 other panels -> strip cards | `FocusStage` assigns `motionIndex` plus `--strip-enter-delay`/`--strip-exit-delay`; `dashboardStripEnter`/`dashboardStripExit` use strip FLIP variables and stagger the rail cards. |
| B2 sidebar collapse | Existing `.shell.focusMode` is preserved and label opacity/transform transitions now follow the B1 `motion-fast` token. |
| B2 settle/focus/live region | `FocusStage` heading has `data-focus-heading` and `tabIndex={-1}`; `WorkOverview.tsx` moves focus to the heading after enter/switch and uses a polite `role="status"` live region for focus enter/switch/return messages. |
| B2 reverse via 返回/Esc | `requestExitFocus()` drives `exiting`/`reduced-exiting`; reverse uses `motion-base` 220ms for stage/strip, scrim out uses `ease-accelerate`, and focus is restored to the originating grid card. |
| B2 rail-card switch | `switchFocus()` uses `switching`/`reduced-switching`; CSS uses `dashboardStageSwitch`. The e2e recorder asserts switching does not trigger `dashboardStageExit`. |
| B2 transform/opacity performance guardrail | The new keyframes animate `transform` and `opacity`; no width/height/top/left/margin animation is introduced for the hero/strip/stage paths. |
| B6 reduced motion | `usePrefersReducedMotion()` sets `data-motion-mode="reduced"` and reduced phases; `.dashboardMotionReduced` forces `transform: none` on stage/strip/content and runs `dashboardReducedFocusAppear`/`dashboardReducedFocusExit` opacity-only animations. |
| B6 reduced path tested separately | `TEST-UIUX-A7-001` records `animationstart` events, asserts full-motion enter/strip/switch/exit animations, then asserts reduced mode uses `dashboardReducedFocusAppear` and does not trigger travel animations. |

## Constraint Check

- Backend/shared/root-api diff: 0. `git diff --name-only -- services shared api backend`
  and `git status --short -- services shared api backend` returned no paths.
- Frontend-only: yes. Changes are limited to the frontend/dashboard motion path
  plus planning/evidence documentation.
- New colors: none. CSS diff color scan found only the locked B2 scrim value:
  `background: rgba(15, 23, 42, .06);`.
- zh-CN: preserved. New visible/live-region copy is Chinese.
- Enum/role comparison values: unchanged. The touched frontend code is dashboard
  presentation/motion plumbing and e2e assertions; it does not change persisted
  enum values, role values, or comparison semantics.
- E2E no weakening/no skip: `TEST-UIUX-A7-001` was strengthened from reduced-only
  snap coverage to normal-motion plus reduced-motion path coverage. The skip/only
  scan returned no matches.
- No backend-needed aggregation or new API behavior: none added.

## Verification

Commands run from `frontend/` on 2026-06-10:

| Command | Result |
|---|---|
| `npx tsc --noEmit` | PASS, exit 0. |
| `npm run build` | PASS. Vite transformed 1633 modules and produced `dist/assets/index-CgVjA1iB.css` and `dist/assets/index-Ctu9eh-s.js`. |
| `npm run test:e2e` | PASS. Playwright ran 49 tests using 5 workers: `49 passed (12.6s)`, 0 skipped. |

Additional read-only checks from repo root:

| Command | Result |
|---|---|
| `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e` | No matches. |
| `git diff --name-only -- services shared api backend` | No output. |
| `git status --short -- services shared api backend` | No output. |
| `git diff -- frontend/src/styles/motion.css frontend/src/styles/design-system.css \| rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}\|rgba?\\(\|hsla?\\(\|color-mix\\()"` | Only `+  background: rgba(15, 23, 42, .06);`, the locked B2 scrim. |

## Build Artifact Hashes

SHA-256 after the 2026-06-10 `npm run build`:

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `262338804e944e92a111d636059c7b5ca48c6c02e57d8de618a743c191bcdb9a` |
| `frontend/dist/assets/index-CgVjA1iB.css` | `79b47f2b5849c86863d2b6d0be467a33f9231fc643656cd2a88307ba1f15d8a9` |
| `frontend/dist/assets/index-Ctu9eh-s.js` | `a45706ca209e1c64400de87e3e22d8b4a73fe0e5d4fe5599896db35c1ccf9502` |

## Handoff

Returned to Claude for independent A7/B2 re-verification:

- Live Card->Focus motion observation.
- Motion token scale check against B1.
- Strengthened `TEST-UIUX-A7-001` behavior and reduced-motion path check.
- Cumulative constraints: 0 backend/shared/root-api diff, no new color, zh-CN
  preserved, enum/role values unchanged, full e2e green with no skips.
