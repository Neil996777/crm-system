# Architecture Risk Register

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Risks

| ID | Severity | Risk | Impact | Mitigation | Owner |
|---|---|---|---|---|---|
| ARCH-RISK-001 | P0 | Service boundaries degrade into shared database or shared business logic. | System becomes distributed monolith; P0/P1 service governance fails. | Enforce service DB users, forbidden cross-service DB access, public contracts, and static/audit checks. | Architecture, Backend Engineer, Audit |
| ARCH-RISK-002 | P0 | Authorization is enforced only at frontend or gateway. | Unauthorized API access and invalid ACC-002 evidence. | Domain services must enforce permission and business rules server-side. QA must test direct API denial. | Architecture, Security Compliance, QA TDD |
| ARCH-RISK-003 | P0 | Audit/history event is lost after successful sensitive mutation. | ACC-014/ACC-022 cannot pass; business traceability breaks. | Use local transaction plus outbox-equivalent reliable publication for P0/P1 history/audit events. | Architecture, Backend Engineer |
| ARCH-RISK-004 | P0 | Self-hosted PostgreSQL backup is stored only on same ECS host. | Host or disk failure may lose database and backups together. | Encrypted daily local backups with 7-day retention are required baseline; production release is blocked until encrypted off-server backup evidence and restore rehearsal exist. | Infrastructure Ops, Architecture |
| ARCH-RISK-005 | P0 | Import/export bypasses domain services. | Invalid data, permission bypass, missing audit, and corrupted records. | Import/export service must call target service APIs and preserve row-level validation. | Architecture, Backend Engineer, QA TDD |
| ARCH-RISK-006 | P1 | Reporting directly queries source service tables. | Breaks service ownership and may leak unauthorized aggregates. | Reporting owns read models and uses events or approved Query APIs only. | Architecture, Backend Engineer |
| ARCH-RISK-007 | P1 | Multi-service deployment exceeds operational maturity on one ECS. | Debugging and recovery become slow. | Docker Compose health checks, logs with correlation ID, simple runbook, and service restart procedures. | Infrastructure Ops |
| ARCH-RISK-008 | P1 | Cross-service timeout/retry/idempotency not implemented consistently. | Duplicate records or failed business flows. | Define contract requirements in G5, model in PSM, test in G7/G10/G11. | Architecture, QA TDD, Integration Owner |
| ARCH-RISK-009 | P0 | Production HTTPS endpoint is not specified before release validation. | ACC-017 deployment evidence and secure login/session evidence may be incomplete. | Allow IP-based pre-release validation only; production release requires HTTPS endpoint, TLS evidence, security group evidence, health check, and smoke-test evidence. | Infrastructure Ops, Security Compliance |
| ARCH-RISK-010 | P1 | OQ-016 migration/seed data remains unresolved late. | Launch data readiness may slip. | Product Manager and Business Analyst must close before production launch planning. | Product Manager, Business Analyst |
| ARCH-RISK-011 | P0 | Service-to-service calls rely only on Docker internal network. | A compromised service could call another service as a bypass. | Require signed/HMAC service tokens, caller service ID, intent, audience, expiry, key rotation, rejection behavior, and audit logs. | Architecture, Security Compliance |
| ARCH-RISK-012 | P0 | Stale edits overwrite newer P0 record updates. | User data can be lost and ACC editable-record flows cannot pass. | Require `version` and `expectedVersion` on editable record DTOs/commands and `VERSION_CONFLICT` recovery. | Architecture, Backend Engineer, QA TDD |
| ARCH-RISK-013 | P0 | Archive hides records with active downstream obligations. | Open tasks, pending signatures, or unpaid payments become invisible or unmanaged. | Require archive eligibility API, active obligation DTO, blocked archive response, and history/operation log events. | Architecture, Business Analyst, UX/UI |
| ARCH-RISK-014 | P1 | CSV export enables spreadsheet formula injection. | Exported files may execute unsafe formulas when opened. | Reject or escape dangerous CSV cells and record export safety metadata. | Security Compliance, Backend Engineer |

## Blocker Interpretation

P0/P1 risks become blockers when they prevent architecture, MDA, task planning,
testing, integration, audit, or release evidence from satisfying the acceptance
matrix.

No risk may be resolved by weakening P0/P1 acceptance items.
