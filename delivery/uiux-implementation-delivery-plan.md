# UI/UX Implementation Delivery Plan

Status: **G7/G8 task package produced for Claude G8 handoff audit**. No
implementation has started.
Date: 2026-06-06
Owner platform: Codex (design/production task package).
Next reviewer: Claude (independent G8 design/planning handoff audit).

## Scope

This package implements the locked UI/UX design into the existing React frontend
after Claude passes G8. It covers the follow-on UI/UX completion charter and the
yardstick in `docs/ux-ui/requirements/uiux-implementation.requirements.md`.

Authoritative inputs:

- `delivery/uiux-completion-charter.md`
- `docs/ux-ui/requirements/uiux-implementation.requirements.md`
- `docs/ux-ui/design-system.md`
- `docs/ux-ui/mockups/*.png`
- `docs/ux-ui/mockups/_src/*.html`
- `docs/ux-ui/requirements/list-opportunities.requirements.md`
- `docs/ux-ui/requirements/batch-archetypes.requirements.md`
- `docs/ux-ui/requirements/special-pages.requirements.md`
- `docs/ux-ui/screen-state-spec.md`
- `docs/ux-ui/interaction-spec.md`
- Existing frontend files under `frontend/src/` and `frontend/e2e/`

## Technical Foundation Decision

The implementation should build a small project-local React style system and
component library inside `frontend/src/`:

- CSS variables in a dedicated design-system stylesheet, imported by
  `frontend/src/styles.css`.
- Project-local UI primitives under `frontend/src/components/ui/`.
- No new runtime dependency is required; keep React, Vite, TypeScript, and
  existing `lucide-react`.
- Icons should use existing `lucide-react` icons where possible.
- Existing API modules, auth provider, labels maps, page routing model, and e2e
  harness are retained.

Recommended file targets for G9:

- `frontend/src/styles.css` as the root import surface.
- `frontend/src/styles/design-system.css` for locked CSS variables and global
  base rules.
- `frontend/src/styles/motion.css` for motion tokens and reduced-motion rules.
- `frontend/src/components/ui/` for Shell, PageHeader, Card, Button, Badge,
  DataTable, Toolbar, FormField, EmptyState, ErrorState, PermissionDenied,
  Skeleton, LiveToggle, Pagination, BulkActionBar, MetricCard, FunnelBars,
  DonutChart, TrendPanel, TimelineRow, Drawer, Toast.
- Existing `frontend/src/app/Shell.tsx` and `Nav.tsx` should be refit to the
  design-system app shell, not replaced with a new routing model.

All values must be copied or directly derived from `docs/ux-ui/design-system.md`.
Do not recolor, invent new visual tokens, or introduce an alternate theme.

## Locked Constraints For Every Task

- **C1 frontend/design only:** no backend, API, data model, business logic, or
  service-boundary change.
- **C2 no downgrade:** do not weaken any P0/P1 behavior or prior G12 security
  repair, including IDOR, durable audit, optimistic concurrency, and idempotency.
  Visual changes must not change authorization behavior or data exposure.
- **C3 zh-CN preserved:** no English UI regression. Do not change enum or role
  comparison values.
- **C4 real enums:** display labels go through `frontend/src/i18n/labels.ts`;
  comparison values remain backend/API true values.
- **C5 e2e green:** update selectors only where DOM structure changes; do not
  skip tests or reduce assertions.
- **C6 locked visual contract:** match locked mockups and design system; tech
  feel comes from components, data presentation, and restrained interaction.

## Implementation Coverage Target

The task plan covers all 14 navigation pages through the 9 locked page types:

| Nav page | Page type / mockup target |
|---|---|
| 工作台 | dashboard-v7-sales, dashboard-v7-manager, dashboard-v7-manager-focus |
| 线索 | list archetype + detail/form reusable CRUD primitives |
| 公司/客户 | list archetype + detail/form reusable CRUD primitives |
| 联系人 | list archetype + detail/form reusable CRUD primitives |
| 商机 | list-opportunities, list-opportunities-sales, detail-opportunity, form-opportunity |
| 报价 | list archetype + detail/form reusable CRUD primitives |
| 合同 | list archetype + detail/form reusable CRUD primitives |
| 回款 | list archetype + detail/form reusable CRUD primitives |
| 任务 | list archetype + detail/form reusable CRUD primitives |
| 提醒中心 | reminders-center |
| 报表 | reports-team |
| 导入/导出 | import-export |
| 管理：用户与角色 | admin-users |
| 操作日志 | operation-log |

The eight CRUD entities are Leads, Accounts/Customers, Contacts, Opportunities,
Quotes, Contracts, Payments, and Tasks/Work Items. They reuse list/detail/form
primitives instead of each page inventing a visual pattern.

## Objective Acceptance

"Design implemented" means Claude can inspect the G9/G11 output against the
locked inputs and verify:

- A1 design tokens and components exist and are used by pages.
- A2 all 14 nav pages map to their page type.
- A3 canonical states are implemented: `loading`, `empty`, `error`, `disabled`,
  `selected`, `focused`, `hover`, `permission-denied`, `optimistic-update`,
  `success`.
- A4 role gates match the reviewed rules.
- A5 accessibility baseline is testable.
- A6 desktop-first layouts are stable around the 1440px reference and degrade
  without overflow at narrower desktop/tablet widths.
- A7 motion is restrained, tokenized, and snaps under `prefers-reduced-motion`.
- A8 Reminder Center status/priority display values are aligned with real
  backend values through display mapping only.
- A9 Import/Export samples and result fields are internally consistent with real
  data fields, especially `archivedIncluded`.

## Verification Expectations For G9-G11

The implementation agent should run, at minimum:

- `npm run build` from `frontend/`.
- `npm run test:e2e` from `frontend/`.
- A static diff review showing no changes under backend service/API/data-model
  paths for this UI/UX package.
- A selector/coverage review showing no Playwright `skip`/`only` and no reduced
  assertions.
- Manual or browser evidence for representative desktop pages: 工作台, 商机列表
  manager and sales variants, 商机详情, 商机表单, 报表, 提醒中心, 导入/导出,
  用户与角色, 操作日志.

