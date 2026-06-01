# PRD

## Document Control

- Product: CRM System
- Owner Agent: Product Manager
- Phase: PRD / Acceptance Matrix
- Current Gate: G5 Architecture Design Required
- Date: 2026-05-26
- Source Documents:
  - `docs/product/project-charter.md`
  - `docs/product/requirements.md`
  - `docs/product/open-questions.md`
  - `docs/product/out-of-scope.md`
  - `docs/product/decision-log.md`
  - `docs/product/acceptance-matrix.md`

## Overview

The CRM System is a production-ready team collaboration CRM for ToB sales. The
committed release must support the complete sales business loop from lead entry through
customer, contact, opportunity, quote, contract, payment, follow-up, task
management, and won/lost closure.

The product is not a demo or prototype. P0/P1 capabilities must use real
persistent data and cannot be satisfied by mock data, static-only screens, TODO
placeholders, or non-persistent behavior.

This PRD defines product intent and release scope. It does not decide final
technology stack, database schema, API contract, infrastructure, or UI visual
design. Those decisions belong to downstream modeling, UX/UI, security, and
architecture work after the relevant gates.

## Goals

- Provide one shared CRM workspace for a ToB sales team.
- Let users manage the full committed CRM loop: lead, qualification, company/customer,
  contact, opportunity, quote, contract, payment, activity, task, and closure.
- Support the three confirmed roles: Administrator, Sales Manager, and Sales.
- Preserve ownership, stage, status, quote, contract, payment, activity, and
  task history for team collaboration and auditability.
- Make P0/P1 scope verifiable through the acceptance matrix before task
  planning.
- Keep implementation blocked until Gate G8 passes.

## Non-Goals

- Contract approval workflow.
- Electronic signature.
- Contract template generation.
- Quote approval, discount approval, or complex approval rules.
- Invoice management.
- Email or calendar synchronization.
- External integrations with Feishu, DingTalk, WeCom, ERP, or finance systems.
- Advanced analytics, sales forecasting, or performance prediction.
- AI customer summary, sales advice, or risk prediction.
- Dedicated mobile app.
- Multi-tenant SaaS organization management unless formally promoted later.
- Technology stack or architecture selection in this PRD.

Detailed exclusions are maintained in `docs/product/out-of-scope.md`.

## Users And Roles

| Role | Description | Primary Goal | Priority |
|---|---|---|---|
| Administrator | System-level operator for users, roles, full CRM data visibility, and critical configuration. | Keep the CRM governed, auditable, and operational. | P0 |
| Sales Manager | Team sales lead who monitors team work, assigns or reviews records, and tracks pipeline progress. | Coordinate team sales work and detect deal, task, contract, or payment risk. | P0 |
| Sales | Sales team member responsible for assigned leads, customers, opportunities, quotes, contracts, payments, activities, and tasks. | Execute daily sales follow-up and keep customer/deal records complete. | P0 |

## Core Business Loop

1. A user logs in as Administrator, Sales Manager, or Sales.
2. A lead is created or assigned.
3. The lead is qualified as valid, invalid, or needing follow-up.
4. A company/customer and contacts are created or linked.
5. An opportunity is created for the qualified business need.
6. The opportunity moves through the sales pipeline.
7. A quote is created and tracked.
8. A contract record is created or linked.
9. Payment plans and actual payment records are tracked.
10. Activities, notes, and follow-up tasks preserve the working history.
11. The opportunity is closed as won or lost.
12. Authorized team members review the historical record for continued customer follow-up.

## Core User Journeys

### Journey 1: Sales Handles A New Lead

Role:
- Sales

Main path:
1. Sales logs in.
2. Sales creates or opens an assigned lead.
3. Sales records source, company/contact details, need summary, owner, and status.
4. Sales qualifies the lead.
5. If valid, Sales links or creates the company/customer and contacts.
6. Sales creates an opportunity from the qualified need.

Success result:
- The lead qualification and downstream customer/opportunity context are persisted and traceable.

Failure and edge behavior:
- Missing required fields prevent save.
- Unauthorized lead access is denied.
- Invalid qualification or conversion behavior follows business rules.

Acceptance coverage:
- ACC-001, ACC-002, ACC-003, ACC-004, ACC-005, ACC-006, ACC-007, ACC-016

### Journey 2: Sales Advances An Opportunity To Quote And Contract

Role:
- Sales

Main path:
1. Sales opens an opportunity.
2. Sales updates stage as the deal progresses.
3. Sales creates a quote linked to the customer and opportunity.
4. Sales records quote amount, validity period, status, and owner.
5. Sales creates a contract record linked to the customer, opportunity, and quote.
6. Sales records contract amount, status, signed/effective date, and attachment or notes.

Success result:
- Quote and contract history are persisted and visible from related records.

Failure and edge behavior:
- Invalid stage transitions are rejected.
- Invalid amount/date/status values are rejected.
- Missing required links block valid quote or contract completion.

Acceptance coverage:
- ACC-007, ACC-008, ACC-009, ACC-010, ACC-014, ACC-016

### Journey 3: Sales Tracks Payment And Closes The Opportunity

Role:
- Sales

Main path:
1. Sales opens a contract.
2. Sales creates payment plan records.
3. Sales records actual payments.
4. Sales reviews payment status.
5. Sales closes the opportunity as won or lost according to the PRD closure rules.

Success result:
- Payment plan, actual payment records, and opportunity closure are persisted and traceable.

Failure and edge behavior:
- Partial, overdue, overpayment, and invalid amount cases follow business rules.
- Close status requires required close data.
- Won/lost history remains preserved.

Acceptance coverage:
- ACC-011, ACC-013, ACC-014, ACC-016

### Journey 4: Sales Manager Reviews Team Work

Role:
- Sales Manager

Main path:
1. Sales Manager logs in.
2. Sales Manager reviews authorized team leads, opportunities, quotes, contracts, payments, tasks, and pipeline status.
3. Sales Manager opens details to inspect ownership, stage, status, and follow-up history.
4. Sales Manager acts on authorized assignments or review needs.

Success result:
- Sales Manager can understand team sales progress and risks without accessing unauthorized data.

Failure and edge behavior:
- Empty data, permission-denied records, and unavailable summaries are handled clearly.
- Visibility rules follow the permission matrix in this PRD.

Acceptance coverage:
- ACC-001, ACC-002, ACC-014, ACC-015, ACC-018, ACC-023

### Journey 5: Administrator Governs Access And Audit

Role:
- Administrator

Main path:
1. Administrator logs in.
2. Administrator reviews users, roles, and authorized CRM data.
3. Administrator reviews key operation logs for ownership, pipeline, quote, contract, payment, and access-sensitive changes.

Success result:
- CRM governance and audit-sensitive activity can be reviewed.

Failure and edge behavior:
- Logs are not editable through normal CRM actions.
- Unauthorized users cannot access administrator-only views.

Acceptance coverage:
- ACC-001, ACC-002, ACC-014, ACC-022

## Functional Requirements

| ID | Priority | Requirement | Source Requirement | Acceptance ID |
|---|---|---|---|---|
| PRD-001 | P0 | Users can log in and operate under Administrator, Sales Manager, or Sales roles. | REQ-001 | ACC-001 |
| PRD-002 | P0 | The system enforces role-based access and action control across P0 CRM records. | REQ-002 | ACC-002 |
| PRD-003 | P0 | Users can create, view, edit, search, filter, and assign leads. | REQ-003 | ACC-003 |
| PRD-004 | P0 | Users can qualify leads and record qualification result and reason where applicable. | REQ-004 | ACC-004 |
| PRD-005 | P0 | Users can manage ToB companies/customers and customer status. | REQ-005 | ACC-005 |
| PRD-006 | P0 | Users can manage multiple contacts under a company/customer. | REQ-006 | ACC-006 |
| PRD-007 | P0 | Users can create and manage opportunities linked to customer, contacts, owner, amount, expected close date, and stage. _(amended 2026-06-01 DEC-020: `status` dimension removed)_ | REQ-007 | ACC-007 |
| PRD-008 | P0 | Users can move opportunities through the sales pipeline. | REQ-008 | ACC-008 |
| PRD-009 | P0 | Users can create and manage quote records linked to opportunities and customers. | REQ-009 | ACC-009 |
| PRD-010 | P0 | Users can create and manage record-based contract records linked to customer, opportunity, and quote. | REQ-010 | ACC-010 |
| PRD-011 | P0 | Users can manage payment plans and actual payment records linked to contracts. | REQ-011 | ACC-011 |
| PRD-012 | P0 | Users can record activities, notes, and follow-up tasks against CRM records. | REQ-012 | ACC-012 |
| PRD-013 | P0 | Users can close opportunities as won or lost while preserving business history. | REQ-013 | ACC-013 |
| PRD-014 | P0 | Authorized users can review collaboration and change history for core CRM records. | REQ-014 | ACC-014 |
| PRD-015 | P0 | Core CRM entities provide list, detail, search, and basic filtering views. | REQ-015 | ACC-015 |
| PRD-016 | P0 | All core CRM data is persisted and survives refresh, logout/login, and service restart. | REQ-016 | ACC-016 |
| PRD-017 | P0 | The system can be deployed and operated with real configuration and real persisted data. | REQ-017 | ACC-017 |
| PRD-018 | P1 | Sales Managers can view a team overview for core sales work. | REQ-018 | ACC-018 |
| PRD-019 | P1 | Users receive likely duplicate warnings for companies, contacts, or leads. | REQ-019 | ACC-019 |
| PRD-020 | P1 | Authorized users can import and export core CRM records. | REQ-020 | ACC-020 |
| PRD-021 | P1 | Users receive reminders for due or overdue follow-up tasks, contracts, and payments. | REQ-021 | ACC-021 |
| PRD-022 | P1 | Administrators can review key operation logs. | REQ-022 | ACC-022 |
| PRD-023 | P1 | Administrators and Sales Managers can view basic sales reports for core CRM records. | REQ-023 | ACC-023 |
| PRD-024 | P2 | The system may support email and calendar integration. | REQ-024 | Not in committed acceptance scope |
| PRD-025 | P2 | The system may support advanced reporting, forecasting, and performance analytics. | REQ-025 | Not in committed acceptance scope |
| PRD-026 | P2 | The system may support quote approval, contract approval, discount approval, and workflow rules. | REQ-026 | Not in committed acceptance scope |
| PRD-027 | P2 | The system may support electronic signature and contract template generation. | REQ-027 | Not in committed acceptance scope |
| PRD-028 | P2 | The system may support invoice management. | REQ-028 | Not in committed acceptance scope |
| PRD-029 | P2 | The system may support external collaboration and finance integrations. | REQ-029 | Not in committed acceptance scope |
| PRD-030 | P2 | The system may support AI sales summaries, next-step suggestions, and risk hints. | REQ-030 | Not in committed acceptance scope |

## Non-Functional Requirements

| ID | Priority | Requirement | Completion Standard | Acceptance ID |
|---|---|---|---|---|
| NFR-001 | P0 | Persistent core data | P0 CRM data survives refresh, logout/login, and service restart; no core path depends on mock, static-only, TODO, in-memory-only, or non-persistent behavior. | ACC-016 |
| NFR-002 | P0 | Authorization enforcement | Protected P0 data and actions are enforced by role and record visibility rules, not only hidden in the UI. | ACC-002 |
| NFR-003 | P0 | Production deployment readiness | The committed system can be deployed with real configuration and connected persistent services. | ACC-017 |
| NFR-004 | P0 | Audit-sensitive history | Ownership, stage, quote, contract, payment, and key business changes are retained for collaboration and audit review. | ACC-014 |
| NFR-005 | P1 | Operational audit query | Administrator can review key operation logs for committed P1 audit needs. | ACC-022 |
| NFR-006 | P1 | Data import/export integrity | Import errors do not corrupt existing data; export reflects authorized persisted records. | ACC-020 |
| NFR-007 | P1 | Report traceability | Basic sales report numbers are traceable to persisted CRM records and documented business definitions. | ACC-023 |
| NFR-008 | P1 | Reminder correctness | Due and overdue reminders reflect configured timing rules and exclude completed/cancelled items. | ACC-021 |

## Permissions Summary

The confirmed PRD-level permission model is a single team workspace with three
roles. Detailed security implementation belongs to Security Compliance and
Architecture, but the following product-level allow/deny rules are required for
G3 acceptance testability.

| Role | Capability | Allowed | Notes |
|---|---|---|---|
| Administrator | Manage users, roles, and critical CRM configuration | Yes | Exact user lifecycle and configuration scope remain downstream. |
| Administrator | View and govern CRM records and operation logs | Yes | Must respect audit and privacy requirements. |
| Sales Manager | View team CRM records | Yes | The committed scope has one sales team; manager scope is all team records. |
| Sales Manager | Assign or transfer team work | Yes | Open tasks and follow-ups transfer with the parent owner unless manually reassigned. |
| Sales | Manage owned or assigned CRM records | Yes | Sales can view and edit records where they are owner or assignee. |
| Sales | View non-owned CRM records | No | Non-owned data is hidden unless later shared by a formal permission rule. |
| Sales | View global audit logs | No | Administrator-only unless later changed by formal decision. |
| Unauthenticated user | Access core CRM data | No | Must be rejected. |

### Minimum Permission Matrix

| Entity / Action | Administrator | Sales Manager | Sales |
|---|---|---|---|
| User and role management | Full | No | No |
| Lead create | Yes | Yes | Yes, assigned to self by default |
| Lead view | All | All team | Owned/assigned only |
| Lead edit | All | All team | Owned/assigned only |
| Lead assign/transfer | All | All team | No |
| Company/customer create | Yes | Yes | Yes |
| Company/customer view | All | All team | Related to owned/assigned records only |
| Company/customer edit | All | All team | Related to owned/assigned records only |
| Contact create/edit | All | All team | Related to owned/assigned company/customer only |
| Opportunity create/edit | All | All team | Owned/assigned only |
| Opportunity sales closure as Won/Lost | All | All team | Owned/assigned only |
| Quote create/edit/status | All | All team | Owned/assigned opportunity only |
| Contract create/edit/status | All | All team | Owned/assigned opportunity only |
| Payment plan/record create/edit | All | All team | Owned/assigned contract only |
| Activity/note/task create/edit | All | All team | Owned/assigned related record only |
| Record-local history view | All | All team | Owned/assigned related record only |
| Admin/global operation log view | Yes | No | No |
| Import/export | Yes | Yes for team records | No |
| Basic reports | Yes | Team reports | No |
| Hard delete core CRM data | No | No | No |
| Archive eligible records | Yes | Yes for team records | No |

Rules:
- The committed scope does not allow hard deletion of core CRM records.
- Administrators and Sales Managers act as themselves; they do not silently act
  as another Sales user.
- Any owner, stage, status, amount, contract, payment, archive, import, or
  permission-sensitive action must create a record-local history event where
  applicable and an admin/global operation-log event where applicable.

## Required Fields

| Entity | Required Fields For Valid Save |
|---|---|
| Lead | lead name or company name, source, status; owner is required before Pending Qualification or later states |
| Company/customer | company name, customer status, owner |
| Contact | contact name, related company/customer, at least one contact method or role note |
| Opportunity | related company/customer, owner, stage, expected amount, expected close date |
| Quote | related opportunity, related company/customer, quote amount, status, validity end date, owner |
| Contract | related customer, related opportunity, accepted quote, contract amount, status, contract note; expected signed date is required while status is Pending Signature; signed/effective date is required when status is Signed, Active, Completed, or Terminated after signing |
| Payment plan | related contract, due amount, due date, status |
| Actual payment | related contract, paid amount, payment date, status |
| Activity/note | related CRM record, activity type or note type, content, actor, timestamp |
| Task | related CRM record, owner, due date, status, title |

## Business State And Transition Rules

These states are product-level rules for G3 acceptance. Business Analyst,
Domain Modeling, UX, Security, and QA may refine them later only by preserving
or strengthening P0/P1 behavior.

| Object | States |
|---|---|
| Lead | Unassigned, Pending Qualification, Valid, Invalid, Converted To Opportunity |
| Opportunity | New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Won, Lost _(amended 2026-06-01 DEC-017: Won = related contract Signed; `Contract Signed` and `Payment In Progress` stages removed; `status` dimension removed per DEC-020)_ |
| Quote | Draft, Sent, Accepted, Rejected, Expired |
| Contract | Pending Signature, Signed, Active, Completed, Terminated |
| Payment | Unpaid, Partially Paid, Paid, Overdue |
| Task | Open, Completed, Cancelled, Overdue |

### Lead Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| Unassigned | Pending Qualification | Administrator, Sales Manager | owner | Owner assigned |
| Pending Qualification | Valid | Sales, Sales Manager | qualification result | Lead qualified |
| Pending Qualification | Invalid | Sales, Sales Manager | invalid reason | Lead disqualified |
| Valid | Converted To Opportunity | Sales, Sales Manager | related company/customer and opportunity | Lead converted |

Forbidden:
- Unassigned leads may exist only before assignment and cannot be qualified,
  edited by Sales, or converted to opportunity.
- Invalid leads cannot convert to opportunity unless first changed back to Pending Qualification by Administrator or Sales Manager.
- Converted leads cannot be deleted or converted again.

### Opportunity Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| New Opportunity | Needs Confirmed | Sales, Sales Manager | related customer/contact | Stage changed |
| Needs Confirmed | Quote | Sales, Sales Manager | need summary and expected amount | Stage changed |
| Quote | Contract Negotiation | Sales, Sales Manager | the quote is Sent or Accepted | Stage changed |
| Contract Negotiation | Won | Sales, Sales Manager | related contract is Signed | Opportunity won |
| Any non-terminal stage | Lost | Sales, Sales Manager | lost reason | Opportunity lost |

Forbidden:
- Won and Lost are terminal for the committed scope. Reopen is not allowed in the committed scope.
- Opportunity reaches Won when the related contract is Signed; it cannot reach Won without a Signed contract.

Amended 2026-06-01 (DEC-017): Won is reached at contract signing, not full payment;
the `Contract Signed` and `Payment In Progress` stages are removed. Post-signing
breach is handled at the contract level (Terminated); the opportunity stays Won.

### Quote Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| Draft | Sent | Sales, Sales Manager | amount and validity end date | Quote sent |
| Sent | Accepted | Sales, Sales Manager | accepted date | Quote accepted |
| Sent | Rejected | Sales, Sales Manager | rejection reason | Quote rejected |
| Sent | Expired | System or authorized user | validity end date passed | Quote expired |

Rules:
- Each opportunity has exactly one quote (DEC-018).
- A contract can only reference the opportunity's Accepted quote.
- Expired quotes cannot be linked to new contracts.

Amended 2026-06-01 (DEC-018): one quote per opportunity; the system records the
quote result, not multi-round negotiation; previous multiple-quotes / one-Accepted
constraint removed.

### Contract Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| Pending Signature | Signed | Sales, Sales Manager | signed/effective date and contract amount | Contract signed |
| Signed | Active | Sales, Sales Manager | effective date reached or manually activated | Contract active |
| Active | Completed | Sales, Sales Manager | full payment recorded | Contract completed |
| Pending Signature, Signed, Active | Terminated | Sales Manager, Administrator | termination reason | Contract terminated |

Rules:
- Contract notes are P0 and required; attachment upload is not required for P0.
- Pending Signature contracts require customer, opportunity, Accepted quote,
  amount, status, contract note, and expected signed date; they do not require
  signed/effective date.
- Expected signed date is the planned signature deadline used by contract
  reminders. It is not a substitute for signed/effective date.
- Signed, Active, Completed, and post-signature Terminated contracts require
  signed/effective date.
- Contract amount may differ from accepted quote amount only when a difference reason is recorded.
- Contract approval, electronic signature, and template generation are not part of committed P0/P1.

### Payment Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| Unpaid | Partially Paid | Sales, Sales Manager | payment amount greater than 0 and less than contract remaining amount | Payment recorded |
| Unpaid | Paid | Sales, Sales Manager | payment amount equals contract amount | Payment recorded |
| Partially Paid | Paid | Sales, Sales Manager | cumulative paid amount equals contract amount | Payment completed |
| Unpaid, Partially Paid | Overdue | System or authorized user | due date is past and unpaid amount remains | Payment overdue |

Rules:
- The committed scope uses one currency for all quote, contract, and payment amounts.
- Tax, discount, and multi-currency automation are out of committed P0/P1 scope.
- Overpayment is blocked.
- Negative or zero actual payment amounts are rejected.
- Partial payment is supported.

### Task Transitions

| From | To | Actor | Required Data | History Event |
|---|---|---|---|---|
| Open | Completed | Task owner, Sales Manager, Administrator | completion timestamp | Task completed |
| Open | Cancelled | Task owner, Sales Manager, Administrator | cancellation reason | Task cancelled |
| Open | Overdue | System or authorized user | due date passed and task not completed/cancelled | Task overdue |

Rules:
- Completed and Cancelled tasks are not active reminders.
- On owner transfer of a parent record, open tasks and follow-ups transfer to
  the new owner unless manually reassigned during transfer.

## P1 Minimum Behavior

| Capability | Minimum Committed Behavior |
|---|---|
| Duplicate warnings | Warn on exact company name match; warn on contact phone or email match; warn on lead company/contact match. Warning does not block save. |
| Import/export | CSV only for the committed release. Import validates required fields row by row and reports failed rows without changing valid existing records. Export includes authorized records only. |
| Reminders | In-app reminders only for due or overdue tasks, contracts pending signature past their required expected signed date, and due/overdue payments. Completed/cancelled tasks, signed contracts, terminated contracts, and fully paid contracts do not create active reminders. |
| Admin operation logs | Global admin query covers login/access failures, owner changes, stage/status changes, quote acceptance, contract status changes, payment records, archive actions, import/export. |
| Basic reports | Counts and sums for leads by status, opportunities by stage, quotes by status/amount, contracts by status/amount, payments by status/amount. Reports are based only on persisted authorized records. |

## History And Audit Distinction

| Item | Priority | Purpose | Visibility |
|---|---|---|---|
| Record-local business history | P0 | Shows business timeline on the related lead, customer, opportunity, quote, contract, payment, activity, or task. | Visible according to record permissions. |
| Admin/global operation log | P1 | Gives Administrators a cross-record operational audit query. | Administrator only. |

Shared event IDs for the committed release:
- EVT-OWNER-CHANGED
- EVT-STAGE-CHANGED
- EVT-STATUS-CHANGED
- EVT-QUOTE-ACCEPTED
- EVT-CONTRACT-SIGNED
- EVT-CONTRACT-TERMINATED
- EVT-PAYMENT-RECORDED
- EVT-PAYMENT-OVERDUE
- EVT-OPPORTUNITY-WON
- EVT-OPPORTUNITY-LOST
- EVT-RECORD-ARCHIVED
- EVT-IMPORT-RUN
- EVT-EXPORT-RUN

## Release Boundary

The committed release is valid only when every P0 and P1 item in this PRD has:

- acceptance matrix coverage
- downstream business, UX, security, modeling, architecture, task, and test
  coverage where applicable
- implementation evidence
- QA evidence
- integration evidence
- audit pass
- no open P0/P1 blocker

`Implemented`, `QA Verified`, `Integration Verified`, or `Audit Passed` alone
does not mean `Done`.

## Risks

| ID | Risk | Impact | Mitigation |
|---|---|---|---|
| RISK-001 | Security and privacy details require downstream design ownership. | Security requirements, permission matrix, data classification, retention, and abuse cases may block later gates if not kept traceable. | OQ-014 has been resolved by Security Design in `docs/security/privacy-requirements.md`; downstream work may strengthen security controls without weakening P0/P1 behavior. |
| RISK-002 | Production provider and operations details require provisioning evidence. | Production readiness cannot pass later gates without configured provider, domain, backup, and environment ownership evidence. | Reopened by architecture reset on 2026-05-29. Architecture must reselect or confirm the production target before downstream MDA and task planning. |
| RISK-003 | P1 reporting/reminder expectations could expand. | Scope pressure may weaken G3 acceptance focus. | Preserve P0/P1/P2 boundaries and require formal scope change for promotions/removals. |
| RISK-004 | Business rules may need downstream strengthening. | Business, modeling, UX, security, and QA may find gaps during G4 work. | Downstream agents may add stricter rules, but cannot downgrade or weaken P0/P1 behavior. |

## Open Questions

Open questions are maintained in `docs/product/open-questions.md`. G3-blocking
product/business/testability questions have been resolved in this PRD and the
acceptance matrix. Remaining questions are downstream gate inputs.

| ID | Question Summary | Owner | Status |
|---|---|---|---|
| OQ-001 | Exact production provider/domain/backup/operations details. | Architecture | Reopened by architecture reset on 2026-05-29; must be resolved by new Architecture Design. |
| OQ-014 | Data classification and retention expectations. | Security Compliance | Resolved by Security Design in `docs/security/privacy-requirements.md` |
| OQ-016 | Data migration or initial seed data before launch. | Product Manager / Business Analyst | Watch item for ACC-017; launch planning pending before release |

Any unresolved question that affects a P0/P1 acceptance item must block the
later gate that requires that decision.

## Restart Traceability Notes

- Every P0/P1 PRD item maps to an acceptance item. New MDA traces, G8 tasks,
  and test model entries must be recreated after the new architecture is
  accepted.
- `docs/product/acceptance-matrix.md` remains the source of truth for
  completion standards.
- Implementation is blocked until the restarted delivery flow reaches and
  passes G8 again.
- No P0/P1 requirement may be weakened, merged away, or accepted as partial
  implementation.
