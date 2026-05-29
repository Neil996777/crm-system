# Project Context: CRM System

## Project Summary

This project is a CRM system for managing leads, customers, contacts, sales
opportunities, quotes, contracts, payments, activities, and follow-ups.

## Current Phase

- Phase: Architecture Reset
- Current Gate: G5 Architecture Design Required
- Current Status: Previous architecture, MDA, task planning, implementation,
  QA, integration, audit, deployment, and project-specific implementation-agent
  artifacts were discarded by user direction on 2026-05-29.

## Current Goal

Restart architecture design from the retained product, business, UX/UI, and
security inputs. New architecture must not reuse discarded engineering
artifacts as design authority. After the new architecture passes G5, the
project must proceed through MDA, traceability/test model, task planning, and
later implementation gates again.

## Reset Boundary

Retained active inputs:

- Product requirements, PRD, acceptance matrix, open questions, out-of-scope,
  decision log, and G4 work plan.
- Business process, rule, scenario, permission, edge-case, and glossary
  documents.
- UX/UI flows, journeys, screen flows, interaction, screen state, UI,
  component, and responsive specifications.
- Security requirements, permission matrix, audit-log specification, privacy
  requirements, abuse cases, and compliance risks.
- Workspace delivery rules, no-downgrade rule, Gate rules, and the project
  sequencing note.

Discarded engineering artifacts:

- Architecture documents.
- MDA modeling documents.
- Task planning documents.
- Implementation source, generated contracts, tests, tooling, deployment,
  migration, dependency, and build files.
- QA, integration, and audit reports created from the discarded implementation
  and engineering design.
- Project-specific implementation agents created for the discarded G8 plan.

Archived historical review evidence may still exist under `archive/`, but it
is not current design authority for the restarted architecture work.

## Current Sequencing Rule

Product, business, UX/UI, and security design define what must be built.
Architecture then defines how the system will implement those requirements and
constraints. MDA Modeling turns product, business, UX/UI, security, and
accepted architecture outputs into traceable CIM/PIM/PSM artifacts. Coding
follows accepted MDA documents and G8 tasks.

## Business Domain

Core CRM concepts:

- Lead
- Customer
- Contact
- Company
- Opportunity
- Quote
- Contract
- Payment
- Deal Stage
- Activity
- Follow-up
- Sales Owner
- Pipeline
- Task
- Note

## Current Active Documents

Product:

- `docs/product/project-charter.md`
- `docs/product/requirements.md`
- `docs/product/prd.md`
- `docs/product/acceptance-matrix.md`
- `docs/product/open-questions.md`
- `docs/product/out-of-scope.md`
- `docs/product/decision-log.md`
- `docs/product/g4-work-plan.md`
- `development-sequencing-change-note.md`

Business:

- `docs/business/business-processes.md`
- `docs/business/business-rules.md`
- `docs/business/user-scenarios.md`
- `docs/business/role-permission-scenarios.md`
- `docs/business/edge-cases.md`
- `docs/business/business-glossary.md`

UX/UI:

- `docs/ux-ui/ux-flows.md`
- `docs/ux-ui/user-journeys.md`
- `docs/ux-ui/screen-flows.md`
- `docs/ux-ui/interaction-spec.md`
- `docs/ux-ui/screen-state-spec.md`
- `docs/ux-ui/ui-spec.md`
- `docs/ux-ui/component-spec.md`
- `docs/ux-ui/responsive-spec.md`

Security:

- `docs/security/security-requirements.md`
- `docs/security/permission-matrix.md`
- `docs/security/audit-log-spec.md`
- `docs/security/privacy-requirements.md`
- `docs/security/abuse-cases.md`
- `docs/security/compliance-risks.md`

Compliance:

- `docs/compliance/README.md`
- `docs/compliance/cmmi-process-standard.md`
- `docs/compliance/evidence-register-template.md`

Intellectual Property:

- `docs/ip/README.md`
- `docs/ip/software-copyright-standard.md`
- `docs/ip/patent-readiness-standard.md`

## Current Open Questions

- OQ-001: Exact production deployment target, domain, database, backup
  location, and environment ownership were reopened by the architecture reset
  and must be resolved by new Architecture Design.
- OQ-016: Data migration or initial seed data requirements remain launch
  planning inputs before production release.

## Current Blockers

- Architecture is not currently accepted. No MDA, task planning, or
  implementation work may proceed until new architecture design passes the
  required Gate.
