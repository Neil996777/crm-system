# UI/UX G12 Live Resume Meta Evidence — BLK-UIUX-G12-016

Date: 2026-06-11  
Owner: Codex (frontend execution + QA execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-016`
- Surface: dashboard live header meta text after `实时更新` / `暂停` toggling.
- Change type: pure frontend bug fix plus e2e regression coverage.
- Out of scope: backend/API/shared changes, new colors/CSS, role or enum semantics, polling cadence changes, manual refresh behavior, B3 pause-buffer behavior beyond the stale meta fix.

## Root Cause

`WorkOverview.tsx` pauses live mode by setting:

```text
refreshNotice = "已暂停 · 更新于 HH:MM"
```

When live mode is resumed, `toggleLiveUpdates()` calls `applyBufferedLive(true)`. If no update was buffered while paused, `pendingLiveSnapshotRef.current` is `null`; the old code returned before resetting the header notice. The button and live state flipped back to real-time, but the visible `.updateMeta` stayed on `已暂停`.

## Implementation

- `frontend/src/pages/WorkOverview.tsx`
  - `applyBufferedLive(true)` now clears `refreshNotice` before returning when there is no buffered snapshot.
  - Clearing the notice lets `DashboardHeader` fall back to the normal live meta: `更新于 HH:MM`.
  - Buffered resume behavior is unchanged: when a buffered snapshot exists, it is still applied and the meta is set to `实时更新 · HH:MM`.
  - Live paused state, `livePausedRef`, live dot, and `aria-pressed` stay aligned with the real live state.
- `frontend/e2e/overview.spec.ts`
  - `TEST-UIUX-B3-001` now covers both resume paths:
    - no buffered updates: pause -> resume, meta must not contain `已暂停`, live dot pulses, toggle has `aria-pressed="true"`;
    - buffered updates: pause -> poll update -> resume, buffered row applies and meta still must not contain `已暂停`.
  - Existing B3 assertions remain: polling applies without reload, `arrived` is change-driven, pause shows the buffer pill, manual refresh shows feedback.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Resume without buffered updates clears stale `已暂停` meta | PASS — e2e pauses and immediately resumes before any buffered update, then asserts `.updateMeta` does not contain `已暂停`. |
| Resume with buffered updates still applies the buffer | PASS — e2e simulates `轮询动态 2`, verifies it is hidden while paused, then visible after resume. |
| Buffered resume meta is not stale | PASS — e2e asserts `.updateMeta` does not contain `已暂停` after applying the buffered update. |
| Live dot and toggle state are consistent | PASS — e2e asserts pulsing live dot and `aria-pressed="true"` after each resume path. |
| Manual refresh feedback unchanged | PASS — B3 test still asserts `刷新中` then `已刷新`. |
| Polling/no reload unchanged | PASS — B3 test still asserts the polling update occurs without `beforeunload`. |
| Pure frontend | PASS — only `frontend/src/pages/WorkOverview.tsx`, `frontend/e2e/overview.spec.ts`, planning docs, and this evidence file changed. |
| 0 backend/shared/root-api diff | PASS — diff/status scans for `services`, `shared`, `api`, and `backend` returned no files. |
| No new color | PASS — no CSS/design-system files touched; added-code color literal scan returned no matches. |
| zh-CN preserved | PASS — all added user-visible assertions/strings are zh-CN. |
| Enum/role comparison values unchanged | PASS — no enum/role logic touched. |
| E2E not weakened / 0 skip | PASS — skip/only/slow scan returned no matches; full suite passed 52/52 with 0 skipped. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npm run test:e2e -- --grep TEST-UIUX-B3-001`  
  Result: PASS, exit 0; 1 passed, 5.4s.
- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built `index-TNHK1XAP.css` and `index-BxWQVDfj.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 52 passed, 0 skipped, 28.2s; Playwright used 4 workers.
- `rg -n "\b(test|describe)\.(only|skip|slow)|test\.skip|describe\.skip|test\.only|describe\.only" frontend/e2e frontend/playwright.config.ts`  
  Result: PASS, no matches.
- `git diff --check`  
  Result: PASS, no whitespace errors.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- docs/ux-ui frontend/src/styles frontend/src/components/ui`  
  Result: PASS, no output.
- `git diff -- frontend/src/pages/WorkOverview.tsx frontend/e2e/overview.spec.ts | rg "^\+[^+].*(#[0-9A-Fa-f]{3,8}|rgba?\(|hsla?\(|color-mix\()"`  
  Result: PASS, no matches.

## Full E2E Tail

```text
> crm-system-frontend@0.1.0 test:e2e
> playwright test

Running 52 tests using 4 workers
...
  ✓  37 [chromium] › e2e/overview.spec.ts:140:1 › TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload (5.7s)
...
  ✓  52 [chromium] › e2e/overview.spec.ts:391:1 › TEST-UIUX-A5-001 main navigation is keyboard reachable (310ms)

  52 passed (28.2s)
```

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `76f2e464ad235cb617f888e1a423e1b6565dc2e1f2cc8adea1ca74c797b069fa` |
| `frontend/dist/assets/index-TNHK1XAP.css` | `d76afb49e74c93175ee45e9ba92228fbbc3732faaffe58143d11941f84b92b8e` |
| `frontend/dist/assets/index-BxWQVDfj.js` | `7f13cfbbb558947556a18e2e0fc559ac65eb7f0f60bd46490d47de1936cf16ae` |

## Handoff

`BLK-UIUX-G12-016` is moved to In Review for Claude. Codex does not self-resolve the blocker and does not self-pass UI/UX G12.
