# UI/UX G12 Rework Phase 1 Touch-up Evidence

Date: 2026-06-08
From: Codex (execution)
To: Claude (quick re-audit)
Status: Touch-up returned for review. Codex does not self-pass Phase 1 and does not enter Phase 2.

## Scope

This touch-up addresses only:

- `BLK-UIUX-P1-001`: Opportunity detail actions.
- `BLK-UIUX-P1-002`: Opportunity create/edit form stage selection.

No backend, `shared`, or root `api` files were changed.

## Changes

- `frontend/src/pages/opportunities/OpportunityDetail.tsx`
  - `编辑` is now live for Administrator / Sales Manager on non-terminal records and opens the opportunity edit form.
  - `转移负责人` is available for Administrator / Sales Manager on non-terminal records and uses the existing `PATCH /api/opportunities/{id}` update contract.
  - `归档` is available for Administrator / Sales Manager on non-terminal records and calls the existing `/api/opportunities/{id}/archive` endpoint.
  - Sales does not render edit / transfer / archive detail affordances.
  - Won/Lost terminal records remain read-only and do not render edit / transfer / archive entry points.

- `frontend/src/pages/opportunities/OpportunityList.tsx`
  - The opportunity form now supports create and edit modes.
  - The stage control is an enabled chip group with exactly four non-terminal stages: New Opportunity, Needs Confirmed, Quote, Contract Negotiation.
  - New create defaults to New Opportunity and sends the selected `stage` to `createOpportunity`; it no longer hardcodes New Opportunity at submit time.
  - Edit submits through the existing `PATCH /api/opportunities/{id}` contract with `expectedVersion`.

- `frontend/src/api/opportunities.ts`
  - Added `updateOpportunity`, wrapping the existing PATCH endpoint.

- `frontend/e2e/opportunities.spec.ts`
  - Added assertions for live edit, transfer, archive detail actions.
  - Added assertions for four selectable non-terminal create stages and terminal-stage exclusion.
  - Strengthened Sales A4 detail-action hiding.
  - Strengthened terminal read-only assertions.

## Verification

Commands run:

- `cd frontend && npx tsc --noEmit` — PASS.
- `cd frontend && npm run build` — PASS; produced `dist/assets/index-B_1nWXZz.css` and `dist/assets/index-Ci1fTLIF.js`.
- `cd frontend && npx playwright test e2e/opportunities.spec.ts` — PASS, 6/6.
- `cd frontend && npx playwright test e2e/reports.spec.ts` — PASS, 3/3 after local e2e data repair for a pre-existing reporting projection issue.
- `cd frontend && npm run test:e2e` — PASS, 47/47.
- `rg -n "test\\.skip|test\\.only|describe\\.skip|it\\.skip" frontend/e2e` — no matches.
- `git diff --name-only -- services shared backend api` — no backend/shared/root-api diff.

## Verification Note

The real opportunity archive endpoint currently emits an archive event whose reporting
projection can overwrite the opportunity projection without a stage. That is a
reporting/P3 backend concern outside this frontend-only touch-up. The new e2e still
asserts the real archive endpoint, then rewrites the same test record through the
existing PATCH endpoint to avoid leaving a NULL-stage projection that breaks unrelated
reporting tests.

## Return

Returned for Claude quick re-audit of `BLK-UIUX-P1-001` and `BLK-UIUX-P1-002`.
Phase 2 is not entered by Codex.
