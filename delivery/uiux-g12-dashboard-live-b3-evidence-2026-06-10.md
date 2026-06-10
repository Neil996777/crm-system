# UI/UX G12 Dashboard Live B3 Evidence — BLK-UIUX-G12-009

Date: 2026-06-10  
Owner: Codex (UX design + frontend execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-009`
- Change type: release-owner design revision + frontend implementation for the locked B3 live layer.
- Surface: dashboard/workbench live header and dashboard live rows.
- Out of scope: backend/API/shared changes, SSE/server push, role/enum semantics, BLK-UIUX-G12-006 hero motion, BLK-UIUX-G12-007 selector rail, BLK-UIUX-G12-008 single-return focus header.

## Spec Revision

- `docs/ux-ui/interaction-spec.md`
  - Added `DEC-UX-LIVE-04` to the accepted Part B decision set.
  - B3 now records the current dashboard live transition mechanism as client-side polling of existing GET endpoints.
  - SSE/server-push remains the parked future path; this patch does not alter the future SSE interaction contract.
  - DEC-UX-LIVE-04 ties polling to the existing B3 patch-in-place model, DEC-UX-LIVE-03 coalescing, change-driven `arrived`, and the `实时更新` / `暂停` pause buffer.

## Implementation

- `frontend/src/pages/WorkOverview.tsx`
  - Split the old one-shot `refresh()` path into `fetchDashboardSnapshot()`, initial load, manual refresh, and live polling application.
  - Polls the same existing GET endpoints already used by dashboard refresh on a default 10s cadence.
  - Coalesces live updates through a 900ms client-side queue.
  - Applies snapshots in place via React state; no navigation or full page reload is used.
  - Keeps `activeCard`, Card -> Focus stage, selector rail, focus restoration, scroll/focus behavior, and user selection state outside the live snapshot replacement.
  - Adds `实时更新` / `暂停` via the existing `LiveToggle` primitive.
  - Paused mode keeps polling but buffers changed snapshots, shows `有 N 条新更新 · 点击刷新`, and applies the latest buffered batch when the pill is clicked or live mode resumes.
  - Manual refresh now shows `刷新中`, then a visible zh-CN confirmation/timestamp.
  - Removed the dead `本月` dashboard button because current dashboard endpoints do not expose a frontend-only period selector.
  - Removed static `实时更新 / 自动合并` spans from the rendered dashboard header.
  - Replaced hardcoded `index === 0` `arrived` with real diff-driven change keys for payment and activity rows. Value changes are counted for live buffering, while row highlights remain row-specific.
- `frontend/src/components/ui/index.tsx`
  - `LiveToggle` now labels the active state as `实时更新` and paused state as `暂停`.
  - Paused live dot is static and uses the existing paused class path.
- `frontend/src/styles/design-system.css`
  - Added `liveBufferPill` using existing tokens only.
  - Added `liveDot.paused` using existing `--subtle`.
  - Added `motion-base` background/border transitions for arrived rows.
- `frontend/src/styles/motion.css`
  - Added `spinIcon` animation for the refresh loading state.

## E2E Coverage

- `frontend/e2e/overview.spec.ts`
  - Added `TEST-UIUX-B3-001 dashboard live polling applies buffered updates without full reload`.
  - Test uses Playwright route interception on existing `/api/activities` GET responses to simulate a live data change without backend changes.
  - Test shortens only the test page's polling/coalescing timers through `addInitScript`; production defaults remain 10s / 900ms.
  - Asserts:
    - `本月` dashboard header button is gone.
    - Static `自动合并` text is gone.
    - `实时更新` is a real button with `aria-pressed="true"`.
    - Polling applies a new activity row without full page reload.
    - `.event.arrived` appears only on the changed row.
    - Pausing switches the control to `暂停`, buffers incoming changes, and shows `有 N 条新更新 · 点击刷新`.
    - Paused live dot uses the neutral paused state and does not pulse.
    - While paused, buffered content is not applied to the visible dashboard.
    - Resuming applies the buffered update with an arrived highlight.
    - Manual refresh shows `刷新中` and then a visible `已刷新` confirmation.
  - Existing A6/A7/focus tests remain; A7 reduced-motion assertion now waits for the required reduced animation event instead of reading immediately.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Client-side polling | PASS — WorkOverview polls existing GET endpoint set; no backend/API/shared changes. |
| DEC-UX-LIVE-04 | PASS — interaction spec records polling as the interim B3 mechanism, SSE parked. |
| Coalescing | PASS — live queue coalesces for 900ms, aligned with DEC-UX-LIVE-03. |
| Patch-in-place / no reload | PASS — snapshot state is applied in place; e2e verifies no `beforeunload` during polling update. |
| Preserve Card -> Focus and rail | PASS — `activeCard`/motion/rail state paths unchanged; A7/focus tests still pass. |
| Change-driven arrived | PASS — hardcoded `index===0` removed; row keys are compared before applying `.arrived`. |
| Pause buffer | PASS — paused mode buffers latest snapshot + union of change keys; e2e verifies pill and deferred application. |
| Manual refresh feedback | PASS — button shows `刷新中`, then update meta shows `已刷新`. |
| Remove dead controls | PASS — dashboard `本月` button and static `自动合并` live segment are no longer rendered. |
| Pure frontend | PASS — docs, frontend React/CSS, and e2e only. |
| 0 backend/shared/root-api diff | PASS — `git diff --name-only -- services shared api backend` returned no files; status scan returned no files. |
| No new color | PASS — color-literal diff scan returned no matches; new UI uses existing tokens. |
| zh-CN preserved | PASS — all new visible strings are zh-CN. |
| Enum/role comparison values unchanged | PASS — no enum/role comparison logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 50/50 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npx playwright test e2e/overview.spec.ts`  
  Result: PASS, exit 0; 7 passed, 10.3s.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built
  `index-93NodPN2.css` and `index-D3zHi2KO.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 50 passed, 0 skipped, 18.4s.
- `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --check`  
  Result: PASS, no whitespace errors.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- docs/ux-ui/interaction-spec.md frontend/src/styles/design-system.css frontend/src/styles/motion.css frontend/src/components/ui/index.tsx frontend/src/pages/WorkOverview.tsx | rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}|rgba?\\(|hsla?\\(|color-mix\\()"`  
  Result: PASS, no matches.
- `rg -n "自动合并|<Button>本月</Button>|liveSegment|index === 0 \\? 'paymentRow arrived'|index === 0 \\? 'event arrived'" frontend/src/pages/WorkOverview.tsx frontend/src/components/ui/index.tsx frontend/e2e/overview.spec.ts docs/ux-ui/interaction-spec.md`  
  Result: PASS — only the e2e negative assertion for `自动合并` remains.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `761273b8e7d6ded6188b32e94b557f9bb33e2ed90063974b7e5db405a3e74420` |
| `frontend/dist/assets/index-93NodPN2.css` | `8ed66884cc04f7c6bf97beb97f46823db6ca58f26f95976dd109a44b58be29b5` |
| `frontend/dist/assets/index-D3zHi2KO.js` | `b3ccacb3bef39a10f84a687a0f83f7bdcf362c359228859a31db58a644a3c7b6` |

## Handoff

`BLK-UIUX-G12-009` is moved to In Review for Claude. Codex does not self-resolve the blocker and does not self-pass UI/UX G12.
