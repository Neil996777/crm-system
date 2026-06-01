# Acceptance → Task Map (ACC-001..023, 23/23 coverage proof)

Every ACC item maps to ≥1 TASK. Foundation tasks map to the ACC they primarily
enable; each capability ACC has at least one backend+frontend task pair (plus
tests). No ACC is unmapped; no P0/P1 ACC is satisfied by a foundation task alone
where a vertical slice is required. Cross-references resolve the PSM Traceability
`Task ID = pending (G8)` placeholders.

| ACC | Priority | Title (current model) | Primary CAP | Owning SVC | Tasks (primary first) |
|---|---|---|---|---|---|
| ACC-001 | P0 | Log in / operate under assigned role | CAP-001 | SVC-002 | TASK-002 (backend auth/session), TASK-006 (frontend sign-in/shell) |
| ACC-002 | P0 | Enforce three-role access control; no hard delete | CAP-001/CAP-012 | SVC-002 (+all) | TASK-003 (permission/scope/S2S/last-admin), TASK-029 (user/role admin UI), TASK-032 (archive governance + no hard delete) |
| ACC-003 | P0 | Manage leads (owner/source/info/status) | CAP-002 | SVC-003 | TASK-007 (backend CRUD/assign), TASK-009 (frontend) |
| ACC-004 | P0 | Qualify leads, preserve history | CAP-002 | SVC-003 | TASK-008 (backend qualify/convert), TASK-009 (frontend) |
| ACC-005 | P0 | Manage companies/customers + status | CAP-003 | SVC-004 | TASK-010 (backend), TASK-012 (frontend) |
| ACC-006 | P0 | Manage multiple contacts under a company | CAP-003 | SVC-004 | TASK-011 (backend), TASK-012 (frontend) |
| ACC-007 | P0 | Manage opportunities (Stage only, no Status) | CAP-004 | SVC-005 | TASK-013 (backend create/edit), TASK-016 (frontend) |
| ACC-008 | P0 | Move opportunities through pipeline | CAP-004 | SVC-005 | TASK-014 (backend stage transitions), TASK-016 (frontend) |
| ACC-009 | P0 | Manage quote (exactly one per opportunity) | CAP-005 | SVC-006 | TASK-017 (backend quote lifecycle), TASK-021 (frontend) |
| ACC-010 | P0 | Manage contracts | CAP-005 | SVC-006 | TASK-018 (backend create from Accepted quote), TASK-019 (backend lifecycle), TASK-022 (frontend) |
| ACC-011 | P0 | Manage payment plans + actual payments | CAP-005 | SVC-006 | TASK-020 (backend payments/overpayment), TASK-023 (frontend) |
| ACC-012 | P0 | Activities, notes, follow-up tasks | CAP-006 | SVC-007 | TASK-024 (backend), TASK-025 (frontend) |
| ACC-013 | P0 | Close opportunities Won/Lost, preserve history | CAP-004/CAP-005 | SVC-005 | TASK-015 (backend close Won/Lost), TASK-016 (frontend), TASK-019 (contract Signed gate), TASK-020 (payment decoupled) |
| ACC-014 | P0 | Review record-local history | CAP-008 | SVC-008 | TASK-004 (backend append/query), TASK-027 (frontend timeline), TASK-038 (classification/masking) |
| ACC-015 | P0 | List/detail/search/basic filter | CAP-007 | SVC-001 | TASK-030 (gateway aggregation + UI list/detail), TASK-005 (gateway spine) |
| ACC-016 | P0 | Persist all core CRM data | CAP-011 | SVC-002..010 | TASK-037 (persistence verification), TASK-001 (persistence platform), TASK-040 (backup durability) |
| ACC-017 | P0 | Deploy/operate with real config + data | CAP-011 | runtime host | TASK-039 (deploy/HTTPS/security/monitoring), TASK-040 (off-server backup + restore rehearsal) |
| ACC-018 | P1 | Manager team overview | CAP-009 | SVC-009 | TASK-033 (backend overview + frontend) |
| ACC-019 | P1 | Duplicate warnings (company/contact/lead) | CAP-002/CAP-003 | SVC-003/SVC-004 | TASK-031 (duplicate warning backend + UI) |
| ACC-020 | P1 | CSV import/export | CAP-010 | SVC-010 | TASK-035 (import), TASK-036 (export) |
| ACC-021 | P1 | In-app reminders (tasks/contracts/payments) | CAP-006/CAP-012 | SVC-007 | TASK-026 (reminder on-read + Reminder Center) |
| ACC-022 | P1 | Admin/global operation logs | CAP-008 | SVC-008 | TASK-028 (oplog query + UI), TASK-004 (append/store) |
| ACC-023 | P1 | Basic sales reports | CAP-009 | SVC-009 | TASK-034 (reports backend + frontend), TASK-033 (projection infra) |

## Coverage Summary

- All 23 ACC items (ACC-001..023) map to ≥1 task. 23/23.
- Each P0/P1 capability ACC has a backend persistence/contract task and a frontend
  screen task (where user-visible), plus the tests named in each task.
- Foundation tasks map to the ACC they primarily enable: TASK-001→ACC-016 platform,
  TASK-002→ACC-001, TASK-003→ACC-002, TASK-004→ACC-014/022 core, TASK-005→ACC-015 spine.
- Deployment/release-evidence: TASK-039 and TASK-040 carry ACC-017 (and reinforce
  ACC-016 durability); their ARCH-ACC-004/008/013/014/015 are `Release-evidence
  pending`, proven at G11 and audited at G12.

## Reverse map (Task → primary ACC)

TASK-001→ACC-016 · TASK-002→ACC-001 · TASK-003→ACC-002 · TASK-004→ACC-014 ·
TASK-005→ACC-015 · TASK-006→ACC-001 · TASK-007→ACC-003 · TASK-008→ACC-004 ·
TASK-009→ACC-003 · TASK-010→ACC-005 · TASK-011→ACC-006 · TASK-012→ACC-005 ·
TASK-013→ACC-007 · TASK-014→ACC-008 · TASK-015→ACC-013 · TASK-016→ACC-008 ·
TASK-017→ACC-009 · TASK-018→ACC-010 · TASK-019→ACC-010 · TASK-020→ACC-011 ·
TASK-021→ACC-009 · TASK-022→ACC-010 · TASK-023→ACC-011 · TASK-024→ACC-012 ·
TASK-025→ACC-012 · TASK-026→ACC-021 · TASK-027→ACC-014 · TASK-028→ACC-022 ·
TASK-029→ACC-002 · TASK-030→ACC-015 · TASK-031→ACC-019 · TASK-032→ACC-002 ·
TASK-033→ACC-018 · TASK-034→ACC-023 · TASK-035→ACC-020 · TASK-036→ACC-020 ·
TASK-037→ACC-016 · TASK-038→ACC-014 · TASK-039→ACC-017 · TASK-040→ACC-017

Note: a few tasks deliver behavior for more than one ACC (e.g. TASK-009 covers
ACC-003 and ACC-004 lead/qualification UI; TASK-012 covers ACC-005 and ACC-006);
the table above lists each task's primary ACC per the per-task schema field 5,
while field-3 objectives note secondary ACC coverage.
