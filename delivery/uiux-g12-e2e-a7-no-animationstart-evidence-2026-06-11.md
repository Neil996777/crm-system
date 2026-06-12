# UI/UX G12 E2E Stability Evidence - A7 No Animationstart Dependency

Date: 2026-06-11
Scope: BLK-UIUX-G12-015 continuation, e2e only.
Status: returned to Claude for independent re-verification. Codex does not self-pass G12 or self-resolve the blocker.

## Reason

Claude re-verification #6 confirmed `workers: 2` stabilised the wider contention class, but `TEST-UIUX-A7-001` still flaked once in 8 full-suite runs. The remaining A7 failure was tied to the test capturing transient `animationstart` events through `installDashboardAnimationRecorder` / `expectDashboardAnimationStarted(...)`.

## Change

Changed `frontend/e2e/overview.spec.ts` only:

- removed the `animationstart` listener, global animation recorder, reset helper, event-capture snapshot helper, and event-based animation-start polling.
- kept `workers: 2` in `frontend/playwright.config.ts` from the prior release-owner decision.
- A7 now verifies active motion state by deterministic DOM/CSS state:
  - `data-transition-phase="entering"` + `.stage` computed `animationName` contains `dashboardStageEnter`.
  - entering rail item computed `animationName` contains `dashboardStripEnter`.
  - `data-transition-phase="switching"` + `.dashboardStageContent` computed `animationName` contains `dashboardStageSwitch`, while `.stage` does not apply `dashboardStageExit`.
  - `data-transition-phase="exiting"` + `.stage` / side-card computed `animationName` contain `dashboardStageExit` / `dashboardStripExit`.
  - reduced path uses `data-transition-phase="reduced-entering"` + `.stage` computed `animationName` contains `dashboardReducedFocusAppear`, with no `dashboardStageEnter` / `dashboardStripEnter`.
- timing assertions still read configured CSS `animationDuration + animationDelay`, not runtime wall-clock duration:
  - StageEnter: 450ms
  - StageSwitch: 220ms
  - StageExit: 310ms
  - ReducedFocusAppear: <=80ms
- no retries, skip/only, `test.slow`, or assertion weakening were introduced.
- no product code, CSS, backend, shared code, or root API behavior changed.

## Verification

Pre-flight:

```text
npx tsc --noEmit
Result: PASS, no output.

npx playwright test e2e/overview.spec.ts -g TEST-UIUX-A7-001
Running 1 test using 1 worker
  ✓  1 [chromium] › e2e/overview.spec.ts:268:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (3.3s)
  1 passed (4.3s)

rg -n "installDashboardAnimationRecorder|expectDashboardAnimationStarted|resetDashboardAnimationRecorder|animationstart|__dashboardAnimation" frontend/e2e/overview.spec.ts
Result: no matches.
```

Eight consecutive full-suite runs, no reruns:

### Run 1

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (4.0s)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (822ms)
  ✓  48 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (381ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (566ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (807ms)
  ✓  51 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (914ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (975ms)

  52 passed (32.6s)
```

### Run 2

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (836ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (564ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (606ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (610ms)
  ✓  46 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.1s)
  ✓  52 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (325ms)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (669ms)

  52 passed (29.4s)
```

### Run 3

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (867ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (546ms)
  ✓  45 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.6s)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (1.2s)
  ✓  50 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (431ms)
  ✓  51 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (638ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (704ms)

  52 passed (29.0s)
```

### Run 4

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  45 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (398ms)
  ✓  41 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.3s)
  ✓  47 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (530ms)
  ✓  46 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (1.1s)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (864ms)
  ✓  49 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (766ms)
  ✓  50 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (531ms)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (1.2s)
  ✓  52 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (916ms)

  52 passed (29.8s)
```

### Run 5

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  45 [chromium] › e2e/retrieval.spec.ts:13:1 › TEST-NAV-RETRIEVE-001 lists and details contacts from the primary navigation (1.1s)
  ✓  46 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (837ms)
  ✓  47 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (466ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (536ms)
  ✓  49 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (915ms)
  ✓  50 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (730ms)
  ✓  51 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (593ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (825ms)

  52 passed (36.3s)
```

### Run 6

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.3s)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (1.0s)
  ✓  48 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (303ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (854ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (1.2s)
  ✓  51 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (1.7s)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (2.0s)

  52 passed (33.3s)
```

### Run 7

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (781ms)
  ✓  43 [chromium] › e2e/overview.spec.ts:351:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.3s)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (487ms)
  ✓  49 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (424ms)
  ✓  50 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (822ms)
  ✓  51 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (699ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (688ms)

  52 passed (34.1s)
```

### Run 8

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  48 [chromium] › e2e/overview.spec.ts:389:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (464ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (950ms)
  ✓  50 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (527ms)
  ✓  49 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (748ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (952ms)
  ✓  51 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (986ms)

  52 passed (30.4s)
```

## Hygiene

```text
rg -n "\b(test|describe)\.(only|skip|slow)|test\.skip|describe\.skip|test\.only|describe\.only" frontend/e2e frontend/playwright.config.ts
Result: no matches.

git diff --check
Result: PASS, no output.

git diff --name-only -- services shared api backend frontend/src
Result: no output.
```

## Handoff

BLK-UIUX-G12-015 is returned to Claude for the requested independent >=8 consecutive full-suite verification. Codex does not self-resolve.
