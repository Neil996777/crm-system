# Frontend zh-CN Localization — Phase 2: Backend-originated messages (Claude → Codex)

## Document Control

- Project: CRM System
- Type: continuation of the zh-CN scope change (phase 1 = `delivery/i18n-zh-CN-handoff.md`, done in
  commit `d15f233`). Phase 1 localized all frontend chrome/forms/enums (display-map pattern), build
  green, e2e updated. Phase 2 closes the remaining "mixed language" surface.
- Approach (user-decided, 2026-06-05): **frontend message catalog ("够用档")** — map
  backend-originated user-visible English to Chinese in the FRONTEND. **Do NOT change backend
  code.** Best-practice principle: backend returns stable locale-neutral codes/messages; the
  presentation layer localizes. Backend stays frozen (already audited + live).
- Roles: Claude plans (this doc) + audits; Codex implements + redeploys.

## What's left after phase 1

1. **63 backend error/validation messages** are shown raw in alert toasts. The frontend receives
   `ApiError { code, category, safeMessage, fieldErrors[].code }` (`api/client.ts`) and components
   render `apiError.safeMessage` directly (e.g. `SessionProvider.tsx:40`, `TaskList.tsx:21/34`,
   `QualificationActions.tsx:26`, `ConvertLeadDialog.tsx:28`, and the other catch blocks). The
   English `safeMessage` therefore reaches the user.
2. **A few raw backend fields** still render English, notably `event.action` in
   `pages/admin/OperationLogs.tsx:57` (e.g. `create_user`). (History `beforeSummary`/`afterSummary`
   already go through `summaryTextZh`; verify those are fully covered too.)

## Implementation

### 1. Central error localizer (single injection point)

Add to `frontend/src/i18n/labels.ts` (or a sibling `errors.ts`) a map keyed by the exact backend
`safeMessage` string → Chinese, plus a helper:

```ts
export function localizeError(err: ApiError | undefined): string {
  if (!err) return '请求失败。';
  return errorMessageZh[err.safeMessage] ?? err.safeMessage; // fallback: show raw, never blank
}
```

Then replace every `apiError.safeMessage || '...'` / `error.safeMessage || '...'` usage with
`localizeError(apiError)`. Keep the existing Chinese fallbacks as the map's default for
`REQUEST_FAILED`. Also localize `fieldErrors[].safeMessage` the same way where field errors are shown.

**Fallback rule (best-practice robustness):** if a message is not in the map, render the raw backend
string (so nothing is ever blank), and the map miss is acceptable — list any newly-seen unmapped
message in the return-to-Claude notes. Do NOT throw on unmapped.

### 2. `event.action` and any remaining raw audit fields

Map `event.action` codes (e.g. `create_user`, `change_role`, `change_status`) → Chinese via a small
`actionLabel` map in labels.ts, fallback to raw. Confirm `summaryTextZh` covers the real
before/after summary shapes; extend if any English leaks.

### 3. Canonical EN → 中文 message table (use verbatim for consistency)

> Planning owns these translations. Use exactly these; for any backend message not listed here,
> translate following the same conventions ("... input is invalid." → "……输入无效。";
> "The requested X transition is not allowed." → "不允许该 X 流转。") and add it to the map.

| Backend message | 中文 |
|---|---|
| The request is invalid. | 请求无效。 |
| The account input is invalid. | 客户输入无效。 |
| The archive input is invalid. | 归档输入无效。 |
| The contact input is invalid. | 联系人输入无效。 |
| The lead input is invalid. | 线索输入无效。 |
| The qualification input is invalid. | 资质评估输入无效。 |
| The conversion input is invalid. | 转化输入无效。 |
| The opportunity input is invalid. | 商机输入无效。 |
| The stage transition input is invalid. | 阶段流转输入无效。 |
| The close-won input is invalid. | 赢单输入无效。 |
| The close-lost input is invalid. | 丢单输入无效。 |
| The quote input is invalid. | 报价输入无效。 |
| The quote status input is invalid. | 报价状态输入无效。 |
| The contract input is invalid. | 合同输入无效。 |
| The contract status input is invalid. | 合同状态输入无效。 |
| The contract quote link is invalid. | 合同与报价的关联无效。 |
| The payment input is invalid. | 回款输入无效。 |
| The payment plan input is invalid. | 回款计划输入无效。 |
| The task input is invalid. | 任务输入无效。 |
| The task status input is invalid. | 任务状态输入无效。 |
| The work item input is invalid. | 工作项输入无效。 |
| The owner transfer input is invalid. | 负责人转移输入无效。 |
| The duplicate check input is invalid. | 查重输入无效。 |
| The duplicate warning confirmation is invalid. | 重复提醒确认无效。 |
| The duplicate warning confirmation was already used. | 该重复提醒确认已被使用。 |
| The obligation query input is invalid. | 未结事项查询输入无效。 |
| The reminder query input is invalid. | 提醒查询输入无效。 |
| The filter is invalid. | 筛选条件无效。 |
| The projection input is invalid. | 报表数据输入无效。 |
| The import input is invalid. | 导入输入无效。 |
| The export input is invalid. | 导出输入无效。 |
| The CSV content is invalid. | CSV 内容无效。 |
| Only CSV import is supported. | 仅支持 CSV 导入。 |
| The object type is not supported for import. | 该对象类型不支持导入。 |
| The object type is not supported for export. | 该对象类型不支持导出。 |
| Export confirmation is required. | 需要确认后才能导出。 |
| The import run could not be saved. | 导入记录保存失败。 |
| The export could not be completed. | 导出未能完成。 |
| The requested stage transition is not allowed. | 不允许该阶段流转。 |
| The requested quote status transition is not allowed. | 不允许该报价状态流转。 |
| The requested contract status transition is not allowed. | 不允许该合同状态流转。 |
| The requested task status transition is not allowed. | 不允许该任务状态流转。 |
| Terminal opportunity records cannot be closed again. | 已终结的商机不能再次关闭。 |
| Terminal opportunity records cannot be edited. | 已终结的商机不能编辑。 |
| Terminal opportunity records cannot change stage. | 已终结的商机不能变更阶段。 |
| A quote already exists for this opportunity. | 该商机已存在报价。 |
| A contract already exists for this quote. | 该报价已存在合同。 |
| A reason is required when contract amount differs from quote amount. | 合同金额与报价金额不一致时必须填写原因。 |
| Signed or effective date is required for this contract status. | 该合同状态需要填写签署或生效日期。 |
| Won requires a Signed related contract. | 赢单需要有已签署的关联合同。 |
| Lost reason is required. | 必须填写丢单原因。 |
| Payment amount must be greater than zero. | 回款金额必须大于零。 |
| Payment exceeds the remaining contract amount. | 回款金额超过合同剩余应收。 |
| Payments use the committed single currency. | 回款只能使用约定的单一币种。 |
| The lead cannot be converted in its current state. | 当前状态的线索无法转化。 |
| The lead has already been converted. | 该线索已转化。 |
| The record changed after it was opened. | 记录在你打开后已被他人修改，请刷新重试。 |
| The requested resource was not found. | 未找到所请求的资源。 |
| Permission denied. | 没有权限执行该操作。 |
| A required service is unavailable. | 依赖的服务暂不可用，请稍后重试。 |
| Service authentication failed. | 服务认证失败。 |
| Audit log failed. | 审计日志写入失败。 |
| The audit event could not be persisted. | 审计事件未能持久化。 |

(If the live code yields any message not in this table, translate it the same way and add it; report
the additions.)

## Constraints

- **No backend change.** Frontend-only. No new runtime dependency.
- No-downgrade: don't weaken/skip/delete tests. e2e specs that assert the English alert text (e.g.
  `quotes.spec.ts:24` `'The quote input is invalid.'`, `leads.spec.ts:22`, `accounts.spec.ts:22`,
  `work.spec.ts:20`, `payments.spec.ts:70/80`, `opportunities.spec.ts:19/32`, `contracts.spec.ts:28/38/63`)
  must be updated to the new Chinese copy and stay green.
- Don't touch enum/role/stage compared VALUES (phase-1 rule still holds).

## Acceptance

- All 63 backend messages render in Chinese in the UI; unmapped messages fall back to raw (never
  blank, never crash). `event.action` and audit summaries render Chinese.
- `npx tsc --noEmit` + `npm run build` green; new bundle hash.
- e2e updated to the Chinese alert copy and passing; 0 skips.
- **Redeploy to production** per `deploy/ops/go-live-runbook.md` (frontend rebuild + `up -d --build`
  is not needed for FE-only, but the served `frontend/dist` must be rebuilt and published; if only
  the SPA changed, rebuilding `frontend/dist` + Nginx serving it is sufficient — capture the new
  bundle hash as evidence). Confirm the live site shows Chinese error toasts.

## Definition of Done

- Phase-2 map implemented; e2e green; build green; deployed; live error toast shows Chinese.
- `delivery/tasks.md` / traceability updated; commit made.
- Return to Claude for audit (string coverage of the 63 + live spot-check of a triggered error in
  Chinese + confirm served bundle is the new one). Do not self-pass.
