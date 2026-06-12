# UI/UX G12 Row Menu Portal Evidence - 2026-06-12

Status: Codex return for Claude re-audit. Codex does not self-resolve
`BLK-UIUX-G12-018`.

Scope: frontend-only fix for Post-G12 spot-fix #14: row action menus were opened
inside the scrollable table and could be covered/clipped by later rows; the
separate `->` quick-view buttons also duplicated `查看`.

## Implementation

Shared menu primitive:

- `frontend/src/components/ui/index.tsx` `ActionMenu` now renders the open
  `role="menu"` through `createPortal(..., document.body)`.
- The menu panel is `position: fixed`, anchored to the trigger button via
  `getBoundingClientRect()`, and recalculates on window resize and captured
  scroll events.
- The panel is constrained to the viewport with `maxHeight`, keeps outside-click
  close, Escape close, disabled menu item reasons, and click propagation guards.
- `frontend/src/styles/design-system.css` changed `.actionMenuPanel` from table
  local `absolute/z-index:90` to body-level `fixed/z-index:1000`, with existing
  color tokens only.

Clickable record names:

- `frontend/src/components/CrudScaffold.tsx` `RecordIdentity` now accepts
  `onTitleClick` and `titleAriaLabel`.
- The primary title renders as a keyboard-focusable `.recordLinkButton` when a
  click handler is supplied; it uses existing `--primary` focus/hover styling and
  `data-row-interactive="true"` so table row click does not double-fire.

Seven list pages:

| List | Removed `->` button | New title control | Menu retained |
|---|---:|---|---|
| 线索 | Yes | `打开线索 {name}` -> `getLead(id)` detail | Yes: `查看` / `转移负责人` / `归档` / `转为商机`, existing lead endpoints and existing gates. |
| 公司/客户 | Yes | `打开客户 {companyName}` -> `getAccount(id)` detail | Yes: `查看` / `归档`, existing account endpoints and Sales/archive gates. |
| 联系人 | Yes | `打开联系人 {contactName}` -> `getContact(id)` detail | No menu; there is no non-view row action endpoint, so no dead menu is reintroduced. |
| 商机 | Yes | `打开商机 {title}` -> `getOpportunity(id)` detail | Yes: `查看` / `编辑` / `推进阶段` / `转移负责人` / `归档`, existing opportunity endpoints and role/terminal/archive gates. |
| 报价 | Yes | `打开报价 {quoteId}` -> `getQuote(id)` detail | Yes: `查看` / `发送` / `接受` / `拒绝` / `标记过期`, existing quote status endpoint and status gates. |
| 合同 | Yes | `打开合同 {contractId}` -> `getContract(id)` detail | Yes: `查看` / `签署` / `启用` / `完成` / `终止` / `归档`, existing contract status/archive endpoints and gates. |
| 回款 | Yes | `打开回款合同 {contractId}` -> payment/contract detail | Yes: `查看` / `新建计划` / `登记回款`, existing contract/payment-plan/payment detail flows. |

## E2E Evidence

New file: `frontend/e2e/list-actions.spec.ts`.

- `TEST-UIUX-G12-018 row action menu portal stays topmost and clickable on first middle last rows`
  creates 30 leads, filters the list, opens the first/middle/last row menus,
  scrolls rows into view, asserts `elementFromPoint` at the menu center is inside
  the `role="menu"` tree instead of a table row/`.rowAction`, and clicks `查看`
  to prove the menu item executes.
- `TEST-UIUX-G12-018 seven lists use clickable record names and no legacy arrow view buttons`
  creates one record path covering lead/account/contact/opportunity/quote/contract/payment,
  asserts each filtered row has no legacy `查看*` arrow button, checks body-portal
  menu topmost behavior plus `查看` execution on the six menu-bearing lists, and
  clicks each record-name control to open the matching detail surface.

Targeted run:

| Command | Result |
|---|---|
| `cd frontend && npx playwright test e2e/list-actions.spec.ts` | PASS, 2/2. |

Full run:

| Command | Result |
|---|---|
| `cd frontend && npx tsc --noEmit` | PASS, clean. |
| `cd frontend && npm run build` | PASS (`tsc -b && vite build`; assets `index-BdXpNNve.css`, `index-KoAr-C7O.js`). |
| `cd frontend && npm run test:e2e` | PASS, 56/56, 0 failed, 0 skipped, `workers: 2`, `retries: 1`; no flaky retry reported. |

## Constraint Check

- Frontend-only: implementation touched frontend source/e2e plus evidence/planning
  docs.
- Backend/shared/root API diff: empty by
  `git diff --name-only -- services shared api packages/shared apps/api` and
  `git status --short -- services shared api packages/shared apps/api`.
- Existing endpoints only: row menu actions still call existing frontend API
  clients and existing detail/form flows; no backend route or shared contract was
  added.
- No new colors: `git diff -U0` over touched frontend files has no added hex or
  rgba literals; CSS additions use existing tokens.
- zh-CN preserved: all new labels/aria names/evidence terms are Chinese.
- Enum/role comparison values unchanged: no opportunity stage, quote/contract
  status, or role comparison values were changed.
- Role gating not widened: Sales/archive/terminal/status-disabled menu gates are
  unchanged; unsupported contact row actions remain absent rather than faked.
- e2e hygiene: `rg -n 'test\.(skip|only|slow)' frontend/e2e frontend/playwright.config.ts`
  has no matches.
- `git diff --check`: PASS.

Pending: Claude re-audit with live `elementFromPoint`/topmost/clipping checks and
manual click verification. Codex does not self-resolve `BLK-UIUX-G12-018`.
