# UI/UX G12 Focus Header Evidence — BLK-UIUX-G12-008

Date: 2026-06-10  
Owner: Codex (UX design + frontend execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-008`
- Change type: release-owner design revision superseding the visible Esc hint chip in the locked focus mockup.
- Surface: dashboard Card -> Focus stage header only.
- Out of scope: backend/API/shared changes, role/enum semantics, BLK-UIUX-G12-006 hero motion, BLK-UIUX-G12-007 persistent selector rail.

## Spec Revision

- `docs/ux-ui/interaction-spec.md`
  - Added `DEC-UX-FOCUSEXIT-01` to the accepted Part B decision set.
  - B2 `Entry/exit + keyboard` now requires exactly one visible focus-stage exit control: the `返回` button with back-chevron icon.
  - Removed the prior binding language that the visible `Esc 返回` hint chip stays visible.
  - Preserved `Esc` as the global keyboard shortcut for exiting focus mode, satisfying B6 keyboard operability.
  - Decision record states this supersedes the mockup's visible `Esc 退出`/`Esc 返回` chip while leaving data-scope badge, hero enter/exit, and selector rail unchanged.
- `docs/ux-ui/design-system.md`
  - Focus layout and Card -> Focus visual sections now describe a single `返回` ghost/secondary button with back-chevron icon.
  - Button/chip section no longer reserves `chip` for a focus-stage Esc hint; generic chips remain available for non-exit metadata.

## Implementation

- `frontend/src/components/ui/index.tsx`
  - `FocusStage` no longer accepts or renders `escapeHint`.
  - Focus-stage `返回` button keeps the existing secondary button treatment and now includes the existing `ChevronLeft` icon.
  - `aria-label={backLabel}` keeps the accessible button name exactly `返回` for the dashboard focus stage.
- `frontend/src/pages/WorkOverview.tsx`
  - Removed the `escapeHint="Esc 返回"` prop.
  - Existing `keydown` Escape handling remains unchanged.
  - Data-scope badge (`全部`/`团队`/`本人`) remains in `tools`.
- No CSS changes were required for this patch; no new color token or color literal was introduced.

## E2E Coverage

- `frontend/e2e/overview.spec.ts`
  - Added `expectSingleFocusExitControl(page)`.
  - Asserts the focus-stage `.stageTools` contains exactly one button.
  - Asserts the single button has accessible name `返回`.
  - Asserts the focus stage contains no visible `Esc` text.
  - Existing Escape-key assertions remain: pressing `Escape` exits focus mode and returns to the dashboard.
  - Existing BLK-UIUX-G12-006/007 assertions remain: hero motion, reduced-motion path, selector count/order, and selected `aria-current` transfer are still covered.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Single visible focus exit control | PASS — focus-stage tools render one `返回` button. |
| Remove visible Esc hint chip | PASS — no `escapeHint` prop exists in `frontend/src`; overview e2e asserts no `Esc` text in the focus stage. |
| Esc key still exits focus mode | PASS — keydown handler unchanged; e2e presses `Escape` and verifies dashboard return. |
| Data-scope badge unchanged | PASS — `tools={<Badge tone="primary">{model.scopeBadge}</Badge>}` unchanged. |
| BLK-UIUX-G12-006 hero motion unchanged | PASS — no motion CSS/state-machine changes; A7 motion test still passes. |
| BLK-UIUX-G12-007 selector rail unchanged | PASS — no rail composition/style changes; persistent rail tests still pass. |
| Pure frontend | PASS — changed docs, React component/page, and e2e only. |
| 0 backend/shared/root-api diff | PASS — `git diff --name-only -- services shared api backend` returned no files; status scan returned no files. |
| No new color | PASS — color-literal diff scan returned no matches. |
| zh-CN preserved | PASS — visible UI text remains Chinese; removed duplicate hint instead of adding copy. |
| Enum/role comparison values unchanged | PASS — no enum/role logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 49/49 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built
  `index-CrP9n7t6.css` and `index-SCNr76DL.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 49 passed, 0 skipped, 11.8s.
- `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --check`  
  Result: PASS, no whitespace errors.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- docs/ux-ui/design-system.md frontend/src/styles/design-system.css frontend/src/components/ui/index.tsx frontend/src/pages/WorkOverview.tsx | rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}|rgba?\\(|hsla?\\(|color-mix\\()"`  
  Result: PASS, no matches.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `2254ac97d84caba3a0fbf5dd4145195523d1bf1922aef9483ec0b84ffafa3451` |
| `frontend/dist/assets/index-CrP9n7t6.css` | `adb10b3113eda7f8d7bc59cd65f3328188eca6450040cf1abb88573f8c06c843` |
| `frontend/dist/assets/index-SCNr76DL.js` | `c3bfa2175a5fbdc5f127510b9408637895a48bab03404bcf457ddd5e7b1f5722` |

## Handoff

`BLK-UIUX-G12-008` is moved to In Review for Claude. Codex does not self-resolve the blocker and does not self-pass UI/UX G12.
