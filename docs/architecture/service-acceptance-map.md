# Service Acceptance Map

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Purpose

This document maps product acceptance items to the G5 architecture service list.
It does not replace `docs/product/acceptance-matrix.md` as the product
completion source of truth. It provides architecture authority for G6 PSM,
G7 traceability/test model, and G8 task planning.

## Acceptance To Service Map

| Acceptance ID | Priority | Primary Services | Supporting Services | Service Owner Agent | Required Contract Families |
|---|---|---|---|---|---|
| ACC-001 | P0 | SVC-002 identity-authz-service | SVC-001 gateway-bff, SVC-008 audit-history-service | backend-engineer | Authn, opaque session, secure cookie, session revocation, active-user check, safe denial, audit event |
| ACC-002 | P0 | SVC-002 identity-authz-service | SVC-001 gateway-bff, all domain services, SVC-008 audit-history-service | backend-engineer | Permission, scope, denial error, role/status recheck, service-to-service authz, operation log |
| ACC-003 | P0 | SVC-003 lead-service | SVC-002 identity-authz-service, SVC-008 audit-history-service | backend-engineer | Lead command, lead query, permission, history event |
| ACC-004 | P0 | SVC-003 lead-service | SVC-004 account-service, SVC-005 opportunity-service, SVC-008 audit-history-service | backend-engineer | Lead transition, lead conversion, account/contact command, opportunity command, history event |
| ACC-005 | P0 | SVC-004 account-service | SVC-002 identity-authz-service, SVC-007 work-service, SVC-006 commercial-service, SVC-008 audit-history-service | backend-engineer | Account command/query, permission, archive eligibility, owner transfer, expectedVersion, history event |
| ACC-006 | P0 | SVC-004 account-service | SVC-002 identity-authz-service, SVC-008 audit-history-service | backend-engineer | Contact command/query, permission, expectedVersion, history event |
| ACC-007 | P0 | SVC-005 opportunity-service | SVC-004 account-service, SVC-002 identity-authz-service, SVC-007 work-service, SVC-008 audit-history-service | backend-engineer | Opportunity command/query, account summary, archive eligibility, owner transfer, expectedVersion, history event |
| ACC-008 | P0 | SVC-005 opportunity-service | SVC-006 commercial-service, SVC-008 audit-history-service | backend-engineer | Opportunity transition, commercial status summary, history event |
| ACC-009 | P0 | SVC-006 commercial-service | SVC-005 opportunity-service, SVC-008 audit-history-service | backend-engineer | Quote command, quote transition, opportunity summary, history event |
| ACC-010 | P0 | SVC-006 commercial-service | SVC-005 opportunity-service, SVC-007 work-service, SVC-008 audit-history-service | backend-engineer | Contract command, contract transition, accepted quote constraint, pending signature obligation, expectedVersion, history event |
| ACC-011 | P0 | SVC-006 commercial-service | SVC-008 audit-history-service | backend-engineer | Payment command, payment status, idempotency, history event |
| ACC-012 | P0 | SVC-007 work-service | Target record-owning services, SVC-002 identity-authz-service, SVC-008 audit-history-service | backend-engineer | Work item command, related-record summary, permission, history event |
| ACC-013 | P0 | SVC-005 opportunity-service | SVC-006 commercial-service, SVC-007 work-service, SVC-008 audit-history-service | backend-engineer | Close Won, Close Lost, required lostReason, terminal edit protection, post-close work path, payment status summary, history event |
| ACC-014 | P0 | SVC-008 audit-history-service | All mutation-producing services, SVC-002 identity-authz-service | backend-engineer | Append-only history event, history query, tamper-evidence hash, permission |
| ACC-015 | P0 | SVC-001 gateway-bff | Record-owning services, SVC-002 identity-authz-service | backend-engineer | List query, detail query, filter error, permission |
| ACC-016 | P0 | SVC-002 to SVC-010 | PostgreSQL runtime, backup job | backend-engineer | Persistence, health, encrypted local backup, off-server production backup evidence, restore rehearsal |
| ACC-017 | P0 | Runtime deployment on Volcengine ECS runtime host (`srv-volcengine-sh-01`), off-server backup on Alibaba Cloud ECS (`srv-aliyun-bj-01`) | SVC-001 to SVC-010, PostgreSQL, reverse proxy, backup job | backend-engineer | Deployment, HTTPS endpoint/TLS evidence, security group evidence, co-location/capacity record, health, monitoring, backup/restore, smoke test; infrastructure-ops provides environment evidence |
| ACC-018 | P1 | SVC-009 reporting-service | SVC-002 identity-authz-service, source domain events | backend-engineer | Overview query, report authz-before-aggregate, BR-017 status/stage breakdown groups, zero-state DTO, report error |
| ACC-019 | P1 | SVC-003 lead-service, SVC-004 account-service | SVC-002 identity-authz-service | backend-engineer | Duplicate normalization, cross lead/account/contact safe lookup, warning token, proceed-after-warning, no merge/overwrite |
| ACC-020 | P1 | SVC-010 import-export-service | Target domain services, SVC-002 identity-authz-service, SVC-008 audit-history-service | backend-engineer | Import run, export run, row error, target routing, CSV formula safety, temp retention/deletion, operation log |
| ACC-021 | P1 | SVC-007 work-service | SVC-006 commercial-service, target record services, SVC-002 identity-authz-service | backend-engineer | Reminder query, Asia/Shanghai business date, due/overdue evaluation, reminder row DTO, permission |
| ACC-022 | P1 | SVC-008 audit-history-service | SVC-002 identity-authz-service, all mutation-producing services | backend-engineer | Append-only operation log event, tamper evidence, operation log query, audit permission |
| ACC-023 | P1 | SVC-009 reporting-service | SVC-002 identity-authz-service, source domain events | backend-engineer | Basic report query, report authz-before-aggregate, metric/breakdown `GroupRow` DTO covering leads-by-status, opportunities-by-stage, quotes/contracts/payments by status/amount (BR-017), active/default archive filter |

## Downstream Requirements

- G6 PSM must represent this map as service, bounded context, aggregate, data,
  event, permission, and contract mappings.
- G7 traceability must map each ACC item to tests and integration scenarios.
- G8 tasks must include service ID, service owner agent, contract references,
  acceptance ID, forbidden boundary access, and validation steps.
- Implementation cannot satisfy any mapped P0/P1 item by bypassing the owning
  service, using direct cross-service database access, or using a generic CRUD
  service.
