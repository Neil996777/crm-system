# UI/UX G12 Focus Rail Selector Evidence

Date: 2026-06-10
From: Codex (UX design + frontend execution)
To: Claude (G12 focus-rail re-verification audit)
Status: `BLK-UIUX-G12-007` returned for review. Codex does not self-pass G12.

## Scope

Blocker: `BLK-UIUX-G12-007` - release-owner design revision for focus mode's
right rail. The approved BLK-UIUX-G12-006 hero enter/exit transition is kept
unchanged; this patch changes the in-focus rail composition and switch model.

Changed files:

- `docs/ux-ui/interaction-spec.md`
- `docs/ux-ui/design-system.md`
- `frontend/src/pages/WorkOverview.tsx`
- `frontend/src/components/ui/index.tsx`
- `frontend/src/styles/design-system.css`
- `frontend/e2e/overview.spec.ts`
- `planning/blockers.md`
- `planning/gate-status.md`

## Spec Revision

UX/design revision completed before implementation:

- `interaction-spec.md` Part B B2 now defines the right rail as a persistent
  focus selector rail, not a re-composing "7 non-active cards" strip.
- Added `DEC-UX-FOCUSRAIL-01` under Part B Decisions, citing release-owner
  direction recorded in `planning/blockers.md` `BLK-UIUX-G12-007`.
- `design-system.md` sections 6.5, 7.8, and 10 were aligned so visual language no
  longer says "the other panels collapse" as the resting in-focus model.

Decision summary:

- The locked manager dashboard focus rail lists all 8 dashboard cards, including
  the currently-focused one.
- The rail order is stable; focus switching does not add, remove, or reorder rail
  items.
- The current item uses an existing-token selected treatment and
  `aria-current="true"`.
- Click/Enter/Space on a rail item switches only the left stage via the existing
  content crossfade.
- Role-scoped variants list their full authorized card set and do not introduce
  unauthorized manager-only cards.

## Implementation Summary

- `WorkOverview.tsx` now maps the full `cards` array into `sideCards` instead of
  filtering out the active card.
- Each side card receives `selected: card.key === active.key`; `switchFocus()`
  still changes only the active stage and preserves the existing
  `dashboardStageSwitch` crossfade path.
- `FocusStage` changed the right rail aria label to `看板选择器` and renders
  `aria-current="true"` on the selected item.
- CSS strengthens `.sideCard.selected` with existing tokens only:
  `--primary`, `--primary-soft`, and `--badge-shadow`.
- Reduced-motion behavior remains unchanged: the existing `data-motion-mode` and
  `.dashboardMotionReduced` path still applies, with no travel animation.

## Acceptance Mapping

| Acceptance point | Evidence |
|---|---|
| Right rail lists all 8 manager dashboard cards including active | `TEST-UIUX-DASHBOARD-001` and `TEST-UIUX-A7-001` assert `看板选择器` has 8 `.sideCard` items in manager focus mode. |
| Rail order is stable | E2E records `data-focus-side-card` order and asserts it equals `['funnel','stage','trend','leaderboard','todo','payments','key-opportunities','activity']`; after switching focus it asserts the same order remains. |
| Selected state + `aria-current` | E2E asserts exactly one `[aria-current="true"]` in the selector rail and verifies the selected key moves from `funnel`/`payments` to `stage` after switch. |
| Switch changes only the left stage | E2E clicks the `商机阶段构成` rail item, asserts the stage heading changes, then asserts the rail count/order remains unchanged. Existing animation recorder still asserts `dashboardStageSwitch` and no `dashboardStageExit` on switch. |
| Hero enter/exit unchanged | BLK-UIUX-G12-006 motion classes/keyframes are retained; A7 e2e still asserts `dashboardStageEnter`, `dashboardStripEnter`, `dashboardStageExit`, and `dashboardStripExit`. |
| Keyboard/reduced-motion retained | Whole-card keyboard entry remains covered; reduced-motion focus still asserts `data-motion-mode="reduced"`, 8 manager selector items, selected item, `transform:none`, and no travel animations. |
| Role constraints preserved | Sales dashboard remains role-scoped: e2e asserts no manager-only cards and the selector contains the sales user's full authorized six-card set. No enum/role values changed. |

## Constraint Check

- Frontend/design-only: yes.
- Backend/shared/root-api diff: 0. `git diff --name-only -- services shared api backend`
  and `git status --short -- services shared api backend` returned no paths.
- New colors: none. Diff color scan over `docs/ux-ui/design-system.md` and
  `frontend/src/styles/design-system.css` returned no added hex/rgb/hsl/color-mix
  literals.
- zh-CN: preserved. New visible/aria label is Chinese (`看板选择器`).
- Enum/role comparison values: unchanged. Sales role scoping remains in place.
- E2E no weakening/no skip: assertions were strengthened for rail count/order
  and selected state. The skip/only scan returned no matches.
- No backend-needed aggregation or new API behavior: none added.

## Verification

Commands run from `frontend/` on 2026-06-10:

| Command | Result |
|---|---|
| `npx tsc --noEmit` | PASS, exit 0. |
| `npx playwright test e2e/overview.spec.ts` | PASS, 6/6. |
| `npm run build` | PASS. Vite transformed 1633 modules and produced `dist/assets/index-CrP9n7t6.css` and `dist/assets/index-BD4d9BTF.js`. |
| `npm run test:e2e` | PASS. Playwright ran 49 tests using 5 workers: `49 passed (14.8s)`, 0 skipped. |

Additional read-only checks from repo root:

| Command | Result |
|---|---|
| `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e` | No matches. |
| `git diff --name-only -- services shared api backend` | No output. |
| `git status --short -- services shared api backend` | No output. |
| `git diff --check` | PASS, no output. |
| `git diff -- docs/ux-ui/design-system.md frontend/src/styles/design-system.css \| rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}\|rgba?\\(\|hsla?\\(\|color-mix\\()"` | No matches. |
| `rg -n "折叠卡片\|card\\.key !== active\|toHaveCount\\(7\\)\|toHaveCount\\(5\\)" frontend/src/pages/WorkOverview.tsx frontend/src/components/ui/index.tsx frontend/e2e/overview.spec.ts docs/ux-ui/interaction-spec.md docs/ux-ui/design-system.md` | No matches. |

## Build Artifact Hashes

SHA-256 after the 2026-06-10 `npm run build`:

| Artifact | SHA-256 |
|---|---|
| `frontend/dist/index.html` | `11f9599efcd2cec5ea481cd295fddf5f9d6eac6e68b52306263987cd8abf0797` |
| `frontend/dist/assets/index-CrP9n7t6.css` | `adb10b3113eda7f8d7bc59cd65f3328188eca6450040cf1abb88573f8c06c843` |
| `frontend/dist/assets/index-BD4d9BTF.js` | `2958c453bed0ffc9182544722d1f2ca8224bf382c3c133c4c24f24384927419c` |

## Handoff

Returned to Claude for independent BLK-UIUX-G12-007 re-verification:

- Spec/DEC consistency with release-owner Option A.
- Live focus rail behavior: manager rail has persistent 8 items, stable order,
  and selected state transfer.
- Strengthened e2e coverage and cumulative constraints.
- Codex does not self-resolve the blocker and does not self-pass G12.
