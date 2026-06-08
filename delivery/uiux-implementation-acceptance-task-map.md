# UI/UX Implementation Acceptance → Task Map

Status: **G9 coverage proof updated for DEC-UIUX-A5-001 token checkpoint**.
Scope: UI/UX completion only. This map does not replace the historical
`delivery/acceptance-task-map.md`; it adds the follow-on design-realization task
coverage for ACC-018 and ACC-023 plus the UI/UX yardstick A1-A9/C1-C6.

## Product Acceptance Binding

| Acceptance | Priority | Capability | UI/UX realization coverage |
|---|---|---|---|
| ACC-018 | P1 | Manager team overview / CAP-009 | UIUX-001, UIUX-002, UIUX-006, UIUX-008, UIUX-010, UIUX-011, UIUX-012, UIUX-014 |
| ACC-023 | P1 | Basic sales reports / CAP-009 | UIUX-001, UIUX-006, UIUX-008, UIUX-010, UIUX-011, UIUX-012, UIUX-014 |

ACC-018/023 functional behavior already exists in the audited implementation.
This package covers the missing design-realization layer without changing
backend/API/reporting semantics.

## Yardstick A1-A9 Coverage

| Yardstick | Requirement | Tasks | Objective audit hook |
|---|---|---|---|
| A1 | React design system from locked tokens/components | UIUX-001, UIUX-002, UIUX-003, UIUX-004, UIUX-005, UIUX-006, UIUX-007 | Token diff against `design-system.md`; pages use shared primitives; no new color token except DEC-UIUX-A5-001 text-only `*-ink` additions. |
| A2 | 14 nav pages mapped to 9 page types; 8 CRUD entities reuse patterns | UIUX-002, UIUX-003, UIUX-004, UIUX-005, UIUX-006, UIUX-007 | Page-by-page inspection against mockups and nav list. |
| A3 | Canonical states implemented | UIUX-008 plus page tasks UIUX-002..007 | Representative state checks for loading/empty/error/disabled/selected/focused/hover/permission-denied/optimistic-update/success. |
| A4 | Reviewed role/permission gates | UIUX-003, UIUX-004, UIUX-007, UIUX-009, UIUX-014 | Sales bulk actions hidden, terminal read-only, form terminal stages excluded, admin-only pages, last-admin disabled affordance. |
| A5 | Accessibility baseline | UIUX-001, UIUX-010, UIUX-014 | Keyboard path, focus ring, labels, landmarks, aria-live, contrast, target size evidence. |
| A6 | Desktop-first responsive no-break layout | UIUX-002, UIUX-003, UIUX-004, UIUX-005, UIUX-006, UIUX-007, UIUX-011 | Screenshots at 1440px and narrower desktop/tablet width; no overlap/overflow. |
| A7 | Conservative motion + reduced-motion snap | UIUX-001, UIUX-002, UIUX-003, UIUX-006, UIUX-012 | Motion tokens, NAV-01, MOTION-02, LIVE-03, `prefers-reduced-motion` evidence. |
| A8 | Reminder status/priority real display values | UIUX-007, UIUX-013, UIUX-014 | Reminder rows display backend values through labels/display mapping only. |
| A9 | Import/export sample consistency | UIUX-007, UIUX-013, UIUX-014 | ImportRun/ExportRun fields align; `archivedIncluded` and confirmation copy consistent. |

## Constraint C1-C6 Coverage

| Constraint | Guardrail | Tasks | Audit evidence |
|---|---|---|---|
| C1 | Frontend/design only | UIUX-001..014 | Diff excludes backend/API/data-model/business-logic/service-boundary changes. |
| C2 | No P0/P1 or G12 security downgrade | UIUX-003, UIUX-004, UIUX-005, UIUX-007, UIUX-009, UIUX-014 | Existing security/permission e2e remains green; no visual affordance widens data scope. |
| C3 | Preserve zh-CN; enum/role comparison values unchanged | UIUX-001..014 | UI text zh-CN; `labels.ts` additions are display-only; role/enum values not modified. |
| C4 | Real enums through labels | UIUX-003, UIUX-004, UIUX-005, UIUX-006, UIUX-007, UIUX-013 | Six opportunity stages and other status labels come from `labels.ts`; no invented values. |
| C5 | E2E green; no skip/assertion reduction | UIUX-014 | `npm run test:e2e`; grep for `skip`/`only`; selector changes preserve assertions. |
| C6 | Locked mockup/design-system consistency | UIUX-001..014 | Screenshot/page inspection against locked mockups and `design-system.md`. |

## Page-Type Coverage

| Page type / locked frame | Tasks | Nav pages covered |
|---|---|---|
| Dashboard / manager / sales / focus | UIUX-002, UIUX-012 | 工作台 |
| List archetype + sales variant | UIUX-003, UIUX-005, UIUX-008, UIUX-009 | 线索, 公司/客户, 联系人, 商机, 报价, 合同, 回款, 任务 |
| Detail archetype | UIUX-004, UIUX-005, UIUX-008, UIUX-009 | 8 CRUD entities where detail exists |
| Form archetype | UIUX-004, UIUX-005, UIUX-008, UIUX-009 | 8 CRUD entities where create/edit exists |
| Reports / team overview | UIUX-006, UIUX-010, UIUX-011, UIUX-012 | 报表 |
| Admin users | UIUX-007, UIUX-009, UIUX-010 | 管理：用户与角色 |
| Reminder Center | UIUX-007, UIUX-013 | 提醒中心 |
| Import/Export | UIUX-007, UIUX-013 | 导入/导出 |
| Operation Log | UIUX-007, UIUX-009, UIUX-010 | 操作日志 |

## Reverse Task Map

| Task | Primary coverage |
|---|---|
| UIUX-001 | A1, C1-C6, ACC-018/023 design foundation |
| UIUX-002 | A2 dashboard/shell, A3/A5/A6/A7 shell paths, ACC-018 |
| UIUX-003 | A2 list archetype, A4 sales/manager/admin list gates |
| UIUX-004 | A2 detail/form, A4 opportunity terminal/stage/owner gates |
| UIUX-005 | A2 8 CRUD entity application, A3 states per entity |
| UIUX-006 | ACC-018/023 primary reports realization, A6/A7 report visuals |
| UIUX-007 | Special pages, admin-only/oplog/read-only/reminder/import-export semantics |
| UIUX-008 | A3 canonical state layer |
| UIUX-009 | A4 role and permission gate review |
| UIUX-010 | A5 accessibility baseline |
| UIUX-011 | A6 desktop-first responsive stability |
| UIUX-012 | A7 conservative motion/reduced-motion |
| UIUX-013 | A8/A9 folded observations |
| UIUX-014 | C5 e2e and final evidence |
