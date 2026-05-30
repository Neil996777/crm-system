# Service Boundary Security Inputs

## Document Control

- Project: CRM System
- Phase: G5 Pre-Architecture Input Supplement
- Owner Agent: Security Compliance
- Status: Ready for Architecture Intake
- Date: 2026-05-29
- Sources:
  - `docs/product/business-capability-map.md`
  - `docs/product/acceptance-matrix.md`
  - `docs/security/security-requirements.md`
  - `docs/security/permission-matrix.md`
  - `docs/security/privacy-requirements.md`
  - `docs/security/audit-log-spec.md`
  - `docs/security/abuse-cases.md`
  - `docs/security/compliance-risks.md`

## Purpose

This document provides security inputs for service-boundary-first architecture.
It does not decide final service boundaries, protocols, deployment units, or
implementation technology.

## Security Principles For Service Design

- Frontend hiding or disabling actions is UX guidance only.
- Every protected capability requires backend authorization.
- Internal service calls must not be trusted only because they are internal.
- Service-to-service calls must define caller identity, allowed action,
  resource, scope, denial behavior, and audit behavior.
- Sensitive payloads must be minimized and classified.
- Cross-service flows must carry correlation IDs for integration and audit
  evidence.
- Service boundaries must not weaken record scope, data masking, retention, or
  audit requirements.

## Trust Boundary Inputs

| Boundary ID | Boundary | Security Requirement | Acceptance IDs |
|---|---|---|---|
| STB-001 | User browser to CRM backend/API surface | Authenticate protected routes/actions; reject unauthenticated and disabled users safely. | ACC-001, ACC-002 |
| STB-002 | Frontend state to backend authorization | Frontend role/action hints are not authoritative; backend evaluates actor/action/resource/condition. | ACC-002 |
| STB-003 | Service candidate to service candidate | Internal calls require service identity, authorization, safe payload, and correlation ID. | ACC-002, ACC-014, ACC-022 |
| STB-004 | Service candidate to persistence ownership | A service must not directly read/write another service's owned data without approved contract/read model. | ACC-016 |
| STB-005 | Import/export boundary | Uploaded/imported and generated/exported CSV content is Restricted; row errors and summaries must be safe. | ACC-020, ACC-022 |
| STB-006 | Reporting boundary | Aggregates must exclude unauthorized data and default archived records. | ACC-018, ACC-023 |
| STB-007 | History/log boundary | Record-local history follows record permission; global operation log is Administrator-only. | ACC-014, ACC-022 |
| STB-008 | Infrastructure and runtime boundary | Secrets, backup credentials, deployment configs, and operational access are not exposed in product UI or public repo. | ACC-017 |

## Sensitive Cross-Service Data Flow Inputs

| Flow ID | Data Flow | Data Classification | Minimum Security Requirements | Acceptance IDs |
|---|---|---|---|---|
| SDF-001 | Identity/session/role context to protected capabilities | Security Critical | Active user recheck, role/scope evaluation, safe denial, correlation ID. | ACC-001, ACC-002 |
| SDF-002 | Lead/customer/contact context to opportunity and commercial work | Confidential | Scope-limited lookup, no unauthorized existence leak, history event on mutation. | ACC-003 to ACC-010, ACC-014 |
| SDF-003 | Quote/contract/payment amount and status | Restricted | Least privilege, exact authorization, safe errors, audit events on sensitive mutations. | ACC-009 to ACC-014, ACC-022 |
| SDF-004 | Reminder source data across tasks/contracts/payments | Confidential / Restricted | Authorized reminder query, inactive record exclusion, no unauthorized related detail. | ACC-021 |
| SDF-005 | Duplicate warning match signal | Confidential | Warning may be shown but must not expose unauthorized matched record details; no automatic merge. | ACC-019 |
| SDF-006 | Import row content and validation result | Restricted | Row-level safe errors; valid rows mutate only through normal authorization and validation paths. | ACC-020 |
| SDF-007 | Export result content | Restricted | Explicit confirmation, authorized scope only, operation log, safe metadata. | ACC-020, ACC-022 |
| SDF-008 | Report aggregate inputs | Confidential / Restricted | Authorization before aggregation; unauthorized rows excluded from sums/counts. | ACC-018, ACC-023 |
| SDF-009 | Record-local history and operation logs | Restricted / Security Critical | Append-only, role/scope-controlled queries, safe before/after values. | ACC-014, ACC-022 |
| SDF-010 | Backup/restore and deployment evidence | Security Critical / Operational | Secrets not committed; restore evidence must preserve CRM data, history, and logs. | ACC-016, ACC-017, ACC-022 |

## Service-To-Service Permission Input Matrix

The final matrix belongs to Architecture and PSM. These inputs must be preserved:

| Caller Candidate | Target Candidate | Allowed Purpose | Required Controls |
|---|---|---|---|
| SVC-CAND-QUERY-EXPERIENCE | Record-owning candidates | Authorized list/detail/search/filter views | Actor context, scope check, archived filter rule, safe error. |
| SVC-CAND-LEAD | SVC-CAND-ACCOUNT / SVC-CAND-OPPORTUNITY | Lead conversion to customer/contact/opportunity context | Conversion-once guard, permission check, transactional or reliable consistency strategy, history event. |
| SVC-CAND-COMMERCIAL | SVC-CAND-OPPORTUNITY | Payment status and Won eligibility | Full-payment verification, no direct unauthorized data access, correlation ID. |
| SVC-CAND-WORK | Record-owning candidates | Reminder related-record display and eligibility | Scope check, inactive record exclusion, safe related record summary. |
| SVC-CAND-REPORTING | Record-owning candidates | Authorized metrics and overview | Authorization before aggregation, no unauthorized rows or aggregates. |
| SVC-CAND-IMPORT-EXPORT | Record-owning candidates | Import/export rows through normal business rules | Authorization per scope, row validation, safe row errors, operation log. |
| Mutation-producing candidates | SVC-CAND-HISTORY-AUDIT | Record-local history and operation log event creation | Append-only event, actor/resource/action/result, correlation ID, failure behavior. |
| SVC-CAND-PLATFORM-OPS | Runtime service candidates | Health, backup, restore, deployment evidence | No secret exposure, auditable environment ownership, restore verification. |

## Security Blockers For G5/G8

| Blocker ID | Severity | Issue | Required Closure |
|---|---|---|---|
| SEC-SVC-BLK-001 | P0 | Final service trust boundaries missing. | Architecture must define boundaries and Security must review. |
| SEC-SVC-BLK-002 | P0 | Service-to-service authorization and audit behavior missing for P0/P1 flows. | Architecture and Security must define contracts/rules before G8. |
| SEC-SVC-BLK-003 | P0 | Data ownership and forbidden direct access rules missing. | Architecture and Domain Modeling must define and represent in PSM. |
| SEC-SVC-BLK-004 | P0 | Error and denial contracts missing for safe UI/API behavior. | Architecture must define safe error contracts; QA must test. |
| SEC-SVC-BLK-005 | P1 | Correlation ID expectations missing for service-chain integration evidence. | Architecture and Integration Owner must define before G8. |
| SEC-SVC-BLK-006 | P1 | Infrastructure secrets, backup, restore, and public exposure ownership missing. | Infrastructure Ops and Architecture must close OQ-001. |

