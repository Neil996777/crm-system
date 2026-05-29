# QA/TDD G8 Task Planning Review

## Decision

Blocked.

QA/TDD cannot recommend G8 pass yet. The `delivery/` task package is broadly structured as an end-to-end delivery plan and most tasks include automated test expectations, manual verification, negative cases, TDD guards, blocker triggers, and MDA trace strings. However, there are unresolved P0/P1 quality-planning gaps that prevent the G8 criterion "tasks are traceable, testable, end-to-end, and no P0/P1 item is uncovered" from being met.

## Reviewed Inputs

- `delivery/tasks.md`
- `delivery/task-dependencies.md`
- `delivery/delivery-plan.md`
- `delivery/acceptance-task-map.md`
- `delivery/blockers.md`
- `modeling/test-model.md`
- `modeling/traceability-matrix.md`
- `docs/product/acceptance-matrix.md`
- `docs/qa/test-plan.md`
- `docs/qa/test-cases.md`
- Workspace standards: `company/operating-model.md`, `standards/acceptance-matrix-standard.md`, `standards/status-and-priority-standard.md`
- QA role rules: `agents/qa-tdd.md`, `agents/crm-qa-tdd-owner.md`

## Findings

| ID | Severity | Area | Finding | Evidence | Required Action |
|---|---|---|---|---|---|
| QA-G8-001 | P0 Blocker | ACC-012 / TASK-012 test coverage | ACC-012 requires activity, note, and task scenarios for lead/customer/contact/opportunity/quote/contract/payment records, but TASK-012 and TM-012 only make lead/customer/opportunity/contract/payment concrete. Contact and quote related-record coverage is missing from planned files/tests/verification, so a P0 acceptance path is uncovered. | `docs/product/acceptance-matrix.md` ACC-012 requires lead/customer/contact/opportunity/quote/contract/payment. `delivery/tasks.md` TASK-012 modifies only lead/company/opportunity/contract/payment pages and tests lead/customer/opportunity/contract/payment contexts. `modeling/test-model.md` TM-012 repeats the same narrower coverage. | Update TASK-012 and TM-012 planning to explicitly include contact and quote contexts in production surfaces, automated API/E2E tests, manual verification, and related-record permission negative cases. |
| QA-G8-002 | P1 Blocker | QA test plan and test cases | Active QA artifacts do not provide a concrete G9/G10 test baseline. `docs/qa/test-plan.md` has `TBD` scope/test data/risks and only ACC-001. `docs/qa/test-cases.md` has one placeholder TEST-001 with TBD scenario/expected result. This does not map P0/P1 acceptance to concrete test cases or future command/evidence expectations. | `docs/qa/test-plan.md` lines 5, 11, 23, and 29; `docs/qa/test-cases.md` line 5. | Before G8 pass, either update QA artifacts or create an accepted G8 QA coverage artifact that maps ACC-001..ACC-023 / TM-001..TM-023 / TASK-001..TASK-023 to concrete test case groups, expected automation layers, manual evidence, and G9/G10 result capture. |
| QA-G8-003 | P2 Improvement | Automation evidence handoff | `delivery/tasks.md` lists planned test files and test categories for all 23 tasks, but it does not define the expected command groups or result artifacts that implementation agents must provide at G9/G10. `delivery/delivery-plan.md` asks for commands run, but the command baseline is still implicit. | `delivery/tasks.md` has 23 `Automated tests` rows. `delivery/delivery-plan.md` requires commands/results in handoff. | Add a compact command/evidence matrix for backend unit/integration, frontend/component, E2E, operational smoke, and manual evidence capture once the concrete toolchain is chosen. |
| QA-G8-004 | P2 Improvement | Dependency precision | Some dependencies and file-change scopes use shorthand ranges such as `TASK-003 to TASK-013`. This is readable, but weaker for QA evidence tracking and blocker routing than explicit dependency IDs. | `delivery/task-dependencies.md` and `delivery/tasks.md` use task ranges for TASK-014, TASK-015, TASK-016, and TASK-017. | Expand range shorthand in evidence-critical rows to explicit task IDs so QA, Integration, and Audit can verify coverage mechanically. |
| QA-G8-005 | P2 Improvement | Manual verification evidence template | Each task has reproducible manual steps and expected results, which is good for G8. The later evidence shape should still require actual result, environment, date, actor role, data fixture/record IDs, and artifact links/screenshots where applicable. | `delivery/tasks.md` has 23 `Manual verification` rows; `delivery/delivery-plan.md` lists handoff evidence requirements. | Add a reusable manual evidence template for G9/G10 task closure to prevent incomplete manual evidence. |

## P0/P1 Blockers

| Blocker ID | Severity | Affected Items | Owner | Blocking Condition | Required Resolution |
|---|---|---|---|---|---|
| QA-G8-BLOCK-001 | P0 | ACC-012, TM-012, TASK-012 | Task Planner + Domain Modeling + QA TDD | Contact and quote related-record activity/note/task coverage is not concretely planned. | Repair TASK-012 and TM-012 coverage so contact and quote contexts are included in automated tests, manual verification, scope checks, and planned UI/API file impacts. |
| QA-G8-BLOCK-002 | P1 | ACC-001..ACC-023, TM-001..TM-023, TASK-001..TASK-023 | QA TDD + Task Planner | QA test plan/test cases remain placeholder-only and do not provide concrete G9/G10 test baseline. | Replace TBD QA plan/cases with concrete acceptance-to-test/task mappings and evidence expectations, or record an accepted equivalent QA coverage artifact for G8 handoff. |

## P2 Improvements

| ID | Improvement |
|---|---|
| QA-G8-P2-001 | Add expected test command groups and result artifact names for unit, integration, E2E, operational smoke, and manual verification. |
| QA-G8-P2-002 | Replace task-range shorthand in dependency/evidence-critical rows with explicit task IDs. |
| QA-G8-P2-003 | Add a manual evidence template requiring environment, role, starting state, steps, expected result, actual result, command/output or screenshot links, and blocker status. |

## Positive Coverage Notes

- `delivery/tasks.md` contains 23 task sections and 23 rows each for automated tests, manual verification, TDD quality guard, MDA trace, planned file changes, and blocker records.
- `delivery/acceptance-task-map.md` maps ACC-001..ACC-023 to TASK-001..TASK-023 and TM-001..TM-023.
- Core no-downgrade controls are present, including explicit rejection of mock, static, localStorage-only, in-memory-only, and frontend-only authorization paths.
- Negative cases are generally present across identity, authorization, lifecycle transitions, money, persistence, import/export, reminders, reports, and operation logs.

## Recommendation

Do not pass G8 from QA/TDD yet.

Repair the two P0/P1 blockers, then re-run QA/TDD G8 review. After repair, if no additional coverage gaps appear, QA/TDD can likely recommend G8 pass with P2 improvements carried into implementation and G9/G10 evidence discipline.
