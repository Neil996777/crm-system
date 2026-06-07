# UI/UX Implementation Task Dependencies

Status: **G7/G8 dependency DAG for Claude G8 audit**.
No implementation may start until Claude passes G8.

## Dependency DAG

Edges point from prerequisite to dependent.

```text
UIUX-001 React design-system foundation
  ├─> UIUX-002 App shell, nav, topbar, dashboard/focus
  │     ├─> UIUX-003 List archetype and opportunity list variants
  │     │     ├─> UIUX-004 Detail/form archetypes and opportunity detail/form
  │     │     │     └─> UIUX-005 Apply CRUD archetypes to 8 entity areas
  │     │     └─> UIUX-006 Reports and manager overview CAP-009 realization
  │     └─> UIUX-007 Special page types
  ├─> UIUX-008 Canonical state layer
  │     ├─> UIUX-003
  │     ├─> UIUX-004
  │     ├─> UIUX-005
  │     ├─> UIUX-006
  │     └─> UIUX-007
  ├─> UIUX-009 Role and permission gate review
  │     ├─> UIUX-003
  │     ├─> UIUX-004
  │     ├─> UIUX-006
  │     └─> UIUX-007
  ├─> UIUX-010 Accessibility baseline
  │     └─> all page tasks (UIUX-002..007)
  ├─> UIUX-011 Desktop-first responsive stability
  │     └─> all page tasks (UIUX-002..007)
  └─> UIUX-012 Conservative motion and reduced-motion path
        └─> UIUX-002, UIUX-003, UIUX-006

UIUX-013 Fold in G8 observations A8/A9
  ├─> depends on UIUX-007 for Reminder Center and Import/Export layout
  └─> depends on UIUX-001 for label/badge primitives

UIUX-014 E2E, regression, and handoff evidence
  └─> depends on UIUX-002..013
```

## Required Build Order

1. **Foundation:** UIUX-001.
2. **Global shell and dashboard:** UIUX-002.
3. **Cross-cutting state/permission scaffolding:** UIUX-008 and UIUX-009 can
   start after UIUX-001, then must be reconciled into page work.
4. **Page-type implementation:** UIUX-003, UIUX-004, UIUX-005, UIUX-006,
   UIUX-007.
5. **Quality cross-cuts:** UIUX-010, UIUX-011, UIUX-012 run across the completed
   page work and may feed small fixes back into UIUX-002..007.
6. **Folded observations:** UIUX-013 after the special pages have their new
   layout.
7. **Regression/evidence:** UIUX-014 last.

## Parallelization Guidance

- UIUX-008 and UIUX-009 may be implemented in parallel after UIUX-001, but they
  are not complete until all page tasks consume their primitives and rules.
- UIUX-006 reports and UIUX-007 special pages may proceed in parallel after
  UIUX-002, UIUX-008, and UIUX-009.
- UIUX-005 should not start until UIUX-003 and UIUX-004 establish the reusable
  list/detail/form patterns.
- UIUX-010/011/012 should be treated as continuous review tasks, but their final
  acceptance is after all page tasks are complete.

## Blocking Conditions

Any of the following blocks G9/G11 completion and returns to planning/audit:

- A task requires backend/API/data-model/business-logic changes to satisfy the
  visual target.
- A locked token is missing or the implementation needs a new color/token.
- A role gate cannot be represented without weakening a server-side invariant.
- A page cannot preserve zh-CN without changing enum/role comparison values.
- Existing e2e coverage would need to be skipped or weakened to pass.
- Any page type remains unmapped to the 14 navigation pages or 9 locked mockups.

