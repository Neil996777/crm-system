# G5 Architecture Repair Note

## Document Control

- Project: CRM System
- Date: 2026-05-30
- Role: Architecture
- Purpose: Record the architecture repair work after blocked G5 review.
- Archive Note: This file records repair evidence only. Active design authority
  remains under `docs/architecture/` and current upstream design documents.

## Scope

No implementation code was written. No P0/P1 product acceptance item was
downgraded, deleted, merged away, weakened, or accepted as partial work.

## Repaired Active Documents

- `docs/architecture/architecture.md`
- `docs/architecture/api-spec.md`
- `docs/architecture/authz-architecture.md`
- `docs/architecture/data-design.md`
- `docs/architecture/integration-design.md`
- `docs/architecture/frontend-backend-contract.md`
- `docs/architecture/deployment-notes.md`
- `docs/architecture/risk-register.md`
- `docs/architecture/service-architecture-adr.md`
- `docs/architecture/service-architecture-acceptance.md`
- `docs/architecture/service-acceptance-map.md`
- `docs/product/open-questions.md`
- `PROJECT_CONTEXT.md`

## Blocker Closure Mapping

| Review Item | Repair Summary | Active Design Reference |
|---|---|---|
| G5-BLK-001 | Added endpoint strategy, HTTPS/TLS evidence, security group evidence, monitoring target, deployment evidence ownership. | `deployment-notes.md`, `PROJECT_CONTEXT.md`, `open-questions.md` |
| G5-BLK-002 | Marked same-host-only backup as a production release blocker and required encrypted off-server backup evidence. | `deployment-notes.md`, `data-design.md`, `service-architecture-adr.md`, `risk-register.md` |
| G5-BLK-003 | Defined concrete service-to-service authentication, credential storage, rotation, rejection behavior, and caller verification. | `authz-architecture.md`, `api-spec.md` |
| G5-BLK-004 | Defined opaque secure session, server-side session store, logout, disabled-user handling, role/status version recheck, and session revocation. | `authz-architecture.md`, `api-spec.md` |
| G5-BLK-005 | Defined backup directory permissions, encryption, secret handling, restore logging, checksum, and restricted restore handling. | `deployment-notes.md`, `data-design.md` |
| G5-BLK-006 | Defined HTTPS-only production ingress, HTTP redirect, secure cookie transport, reverse proxy security headers, and no public debug/admin exposure. | `deployment-notes.md`, `authz-architecture.md` |
| G5-BLK-007 | Added archive eligibility API, active obligation DTO, blocked archive response, retry/refresh behavior, and history/log events. | `architecture.md`, `api-spec.md`, `frontend-backend-contract.md`, `integration-design.md` |
| G5-BLK-008 | Added `version`/`expectedVersion` concurrency contract and `VERSION_CONFLICT` response. | `api-spec.md`, `data-design.md`, `frontend-backend-contract.md` |
| G5-BLK-009 | Added owner transfer contract, OwnerChanged/OpenWorkTransferred flow, manual exception, retry/failure behavior, and history events. | `architecture.md`, `api-spec.md`, `integration-design.md` |
| G5-BLK-010 | Added Close Lost contract, required lost reason, terminal edit protection, post-close work-service path, and audit/reporting events. | `architecture.md`, `api-spec.md`, `data-design.md` |
| G5-ISS-001 | Added duplicate normalization, safe lookup, warning token, proceed-after-warning, and no merge/overwrite rule. | `architecture.md`, `api-spec.md`, `frontend-backend-contract.md` |
| G5-ISS-002 | Added report metric DTO, groupings, amount fields, zero state, active/default archive filter, and auth-before-aggregation rule. | `api-spec.md`, `frontend-backend-contract.md`, `service-acceptance-map.md` |
| G5-ISS-003 | Added reminder timezone/business-date rules and reminder row DTO. | `api-spec.md`, `frontend-backend-contract.md`, `service-acceptance-map.md` |
| G5-ISS-004 | Added import/export object scope, target service routing, row validation, CSV formula safety, temp retention/deletion, auth, and operation logs. | `api-spec.md`, `integration-design.md`, `frontend-backend-contract.md` |
| G5-ISS-005 | Added append-only API rule, no update/delete admin UI, DB permission constraints, and tamper-evidence hash fields. | `data-design.md`, `service-acceptance-map.md` |
| G5-ISS-006 | Added runtime paths, resource threshold expectations, monitoring/alerting evidence, and upgrade trigger. | `deployment-notes.md`, `risk-register.md` |
| G5-ISS-007 | Added deploy/ops user plan and root SSH boundary before production release. | `deployment-notes.md` |

## Remaining Gate Status

Architecture repair is ready for G5 re-review. G5 is not passed until required
reviewers approve the revised architecture without open P0/P1 blockers.
