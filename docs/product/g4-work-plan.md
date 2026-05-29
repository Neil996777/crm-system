# G4 Work Plan

## Purpose

G4 moves the project from stable product acceptance into business, UX/UI, and
security design.

Current sequencing rule:

- Product acceptance defines what must be completed.
- Business, UX/UI, and Security define the detailed design and constraints.
- Architecture comes after these design inputs and defines how to implement them.
- MDA Modeling comes after Architecture and turns the accepted design and
  architecture into CIM/PIM/PSM.
- Task Planning and implementation start only after G8 passes.

## Gate Context

- Current Gate: G5 Architecture Design Required; G4 completed and preserved as design input
- Transition: Acceptance Matrix -> Business/UX/UI/Security Design
- Owner: Product Manager
- Required Reviewers: Business Analyst, UX Designer, UI Designer, Security Compliance
- Pass Condition: Product acceptance is stable enough for detailed design.

## Sequencing

| Step | Workstream | Output | Sequencing Rule |
|---|---|---|---|
| 1 | Business Design | Business processes, business rules, user scenarios, role-permission scenarios, edge cases, glossary | Defines how the CRM business operates. Does not define technical implementation. |
| 2 | UX Design | User journeys, UX flows, screen flows, interaction specs, screen-state specs | Defines user tasks, interaction behavior, screen states, feedback, and recovery. Does not define API/database/service structure. |
| 3 | UI Design | UI spec, component spec, responsive spec, visual states | Defines visual structure, components, states, and responsive behavior. Does not define frontend architecture. |
| 4 | Security Design | Security requirements, permission matrix, audit-log spec, privacy requirements, abuse cases, compliance risks | Defines permission, privacy, audit, abuse-case, and compliance constraints. Does not define UX/UI or architecture. |
| 5 | Architecture Design | Architecture, module boundaries, API spec, data design, authz architecture, integration design, frontend/backend contract | Defines how to implement accepted product, business, UX/UI, and security requirements without downgrading them. |
| 6 | MDA Modeling | CIM, PIM, PSM, domain model, state machines, domain events | Turns accepted design and architecture into traceable engineering models. PSM represents platform-specific architecture. |
| 7 | Traceability + Test Model | Traceability matrix, test model, test plan | Proves P0/P1 acceptance items are traceable and testable before task planning. |
| 8 | Task Planning | Tasks, dependencies, delivery plan, acceptance-task map, blockers | Creates end-to-end implementation tasks after accepted MDA and test model exist. |

## No-Downgrade Rule

- P0/P1 acceptance items from `docs/product/acceptance-matrix.md` remain the
  source of truth.
- Downstream artifacts may clarify, decompose, or strengthen P0/P1 behavior.
- Downstream artifacts must not downgrade, delete, merge away, weaken, or
  accept partial P0/P1 behavior.
- No mock, static-only, TODO, in-memory-only, or non-persistent behavior can
  satisfy core CRM paths.

## Architecture Acceptance

Architecture acceptance is represented inside `modeling/PSM.md` by default.

It verifies that architecture decisions satisfy:

- product acceptance IDs
- NFRs
- business rules
- UX/UI requirements
- security requirements
- PSM elements
- tasks and tests

It must not:

- create new product scope
- replace `docs/product/acceptance-matrix.md`
- weaken or reinterpret P0/P1 product acceptance
- allow mock, static-only, TODO, in-memory-only, or non-persistent behavior to satisfy core paths

A separate architecture acceptance matrix is optional only for complex projects.

## Immediate Next Work

1. Business Analyst prepares business processes, business rules, user scenarios,
   role-permission scenarios, edge cases, and glossary for the complete ToB CRM
   loop.
2. UX Designer prepares user journeys, UX flows, screen flows, interaction specs,
   and screen-state specs for Administrator, Sales Manager, and Sales.
3. UI Designer prepares UI spec, component spec, responsive spec, and visual
   state requirements after UX flows are stable.
4. Security Compliance prepares permission, audit, privacy, abuse-case, and
   compliance documents.
5. Architecture designs module boundaries, API/data design, authz architecture,
   integration design, and frontend/backend contract from the accepted design
   inputs.
6. Domain Modeling creates CIM/PIM/PSM, domain model, state machines, domain
   events, traceability matrix, and test model after Architecture.
7. Task Planning and coding start only after G8; coder work follows accepted MDA
   and acceptance-linked tasks.

## Implementation Boundary

Implementation is blocked until the restarted delivery flow reaches and passes
G8 again.
