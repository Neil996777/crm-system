# CIM

Computation Independent Model.

## Document Control

- Project: CRM System
- Date: 2026-05-31
- Role: Domain Modeling
- Gate: G6 (MDA Modeling)
- Scope note: Models accepted design; not design authority for scope. This CIM
  traces the accepted G4 business design and G5 architecture inputs. It does not
  introduce new scope, rules, fields, or concepts. Where the domain authority
  needs something the accepted sources do not provide, it is recorded under
  `## Open / Blocked` rather than invented. Revised 2026-05-31 per independent
  4-role CIM audit (BA, Security, QA-TDD, PM); altitude discipline applied.
  Data classification/retention (COMP-013, privacy-requirements.md) added
  2026-06-01 to close the pre-audit coverage gap; durations referenced not copied.
  Amended 2026-06-01 per Formal Scope Change by User (decision-log.md
  DEC-017..020): Won = related contract Signed (not full payment); exactly one
  quote per opportunity; payment tracking retained but decoupled from Won;
  Opportunity Status field removed (Pipeline Stage is the sole lifecycle
  dimension). `Payment In Progress` and `Contract Signed` pipeline stages removed.
  Affected IDs retired in place (not renumbered) to preserve cross-references.

## Business Concepts

| ID | Concept | Definition | Source | Acceptance ID |
|---|---|---|---|---|
| CIM-001 | User | A person with a CRM account holding exactly one assigned role, authenticated per session and authorized on protected routes/actions. | business-glossary (Owner), PRD-001, DEC-005, ACC-001 | ACC-001, ACC-002 |
| CIM-002 | Role | The access classification assigned to a user; committed model is exactly three roles. | business-glossary, PRD-002, DEC-005 | ACC-001, ACC-002 |
| CIM-003 | Administrator | Role responsible for user/role governance, full governed CRM visibility, and global operation logs. | business-glossary, PRD-001, PRD-002, PERM-002 | ACC-001, ACC-002, ACC-022 |
| CIM-004 | Sales Manager | Role responsible for team record visibility, assignment/transfer, team overview, and team reports; cannot manage users/roles or view global audit logs. | business-glossary, PRD-002, PRD-018, DEC-010, PERM-003 | ACC-002, ACC-018, ACC-023 |
| CIM-005 | Sales | Role responsible for owned/assigned leads, customers, opportunities, quotes, contracts, payments, activities, notes, and tasks only; cannot archive, import/export, or view global logs/manager views. | business-glossary, PRD-002, DEC-010, PERM-004, PERM-012, PERM-026 | ACC-002 |
| CIM-006 | Team / Organization | The single committed team/organization that owns all CRM records; defines the boundary of Sales Manager "team scope". | DEC-002, DEC-009 | ACC-002, ACC-018 |
| CIM-007 | Ownership / Assignment | The relation binding a CRM record to a responsible/assigned user; drives Sales visibility and scope. | business-glossary (Owner, Assigned Record), PRD-002, DEC-010, PERM-007 | ACC-002, ACC-003 |
| CIM-008 | Lead | A potential ToB sales target before qualification or conversion; may be Unassigned before assignment. | business-glossary, PRD-003, PRD-004, DEC-001 | ACC-003, ACC-004 |
| CIM-009 | Lead Status | The qualification/lifecycle state of a lead: Unassigned, Pending Qualification, Valid, Invalid, Converted To Opportunity. | business-glossary (Unassigned/Qualified/Invalid Lead), PRD-003, PRD-004 | ACC-003, ACC-004 |
| CIM-010 | Lead Source | The origin attribute recorded on a lead, required at creation. | PRD-003, ACC-003 | ACC-003 |
| CIM-011 | Company / Customer | A ToB account or organization record used for contacts, opportunities, quotes, contracts, and payments; carries company name, customer status, and owner. | business-glossary, PRD-005 | ACC-005 |
| CIM-012 | Customer Status | The status attribute persisted on a company/customer. | business-glossary (Company/Customer), PRD-005, ACC-005 | ACC-005 |
| CIM-013 | Contact | A person or role linked under a company/customer; requires contact name, related company/customer, and at least one contact method or role note. | business-glossary, PRD-006 | ACC-006 |
| CIM-014 | Opportunity | A sales deal record linked to customer, contacts, owner, expected amount, expected close date, and stage; moves through the pipeline. (Status dimension removed 2026-06-01 per DEC-020; Pipeline Stage is the sole lifecycle dimension.) | business-glossary, PRD-007, DEC-020 | ACC-007, ACC-008 |
| CIM-015 | Opportunity / Pipeline Stage | A business state describing opportunity progress through the sales pipeline, named: New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost (Won/Lost terminal). (`Contract Signed` and `Payment In Progress` stages removed 2026-06-01 per DEC-017.) The stage transition matrix is modeled in PIM. | business-glossary (Pipeline Stage), PRD-008, DEC-017 | ACC-008 |
| CIM-016 | Opportunity Status | **[RETIRED 2026-06-01 DEC-020 — Opportunity `Status` field removed; Pipeline Stage (CIM-015), including terminal Won/Lost, is the sole opportunity lifecycle dimension.]** ID retained for cross-reference stability; not renumbered. | PRD-007, DEC-020 | (retired) |
| CIM-017 | Won | Terminal opportunity outcome reached when the related contract is Signed (DEC-017); cannot be reopened in the committed scope. Full payment is not a Won precondition (payment is post-sale follow-up, DEC-019). | business-glossary, PRD-013, DEC-017 | ACC-013 |
| CIM-018 | Lost | Terminal opportunity outcome requiring a lost reason; cannot be reopened in the committed scope. | business-glossary, PRD-013, DEC-012, DEC-017 | ACC-013 |
| CIM-019 | Quote | A commercial offer linked to opportunity and customer, with amount, status, validity end date, and owner; exactly one quote per opportunity (DEC-018). | business-glossary, PRD-009, DEC-018 | ACC-009 |
| CIM-020 | Quote Status | The lifecycle state of a quote, named: Draft, Sent, Accepted, Rejected, Expired; exactly one quote per opportunity (DEC-018), so at most one Accepted quote follows trivially. Transition rules are modeled in PIM. | business-glossary (Accepted/Expired Quote), PRD-009, DEC-018 | ACC-009 |
| CIM-021 | Contract | A record-based contract linked to customer, opportunity, and Accepted quote; carries amount, status, required contract note, expected signed date, and (for signed states) signed/effective date. No approval, e-signature, or template generation. | business-glossary, PRD-010, DEC-006, DEC-007, DEC-016 | ACC-010 |
| CIM-022 | Contract Status | The lifecycle state of a contract: Pending Signature, Signed, Active, Completed, Terminated. | business-glossary (Pending Signature/Signed/Effective), PRD-010 | ACC-010 |
| CIM-023 | Expected Signed Date | Planned signature deadline on a Pending Signature contract; drives signature reminders; not a substitute for signed/effective date. | business-glossary, PRD-010, PRD-021 | ACC-010, ACC-021 |
| CIM-024 | Signed / Effective Date | Date required once a contract becomes Signed, Active, Completed, or post-signature Terminated. | business-glossary, PRD-010 | ACC-010 |
| CIM-025 | Contract Note | The textual note that is P0-required on a contract record. | business-glossary, PRD-010, DEC-016 | ACC-010 |
| CIM-026 | Payment Plan | A planned payment record linked to a contract with due amount, due date, and status; basis for due/overdue tracking. | business-glossary, PRD-011 | ACC-011 |
| CIM-027 | Actual Payment | A recorded payment event linked to a contract; zero, negative, and overpayment amounts are rejected. | business-glossary, PRD-011, DEC-014 | ACC-011 |
| CIM-028 | Payment Status | The state of payment against a contract/plan: Unpaid, Partially Paid, Paid, Overdue. | business-glossary (Partial/Overdue Payment), PRD-011, DEC-014 | ACC-011 |
| CIM-029 | Amount / Money Model | Single-currency monetary value used across opportunity, quote, contract, and payment; excludes tax, discount, and multi-currency automation. | business-glossary (the *Amount terms), PRD-007/009/010/011, DEC-013 | ACC-007, ACC-009, ACC-010, ACC-011 |
| CIM-030 | Activity | A business interaction record linked to a CRM record, preserving follow-up history. | business-glossary, PRD-012 | ACC-012 |
| CIM-031 | Note | A textual business note linked to a CRM record, visible by record permission. | business-glossary, PRD-012 | ACC-012 |
| CIM-032 | Task | A follow-up work item with owner, due date, status, and title; due/overdue tasks create reminders. Task Status is named: Open, Completed, Cancelled, Overdue. Transition rules are modeled in PIM. | business-glossary, PRD-012 | ACC-012, ACC-021 |
| CIM-033 | Reminder / Follow-up | An in-app notification for due/overdue tasks, Pending Signature contracts past expected signed date, and due/overdue payments; in-app only in the committed release. | business-glossary, PRD-021, DEC-015 | ACC-021 |
| CIM-034 | Business Date | Workspace-local date used to evaluate due/overdue for reminders until Architecture defines timezone handling. The overdue-evaluation trigger (on-read vs scheduled) is an Architecture-pending input (see Open / Blocked) needed for deterministic overdue test design. | business-glossary, PRD-021 | ACC-021 |
| CIM-035 | Record-Local History | The permitted business timeline of events (owner, stage, status, quote, contract, payment, task, archive changes) visible from a related CRM record by record permission; not editable through normal CRM actions. | business-glossary, PRD-014, NFR-004 | ACC-014 |
| CIM-036 | Admin / Global Operation Log | Administrator-only operational audit query across records and access-sensitive actions, recording the acting user, the action, the affected resource, when it occurred, the result, and before/after values; not editable through normal CRM actions. Concrete log fields, event identifiers, and schema are modeled in PSM. | business-glossary, PRD-022, NFR-005 | ACC-022 |
| CIM-037 | Archive (action) | A non-delete lifecycle action available to Administrator and Sales Manager for eligible records; Sales cannot archive. | business-glossary, PRD-002, PERM-012, PERM-013 | ACC-002 |
| CIM-038 | Archived Record | A record removed from active work views, active reminders, and default operational views without hard deletion; retrievable via explicit archived filters and audit/history views. A record with unresolved active downstream obligations (open tasks, pending-signature contracts, unpaid payments) cannot be archived until those obligations are resolved or archived first (see CIM-PROC-020). | business-glossary, ACC-014, ACC-015, EDGE-032, BR-016 | ACC-014, ACC-015, ACC-023 |
| CIM-039 | Hard Delete (prohibited) | Permanent deletion of core CRM records; not allowed in the committed scope. | business-glossary, PRD-002, DEC-011 | ACC-002 |
| CIM-040 | Duplicate Warning | A non-blocking warning raised when company/contact/lead data matches configured duplicate rules; no silent overwrite or automatic merge. | business-glossary, PRD-019 | ACC-019 |
| CIM-041 | CSV Import Job | An authorized bulk-entry operation in CSV format that validates rows and reports row-level errors; Sales cannot import. | business-glossary (CSV Import), PRD-020, DEC-015 | ACC-020 |
| CIM-042 | CSV Export Job | An authorized export operation in CSV format including only authorized records; Sales cannot export. | business-glossary (CSV Export), PRD-020, DEC-015 | ACC-020 |
| CIM-043 | Team Overview | A Sales Manager view of all team leads, opportunities, quotes, contracts, payments, tasks, and pipeline status. | PRD-018, PERM-019 | ACC-018 |
| CIM-044 | Basic Report | Counts and sums for committed CRM groupings (leads by status, opportunities by stage, quotes/contracts/payments by status/amount) using persisted authorized records; default excludes archived records. | business-glossary (Basic/Active Report), PRD-023, NFR-007 | ACC-023 |
| CIM-045 | Persistent CRM Record / Persistence | The requirement that all core CRM data is persisted (not mock/static/in-memory) and survives refresh, re-login, and service restart. | PRD-016, NFR-001, DEC-008 | ACC-016 |
| CIM-046 | Deployed CRM Instance | The CRM operated in the target environment with real configuration connected to persistent services. | PRD-017, NFR-003, DEC-004 | ACC-017 |
| CIM-047 | Core Record Navigation / Retrieval View | The role-scoped business view for finding and opening committed CRM records via list, detail, search, and basic filter across the P0 entities (leads, companies/customers, contacts, opportunities, quotes, contracts, payments, activities, tasks); shows empty state on no results, validation feedback on invalid filter, hides permission-out-of-scope records, and denies detail access to unauthorized records. Query/API design is deferred to PSM. | business-glossary, PRD-015, CAP-007 | ACC-015 |
| CIM-048 | Data Classification | The committed sensitivity class (Security Critical / Confidential / Restricted) on each CRM data category, governing visibility, masking, and retention tier; accepted vocabulary, not a new rule. | privacy-requirements.md PRIV-001..016, COMP-013 | ACC-014, ACC-022 (and the per-category PRIV→ACC mapping) |
| CIM-049 | Data Retention Expectation | The committed minimum retention per CRM data category, anchored to lifecycle events (archive, opportunity closure, contract completion/termination, full payment, user deactivation), never shortened below the committed policy, consistent with no-hard-delete; concrete durations live in privacy-requirements.md (Retention Policy), storage/TTL deferred to PSM/data-design. | privacy-requirements.md Retention Policy + PRIV-001..016, COMP-013, DEC-011 | ACC-014, ACC-022 |

## Business Processes

| ID | Process | Actors | Trigger | Outcome |
|---|---|---|---|---|
| CIM-PROC-001 | Login and role entry (BP-001) | Administrator, Sales Manager, Sales | A user signs in to the CRM. | Authenticated identity is persisted for the session and the assigned role governs all subsequent record-visibility and action decisions; invalid/disabled/unauthenticated access is denied without exposing data. (ACC-001, ACC-002) |
| CIM-PROC-002 | Access-control enforcement (BP-001/BP-008) | Administrator, Sales Manager, Sales | A user attempts create/view/edit/assign/close/archive/import/export/report/audit. | Administrator governs all records, Sales Manager manages all team records, Sales manages only owned/assigned records; unauthorized actions are denied without exposing or mutating data; no hard delete. (ACC-002) |
| CIM-PROC-003 | Lead capture and assignment (BP-002) | Sales, Sales Manager, Administrator | A lead is created, imported, or assigned. | Lead is persisted with name/company name, source, and status; owner is required before Pending Qualification or later; Unassigned leads may exist; owner-change history is preserved. (ACC-003, ACC-014, ACC-016) |
| CIM-PROC-004 | Lead qualification and conversion (BP-002) | Sales, Sales Manager | A lead in Pending Qualification is acted on. | Lead is marked Valid, Invalid (with reason), or Converted To Opportunity (with downstream link); conversion preserves original lead history; Sales cannot act on Unassigned leads; converted leads cannot convert again; Invalid leads convert only after Administrator/Sales Manager restores them. (ACC-004, ACC-014) |
| CIM-PROC-005 | Duplicate detection on create/edit (BP-002/BP-003) | Sales, Sales Manager | A user enters lead/company/contact data matching configured duplicate rules. | A non-blocking duplicate warning is shown (company name, contact phone/email, lead company/contact); legitimate save proceeds; no silent overwrite or automatic merge. (ACC-019) |
| CIM-PROC-006 | Customer and contact setup (BP-003) | Sales, Sales Manager | A valid lead is converted or a user creates customer/contact records directly. | Company/customer is persisted with company name, customer status, and owner; one or more contacts (name, related company/customer, contact method or role note) are linked and visible in authorized views; missing required links block save. (ACC-005, ACC-006, ACC-014, ACC-016) |
| CIM-PROC-007 | Opportunity pipeline management (BP-004) | Sales, Sales Manager | A valid business need is tracked as an opportunity / stage change requested. | Opportunity is persisted with customer, owner, stage, expected amount, expected close date (no separate Status dimension — DEC-020); stage changes persist and create record-local history; the allowed stage-transition matrix (New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost) is modeled in PIM. (ACC-007, ACC-008, ACC-014) |
| CIM-PROC-008 | Quote lifecycle including accept/reject (BP-005) | Sales, Sales Manager | An opportunity reaches the quote stage. | Exactly one quote per opportunity (DEC-018) is persisted with amount, status, validity end date, and owner; the quote can be sent/accepted/rejected/expired; an expired quote cannot link to a new contract. (ACC-009, ACC-014) |
| CIM-PROC-009 | Contract from accepted quote (BP-005) | Sales, Sales Manager | A quote is Accepted and a contract is created. | A Pending Signature contract is persisted linked to customer, opportunity, and Accepted quote with amount, required note, and expected signed date; signed/effective date is required when status becomes Signed/Active/Completed/post-signature Terminated; amount differing from the quote requires a recorded difference reason. (ACC-010, ACC-014) |
| CIM-PROC-010 | Payment recording with overpayment block (BP-006) | Sales, Sales Manager | A contract has payment terms or an actual payment to record. | Payment plans and actual payments are persisted with contract, amount, date, and status; partial payment sets Partially Paid, full payment sets Paid, past-due unpaid sets Overdue; zero/negative/overpayment amounts are rejected; single currency applies. Payment tracking is retained as post-sale collection follow-up and does NOT gate Opportunity Won (DEC-019). (ACC-011, ACC-014) |
| CIM-PROC-011 | Close opportunity Won or Lost (BP-004/BP-006) | Sales, Sales Manager | A user closes an opportunity. | Won is reached when the related contract is Signed (DEC-017) — full payment is not a precondition; Lost requires a lost reason; close date and related quote/contract/payment/activity/task history are preserved; Won and Lost are terminal and cannot be reopened. Post-signing breach is handled at contract level (Terminated); the opportunity stays Won. (ACC-013, ACC-014) |
| CIM-PROC-012 | Activities, notes, and tasks (BP-007) | Sales, Sales Manager | A user records follow-up, collaboration, or a work item against a CRM record. | Activity/note/task is persisted with related entity, actor/owner, timestamp, content/title, due date where applicable, and status; missing related record or required fields block save; unauthorized create/view is denied. (ACC-012) |
| CIM-PROC-013 | In-app reminders for due/overdue work (BP-007) | Sales, Sales Manager | A due/overdue task, a Pending Signature contract past expected signed date, or a due/overdue payment exists; or the user opens the reminder area. | Authorized in-app reminders are shown; completed/cancelled tasks, signed/terminated/fully-paid contracts do not create active reminders; unauthorized records are hidden. (ACC-021) |
| CIM-PROC-014 | Owner assignment and transfer (BP-008) | Sales Manager, Administrator | A manager assigns or transfers ownership of a team record. | Ownership transfers within team scope; open tasks and follow-ups transfer with the parent owner unless manually reassigned; owner-change history is recorded; managers/admins act as themselves. (ACC-002, ACC-014) |
| CIM-PROC-015 | Team overview (BP-008) | Sales Manager | Sales Manager opens the team overview. | Manager sees all team leads, opportunities, quotes, contracts, payments, tasks, and pipeline status from persisted authorized records; Sales users are denied; empty data shows empty state. (ACC-018) |
| CIM-PROC-016 | CSV import and export (BP-009) | Administrator, Sales Manager | An authorized user runs a CSV import or export. | Valid CSV rows are imported with success summary; invalid rows reported with row-level errors and no corruption; export includes only authorized records; Sales cannot import/export; operation is logged. (ACC-020) |
| CIM-PROC-017 | Record-local history review (BP-010) | Administrator, Sales Manager, Sales | An authorized user opens record-local history/collaboration context. | Permitted history events (acting user, what changed on the resource, when, and before/after values) are shown by record permission; unauthorized history is denied; history is not editable via normal CRM actions. Concrete event fields/identifiers are modeled in PSM. (ACC-014) |
| CIM-PROC-018 | Admin/global operation-log review (BP-010) | Administrator | Administrator opens the global operation-log query after key actions. | Administrator sees global operation events with required fields, covering the business event classes: login success and login/access failures, user role change, user status change, last-Administrator-blocked, owner changes, stage/status changes, quote acceptance, contract changes, payments, archive, import, and export; Sales Manager and Sales are denied; logs are not editable via normal CRM actions. (Traceability: audit-log-spec.md Event Catalog and security-requirements.md Sensitive Operations; concrete event IDs and event ID schema/routing are modeled in PSM.) (ACC-001, ACC-002, ACC-022) |
| CIM-PROC-019 | Basic reporting / manager overview metrics (BP-010) | Administrator, Sales Manager | An authorized user opens basic reports. | Counts and sums (leads by status, opportunities by stage, quotes/contracts/payments by status/amount) traceable to persisted authorized records; default excludes archived records; Sales users are denied. (ACC-023) |
| CIM-PROC-020 | Archive and active-work filtering (BP-011) | Administrator, Sales Manager | An authorized user archives an eligible record or applies an archived filter. | Record is archived (no hard delete), archive history/operation-log event recorded, and removed from active lists/reminders/default views; archived records retrievable via explicit archived filter and audit/history views; archiving a record with unresolved active downstream obligations (open tasks, pending-signature contracts, unpaid payments) is blocked, or requires resolving/archiving those obligations first (EDGE-032, BR-016, PERM-013; enforcement mechanism modeled in PIM/PSM); Sales cannot archive. (ACC-002, ACC-014, ACC-015, ACC-021, ACC-023) |
| CIM-PROC-021 | Persist all core CRM data (BP-002..011) | Administrator, Sales Manager, Sales | A user refreshes, logs out/in, or the service restarts after data is created/changed. | Previously saved data remains available per permissions; failed saves are surfaced and never appear as successful persistent changes; no P0 path relies on mock/static/in-memory behavior. (ACC-016) |
| CIM-PROC-022 | Deploy and operate the CRM | Administrator / Operator | Architecture has selected the production target and configuration; the system is deployed and smoke-tested. | CRM is reachable, configured, and connected to persistent services; misconfiguration or unavailable dependencies block production readiness. Backup/restore and exact provider/domain details are Architecture-owned operations evidence (RISK-002, OQ-001). (ACC-017) |
| CIM-PROC-023 | Core record navigation and retrieval (BP-003/BP-011) | Administrator, Sales Manager, Sales | An authenticated user lists, opens, searches, or filters committed CRM records. | Authorized list/detail/search/basic-filter views are available across the P0 entities (leads, companies/customers, contacts, opportunities, quotes, contracts, payments, activities, tasks); an empty result shows an empty state; an invalid filter shows validation feedback; records outside the user's permission scope are hidden from lists, and unauthorized detail access is denied. Query/API design is deferred to PSM. (ACC-015) |
| CIM-PROC-024 | User and role administration (BP-001) | Administrator | Administrator creates a user, changes a user's role, or changes a user's status (enable/disable). | The user/role change is persisted and governs subsequent authorization; the business invariant holds that a change cannot remove or deactivate the last active Administrator (such a change is rejected); role changes, status changes, and last-Administrator-blocked outcomes are recorded as operation-log events (see CIM-PROC-018). Enforcement mechanism is modeled in PIM/PSM. (PM-007, ABUSE-005, SEC-008) (ACC-001, ACC-002, ACC-022) |

## Open / Blocked

- Backup/restore operations: ACC-016/ACC-017 reference a `CONTRACT-CAND-BACKUP-RESTORE`
  contract candidate, but the accepted business-process and PRD sources do not
  define a business-level backup/restore process with actors, triggers, and
  outcomes. The production deployment target and backup configuration are now
  decided and recorded in `docs/architecture/deployment-notes.md`, so OQ-001's
  architecture decision is no
  longer open; what remains is release-time operations evidence (encrypted
  off-server backup copy + restore rehearsal), which is a release gate, not a CIM
  modeling input. CIM-PROC-022 models deployment/operation at the business level
  and flags backup/restore as Architecture-owned rather than inventing a business
  process. Refs: ACC-016, ACC-017, `docs/architecture/deployment-notes.md`; OQ-001
  (architecture decision recorded, release evidence pending); PRD-017.
- Second-quote-accept observable outcome — RESOLVED 2026-06-01 by DEC-018 (exactly
  one quote per opportunity): with only one quote possible there is no second quote
  to accept, so the reject-vs-auto-demote ambiguity no longer exists. CIM-PROC-008
  now states the one-quote rule. Refs: DEC-018, EDGE-012, ACC-009.
- Overdue-evaluation trigger (Architecture-pending): deterministic overdue test
  design for tasks and payments needs to know whether overdue is evaluated on-read
  vs on a schedule. This is an Architecture/PIM mechanism concern, not CIM
  vocabulary; CIM-034 names the Business Date basis only. Pending input refs:
  BR-021, CIM-034.
- Audit as a distinct "process": there is no standalone audit business process in
  the accepted sources beyond record-local history (CIM-PROC-017) and the
  admin/global operation log (CIM-PROC-018). The G12 independent audit is a
  platform governance activity, not an in-product CRM process, so it is not
  modeled as a CIM business process.
