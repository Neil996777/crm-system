# UI/UX G12 Functional Completeness Evidence - 2026-06-11

Status: Codex return for Claude functional re-verification. Codex does not
self-resolve `BLK-UIUX-G12-017`.

Scope: functional-completeness polish across non-dashboard pages. This evidence
records the already-landed implementation and e2e verification from the previous
Codex pass. This docs-only return did not change implementation, CSS, or e2e.

## A. Row Action Menus

Shared primitive: `frontend/src/components/ui/index.tsx` now provides
`ActionMenu` with a real trigger button (`aria-haspopup="menu"`), `role="menu"`,
keyboard Escape close, outside-click close, disabled menu items with reasons, and
no-op prevention for disabled entries. `DataTable` row clicks ignore interactive
targets, so menu and arrow buttons do not accidentally trigger the row open path.

| List page | Row menu/actions | Role and state gates | Existing endpoint/use |
|---|---|---|---|
| 线索 `LeadList` | `查看`; `转移负责人`; `归档`; `转为商机` | `转移负责人` and `归档` disabled for `Sales` and archived records. `转为商机` disabled unless status is `Valid`, owner exists, and record is not archived; it opens detail for the existing conversion form. | `GET /api/leads/{id}`; `POST /api/leads/{id}/owner-transfer`; `POST /api/leads/{id}/archive`; conversion remains detail flow via `POST /api/leads/{id}/convert`. |
| 公司/客户 `AccountList` | `查看`; `归档` | `归档` disabled for `Sales` and archived records. Owner transfer remains disabled-with-reason because no customer owner-transfer endpoint exists. | `GET /api/accounts/{id}`; `POST /api/accounts/{id}/archive`. |
| 联系人 `ContactList` | No `···` menu is rendered. The former dead menu was removed because the only available row action is `查看`; arrow button and row click open detail. | Bulk transfer/archive stay disabled-with-reason because no contact transfer/archive endpoint exists. | `GET /api/contacts/{id}` via `getContact`. |
| 商机 `OpportunityList` | `查看`; `编辑`; `推进阶段`; `转移负责人`; `归档` | `编辑`, `转移负责人`, and `归档` disabled for `Sales`, terminal `Won`/`Lost`, and archived records. `推进阶段` disabled for terminal or archived records. Terminal records remain read-only. | `GET /api/opportunities/{id}`; edit/transfer detail uses `PATCH /api/opportunities/{id}`; stage uses `POST /api/opportunities/{id}/stage`; archive uses `POST /api/opportunities/{id}/archive`. |
| 报价 `QuoteList` | `查看`; `发送`; `接受`; `拒绝`; `标记过期` | Lifecycle gates are status-based: `Draft -> Sent`, `Sent -> Accepted/Rejected`, `Draft/Sent -> Expired`. Non-eligible actions are disabled with zh-CN reasons. | `GET /api/quotes/{id}`; `POST /api/quotes/{id}/status`. |
| 合同 `ContractList` | `查看`; `签署`; `启用`; `完成`; `终止`; `归档` | `签署` only from `Pending Signature`; `启用` only from `Signed`; `完成` only from `Active`; `终止` disabled for `Completed`/`Terminated`; `归档` disabled for `Sales` and archived records. | `GET /api/contracts/{id}`; `POST /api/contracts/{id}/status`; `POST /api/contracts/{id}/archive`. |
| 回款 `PaymentList` | `查看`; `新建计划`; `登记回款` | Row is contract-backed. `新建计划` opens the existing payment-plan form with `contractId`; `登记回款` opens payment detail for the existing payment form. Bulk transfer/archive stay disabled-with-reason because this page does not own those endpoints. | `GET /api/contracts/{id}`; `POST /api/contracts/{contractId}/payment-plans`; detail flow records payments via `POST /api/contracts/{contractId}/payments`. |

Additional already-functional list: `TaskList` has `查看` and `完成任务`; completion
uses `POST /api/tasks/{id}/status` and is disabled for `Completed`/`Cancelled`.

## B. Row Click Navigation

Shared `DataTable` row click wiring:

- `onRowClick(row, index)` added to `DataTable`.
- Click, Enter, and Space open rows when `isRowClickable` allows it.
- `getRowAriaLabel` names row-open targets.
- Interactive children (`button`, `a`, form controls, `[role="menu"]`,
  `[data-row-interactive="true"]`) are excluded from row-open propagation.

Shell-level record navigation:

- `frontend/src/app/navigation.ts` maps `Opportunity`, `Lead`, `Account`,
  `Contact`, `Quote`, `Contract`, `Payment`, and `Task` related records to
  `RecordNavigationTarget`.
- `Shell` receives the target, switches `view`, passes `targetRecordId` into the
  owning list, and each list loads the detail record then clears the target.

List row click targets:

| Surface | Row click opens |
|---|---|
| `LeadList` | lead detail via `getLead(id)`. |
| `AccountList` | account detail via `getAccount(id)`. |
| `ContactList` | contact detail via `getContact(id)`. |
| `OpportunityList` | opportunity detail via `getOpportunity(id)`. |
| `QuoteList` | quote detail via `getQuote(id)`. |
| `ContractList` | contract detail via `getContract(id)`. |
| `PaymentList` | payment/contract detail via `getContract(contractId)`. |
| `TaskList` | task detail in the task list component. |

Workbench/focus row targets:

| Workbench surface | Row behavior |
|---|---|
| `FunnelFocus` table `聚焦商机明细` | Opportunity rows open the opportunity detail. |
| `KeyOpportunityOverview` dashboard rows | Row button opens the opportunity detail. |
| `KeyOpportunityFocus` table `重点商机明细` | Opportunity rows open the opportunity detail. |
| `TodoOverview` dashboard rows | Row button opens the related task/contract/payment/etc. through `targetForRelatedRecord`. |
| `TodoFocus` table `待办与预警明细` | Rows with a target open the related record. |
| `PaymentsOverview` rows | Rows with `recordId` open the payment/contract detail. |
| `PaymentsFocus` table `回款到账明细` | Rows with `recordId` open the payment/contract detail; rows without a record target are not clickable. |
| `ActivityOverview` rows | Rows with a supported `relatedType/relatedId` open that related record. |
| `ActivityFocus` table `最近活动明细` | Rows with a supported related record open that record; unsupported rows are not clickable. |
| `StageFocus`, `TrendFocus`, `LeaderboardFocus` | Aggregate-only rows do not get row-click affordance or navigation. |

## C. Former Dead Buttons

| Surface/button | Current behavior |
|---|---|
| `ReminderCenter` - `按到期排序` | Toggles `sortDueFirst`, updates `aria-pressed`, and sorts visible reminders by due date. |
| `ReminderCenter` - row `查看` | Calls `openReminder`, maps `relatedRecord.type/id`, and navigates to the related task/contract/payment/etc. |
| `BasicReports` - `本月` | Sets a zh-CN status notice and refreshes the existing basic report endpoint. |
| `BasicReports` - `按负责人分组` | Scrolls/focuses the `负责人分组` panel and sets a zh-CN status notice. |
| `UserManagement` - `导出` | Exports the currently filtered user rows to `users-filtered.csv` client-side from existing `GET /admin/users` data. |
| `OperationLogs` - `今天` | Sets `timeFilter='today'`, resets pagination, and syncs the toolbar select. |
| `OperationLogs` - `导出` | Exports filtered safe-summary audit rows to `operation-logs-filtered.csv` client-side from existing `GET /api/operation-log` data. |
| `ImportExportPage` - `最近批次` | Removed; e2e asserts no `最近批次` button remains. |
| `ImportExportPage` - `新建导入` | Scrolls/focuses the CSV file input in the existing import form; submit remains `POST /api/imports`. |

## Constraint Check

- Backend/shared/root API diff: empty. Checked with
  `git diff --name-only -- services shared api packages/shared apps/api` and
  `git status --short -- services shared api packages/shared apps/api`.
- Existing endpoints only: all actions above call existing frontend API clients;
  no backend route, shared contract, or root API file was added.
- No new color: the functional styles use existing tokens such as `--card`,
  `--section`, `--primary`, `--border`, `--text`, and `--muted`; the CSS diff
  changes one old `#fff` use to `var(--card)` and adds no new functional color
  literal.
- zh-CN preserved: new labels, notices, disabled reasons, menu items, and exports
  are Chinese.
- Enum/role values unchanged: comparisons still use the existing values
  `Administrator`, `Sales Manager`, `Sales`, the six opportunity stages, and
  existing quote/contract/task status values. Display labels still go through
  `labelFor` where applicable.
- Role gates not widened: Sales restrictions on owner-transfer/archive/bulk
  management remain hidden or disabled; admin-only pages remain nav/API gated;
  frontend affordances do not replace server authorization.
- Unsupported backend behavior is not faked: missing owner-transfer/archive
  endpoints for contact/customer/quote/payment/task bulk paths remain omitted or
  disabled-with-reason instead of simulated.

## Verification

Recorded full verification from the completed implementation pass:

| Command | Result |
|---|---|
| `cd frontend && npx tsc --noEmit` | PASS, clean. |
| `cd frontend && npm run build` | PASS. |
| `cd frontend && npm run test:e2e` | PASS, 54/54, 0 failed, 0 skipped, `workers: 2`. |

Playwright worker configuration is set in `frontend/playwright.config.ts`:
`workers: 2`.

Functional assertions added or strengthened include:

- `TEST-UIUX-FUNC-ROWMENU-001`: lead row menu opens, `查看` runs, and row click
  opens the record detail.
- `TEST-UIUX-FUNC-ROWNAV-001`: dashboard key-opportunity row and focus-stage
  opportunity row open the opportunity detail.
- Reminder coverage asserts `按到期排序` toggles `aria-pressed` and row `查看`
  opens the related contract detail.
- Reports coverage asserts `本月` status notice, `按负责人分组` notice, and focus.
- User-admin coverage asserts `导出` downloads `users-filtered.csv`.
- Operation-log coverage asserts `今天` filter and `导出` downloads
  `operation-logs-filtered.csv`.
- Import coverage asserts `最近批次` is absent and `新建导入` focuses the CSV input.

Pending: Claude functional re-verification by clicking every touched control.
Codex does not self-resolve `BLK-UIUX-G12-017`.
