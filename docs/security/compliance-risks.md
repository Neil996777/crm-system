# Compliance Risks

## Document Control

- Project: CRM System
- Phase: G4 Security Design
- Owner Agent: Security Compliance
- Status: Accepted as Architecture Input

## Risk Register

| ID | Priority | Risk | Impact | Mitigation / Required Downstream Input | Owner Agent | Verification Gate | Status |
|---|---|---|---|---|---|---|---|
| COMP-001 | P0 | Authorization is implemented only in UI or route guards. | Unauthorized API access, IDOR, data leakage, and invalid acceptance evidence. | Architecture must define backend authz enforcement points for every protected action; QA must test direct API denial. | Architecture, Backend, QA | G5, G7, Audit | Input to Architecture |
| COMP-002 | P0 | Sales ownership/assignment relation is modeled too weakly to enforce Sales visibility. | Sales may see or mutate non-owned/non-assigned records. | Architecture and MDA must model owner, assignee, related-parent scope, and child-record visibility explicitly. | Architecture, Domain Modeling | G5, G6 | Input to Architecture |
| COMP-003 | P0 | User role/status changes do not invalidate or re-evaluate existing sessions. | Disabled or downgraded users may retain access. | Architecture must define session and role re-evaluation behavior; QA must test stale role/status. | Architecture, Backend, QA | G5, G7 | Input to Architecture |
| COMP-004 | P0 | Last Administrator protection is not enforced server-side. | System governance can be locked out. | Backend authorization must block disable/downgrade of the last active Administrator and log the blocked result. | Architecture, Backend, QA, Audit | G5, G7, Audit | Input to Architecture |
| COMP-005 | P0 | Audit events are not written atomically with sensitive business changes. | Business state changes may lack required history or operation-log evidence. | Architecture must define durable event write behavior and failure handling for sensitive operations. | Architecture, Backend, QA | G5, G7 | Input to Architecture |
| COMP-006 | P0 | Record-local history or global operation logs are editable through normal CRM workflows. | Auditability and collaboration history are unreliable. | Architecture must define append-only or tamper-evident behavior; QA/Audit must verify no edit/delete surface exists. | Architecture, QA, Audit | G5, Audit | Input to Architecture |
| COMP-007 | P0/P1 | Sensitive before/after values leak in errors, import row results, toasts, or log summaries. | Contact, contract, payment, or restricted business data may leak to unauthorized users. | Frontend/backend contracts must distinguish safe summaries from authorized detail payloads. | Architecture, Frontend, Backend, QA | G5, G7 | Input to Architecture |
| COMP-008 | P1 | CSV export can include unauthorized or archived records by default. | Bulk data leakage. | Export authorization must apply before query/output; archived inclusion requires explicit authorized filter; export run logged. | Architecture, Backend, QA, Audit | G5, G7, Audit | Input to Architecture |
| COMP-009 | P1 | CSV import allows formula injection or reference-based permission bypass. | Spreadsheet execution risk or unauthorized data mutation. | Import validation must handle dangerous cell prefixes, row permissions, related references, and partial failure behavior. | Architecture, Backend, QA | G5, G7 | Input to Architecture |
| COMP-010 | P1 | Reports aggregate unauthorized records before filtering. | Unauthorized users infer sensitive business totals. | Report queries must enforce authorization before aggregation and drill-in. | Architecture, Backend, QA | G5, G7 | Input to Architecture |
| COMP-011 | P0/P1 | Archive is treated like deletion or breaks related history. | P0 history, reports, reminders, and audit evidence become inconsistent. | Architecture/MDA must model archived state, active/default filters, explicit archived filters, and no-hard-delete behavior. | Architecture, Domain Modeling, QA | G5, G6, G7 | Input to Architecture |
| COMP-012 | P0 | Production backup and restore handling is undefined for restricted data. | Recovery may lose audit evidence or expose restricted data. | Architecture must define backup location, retention, access control, restore process, and relation to privacy retention. | Architecture, Integration, Audit | G5, Integration, Audit | Input to Architecture |
| COMP-013 | P0/P1 | Retention expectations are not represented in data design. | Contract/payment/log retention may be shortened accidentally. | Architecture and MDA must carry the committed retention policy into data design and lifecycle states. | Architecture, Domain Modeling, Audit | G5, G6, Audit | Input to Architecture |
| COMP-014 | P0 | Direct business-rule bypass through API is possible. | Invalid Won, overpayment, expired quote contracts, or forbidden transitions can corrupt core CRM state. | Backend must enforce business rules and authorization together; QA must test direct API negative paths. | Architecture, Backend, QA | G5, G7 | Input to Architecture |
| COMP-015 | P1 | Duplicate warning exposes unauthorized matched record details. | User can enumerate customers or contacts. | Duplicate warning UI/API must return safe match signals without restricted matched record details. | Architecture, Frontend, Backend, QA | G5, G7 | Input to Architecture |
| COMP-016 | P0/P1 | Audit and security tests are not traceable to ACC items. | P0/P1 items may be marked complete without security evidence. | QA/Test Model must map permission, abuse-case, privacy, and audit tests to ACC IDs and security requirement IDs. | QA TDD, Domain Modeling, Audit | G7, Audit | Input to QA |

## Downstream Validation Inputs

Architecture must explicitly address:

- Authentication and session strategy.
- Backend authorization enforcement points.
- Owner/assignee and related-parent scope modeling.
- Admin user/role lifecycle and last-Administrator protection.
- Audit event storage, append-only behavior, and query model.
- Import/export validation and storage handling.
- Report authorization before aggregation.
- Retention, archive, backup, and restore boundaries.

QA TDD must explicitly cover:

- Positive and negative permission tests for every actor/action/resource class.
- IDOR and direct API mutation tests.
- Login enumeration and stale-session/role tests.
- CSV import/export abuse tests.
- Audit event creation and global-log denial tests.
- Privacy masking and safe-summary tests.

Integration must explicitly prove:

- Frontend, backend, and persistence enforce the same role and record scope.
- Data survives refresh, logout/login, and service restart with permissions
  intact.
- History and operation logs are created for sensitive end-to-end flows.
- Exported and reported data match authorized persisted records only.

Audit must explicitly verify:

- No P0/P1 item is marked Done without implementation, QA, integration, and
  audit evidence.
- No core CRM path is satisfied by mock, static-only, in-memory-only, or
  non-persistent behavior.
- No P0/P1 security, privacy, permission, retention, or audit requirement is
  downgraded during architecture, modeling, planning, or implementation.
