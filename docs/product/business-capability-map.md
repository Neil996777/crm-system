# Business Capability Map

## Document Control

- Project: CRM System
- Phase: G5 Pre-Architecture Input Supplement
- Owner Agents: Product Manager, Business Analyst
- Status: Ready for Architecture Intake
- Date: 2026-05-29

## Purpose

This document maps committed P0/P1 CRM acceptance items to business
capabilities so Architecture can define service boundaries without guessing.

The `SVC-CAND-*` values below are service-boundary inputs only. They are not
approved final services, final deployment units, database ownership decisions,
or implementation tasks. Architecture must confirm, split, merge, rename, or
replace them through G5 design while preserving P0/P1 acceptance.

## Capability Map

| Capability ID | Business Capability | Purpose | Primary Actors | Acceptance IDs | Service Candidate Inputs | Notes for Architecture |
|---|---|---|---|---|---|---|
| CAP-001 | Identity and role access | Authenticate users and provide role context for all protected CRM work. | Administrator, Sales Manager, Sales | ACC-001, ACC-002 | SVC-CAND-IDENTITY, SVC-CAND-AUTHZ | Must support active/disabled user checks, session behavior, and role version changes. |
| CAP-002 | Lead intake and qualification | Capture leads, assign owners, qualify leads, and convert valid leads into downstream sales context. | Sales, Sales Manager, Administrator | ACC-003, ACC-004, ACC-019 | SVC-CAND-LEAD, SVC-CAND-DUPLICATE | Must preserve unassigned-lead behavior, restore rules, conversion-once guard, and duplicate warning behavior. |
| CAP-003 | Account and contact management | Maintain ToB company/customer records and multiple contacts under each customer context. | Sales, Sales Manager, Administrator | ACC-005, ACC-006, ACC-019 | SVC-CAND-ACCOUNT, SVC-CAND-DUPLICATE | Must preserve ownership, related child visibility, duplicate warnings, and no-hard-delete rule. |
| CAP-004 | Opportunity pipeline | Track sales opportunities, allowed stage transitions, terminal outcomes, and closure rules. | Sales, Sales Manager | ACC-007, ACC-008, ACC-013 | SVC-CAND-OPPORTUNITY | Must enforce forbidden transitions, Won/Lost terminal behavior, lost reason, and full-payment-before-Won. |
| CAP-005 | Commercial execution | Manage quote, contract, payment plan, actual payment, and commercial lifecycle integrity. | Sales, Sales Manager | ACC-009, ACC-010, ACC-011, ACC-013 | SVC-CAND-COMMERCIAL | Candidate may be refined by Architecture; must preserve quote acceptance, contract note/date rules, amount difference reason, overpayment block, and payment status rules. |
| CAP-006 | Work activity and reminders | Record activities, notes, tasks, and in-app reminders for due/overdue work. | Sales, Sales Manager | ACC-012, ACC-021 | SVC-CAND-WORK | Must preserve related-record requirement, active/inactive reminder rules, and authorized reminder filtering. |
| CAP-007 | Core CRM navigation and record retrieval | Provide list, detail, search, filter, and role-scoped navigation for committed CRM records. | Administrator, Sales Manager, Sales | ACC-015 | SVC-CAND-QUERY-EXPERIENCE | Architecture may implement through query/read models or service APIs; must preserve permission filtering and empty/error states. |
| CAP-008 | Collaboration history and operation audit | Provide record-local business history and Administrator global operation logs. | Administrator, Sales Manager, Sales | ACC-014, ACC-022 | SVC-CAND-HISTORY-AUDIT | Must preserve append-only history/log behavior and role/scope visibility. |
| CAP-009 | Team overview and reports | Provide team overview and basic sales reports based on persisted authorized records. | Sales Manager, Administrator | ACC-018, ACC-023 | SVC-CAND-REPORTING | Must exclude unauthorized and default archived records unless explicit authorized archived filter is provided. |
| CAP-010 | Data import and export | Support authorized CSV import/export with validation, row-level errors, and safe summaries. | Administrator, Sales Manager | ACC-020 | SVC-CAND-IMPORT-EXPORT | Must use authorization-before-mutation/export and preserve partial failure behavior. |
| CAP-011 | Persistence and production operation | Ensure committed CRM data persists and the v1 CRM can be deployed and operated with real configuration. | Administrator / Operator | ACC-016, ACC-017 | SVC-CAND-PLATFORM-OPS | Must resolve OQ-001 and define production environment, backup, restore, and ownership. |
| CAP-012 | Archive and lifecycle governance | Preserve no-hard-delete behavior, eligible archive behavior, active/default filtering, and lifecycle auditability. | Administrator, Sales Manager | ACC-002, ACC-014, ACC-015, ACC-021, ACC-023 | SVC-CAND-LIFECYCLE-GOVERNANCE | Candidate may map into other services; Architecture must define ownership and forbidden data access. |

## Cross-Capability Flow Inputs

| Flow ID | Flow | Capabilities | Primary Acceptance IDs | Primary Flow Owner Candidate | Notes |
|---|---|---|---|---|---|
| FLOW-CAP-001 | User signs in and works under role scope | CAP-001, all protected capabilities | ACC-001, ACC-002 | Architecture to assign | All protected actions depend on identity and authorization context. |
| FLOW-CAP-002 | Lead becomes sales opportunity | CAP-002, CAP-003, CAP-004, CAP-008 | ACC-003, ACC-004, ACC-005, ACC-006, ACC-007, ACC-014 | Architecture to assign | Must preserve conversion history and prevent duplicate conversion. |
| FLOW-CAP-003 | Opportunity advances to quote and contract | CAP-004, CAP-005, CAP-008 | ACC-007, ACC-008, ACC-009, ACC-010, ACC-014 | Architecture to assign | Must preserve accepted quote constraint and expired quote block. |
| FLOW-CAP-004 | Contract payment closes opportunity | CAP-004, CAP-005, CAP-008 | ACC-011, ACC-013, ACC-014 | Architecture to assign | Won requires full payment; overpayment is blocked. |
| FLOW-CAP-005 | Activity and reminder loop | CAP-006, CAP-007, CAP-008, CAP-012 | ACC-012, ACC-021 | Architecture to assign | Reminders must hide unauthorized and inactive records. |
| FLOW-CAP-006 | Import/export operational flow | CAP-001, CAP-010, CAP-008, CAP-011 | ACC-020, ACC-022, ACC-016 | Architecture to assign | Import/export must be durable, authorized, and auditable. |
| FLOW-CAP-007 | Reports and overview | CAP-001, CAP-007, CAP-009, CAP-012 | ACC-018, ACC-023 | Architecture to assign | Reports must use persisted authorized records and default active scope. |
| FLOW-CAP-008 | Production operation and recovery | CAP-011, CAP-008 | ACC-016, ACC-017, ACC-022 | Infrastructure Ops + Architecture to coordinate | Deployment, persistence, backup, restore, logs, and ownership must be defined. |

## Architecture Handoff Notes

- Architecture must produce the final service list, owner agents, contracts,
  data ownership, and deployment boundaries.
- Candidate IDs in this document are allowed to change during G5, but the
  P0/P1 acceptance items they support must not be downgraded or weakened.
- Missing final service owner, contract, data ownership, or service-chain test
  path remains a blocker before G8.

