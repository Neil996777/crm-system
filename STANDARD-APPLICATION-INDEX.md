# Standard Application Index

Date: 2026-05-29
Status: Required before new architecture approval

## Purpose

This index tells CRM project agents how to apply the current company standards
after the architecture reset.

It is a routing and review index only. It does not define CRM service
boundaries, service owner names, contracts, PSM elements, tasks, tests, or
release decisions. Those conclusions must be produced by the responsible CRM
project agents from the active CRM inputs listed below.

## Required Company Standards

Before reviewing or changing CRM project artifacts, every project agent must
read:

- `../../AGENTS.md`
- `../../company/operating-model.md`
- `../../standards/acceptance-matrix-standard.md`
- `../../standards/status-and-priority-standard.md`
- `../../workflows/project-initialization.md`
- `../../workflows/software-delivery.md`
- `../../company/policy-changes/2026-05-29-default-microservice-governance.md`
- `../../templates/architecture.md`
- `../../templates/modeling/PSM.md`
- `../../templates/tasks.md`
- `../../templates/service-architecture-acceptance-matrix.md`

## Current Project Authority

Project agents must treat these files as the active CRM input set:

- `AGENTS.md`
- `PROJECT_CONTEXT.md`
- `README.md`
- `development-sequencing-change-note.md`
- `docs/product/project-charter.md`
- `docs/product/requirements.md`
- `docs/product/prd.md`
- `docs/product/acceptance-matrix.md`
- `docs/product/open-questions.md`
- `docs/product/out-of-scope.md`
- `docs/product/decision-log.md`
- `docs/product/g4-work-plan.md`
- `docs/business/business-processes.md`
- `docs/business/business-rules.md`
- `docs/business/user-scenarios.md`
- `docs/business/role-permission-scenarios.md`
- `docs/business/edge-cases.md`
- `docs/business/business-glossary.md`
- `docs/ux-ui/ux-flows.md`
- `docs/ux-ui/user-journeys.md`
- `docs/ux-ui/screen-flows.md`
- `docs/ux-ui/interaction-spec.md`
- `docs/ux-ui/screen-state-spec.md`
- `docs/ux-ui/ui-spec.md`
- `docs/ux-ui/component-spec.md`
- `docs/ux-ui/responsive-spec.md`
- `docs/security/security-requirements.md`
- `docs/security/permission-matrix.md`
- `docs/security/audit-log-spec.md`
- `docs/security/privacy-requirements.md`
- `docs/security/abuse-cases.md`
- `docs/security/compliance-risks.md`

Historical files under `archive/` are reference material only. They are not
current architecture, MDA, task, test, integration, audit, or release authority.

## Current Gate Position

- Current phase: Architecture Reset
- Current gate: G5 Architecture Design Required
- Implementation is blocked until G8 passes.
- New MDA, task planning, implementation, QA, integration, and audit artifacts
  must be rebuilt from the active input set.

## What Agents Must Produce

| Area | Required Output | Responsible Agent |
|---|---|---|
| Business capability map | Capability list mapped to P0/P1 acceptance items | product-manager, business-analyst |
| Product acceptance update | `docs/product/acceptance-matrix.md` with business capability, related services, service owner agents, and contracts for P0/P1 items | product-manager |
| Business rule review | Updated business rules, edge cases, permissions, and operational constraints needed by service design | business-analyst |
| UX flow review | UX paths mapped to capabilities and service-backed states where applicable | ux-designer |
| UI state review | UI components, empty/loading/error/permission states mapped to accepted flows | ui-designer |
| Security review | Permission, trust boundary, audit, privacy, and abuse-case constraints for service-backed capabilities | security-compliance |
| Service architecture | `architecture.md`, service list, service owner agents, data ownership, contracts, deployment notes, risk register | architecture |
| MDA/PSM rebuild | PSM service mapping, contracts, data ownership, events, permissions, traceability | domain-modeling |
| Task plan | `tasks.md` with service, service owner agent, contract reference, acceptance ID, tests, and forbidden boundary access | task-planner |
| Test model | Service contract tests, boundary tests, P0/P1 acceptance tests, integration scenarios | qa-tdd |
| Integration plan | Service-chain evidence plan, environment checks, correlation ID expectations | integration-owner |
| Infrastructure review | Deployment environment request, server/database/domain/port requirements, backup and monitoring expectations | infrastructure-ops |
| Final audit plan | Reverse trace from acceptance to service, model, task, test, integration, and infrastructure evidence | audit |

## Agent Review Checklist

### product-manager

- Verify every P0/P1 CRM acceptance item remains intact.
- Add or request the required service governance fields for P0/P1 items.
- Do not decide technical service boundaries alone.
- Mark unresolved P0/P1 mapping as `Blocked`, not `Done`.

### business-analyst

- Review business processes, rules, roles, exceptions, and operational
  constraints.
- Identify business capabilities and bounded-context candidates.
- Escalate contradictions between business rules and current acceptance items.

### ux-designer

- Review user journeys, UX flows, screen flows, and interaction behavior.
- Identify cross-capability user paths that may require service-chain evidence.
- Preserve required empty, error, permission, and recovery states.

### ui-designer

- Review UI specification, components, responsive behavior, and visual states.
- Identify UI states that need API, permission, loading, empty, or error
  contract support.

### security-compliance

- Review permissions, privacy, audit logs, abuse cases, and compliance risks.
- Define service-to-service trust boundaries and sensitive data handling
  constraints.
- Treat public exposure, privileged access, secrets, and audit gaps as blockers
  when they affect P0/P1 scope.

### architecture

- Produce the CRM service-boundary-first architecture from active inputs.
- Assign exactly one `Service Owner Agent` to every service candidate.
- Define contracts, data ownership, deployment boundaries, observability, and
  reliability behavior.
- If physical microservice separation is delayed, create an ADR without
  removing service boundaries, owners, contracts, data ownership, or tests.

### domain-modeling

- Rebuild MDA/PSM after architecture is accepted.
- Represent CRM service mapping, aggregate ownership, events, permissions,
  contracts, and traceability.
- Do not reuse discarded MDA artifacts as current authority.

### task-planner

- Create tasks only after required architecture and MDA inputs exist.
- Every service-backed task must include service, service owner agent,
  acceptance ID, contract reference, tests, and forbidden boundary access.

### frontend-engineer

- Do not start implementation before G8.
- After G8, implement only against accepted contracts, UI specs, and task
  boundaries.

### backend-engineer

- Do not start implementation before G8.
- After G8, implement only against accepted service ownership, data ownership,
  API/event/error/permission contracts, and tasks.

### qa-tdd

- Define tests from P0/P1 acceptance, service contracts, permissions,
  persistence, and boundary rules.
- Treat missing testability for P0/P1 service-backed capabilities as blocked.

### integration-owner

- Plan end-to-end service-chain evidence for P0/P1 flows.
- Verify environment readiness only after infrastructure requirements are
  documented.

### infrastructure-ops

- Do not decide CRM business architecture.
- Review only deployment environment needs: servers, database platform or
  self-hosting choice, domains, ports, ingress, secrets metadata, backups,
  monitoring, and operational ownership.
- Record infrastructure blockers separately from product acceptance blockers.

### audit

- Reverse-check that CRM acceptance items trace to capability, service, owner,
  contract, model, task, test, integration evidence, and infrastructure evidence
  where applicable.
- Block release when P0/P1 evidence is missing or service governance is
  incomplete.

## Required Gap Closure Before Implementation

G8 may not pass until the CRM project has:

| Item | Priority |
|---|---|
| Business capability map for P0/P1 acceptance items | P0 |
| Related service or service candidate for every software-backed P0/P1 item | P0 |
| Exactly one `Service Owner Agent` per service | P0 |
| API/event/error/permission contracts needed for P0/P1 implementation | P0 |
| Data ownership and forbidden cross-service access rules | P0 |
| Service-aware PSM and traceability | P0 |
| Task plan with service, owner, contract, acceptance ID, tests, and boundaries | P0 |
| QA test model for service contracts and P0/P1 acceptance | P0 |
| Integration evidence plan for service-backed P0/P1 flows | P1 |
| Infrastructure environment requirements and ownership record | P1 |

## Explicit Non-Authority

This index does not:

- approve architecture
- approve implementation
- approve release
- choose the CRM service split
- choose the CRM database schema
- choose deployment servers or ports
- create service contracts
- mark any acceptance item done
- override project gates or company standards
