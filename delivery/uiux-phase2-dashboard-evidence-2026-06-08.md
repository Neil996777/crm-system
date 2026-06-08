# UI/UX G12 Rework Phase 2 — Dashboard Evidence

Date: 2026-06-08
From: Codex (execution)
To: Claude (screenshot fidelity re-audit)
Status: Phase 2 returned for review. Codex does not self-pass G12 and does not enter Phase 3.

## Scope

Phase 2 implements the dashboard work from `delivery/uiux-g12-rework-1.md`:

- Manager / Administrator dashboard: 8 cards for sales funnel, stage donut, win trend, sales leaderboard, todo alerts, payments, key opportunities, and recent activity.
- Sales dashboard variant: personal 6-card layout without manager-only leaderboard and key-opportunity team card.
- Per-card Card -> Focus: each dashboard card has its own expand action and renders a focus stage with the remaining cards collapsed into the right-side `sideCards` strip.
- Conservative motion / reduced-motion path is preserved by using the existing focus-stage state and CSS motion tokens.

Phase 2 delta from Claude's `/tmp/uiux-p1done-baseline-src` snapshot is limited to:

- `frontend/src/pages/WorkOverview.tsx` — 990 lines.
- `frontend/src/app/Shell.tsx`.
- `frontend/src/styles/design-system.css`.

## Implementation Summary

`WorkOverview.tsx` now composes the Phase 0 primitives instead of the old
simplified dashboard:

- `MetricCard` renders the live 4-metric strip.
- `FunnelBars` renders the sales funnel card.
- `StageDonut` renders stage composition.
- Inline SVG trend visual plus `DataTable` renders the trend card and focus detail.
- `Leaderboard` renders the manager sales leaderboard.
- `DataTable` renders focus tables for funnel, stage, trend, leaderboard, todo,
  payments, key opportunities, and activity.
- `FocusStage` receives a per-card active card and generated `sideCards`; `Esc`
  and `返回` close focus mode.

`Shell.tsx` keeps only shell-level focus styling state. `WorkOverview` owns the
active card and reports focus state through `onFocusChange`.

## Data Sources

No dashboard metric is implemented as a static placeholder or hidden hardcoded
zero. Values come from existing frontend API adapters only:

- Manager/report aggregations where they already exist:
  `getManagerOverview()`, `getBasicReport()`.
- Existing authorized list endpoints used for frontend aggregation:
  `listOpportunities`, `listLeads`, `listQuotes`, `listContracts`,
  `listPaymentContracts`, `listTasks`, `listReminders`, and `listActivities`.

Cards without a dedicated existing backend aggregation use deterministic
frontend aggregation over those authorized records:

- Win trend: grouped from won opportunities over the last six months.
- Leaderboard: grouped from won opportunities by `ownerId`.
- Key opportunities: sorted from authorized opportunity rows by non-terminal
  status and amount.
- Todo alerts: merged from reminders and open tasks.
- Payments: uses `getBasicReport().breakdowns.paymentsByStatus` when present,
  otherwise existing payment/contract records.

This Phase 2 found no metric that required a new backend endpoint or new backend
aggregation. No blocker or Formal Scope Change is raised.

## Constraint Check

- Frontend-only: yes. `git diff --name-only -- services shared api backend`
  returned no paths; `git status --short -- services shared api backend` also
  returned no paths.
- Backend/shared diff: 0.
- New colors: none. The Phase 2 CSS diff does not add or change `:root` color
  tokens. Literal colors seen in the Phase 2 files are locked design-system
  values (`#fff`, `#ECEDFE`, `#F2ECFD`, `#2563EB`) or existing token references.
- zh-CN: dashboard display copy is Chinese. Enum and role values remain real
  comparison values; display text uses existing labels.
- Per-card focus: yes. `TEST-UIUX-DASHBOARD-001` asserts manager 8 cards and
  per-card focus. `TEST-UIUX-DASHBOARD-002` asserts the sales variant. `TEST-UIUX-A7-001`
  asserts focus behavior under reduced-motion mode.
- No assertion downgrade: no e2e assertion was removed for Phase 2; no
  `test.skip`, `test.only`, `describe.skip`, or `describe.only` was found.

## Verification

Commands run from `frontend/` unless noted:

- `npx tsc --noEmit` — PASS.
- `npm run build` — PASS; produced `dist/assets/index-nKPfxN8h.css` and
  `dist/assets/index-BPV95IRy.js`.
- `npm run test:e2e` — all 48 non-persistence tests PASS; the known
  `TEST-PERSISTENCE-001..005` flake timed out waiting for the sign-in email field
  because the page was already on the logged-in dashboard. No skipped tests.
- `npx playwright test e2e/persistence.spec.ts` — PASS, 1/1. Re-run twice after
  the full-suite timeout; both passed.
- `npm run test:e2e -- --workers=1` — same known persistence timeout after the
  other 48 tests passed; targeted persistence re-run above passed.
- `rg -n "test\\.(skip|only)|describe\\.(skip|only)" e2e` — no matches.
- From repo root, `diff -qr /tmp/uiux-p1done-baseline-src frontend/src` — only
  `app/Shell.tsx`, `pages/WorkOverview.tsx`, and `styles/design-system.css`
  differ from the post-Phase-1 snapshot.
- From repo root, `git diff --name-only -- services shared api backend` — no
  output.

The persistence behavior matches the Phase 2 instruction's known-flake protocol:
the spec failed only in full-suite context and passed on targeted re-run; no
assertion was reduced.

## Return

Returned to Claude for screenshot-vs-mockup fidelity re-audit against:

- `docs/ux-ui/mockups/dashboard-v7-manager.png`
- `docs/ux-ui/mockups/dashboard-v7-manager-focus.png`
- `docs/ux-ui/mockups/dashboard-v7-sales.png`

Codex does not self-pass Phase 2 and does not proceed to Phase 3.
