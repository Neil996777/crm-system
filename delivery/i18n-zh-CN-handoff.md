# Frontend Chinese (zh-CN) Localization — Planning Handoff (Claude → Codex)

## Document Control

- Project: CRM System
- Type: **post-release scope change** (user request, 2026-06-05). The UI language was never
  specified in the design baseline (`docs/`, `modeling/` contain no UI-language requirement), so
  this is a NEW requirement, not a defect.
- Decision (user): **Full single-language Chinese (zh-CN)** — replace the hardcoded English UI text
  with Chinese; NO language switcher, NO i18n framework.
- Roles: Claude plans (this doc) + audits; Codex implements. Claude does not write the impl.
- Current state: frontend has no i18n library; UI strings are hardcoded English across 37 `.tsx`
  files; source contains no Chinese today.

## Goal

Every user-visible string in the CRM web UI renders in natural Simplified Chinese, with NO change
to application behavior, data, API contracts, or access-control logic.

## ⚠️ Two hard landmines (the whole reason this needs planning)

### 1. Do NOT translate logic/contract VALUES — only their DISPLAYED label

Many "English strings" are not display text — they are values compared against backend data or
used as keys. Changing the string literal breaks authz/state logic. These MUST keep their exact
current value; translate only what is rendered to the screen via a display-map layer:

- **Role values**: `'Administrator'`, `'Sales Manager'`, `'Sales'` (e.g. `app/Nav.tsx` `roles: [...]`,
  role comparisons). Keep the value; map for display.
- **Status / stage / lifecycle enums** compared in code: e.g. `lead.status === 'Pending Qualification'`,
  `'Unassigned'`, opportunity stages, contract `'Signed'`, user `'Active'`, etc.
- **View keys**: `'overview'`, `'leads'`, … in `AppView`.
- **API field names, object keys, route/query params, test ids, `aria` role strings, icon names.**

Rule: introduce a centralized display-map module (e.g. `frontend/src/i18n/labels.ts`) exporting
maps like `roleLabel`, `leadStatusLabel`, `opportunityStageLabel`, `contractStatusLabel`,
`paymentStatusLabel`, etc. Render `roleLabel[user.role]` etc. NEVER replace the compared literal.
If a backend enum value reaches the UI without a known mapping, show the raw value (don't crash) and
list it for follow-up — do not invent a translation that could mask an unknown state.

### 2. e2e selectors depend on the English strings — update them in lockstep

`frontend/e2e/*.spec.ts` locate elements by English text and assert English copy, e.g.:
`getByLabel('Email')`, `getByLabel('Password')`, `getByRole('button', { name: 'Sign in' })`,
`getByRole('heading', { name: 'Work Overview' })`, `getByRole('button', { name: 'Quotes' })`,
`getByRole('alert')).toContainText('The quote input is invalid.')`.

Translating the UI WILL break every one of these. For each translated string that an e2e test
targets or asserts, update the spec to the new Chinese string (or, preferred where practical,
switch the locator to a stable `data-testid` and keep the assertion on the Chinese copy). e2e must
be green after the change — do not delete or skip tests to make them pass (no-downgrade).

## Translation surface (Codex does the exhaustive sweep)

- Chrome/shell: `app/Nav.tsx`, `app/Shell.tsx`, `pages/SignIn.tsx`, `auth/*`.
- 14 components in `components/*.tsx` (dialogs, tables, timelines, steppers, warnings).
- 21 page files under `pages/{accounts,admin,contracts,importexport,leads,opportunities,payments,quotes,reminders,reports}`.
- Translate: headings, nav/menu labels, buttons, form labels & placeholders, table headers,
  helper/empty-state text, validation/toast/error messages shown to users, dialog titles, tooltips,
  `aria-label` user-facing text, and the document `<title>`.

## Canonical glossary (use these for consistency; extend in `labels.ts`)

Navigation: Work Overview→工作台 · Leads→线索 · Companies/Customers→公司/客户 · Contacts→联系人 ·
Opportunities→商机 · Quotes→报价 · Contracts→合同 · Payments→回款 · Tasks→任务 ·
Reminder Center→提醒中心 · Reports→报表 · Import/Export→导入/导出 · Admin: Users/Roles→管理：用户与角色 ·
Operation Logs→操作日志.

Auth: Sign in→登录 · Sign out→退出登录 · Sign in to continue→登录以继续 · Email→邮箱 · Password→密码.

Roles (display map, value unchanged): Administrator→管理员 · Sales Manager→销售经理 · Sales→销售.

Common statuses (display map; confirm exact enum set from code/api before mapping):
Unassigned→未分配 · Pending Qualification→待确认 · Qualified→已确认 · Won→赢单 · Lost→丢单 ·
Signed→已签署 · Active→启用 · Inactive→停用 · Archived→已归档.

Common actions: Save→保存 · Cancel→取消 · Create/New→新建 · Edit→编辑 · Delete→删除 · Confirm→确认 ·
Search→搜索 · Export→导出 · Import→导入 · Close→关闭.

(For any term not above, pick natural CRM Chinese and add it to `labels.ts` so usage stays uniform.)

## Constraints

- No behavior/logic/data change; no API field or enum VALUE change; no route change.
- No new runtime dependency (single-language; do NOT add react-i18next or similar — the user chose
  direct replacement). A plain `labels.ts` map module is fine.
- No-downgrade: do not weaken/skip/delete tests; update e2e to the Chinese copy and keep them green.
- Keep numbers, IDs, currency, dates as-is (formatting is out of scope unless trivially natural).

## Acceptance criteria

- Every user-visible UI string renders in Simplified Chinese (login, nav, all 10 feature areas,
  dialogs, tables, validation/toasts, empty states, admin pages).
- Role/status/stage values shown in Chinese via display maps; underlying compared values unchanged;
  access control and state transitions behave exactly as before.
- `npx tsc --noEmit` and `npm run build` green.
- `frontend/e2e/*` updated to the Chinese copy (or stable test-ids) and passing; 0 skips; no test
  deleted/weakened.
- A short follow-up list of any backend-originated strings still surfacing in English (e.g. server
  error messages, seeded enum labels) that would need a separate backend/display decision.

## Definition of Done

- All UI text Chinese; landmines #1/#2 handled; build + e2e green; no scope/contract/behavior change.
- `delivery/tasks.md` / traceability updated for the new requirement; commit made.
- Return to Claude for an independent audit (string-coverage sweep + verify no enum VALUE was
  changed + e2e green + a live spot-check on the deployed UI). Do not self-pass.
