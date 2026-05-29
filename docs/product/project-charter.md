# Project Charter

## Project

- Name: CRM System
- Sponsor: User / Project Sponsor
- Date: 2026-05-26
- Status: Accepted as Architecture Input

## Business Goal

Build a production-ready ToB CRM system for a sales team. The v1 release must
cover the complete sales business loop from lead capture through customer,
opportunity, quote, contract, payment, follow-up, and final win/loss tracking.

The system is not a demo or prototype. Core CRM paths must use persistent data
and must not be satisfied by mock data, static-only screens, TODO placeholders,
or non-persistent behavior.

## Target Users

| Role | Goal | Key Pain | Priority |
|---|---|---|---|
| Administrator | Manage system users, roles, and critical CRM configuration. | Without centralized access control, customer data and ownership are hard to govern. | P0 |
| Sales Manager | Track team pipeline, assign work, monitor progress, and review sales outcomes. | Team sales work is hard to coordinate without shared visibility and accountability. | P0 |
| Sales | Manage assigned leads, customers, opportunities, quotes, contracts, payments, activities, and tasks. | Sales follow-up and deal history can be lost or fragmented across personal notes and messages. | P0 |

## Core Business Loop

1. A team member logs in with an assigned role.
2. A lead is created or assigned to a sales owner.
3. The lead is qualified as valid, invalid, or needing follow-up.
4. The user creates or links the relevant company/customer and contacts.
5. A sales opportunity is created for the qualified business need.
6. The opportunity moves through the sales pipeline.
7. A quote is recorded for the opportunity.
8. A contract record is created or linked when the deal reaches contract stage.
9. Payment plans and actual payment records are tracked against the contract.
10. Activities, notes, and follow-up tasks preserve the sales history.
11. The opportunity is closed as won or lost with a recorded reason or outcome.
12. The team can review authorized historical records and continue customer follow-up.

## In Scope

- Team login and three-role access model: Administrator, Sales Manager, Sales.
- Lead creation, editing, qualification, assignment, owner changes, search, and filtering.
- ToB company/customer management.
- Contact management with multiple contacts per company/customer.
- Opportunity management with owner, amount, expected close date, pipeline stage, and close status.
- Sales pipeline progression from new opportunity through quote, contract, payment, won, or lost.
- Quote records linked to customer, opportunity, and owner.
- Contract records with status, amount, effective or signed date, related quote/opportunity/customer, and attachment or notes.
- Payment plans and actual payment records linked to contracts.
- Activity, note, and follow-up task records linked to CRM entities.
- Team collaboration history for authorized users.
- Core list/detail/search/filter experiences for CRM entities.
- Data persistence for all core CRM records.
- Production deployment readiness for the committed v1 scope.

## Out Of Scope Summary

- Contract approval workflow.
- Electronic signature.
- Contract template generation.
- Advanced reporting and forecasting.
- Email/calendar synchronization.
- External collaboration integrations such as Feishu, DingTalk, WeCom, ERP, or finance systems.
- AI sales recommendations or summaries.

Detailed scope exclusions are maintained in `docs/product/out-of-scope.md`.

## Success Metrics

| Metric | Target | Verification |
|---|---|---|
| P0 loop coverage | Every P0 step in the core business loop has an acceptance item. | Acceptance matrix review at G3. |
| P0 traceability | Every P0 requirement maps to at least one acceptance item before implementation planning. | Requirements and acceptance matrix review. |
| Data persistence | Core CRM records survive refresh, logout/login, and service restart in the target environment. | QA, integration, and audit evidence before Done. |
| Role correctness | Administrator, Sales Manager, and Sales can perform only authorized actions. | Permission tests and manual verification. |
| Production readiness | The v1 system can be deployed and operated with real data. | Deployment runbook, environment verification, and integration evidence. |

## Constraints

- Business: v1 must support team collaboration for ToB sales.
- Technical: frontend and backend must remain separated; shared contracts/types belong in `packages/shared/`.
- Time: no implementation work may begin before Gate G8 passes.
- Compliance: customer, contact, contract, and payment data require access control and audit-sensitive handling.
- Governance: P0/P1 items cannot be downgraded, deleted, weakened, or accepted as partial work.

## Initial Risks

| ID | Risk | Impact | Owner |
|---|---|---|---|
| RISK-001 | Payment and contract business rules may require more policy detail than currently confirmed. | P0 acceptance items may become blocked until rules are clarified. | Product Manager / Business Analyst |
| RISK-002 | Role visibility rules are confirmed only at a high level. | Permission design and verification may be blocked if manager/sales visibility boundaries are unclear. | Product Manager / Security Compliance |
| RISK-003 | Production deployment target is not yet defined. | Architecture and release readiness cannot be completed without environment decisions. | Architecture |
| RISK-004 | Reporting expectations may expand beyond core v1 loop. | Scope pressure may affect G2/G3 unless P1/P2 boundaries are preserved. | Product Manager |
