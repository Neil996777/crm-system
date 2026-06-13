# Audit Deny Event Fix - Traceability

Status: G5/G7/G8 handoff traceability for Claude audit.

## Deliverables

| ID | Deliverable | Purpose |
|---|---|---|
| DLV-AUDIT-YARDSTICK | `delivery/audit-deny-event-fix-acceptance.md` | Claude-owned ACC-AUDIT-001..006 and C1-C6 yardstick. |
| DLV-AUDIT-G5 | `docs/architecture/audit-deny-event-fix-design.md` | G5 root cause, architecture, affected services/contracts, and non-goals. |
| DLV-AUDIT-G8 | `delivery/audit-deny-event-fix-g8-task-package.md` | M1-M6 implementation, QA, and integration task package. |
| DLV-AUDIT-TRACE | `delivery/audit-deny-event-fix-traceability.md` | Acceptance and constraint traceability. |
| DLV-AUDIT-GATE | `planning/gate-status.md`, `planning/blockers.md` | Sync register for Claude G8 handoff audit and BLK-PROD-AUDIT-001 status. |

## Root Cause Traceability

| Finding | Evidence | Design response |
|---|---|---|
| `UserSignedIn` carries actor role/display and is the healthy reference. | `services/identity-authz/internal/handler/auth.go` appends `UserSignedIn` with `actorId`, `actorRole`, `actorDisplay`, `role`, `result`. | M1 preserves this shape and adds regression tests. |
| `UserAccessDenied` generic helper underfills the event. | `appendAccessDenied` appends only `actorId`, `reason`, `result`. | M1/M2 replace ad hoc construction with complete envelope. |
| Permission denied branch underfills even when actor is known. | `permissionCheck` denied branch has loaded `actor` but appends no `actorRole`/`actorDisplay`. | M2 passes known actor into the audit envelope and leaves decision unchanged. |
| User-admin denied branch loses known actor context. | `requireAdministrator` calls `appendAccessDenied(actor.ID, "user_admin_denied")`. | M2 emits `EVT-AUTH-ACCESS-DENIED` with Sales actor role/display and 403 unchanged. |
| Some auth-denied cases can have empty aggregate/resource id. | Unknown login, unauthenticated, and invalid-session paths call the generic helper with empty user id. | M1/M2 define non-empty safe audit resource id for unknown/unauthenticated cases. |
| `reason` is not promoted to `reasonCode`. | `auditAppendBody` maps event id/action but omits request `reasonCode`. | M1 adds reason-code mapping and tests. |
| Denied event stores broad payload in `afterSummary`. | `auditAppendBody` uses `afterSummary: item.Payload`. | M1/M4 require denied events to store `{}` before/after and safe summary only. |
| Audit-history 400 lacks structured reason. | `appendEvent` returns the same safe `INVALID_REQUEST` for missing body fields, missing actor headers, and repo append failure. | M3 adds structured server-side rejection diagnostics without relaxing validation. |

## ACC-AUDIT Mapping

| ACC | Required proof | Delivery mapping |
|---|---|---|
| ACC-AUDIT-001 | A denied action produces `UserAccessDenied`, audit-history accepts it, and Administrator sees it in Operation Logs. | DLV-AUDIT-G5 target contract; DLV-AUDIT-G8 M1/M2/M5/M6. |
| ACC-AUDIT-002 | Stored denied event includes actor id, actor role, action, result, and reason code with correct semantics. | DLV-AUDIT-G5 target contract; DLV-AUDIT-G8 M1/M2/M4/M5/M6. |
| ACC-AUDIT-003 | Well-formed audit appends no longer 400 and identity-authz outbox drains. | DLV-AUDIT-G5 producer validation; DLV-AUDIT-G8 M1/M2/M6. |
| ACC-AUDIT-004 | Malformed audit-history append rejections log structured reason. | DLV-AUDIT-G5 diagnostics design; DLV-AUDIT-G8 M3/M6. |
| ACC-AUDIT-005 | Deny-to-audit e2e/integration test exists and runs green. | DLV-AUDIT-G8 M5/M6. |
| ACC-AUDIT-006 | Existing audit events still append and hash chain remains valid. | DLV-AUDIT-G5 hash-chain design; DLV-AUDIT-G8 M1/M2/M4/M6. |

## Constraint Mapping

| Constraint | Design enforcement | Task enforcement |
|---|---|---|
| C1 safe-summary only; no raw before/after | Denied events store empty before/after; Operation Logs continues to render `safeSummary` through `AuditEventCard`. | M1, M2, M4, M5, M6. |
| C2 hash chain intact | Events continue through `EventRepo.Append`; no update/delete or hash bypass. | M4, M6. |
| C3 no role gates / enum / access-decision changes | Audit-only actor labels remain display data; permission outcomes unchanged. | M1, M2, M5, M6. |
| C4 zh-CN via `labels.ts` | Existing action labels used; new labels, if any, are display-only. | M5, M6. |
| C5 no-downgrade | `UserSignedIn`, `UserRoleStatusChanged`, `LastAdministratorBlocked`, existing operation-log access controls stay covered. | M1, M2, M4, M6. |
| C6 full gates | This package stops at G8 review; no code implementation before Claude pass. | Gate/status update and M6 evidence. |

## Service Boundary Mapping

| Service | Role | Owner agent | Contract |
|---|---|---|---|
| SVC-002 identity-authz-service | Producer of auth/user audit events and denied access events. | backend-engineer | `CONTRACT-002` identity event contract; complete audit envelope to SVC-008. |
| SVC-008 audit-history-service | Consumer/storage/query for operation logs and hash-chain events. | backend-engineer | `CONTRACT-013` append/query contract; strict validation plus diagnostics. |
| SVC-001 gateway-bff | Existing operation-log aggregate and browser session path. | backend-engineer | No contract change expected; E2E uses it for realistic deny/log flow. |
| frontend | Admin operation-log display. | frontend-engineer | SafeSummary-only `AuditEventCard`; zh-CN labels only. |

## Reviewer Hooks

| Reviewer | What to inspect |
|---|---|
| Claude G8 handoff audit | DLV-AUDIT-G5, DLV-AUDIT-G8, this traceability table, C1-C6, and no implementation-code diff. |
| Security Compliance | No credential/raw before-after leakage; no permission widening; operation-log admin-only behavior preserved. |
| Architecture | Producer/consumer contract completeness, service boundaries, hash-chain preservation. |
| QA Test Design | Deny-to-audit e2e/integration test plan and existing audit regression coverage. |
| Integration Owner | G11 evidence plan includes outbox drain, audit-history row, operation-log UI, and hash-chain proof. |
