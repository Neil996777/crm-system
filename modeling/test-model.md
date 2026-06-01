# Test Model

The platform-independent test design that maps the accepted acceptance spine,
domain behaviors, state machines, invariants, edge cases, permission rules, and
abuse cases to concrete `TEST-*` test concepts (what to verify, test type, and
expected result) — so QA Execution (G9/G10) can implement them without redesign.
This is test DESIGN, not test implementation: no test code and no production code.

## Document Control

- Project: CRM System
- Date: 2026-06-01
- Role: QA Test Design
- Gate: G6 (MDA Modeling) — final MDA-package artifact
- Scope note: This artifact is test DESIGN only. It maps acceptance items,
  behaviors, state transitions, permissions, invariants, edges, and abuse cases
  to `TEST-*` concepts with test type and expected result. Test implementation
  and execution are G9/G10 (QA Execution). It maps FROM the accepted
  `acceptance-matrix.md` (ACC-001..023), the accepted `PIM.md`
  (PIM-SM-001..011, PIM-INV-001..052, PIM-BEH-001..034 and their TEST-* concept
  tags), `CIM.md`, `PSM.md` (SVC-001..010, CONTRACT-001..020, FLOW-001..013, and
  the architecture-resolved on-read overdue evaluation), and the accepted G4
  business/security design (`business-rules.md` BR-*, `edge-cases.md` EDGE-*,
  `permission-matrix.md` PM-*, `role-permission-scenarios.md` PERM-*,
  `abuse-cases.md` ABUSE-*, `security-requirements.md` SEC-*, `decision-log.md`
  DEC-*). It invents no expected result for product behavior that the accepted
  sources leave undecided.
  Amended 2026-06-01 per Formal Scope Change by User (decision-log.md
  DEC-017..020): Won = related contract Signed (not full payment, DEC-017);
  exactly one quote per opportunity (DEC-018); payment tracking retained but
  decoupled from Won (DEC-019); Opportunity Status field removed, Pipeline Stage is
  the sole lifecycle dimension (DEC-020); `Contract Signed` and `Payment In
  Progress` opportunity stages removed. This change set resolved all three prior G6
  blockers (BLK-001 by DEC-020, BLK-002 by DEC-017/019, BLK-003 by DEC-018); the
  three formerly-`Pending` blocker tests are retired in place and no `Pending` test
  remains. Affected test IDs are retired in place (marked, not renumbered) to
  preserve cross-references.

### Test-type convention (PSM-informed)

- **Unit** — single-aggregate guard/invariant logic inside one service boundary
  (e.g. overpayment ceiling inside the Contract aggregate, SVC-006).
- **Integration** — service contract/API + persistence + permission + emitted
  history/operation-log event, and cross-service flows (FLOW-001..013). Per
  `abuse-cases.md` and SEC-017, P0 security/business-rule paths require
  backend/API negative tests, not only UI.
- **E2E** — user-visible path across gateway-bff (SVC-001) through to persisted
  result and back (lists, detail, reminders, reports).
- **Manual** — exploratory/visual/operational confirmation where automation is
  not the primary evidence (e.g. deployment smoke, ACC-017).

`TEST-*` IDs reuse the PIM behavior concept tags (TEST-AUTH-LOGIN,
TEST-AUTHZ-SCOPE, TEST-OPP-CLOSE, …) as the family name and number concrete cases
within the family (e.g. `TEST-OPP-CLOSE-001`). Negative/reject cases are numbered
within the same family. All statuses are `Draft` (G6 design) unless marked
`Pending` (blocked on a product/BA decision).

---

## Acceptance To Test Mapping

Every ACC-001..023 maps to ≥1 `TEST-*`. P0 items carry both positive and negative
coverage. Multiple tests are listed where the ACC carries several behaviors.

| Acceptance ID | Rule / Behavior | Test Type | Test ID(s) | Status |
|---|---|---|---|---|
| ACC-001 (P0) | Login binds assigned role to session; invalid creds / disabled / unauthenticated denied; session persists (BR-001, EDGE-001, PIM-BEH-001, SEC-001/AUTH-001..006, DEC-005) | Integration, E2E, Manual | TEST-AUTH-LOGIN-001 (valid login, role bound), TEST-AUTH-LOGIN-002 (invalid creds rejected, unified message), TEST-AUTH-LOGIN-003 (disabled user denied), TEST-AUTH-LOGIN-004 (unauthenticated protected route/API denied), TEST-AUTH-LOGIN-005 (session persists across refresh/re-login), TEST-AUTH-LOGIN-006 (stale session re-evaluated after disable/role change) | Draft |
| ACC-002 (P0) | Three-role scope on every protected action; deny without expose/mutate; no hard delete; Sales no global logs (BR-001/002, PIM-BEH-002, PM-001..029, PERM-001..012, SEC-002..006) | Integration, E2E, Manual | TEST-AUTHZ-SCOPE-001 (admin governs all), TEST-AUTHZ-SCOPE-002 (manager team scope), TEST-AUTHZ-SCOPE-003 (sales owned only), TEST-AUTHZ-SCOPE-004 (cross-owner view denied, no leakage), TEST-AUTHZ-SCOPE-005 (denied mutation creates no change), TEST-AUTHZ-SCOPE-006 (no hard delete any role) | Draft |
| ACC-003 (P0) | Lead CRUD/search/filter/assign; required fields; owner before Pending Qualification; Sales no non-owned; Unassigned restricted (BR-003/004, PIM-BEH-004/005, PIM-SM-001, EDGE-003/004) | Integration, E2E | TEST-LEAD-CREATE-001 (create w/ required fields persisted), TEST-LEAD-CREATE-002 (missing name/source/status blocked), TEST-LEAD-CREATE-003 (create Unassigned allowed), TEST-LEAD-ASSIGN-001 (manager assigns owner, history recorded), TEST-LEAD-ASSIGN-002 (Sales assign/transfer denied), TEST-AUTHZ-SCOPE-004 (Sales non-owned denied), TEST-NAV-RETRIEVE-002 (search/filter scoped) | Draft |
| ACC-004 (P0) | Qualify Valid/Invalid/Convert; preserve history; Invalid no convert until restored; converted once (BR-004, PIM-BEH-006, PIM-SM-001, EDGE-004/005/006) | Integration, E2E | TEST-LEAD-QUALIFY-001 (Pending→Valid w/ authz), TEST-LEAD-QUALIFY-002 (Pending→Invalid w/ reason), TEST-LEAD-QUALIFY-003 (Valid→Converted creates opportunity, preserves history), TEST-LEAD-QUALIFY-004 (Unassigned qualify rejected), TEST-LEAD-QUALIFY-005 (Invalid convert rejected until restored), TEST-LEAD-QUALIFY-006 (Invalid→Pending restore by admin/manager only), TEST-LEAD-QUALIFY-007 (re-convert rejected) | Draft |
| ACC-005 (P0) | Company/customer CRUD/search/filter; required fields; no hard delete; authz (BR-002/003, PIM-BEH-007, PIM-SM-010) | Integration, E2E | TEST-CUSTOMER-CRUD-001 (create w/ company name+status+owner), TEST-CUSTOMER-CRUD-002 (missing required fields blocked), TEST-CUSTOMER-CRUD-003 (edit persisted), TEST-CUSTOMER-CRUD-004 (Sales unrelated denied), TEST-INV-NODELETE-001 (hard delete unavailable) | Draft |
| ACC-006 (P0) | Multiple contacts under company; required link + contact method/role note; authz (BR-003, PIM-BEH-008, EDGE-007) | Integration, E2E | TEST-CONTACT-LINK-001 (create w/ company + method/note), TEST-CONTACT-LINK-002 (save without company blocked), TEST-CONTACT-LINK-003 (multiple contacts visible in company context), TEST-CONTACT-LINK-004 (unrelated Sales denied) | Draft |
| ACC-007 (P0) | Opportunity CRUD; required links/fields; persisted Stage (sole lifecycle dimension — DEC-020, no separate Status); authz; no hard delete (BR-003, PIM-BEH-009, PIM-SM-002) | Integration, E2E | TEST-OPP-CREATE-001 (create w/ customer/owner/stage/amount/close date, no Status field), TEST-OPP-CREATE-002 (missing required link/field blocked), TEST-OPP-CREATE-003 (plain opportunity create persisted; Stage is sole lifecycle dimension, no separate Status), TEST-OPP-CREATE-004 (non-owned edit denied) [TEST-OPP-STATUS-ENUM-001 RETIRED 2026-06-01 — resolved by DEC-020 (Opportunity Status field removed)] | Draft |
| ACC-008 (P0) | Stage transitions allowed/forbidden; Won/Lost terminal; Won requires Signed contract (DEC-017); Lost needs reason; history (BR-005/020, PIM-BEH-010/011, PIM-SM-002/009, EDGE-008/009/010/011/036) | Integration, E2E | TEST-OPP-STAGE-001 (allowed forward transition persisted + history), TEST-OPP-STAGE-002 (forbidden transition rejected, no mutation), TEST-OPP-STAGE-003 (arbitrary rollback rejected), TEST-OPP-CLOSE-001/002/003 (see closure), TEST-HISTORY-002 (stage-change history event) | Draft |
| ACC-009 (P0) | Quote lifecycle; required fields; exactly one quote per opportunity (DEC-018); expired no contract link; authz (BR-006, PIM-BEH-012/013, PIM-SM-004, EDGE-012/013, DEC-018) | Integration, E2E | TEST-QUOTE-LIFECYCLE-001 (create Draft w/ required fields), TEST-QUOTE-LIFECYCLE-002 (missing amount/status/validity blocked), TEST-QUOTE-LIFECYCLE-003 (send/reject/expire), TEST-QUOTE-ACCEPT-001 (accept the opportunity's single quote), TEST-QUOTE-ACCEPT-002 (the single quote is accepted; one-quote-per-opportunity enforced), TEST-INV-ONEACCEPT-001 (invariant: exactly one quote per opportunity), TEST-CONTRACT-CREATE-003 (expired quote contract link rejected) [TEST-QUOTE-ACCEPT-003 RETIRED 2026-06-01 — resolved by DEC-018 (one quote per opportunity, no second-accept)] | Draft |
| ACC-010 (P0) | Contract lifecycle from Accepted quote; note + expected signed date; signed/effective date for signed states; amount-diff reason; authz (BR-007, PIM-BEH-014/015/016, PIM-SM-005, EDGE-013/014/015/016/017, DEC-006/007/016) | Integration, E2E | TEST-CONTRACT-CREATE-001 (Pending Signature w/ note+expected signed date, no signed date), TEST-CONTRACT-CREATE-002 (missing note/link/amount/expected signed date rejected), TEST-CONTRACT-CREATE-003 (expired/non-Accepted quote link rejected), TEST-CONTRACT-LIFECYCLE-001 (sign/activate/complete require signed date), TEST-CONTRACT-LIFECYCLE-002 (signed state w/o signed date rejected), TEST-CONTRACT-LIFECYCLE-003 (terminate pre/post signature), TEST-CONTRACT-AMOUNT-DIFF-001 (diff reason required), TEST-INV-NOAPPROVAL-001 (no approval/e-sign/template) | Draft |
| ACC-011 (P0) | Payment plan + actual payment; partial/full status; reject zero/neg/overpayment; overdue; single currency; authz (BR-008, PIM-BEH-017/018/019, PIM-SM-006, EDGE-018/019/020, DEC-013/014) | Integration, E2E, Unit | TEST-PAYMENT-RECORD-001 (plan create Unpaid), TEST-PAYMENT-RECORD-002 (partial → Partially Paid), TEST-PAYMENT-RECORD-003 (full single plan → Paid), TEST-PAYMENT-GUARD-001 (zero rejected), TEST-PAYMENT-GUARD-002 (negative rejected), TEST-PAYMENT-GUARD-003 (contract overpayment rejected), TEST-PAYMENT-GUARD-004 (single currency), TEST-PAYMENT-OVERDUE-001 (overdue on business date) [TEST-PAYMENT-FULLPAID-AGG-001 RETIRED 2026-06-01 — resolved by DEC-017/019 (Won = contract Signed; payment decoupled from Won, no full-payment aggregation gate)] | Draft |
| ACC-012 (P0) | Activity/note/task against related records; required fields; completed/cancelled not active reminders; authz (BR-003/009, PIM-BEH-020/021, PIM-SM-007, EDGE-023/024) | Integration, E2E | TEST-ACTIVITY-NOTE-001 (activity/note on lead/customer/opportunity/contract/payment), TEST-ACTIVITY-NOTE-002 (missing related record/fields blocked), TEST-TASK-LIFECYCLE-001 (create Open), TEST-TASK-LIFECYCLE-002 (complete/cancel), TEST-TASK-LIFECYCLE-003 (overdue on business date), TEST-TASK-LIFECYCLE-004 (completed/cancelled not reminder), TEST-ACTIVITY-NOTE-003 (permission-denied create) | Draft |
| ACC-013 (P0) | Close Won (related contract Signed — DEC-017) / Lost (reason); preserve history; terminal non-reopen (BR-005/020, PIM-BEH-011, PIM-SM-009, EDGE-009/010/011/036, DEC-017) | Integration, E2E | TEST-OPP-CLOSE-001 (Won after related contract Signed, history preserved), TEST-OPP-CLOSE-002 (Won rejected without a Signed contract), TEST-OPP-CLOSE-003 (Lost with reason), TEST-OPP-CLOSE-004 (Lost without reason rejected), TEST-OPP-CLOSE-005 (reopen/rollback/re-close rejected), TEST-OPP-CLOSE-006 (post-close notes/tasks allowed, stage edit rejected) | Draft |
| ACC-014 (P0) | Record-local history visible by record permission; non-editable (BR-010, PIM-BEH-028, PM-024/025, PERM-017/018, EDGE-029, SEC-009) | Integration, E2E | TEST-HISTORY-001 (history visible to owner w/ actor/event/resource/timestamp/before-after), TEST-HISTORY-002 (events emitted for owner/stage/quote/contract/payment/task/archive), TEST-HISTORY-003 (non-owned history denied, no leakage), TEST-HISTORY-004 (history not editable via normal CRM actions) | Draft |
| ACC-015 (P0) | List/detail/search/basic filter across P0 entities; empty/invalid-filter/permission-hidden states (PIM-BEH-030, PM-008..033, EDGE-002/031) | E2E, Manual | TEST-NAV-RETRIEVE-001 (list+detail for all P0 entities), TEST-NAV-RETRIEVE-002 (search/filter happy path), TEST-NAV-RETRIEVE-003 (empty-state), TEST-NAV-RETRIEVE-004 (invalid filter feedback), TEST-NAV-RETRIEVE-005 (unauthorized records hidden), TEST-NAV-RETRIEVE-006 (archived excluded from active views) | Draft |
| ACC-016 (P0) | All core CRM data persists across refresh/re-login/restart; failed saves surfaced (BR-015, PIM-BEH-033, DEC-008, EDGE-030, SEC-018) | Integration, E2E, Manual | TEST-PERSISTENCE-001 (data survives refresh), TEST-PERSISTENCE-002 (survives logout/login), TEST-PERSISTENCE-003 (survives service restart), TEST-PERSISTENCE-004 (failed save surfaced, not false success), TEST-PERSISTENCE-005 (no mock/in-memory-only path satisfies P0) | Draft |
| ACC-017 (P0) | Deploy and operate with real config + persisted data; reachable/configured (BR-015, NFR-003, DEC-004, PSM-014, ARCH-ACC-004/008/013/014/015) | Manual, Integration | TEST-DEPLOY-SMOKE-001 (deployed CRM reachable+configured, persistent services connected), TEST-DEPLOY-SMOKE-002 (misconfiguration/unavailable dependency blocks readiness). Release-evidence items (HTTPS/TLS, off-server backup+restore, security-group, monitoring) verified at G11/G12 per ARCH-ACC carried blockers. | Draft |
| ACC-018 (P1) | Manager team overview; deny Sales; empty state (BR-014, PIM-BEH-031, PM-044/045, PERM-019/020, EDGE-028) | E2E, Manual | TEST-TEAM-OVERVIEW-001 (manager sees team leads/opps/quotes/contracts/payments/tasks/pipeline), TEST-TEAM-OVERVIEW-002 (Sales denied), TEST-TEAM-OVERVIEW-003 (empty-state), TEST-TEAM-OVERVIEW-004 (non-team records excluded) | Draft |
| ACC-019 (P1) | Duplicate warning on company/contact/lead match; non-blocking; no merge (BR-011/019, PIM-BEH-025, PM-048, EDGE-025/033/034/035) | Integration, E2E | TEST-DUPLICATE-WARN-001 (exact company name match), TEST-DUPLICATE-WARN-002 (contact phone match normalized), TEST-DUPLICATE-WARN-003 (contact email match normalized), TEST-DUPLICATE-WARN-004 (lead company/contact match), TEST-DUPLICATE-WARN-005 (proceed-after-warning creates record, no merge/overwrite), TEST-DUPLICATE-WARN-006 (unique data → no warning) | Draft |
| ACC-020 (P1) | CSV import/export; row-level errors; no corruption; authorized records only; Sales denied (BR-012/018, PIM-BEH-026/027, PM-034..039, PERM-024..026, EDGE-026/027) | Integration, E2E | TEST-CSV-IMPORT-001 (valid rows imported w/ summary), TEST-CSV-IMPORT-002 (mixed valid/invalid → invalid reported, existing not corrupted), TEST-CSV-IMPORT-003 (unsupported format rejected), TEST-CSV-IMPORT-004 (Sales import denied), TEST-CSV-EXPORT-001 (export authorized records only, archived excluded by default), TEST-CSV-EXPORT-002 (Sales export denied) | Draft |
| ACC-021 (P1) | In-app reminders for due/overdue tasks, pending-signature contracts past expected signed date, due/overdue payments; inactive suppressed; permission-filtered (BR-013/021, PIM-BEH-019/022, PIM-SM-005/006/007, EDGE-020/021/022/023, DEC-015) | Integration, E2E | TEST-REMINDER-001 (due/overdue task reminder), TEST-REMINDER-002 (pending-signature contract past expected signed date reminder), TEST-REMINDER-003 (due/overdue payment reminder), TEST-REMINDER-004 (completed/cancelled task, signed/terminated/fully-paid contract suppressed), TEST-REMINDER-005 (unauthorized reminders hidden) | Draft |
| ACC-022 (P1) | Admin global operation logs; required events; non-editable; Sales/Manager denied (BR-010, NFR-005, PIM-BEH-003/029, PM-040..042, PERM-021..023, EDGE-029, SEC-010) | Integration, Manual | TEST-OPLOG-001 (admin sees global events: login/access failures, owner/stage changes, quote accept, contract changes, payments, archive, import, export), TEST-OPLOG-002 (event has id/actor/action/resource/timestamp/result/before-after), TEST-OPLOG-003 (Manager denied), TEST-OPLOG-004 (Sales denied), TEST-OPLOG-005 (logs not editable via normal CRM actions), TEST-USER-ADMIN-001..004 (see invariant/permission) | Draft |
| ACC-023 (P1) | Basic reports — counts/sums per committed groupings over authorized persisted records; empty state; Sales denied (BR-014/017, PIM-BEH-032, PM-043/044/045, PERM-027..029, EDGE-028/031) | Integration, E2E | TEST-BASIC-REPORT-001 (leads by status / opps by stage / quotes / contracts / payments groupings traceable to records), TEST-BASIC-REPORT-002 (empty data → zero/empty state), TEST-BASIC-REPORT-003 (Sales denied), TEST-BASIC-REPORT-004 (unauthorized records excluded from aggregate), TEST-BASIC-REPORT-005 (archived excluded by default) | Draft |

---

## State Transition Tests

For each PIM-SM-001..011: representative valid transitions (guard satisfied) AND
forbidden/guard-failure transitions that MUST be rejected without data mutation
(PIM-INV-006, EDGE-008). Type is Integration unless noted (single-aggregate guard
logic is Unit-eligible).

| State Machine | From | To | Trigger | Expected (incl. reject cases) | Test ID |
|---|---|---|---|---|---|
| PIM-SM-001 Lead | (none) | Unassigned | Create unassigned | Persisted; name/source/status present | TEST-LEAD-CREATE-003 |
| PIM-SM-001 Lead | (none) | Pending Qualification | Create w/ owner | Persisted; owner present | TEST-LEAD-CREATE-001 |
| PIM-SM-001 Lead | Unassigned | Pending Qualification | Assign owner | Persisted; authorized actor | TEST-LEAD-ASSIGN-001 |
| PIM-SM-001 Lead | Pending Qualification | Valid | Qualify valid | Persisted; not Unassigned | TEST-LEAD-QUALIFY-001 |
| PIM-SM-001 Lead | Pending Qualification | Invalid | Mark invalid | Persisted; reason recorded | TEST-LEAD-QUALIFY-002 |
| PIM-SM-001 Lead | Invalid | Pending Qualification | Restore | Allowed only admin/manager; **Sales restore rejected** | TEST-LEAD-QUALIFY-006 |
| PIM-SM-001 Lead | Valid | Converted | Convert | Opportunity link created, history preserved | TEST-LEAD-QUALIFY-003 |
| PIM-SM-001 Lead | Unassigned | (rejected) | Qualify/convert | **Rejected, no mutation** (PIM-INV-001) | TEST-LEAD-QUALIFY-004 |
| PIM-SM-001 Lead | Invalid | (rejected) | Convert without restore | **Rejected** (PIM-INV-002) | TEST-LEAD-QUALIFY-005 |
| PIM-SM-001 Lead | Converted | (rejected) | Convert again | **Rejected** (PIM-INV-003) | TEST-LEAD-QUALIFY-007 |
| PIM-SM-002 Opp Stage | (none) | New Opportunity | Create | Persisted w/ required fields | TEST-OPP-CREATE-001 |
| PIM-SM-002 Opp Stage | New Opportunity → Needs Confirmed → Quote → Contract Negotiation | next | Advance | Each allowed forward transition persisted + history (Contract Signed / Payment In Progress stages removed — DEC-017) | TEST-OPP-STAGE-001 |
| PIM-SM-002 Opp Stage | any | (rejected) | Forbidden transition | **Rejected, no mutation** (PIM-INV-006) | TEST-OPP-STAGE-002 |
| PIM-SM-002 Opp Stage | later stage | (rejected) | Arbitrary rollback | **Rejected** (PIM-INV-006) | TEST-OPP-STAGE-003 |
| PIM-SM-003 Opp Status | — | — | — | **[RETIRED 2026-06-01 — resolved by DEC-020]** Opportunity Status dimension removed; Pipeline Stage (PIM-SM-002) is the sole lifecycle dimension. TEST-OPP-CREATE-003 re-aimed to plain opportunity create (no Status); TEST-OPP-STATUS-ENUM-001 retired. | — |
| PIM-SM-004 Quote | (none) | Draft | Create | Persisted w/ required fields | TEST-QUOTE-LIFECYCLE-001 |
| PIM-SM-004 Quote | Draft | Sent | Send | Authorized | TEST-QUOTE-LIFECYCLE-003 |
| PIM-SM-004 Quote | Sent | Accepted | Accept | Opportunity's single quote accepted (exactly one quote per opportunity — DEC-018) | TEST-QUOTE-ACCEPT-001 |
| PIM-SM-004 Quote | Sent | Rejected | Reject | Persisted | TEST-QUOTE-LIFECYCLE-003 |
| PIM-SM-004 Quote | Draft/Sent | Expired | Expire | Validity end passed | TEST-QUOTE-LIFECYCLE-003 |
| PIM-SM-004 Quote | — | — | (second quote) | **[Row removed 2026-06-01 — resolved by DEC-018]** Exactly one quote per opportunity; no second quote exists to accept, so no second-accept transition. TEST-QUOTE-ACCEPT-002 re-aimed to "the single quote is accepted"; TEST-QUOTE-ACCEPT-003 retired. | TEST-QUOTE-ACCEPT-002 |
| PIM-SM-005 Contract | (none) | Pending Signature | Create | Note + expected signed date present; Accepted non-Expired quote; **expired quote link rejected** | TEST-CONTRACT-CREATE-001 / -003 |
| PIM-SM-005 Contract | Pending Signature | Signed | Sign | Signed/effective date present; **missing date rejected** | TEST-CONTRACT-LIFECYCLE-001 / -002 |
| PIM-SM-005 Contract | Signed | Active | Activate | Signed/effective date present | TEST-CONTRACT-LIFECYCLE-001 |
| PIM-SM-005 Contract | Active | Completed | Complete | Signed/effective date present | TEST-CONTRACT-LIFECYCLE-001 |
| PIM-SM-005 Contract | Pending Signature | Terminated | Terminate (pre-sig) | Allowed | TEST-CONTRACT-LIFECYCLE-003 |
| PIM-SM-005 Contract | Signed/Active | Terminated | Terminate (post-sig) | Signed/effective date present | TEST-CONTRACT-LIFECYCLE-003 |
| PIM-SM-006 Payment | (none) | Unpaid | Create plan | Persisted w/ required fields | TEST-PAYMENT-RECORD-001 |
| PIM-SM-006 Payment | Unpaid/Partially Paid | Partially Paid | Record partial | Amount > 0, cumulative < due | TEST-PAYMENT-RECORD-002 |
| PIM-SM-006 Payment | Unpaid/Partially Paid/Overdue | Paid | Record full | Amount > 0, cumulative == due, no overpayment | TEST-PAYMENT-RECORD-003 |
| PIM-SM-006 Payment | Unpaid/Partially Paid | Overdue | Due date passes unpaid | Business date > due AND unpaid remains (on-read, deterministic per PSM) | TEST-PAYMENT-OVERDUE-001 |
| PIM-SM-006 Payment | any | (rejected) | Zero/negative/overpayment | **Rejected, no mutation** (PIM-INV-022/023) | TEST-PAYMENT-GUARD-001/002/003 |
| PIM-SM-007 Task | (none) | Open | Create | Related record/owner/due/status/title present | TEST-TASK-LIFECYCLE-001 |
| PIM-SM-007 Task | Open/Overdue | Completed | Complete | Authorized | TEST-TASK-LIFECYCLE-002 |
| PIM-SM-007 Task | Open/Overdue | Cancelled | Cancel | Authorized | TEST-TASK-LIFECYCLE-002 |
| PIM-SM-007 Task | Open | Overdue | Due date passes incomplete | Business date > due AND incomplete (on-read) | TEST-TASK-LIFECYCLE-003 |
| PIM-SM-007 Task | Completed/Cancelled | (no reminder) | Due date reached | **No active reminder** (PIM-INV-028) | TEST-TASK-LIFECYCLE-004 |
| PIM-SM-008 Ownership | Unassigned | Assigned | Assign | Manager/admin, target in team scope; **Sales denied; out-of-scope target rejected** | TEST-OWNER-TRANSFER-001 / -003 |
| PIM-SM-008 Ownership | Assigned | Assigned(new) | Transfer | Manager/admin; open tasks/follow-ups cascade to new owner unless reassigned (EDGE-024) | TEST-OWNER-TRANSFER-002 (transfer) / TEST-OWNER-TRANSFER-004 (open-work cascade, PIM-INV-030/033) |
| PIM-SM-009 Opp Close | Contract Negotiation | Won | Close won | **Related contract Signed** (DEC-017); history preserved; **Won without a Signed contract rejected** | TEST-OPP-CLOSE-001 / -002 |
| PIM-SM-009 Opp Close | any open | Lost | Close lost | **Lost reason recorded**; **no-reason rejected** | TEST-OPP-CLOSE-003 / -004 |
| PIM-SM-009 Opp Close | Won/Lost | (rejected) | Reopen/rollback/re-close | **Rejected — terminal** (PIM-INV-009/037) | TEST-OPP-CLOSE-005 |
| PIM-SM-010 Archive | Active | Archived | Archive | Admin/manager; eligible; obligations resolved; **Sales denied; obligation-blocked rejected** | TEST-ARCHIVE-001 / -003 / -004 |
| PIM-SM-010 Archive | Archived | Archived(read) | Retrieve via filter | Authorized, explicit archived filter | TEST-ARCHIVE-002 |
| PIM-SM-010 Archive | any | (rejected) | Hard delete | **Rejected/unavailable any role** (PIM-INV-040) | TEST-INV-NODELETE-001 |
| PIM-SM-011 User | (none) | Active | Create user | Admin actor; exactly one role | TEST-USER-ADMIN-001 |
| PIM-SM-011 User | Active | Active(new role) | Change role | Admin; **not removing last active admin** | TEST-USER-ADMIN-002 |
| PIM-SM-011 User | Active | Disabled | Disable | Admin; **not last active admin** | TEST-USER-ADMIN-003 |
| PIM-SM-011 User | Disabled | Active | Enable | Admin actor | TEST-USER-ADMIN-004 |
| PIM-SM-011 User | Active | (rejected) | Disable/downgrade last admin | **Rejected before save, no mutation** (PIM-INV-046) | TEST-INV-LASTADMIN-001 |

Reject-case count in this section: 18 (TEST-LEAD-QUALIFY-004/005/006/007,
TEST-OPP-STAGE-002/003, TEST-CONTRACT-LIFECYCLE-002,
TEST-CONTRACT-CREATE-003, TEST-PAYMENT-GUARD-001/002/003, TEST-OPP-CLOSE-002/004/005,
TEST-ARCHIVE-003/004, TEST-INV-NODELETE-001, TEST-INV-LASTADMIN-001).
(TEST-QUOTE-ACCEPT-002 re-aimed from a second-accept reject case to a positive
single-quote-accepted case 2026-06-01 per DEC-018; TEST-OPP-CLOSE-002 is now
"Won without a Signed contract rejected" per DEC-017.)

---

## Permission Tests

Three roles (Administrator / Sales Manager / Sales) × key actions × resource scope
(owned / team / governed-all) × expected allow/deny, traceable to PM-* / PERM-*.
All Integration (backend/API authorization), with E2E where the action is
user-visible. Every denied test asserts no data exposure and no mutation
(AUTHZ-008, SEC-004).

| Role | Action | Resource / Scope | Expected | Trace | Test ID |
|---|---|---|---|---|---|
| Administrator | Manage users/roles | governed | Allow | PM-003, PERM-002 | TEST-PERM-USERADMIN-001 |
| Sales Manager | Manage users/roles | any | Deny | PM-004, PERM-003 | TEST-PERM-USERADMIN-002 |
| Sales | Manage users/roles | any | Deny | PM-005 | TEST-PERM-USERADMIN-003 |
| Administrator | Disable/downgrade | last active Administrator | Deny (blocked before save) | PM-007, PIM-INV-046 | TEST-INV-LASTADMIN-001 |
| Administrator | CRUD all entities | governed | Allow | PM-008 | TEST-PERM-CRUD-ADMIN-001 |
| Sales Manager | CRUD all entities | team | Allow | PM-009, PERM-006 | TEST-PERM-CRUD-MGR-001 |
| Sales | CRUD lead/customer/opp/quote/contract/payment | owned/assigned | Allow | PM-010/011/016/018/020/021/022, PERM-004/008/010/014/015/016 | TEST-PERM-CRUD-SALES-001 |
| Sales | View/edit | non-owned record (IDOR via id) | Deny, no leakage/mutation | PM-012/017/019, PERM-007/009, ABUSE-002 | TEST-AUTHZ-SCOPE-004 / TEST-ABUSE-IDOR-001 |
| Sales | Qualify/edit/convert | Unassigned lead | Deny | PM-013, PERM-005 | TEST-LEAD-QUALIFY-004 |
| Sales | Assign/transfer owner | any | Deny | PM-015, PERM-... | TEST-LEAD-ASSIGN-002 / TEST-OWNER-TRANSFER-003 |
| Sales Manager / Administrator | Assign/transfer owner | team/governed | Allow, history recorded | PM-014, PERM-006 | TEST-OWNER-TRANSFER-001/002 |
| Sales | Archive | any record | Deny | PM-028, PERM-012 | TEST-ARCHIVE-004 |
| Sales Manager | Archive | eligible team record, obligations resolved | Allow | PM-027, PERM-013 | TEST-ARCHIVE-001 |
| Sales Manager | Archive | record w/ unresolved obligations | Deny/block | PM-027, EDGE-032 | TEST-ARCHIVE-003 |
| Any role | Hard delete | core record | Deny/unavailable | PM-029, ABUSE-017 | TEST-INV-NODELETE-001 |
| Sales | View | record-local history (owned) | Allow | PM-024, PERM-017 | TEST-HISTORY-001 |
| Sales | View | record-local history (non-owned) | Deny | PM-025, PERM-018, ABUSE-013 | TEST-HISTORY-003 |
| Administrator | View/query | global operation logs | Allow | PM-040, PERM-021 | TEST-OPLOG-001 |
| Sales Manager | View/query | global operation logs | Deny | PM-041, PERM-022, ABUSE-012 | TEST-OPLOG-003 |
| Sales | View/query | global operation logs | Deny | PM-042, PERM-023 | TEST-OPLOG-004 |
| Administrator / Sales Manager | Import/export | governed/team | Allow | PM-034/035/037/038, PERM-024/025 | TEST-CSV-IMPORT-001 / TEST-CSV-EXPORT-001 |
| Sales | Import/export | any | Deny | PM-036/039, PERM-026, ABUSE-010 | TEST-CSV-IMPORT-004 / TEST-CSV-EXPORT-002 |
| Sales Manager | View | team overview / team reports | Allow (team) | PM-044, PERM-019/028 | TEST-TEAM-OVERVIEW-001 / TEST-BASIC-REPORT-001 |
| Sales | View | team overview / manager-admin reports | Deny | PM-045, PERM-020/029, ABUSE-015 | TEST-TEAM-OVERVIEW-002 / TEST-BASIC-REPORT-003 |
| Administrator/Manager/Sales | View | reminders (authorized related only) | Allow authorized / Deny unauthorized | PM-046/047, ABUSE-021 | TEST-REMINDER-005 |
| Administrator/Manager/Sales | View | archived (scope + explicit filter) | Allow in scope / Deny out of scope | PM-030/031/032/033, PERM-030/031 | TEST-NAV-RETRIEVE-006 / TEST-ABUSE-ARCHIVED-001 |

---

## Invariant Tests

Each load-bearing PIM-INV-* gets a test that proves the invariant holds (positive)
and that the violation is rejected (negative, no mutation). Type: Unit for
single-aggregate guards, Integration where persistence/cross-aggregate.

| Invariant | Statement | Positive / Negative | Test ID |
|---|---|---|---|
| PIM-INV-012 | Each opportunity has exactly one quote (DEC-018) | Single quote per opportunity holds; a second quote on the same opportunity is not created | TEST-INV-ONEACCEPT-001 |
| PIM-INV-007/025/035 | Won requires Signed contract (DEC-017; full payment not a precondition) | Won succeeds when related contract is Signed; Won without a Signed contract rejected | TEST-INV-WONAFTERPAY-001 |
| PIM-INV-009/037 | Won/Lost terminal, non-reopenable | Post-close reopen/rollback/re-close rejected; notes/tasks allowed | TEST-INV-TERMINAL-001 |
| PIM-INV-040 | No hard delete of any core CRM record | Delete endpoint/route unavailable or rejected; record stays persisted | TEST-INV-NODELETE-001 |
| PIM-INV-041 | Cannot archive record with unresolved active downstream obligations | Archive blocked when open tasks / pending-signature contracts / unpaid payments exist; allowed once resolved | TEST-INV-ARCHIVEBLOCK-001 |
| PIM-INV-022 | Actual payment amount must be > 0 | Zero and negative rejected; positive accepted | TEST-INV-PAYAMOUNT-001 |
| PIM-INV-023 | Cumulative payment cannot exceed contract remaining (overpayment rejected, contract-level ceiling) | Overpayment at contract level rejected even across plans | TEST-INV-OVERPAY-001 |
| PIM-INV-024 | Single currency across all amounts | Non-matching currency rejected / single-currency enforced | TEST-INV-CURRENCY-001 |
| PIM-INV-046 | Cannot remove/deactivate last active Administrator | Disable/downgrade of last admin rejected before save, no mutation | TEST-INV-LASTADMIN-001 |
| PIM-INV-013/014/019 | Contract references exactly one Accepted, non-Expired quote | Accepted quote link allowed; expired/non-Accepted rejected | TEST-INV-CONTRACTQUOTE-001 |
| PIM-INV-016/017 | Pending Signature needs expected signed date (not signed date); signed states need signed/effective date | Each guard proven both directions | TEST-INV-CONTRACTDATE-001 |
| PIM-INV-028 | Completed/Cancelled tasks not active reminders | Inactive tasks suppressed from reminders | TEST-INV-TASKREMINDER-001 |
| PIM-INV-050/051 | Retention ≥ committed minimum, anchored to lifecycle events, never shortened, never hard-deleted; logs append-only | Records/logs retained past lifecycle event; no path shortens below committed minimum | TEST-RETENTION-001 |
| PIM-INV-049 | Every record/log carries committed Data Classification governing visibility/masking | Sensitive before/after values masked per classification in denied/error/log states | TEST-RETENTION-002 |

---

## Edge-Case Tests

EDGE-* mapped to TEST-* IDs. Overdue/business-date evaluation is deterministic
(PSM resolved on-read against the supplied `businessDate` in `Asia/Shanghai`), so
EDGE-020/021/023/037 are designable as deterministic tests.

| Edge | Scenario | Test ID |
|---|---|---|
| EDGE-001 | Unauthenticated route access denied | TEST-AUTH-LOGIN-004 |
| EDGE-002 | Sales accesses non-owned record → hidden/denied | TEST-AUTHZ-SCOPE-004 |
| EDGE-003 | Lead save missing required fields blocked | TEST-LEAD-CREATE-002 |
| EDGE-004 | Sales qualifies Unassigned lead → denied | TEST-LEAD-QUALIFY-004 |
| EDGE-005 | Invalid lead converted → rejected until restored | TEST-LEAD-QUALIFY-005 |
| EDGE-006 | Converted lead re-converted → rejected | TEST-LEAD-QUALIFY-007 |
| EDGE-007 | Contact without company → blocked | TEST-CONTACT-LINK-002 |
| EDGE-008 | Forbidden opportunity transition → rejected | TEST-OPP-STAGE-002 |
| EDGE-009 | Won without a Signed contract → rejected (DEC-017) | TEST-OPP-CLOSE-002 |
| EDGE-010 | Lost without reason → rejected | TEST-OPP-CLOSE-004 |
| EDGE-011 | Won/Lost reopened → rejected | TEST-OPP-CLOSE-005 |
| EDGE-012 | Exactly one quote per opportunity → the single quote is accepted; no second quote exists (DEC-018; second-accept ambiguity removed) | TEST-QUOTE-ACCEPT-002 |
| EDGE-013 | Expired quote linked to contract → rejected | TEST-CONTRACT-CREATE-003 |
| EDGE-014 | Pending Signature lacks expected signed date → rejected | TEST-CONTRACT-CREATE-002 |
| EDGE-015 | Pending Signature lacks signed date → allowed | TEST-CONTRACT-CREATE-001 |
| EDGE-016 | Signed/Active/Completed lacks signed date → rejected | TEST-CONTRACT-LIFECYCLE-002 |
| EDGE-017 | Contract amount differs from quote → reason required | TEST-CONTRACT-AMOUNT-DIFF-001 |
| EDGE-018 | Zero/negative payment → rejected | TEST-PAYMENT-GUARD-001/002 |
| EDGE-019 | Overpayment beyond remaining → rejected | TEST-PAYMENT-GUARD-003 |
| EDGE-020 | Payment due date passes unpaid → Overdue + reminder | TEST-PAYMENT-OVERDUE-001 / TEST-REMINDER-003 |
| EDGE-021 | Pending Signature past expected signed date → reminder | TEST-REMINDER-002 |
| EDGE-022 | Signed/terminated/fully-paid contract in reminders → none | TEST-REMINDER-004 |
| EDGE-023 | Completed/cancelled task at due date → no reminder | TEST-TASK-LIFECYCLE-004 |
| EDGE-024 | Parent owner change with open tasks → tasks/follow-ups transfer unless reassigned | TEST-OWNER-TRANSFER-004 |
| EDGE-025 | Duplicate data entered → warning, save may continue | TEST-DUPLICATE-WARN-005 |
| EDGE-026 | CSV mixed valid/invalid rows → valid import, invalid reported, no corruption | TEST-CSV-IMPORT-002 |
| EDGE-027 | Sales CSV export → denied | TEST-CSV-EXPORT-002 |
| EDGE-028 | Report with no data → empty/zero state | TEST-BASIC-REPORT-002 |
| EDGE-029 | Edit history/operation logs → rejected/unavailable | TEST-HISTORY-004 / TEST-OPLOG-005 |
| EDGE-030 | Refresh/re-login/restart after save → data persists | TEST-PERSISTENCE-001/002/003 |
| EDGE-031 | Archived would appear in active list/reminder/report → hidden | TEST-NAV-RETRIEVE-006 / TEST-REMINDER-004 / TEST-BASIC-REPORT-005 |
| EDGE-032 | Archive with unresolved obligations → blocked | TEST-ARCHIVE-003 |
| EDGE-033 | Company name differs by case/spaces → warning | TEST-DUPLICATE-WARN-001 |
| EDGE-034 | Contact phone differs by spaces/hyphens/parens → warning | TEST-DUPLICATE-WARN-002 |
| EDGE-035 | Contact email differs by case/spaces → warning | TEST-DUPLICATE-WARN-003 |
| EDGE-036 | Edit opportunity stage after Won/Lost → rejected; notes/tasks allowed | TEST-OPP-CLOSE-006 |
| EDGE-037 | Reminder evaluated around local date boundary → workspace business date used | TEST-REMINDER-BOUNDARY-001 |

---

## Abuse / Negative-Security Tests

ABUSE-* mapped to TEST-* IDs. Per `abuse-cases.md` and SEC-017, every P0 abuse
case requires a backend/API negative-path test (direct-API, not UI-only).
References SEC-* and ARCH-ACC-002 (direct-API denial).

| Abuse | Scenario | Required protection | Test ID |
|---|---|---|---|
| ABUSE-001 (P0) | Unauthenticated CRM data access (direct route/API) | Reject, no data, failure logged | TEST-ABUSE-UNAUTH-001 (SEC-001) |
| ABUSE-002 (P0) | IDOR — Sales uses non-owned id across each P0 resource group | Backend ownership/scope check; deny w/o existence or value leakage | TEST-ABUSE-IDOR-001 (SEC-003/004) |
| ABUSE-003 (P0) | Unauthorized mutation (edit/transfer/archive/close non-owned) | Deny; no mutation; no false success | TEST-ABUSE-MUTATE-001 (AUTHZ-008) |
| ABUSE-004 (P0) | Privilege escalation via user/role API or payload role fields | Deny; role fields ignored/rejected | TEST-ABUSE-PRIVESC-001 (SEC-008/AUTHZ-005) |
| ABUSE-005 (P0) | Last Administrator removal | Blocked before save; no mutation; logged | TEST-INV-LASTADMIN-001 (SEC-008) |
| ABUSE-006 (P0) | Login enumeration | Unified failure message; no account-state disclosure | TEST-ABUSE-ENUM-001 (AUTH-002) |
| ABUSE-013 (P0) | Record-local history leakage (non-owned record/child event) | Backend related-record permission check on history query | TEST-HISTORY-003 (SEC-009) |
| ABUSE-016 (P0/P1) | Archived record bypass via active list/search/reminder/report/direct id | Excluded from active/default; explicit filter + authz required | TEST-ABUSE-ARCHIVED-001 (SEC-015) |
| ABUSE-017 (P0) | Hard-delete attempt via endpoint/route | No hard delete; reject; data persists | TEST-INV-NODELETE-001 (SEC-006) |
| ABUSE-018 (P0) | Business-rule bypass via direct API (forbidden stage, Won w/o Signed contract, expired-quote contract link, overpayment, post-close edit) | Backend validates rule + authz together; no mutation on failure | TEST-ABUSE-BRBYPASS-001 (covers TEST-OPP-STAGE-002, TEST-OPP-CLOSE-002, TEST-CONTRACT-CREATE-003, TEST-PAYMENT-GUARD-003, TEST-OPP-CLOSE-006) |
| ABUSE-019 (P0/P1) | Acting on behalf of Sales (payload claims other actor) | Authenticated actor recorded as actor; owner/assignee separate | TEST-ABUSE-ACTAS-001 (AUTHZ-009) |
| ABUSE-022 (P0) | Stale session / stale role after disable or role change | Protected requests re-evaluate active account/role; disabled denied | TEST-AUTH-LOGIN-006 (AUTH-004) |
| ABUSE-007 (P1) | CSV formula injection (import/export) | Validation flags/escapes dangerous cells; export safely encodes | TEST-ABUSE-CSVINJECT-001 (SEC-011, ARCH-ACC-012) |
| ABUSE-008 (P1) | Import row authorization bypass (out-of-scope refs) | Per-row permission/reference check before mutation | TEST-ABUSE-IMPORTAUTHZ-001 (SEC-011) |
| ABUSE-009 (P1) | Import partial-failure corruption | Valid rows handled, invalid reported, existing not corrupted, no silent overwrite | TEST-CSV-IMPORT-002 |
| ABUSE-010 (P1) | Export data leakage (Sales export; manager filter expands scope/archived) | Sales denied; manager limited to team; archived needs explicit authorized filter; logged | TEST-CSV-EXPORT-002 / TEST-ABUSE-EXPORTLEAK-001 |
| ABUSE-011 (P1) | Export confirmation bypass via direct API | Backend requires authorized export semantics + logs | TEST-ABUSE-EXPORTCONFIRM-001 (SEC-012) |
| ABUSE-012 (P1) | Global operation-log leakage (route/API/id guessing) | Deny; no event/target/diff/user disclosure | TEST-OPLOG-003 / TEST-OPLOG-004 |
| ABUSE-014 (P1) | Sensitive before/after value leakage in errors/toasts/row errors/log summaries | Safe summaries by default; details require authz + classification | TEST-RETENTION-002 (SEC-014) |
| ABUSE-015 (P1) | Report inference (Sales access / unauthorized records pre-filter) | Authz before aggregation; Sales denied | TEST-BASIC-REPORT-003 / TEST-BASIC-REPORT-004 |
| ABUSE-020 (P1) | Duplicate-warning enumeration | Warnings reveal no restricted matched-record details; no merge | TEST-ABUSE-DUPENUM-001 (SEC-014) |
| ABUSE-021 (P0/P1) | Reminder leakage of unauthorized due/overdue items | Reminder queries apply record permissions; hide unauthorized | TEST-REMINDER-005 |

Cross-service / service-boundary negative tests (PSM CONTRACT-019, ARCH-ACC-009):
TEST-ABUSE-S2S-001 — service-to-service call without a valid signed token returns
`SERVICE_AUTH_FAILED`; network trust alone does not authorize (SVC-ACC-008).

---

## Traceability note

- Every acceptance item ACC-001..023 has ≥1 `TEST-*` in `## Acceptance To Test
  Mapping` (confirmed: 23/23). P0 items carry both positive and negative cases.
- These `TEST-*` IDs are what the PSM `traceability-matrix.md` / PSM Traceability
  `Test ID = pending (G7)` column will resolve to: each ACC row's `pending (G7)`
  Test placeholder maps to this artifact's `TEST-*` family for that ACC.
- TEST families reuse the PIM behavior concept tags (PIM-BEH-001..034 TEST-*
  tags) verbatim as family names and number concrete cases within them, plus
  invariant (`TEST-INV-*`), permission (`TEST-PERM-*`), abuse (`TEST-ABUSE-*`),
  and deployment (`TEST-DEPLOY-*`) families introduced here for cross-cutting
  coverage.
- Test types are assigned per the PSM service/contract/flow design: cross-service
  flows FLOW-001..013 drive Integration coverage; single-aggregate guards are
  Unit-eligible; user-visible paths are E2E.

---

## Open / Blocked

All three prior G6 blockers were resolved 2026-06-01 by the Formal Scope Change by
User (decision-log.md DEC-017..020). The previously `Pending` tests are retired in
place; the formerly-undecided observables are now decided, so no `Pending` test
remains in this artifact.

- **BLK-001 — Opportunity Status enumerated value set (ACC-007, P0) — RESOLVED
  2026-06-01 (DEC-020).** The separate Opportunity `Status` field is removed;
  Pipeline Stage (PIM-SM-002), including terminal Won/Lost, is the sole lifecycle
  dimension. TEST-OPP-STATUS-ENUM-001 is **[RETIRED 2026-06-01 — resolved by
  DEC-020]**; TEST-OPP-CREATE-003 re-aimed to plain opportunity create (no Status
  field). No distinct Status enumeration exists to test. Refs: DEC-020,
  PIM-OPEN-003 (resolved), PIM-SM-003 (retired).
- **BLK-002 — Multi-plan "fully paid" aggregation (ACC-011, ACC-013, P0) —
  RESOLVED 2026-06-01 (DEC-017/019).** Won no longer depends on full payment (Won =
  related contract Signed, DEC-017; payment tracking retained but decoupled,
  DEC-019), so no contract "fully paid" aggregation gates closure.
  TEST-PAYMENT-FULLPAID-AGG-001 is **[RETIRED 2026-06-01 — resolved by
  DEC-017/019]**. Payment coverage retained and unchanged: single-plan full payment
  (TEST-PAYMENT-RECORD-003), contract-level overpayment ceiling
  (TEST-PAYMENT-GUARD-003 / TEST-INV-OVERPAY-001), partial/overdue
  (TEST-PAYMENT-RECORD-002 / TEST-PAYMENT-OVERDUE-001), single currency
  (TEST-PAYMENT-GUARD-004 / TEST-INV-CURRENCY-001) — none of these gate Won. Refs:
  DEC-017, DEC-019, PIM-OPEN-005 (resolved).
- **BLK-003 — Second-quote-accept observable outcome (ACC-009, P0) — RESOLVED
  2026-06-01 (DEC-018).** Each opportunity has exactly one quote, so there is no
  second quote to accept and no losing-path ambiguity. TEST-QUOTE-ACCEPT-003 is
  **[RETIRED 2026-06-01 — resolved by DEC-018]**; TEST-QUOTE-ACCEPT-002 re-aimed
  to "the single quote is accepted" and TEST-INV-ONEACCEPT-001 to "each
  opportunity has exactly one quote". Refs: DEC-018, EDGE-012, PIM-OPEN-001
  (resolved).

Determinism note (NOT blocked): overdue/business-date evaluation is deterministic.
PSM resolved PIM-OPEN-002 / BLK-A01 as on-read evaluation against the supplied
`businessDate` (`Asia/Shanghai`) at query time (PSM Resolved Mechanisms;
api-spec.md "Reminder Query"; FLOW-005). Therefore TEST-PAYMENT-OVERDUE-001,
TEST-TASK-LIFECYCLE-003, TEST-REMINDER-001/002/003, and TEST-REMINDER-BOUNDARY-001
CAN be designed and executed deterministically by supplying the business date.

ACC-017 note (NOT blocked, evidence-staged): TEST-DEPLOY-SMOKE-001/002 are
designable as deployment smoke/Manual now; the HTTPS/TLS, off-server
backup+restore rehearsal, security-group, and monitoring items are
release-evidence (ARCH-ACC-004/008/013/014, carried release blockers) verified at
G11/G12, not G6 design gaps.
