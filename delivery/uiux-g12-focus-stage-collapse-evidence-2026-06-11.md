# UI/UX G12 Focus Stage Collapse Evidence — BLK-UIUX-G12-014

Date: 2026-06-11  
Owner: Codex (frontend execution)  
Status: Returned to Claude for independent re-verification. Codex does not self-pass G12.

## Scope

BLK-UIUX-G12-014 fixes the regression introduced by BLK-UIUX-G12-013: focus-mode `.sidebar` was changed to `position: fixed`, which removed it from `.shell.focusMode` grid flow and allowed the workspace to auto-place into the 72px first column. The focus `.stage` then collapsed to near-zero width while the 300px selector rail overflowed.

Frontend-only files changed:

- `frontend/src/styles/design-system.css`
- `frontend/e2e/overview.spec.ts`

## Implementation

- `.shell.focusMode .sidebar` is back in grid flow with `position: sticky; top: 0`; the 72px shell track is preserved and the workspace stays in the second grid column.
- The focus rail keeps the G12-013 behavior: compact brand, true sr-only brand/nav labels, nav vertical scrolling on short heights, and NAV-01 zh-CN flyout above the stage with the 360ms hover dwell.
- `.shell.focusMode .nav { width: 220px }` remains as the flyout carrier. It no longer affects shell placement because the sidebar stays in the fixed 72px grid track and uses visible overflow.
- Hero motion (G12-006/012), persistent selector rail (G12-007), single-return header (G12-008), B3 polling (G12-009), and dashboard responsive fixes (G12-010/011) were not changed.

## E2E Strengthening

`overview.spec.ts` adds `TEST-UIUX-FOCUS-LAYOUT-001`:

- Checks focus mode at 1280, 1366, and 1440px desktop widths.
- Asserts `.stage` width is greater than 500px.
- Asserts `.stage` is more than 45% of `.focus`.
- Asserts `.focus` does not compute to `0px 300px`.
- Asserts the sidebar is not `position: fixed`.
- Asserts the 300px selector rail remains present.
- Keeps the persistent 8-item selector and selected state checks.

`TEST-UIUX-NAV-01` is retained and now also checks the real stage width before verifying:

- 14 admin nav icons reachable, including `操作日志` on short height.
- NAV-01 hover/focus flyout is topmost via `elementFromPoint`.
- Flyout is not covered by `.stage`.

## Constraint Check

- Pure frontend: yes.
- Backend/shared/root-api diff: none.
- New colors: none; CSS diff adds no color literals.
- zh-CN: preserved; labels continue from existing nav `aria-label` values.
- Enum/role comparison values: unchanged.
- E2E weakening/skips: none; full suite increased to 52 tests, 0 skipped.
- G12 self-resolution: no; returned to Claude.

## Verification

- `npx playwright test e2e/overview.spec.ts -g "TEST-UIUX-FOCUS-LAYOUT-001|TEST-UIUX-NAV-01"`: PASS, 2/2.
- `npx playwright test e2e/overview.spec.ts`: PASS, 9/9.
- `npx tsc --noEmit`: PASS.
- `npm run build`: PASS.
  - Output: `dist/assets/index-TNHK1XAP.css`, `dist/assets/index-CIQmHaIF.js`.
- `npm run test:e2e`: PASS, 52/52, 0 skipped.
- `git diff --check`: PASS.
- `rg -n "\b(test|describe)\.(only|skip)|test\.skip|describe\.skip" frontend/e2e`: no matches.
- `git diff --name-only -- services shared api backend`: no output.
- New-color diff scan on `design-system.css`/`overview.spec.ts`: no matches.

Artifact SHA-256:

- `frontend/dist/index.html`: `379e8bb971a5dfc61dc794124f167f96d41a6321196cc37f79d7f255502a9e5e`
- `frontend/dist/assets/index-TNHK1XAP.css`: `d76afb49e74c93175ee45e9ba92228fbbc3732faaffe58143d11941f84b92b8e`
- `frontend/dist/assets/index-CIQmHaIF.js`: `3e147c5d22e7ee7c91455d17bf8ef2ba93541bc1ae59973587107974c39b7d2d`

## Handoff

BLK-UIUX-G12-014 is returned to Claude for independent re-verification. Per the strengthened audit standard recorded in `planning/blockers.md`, Claude should include full-view focus screenshots at 1280/1366/1440 and confirm the automated stage-width guard.
