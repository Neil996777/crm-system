# UI/UX G12 Rework Phase 1 — CRUD Archetype Evidence

Date: 2026-06-08
From: Codex (execution)
To: Claude (independent screenshot fidelity re-audit)
Status: Phase 1 returned for review. Codex does not self-pass G12 or enter Phase 2.

## Scope

Phase 1 implements the CRUD list/detail/form archetypes from
`delivery/uiux-g12-rework-1.md`, using the Phase 0 locked primitives in
`frontend/src/components/ui/` and the shared `CrudScaffold` layer.

Covered entities:

- Opportunities
- Leads
- Accounts
- Contacts
- Quotes
- Contracts
- Payments
- Tasks

## Implemented Archetypes

Lists now compose from the Phase 0 primitives:

- `Toolbar`: search, entity filters, clear action, and active-filter summary.
- `DataTable`: multi-column table, row selection, select-all, row actions, and empty
  state.
- `BulkActionBar`: selected-count affordance, bulk action surface, clear selection.
- `CrudPagination`: page controls and page-size controls.

Sales A4 gating is preserved: Sales users do not render bulk transfer/archive
affordances. Manager/Admin bulk complete/archive/transfer behavior is connected to
existing single-record frontend actions where an existing backend endpoint exists;
otherwise the disabled action includes an A3 zh-CN reason. No new backend endpoint or
request/response contract was introduced.

Details/forms now compose from the detail/form archetype primitives:

- `DetailHero`, `DetailStat`, `Panel`, and status/stage badges for read-only detail.
- `StageStepper` for opportunity lifecycle display.
- `FormShell` and `FormSection` for create/edit surfaces.
- Terminal opportunity states render as read-only with disabled close/edit affordances.
- Opportunity create/edit excludes terminal stages.
- Sales owner locking is preserved on owner fields.

## Backend Contract Guard

This phase is frontend-only. No backend, `shared`, or root `api` files were changed.
Frontend API adapter changes only call existing backend contracts verified in the
current services:

- Existing archived-list filters: `includeArchived=true`.
- Existing archive endpoints for opportunity, lead, contract, and account flows.
- Existing lead owner-transfer endpoint.
- Existing task endpoints: `listTasks({ activeOnly, businessDate })`,
  `createTask`, and `changeTaskStatus`.

Task list search/status filtering is client-side over the existing task list result;
no unsupported task query parameters were invented.

## Verification

Commands run from `frontend/`:

- `npx tsc --noEmit` — PASS.
- `npm run build` — PASS; produced `dist/assets/index-Bozf24-W.css` and
  `dist/assets/index-BQLTU3YE.js`.
- `npm run test:e2e` — PASS, 45/45 tests.
- `rg "test\\.skip|test\\.only" frontend/e2e` from repo root — no matches.
- `git diff --name-only -- services shared backend api` from repo root — no
  backend/shared/root-api diff.

The e2e suite was updated to assert the designed Phase 1 structures, including table
columns, filter toolbar, row selection, bulk bar, and pagination for the task list,
while preserving the existing business-flow assertions.

## Return

Phase 1 is returned to Claude for screenshot-vs-mockup re-audit against:

- `docs/ux-ui/mockups/list-opportunities.png`
- `docs/ux-ui/mockups/list-opportunities-sales.png`
- `docs/ux-ui/mockups/detail-opportunity.png`
- `docs/ux-ui/mockups/form-opportunity.png`

Codex does not mark Phase 1 passed and does not proceed to Phase 2.
