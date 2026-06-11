# UI/UX G12 E2E Stability Evidence - A7 Configured Duration Fix

Date: 2026-06-11
Scope: BLK-UIUX-G12-015 continuation, e2e-only.
Status: returned to Claude for independent re-verification. Codex does not self-pass G12 or self-resolve the blocker.

## Reason

Claude re-opened BLK-UIUX-G12-015 after `TEST-UIUX-A7-001` flaked again under the capped full suite (`workers: 4`). The remaining unstable part was the A7 timing assertion path: the test still waited for runtime animation completion / timing records before asserting the 450/310/220/80ms totals. That can still be sensitive to residual CPU contention.

The deterministic replacement is to assert the configured CSS timing, not the wall-clock runtime. This still catches the BLK-UIUX-G12-012 regression class (for example 220ms accidentally replacing 450ms) while removing scheduler jitter from the assertion.

## E2E Change

Changed only `frontend/e2e/overview.spec.ts`.

- Kept named animation trigger checks: `dashboardStageEnter`, `dashboardStripEnter`, `dashboardStageSwitch`, `dashboardStageExit`, `dashboardStripExit`, and `dashboardReducedFocusAppear` must still fire.
- Replaced `expectDashboardAnimationTotal*` with `expectDashboardAnimationConfiguredTotal*`.
- The recorder now stores `getComputedStyle(target).animationDuration` and `animationDelay` at the animation trigger/snapshot point.
- Removed the `getAnimations().finished` wait from the duration assertion path.
- Tightened configured-duration tolerance to 1ms for full-motion totals:
  - Stage enter: configured total 450ms (`410ms` duration + `40ms` delay).
  - Stage switch: configured total 220ms.
  - Stage exit: configured total 310ms.
  - Reduced appear: configured total <= 80ms.
- Reduced-motion still asserts `transform:none` / identity and no full-motion travel keyframes.

No product code, CSS, backend, shared code, or root API behavior changed.

## Verification

Targeted A7:

```text
> npm run test:e2e -- --grep TEST-UIUX-A7-001

Running 1 test using 1 worker

  ✓  1 [chromium] › e2e/overview.spec.ts:268:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (3.1s)

  1 passed (4.2s)
```

Full e2e smoke:

```text
> npm run test:e2e

Running 52 tests using 4 workers

  ✓  48 [chromium] › e2e/overview.spec.ts:268:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (2.9s)
  ✓  50 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (1.0s)
  ✓  51 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.3s)
  ✓  52 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (425ms)

  52 passed (34.6s)
```

Hygiene:

- `rg -n "\b(test|describe)\.(only|skip|slow)|test\.skip|describe\.skip|test\.only|describe\.only" frontend/e2e frontend/playwright.config.ts`
  - Result: no matches.
- `git diff --check`
  - Result: PASS, no output.
- `git diff --name-only -- services shared api backend frontend/src`
  - Result: no output.

## Handoff

BLK-UIUX-G12-015 is returned to Claude for the requested independent >=6 consecutive full-suite verification. Codex does not self-resolve.
