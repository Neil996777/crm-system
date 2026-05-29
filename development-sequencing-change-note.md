# Development Sequencing Change Note

## Purpose

This note records the active development-process sequencing decision for the CRM
project and the workspace-level process.

## Current Sequencing Decision

Product, business, UX/UI, and security design define what must be built.

Architecture comes after these design inputs and defines how the system will
implement them.

MDA Modeling comes after Architecture and turns product, business, UX/UI,
security, and architecture outputs into traceable CIM/PIM/PSM artifacts.

Task Planning starts only after product acceptance, accepted MDA, traceability,
and test model exist.

Implementation starts only after G8 passes.

## Correct Sequence

1. Product Manager creates PRD and Product Acceptance Matrix.
2. Business Analyst creates business processes, business rules, scenarios,
   permissions, edge cases, and glossary.
3. UX Designer creates user journeys, flows, interaction specs, and screen-state
   specs.
4. UI Designer creates UI spec, component spec, responsive spec, and visual
   state specs.
5. Security Compliance creates security requirements, permission matrix,
   privacy requirements, audit-log spec, abuse cases, and compliance risks.
6. Architecture creates the technical architecture that implements all accepted
   product, business, UX/UI, and security requirements.
7. Domain Modeling creates CIM, PIM, PSM, domain model, state machines, domain
   events, traceability matrix, and test model.
8. Task Planner creates tasks from product acceptance, accepted MDA,
   traceability, and test model.
9. Frontend and Backend implement only after G8 passes.
10. QA, Integration, and Audit verify completion.

## Architecture Acceptance

Architecture acceptance is represented inside `modeling/PSM.md` by default.

A separate architecture acceptance matrix is optional only for complex projects.

Architecture acceptance must map to:

- PRD IDs
- product acceptance IDs
- NFR IDs where relevant
- business rules
- UX/UI requirements
- security requirements
- PSM elements
- tasks
- tests

It must not:

- create new product scope
- replace `docs/product/acceptance-matrix.md`
- weaken or reinterpret P0/P1 product acceptance
- allow mock, static-only, TODO, in-memory-only, or non-persistent behavior to
  satisfy core CRM paths

## Current Boundary

Implementation remains blocked until G8 passes.

No P0/P1 item may be downgraded, deleted, merged away, weakened, or accepted as
partial work.
