# PIM

Platform Independent Model — the domain objects, behaviors, and per-aggregate
state machines (states, transitions, guards, invariants) that refine the accepted
CIM concepts and processes into platform-neutral domain structure. No platform,
service, contract, schema, or persistence detail.

## Document Control

- Project: CRM System
- Date: 2026-05-31
- Role: Domain Modeling
- Gate: G6 (MDA Modeling)
- Scope note: Models accepted design; not design authority for scope. This PIM
  traces the accepted CIM (its direct input), the accepted G4 business design
  (business-rules.md, edge-cases.md, business-processes.md, business-glossary.md,
  decision-log.md), and the accepted G5 architecture set. It does not introduce
  new scope, rules, fields, or concepts. It refines CIM concepts/processes into
  domain objects, behaviors, and state machines (the transition matrices CIM
  explicitly deferred to PIM). Where domain modeling needs something the accepted
  sources do not provide, it is recorded under `## Open / Blocked` rather than
  invented. PSM-tier detail (services, API/event/error/permission contracts,
  event IDs/schema/routing, persistence/query design, reliability mechanisms,
  deployment) is named at PIM altitude and deferred to PSM with a one-line note.
  Amended 2026-06-01 per Formal Scope Change by User (decision-log.md
  DEC-017..020): Won = related contract Signed (not full payment); exactly one
  quote per opportunity; payment tracking retained but decoupled from Won;
  Opportunity Status dimension removed (Pipeline Stage is the sole lifecycle
  dimension); `Payment In Progress` and `Contract Signed` pipeline stages removed.
  Affected IDs retired in place (not renumbered) to preserve cross-references.

## Tier Altitude Statement

This document stays at platform-independent altitude:

- It models domain objects (aggregates/entities/value objects), their
  responsibilities, behaviors, and state machines with transitions, guards, and
  invariants.
- It does NOT name services, bounded-context engineering splits, API/event/error/
  permission contracts, concrete event IDs/schema/routing, database schemas,
  field types/DDL, persistence/query design, reliability mechanisms (idempotency/
  retry/timeout), or deployment/host/provider detail. Those are PSM and are
  deferred with explicit one-line notes.
- It does NOT re-derive business vocabulary or re-state processes; it references
  CIM IDs (CIM-001…CIM-049, CIM-PROC-001…CIM-PROC-024) as authority.

## Domain Objects

Aggregates are the consistency boundaries that enforce invariants. Entities and
value objects live inside an aggregate. Concrete attribute types, schema, and
persistence ownership are modeled in PSM.

| ID | Type | Name | Responsibility | Acceptance ID |
|---|---|---|---|---|
| PIM-001 | Aggregate | User | A CRM account holding exactly one Role; governs authentication identity and authorization scope (CIM-001, CIM-002). Holds active/disabled status. Enforces the last-active-Administrator invariant. | ACC-001, ACC-002, ACC-022 |
| PIM-002 | Value Object | Role | The access classification (Administrator, Sales Manager, Sales) bound to a User; exactly three roles committed (CIM-002..CIM-005, DEC-005). | ACC-001, ACC-002 |
| PIM-003 | Value Object | Ownership / Assignment | The relation binding a CRM record to a responsible/assigned User; drives Sales visibility and Sales Manager team scope (CIM-007). Transferable within team scope. | ACC-002, ACC-003 |
| PIM-004 | Aggregate | Lead | A pre-qualification ToB sales target carrying lead/company name, source, status, and owner; may be Unassigned before assignment; convertible to Opportunity once Valid (CIM-008, CIM-009, CIM-010). | ACC-003, ACC-004 |
| PIM-005 | Aggregate | Company / Customer | A ToB account record carrying company name, customer status, and owner; parent of Contacts and downstream commercial records (CIM-011, CIM-012). | ACC-005 |
| PIM-006 | Entity | Contact | A person/role under a Company/Customer with name, related company/customer, and at least one contact method or role note (CIM-013). Belongs to the Company/Customer aggregate. | ACC-006 |
| PIM-007 | Aggregate | Opportunity | A sales deal linked to customer, contacts, owner, expected amount, expected close date; carries Pipeline Stage (the sole lifecycle dimension) and the terminal Won/Lost outcome (CIM-014, CIM-015, CIM-017, CIM-018). Status dimension removed 2026-06-01 (DEC-020). | ACC-007, ACC-008, ACC-013 |
| PIM-008 | Aggregate | Quote | A commercial offer linked to opportunity and customer with amount, status, validity end date, and owner; exactly one quote per opportunity (DEC-018), so at most one Accepted follows trivially (CIM-019, CIM-020). | ACC-009 |
| PIM-009 | Aggregate | Contract | A record-based contract linked to customer, opportunity, and the Accepted quote; carries amount, status, required contract note, expected signed date, and signed/effective date for signed states; no approval/e-signature/template (CIM-021..CIM-025). | ACC-010 |
| PIM-010 | Entity | Payment Plan | A planned payment under a Contract with due amount, due date, and status; basis for due/overdue tracking (CIM-026, CIM-028). Belongs to the Contract aggregate. | ACC-011 |
| PIM-011 | Entity | Actual Payment | A recorded payment event under a Contract with paid amount, payment date, and status; zero/negative/overpayment rejected; single currency (CIM-027, CIM-029). Belongs to the Contract aggregate. | ACC-011, ACC-013 |
| PIM-012 | Value Object | Amount / Money | Single-currency monetary value used across Opportunity, Quote, Contract, Payment; excludes tax, discount, multi-currency (CIM-029, DEC-013). | ACC-007, ACC-009, ACC-010, ACC-011 |
| PIM-013 | Entity | Activity | A business interaction record linked to a CRM record, preserving follow-up history (CIM-030). | ACC-012 |
| PIM-014 | Entity | Note | A textual business note linked to a CRM record, visible by record permission (CIM-031). | ACC-012 |
| PIM-015 | Aggregate | Task | A follow-up work item with owner, due date, status, and title; due/overdue creates reminders; transfers with parent owner unless reassigned (CIM-032). | ACC-012, ACC-021 |
| PIM-016 | Value Object | Reminder / Follow-up | An in-app follow-up signal derived from due/overdue Tasks, Pending Signature Contracts past expected signed date, and due/overdue Payments; in-app only; not independently persisted state — derived from source aggregates by their guards (CIM-033). The evaluation-trigger mechanism is PSM/Architecture (see Open / Blocked). | ACC-021 |
| PIM-017 | Value Object | Business Date | Workspace-local date used to evaluate due/overdue (CIM-034, BR-021). The on-read-vs-scheduled evaluation mechanism is deferred (see Open / Blocked). | ACC-021 |
| PIM-018 | Value Object | Record-Local History Event | A permitted, non-editable timeline event (owner, stage, status, quote, contract, payment, task, archive change) emitted as a side effect of a domain behavior and read by record permission (CIM-035). Concrete event fields/identifiers/schema are PSM. | ACC-014 |
| PIM-019 | Value Object | Operation-Log Event | An Administrator-only operational audit event for access-sensitive and record-mutation actions, non-editable through normal CRM actions (CIM-036). Concrete log fields, event IDs, and schema are PSM. | ACC-022 |
| PIM-020 | Value Object | Archive State | A non-delete lifecycle marker that removes a record from active work views, active reminders, and default reports while preserving it for archived filters and history; no hard delete (CIM-037, CIM-038, CIM-039). Applies to eligible records (BR-016). | ACC-002, ACC-014, ACC-015, ACC-023 |
| PIM-021 | Value Object | Duplicate Warning | A non-blocking warning raised on lead/company/contact data matching configured duplicate rules; no silent overwrite or merge (CIM-040). Normalization rules per BR-019. | ACC-019 |
| PIM-022 | Aggregate | CSV Import Job | An authorized bulk-entry operation that validates rows and reports row-level errors; valid existing records not corrupted by failed rows; Sales cannot import (CIM-041). | ACC-020 |
| PIM-023 | Aggregate | CSV Export Job | An authorized export operation including only authorized records, excluding archived by default; Sales cannot export (CIM-042). | ACC-020 |
| PIM-024 | Value Object | Team Overview | A Sales Manager read view aggregating team records and pipeline status (CIM-043). Query design is PSM. | ACC-018 |
| PIM-025 | Value Object | Basic Report | Counts and sums for committed groupings over persisted authorized records; default excludes archived (CIM-044, BR-014, BR-017). Query/computation design is PSM. | ACC-018, ACC-023 |
| PIM-026 | Value Object | Core Retrieval View | A role-scoped list/detail/search/basic-filter view across the P0 entities with empty-state, invalid-filter feedback, and permission hiding (CIM-047). Query/API design is PSM. | ACC-015 |
| PIM-027 | Value Object | Data Classification | The sensitivity class on a record/log governing visibility/masking and retention tier (CIM-048, PRIV-001..016). | ACC-014, ACC-022 |
| PIM-028 | Value Object | Retention Policy | The committed minimum-retention expectation per data category anchored to lifecycle events; never shortened below committed; concrete durations/storage are PSM (CIM-049, PRIV-*, COMP-013). | ACC-014, ACC-022 |

Aggregate boundary note (modeled in PSM): the engineering split into services/
bounded contexts and data ownership is PSM; PIM only fixes the consistency
boundaries above (e.g., Contract owns its Payment Plans and Actual Payments so the
overpayment and full-payment invariants are enforced inside one boundary).

## State Machines

Each subsection gives a transition table (From State | Event/Action | To State |
Guard) and an Invariants list. State NAMES are taken from the CIM as accepted
vocabulary; transition and guard rules are taken from business-rules.md,
edge-cases.md, and the decision log. Forbidden transitions are rejected without
data mutation (BR-005, EDGE-008). Concrete denial-error/event mechanisms are PSM.

### PIM-SM-001 — Lead Status (CIM-009)

States: Unassigned, Pending Qualification, Valid, Invalid, Converted To
Opportunity.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create lead (unassigned) | Unassigned | Lead name or company name, source, status present (BR-003, EDGE-003) |
| (none) | Create lead with owner | Pending Qualification | Required fields present AND owner assigned (BR-003, BR-004) |
| Unassigned | Assign owner | Pending Qualification | Owner assigned before Pending Qualification or later (BR-004); actor authorized (BR-001) |
| Pending Qualification | Qualify valid | Valid | Actor authorized; lead not Unassigned (BR-004, EDGE-004) |
| Pending Qualification | Mark invalid | Invalid | Invalid reason recorded (BR-004) |
| Invalid | Restore | Pending Qualification | Actor is Administrator or Sales Manager (BR-004, EDGE-005) |
| Valid | Convert to opportunity | Converted To Opportunity | Downstream Opportunity link created; lead not already converted (BR-004, EDGE-006); preserves lead history (CIM-PROC-004) |

Invariants:
- PIM-INV-001: An Unassigned lead cannot be qualified, edited by Sales, or
  converted (BR-004, EDGE-004).
- PIM-INV-002: An Invalid lead cannot convert unless first restored to Pending
  Qualification by Administrator or Sales Manager (BR-004, EDGE-005).
- PIM-INV-003: A Converted lead cannot be converted again (BR-004, EDGE-006).
- PIM-INV-004: Owner is required before Pending Qualification or any later state
  (BR-003, BR-004).
- PIM-INV-005: Owner assignment, qualification, disqualification, and conversion
  each emit a record-local history event (PIM-018, BR-004, CIM-035).

### PIM-SM-002 — Opportunity Pipeline Stage (CIM-015)

States: New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost.
Won and Lost are terminal. (`Contract Signed` and `Payment In Progress` stages
removed 2026-06-01 per DEC-017.)

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create opportunity | New Opportunity | Related customer, owner, stage, expected amount, expected close date present (no separate Status — DEC-020) (BR-003) |
| New Opportunity | Advance | Needs Confirmed | Allowed forward transition (BR-005) |
| Needs Confirmed | Advance | Quote | Allowed forward transition (BR-005) |
| Quote | Advance | Contract Negotiation | Allowed forward transition (BR-005) |
| Contract Negotiation | Close won | Won | Related contract is Signed (DEC-017, BR-005, EDGE-009); see PIM-SM-009 |
| New Opportunity / Needs Confirmed / Quote / Contract Negotiation | Close lost | Lost | Lost reason recorded (BR-005, EDGE-010); see PIM-SM-009 |

Invariants:
- PIM-INV-006: Only transitions in the accepted PRD transition table are allowed;
  forbidden transitions (including arbitrary stage rollback) are rejected without
  data mutation (BR-005, BR-020, EDGE-008).
- PIM-INV-007: Won is reached when the related contract is Signed (DEC-017,
  BR-005, EDGE-009). Full payment is NOT a Won precondition (payment is decoupled
  post-sale follow-up — DEC-019).
- PIM-INV-008: Lost requires a recorded lost reason (BR-005, EDGE-010).
- PIM-INV-009: Won and Lost are terminal and non-reopenable in the committed
  scope; reopen, stage rollback after close, early Won, and Lost-without-reason
  are rejected (BR-005, BR-020, EDGE-011, EDGE-036, DEC-012, DEC-017).
- PIM-INV-010: Each stage change emits a record-local history event (PIM-018,
  CIM-PROC-007, ACC-008).

Note: Pipeline Stage is the sole opportunity lifecycle dimension; the former
separate Opportunity Status dimension was removed 2026-06-01 (DEC-020, see retired
PIM-SM-003).

### PIM-SM-003 — Opportunity Status (CIM-016) — [RETIRED 2026-06-01 DEC-020]

**[RETIRED 2026-06-01 DEC-020 — the separate Opportunity `Status` dimension is
removed; Pipeline Stage (PIM-SM-002), including terminal Won/Lost, is the sole
opportunity lifecycle dimension.]** ID retained for cross-reference stability; not
renumbered. The open→closed-by-outcome behavior this section formerly modeled is
fully carried by Pipeline Stage (PIM-SM-002) and Close (PIM-SM-009); there is no
separate Status state machine.

Invariants:
- PIM-INV-011: **[RETIRED 2026-06-01 DEC-020 — folded into stage closure.]** An
  opportunity closed Won/Lost (terminal Pipeline Stage) is not editable through
  normal Sales workflow except added notes/tasks; this is now wholly enforced by
  PIM-SM-002/PIM-SM-009 and PIM-INV-037 (BR-020, EDGE-036).

### PIM-SM-004 — Quote Status (CIM-020)

States: Draft, Sent, Accepted, Rejected, Expired.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create quote | Draft | Related opportunity, customer, amount, status, validity end date, owner present (BR-003) |
| Draft | Send | Sent | Actor authorized (BR-001) |
| Sent | Accept | Accepted | Quote is the opportunity's single quote (DEC-018, BR-006) |
| Sent | Reject | Rejected | — |
| Draft / Sent | Expire | Expired | Validity end date passed (BR-006) |

Invariants:
- PIM-INV-012: Each opportunity has exactly one quote (DEC-018, BR-006, CIM-019),
  so at most one Accepted quote per opportunity holds trivially. This invariant is
  owned by the Opportunity consistency boundary.
- PIM-INV-013: An Expired quote cannot be linked to a new Contract (BR-006,
  EDGE-013).
- PIM-INV-014: A Contract may reference only an Accepted quote (BR-006, PIM-SM-005).
- PIM-INV-015: Quote acceptance emits a record-local history event and is an
  operation-log event class (PIM-018, PIM-019, BR-010).

Note: with exactly one quote per opportunity (DEC-018) there is no second quote to
accept, so the former reject-vs-auto-demote ambiguity (previously parked as
PIM-OPEN-001) no longer exists.

### PIM-SM-005 — Contract Status (CIM-022)

States: Pending Signature, Signed, Active, Completed, Terminated. State names are
taken from CIM-022; the linear Pending Signature → Signed → Active → Completed
ordering is refined from that named lifecycle (BR-007 governs the per-state date
guards, not the existence of the transitions).

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create contract | Pending Signature | Related customer, opportunity, Accepted quote, amount, status, contract note, expected signed date present; quote is Accepted and not Expired (BR-003, BR-006, BR-007, EDGE-013, EDGE-014) |
| Pending Signature | Sign | Signed | Signed/effective date present (BR-007, EDGE-016) |
| Signed | Activate | Active | Signed/effective date present (BR-007) |
| Active | Complete | Completed | Signed/effective date present (BR-007) |
| Pending Signature | Terminate | Terminated | — (pre-signature termination) |
| Signed / Active | Terminate | Terminated | Signed/effective date present (post-signature Terminated, BR-007) |

Invariants:
- PIM-INV-016: Pending Signature requires expected signed date and contract note;
  it does not require signed/effective date (BR-007, EDGE-014, EDGE-015,
  DEC-016).
- PIM-INV-017: Signed, Active, Completed, and post-signature Terminated states
  require a signed/effective date (BR-007, EDGE-016).
- PIM-INV-018: Contract amount may differ from the Accepted quote amount only with
  a recorded difference reason (BR-007, DEC-014, EDGE-017).
- PIM-INV-019: A Contract references exactly one Accepted, non-Expired quote
  (BR-006, PIM-INV-013, PIM-INV-014).
- PIM-INV-020: No approval workflow, electronic signature, or template generation
  exists in committed scope (BR-007, DEC-007).
- PIM-INV-021: Contract signature and termination are record-local history events
  and operation-log event classes (PIM-018, PIM-019, BR-010).

### PIM-SM-006 — Payment Status (CIM-028)

States: Unpaid, Partially Paid, Paid, Overdue. Payment status is evaluated for the
Payment Plan against recorded Actual Payments within the Contract aggregate.
Plan-level Paid/Partially Paid/Overdue is a derived position inside the Contract
aggregate; the overpayment ceiling is the contract's remaining amount (EDGE-019).
Payment tracking is retained as post-sale collection follow-up but is decoupled
from Opportunity Won (DEC-019): full payment is no longer a Won precondition (Won
= contract Signed, DEC-017), so the former multi-plan "fully paid" aggregation
question (previously PIM-OPEN-005) no longer gates closure.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create payment plan | Unpaid | Related contract, due amount, due date, status present (BR-003) |
| Unpaid | Record actual payment (partial) | Partially Paid | Paid amount > 0 AND cumulative paid < due (BR-008, EDGE-018) |
| Unpaid / Partially Paid | Record actual payment (full) | Paid | Paid amount > 0 AND cumulative paid == due; no contract overpayment (PIM-INV-023) (BR-008, EDGE-018, EDGE-019) |
| Partially Paid | Record actual payment (partial) | Partially Paid | Paid amount > 0 AND cumulative paid < due (BR-008) |
| Unpaid / Partially Paid | Due date passes with unpaid amount | Overdue | Business date past due date AND unpaid amount remains (BR-008, BR-021, EDGE-020) |
| Overdue | Record actual payment (full) | Paid | Paid amount > 0 AND cumulative paid == due; no contract overpayment (PIM-INV-023) (BR-008, EDGE-019) |
| Overdue | Record actual payment (partial) | Partially Paid | Paid amount > 0 AND cumulative paid < due (BR-008) |

Invariants:
- PIM-INV-022: Actual payment amount must be greater than zero; zero and negative
  amounts are rejected (BR-008, EDGE-018).
- PIM-INV-023: Cumulative recorded payment cannot exceed the contract's remaining
  amount; overpayment is rejected. Payment Plan due amounts are a planning
  breakdown within the contract, not independent overpayment ceilings (BR-008,
  DEC-014, EDGE-019).
- PIM-INV-024: A single currency applies to all amounts (BR-008, DEC-013,
  PIM-012).
- PIM-INV-025: Won is reached when the related contract is Signed (DEC-017,
  PIM-INV-007); full payment is NOT a precondition for Opportunity Won. Payment
  tracking is retained but decoupled (DEC-019); the overpayment ceiling and
  payment status remain enforced for post-sale collection, independent of closure.
- PIM-INV-026: Overdue requires due date passed AND unpaid amount remaining
  (BR-008, BR-021, EDGE-020). The overdue-evaluation trigger (on-read vs
  scheduled) is deferred (see Open / Blocked PIM-OPEN-002).
- PIM-INV-027: Payment records and overdue transitions are record-local history
  events; payments and overdue payment are operation-log event classes (PIM-018,
  PIM-019, BR-010).

### PIM-SM-007 — Task Status (CIM-032)

States: Open, Completed, Cancelled, Overdue.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create task | Open | Related CRM record, owner, due date, status, title present (BR-003) |
| Open | Complete | Completed | Actor authorized (BR-001) |
| Open | Cancel | Cancelled | Actor authorized (BR-001) |
| Open | Due date passes while incomplete | Overdue | Business date past due date AND task incomplete (BR-013, BR-021, EDGE-023) |
| Overdue | Complete | Completed | Actor authorized (BR-001) |
| Overdue | Cancel | Cancelled | Actor authorized (BR-001) |

Invariants:
- PIM-INV-028: Completed and Cancelled tasks are not active reminders (BR-009,
  BR-013, EDGE-023).
- PIM-INV-029: A task must link to a related CRM record (BR-009, BR-003).
- PIM-INV-030: On owner transfer of a parent record, Open tasks and follow-ups
  transfer to the new owner unless manually reassigned during transfer (BR-009,
  EDGE-024; see PIM-SM-008).
- PIM-INV-031: Task overdue evaluation uses the Business Date; trigger mechanism
  deferred (BR-021, PIM-OPEN-002).

### PIM-SM-008 — Owner Assignment / Transfer (CIM-007, CIM-PROC-014)

States (of the Ownership relation on a record): Unassigned, Assigned. Transfer
re-points Assigned ownership to another team member within scope.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| Unassigned | Assign owner | Assigned | Actor is Sales Manager or Administrator; target within team scope (BR-001, DEC-010) |
| Assigned | Transfer owner | Assigned (new owner) | Actor is Sales Manager or Administrator; target within team scope (BR-001, CIM-PROC-014) |

Invariants:
- PIM-INV-032: Ownership transfers only within team scope (CIM-006, CIM-PROC-014,
  DEC-009).
- PIM-INV-033: On transfer, Open tasks and follow-ups transfer with the parent
  owner unless manually reassigned during transfer (BR-009, EDGE-024,
  PIM-INV-030).
- PIM-INV-034: Owner changes are record-local history events and operation-log
  event classes; managers/admins act as themselves and do not silently act on
  behalf of Sales (BR-010, CIM-PROC-014).

### PIM-SM-009 — Opportunity Close (Won / Lost) (CIM-017, CIM-018, CIM-PROC-011)

This is the closure sub-machine of the Opportunity aggregate, factored out for the
terminal/non-reopen invariants. Reachable from the open Pipeline Stages of
PIM-SM-002.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| Contract Negotiation | Close won | Won | Related contract is Signed; close date and related quote/contract/payment/activity/task history preserved (BR-005, DEC-017, EDGE-009) |
| Any open stage | Close lost | Lost | Lost reason recorded; close date and related history preserved (BR-005, EDGE-010) |
| Won / Lost | Reopen / stage rollback / re-close | (rejected) | Terminal — not allowed in committed scope (BR-005, BR-020, EDGE-011, EDGE-036) |

Invariants:
- PIM-INV-035: Won requires the related contract to be Signed (DEC-017,
  PIM-INV-007, PIM-INV-025); full payment is not a Won precondition (DEC-019).
- PIM-INV-036: Lost requires a lost reason (PIM-INV-008).
- PIM-INV-037: Won and Lost are terminal and non-reopenable; post-close
  stage and close data are not editable through normal Sales workflow,
  though notes and tasks may still be added under normal permissions (BR-020,
  EDGE-036).
- PIM-INV-038: Closure preserves related quote/contract/payment/activity/task
  history (CIM-PROC-011, ACC-013).
- PIM-INV-039: Opportunity Won/Lost are operation-log event classes (PIM-019,
  BR-010).

### PIM-SM-010 — Archive Lifecycle (CIM-037, CIM-038, CIM-PROC-020)

States: Active, Archived. Applies to eligible record types (BR-016). No hard
delete in any state.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| Active | Archive | Archived | Actor is Administrator or Sales Manager (Sales cannot); record eligible; no unresolved active downstream obligations (BR-002, BR-016, EDGE-032, PERM-013) |
| Archived | Retrieve via archived filter | Archived (read) | Actor authorized; explicit archived filter applied (BR-016, EDGE-031) |
| Active / Archived | Hard delete | (rejected) | Hard delete is not allowed in committed scope (BR-002, DEC-011, EDGE-... n/a) |

Invariants:
- PIM-INV-040: No core CRM record may be hard-deleted (BR-002, DEC-011, CIM-039).
- PIM-INV-041: A record with unresolved active downstream obligations (open tasks,
  pending-signature contracts, unpaid payments) cannot be archived until those
  obligations are resolved or archived first (BR-016, EDGE-032, CIM-PROC-020).
- PIM-INV-042: Archived records are excluded from active work lists, active
  reminders, and default operational views and reports, and remain retrievable via
  explicit archived filters, record-local history, operation logs, and audit/
  report evidence (BR-016, BR-017, EDGE-031).
- PIM-INV-043: Won and Lost opportunities may be archived but remain terminal
  (BR-016, PIM-INV-037).
- PIM-INV-044: Archive emits a record-local history event and an operation-log
  event class (PIM-018, PIM-019, BR-010, CIM-PROC-020).

### PIM-SM-011 — User Account Status (CIM-001, CIM-PROC-024)

States: Active, Disabled. Plus Role assignment as a governed change.

| From State | Event/Action | To State | Guard |
|---|---|---|---|
| (none) | Create user | Active | Administrator actor; exactly one Role assigned (CIM-PROC-024, DEC-005) |
| Active | Change role | Active (new role) | Administrator actor; change does not remove the last active Administrator (CIM-PROC-024, BR-001) |
| Active | Disable | Disabled | Administrator actor; change does not deactivate the last active Administrator (CIM-PROC-024) |
| Disabled | Enable | Active | Administrator actor (CIM-PROC-024) |

Invariants:
- PIM-INV-045: A user holds exactly one assigned Role (CIM-001, DEC-005).
- PIM-INV-046: A role change or status change that would remove or deactivate the
  last active Administrator is rejected (CIM-PROC-024, ACC-002, ACC-022).
- PIM-INV-047: Disabled and unauthenticated users cannot access core CRM data
  (BR-001, EDGE-001, ACC-001).
- PIM-INV-048: Role changes, status changes, and last-Administrator-blocked
  outcomes are operation-log event classes (PIM-019, BR-010, CIM-PROC-024).

## Data Classification & Retention (cross-cutting)

This cross-cutting section carries the committed data-classification and
retention policy (privacy-requirements.md PRIV-001..016 + Retention Policy,
COMP-013) into the domain model as invariants anchored to the lifecycle events
already modeled in the state machines above. Concrete durations and storage/TTL
mechanism are deferred to PSM/data-design; PIM fixes only the invariant that
retention is anchored to lifecycle events and never shortened below the
committed minimum.

Invariants:
- PIM-INV-049: Every core CRM record and log carries a committed Data
  Classification (Security Critical / Confidential / Restricted) governing
  visibility, masking, and retention tier (PRIV-001..016, COMP-013).
- PIM-INV-050: Each data category is retained for at least its committed
  minimum, anchored to the modeled lifecycle events (archive PIM-SM-010,
  Won/Lost PIM-SM-009, contract completion/termination PIM-SM-005, full payment
  PIM-SM-006, user deactivation PIM-SM-011), never shortened below the committed
  policy and never hard-deleted (PRIV-*, COMP-013, DEC-011, PIM-INV-040);
  concrete durations and storage/TTL are PSM/data-design.
- PIM-INV-051: Operation logs are retained per their committed tiers
  (business-sensitive vs access/login-failure) and are append-only, not
  hard-deleted (PRIV-011, PRIV-016).
- PIM-INV-052: Transient artifacts (raw import files, generated export files,
  report snapshots) are not retained server-side beyond their committed short
  window unless Architecture defines secure temporary storage with short
  expiration (PRIV-012, PRIV-013, PRIV-014); mechanism is PSM.

## Business Behavior

Operations on domain objects with the governing accepted rule and a test-concept
tag (TEST-* is a concept placeholder only; the test model is authored later).
Permission/contract/error/event mechanisms are PSM.

| ID | Object | Behavior | Rule | Test Concept |
|---|---|---|---|---|
| PIM-BEH-001 | User | Authenticate and bind assigned role to session for authorization | BR-001, EDGE-001, CIM-PROC-001 | TEST-AUTH-LOGIN |
| PIM-BEH-002 | User / Ownership | Enforce three-role visibility scope (admin all / manager team / sales owned) on every protected action | BR-001, DEC-010, CIM-PROC-002, EDGE-002 | TEST-AUTHZ-SCOPE |
| PIM-BEH-003 | User | Create user / change role / change status, preserving last-active-Administrator invariant | CIM-PROC-024, PIM-INV-046 | TEST-USER-ADMIN |
| PIM-BEH-004 | Lead | Create lead (Unassigned or owned) with required fields | BR-003, BR-004, EDGE-003 | TEST-LEAD-CREATE |
| PIM-BEH-005 | Lead | Assign / transfer lead owner; preserve owner-change history | BR-004, CIM-PROC-003, EDGE-024 | TEST-LEAD-ASSIGN |
| PIM-BEH-006 | Lead | Qualify Valid / Invalid (with reason) / Convert; enforce qualification guards | BR-004, EDGE-004, EDGE-005, EDGE-006 | TEST-LEAD-QUALIFY |
| PIM-BEH-007 | Company / Customer | Create/edit company with required fields; no hard delete | BR-002, BR-003 | TEST-CUSTOMER-CRUD |
| PIM-BEH-008 | Contact | Create/link contact requiring related company and a contact method or role note | BR-003, EDGE-007 | TEST-CONTACT-LINK |
| PIM-BEH-009 | Opportunity | Create opportunity with required links/fields | BR-003 | TEST-OPP-CREATE |
| PIM-BEH-010 | Opportunity | Move stage along allowed transitions; reject forbidden transitions | BR-005, EDGE-008, PIM-INV-006 | TEST-OPP-STAGE |
| PIM-BEH-011 | Opportunity | Close Won (related contract Signed — DEC-017) / Lost (reason); enforce terminal non-reopen | BR-005, BR-020, DEC-017, EDGE-009, EDGE-010, EDGE-011, EDGE-036 | TEST-OPP-CLOSE |
| PIM-BEH-012 | Quote | Create/send/reject/expire quote with required fields | BR-003, BR-006 | TEST-QUOTE-LIFECYCLE |
| PIM-BEH-013 | Quote | Accept the opportunity's single quote (exactly one quote per opportunity — DEC-018) | BR-006, DEC-018, PIM-INV-012 | TEST-QUOTE-ACCEPT |
| PIM-BEH-014 | Contract | Create Pending Signature contract from Accepted (non-Expired) quote with note + expected signed date | BR-003, BR-006, BR-007, EDGE-013, EDGE-014 | TEST-CONTRACT-CREATE |
| PIM-BEH-015 | Contract | Sign/activate/complete/terminate requiring signed/effective date for signed states | BR-007, EDGE-015, EDGE-016 | TEST-CONTRACT-LIFECYCLE |
| PIM-BEH-016 | Contract | Record contract-amount difference reason when amount differs from quote | BR-007, DEC-014, EDGE-017 | TEST-CONTRACT-AMOUNT-DIFF |
| PIM-BEH-017 | Payment Plan / Actual Payment | Create plan; record actual payment updating Unpaid/Partially Paid/Paid | BR-008, CIM-PROC-010 | TEST-PAYMENT-RECORD |
| PIM-BEH-018 | Actual Payment | Reject zero/negative and overpayment; enforce single currency | BR-008, DEC-013, DEC-014, EDGE-018, EDGE-019 | TEST-PAYMENT-GUARD |
| PIM-BEH-019 | Payment Plan | Evaluate Overdue when due date passes with unpaid amount | BR-008, BR-021, EDGE-020 | TEST-PAYMENT-OVERDUE |
| PIM-BEH-020 | Activity / Note | Record activity/note against a related CRM record with required fields | BR-003, BR-009, CIM-PROC-012 | TEST-ACTIVITY-NOTE |
| PIM-BEH-021 | Task | Create/complete/cancel task; evaluate Overdue; transfer with parent owner | BR-009, BR-013, EDGE-023, EDGE-024 | TEST-TASK-LIFECYCLE |
| PIM-BEH-022 | Reminder | Derive in-app reminders for due/overdue tasks, pending-signature contracts past expected signed date, due/overdue payments; suppress inactive | BR-013, BR-021, EDGE-020, EDGE-021, EDGE-022, EDGE-023 | TEST-REMINDER |
| PIM-BEH-023 | Ownership | Assign/transfer ownership within team scope; cascade open tasks | BR-001, BR-009, CIM-PROC-014, EDGE-024 | TEST-OWNER-TRANSFER |
| PIM-BEH-024 | Archive State | Archive eligible record blocking on unresolved downstream obligations; exclude from active views; no hard delete | BR-002, BR-016, DEC-011, EDGE-031, EDGE-032 | TEST-ARCHIVE |
| PIM-BEH-025 | Duplicate Warning | Raise non-blocking duplicate warning on normalized company/contact/lead match; allow save | BR-011, BR-019, EDGE-025, EDGE-033, EDGE-034, EDGE-035 | TEST-DUPLICATE-WARN |
| PIM-BEH-026 | CSV Import Job | Validate rows, import valid, report row-level errors without corrupting existing records; Sales denied | BR-012, BR-018, EDGE-026, EDGE-027 | TEST-CSV-IMPORT |
| PIM-BEH-027 | CSV Export Job | Export authorized records only, excluding archived by default; Sales denied | BR-012, BR-018, EDGE-027 | TEST-CSV-EXPORT |
| PIM-BEH-028 | Record-Local History Event | Emit non-editable history events for owner/stage/status/quote/contract/payment/task/archive changes; read by record permission | BR-010, EDGE-029, CIM-PROC-017 | TEST-HISTORY |
| PIM-BEH-029 | Operation-Log Event | Emit Administrator-only non-editable operation-log events for access-sensitive and mutation actions | BR-010, EDGE-029, CIM-PROC-018 | TEST-OPLOG |
| PIM-BEH-030 | Core Retrieval View | List/detail/search/basic-filter across P0 entities with empty-state, invalid-filter feedback, permission hiding | BR-001, CIM-PROC-023, EDGE-002 | TEST-NAV-RETRIEVE |
| PIM-BEH-031 | Team Overview | Aggregate team records and pipeline status for Sales Manager; deny Sales | BR-014, CIM-PROC-015, EDGE-028 | TEST-TEAM-OVERVIEW |
| PIM-BEH-032 | Basic Report | Compute counts/sums for committed groupings over authorized persisted records; default excludes archived; deny Sales | BR-014, BR-017, EDGE-028, EDGE-031 | TEST-BASIC-REPORT |
| PIM-BEH-033 | All aggregates | Persist all core CRM data so it survives refresh/re-login/restart; surface failed saves | BR-015, DEC-008, EDGE-030, CIM-PROC-021 | TEST-PERSISTENCE |
| PIM-BEH-034 | Data Classification / Retention Policy (cross-cutting) | Carry committed classification and minimum retention across each record/log lifecycle; block shortening below the committed minimum; defer durations/storage to PSM | PRIV-001..016, COMP-013, DEC-011 | TEST-RETENTION |

Behavior-mechanism deferral note (modeled in PSM): the concrete command/query/
event/error/permission contracts, event IDs and routing, denial-error shapes,
report computation/query design, persistence design, and reliability mechanisms
that realize the behaviors above are PSM. PIM fixes only the object, the behavior,
and the governing rule/guard/invariant.

## Traceability

Acceptance coverage — every P0/P1 acceptance item with domain-object/behavior/
state content maps to PIM objects, state machines, and behaviors:

| Acceptance | PIM Objects | State Machines | Behaviors |
|---|---|---|---|
| ACC-001 | PIM-001, PIM-002 | PIM-SM-011 | PIM-BEH-001 |
| ACC-002 | PIM-001, PIM-002, PIM-003, PIM-020 | PIM-SM-008, PIM-SM-010, PIM-SM-011 | PIM-BEH-002, PIM-BEH-024 |
| ACC-003 | PIM-004, PIM-003 | PIM-SM-001, PIM-SM-008 | PIM-BEH-004, PIM-BEH-005 |
| ACC-004 | PIM-004 | PIM-SM-001 | PIM-BEH-006 |
| ACC-005 | PIM-005 | PIM-SM-010 | PIM-BEH-007 |
| ACC-006 | PIM-006 | — | PIM-BEH-008 |
| ACC-007 | PIM-007, PIM-012 | PIM-SM-002 (PIM-SM-003 retired — DEC-020) | PIM-BEH-009 |
| ACC-008 | PIM-007 | PIM-SM-002, PIM-SM-009 | PIM-BEH-010, PIM-BEH-011 |
| ACC-009 | PIM-008, PIM-012 | PIM-SM-004 | PIM-BEH-012, PIM-BEH-013 |
| ACC-010 | PIM-009, PIM-012 | PIM-SM-005 | PIM-BEH-014, PIM-BEH-015, PIM-BEH-016 |
| ACC-011 | PIM-010, PIM-011, PIM-012 | PIM-SM-006 | PIM-BEH-017, PIM-BEH-018, PIM-BEH-019 |
| ACC-012 | PIM-013, PIM-014, PIM-015 | PIM-SM-007 | PIM-BEH-020, PIM-BEH-021 |
| ACC-013 | PIM-007, PIM-011 | PIM-SM-009 | PIM-BEH-011 |
| ACC-014 | PIM-018, PIM-020 | (all SMs emit history) | PIM-BEH-028 |
| ACC-015 | PIM-026, PIM-020 | PIM-SM-010 | PIM-BEH-030 |
| ACC-016 | PIM-001..PIM-028 (all persisted records, incl. classification/retention metadata) | — | PIM-BEH-033 |
| ACC-017 | (deployment — PSM/Architecture) | — | (operations evidence; see Open / Blocked) |
| ACC-018 | PIM-024 | — | PIM-BEH-031 |
| ACC-019 | PIM-021 | — | PIM-BEH-025 |
| ACC-020 | PIM-022, PIM-023 | — | PIM-BEH-026, PIM-BEH-027 |
| ACC-021 | PIM-015, PIM-016, PIM-017 | PIM-SM-006, PIM-SM-007, PIM-SM-005 | PIM-BEH-019, PIM-BEH-022 |
| ACC-022 | PIM-019, PIM-001 | PIM-SM-011 | PIM-BEH-003, PIM-BEH-029 |
| ACC-023 | PIM-025, PIM-020 | PIM-SM-010 | PIM-BEH-032 |

Business-rule coverage: BR-001 (PIM-BEH-001/002, SM-008/011), BR-002
(SM-010, INV-040), BR-003 (all create behaviors), BR-004 (SM-001), BR-005
(SM-002/009), BR-006 (SM-004/005), BR-007 (SM-005), BR-008 (SM-006), BR-009
(SM-007, INV-030/033), BR-010 (PIM-018/019, BEH-028/029), BR-011/BR-019 (BEH-025),
BR-012/BR-018 (BEH-026/027), BR-013/BR-021 (BEH-022, INV-026/031), BR-014/BR-017
(BEH-031/032), BR-015 (BEH-033), BR-016 (SM-010), BR-020 (SM-009, INV-037).

Accepted invariants explicitly modeled (with sources cited above): exactly one
quote per opportunity (PIM-INV-012; DEC-018, BR-006); Won when related contract
Signed (PIM-INV-007/025/035; BR-005, DEC-017, EDGE-009; payment decoupled per
DEC-019); Won/Lost terminal and non-reopenable (PIM-INV-009/037; BR-005/BR-020,
EDGE-011/EDGE-036, DEC-012/DEC-017); no hard delete (PIM-INV-040; BR-002,
DEC-011); cannot archive a record
with unresolved active downstream obligations (PIM-INV-041; BR-016, EDGE-032);
cannot remove/deactivate the last active Administrator (PIM-INV-046; CIM-PROC-024);
zero/negative/overpayment rejected (PIM-INV-022/023; BR-008, EDGE-018/EDGE-019);
single currency (PIM-INV-024; BR-008, DEC-013).

Data classification & retention: PRIV-001..016 / COMP-013 retention+classification
are now carried by PIM-027/028, PIM-INV-049..052, and PIM-BEH-034 (anchored to the
lifecycle state machines PIM-SM-005/006/009/010/011), with concrete durations and
storage/TTL deferred to PSM/data-design. These reinforce ACC-014 and ACC-022 (and
the broader PRIV→ACC mapping) without changing any existing ACC mapping.

Counts: 28 domain objects (PIM-001..PIM-028); 11 state machine IDs
(PIM-SM-001..011), of which 10 active and 1 retired (PIM-SM-003 retired 2026-06-01
per DEC-020); 52 invariant IDs (PIM-INV-001..052), of which 51 active and 1
retired (PIM-INV-011 retired 2026-06-01 per DEC-020, folded into stage closure);
34 business behaviors (PIM-BEH-001..034). IDs are retired in place, not renumbered,
to preserve cross-references. All 23 acceptance items ACC-001..ACC-023 are traced
(ACC-017 is deployment/operations, carried as PSM/Architecture + release evidence
per Open / Blocked).

## Open / Blocked

- PIM-OPEN-001 — Second-quote-accept observable outcome — RESOLVED 2026-06-01
  (DEC-018). With exactly one quote per opportunity there is no second quote to
  accept, so the reject-vs-auto-demote ambiguity no longer exists. PIM-SM-004 /
  PIM-INV-012 now state the one-quote rule. Refs: DEC-018, EDGE-012, ACC-009.
- PIM-OPEN-002 — Overdue-evaluation trigger (PSM/Architecture mechanism): the
  Overdue transitions in PIM-SM-006 (Payment) and PIM-SM-007 (Task), and the
  reminder derivation in PIM-BEH-022, fire when "the Business Date passes the due
  date while the item remains unresolved" (BR-021, CIM-034). At PIM altitude this
  is modeled as a guard condition over the Business Date value object (PIM-017):
  the Overdue state is true whenever business-date > due-date AND the item is
  unresolved, independent of how the system observes that moment. WHETHER that
  guard is evaluated on-read or by a scheduled sweep is a mechanism choice
  (deterministic timing for test design) that belongs to PSM/Architecture, not to
  the platform-independent state model. Deferred to PSM with this note. Refs:
  BR-021, CIM-034, EDGE-020, EDGE-021, EDGE-037.
- PIM-OPEN-003 — Opportunity Status enumerated value set — RESOLVED 2026-06-01
  (DEC-020). The separate Opportunity Status dimension is removed; Pipeline Stage
  (PIM-SM-002), including terminal Won/Lost, is the sole lifecycle dimension.
  PIM-SM-003 and PIM-INV-011 are retired in place. No distinct Status enumeration
  exists. Refs: DEC-020, CIM-016 (retired), PRD-007, ACC-007.
- PIM-OPEN-004 — ACC-017 deployment/operations: ACC-017 (deploy and operate) has
  no platform-independent domain-object or state-machine content; deployment
  target, configuration, and backup/restore are Architecture-owned (recorded in
  `docs/architecture/deployment-notes.md`) and release-time operations evidence.
  Modeled here only as the Deployed CRM Instance concept (CIM-046) with no PIM
  state machine; mechanism and evidence are PSM/Architecture/release. Refs:
  ACC-017, CIM-046, CIM-PROC-022.
- PIM-OPEN-005 — Multi-plan contract "fully paid" aggregation — RESOLVED
  2026-06-01 (DEC-017 + DEC-019). Won no longer depends on full payment (Won =
  contract Signed, DEC-017), so "contract fully paid" is no longer a closure
  precondition and the multi-plan aggregation question no longer gates Won. The
  contract-level overpayment ceiling (PIM-INV-023, EDGE-019) and payment status
  remain enforced for post-sale collection (payment retained but decoupled,
  DEC-019), independent of closure. Refs: DEC-017, DEC-019, EDGE-019, CIM-026,
  CIM-028, ACC-011, ACC-013.
