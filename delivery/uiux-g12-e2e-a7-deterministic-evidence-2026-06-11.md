# UI/UX G12 e2e stability evidence — A7 deterministic reset

Date: 2026-06-11  
Scope: BLK-UIUX-G12-015 continuation, `TEST-UIUX-A7-001` only plus existing e2e/config work  
Return: Codex -> Claude re-verification. Codex does not self-resolve G12.

## Change

`TEST-UIUX-A7-001` no longer observes transient animation-in-progress state. The test now avoids:

- `animationstart` event capture / global animation recorder.
- `data-transition-phase` polling.
- runtime `getComputedStyle(element).animationName` polling during the 450ms/310ms/220ms motion window.
- WAAPI/runtime wall-clock duration measurement.

The replacement is deterministic:

- Always-present configuration check: reads `getComputedStyle(document.documentElement)` and asserts `--motion-hero: 450ms`, `--motion-hero-exit: 310ms`, `--motion-base: 220ms`, and `--motion-instant: 80ms`.
- Static stylesheet-source check: fetches same-origin stylesheet text and asserts the configured selectors still wire `dashboardStageEnter`, `dashboardStripEnter`, `dashboardStageSwitch`, `dashboardStageExit`, `dashboardStripExit`, `dashboardReducedFocusAppear`, and reduced `transform: none`.
- Stable terminal-state checks: after entry, switch, exit, and reduced entry, asserts final DOM, focus, rail count/order/selection, stage heading, dashboard restoration, restored card focus, and reduced final transform.

## Acceptance mapping

- A7 / BLK-012 timing regression guard: covered by constant CSS token assertions and stylesheet rule wiring, not runtime clock measurement.
- B2 focus behavior: entry reaches `[data-uiux="dashboard-focus"]`, right rail has 8 items in stable order, selected card has `aria-current`, stage heading is focused, switch changes only selection/stage content, and Esc returns to dashboard with DOM focus restored to the originating card.
- B6 reduced motion: reduced mode reaches the same terminal focus state, right rail remains 8 items, and `.stage` final transform is `none` / identity.
- C5 stability: removes all remaining A7 assertions that depended on a short-lived transition window.

## Constraints

- Product code unchanged in this continuation.
- `workers: 2` retained in `frontend/playwright.config.ts`.
- No retry, no `test.skip`, no `test.only`, no `test.slow`, no assertion deletion for the product behavior under test.
- 0 backend / `shared` / root-api / `frontend/src` diff in this continuation check.
- zh-CN, enum values, role comparison values, colors, and UI behavior unchanged.

## Verification

Commands executed from `frontend/` unless noted.

### TypeScript

Command:

```bash
npx tsc --noEmit
```

Result: PASS, no output.

### Build

Command:

```bash
npm run build
```

Result: PASS.

Tail:

```text
dist/index.html                   0.40 kB │ gzip:   0.28 kB
dist/assets/index-TNHK1XAP.css   56.31 kB │ gzip:  10.80 kB
dist/assets/index-BxWQVDfj.js   414.07 kB │ gzip: 113.28 kB
✓ built in 825ms
```

Artifact SHA-256:

```text
d76afb49e74c93175ee45e9ba92228fbbc3732faaffe58143d11941f84b92b8e  dist/assets/index-TNHK1XAP.css
7f13cfbbb558947556a18e2e0fc559ac65eb7f0f60bd46490d47de1936cf16ae  dist/assets/index-BxWQVDfj.js
```

### Targeted A7

Command:

```bash
npx playwright test e2e/overview.spec.ts -g TEST-UIUX-A7-001
```

Result:

```text
Running 1 test using 1 worker
✓  1 [chromium] › e2e/overview.spec.ts:267:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (2.9s)
1 passed (3.6s)
```

### Full e2e, 8 consecutive runs

Initial sandboxed loop could not bind the Vite web server (`listen EPERM 127.0.0.1:5173`), so the full-suite loop was rerun with approved local web-server permission. The following eight runs are one continuous verification pass after the implementation change, with no rerun-on-failure.

Command:

```bash
npm run test:e2e
```

Result tails:

```text
===== FULL E2E RUN 1 =====
Running 52 tests using 2 workers
52 passed (27.0s)

===== FULL E2E RUN 2 =====
Running 52 tests using 2 workers
52 passed (26.7s)

===== FULL E2E RUN 3 =====
Running 52 tests using 2 workers
52 passed (27.6s)

===== FULL E2E RUN 4 =====
Running 52 tests using 2 workers
52 passed (26.6s)

===== FULL E2E RUN 5 =====
Running 52 tests using 2 workers
52 passed (26.9s)

===== FULL E2E RUN 6 =====
Running 52 tests using 2 workers
52 passed (26.2s)

===== FULL E2E RUN 7 =====
Running 52 tests using 2 workers
52 passed (26.0s)

===== FULL E2E RUN 8 =====
Running 52 tests using 2 workers
52 passed (26.3s)
```

Summary: 8/8 full-suite runs passed, 52/52 each, 0 failed, 0 skipped, no rerun.

### Guard scans

Command:

```bash
rg -n "expectDashboardTransitionPhase|expectDashboardAnimationApplied|expectDashboardAnimationConfiguredTotal|dashboardAnimationState|animationstart|__dashboardAnimation|dashboardAnimationTimeoutMs|\bLocator\b" frontend/e2e/overview.spec.ts
```

Result: PASS, no matches.

Command:

```bash
rg -n "\b(test|describe)\.(only|skip|slow)|test\.skip|describe\.skip|test\.only|describe\.only" frontend/e2e frontend/playwright.config.ts
```

Result: PASS, no matches.

Command:

```bash
git diff --check
```

Result: PASS, no output.

Command:

```bash
git diff --name-only -- services shared api backend frontend/src
```

Result: PASS, no output.

## Handoff

Status is returned to Claude for independent >=8-run full-suite verification. Codex does not self-pass UI/UX G12 or resolve BLK-UIUX-G12-015.
