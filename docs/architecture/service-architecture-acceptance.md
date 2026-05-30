# Service Architecture Acceptance

This matrix verifies service architecture governance. It is subordinate to the
product acceptance matrix and may not create, downgrade, delete, weaken, or
reinterpret P0/P1 product acceptance.

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Service Architecture Acceptance

| ID | Priority | Acceptance Item | Completion Standard | Verification Method | Evidence | Owner Agent | Status | Gate |
|---|---|---|---|---|---|---|---|---|
| SVC-ACC-001 | P0 | Project uses service-boundary-first governance. | Business capability map, final service list, and service architecture ADR exist. | Document review | `docs/product/business-capability-map.md`, `docs/architecture/architecture.md`, `docs/architecture/service-architecture-adr.md` | Architecture | Ready for G5 Re-review | G5 |
| SVC-ACC-002 | P0 | Every service has a `Service Owner Agent`. | Service list contains exactly one owner agent per service. | Service list review | `docs/architecture/architecture.md` | Architecture | Ready for G5 Re-review | G5/G8 |
| SVC-ACC-003 | P0 | P0/P1 acceptance maps to services. | Every P0/P1 item maps to one or more services or service candidates. | Traceability review | `docs/product/acceptance-matrix.md`, `docs/architecture/architecture.md` | Domain Modeling | Pending G6/G7 | G7/G8 |
| SVC-ACC-004 | P0 | Service data ownership is explicit. | Each service declares owned data and forbidden direct data access. | Architecture review and audit | `docs/architecture/data-design.md` | Architecture | Ready for G5 Re-review | G6/G12 |
| SVC-ACC-005 | P0 | Contracts exist before implementation. | API/event/error/permission contracts cover P0/P1 service-backed acceptance and are referenced by implementation tasks. | Contract and task review | `docs/architecture/api-spec.md`, `docs/architecture/integration-design.md`, `docs/architecture/authz-architecture.md` | Architecture | Revised for G5 re-review, requires PSM/task refinement | G8 |
| SVC-ACC-006 | P0 | Cross-service internal dependency is prohibited. | No cross-service internal imports or shared business implementation. | Static check and audit | `docs/architecture/module-boundaries.md` | Backend Engineer / Frontend Engineer | Draft, requires implementation evidence | G9/G12 |
| SVC-ACC-007 | P0 | Cross-service database access is prohibited. | No service directly reads or writes another service's database tables. | Code and data-access audit | `docs/architecture/data-design.md` | Backend Engineer | Draft, requires implementation evidence | G9/G12 |
| SVC-ACC-008 | P0 | Service-to-service calls enforce security boundaries. | Authentication, authorization, audit, credential storage, rotation, rejection, and sensitive data rules exist and are tested. | Security review and tests | `docs/architecture/authz-architecture.md`, `docs/security/service-boundary-security.md` | Security Compliance | Revised for G5 re-review, requires tests later | G10/G12 |
| SVC-ACC-009 | P0 | Cross-service reliability is designed and tested. | Idempotency, retry, timeout, compensation, correlation ID, outbox/inbox or equivalent, owner-transfer recovery, duplicate warning, archive blocking, and terminal lifecycle behavior exist where needed. | Integration tests | `docs/architecture/integration-design.md` | Integration Owner | Revised for G5 re-review, requires test/integration evidence | G11 |
| SVC-ACC-010 | P0 | AI tasks are constrained by service boundaries. | Tasks include service, owner agent, contract, acceptance ID, and forbidden boundary access. | Task review | Pending G8 task planning | Task Planner | Pending G8 | G8 |
| SVC-ACC-011 | P0 | Fake core implementation cannot pass. | No mock, stub, TODO, fake data, static-only, or non-persistent path satisfies P0/P1. | QA and audit | `docs/architecture/architecture.md`, `docs/product/acceptance-matrix.md` | QA TDD / Audit | Draft, requires implementation evidence | G10/G12 |
| SVC-ACC-012 | P0 | End-to-end evidence exists for P0/P1 cross-service flows. | Integration evidence includes acceptance ID, environment, steps, actual result, service chain, correlation ID, deployment endpoint, monitoring, and backup/restore evidence where applicable. | Integration report and audit | `docs/architecture/integration-design.md`, `docs/architecture/deployment-notes.md` | Integration Owner / Audit | Pending G11/G12 | G11/G12 |

## Traceability Summary

| Service Acceptance ID | Product Acceptance IDs | Services | Contracts | Downstream Requirement |
|---|---|---|---|---|
| SVC-ACC-001 | ACC-001 to ACC-023 | SVC-001 to SVC-010 | ADR-ARCH-001 to ADR-ARCH-005 | PSM service mapping |
| SVC-ACC-002 | ACC-001 to ACC-023 | SVC-001 to SVC-010 | Service list | G8 tasks must preserve owner agent |
| SVC-ACC-003 | ACC-001 to ACC-023 | SVC-001 to SVC-010 | API/event contracts | G6/G7 traceability matrix |
| SVC-ACC-004 | ACC-016, ACC-017 | SVC-002 to SVC-010 | data ownership map | PSM data ownership and audit checks |
| SVC-ACC-005 | ACC-001 to ACC-023 | All service-backed items | Command, Query, Event, Error, Permission contracts | Contract tests before implementation completion |
| SVC-ACC-008 | ACC-001, ACC-002, ACC-014, ACC-020, ACC-022, ACC-023 | All protected services | Permission, service auth, session invalidation, and denial contracts | Security tests and audit |
| SVC-ACC-009 | ACC-004, ACC-005, ACC-007, ACC-011, ACC-013, ACC-014, ACC-019, ACC-020, ACC-021, ACC-022 | Cross-service flows | Idempotency, timeout, retry, owner transfer, archive obligation, duplicate warning, report/reminder/import-export event contracts | Integration evidence |
| SVC-ACC-012 | ACC-001 to ACC-023 | All flow services | Correlation ID, endpoint, monitoring, backup/restore, and integration contracts | G11 integration report |
