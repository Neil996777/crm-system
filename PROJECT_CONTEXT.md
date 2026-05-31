# Project Context: CRM System

## Project Summary

This project is a CRM system for managing leads, customers, contacts, sales
opportunities, quotes, contracts, payments, activities, and follow-ups.

## Current Phase

- Phase: G5 Architecture Design Passed
- Current Gate: G6 MDA Modeling Required
- Current Status: Previous architecture, MDA, task planning, implementation,
  QA, integration, audit, deployment, and project-specific implementation-agent
  artifacts were discarded by user direction on 2026-05-29. Product,
  Business Analyst, UX/UI, and Security added service-boundary-first input
  supplements. Architecture produced a new G5 design, repaired it through two
  review rounds, and G5 passed on 2026-05-30 with all required reviewers
  (including project-added Infrastructure Ops) returning Pass. See
  `archive/reviews/g5-architecture/g5-architecture-final-decision-2026-05-30.md`.

## Current Goal

G5 Architecture Design has passed. The project proceeds to G6 MDA Modeling:
Domain Modeling and Architecture build CIM/PIM/PSM, service mapping, state
machines, domain events, and traceability that faithfully represent the accepted
architecture. Per project interpretation (GAP-PROC-003), G6 does not re-litigate
accepted architecture decisions. Implementation remains blocked until G8.

The current immediate goal is the G6 MDA package, reviewed by Product Manager,
Business Analyst, UX Designer, UI Designer, Security Compliance, and QA Test
Design. The
accepted architecture defines the service list, owner agents, contracts, data
ownership, and deployment decisions that MDA must trace, not redefine.

Production deployment target (accepted at G5): runtime host
`srv-volcengine-sh-01` (Volcengine ECS, Shanghai, 4 vCPU / 8 GiB); off-server
backup target `srv-aliyun-bj-01` (Alibaba Cloud ECS, Beijing). Off-server backup
copy evidence, restore rehearsal, and HTTPS/TLS endpoint evidence remain release
blockers, correctly deferred to release.

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
- `docs/product/business-capability-map.md`
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
- `docs/business/service-governance-inputs.md`

UX/UI:

- `docs/ux-ui/ux-flows.md`
- `docs/ux-ui/user-journeys.md`
- `docs/ux-ui/screen-flows.md`
- `docs/ux-ui/interaction-spec.md`
- `docs/ux-ui/screen-state-spec.md`
- `docs/ux-ui/ui-spec.md`
- `docs/ux-ui/component-spec.md`
- `docs/ux-ui/responsive-spec.md`
- `docs/ux-ui/service-state-mapping.md`

Security:

- `docs/security/security-requirements.md`
- `docs/security/permission-matrix.md`
- `docs/security/audit-log-spec.md`
- `docs/security/privacy-requirements.md`
- `docs/security/abuse-cases.md`
- `docs/security/compliance-risks.md`
- `docs/security/service-boundary-security.md`

Architecture:

- `docs/architecture/architecture.md`
- `docs/architecture/module-boundaries.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`
- `docs/architecture/service-architecture-adr.md`
- `docs/architecture/service-architecture-acceptance.md`
- `docs/architecture/service-acceptance-map.md`

## Current Open Questions

- OQ-001: Architecture decision recorded; production release evidence remains
  pending. The draft selects Alibaba Cloud ECS, Docker Compose, multiple Go
  service containers, self-hosted PostgreSQL, encrypted local automatic backup
  with 7-day retention, HTTPS-only production ingress, and Architecture plus
  Infrastructure Ops ownership. Pre-release validation may use the ECS IP.
  Production release requires HTTPS endpoint/TLS evidence, security group
  evidence, monitoring evidence, restore rehearsal, and encrypted off-server
  backup evidence. Same-host-only backup is a release blocker, not accepted
  P0/P1 completion.
- OQ-016: Data migration or initial seed data requirements remain launch
  planning inputs before production release.

## Current Blockers

- No open G5 architecture blocker. G5 passed on 2026-05-30.
- G6 MDA Modeling has not yet started; task planning and implementation remain
  gated until G6, G7, and G8 pass.
- Carried-forward release blockers (not gate blockers now): encrypted off-server
  backup copy + restore rehearsal evidence, HTTPS/TLS endpoint evidence,
  security-group and monitoring evidence. These block production release, not G6.
- Company-layer follow-up (deferred to workspace discussion): infrastructure
  registers should name CRM as the consuming project for the runtime host. See
  `process/process-gap-register.md`.
