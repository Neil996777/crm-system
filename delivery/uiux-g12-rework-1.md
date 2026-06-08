# UI/UX G12 Audit — Rework #1 (FAIL: design not realized)

Date: 2026-06-08
From: Claude (independent G12 audit)
To: Codex (G9–G11 execution)
Decision: **G12 FAILED — full design-realization kickback (phased rebuild).**

## Correction of the first-pass verdict (integrity note)

An earlier same-day pass recorded G12 as "substantive PASS with 2 minor display
findings." **That verdict is overturned.** It was produced by checks that covered
design-SYSTEM concerns (tokens, role-gating logic, zh-CN/enum integrity) but did
**not** perform the A2/C6 requirement: per-page visual comparison of the running
implementation against the locked mockups. When that comparison was finally done
(live screenshots of all 14 nav pages at 1440px vs `docs/ux-ui/mockups/*.png`),
**every page is a MAJOR deviation** from the locked design. The design-realization
layer — the entire point of the UI/UX completion change — was not delivered.

Codex's G11 evidence claimed "14 nav pages reskinned to their mockups"; that claim
does not hold. The e2e added in G9–G11 asserts the as-built simplified shapes (e.g.
`TEST-UIUX-DASHBOARD-001` locks the two-mode toggle), so a green 45/45 e2e gave
false assurance of design fidelity. Going forward, e2e must assert the DESIGNED
structures (table columns, pagination controls, dashboard cards, reminder badges,
safe-summary audit rows), never the reduced placeholder shapes.

## What DID pass (keep, do not redo)

- Design-system foundation (A1): CSS tokens byte-match `design-system.md`; the four
  DEC-UIUX-A5-001 `*-ink` tokens are text-only; primitives in `components/ui/`.
- App shell + role-filtered nav (admin-only 用户与角色 / 操作日志; reports hidden from
  Sales) and topbar.
- Role-gating LOGIC and no data-exposure widening (C2) at the affordance level;
  zh-CN preserved except the one regression below (C3); enum/role comparison values
  unchanged (C4); build green; no backend/`shared` diff (C1).

These are necessary but not sufficient — they are the SYSTEM, not the page DESIGNS.

## Per-page fidelity findings (all MAJOR unless noted)

Evidence: live screenshots `/tmp/uiux-g12-shots/*` vs `docs/ux-ui/mockups/*.png`.

1. **Dashboard (工作台, manager/admin)** — locked: 8 data-rich cards (team funnel,
   stage donut, trend line, leaderboard, todo/alerts, payments, key opportunities,
   recent activity) + 4-metric live strip + the signature **Card→Focus** per-card
   expand (one card expands to a detail+table view while the others collapse to a
   right strip; design-system §10, A7). Implemented: 3 KPI tiles + 2 empty panels +
   a **global two-mode toggle** (`focusMode` boolean, fixed 3 generic side cards).
   No charts, no 8-card grid, no per-card expand. Sales variant
   (`dashboard-v7-sales.png`) also not realized.
2. **List pages ×8 (商机/线索/公司客户/联系人/报价/合同/回款/任务)** — locked
   (`list-opportunities.png` + sales variant): full-width multi-column data table
   (designed columns + colored stage/status badges), filter/toolbar row (search +
   relevant filters + 清除筛选 + active-filter summary), per-row selection + bulk-
   action bar (Sales-hidden per A4), pagination footer (page numbers + size
   selector). Implemented: the **old two-pane master-detail card list** (narrow
   left card column + right detail placeholder) on every entity — no table columns,
   no full filter toolbar, no pagination, no per-row select. Internally consistent
   (good reuse) but the wrong page-type. Includes folding **BLK-UIUX-G12-002**
   (bulk buttons hardcoded `disabled`).
3. **Reports (报表)** — locked (`reports-team.png`): 8 formatted KPI cards + a
   discrete 管道分析 funnel + 负责人分组 table + a row of 5 status/stage breakdown
   cards. Implemented: title drifted (经理团队总览), KPI values raw/unformatted, the
   discrete 管道分析 + 负责人分组 sections missing, layout reflowed into a long
   single-column stack. (Funnel bar viz partially present — the one bright spot.)
4. **Admin users (管理：用户与角色)** — locked (`admin-users.png`): bounded user
   roster table, last-admin guard banner, 角色/状态/仅管理员 filters, footer summary.
   Implemented: unpaginated long table, **last-admin guard banner not visibly
   rendered**, no footer summary, no pagination. (Role/status badges + 操作 actions
   do match.) Re-confirm the last-admin disabled affordance is present (e2e
   TEST-INV-LASTADMIN passes, so logic exists — the banner/treatment is the gap).
5. **Reminder Center (提醒中心)** — locked (`reminders-center.png`): 5 summary stat
   cards + grouped reminder cards with type icons/badges + right-rail 数据范围/空态.
   Implemented: flat text rows, no stat cards, no badges/icons/grouping, no right
   rail, and **raw English keys shown** (e.g. `reminder task reminder_…`) — a C3/A8
   regression on top of the layout gap.
6. **Import/Export (导入/导出)** — locked (`import-export.png`): two cards with full
   run-result reporting (总行数/成功/失败 + per-row error table; 导出行数/含归档/范围安全)
   + 审计与清理 status + header action bar. Implemented: object-type select + native
   file input + start buttons present, but **all result/audit-field reporting
   missing** and controls unstyled. Includes folding **BLK-UIUX-G12-001**
   (`fileSafety` raw English token at `Export.tsx:75` → add `labels.ts` map).
7. **Operation Log (操作日志)** — locked (`operation-log.png`): card-style read-only
   audit list + filter bar + **safe-summary** rows + pagination + right-rail
   防篡改/门控 cards. Implemented: an unpaginated raw table (~80k px tall), no
   filter bar, no pagination, no card/safe-summary layout. **Re-verify the rendered
   rows use `safeSummary`/sanitized fields and do NOT expose raw before/after
   values** (the prior token/gate audit said `summaryTextZh` is used — confirm in
   the rebuilt layout; if raw sensitive values are shown, that is also a C2 concern,
   not just fidelity).

## Phased rebuild plan (release owner elected phased kickback)

Frontend-only (C1). If realizing a locked card/section genuinely requires backend
data aggregation that doesn't exist, **raise a blocker (kickback) — do not silently
reduce the design**; the release owner decides scope (possible Formal Scope Change
by User per page). Each phase returns to Claude for a per-page fidelity re-audit
(screenshot vs mockup) before the next begins; G12 passes only after all phases
plus a final full sweep.

- **Phase 0 — Foundation/component-inventory confirmation (do FIRST).** The
  blocker is NOT missing tokens (A1 tokens already pass byte-match; `*-ink`
  re-audited) and NOT an empty primitive library — `frontend/src/components/ui/`
  already exports a rich set (`DataTable`, `Pagination`, `Toolbar`,
  `BulkActionBar`, `FunnelBars`, `TrendPanel`, form fields, `Card/Panel/
  MetricCard/Badge/StatusBadge`, `EmptyState/ErrorState/PermissionDenied/
  SkeletonBlock`, `Drawer/Toast/LiveToggle/InlineLoading`). The failure is that
  pages did not COMPOSE from these primitives to match the mockups. So before any
  page work, Codex must: (1) enumerate every component the locked mockups require
  and map each to an existing primitive; (2) verify each existing primitive matches
  its `design-system.md` §7 spec and is actually mockup-capable (e.g. `DataTable` =
  real multi-column + per-row select + sort; `Pagination` = page numbers + size
  selector; `Toolbar` = filters + 清除筛选 + active-filter summary; `FunnelBars`/
  `TrendPanel` match the funnel/trend viz); (3) build the genuinely MISSING
  primitives to spec — initial gap list to confirm: stage **donut/share viz**
  (CMP-014), **leaderboard**, **reminder row card** (CMP-011), the **Card→Focus**
  stage container (§10), and the **read-only audit/event card** with safe-summary.
  Deliver a short inventory doc (required component → exists/missing → primitive
  name → spec-match note) and lock the primitive API + tokens. Returns to Claude for
  a primitive-layer re-audit before Phase 1. No page composition in this phase.
- **Phase 1 — CRUD archetypes (list + detail + form), highest ROI (reused by 8
  entities).** COMPOSE from the Phase-0-confirmed primitives. Build the locked list
  table archetype (columns, badges, filter toolbar, per-row select, bulk bar with A4
  Sales-hiding, pagination) on opportunity as reference + the detail
  (`detail-opportunity.png`) and form (`form-opportunity.png`) archetypes, then
  apply across 线索/客户/联系人/报价/合同/回款/任务. Fold BLK-UIUX-G12-002.
- **Phase 2 — Dashboard.** Manager 8-card overview + sales variant + the per-card
  Card→Focus expand interaction + conservative motion (A7). Reuse existing reporting
  APIs for card data where available; raise a blocker for any metric that needs new
  backend aggregation.
- **Phase 3 — Reports + special pages.** Reports (管道分析 + 负责人分组 + KPI strip +
  breakdown cards), admin (guard banner + pagination + summary), reminders (stat
  cards + badged grouped cards + right rail + fix English keys / A8), import/export
  (result + audit reporting; fold BLK-UIUX-G12-001), operation-log (card list +
  filter + pagination + safe-summary). 

For every phase: update/strengthen e2e to assert the DESIGNED structures; never
reduce assertions or assert the placeholder shapes; keep build + e2e green, 0 skips;
no backend/`shared` diff; no enum/role comparison-value change; zh-CN preserved.

## Re-audit method on each phase return

Claude re-screenshots the affected pages at 1440px on the live app and compares to
the locked mockups (structure/components/interaction), in addition to token/gate/
i18n spot-checks and an independent build + e2e run. Codex does not self-pass G12.
