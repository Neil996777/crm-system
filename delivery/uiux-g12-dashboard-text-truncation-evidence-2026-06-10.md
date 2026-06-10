# UI/UX G12 Dashboard Text / Icon Overlap Evidence — BLK-UIUX-G12-011

Date: 2026-06-10  
Owner: Codex (UX design + frontend execution)  
Status: Re-kick returned to Claude for independent G12 re-verification. Codex does not self-pass G12.

## Scope

- Blocker: `BLK-UIUX-G12-011`
- Acceptance basis: `docs/ux-ui/requirements/uiux-implementation.requirements.md` A6 (`桌面优先 + 响应式不破版`) plus Claude re-verification FAIL on 2026-06-10.
- Surface: dashboard/workbench cards in the normal dashboard grid, especially `团队回款到账`, `团队商机阶段构成`, and `团队最近活动`.
- Break band covered: `1440`, `1512`, and `1600` px, plus A6 coverage at `1700`, `1680`, `1280`, `1180`, `1024`, and `900` px.
- Out of scope: backend/API/shared changes, role/enum semantics, Card -> Focus hero motion (`BLK-UIUX-G12-006`), persistent focus selector rail (`BLK-UIUX-G12-007`), single-return focus header (`BLK-UIUX-G12-008`), and B3 live polling (`BLK-UIUX-G12-009`).

## Spec / Decision Note

- No new UX DEC was introduced. This patch is an implementation correction under the existing A6 responsive/no-overflow requirement.
- The first return fixed over-truncation by reflowing readable dashboard text, but Claude re-verification found a new icon/text overlap in the `1024px–1699px` two-column payment/activity rows.
- This re-kick correction keeps the same responsive model and aligns the payment/activity row icon grid track to the real `.flowIcon` width in both the base row layout and the `1024px–1699px` constrained layout.

## Implementation

- `frontend/src/styles/design-system.css`
  - Existing first-return fix remains: at desktop widths `>=1024px`, payment meta, activity title/meta/time, and dashboard row title/meta use visible wrapping instead of ellipsis.
  - Existing first-return fix remains: in the constrained `1024px–1699px` band, dashboard donut cards stack the legend below the ring so all six stage labels fit in full.
  - Re-kick fix: `.paymentRow` and `.event` now use a 42px icon track in the base three-column row layout and in the constrained two-column `1024px–1699px` layout, so the 42px `.flowIcon` no longer overflows a 22px/24px track. The sales dashboard override is also corrected to 42px.
  - The 22px legend swatch track is unchanged; no color tokens or literals were added.
- `frontend/e2e/overview.spec.ts`
  - `TEST-UIUX-A6-001` now covers `1700/1680/1600/1512/1440/1280/1180/1024/900`.
  - Existing no-truncation checks remain for payment row meta, payment amount/status, donut legend labels, and activity title/meta/time at all tested widths `>=1024px`.
  - New `expectDashboardFlowRowsDoNotOverlap` asserts that `.flowIcon` bounding boxes do not intersect any same-row text/content bounding box in payment rows or activity rows. This covers the requested `1440/1512/1600` widths and the rest of the desktop A6 set.

## Acceptance / Constraint Check

| Requirement | Result |
|---|---|
| Payment row meta, amount, and status render in full at desktop widths | PASS — e2e checks `.paymentRow small`, `.paymentRight .money`, and `.paymentRight .badge` for no horizontal truncation at `1700/1680/1600/1512/1440/1280/1180/1024`. |
| Donut legend labels render in full at desktop widths | PASS — legend stacks below the donut through the constrained band, and e2e checks all stage legend labels for no truncation. |
| Activity title/meta/time render without clipping | PASS — e2e checks `.event p`, `.event small`, and `.eventTime` for no horizontal truncation at the same desktop widths. |
| Flow icons do not overlap payment/activity text | PASS — e2e compares `.flowIcon` bounding boxes against same-row text/content boxes for payment and activity rows, including `1440/1512/1600` and the 1700px base-layout boundary. |
| G12-010 responsive degradation preserved | PASS — existing grid breakpoints and no-overflow checks remain in `TEST-UIUX-A6-001`. |
| Hero/focus/live behavior unaffected | PASS — overview e2e still passes dashboard focus, A7, and B3 live tests. |
| Pure frontend | PASS — changed CSS and e2e only, plus delivery/planning records. |
| 0 backend/shared/root-api diff | PASS — backend/shared/root-api diff and status checks returned no files. |
| No new color | PASS — CSS/e2e diff color-literal scan returned no matches. |
| zh-CN preserved | PASS — no new visible English UI text added. |
| Enum/role comparison values unchanged | PASS — no enum/role comparison logic touched. |
| E2E not weakened / 0 skip | PASS — full suite 50/50 passed; skip/only scan returned no matches. |

## Verification Commands

Run from `/Users/neil/practice/software/projects/crm-system`.

- `cd frontend && npx playwright test e2e/overview.spec.ts`  
  Result: PASS, exit 0; 7 passed, 10.3s.
- `cd frontend && npx tsc --noEmit`  
  Result: PASS, exit 0.
- `cd frontend && npm run build`  
  Result: PASS, exit 0; Vite transformed 1633 modules; built `index-Dw0l66tN.css` and `index-Co62K16j.js`.
- `cd frontend && npm run test:e2e`  
  Result: PASS, exit 0; 50 passed, 0 skipped, 19.4s.
- `git diff --check`  
  Result: PASS, no whitespace errors.
- `rg -n "\\b(test|describe)\\.(only|skip)|test\\.skip|describe\\.skip" frontend/e2e`  
  Result: PASS, no matches.
- `git diff --name-only -- services shared api backend`  
  Result: PASS, no output.
- `git status --short -- services shared api backend`  
  Result: PASS, no output.
- `git diff -- frontend/src/styles/design-system.css frontend/e2e/overview.spec.ts | rg -n "^\\+.*(#|rgba?\\(|hsla?\\(|color-mix\\(|linear-gradient\\()"`
  Result: PASS, no matches.

## Build Artifact SHA-256

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `e65a70a71c53fb4687139047c1474b7ad3225b5f1339489ce55115228ca3d971` |
| `frontend/dist/assets/index-Dw0l66tN.css` | `6a3c73316433481678e1cd08b1257e7790c18806de016b46b51d84580a0aed10` |
| `frontend/dist/assets/index-Co62K16j.js` | `b3ccacb3bef39a10f84a687a0f83f7bdcf362c359228859a31db58a644a3c7b6` |

## Handoff

BLK-UIUX-G12-011 is returned to Claude for independent re-verification after the re-kick. Codex does not self-resolve or self-pass G12.
