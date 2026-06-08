# UI/UX G12 Rework — Final Touch-up Evidence

Date: 2026-06-08
From: Codex (execution)
To: Claude (G12 closing audit)
Status: Final touch-up returned. Codex does not self-pass G12.

## Scope

This is the final touch-up after Claude marked Phase 3 structures PASS and asked
for one required C5 stabilization plus cosmetic closure before the global G12
sweep.

Locked mockup references:

- `docs/ux-ui/mockups/reports-team.png`
- `docs/ux-ui/mockups/admin-users.png`
- `docs/ux-ui/mockups/detail-opportunity.png`

Final delta from Claude's `/tmp/uiux-p3done-baseline-src` snapshot is limited to:

- `frontend/src/pages/reports/ManagerOverview.tsx`
- `frontend/src/pages/reports/BasicReports.tsx`
- `frontend/src/pages/admin/UserManagement.tsx`
- `frontend/src/components/ui/index.tsx`
- `frontend/src/components/ActivityNoteTaskPanel.tsx`
- `frontend/src/styles/design-system.css`

E2E assertion updates are limited to:

- `frontend/e2e/persistence.spec.ts`
- `frontend/e2e/reports.spec.ts`
- `frontend/e2e/user-admin.spec.ts`
- `frontend/e2e/opportunities.spec.ts`
- `frontend/e2e/work.spec.ts`
- `frontend/e2e/overview.spec.ts`

## Implementation Summary

Required C5 stabilization:

- `TEST-PERSISTENCE-001..005` now waits for an interactive auth state before
  filling credentials.
- After logout, the spec waits for the login form instead of racing against the
  still-visible stale logout button.
- After login, the spec polls `/api/leads` until the authenticated gateway read
  returns `200`.
- After `docker compose restart lead`, the spec polls `/api/leads` until the
  restarted service is ready and the just-created lead is visible, then reloads
  and reasserts persistence.
- No persistence assertion was weakened, removed, skipped, or marked `only`.

Final reliability follow-up:

- `work.spec.ts` now gives the post-save activity timeline refresh assertions a
  15s timeout window for the full-suite load path.
- The affected assertions still require the same saved note text, activity
  count/text, task title, and pending status to become visible in `活动时间线`.
- No activity/task assertion was weakened, removed, skipped, or marked `only`.

Reports:

- Removed the extra stacked `经理团队总览` panel by making the report route render
  the Phase 3 `BasicReports` team report directly.
- Kept the existing `team-overview` API available for dashboard/report role
  verification; `overview.spec.ts` now asserts that API separately.
- KPI money values now use compact zh-CN currency formatting (`¥...万` / `¥...亿`)
  plus a fixed non-overflow card layout.

Admin users:

- The table action column now renders three readable buttons per row:
  `编辑`, `停用/启用`, and `改角色`.
- The only active administrator's disable button is disabled in-row, preserving
  the existing last-admin guard.
- Existing role/status confirmation flow remains unchanged.

Forms/details:

- `TextField type="date"` now adds `dateField` / `dateControl` styling hooks.
- Native date inputs are styled through the design system using existing tokens.
- Opportunity detail's right-side work panel now renders notes, activities, and
  tasks as an activity timeline.

## Constraint Check

- Frontend-only: yes.
- Backend/shared/root-api diff: 0. Both `git diff --name-only -- services shared api backend` and `git status --short -- services shared api backend` returned no paths.
- Existing endpoints only: yes. Persistence/report checks use existing `/api/leads`, `/api/reports/team-overview`, and existing page adapters.
- New colors: none. Diff from `/tmp/uiux-p3done-baseline-src/styles/design-system.css` adds no hex/rgb/hsl color literals.
- zh-CN: yes. New visible copy is Chinese.
- Enum/role comparison values unchanged: yes. Role/status values are display-mapped only where already mapped.
- No backend-needed aggregation: none identified; no Kickback blocker raised.
- No e2e downgrade: skip/only scan returned no matches.

## Verification

Commands run from `frontend/` unless noted:

- `npx tsc --noEmit` — PASS.
- `npm run build` — PASS; produced `dist/assets/index-CLRf63xe.css` and `dist/assets/index-DEY6x1ha.js`.
- `npx playwright test e2e/reports.spec.ts e2e/user-admin.spec.ts e2e/opportunities.spec.ts e2e/work.spec.ts e2e/persistence.spec.ts` — PASS, 15/15.
- `npx playwright test e2e/work.spec.ts` — PASS, 2/2 after final activity timeline stabilization.
- `npx playwright test e2e/persistence.spec.ts` — PASS, 1/1 after final auth-readiness stabilization.
- `npm run test:e2e` — PASS, 49/49, 0 skipped, in one complete run.
- `npm run test:e2e` — PASS, 49/49, 0 skipped, second consecutive complete run.
- `npm run test:e2e` — PASS, 49/49, 0 skipped, third consecutive complete run.
- From repo root, `rg -n "test\\.(skip|only)|describe\\.(skip|only)|it\\.(skip|only)" frontend/e2e` — no matches.
- From repo root, `git diff --check` — PASS.
- From repo root, `git diff --name-only -- services shared api backend` — no output.
- From repo root, `git status --short -- services shared api backend` — no output.
- From repo root, `diff -qr /tmp/uiux-p3done-baseline-src frontend/src` — only the six final touch-up source files listed in Scope differ.
- From repo root, `diff -u /tmp/uiux-p3done-baseline-src/styles/design-system.css frontend/src/styles/design-system.css | rg "^\\+[^+].*(#[0-9A-Fa-f]{3,8}|rgba?\\(|hsla?\\()"` — no matches.

Browser plugin note:

- Local Vite service started successfully on `http://127.0.0.1:5174/` for a
  Browser-plugin check, then was stopped.
- The Browser runtime was callable, but `agent.browsers.list()` returned `[]` and
  `iab` was unavailable (`Browser is not available: iab`), so no in-app Browser
  screenshot/DOM check could be captured in this session.
- Playwright e2e is therefore the runtime UI evidence for this return.

## Return

Returned to Claude for the final UI/UX G12 closing sweep: all 14 pages screenshot
fidelity, reliably green full e2e, and cumulative constraint verification. Codex
does not self-pass G12 and does not proceed beyond this return.
