# UI/UX G12 Rework Phase 0 — Primitive Inventory

Date: 2026-06-08
Scope: Phase 0 only, per `delivery/uiux-g12-rework-1.md` and
`planning/blockers.md` BLK-UIUX-G12-003.

## Scope Guard

- Frontend/design primitive layer only: `frontend/src/components/ui/index.tsx` and
  `frontend/src/styles/design-system.css`.
- No page composition in this phase. Existing page/e2e diffs from the prior
  failed G9-G11 attempt remain outside this Phase 0 delta.
- No backend, API, data model, service-boundary, or `shared/` change.
- No new color tokens. New primitive styles reference locked design-system tokens
  and the already approved text-only `*-ink` tokens; `*-ink` remains text-only.
- No e2e assertion edits in this phase.

## Required Component Inventory

| Required component from locked mockups | Mockup/page coverage | Exists before Phase 0 | Primitive after Phase 0 | Spec compliance |
|---|---|---:|---|---|
| App shell / side nav / topbar / live scope affordance | All nav pages, dashboard variants | Exists | Shell/Nav page layer + `LiveToggle` | Existing shell primitive support kept; no Phase 0 page edits. |
| Card / panel containers | All page types | Exists | `Card`, `Panel`, `PanelHeader` | §7 card/panel hierarchy supported; `PanelHeader` added for title/meta/actions. |
| KPI / metric cards | Dashboard, reports, reminders, import-export | Exists | `MetricCard` | Token-bound metric tile with icon, label, value, delta slot. |
| Badges / status pills / chips | Lists, details, reports, reminders, admin | Exists | `Badge`, `StatusBadge` | Badge background colors stay locked soft tokens; readable text uses approved ink where applicable. |
| Multi-column data table | Lists x8, reports owner table, admin users, import errors | Partial | `DataTable<T>` | Now supports columns, headers, per-row selection, select-all, sort controls, row actions, empty state, and badge/status slots through `render`. Legacy children-table mode remains. |
| Pagination footer | Lists x8, admin users, operation log | Partial | `Pagination` | Now supports page numbers, previous/next, total count, page-size selector, and legacy page-chip fallback. |
| Filter/search toolbar | Lists x8, admin users, operation log | Partial | `Toolbar` | Now supports search, multiple filters, actions, clear filters, and active-filter summary/chips. Legacy children wrapper remains. |
| Bulk action bar | Lists x8 | Exists | `BulkActionBar` | Existing primitive retained for selected count/actions; Phase 1 will compose Sales-hidden actions per A4. |
| Funnel / pipeline bars | Dashboard, reports | Exists | `FunnelBars` | Existing discrete rows, labels, tracks, values retained and token-bound. |
| Trend line panel | Dashboard | Partial | `TrendPanel` | Now renders SVG trend line/area from point data while retaining old custom-children mode. |
| Stage/share donut CMP-014 | Dashboard, reports breakdown cards | Missing | `DonutChart`, alias `StageDonut` | Added SVG donut with token tone classes, percent legend, center label, no new colors. |
| Performance leaderboard | Dashboard | Missing | `Leaderboard` | Added ranked row primitive with metric, meta, and token-bound progress track. |
| Reminder row card CMP-011 | Reminder center, dashboard todo/alerts | Missing | `ReminderRowCard` | Added icon, title, description, badges, meta/time, actions, overdue treatment. |
| Card -> Focus stage container §10 | Dashboard manager focus | Partial CSS only | `FocusStage`, alias `CardFocusStage` | Added stage + right collapsed card strip API, tools/back/escape affordance slots; uses existing `.focus/.stage/.sideCard` tokens. |
| Read-only audit/event card | Operation log, import/export audit status | Missing | `AuditEventCard` | Added safe-summary-only card. API has `safeSummary` and metadata fields only; no before/after/raw diff props are accepted or rendered. |
| Form fields | Opportunity form, admin user edit/create, import/export | Exists | `TextField`, `SelectField`, `TextAreaField` | Existing label/hint/error structure retained. |
| State surfaces | All page types | Exists | `EmptyState`, `ErrorState`, `PermissionDenied`, `SkeletonBlock`, `InlineLoading`, `Toast`, `Drawer` | Existing A3 state primitives retained. |

## Locked Primitive APIs

- `DataTable<T>`: `columns`, `rows`, `rowKey`, `selectedRowKeys`,
  `onToggleRow`, `onToggleAll`, `getRowClassName`, `empty`, `actions`; each
  column has `header`, `render`, `sortable`, `sortDirection`, `onSort`, `align`,
  `width`, and class hooks.
- `Toolbar`: `searchValue`, `onSearchChange`, `searchPlaceholder`, `filters`,
  `activeFilters`, `onClearFilters`, `summary`, `actions`.
- `Pagination`: `page`, `totalPages`, `totalItems`, `pageSize`,
  `pageSizeOptions`, `onPageChange`, `onPageSizeChange`, plus legacy
  `onPrevious` / `onNext`.
- `TrendPanel`: `points` + optional `valueLabel`, or legacy custom children.
- `DonutChart` / `StageDonut`: `segments`, `label`, `center`; segment tone is
  limited to existing token tones.
- `Leaderboard`: `items` with `label`, `value`, `meta`, `suffix`, token tone.
- `ReminderRowCard`: icon/title/description/meta/time/badges/actions/tone/overdue.
- `FocusStage` / `CardFocusStage`: stage title/subtitle/icon/tools/back action,
  right-side collapsed `sideCards`, and content slot.
- `AuditEventCard`: `actor`, `action`, `resource`, `result`, `occurredAt`,
  `safeSummary`, optional ids/hash/badges. It intentionally does not model or
  render raw `before` / `after` payloads.

## Phase 0 Verification

- `npm run build` in `frontend/`: PASS.
- Build output: `dist/assets/index-C0OjiDaU.css`,
  `dist/assets/index-CP_r3b9e.js`.
- Backend/shared diff: none introduced by Phase 0.
- Page/e2e composition: intentionally not changed in Phase 0; reserved for
  Phase 1+ after Claude primitive-layer re-audit.
