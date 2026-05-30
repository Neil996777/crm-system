# Business Service Governance Inputs

## Document Control

- Project: CRM System
- Phase: G5 Pre-Architecture Input Supplement
- Owner Agent: Business Analyst
- Status: Ready for Architecture Intake
- Date: 2026-05-29
- Sources:
  - `docs/product/business-capability-map.md`
  - `docs/product/acceptance-matrix.md`
  - `docs/business/business-processes.md`
  - `docs/business/business-rules.md`
  - `docs/business/edge-cases.md`
  - `docs/business/role-permission-scenarios.md`
  - `docs/business/business-glossary.md`

## Purpose

This document gives Architecture and Domain Modeling business inputs for
service-boundary-first design. It does not decide final technical services,
deployment boundaries, APIs, database schemas, or implementation modules.

## Cross-Capability Flow Map

| Flow ID | Business Flow | Trigger | Business Capabilities | Required Business Outcome | Failure / Recovery Behavior | Acceptance IDs |
|---|---|---|---|---|---|---|
| BSG-FLOW-001 | Authenticated role-scoped work | User signs in or accesses protected CRM route | CAP-001 plus target capability | User can act only within assigned role and record scope. | Invalid/disabled/unauthenticated users are denied without data exposure. | ACC-001, ACC-002 |
| BSG-FLOW-002 | Lead to opportunity | Lead is qualified as valid and converted | CAP-002, CAP-003, CAP-004, CAP-008 | Lead qualification, customer/contact context, opportunity, and history remain traceable. | Invalid lead cannot convert before restore; converted lead cannot convert again. | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007, ACC-014 |
| BSG-FLOW-003 | Opportunity to quote | Opportunity reaches quote work | CAP-004, CAP-005, CAP-008 | Quote is linked to authorized opportunity/customer and status changes are traceable. | Missing amount/status/validity or unauthorized access blocks mutation. | ACC-007, ACC-008, ACC-009, ACC-014 |
| BSG-FLOW-004 | Accepted quote to contract | User creates contract from accepted quote | CAP-005, CAP-008 | Contract links customer, opportunity, accepted quote, amount, note, expected signed date, and status. | Expired quote blocked; amount mismatch requires difference reason. | ACC-009, ACC-010, ACC-014 |
| BSG-FLOW-005 | Contract to payment to Won | User records payment and closes opportunity | CAP-004, CAP-005, CAP-008 | Payment state reflects actual payments; opportunity can become Won only after full payment. | Zero, negative, overpayment, early Won, or missing lost reason are rejected. | ACC-011, ACC-013, ACC-014 |
| BSG-FLOW-006 | Work item and reminder loop | User creates task or due date arrives | CAP-006, CAP-007, CAP-012 | Authorized user sees active due/overdue tasks, pending-signature contract reminders, and due/overdue payment reminders. | Inactive, unauthorized, signed, terminated, fully paid, archived/default-hidden items do not appear as active reminders. | ACC-012, ACC-021 |
| BSG-FLOW-007 | Record-local history and admin operation logs | Sensitive or business mutation succeeds or denied access occurs where required | CAP-008 plus source capability | Record-local history and/or admin operation log is created and queryable by allowed roles. | Unauthorized history/log access is denied; logs/history are not editable through normal CRM actions. | ACC-014, ACC-022 |
| BSG-FLOW-008 | Team overview and reports | Manager/Admin opens overview or report | CAP-007, CAP-009, CAP-012 | Counts/sums reflect persisted authorized records and default active scope. | Empty state returns zero/empty; unauthorized records are excluded from rows and aggregates. | ACC-018, ACC-023 |
| BSG-FLOW-009 | CSV import/export | Admin/Manager imports or exports CSV | CAP-010, CAP-008, CAP-011 | Valid rows mutate through normal rules; invalid rows are reported safely; exports include authorized records only. | Unsupported format rejected; Sales denied; invalid rows do not corrupt existing records. | ACC-020, ACC-022, ACC-016 |
| BSG-FLOW-010 | Archive and active filtering | Manager/Admin archives eligible record or opens active/default view | CAP-012, CAP-007, CAP-008, CAP-009 | No hard delete; archived records leave active/default work views and remain available through explicit authorized paths. | Active downstream obligations block archive or require resolution/archive first. | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 |
| BSG-FLOW-011 | Production operation and recovery | System is deployed, restarted, backed up, or restored | CAP-011, CAP-008 | CRM remains reachable, configured, persistent, recoverable, and auditable. | Misconfiguration, missing dependency, missing restore proof, or missing environment ownership blocks readiness. | ACC-016, ACC-017, ACC-022 |

## Business Event Inputs

These are business events for Architecture and Domain Modeling to consider.
They are not final event contracts.

| Event Candidate ID | Business Event | Trigger | Required Consumers / Effects | Acceptance IDs |
|---|---|---|---|---|
| BIZ-EVT-001 | User signed in | Successful authentication | Session/role context available; security log candidate. | ACC-001, ACC-022 |
| BIZ-EVT-002 | User access denied | Invalid credential, disabled user, unauthorized route/action | Safe denial; security-relevant operation log where applicable. | ACC-001, ACC-002, ACC-022 |
| BIZ-EVT-003 | User role/status changed | Administrator changes account state | Authorization context invalidation/recheck; operation log. | ACC-001, ACC-002, ACC-022 |
| BIZ-EVT-004 | Lead owner changed | Admin/Manager assigns or transfers lead | Visibility changes; open work transfer where applicable; history/log. | ACC-003, ACC-014, ACC-022 |
| BIZ-EVT-005 | Lead qualified | Lead marked Valid/Invalid/Converted | Lead lifecycle history; conversion eligibility changes. | ACC-004, ACC-014 |
| BIZ-EVT-006 | Lead converted | Valid lead converted to customer/contact/opportunity | Customer/contact/opportunity context created/linked; conversion-once guard. | ACC-004, ACC-005, ACC-006, ACC-007, ACC-014 |
| BIZ-EVT-007 | Opportunity stage changed | Allowed pipeline transition | History; report/overview source state changes. | ACC-008, ACC-014, ACC-018, ACC-023 |
| BIZ-EVT-008 | Opportunity closed | Won or Lost closure succeeds | Terminal state; report/overview update; history. | ACC-013, ACC-014, ACC-023 |
| BIZ-EVT-009 | Quote accepted | Quote status becomes Accepted | Accepted quote uniqueness; contract eligibility; history/log. | ACC-009, ACC-014, ACC-022 |
| BIZ-EVT-010 | Contract status changed | Contract signed, activated, completed, or terminated | Reminder eligibility changes; history/log. | ACC-010, ACC-014, ACC-021, ACC-022 |
| BIZ-EVT-011 | Payment recorded | Actual payment succeeds | Payment status recalculation; opportunity Won eligibility; history/log. | ACC-011, ACC-013, ACC-014, ACC-022 |
| BIZ-EVT-012 | Payment overdue | Due date passes with unpaid amount | Reminder eligibility; history/log where applicable. | ACC-011, ACC-021, ACC-022 |
| BIZ-EVT-013 | Task status changed | Task completed or cancelled | Reminder eligibility changes; history. | ACC-012, ACC-014, ACC-021 |
| BIZ-EVT-014 | Record archived | Eligible record archived | Active/default lists, reminders, reports exclude by default; history/log. | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 |
| BIZ-EVT-015 | Import run completed | CSV import processing completes | Import result summary; row errors; operation log. | ACC-020, ACC-022 |
| BIZ-EVT-016 | Export run completed | CSV export processing completes | Export result; operation log; safe scope/count summary. | ACC-020, ACC-022 |
| BIZ-EVT-017 | Duplicate warning raised | Exact duplicate rule matches lead/company/contact | User warning; no automatic merge or overwrite. | ACC-019 |

## Data Responsibility Inputs

These rows describe business responsibility, not final database ownership.

| Data Area | Business Meaning Owner | Lifecycle / Integrity Rules | Visibility / Security Notes | Related Capabilities |
|---|---|---|---|---|
| User and role | CAP-001 | One assigned role; active/disabled state; last active Administrator protection. | Administrator governs users; non-admin cannot manage users. | CAP-001 |
| Lead | CAP-002 | Owner required before Pending Qualification or later; Invalid restore before conversion; conversion once only. | Sales sees owned/assigned leads only; manager/admin broader scope. | CAP-002, CAP-012 |
| Company/customer | CAP-003 | Required company name, status, owner; no hard delete. | Access through role/scope and related owned/assigned context. | CAP-003, CAP-012 |
| Contact | CAP-003 | Must link to company/customer; contact name and contact method or role note required. | Contact methods must not leak in unauthorized states or unsafe summaries. | CAP-003 |
| Opportunity | CAP-004 | Required customer, owner, stage, status, amount, expected close date; Won/Lost terminal. | Unauthorized records excluded from rows and aggregates. | CAP-004, CAP-009 |
| Quote | CAP-005 | Multiple quotes allowed; only one Accepted quote; expired quote cannot create contract. | Quote amount and status are restricted by record scope. | CAP-005 |
| Contract | CAP-005 | Pending Signature requires expected signed date and note; signed lifecycle states require signed/effective date; amount mismatch requires reason. | Contract amount, notes, and dates are restricted data. | CAP-005, CAP-006 |
| Payment plan / actual payment | CAP-005 | Positive amounts only; no overpayment; partial/full/overdue status rules. | Payment values are restricted and must not leak through reports/errors/import summaries. | CAP-005, CAP-009 |
| Activity / note / task | CAP-006 | Related record required; active reminders only for eligible open/due items. | Visibility follows related record permission. | CAP-006, CAP-012 |
| Record-local history | CAP-008 | Append-only through normal CRM workflows; follows related record lifecycle. | Visible only by related record permission. | CAP-008 |
| Admin/global operation log | CAP-008 | Append-only; Administrator query only; includes access-sensitive and operation-sensitive events. | Sales Manager and Sales denied. | CAP-008 |
| Import/export run | CAP-010 | Authorized CSV only; row-level validation; export scope/count evidence. | Sales denied; row errors and export summaries must be safe. | CAP-010 |
| Report/overview metrics | CAP-009 | Derived from persisted authorized records; active records by default. | Unauthorized rows and aggregates excluded. | CAP-009 |
| Deployment, backup, restore evidence | CAP-011 | Must prove persistent CRM data and audit evidence survive required operations. | Secrets and backup credentials are not product data and must not be exposed. | CAP-011 |

## Architecture Handoff Blockers

| Blocker ID | Severity | Issue | Required Owner |
|---|---|---|---|
| BSG-BLK-001 | P0 | Final service boundaries, owner agents, contracts, and data ownership do not exist yet. | Architecture |
| BSG-BLK-002 | P0 | OQ-001 production target and environment ownership remains reopened. | Architecture, Infrastructure Ops |
| BSG-BLK-003 | P1 | OQ-016 initial seed/migration requirement remains open for launch planning. | Product Manager, Business Analyst |

