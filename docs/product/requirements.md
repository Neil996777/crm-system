# Requirements

## Requirement Notes

This accepted baseline captures the confirmed G1/G2 discussion for a production ToB CRM
system. The v1 scope is team-based and must support a complete business loop:
lead, qualification, customer/company, contact, opportunity, quote, contract,
payment, activity, task, and win/loss closure.

P0 and P1 requirements are governed by the workspace no-downgrade rule. Once
accepted, they cannot be downgraded, deleted, merged away, weakened, or marked
complete without the required evidence.

## Priority Definition For This Project

| Priority | Meaning In This Project |
|---|---|
| P0 | Required for the v1 CRM to be valid and production-launchable. |
| P1 | Required for the committed v1 release quality, but not the minimum existence of the CRM loop. |
| P2 | Important enhancement planned after the core v1 loop unless promoted by the sponsor. |
| P3 | Future improvement or convenience capability. |

## Accepted Requirements

| ID | Priority | Requirement | Source | Status |
|---|---|---|---|---|
| REQ-001 | P0 | Users must be able to log in and operate under one of three roles: Administrator, Sales Manager, or Sales. | Sponsor discussion | Accepted as Architecture Input |
| REQ-002 | P0 | The system must enforce role-based access and actions for Administrator, Sales Manager, and Sales users. | Sponsor discussion | Accepted as Architecture Input |
| REQ-003 | P0 | Users must be able to create, view, edit, search, filter, and assign leads with owner, source, company/contact information, need summary, and status. | Sponsor discussion | Accepted as Architecture Input |
| REQ-004 | P0 | Users must be able to qualify leads as valid, invalid, or needing follow-up, with recorded result and reason where relevant. | Sponsor discussion | Accepted as Architecture Input |
| REQ-005 | P0 | Users must be able to manage ToB companies/customers and distinguish prospects from converted or active customers. | Sponsor discussion | Accepted as Architecture Input |
| REQ-006 | P0 | Users must be able to manage multiple contacts under a company/customer, including role, title, contact method, and notes. | Sponsor discussion | Accepted as Architecture Input |
| REQ-007 | P0 | Users must be able to create and manage sales opportunities linked to company/customer, contacts, owner, amount, expected close date, stage, and status. | Sponsor discussion | Accepted as Architecture Input |
| REQ-008 | P0 | Users must be able to move opportunities through a sales pipeline covering new opportunity, qualification, quote, contract, payment, won, and lost outcomes. | Sponsor discussion | Accepted as Architecture Input |
| REQ-009 | P0 | Users must be able to create and manage quote records linked to opportunities, customers, amount, validity period, status, and owner. | Sponsor discussion | Accepted as Architecture Input |
| REQ-010 | P0 | Users must be able to create and manage contract records linked to customer, opportunity, and quote, including amount, status, signed/effective date, attachment or notes. | Sponsor decision DEC-006 | Accepted as Architecture Input |
| REQ-011 | P0 | Users must be able to manage payment plans and actual payment records linked to contracts, including due amount, due date, paid amount, payment date, and payment status. | Sponsor discussion | Accepted as Architecture Input |
| REQ-012 | P0 | Users must be able to record activities, notes, and follow-up tasks against leads, customers, contacts, opportunities, quotes, contracts, or payments where applicable. | Sponsor discussion | Accepted as Architecture Input |
| REQ-013 | P0 | Users must be able to close opportunities as won or lost, preserving quote, contract, payment, activity, and task history. | Sponsor discussion | Accepted as Architecture Input |
| REQ-014 | P0 | Authorized team members must be able to review collaboration history, ownership changes, stage changes, and key business updates. | Sponsor discussion | Accepted as Architecture Input |
| REQ-015 | P0 | Core CRM entities must provide list, detail, search, and basic filtering views. | Sponsor discussion | Accepted as Architecture Input |
| REQ-016 | P0 | All core CRM data must be persisted and must survive refresh, logout/login, and service restart. | Workspace no-downgrade rule | Accepted as Architecture Input |
| REQ-017 | P0 | The v1 system must be deployable to a production target with real configuration and real persisted data. | Sponsor discussion | Accepted as Architecture Input |
| REQ-018 | P1 | Sales Managers should have a team overview for leads, opportunities, quotes, contracts, payments, tasks, and pipeline status. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-019 | P1 | The system should warn users about likely duplicate companies, contacts, or leads during creation or update. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-020 | P1 | The system should support data import/export for core CRM records. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-021 | P1 | The system should provide reminders for due or overdue follow-up tasks, contracts, and payments. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-022 | P1 | Administrators should be able to review key operation logs for access, ownership, pipeline, quote, contract, and payment changes. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-023 | P1 | Sales Managers and Administrators should have basic sales reports for leads, opportunities, quotes, contracts, and payments. | Product recommendation accepted by sponsor | Accepted as Architecture Input |
| REQ-024 | P2 | The system may support email and calendar integration. | Product recommendation | Accepted as Architecture Input |
| REQ-025 | P2 | The system may support advanced reporting, forecasting, and sales performance analytics. | Product recommendation | Accepted as Architecture Input |
| REQ-026 | P2 | The system may support quote approval, contract approval, discount approval, and related workflow rules. | Sponsor decision DEC-006 | Accepted as Architecture Input |
| REQ-027 | P2 | The system may support electronic signature and contract template generation. | Sponsor decision DEC-006 | Accepted as Architecture Input |
| REQ-028 | P2 | The system may support invoice management. | Product recommendation | Accepted as Architecture Input |
| REQ-029 | P2 | The system may support external collaboration and finance integrations. | Product recommendation | Accepted as Architecture Input |
| REQ-030 | P2 | The system may support AI sales summaries, next-step suggestions, and risk hints. | Product recommendation | Accepted as Architecture Input |

## P0 Business Loop Coverage

| Loop Step | Requirement IDs |
|---|---|
| Login and role access | REQ-001, REQ-002 |
| Lead entry and assignment | REQ-003 |
| Lead qualification | REQ-004 |
| Company/customer and contacts | REQ-005, REQ-006 |
| Opportunity and pipeline | REQ-007, REQ-008 |
| Quote | REQ-009 |
| Contract | REQ-010 |
| Payment | REQ-011 |
| Activities and tasks | REQ-012 |
| Win/loss closure | REQ-013 |
| Collaboration history | REQ-014 |
| Lists, details, search, filters | REQ-015 |
| Persistence and launch readiness | REQ-016, REQ-017 |

## Initial Business State Candidates

These states are draft business-analysis inputs and must be confirmed or
refined before G3/G4.

| Object | States |
|---|---|
| Lead | Unassigned, Pending Qualification, Valid, Invalid, Converted To Opportunity |
| Opportunity | New Opportunity, Needs Confirmed, Quote, Contract Negotiation, Contract Signed, Payment In Progress, Won, Lost |
| Quote | Draft, Sent, Accepted, Rejected, Expired |
| Contract | Pending Signature, Signed, Active, Completed, Terminated |
| Payment | Unpaid, Partially Paid, Paid, Overdue |
