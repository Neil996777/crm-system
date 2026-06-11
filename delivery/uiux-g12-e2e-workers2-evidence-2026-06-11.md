# UI/UX G12 E2E Stability Evidence - Workers:2 Return

Date: 2026-06-11
Scope: BLK-UIUX-G12-015 continuation, e2e/config only.
Status: returned to Claude for independent re-verification. Codex does not self-pass G12 or self-resolve the blocker.

## Reason

Claude re-verification #5 found that the A7 configured-duration fix was correct, but the full suite still flaked under `workers:4` because of systemic CPU contention. The release owner decided to lower Playwright concurrency to `workers:2` for both local and CI runs.

## Change

Changed only `frontend/playwright.config.ts`:

- `workers: 4` -> `workers: 2`
- no retries added
- no `test.skip`, `test.only`, or `test.slow`
- no assertion weakening
- no product code, CSS, backend, shared code, or root API behavior changed

Existing stabilizations remain in place:

- A7 animation timing checks read configured CSS values via `getComputedStyle(...).animationDuration` + `animationDelay`, not runtime wall-clock duration.
- B3 dashboard live polling route teardown remains guarded.
- Task/activity readiness waits remain deterministic.

## Verification

Six consecutive full-suite runs, no reruns:

### Run 1

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (343ms)
  ✓  45 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (607ms)
  ✓  44 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (821ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (702ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (348ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (469ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (480ms)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (623ms)
  ✓  46 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.2s)
  ✓  52 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (264ms)

  52 passed (24.3s)
```

### Run 2

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  42 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (1.3s)
  ✓  44 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (473ms)
  ✓  46 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (917ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (858ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (400ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (561ms)
  ✓  45 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.6s)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (732ms)
  ✓  51 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (301ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (661ms)

  52 passed (27.7s)
```

### Run 3

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/retrieval.spec.ts:13:1 › TEST-NAV-RETRIEVE-001 lists and details contacts from the primary navigation (594ms)
  ✓  42 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (891ms)
  ✓  44 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (390ms)
  ✓  46 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (880ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (757ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (442ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (542ms)
  ✓  45 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.2s)
  ✓  51 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (451ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (829ms)
  ✓  52 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (608ms)

  52 passed (26.9s)
```

### Run 4

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  42 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (816ms)
  ✓  43 [chromium] › e2e/retrieval.spec.ts:13:1 › TEST-NAV-RETRIEVE-001 lists and details contacts from the primary navigation (673ms)
  ✓  45 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (556ms)
  ✓  46 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (1.4s)
  ✓  44 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.3s)
  ✓  48 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (435ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (1.8s)
  ✓  49 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (606ms)
  ✓  50 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (507ms)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (795ms)
  ✓  52 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (856ms)

  52 passed (26.5s)
```

### Run 5

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (370ms)
  ✓  40 [chromium] › e2e/overview.spec.ts:268:1 › TEST-UIUX-A7-001 card focus respects reduced-motion mode and still snaps between states (2.2s)
  ✓  44 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (1.1s)
  ✓  45 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (1.1s)
  ✓  46 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (853ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (460ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (607ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (542ms)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (678ms)
  ✓  47 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.0s)
  ✓  52 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (282ms)

  52 passed (25.7s)
```

### Run 6

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 2 workers

  ✓  43 [chromium] › e2e/retrieval.spec.ts:41:1 › TEST-NAV-RETRIEVE-003/004 shows empty state and invalid filter feedback (391ms)
  ✓  44 [chromium] › e2e/overview.spec.ts:337:1 › TEST-UIUX-FOCUS-LAYOUT-001 focus stage keeps real desktop width after rail collapse (856ms)
  ✓  45 [chromium] › e2e/retrieval.spec.ts:55:1 › TEST-NAV-RETRIEVE-005 hides unauthorized records from sales lists (847ms)
  ✓  47 [chromium] › e2e/user-admin.spec.ts:13:1 › TEST-USER-ADMIN-001 creates user and changes role/status with confirmation (733ms)
  ✓  48 [chromium] › e2e/user-admin.spec.ts:56:1 › TEST-INV-LASTADMIN-001 blocks disabling or downgrading the last active Administrator (435ms)
  ✓  49 [chromium] › e2e/user-admin.spec.ts:68:1 › TEST-PERM-USERADMIN-002/003 sales is denied user administration (637ms)
  ✓  50 [chromium] › e2e/work.spec.ts:22:1 › TEST-ACTIVITY-NOTE-002 validates missing fields and creates note and activity in record detail (555ms)
  ✓  46 [chromium] › e2e/overview.spec.ts:353:1 › TEST-UIUX-NAV-01 focus rail keeps nav reachable and flyout above the stage (3.0s)
  ✓  51 [chromium] › e2e/work.spec.ts:41:1 › TEST-TASK-LIFECYCLE-002 creates task and completes it from standalone list (660ms)
  ✓  52 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (280ms)

  52 passed (24.8s)
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

BLK-UIUX-G12-015 is returned to Claude for the requested independent >=6-8 consecutive full-suite verification. Codex does not self-resolve.
