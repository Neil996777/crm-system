# UI/UX G12 E2E Stability Evidence - BLK-UIUX-G12-015

Date: 2026-06-11
Scope: UI/UX G12 post-fix #11, e2e/config stability only.
Status: returned to Claude for independent re-verification. Codex does not self-pass G12 or self-resolve the blocker.

## Scope

BLK-UIUX-G12-015 covers full-suite Playwright flakiness under parallel load. Claude re-verification #3 showed the failures were systemic CPU/worker contention rather than one remaining test bug:

- `e2e/work.spec.ts` `TEST-ACTIVITY-NOTE-002` timed out in one full run under load.
- `e2e/overview.spec.ts` `TEST-UIUX-A7-001` still failed in another full run under load.
- Earlier rounds also exposed `TEST-TASK-LIFECYCLE-002` and `TEST-UIUX-B3-001`, both retained as fixed.

This pass changes only Playwright config/e2e tests plus delivery/planning evidence. It does not change product TSX/CSS behavior, backend, shared code, or root API contracts.

## Stabilization Changes

Systemic re-kick #3 closure:

- `frontend/playwright.config.ts` now sets `workers: 4`.
- The cap is unconditional, so local and CI runs use the same stable worker ceiling.
- No `retries` setting was added; Playwright remains at its default 0 retries.
- No `test.skip`, `test.only`, `test.slow`, assertion deletion, or product-behavior workaround was added.

A7 contention hardening retained:

- `TEST-UIUX-A7-001` records dashboard animations through an init-script listener plus a current-document hook.
- The recorder also scans the focus DOM for current computed animation timing when CPU contention delays or misses a test-side `animationstart` observation.
- Animation total assertions wait for the target CSS animation to settle through the Web Animations API (`animation.finished`) before reading the recorded/computed total.
- Full-motion assertions still distinguish the real motion timings: `dashboardStageEnter` ~= 450ms, `dashboardStageExit` ~= 310ms, `dashboardStageSwitch` ~= 220ms, with the existing 24ms tolerance.
- Reduced-motion keeps the effective checks: `dashboardReducedFocusAppear` fires, total <= 80ms, stage transform is `none`/identity, and `dashboardStageEnter`/`dashboardStripEnter` do not fire.
- The focus-stage heading readiness window remains 8s for CPU-contention slack without changing the expected focused element.

Previously returned fixes retained:

- `TEST-TASK-LIFECYCLE-002` waits on real task create/list/status API responses carrying the target task before list/detail/completion assertions.
- `TEST-UIUX-B3-001` uses a named `**/api/activities*` route handler with an `activitiesRouteActive` teardown flag, in-flight handler tracking, guarded `route.fetch()` handling for context/request-abort teardown, `page.unroute(...)` in `finally`, and in-flight route drain before context close.
- `TEST-ACTIVITY-NOTE-002` keeps its real timeline assertions with the existing deterministic 15s refreshed-timeline wait.

## Acceptance Check

| Requirement | Evidence |
|---|---|
| Fix systemic worker contention | Playwright full-suite workers capped at 4 globally in `playwright.config.ts`; no retry/skip/slow masking. |
| Keep A7 timing assertions meaningful | Stage enter/exit/switch exact-total assertions remain 450/310/220ms with 24ms tolerance after WAAPI settled wait. |
| Preserve reduced-motion assertions | Reduced path still requires the dedicated reduced keyframe, no travel transform, no full-motion travel keyframes, and <=80ms snap timing. |
| B3 teardown stays fixed | Route interceptor teardown remains guarded with `page.unroute(...)` and in-flight drain. |
| Task/activity tests keep real readiness | Work tests wait for real API/timeline readiness; assertions were not removed. |
| No rerun-dependent pass | Verification includes five consecutive full `npm run test:e2e` runs after the worker cap, all first-try green. |
| 0 backend/shared/root-api/product-source diff | `git diff --name-only -- services shared api backend frontend/src` empty; `git status --short -- services shared api backend frontend/src` empty. |
| Product behavior unchanged | Source diff is Playwright config/e2e plus planning/evidence only; no frontend product TSX/CSS changed in this pass. |

## Verification

Typecheck/build:

- `npx tsc --noEmit`
  - Result: PASS, exit 0.
- `npm run build`
  - Result: PASS, Vite built 1633 modules.
  - Artifacts:
    - `frontend/dist/index.html` SHA-256 `379e8bb971a5dfc61dc794124f167f96d41a6321196cc37f79d7f255502a9e5e`
    - `frontend/dist/assets/index-TNHK1XAP.css` SHA-256 `d76afb49e74c93175ee45e9ba92228fbbc3732faaffe58143d11941f84b92b8e`
    - `frontend/dist/assets/index-CIQmHaIF.js` SHA-256 `3e147c5d22e7ee7c91455d17bf8ef2ba93541bc1ae59973587107974c39b7d2d`

Five consecutive full e2e runs after the worker cap:

| Run | Command | Result |
|---|---|---|
| 1 | `npm run test:e2e` | PASS: 52 passed, 0 failed, 0 skipped, no rerun, 4 workers, 26.5s |
| 2 | `npm run test:e2e` | PASS: 52 passed, 0 failed, 0 skipped, no rerun, 4 workers, 24.8s |
| 3 | `npm run test:e2e` | PASS: 52 passed, 0 failed, 0 skipped, no rerun, 4 workers, 25.2s |
| 4 | `npm run test:e2e` | PASS: 52 passed, 0 failed, 0 skipped, no rerun, 4 workers, 29.3s |
| 5 | `npm run test:e2e` | PASS: 52 passed, 0 failed, 0 skipped, no rerun, 4 workers, 27.3s |

Full-suite output tails:

Run 1 tail:

```text
  ✓  46 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (716ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (534ms)
  ✓  41 [chromium] › e2e/overview.spec.ts:140:1 › TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload (5.1s)
  ✓  49 [chromium] › e2e/overview.spec.ts:254:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (1.9s)
  ✓  50 [chromium] › e2e/overview.spec.ts:323:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (761ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:339:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (2.9s)
  ✓  52 [chromium] › e2e/overview.spec.ts:377:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (285ms)

  52 passed (26.5s)
```

Run 2 tail:

```text
  ✓  47 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (674ms)
  ✓  46 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (895ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (518ms)
  ✓  42 [chromium] › e2e/overview.spec.ts:140:1 › TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload (5.7s)
  ✓  49 [chromium] › e2e/overview.spec.ts:254:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (1.9s)
  ✓  50 [chromium] › e2e/overview.spec.ts:323:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (786ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:339:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (2.9s)
  ✓  52 [chromium] › e2e/overview.spec.ts:377:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (395ms)

  52 passed (24.8s)
```

Run 3 tail:

```text
  ✓  47 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (745ms)
  ✓  48 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (846ms)
  ✓  42 [chromium] › e2e/overview.spec.ts:140:1 › TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload (5.2s)
  ✓  49 [chromium] › e2e/overview.spec.ts:254:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (1.8s)
  ✓  50 [chromium] › e2e/overview.spec.ts:323:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (469ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:339:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (2.9s)
  ✓  52 [chromium] › e2e/overview.spec.ts:377:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (286ms)

  52 passed (25.2s)
```

Run 4 tail:

```text
  ✓  44 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (2.1s)
  ✓  46 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (1.8s)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (1.4s)
  ✓  49 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (1.3s)
  ✓  47 [chromium] › e2e/overview.spec.ts:254:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (2.3s)
  ✓  50 [chromium] › e2e/overview.spec.ts:323:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (804ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:339:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (2.6s)
  ✓  52 [chromium] › e2e/overview.spec.ts:377:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (303ms)

  52 passed (29.3s)
```

Run 5 tail:

```text
  ✓  47 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (1.1s)
  ✓  40 [chromium] › e2e/overview.spec.ts:140:1 › TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload (4.8s)
  ✓  48 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (835ms)
  ✓  49 [chromium] › e2e/overview.spec.ts:254:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (1.9s)
  ✓  50 [chromium] › e2e/overview.spec.ts:323:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (786ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:339:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (2.9s)
  ✓  52 [chromium] › e2e/overview.spec.ts:377:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (288ms)

  52 passed (27.3s)
```

Hygiene:

- `rg -n "\b(test|describe)\.(only|skip|slow)|test\.skip|describe\.skip|test\.only|describe\.only" frontend/e2e frontend/playwright.config.ts`
  - Result: no matches.
- `git diff --check`
  - Result: PASS, no output.
- `git diff --name-only -- services shared api backend frontend/src`
  - Result: no output.
- `git status --short -- services shared api backend frontend/src`
  - Result: no output.

## Handoff

BLK-UIUX-G12-015 remains In Review for Claude re-verification. Pending Claude re-verification; Codex does not self-resolve.
