# UI/UX G12 Focus Rail Nav Evidence — BLK-UIUX-G12-013

Date: 2026-06-10  
Owner: Codex (frontend execution)  
Status: Returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-013`
- Release-owner finding: in focus mode, the collapsed 72px rail clips the last admin nav item (`操作日志`) and the accepted `DEC-UX-NAV-01` label flyout is hidden under the focus stage.
- Surfaces: shell left navigation only while `shell.focusMode` is active.
- Out of scope: expanded sidebar behavior, dashboard hero timing (`BLK-UIUX-G12-006/012`), persistent focus selector (`BLK-UIUX-G12-007`), single return header (`BLK-UIUX-G12-008`), B3 polling (`BLK-UIUX-G12-009`), and responsive dashboard card fixes (`BLK-UIUX-G12-010/011`).

## Spec Binding

- Existing binding: `docs/ux-ui/interaction-spec.md` `DEC-UX-NAV-01` and `docs/ux-ui/design-system.md` rail flyout section already require the collapsed rail to show a temporary hover/focus label overlay above the focus stage.
- This patch does not introduce a new design decision. It implements the existing NAV-01 requirement and the release-owner `BLK-UIUX-G12-013` acceptance.

## Implementation

- `frontend/src/styles/design-system.css`
  - In focus mode, brand/nav text now becomes true sr-only content instead of transparent layout content, preventing the hidden Chinese subtitle from wrapping vertically and consuming rail height.
  - Focus-mode brand footprint is reduced (`36px` min-height, smaller `CRM` mark, lower padding/margin).
  - Focus-mode sidebar is pinned as a fixed 72px viewport rail, with the grid track still reserving the same 72px column for the workspace.
  - Focus-mode nav uses a compact 4px gap and a scrollable vertical rail region for short viewports; the 44px icon hit area is preserved.
  - NAV-01 flyout stacking is raised above the stage (`sidebar`/`nav`/`navItem`/`::after` z-index path), using existing `--card`, `--border`, `--text`, and `--shadow` tokens only.
  - Hover reveal has a 360ms dwell delay. Keyboard `focus-visible` opens immediately. Reduced-motion removes the slide and delay.
  - Mobile single-column fallback keeps the sidebar static so the focus-only fixed rail does not leak into `<760px` layout.
- `frontend/e2e/overview.spec.ts`
  - Added `TEST-UIUX-NAV-01`.
  - Asserts admin focus mode has all 14 nav items and that `操作日志` is within the viewport at 1440x900.
  - Re-enters focus at 1440x720 and asserts the rail scrolls to make `操作日志` reachable.
  - Hovers `联系人` and focuses `报价`, then checks the `::after` label text, opacity/visibility, and `elementFromPoint` topmost behavior: the point resolves to the nav button and not `.stage`.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Focus-mode brand no longer pushes icons below the fold | PASS — hidden brand text is sr-only, brand footprint is compact, and `TEST-UIUX-NAV-01` verifies all 14 admin nav items with `操作日志` reachable. |
| Short screens scroll instead of clipping | PASS — 1440x720 branch scrolls the focus nav container and verifies the last item is reachable within the viewport. |
| NAV-01 flyout visible above stage | PASS — hover/focus flyout is raised above the focus stage; e2e verifies `elementFromPoint` at the flyout point is the nav button, not `.stage`. |
| Hover delay and focus parity | PASS — hover uses 360ms transition delay; focus-visible opens with 0ms delay; e2e covers both hover and focus. |
| Reduced-motion no slide | PASS — reduced-motion media rule sets the flyout transform to the final position and removes delay. |
| Existing zh-CN labels reused | PASS — flyout text comes from current nav `aria-label` values (`联系人`, `报价`, etc.); no new English UI text. |
| No expanded-sidebar regression | PASS — changes are scoped to `.shell.focusMode`; existing main-navigation keyboard test still passes. |
| Pure frontend | PASS — changed CSS and e2e only, plus delivery/planning records. |
| 0 backend/shared/root-api diff | PASS — backend/shared/root-api diff and status checks returned no files. |
| No new color | PASS — exact color-literal diff scan returned no matches. |
| zh-CN preserved | PASS — all added visible assertions use existing Chinese labels. |
| Enum/role comparison values unchanged | PASS — no enum/role comparison logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 51/51 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx playwright test e2e/overview.spec.ts -g TEST-UIUX-NAV-01`  
  Result: PASS, exit 0; 1 passed, 3.5s.
- `cd frontend && npx playwright test e2e/overview.spec.ts`  
  Result: PASS, exit 0; 8 passed, 13.6s.
- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built `index-B5ZIliaI.css` and `index-Cas16Yno.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 51 passed, 0 skipped, 20.9s.
- `git diff --check`  
  Result: PASS, no output.
- `rg -n "\\b(test|describe)\\.(only|skip)|test\\.skip|describe\\.skip" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- frontend/src/styles/design-system.css frontend/e2e/overview.spec.ts | rg -n "^\\+.*(#[0-9A-Fa-f]{3,8}\\b|rgba?\\(|hsla?\\(|color-mix\\(|linear-gradient\\()"`
  Result: PASS, no matches.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `d53518b37aeed5dbaa63eea02ee3cdd178e44b1ba43351a868a71e9ba1e76c14` |
| `frontend/dist/assets/index-B5ZIliaI.css` | `9684f71609137630110865ef020b649fe1069dcabde04b11f1c90a6aabed0405` |
| `frontend/dist/assets/index-Cas16Yno.js` | `3e147c5d22e7ee7c91455d17bf8ef2ba93541bc1ae59973587107974c39b7d2d` |

## Handoff

BLK-UIUX-G12-013 is returned to Claude for independent re-verification. Codex does not self-resolve or self-pass G12.
