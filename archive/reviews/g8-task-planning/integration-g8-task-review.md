# Integration Owner G8 Task Planning Review

## Decision

Blocked.

Integration Owner does not recommend G8 pass yet. The delivery artifacts cover
all ACC/TM/TASK IDs, but the current task ordering and integration blocker
routing contain executable-sequencing defects that would prevent reproducible
end-to-end verification for P0/P1 scope.

## Reviewed Inputs

- `delivery/tasks.md`
- `delivery/task-dependencies.md`
- `delivery/delivery-plan.md`
- `delivery/acceptance-task-map.md`
- `delivery/blockers.md`
- `docs/integration/integration-report.md`
- `docs/integration/acceptance-status.md`
- `docs/integration/blocker-list.md`
- `modeling/test-model.md`
- `modeling/traceability-matrix.md`
- G8 pass criteria from `../../company/operating-model.md`
- Integration Owner role rules from `../../agents/integration-owner.md`
- CRM integration agent rules from `agents/crm-integration-operations-delivery.md`

## Findings

| ID | Severity | Finding | Evidence | Impact | Required Action / Routing |
|---|---|---|---|---|---|
| INT-G8-001 | P0 Blocker | Record-local history is required by earlier P0 task completion standards, but the history task is sequenced after those tasks. | `TASK-003` requires owner-change history in its completion standard, and `TASK-004` requires conversion history. `TASK-014` is the first task that adds record-local history infrastructure and modifies `TASK-003` to `TASK-013` mutation services, but it depends on `TASK-003 to TASK-013`. | Earlier P0 tasks cannot be completed with their own stated evidence before the later history task exists. This breaks executable order and makes integration closure unreproducible for ACC-003, ACC-004, ACC-008, ACC-013, ACC-014, and downstream ACC-016. | Route to Task Planner with Domain Modeling and Architecture. Either make the history/audit write path an earlier explicit dependency before closing mutation tasks, or define a staged closure rule where mutation tasks remain not complete until `TASK-014` integration evidence is added. |
| INT-G8-002 | P0 Blocker | `TASK-017` operational restore evidence requires reports and reminders, but its prerequisites omit the report/reminder tasks. | `TASK-017` completion standard requires restore rehearsal to prove login, disabled denial, core records, history, logs, reports, and reminders. Its prerequisites are only `TASK-001 to TASK-016, TASK-022`; report and reminder capabilities are in `TASK-018`, `TASK-021`, and `TASK-023`. | ACC-017 is P0. As written, production-equivalent restore verification can start before required evidence-producing capabilities exist, so backup/restore smoke cannot be reproduced against the stated standard. | Route to Task Planner and CRM Integration/Operations Delivery. Add the missing operational evidence dependencies for `TASK-017`, or narrow the `TASK-017` restore standard only if Product/Modeling formally confirms no P0/P1 downgrade occurs. |
| INT-G8-003 | P1 Blocker | Import/export and global operation logs have a circular/unsatisfiable evidence relationship. | `TASK-020` depends on `TASK-022` for operation-log evidence. `TASK-022` requires auth, owner/stage/status, quote, contract, payment, archive, import, and export services to write operation events and its manual verification includes import/export. | `TASK-022` cannot fully verify import/export operation-log events before import/export exists, while `TASK-020` cannot close before `TASK-022`. This blocks ACC-020 and ACC-022 integration sequencing and may later affect ACC-017 restore evidence. | Route to Task Planner and crm-backend-operations-reporting. Split log framework/read access from per-feature log event completion, or sequence `TASK-020` and `TASK-022` with explicit partial/combined integration closure rules that do not mark either Done prematurely. |
| INT-G8-004 | P1 Blocker | Active integration docs contain placeholder open blockers and incomplete acceptance status that conflict with the delivery blocker register. | `delivery/blockers.md` says no open P0/P1 blocker is known for G8 review. `docs/integration/blocker-list.md` contains `INT-BLK-001` with severity `P0/P1/P2/P3`, blocker `TBD`, status `Open`; `docs/integration/acceptance-status.md` covers only `ACC-001`; `docs/integration/integration-report.md` has `TBD` issue rows. | Integration blocker routing is ambiguous. If `docs/integration/*` are active, they currently assert an open unresolved blocker without owner, acceptance meaning, or severity. If they are only placeholders, they need to be labeled as such so G8 blocker status is not contradictory. | Route to Integration Owner and Task Planner. Before G8 pass, either archive/label these as pre-G11 placeholders or replace placeholders with a full ACC-001 to ACC-023 pending template and remove fake open blocker rows. |
| INT-G8-005 | P2 Issue | Planned integration evidence is present for every ACC item, but evidence artifact naming and reproducible manual evidence capture are not yet standardized. | `acceptance-task-map.md` records planned integration evidence for ACC-001 to ACC-023, and `delivery-plan.md` requires commands, manual environment, and screenshots/artifacts in handoff. Individual manual steps are mostly scenario-level and do not define evidence IDs, seed fixture names, reset commands, or output artifact paths. | This is not the primary G8 blocker because the task plan does identify manual paths and tests, but it will slow G10/G11 evidence collection and make audit reverse-verification harder. | Route to CRM Integration/Operations Delivery and QA TDD. Add a G10/G11 evidence convention before implementation closure: `ACC`, `TASK`, environment, seed fixture, command output, screenshots/logs, backup artifact, restore artifact, and blocker IDs. |

## P0/P1 Blockers

| Blocker ID | Severity | Affected Acceptance | Affected Tasks | Owner Route | Required Resolution Before G8 Pass |
|---|---|---|---|---|---|
| INT-G8-001 | P0 | ACC-003, ACC-004, ACC-008, ACC-013, ACC-014, ACC-016 | TASK-003, TASK-004, TASK-008, TASK-013, TASK-014, TASK-016 | Task Planner + Domain Modeling + Architecture | Repair history dependency sequencing so no P0 task claims completion before required transaction-linked history evidence can exist. |
| INT-G8-002 | P0 | ACC-017 | TASK-017, TASK-018, TASK-021, TASK-023 | Task Planner + CRM Integration/Operations Delivery | Align `TASK-017` prerequisites with every capability its restore/smoke standard must prove. |
| INT-G8-003 | P1 | ACC-020, ACC-022 | TASK-020, TASK-022 | Task Planner + crm-backend-operations-reporting | Break the import/export and operation-log evidence cycle with explicit sequencing or combined closure rules. |
| INT-G8-004 | P1 | ACC-001 to ACC-023 blocker routing | Integration docs and `delivery/blockers.md` | Integration Owner + Task Planner | Remove contradictory placeholder open blocker rows or mark integration docs as non-current templates before using them as G8 inputs. |

## P2 Improvements

| ID | Improvement | Recommendation |
|---|---|---|
| INT-G8-P2-001 | Evidence artifact convention | Define evidence IDs and artifact paths for each ACC/TASK before G10/G11 so Integration and Audit can reproduce the same proof set. |
| INT-G8-P2-002 | Manual verification precision | Add fixture/reset references and exact environment assumptions to each manual verification path during implementation handoff. |
| INT-G8-P2-003 | Operational watch items | Keep `WATCH-001` and `WATCH-002` promoted-ready for ACC-017, because seed data, target host, backup target, restore environment, and domain decisions become blockers before production-equivalent evidence is accepted. |

## Recommendation

Do not pass G8 from the Integration Owner review until the P0/P1 blockers above
are repaired in the delivery plan and blocker routing is made unambiguous.

After repair, re-review should confirm:

- no P0/P1 task depends on evidence produced only by a later task unless a
  formal combined closure rule prevents premature Done status;
- `TASK-017` cannot start operational evidence before every capability required
  by its restore/smoke standard exists;
- integration blocker/status artifacts no longer contain fake open blocker
  rows or one-row placeholder acceptance coverage;
- ACC-001 to ACC-023 remain mapped to concrete tasks, tests, manual
  verification, integration evidence targets, and audit-ready traceability.
