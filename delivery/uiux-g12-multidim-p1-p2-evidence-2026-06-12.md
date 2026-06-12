# UI/UX G12 Multi-Dim P1/P2 Evidence — 2026-06-12

Scope: BLK-UIUX-G12-019 + BLK-UIUX-G12-020. Frontend-only repair package for the multi-dimensional audit findings in `delivery/uiux-multidim-audit-findings-2026-06-12.md` against `docs/ux-ui/requirements/uiux-audit-matrix.md`.

Status: returned to Claude for independent functional, forced-failure, responsive, and copy re-audit. Codex does not self-resolve.

## P1

| Item | Implementation | Evidence |
|---|---|---|
| F1 list `refresh()` / `select*()` errors are not silent | Leads, accounts, contacts, opportunities, quotes, contracts, and payments now wrap list refresh and record-open calls in `try/catch`, call `setError(localizeError(...))`, and render the list-level error via `CrudListShell.feedback`. `AccountDetail.refreshContacts()` now catches `listContacts` errors and shows the passed error alert. | `frontend/src/pages/{leads,accounts,opportunities,quotes,contracts,payments}/*List.tsx`; `frontend/src/pages/accounts/AccountDetail.tsx`; `frontend/src/components/CrudScaffold.tsx`. E2E `TEST-UIUX-G12-019 list refresh and record-open failures surface errors` forces `/api/leads` and `/api/leads/{id}` failures and verifies `请求失败。` appears. |
| F2 list async loading state | Seven list pages now track `loading` and `selectingId`; refresh/apply buttons show busy/disabled state, list tables render `ListTableLoading`, and opening a record shows `正在打开记录...`. | `ListAsyncFeedback` and `ListTableLoading` in `CrudScaffold`; seven list pages. |
| F3 QuoteDetail / ContractDetail errors rendered | Verified existing detail components already render passed `error` with `role="alert"`; parent list selection/update error paths now keep feeding localized errors. | `frontend/src/pages/quotes/QuoteDetail.tsx`, `frontend/src/pages/contracts/ContractDetail.tsx`. |
| F4 CloseOpportunityDialog required gates | Submit is disabled unless Won has `contractId`, or Lost has both `reasonCode` and `reasonDetail`; submit handler also returns early if invalid. | `frontend/src/components/CloseOpportunityDialog.tsx`; e2e strengthened in `opportunities.spec.ts` for Won and Lost disabled/enabled transitions. |
| F5 Import no-file disabled | `开始导入` is disabled with `disabled={busy || !file}` until a CSV file is selected. | `frontend/src/pages/importexport/Import.tsx`; `TEST-CSV-IMPORT-001/002` asserts disabled before file and enabled after file. |
| F6 BasicReports fake buttons removed | Removed `本月` and `按负责人分组` because the current frontend has only aggregate report data and no effective frontend dimension to perform a truthful client-side month filter. No notice-only fake control remains. | `frontend/src/pages/reports/BasicReports.tsx`; `reports.spec.ts` asserts both buttons are absent. |
| F7 OpportunityList row menu no fake inline transfer | Row menu `转移负责人` now performs a real `updateOpportunity` PATCH with prompted owner ID and existing versioned payload. It no longer only opens detail while promising an inline action. | `frontend/src/pages/opportunities/OpportunityList.tsx`. Existing role/terminal/archive gates preserved. |
| Remove no-backend disabled bulk buttons | Deleted hardcoded disabled buttons with titles such as `...无接口；按 A3 禁用` from PaymentList, TaskList, AccountList, ContactList, OpportunityList, QuoteList, and ContractList. Real actions such as lead/account/contract/opportunity archive, task complete, export, and clear selection remain. | `frontend/e2e/list-actions.spec.ts` asserts no `button[title*="按 A3 禁用"]` and page-specific absent fake bulk buttons. |

## P2

| Item | Implementation | Evidence |
|---|---|---|
| P2-1 PaymentList title identity | Payment record title now displays `contract.id`; aria-label and click target also use `contract.id`; opportunity id moved to subtitle. | `frontend/src/pages/payments/PaymentList.tsx`; `list-actions.spec.ts` asserts the 回款 title button contains the contract id. |
| P2-3 OperationLogs label | Action filter option changed from `操作人：全部` to `操作：全部`. | `frontend/src/pages/admin/OperationLogs.tsx`; `oplog.spec.ts` asserts the action filter text. |
| P2-4 UserManagement last-admin label | `唯一启用管理员` row label now renders only when `isLastActiveAdministrator(user)` is true. Other administrators show the normal role label. | `frontend/src/pages/admin/UserManagement.tsx`; `user-admin.spec.ts` asserts a new Sales row does not show the label and the seed last admin row does. |
| P2-5 Import/Export type badges | Removed non-selectable 客户/联系人/商机/报价/合同 badges; only the selectable 线索 badge is shown. | `frontend/src/pages/importexport/Import.tsx`, `Export.tsx`; import/export e2e assert one object-type badge. |
| P2-6 Export fileSafety mapping | `fileSafetyLabel` now includes the real backend value `dangerous_cells_prefixed`, safe variants, `none_required`, and a mapped `unknown` fallback; `Export.tsx` no longer hardcodes the fallback string locally. | `frontend/src/i18n/labels.ts`, `frontend/src/pages/importexport/Export.tsx`. |
| P2-9 Export action row containment | Export form grid changed from four forced columns to two constrained columns; the confirmation row spans the grid, preventing the 1280px right-edge overflow. | `frontend/src/styles/design-system.css`. |

Optional P2-7/P2-8 were not taken: extracting status tone helpers and column reordering would be broader refactors with no functional need for this release-blocking return.

## Verification

- `npx tsc --noEmit` — PASS.
- Targeted e2e: `npx playwright test e2e/list-actions.spec.ts e2e/import.spec.ts e2e/export.spec.ts e2e/reports.spec.ts e2e/opportunities.spec.ts e2e/oplog.spec.ts e2e/user-admin.spec.ts` — PASS, 20/20.
- `npm run build` — PASS; Vite output `dist/assets/index-CnyUMR_S.css`, `dist/assets/index-CGA8N7yU.js`.
- Full e2e: `npm run test:e2e` — PASS, 58/58, 0 failed, workers:2, retries:1, no flaky retry output.
- `git diff --check` — PASS.
- `rg -n "test\\.(skip|only|slow)" frontend/e2e frontend/playwright.config.ts` — no matches.
- `git diff --name-only -- services shared api packages/shared apps/api` — no output.
- `git status --short -- services shared api packages/shared apps/api` — no output.
- Added color literal scan: `git diff -U0 -- frontend/src frontend/e2e | rg -n "^\\+.*(#[0-9A-Fa-f]{3,8}|rgba?\\()"` — no matches.

## Handoff

Pending Claude re-audit: forced refresh/select failures, each touched control, required field gates, removed fake buttons, P2 copy/display/layout checks, and any desired elementFromPoint checks. Codex does not self-resolve BLK-UIUX-G12-019 or BLK-UIUX-G12-020.
