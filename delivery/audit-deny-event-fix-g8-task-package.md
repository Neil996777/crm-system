# Audit Deny Event Fix - G7/G8 Task Package

Status: Ready for Claude G8 handoff audit. Implementation is not started.

This package turns the audit-deny G5 design into executable G9-G11 work after
Claude passes the G8 handoff audit. Codex does not self-pass G8 or G12.

## Inputs

- Yardstick: `delivery/audit-deny-event-fix-acceptance.md`
- G5 design: `docs/architecture/audit-deny-event-fix-design.md`
- Product acceptance: `docs/product/acceptance-matrix.md` ACC-022
- Security audit spec: `docs/security/audit-log-spec.md`
- Current blocker: `planning/blockers.md` BLK-PROD-AUDIT-001

## Global Constraints

- This is an application-code change only after G8 passes.
- Allowed implementation areas after G8: `services/identity-authz`,
  `services/audit-history`, related tests, `frontend/e2e`, and display-only
  `frontend/src/i18n/labels.ts`; `shared/contracts` only for additive audit
  contract typing if needed.
- Forbidden changes: role gates, authorization decisions, enum values, role
  comparison values, permission policy semantics, access-denied HTTP outcomes,
  operation-log access control, audit-history append-only semantics, and test
  weakening.
- Denied events must store safe summary and reason code only; no raw
  before/after payload, credentials, session tokens, submitted passwords,
  restricted record names, or secret values.
- Operation Log UI must continue to render `safeSummary` only.
- G9 cannot start until Claude passes this G8 handoff package.

## Task DAG

```text
M1
├── M2
│   └── M5
├── M3
│   └── M5
└── M4
    └── M6

M5 └── M6
```

M1 establishes the local audit envelope contract. M2 and M3 can proceed in
parallel after M1. M4 can proceed after M1 and should run before final
integration. M5 consumes the producer and consumer fixes. M6 produces G11
evidence.

## Tasks

| ID | Owner agent | Objective | Primary files after G8 | Verification |
|---|---|---|---|---|
| M1 | Backend Engineer + Architecture | Define and test the complete identity-authz audit envelope builder and catalog body mapping. | `services/identity-authz/internal/event/*`, possibly `shared/contracts/*` | Unit tests fail first for missing actor role/display, missing reasonCode, non-empty resource id, empty deny before/after. |
| M2 | Backend Engineer | Fill all identity-authz denied and underfilled identity event producers with the complete envelope. | `services/identity-authz/internal/handler/auth.go`, `permission.go`, `user_admin.go`, tests | Known and unknown actor deny paths append well-formed outbox rows; access decisions remain denied. |
| M3 | Backend Engineer | Add audit-history structured rejection diagnostics while preserving strict 400 validation. | `services/audit-history/internal/handler/server.go`, tests | Malformed append returns safe 400 and logs structured missing-field reason without secrets. |
| M4 | Backend Engineer + QA Execution | Preserve audit-history append-only hash-chain and safe-summary behavior under the new deny events. | `services/audit-history/internal/*`, identity dispatcher tests | Accepted events maintain `prevHash/eventHash`; denied events have empty before/after and safe summary. |
| M5 | QA Execution + Frontend Engineer | Add deny-to-audit integration/E2E coverage: trigger denial, then verify Administrator operation log visibility and Sales/Manager denial. | `frontend/e2e/oplog.spec.ts`, service integration tests, optional `labels.ts` and safe operation-log DTO fields | E2E proves denied action stays denied and the operation log shows the event via zh-CN safe-summary display. |
| M6 | Integration Owner + QA Execution | Produce G11 evidence and regression proof for ACC-AUDIT-001..006 and C1-C6. | `delivery/audit-deny-event-fix-g11-evidence.md` after implementation | Full relevant Go tests, frontend e2e, outbox drain proof, operation-log proof, hash-chain proof, no skip/only, no authz diff beyond audit fields. |

## Detailed Task Contracts

### M1 - Identity Audit Envelope Contract

Acceptance: ACC-AUDIT-001, ACC-AUDIT-002, ACC-AUDIT-003, ACC-AUDIT-006.

Objective:

- Add one local identity-authz construction path for audit actor envelope,
  reason code, safe summary, resource type/id, and deny before/after behavior.
- Add tests that fail on the current code because `UserAccessDenied` lacks
  producer-supplied `actorRole`/`actorDisplay`, lacks `reasonCode`, and can emit
  empty `resourceId`.

Required test cases:

- `UserSignedIn` remains accepted with actor id/role/display.
- `UserAccessDenied` for known actor includes actor id/role/display and
  `reasonCode`.
- `UserAccessDenied` for unauthenticated/unknown actor has a non-empty audit-only
  actor id/role/display and non-empty `resourceId`.
- Denied event append body has `{}` before/after summaries.
- Catalog ids remain `EVT-AUTH-LOGIN-FAILED` and `EVT-AUTH-ACCESS-DENIED` for
  the relevant denied cases.

Forbidden:

- Adding audit-only actor roles to domain authorization comparisons.
- Changing role enum values or permission policy logic.

### M2 - Producer Path Fixes

Acceptance: ACC-AUDIT-001, ACC-AUDIT-002, ACC-AUDIT-003, ACC-AUDIT-006.

Objective:

- Replace ad hoc `appendAccessDenied(ctx, userID, reason)` construction with the
  complete envelope helper or equivalent typed path.
- Fill the underfilled `UserSignedOut` event.
- Preserve current request outcomes:
  - bad login stays 401
  - unauthenticated current-user stays 401
  - stale authz stays the existing safe authz error
  - Sales denied from `/admin/users` stays 403
  - permission-check denied response still says `allowed=false`

Required covered paths:

- login failed for unknown email
- login failed for known user
- unauthenticated session check
- invalid/stale session
- inactive user
- user-admin denied
- permission-check denied
- actor-not-found permission check
- sign-out event envelope

Verification:

- Service tests inspect outbox payloads for complete envelope fields.
- Dispatcher integration sends events to a test audit-history server and drains
  unpublished outbox rows after a 2xx append.
- Tests assert no permission decision changed.

### M3 - Audit-History Rejection Diagnostics

Acceptance: ACC-AUDIT-004.

Objective:

- Add structured diagnostic logging for append rejection.
- Preserve safe generic HTTP responses.
- Keep strict validation and S2S authentication.

Required diagnostics:

- missing body fields, e.g. `eventId`, `resourceId`, `safeSummary`
- missing actor headers, e.g. `actorRole`, `actorDisplay`
- repo append/hash/DB validation failure category
- producer service, event uid, and correlation id when available

Forbidden:

- Logging bearer token, service token, cookie, password, raw request body, raw
  before/after payload, or submitted credentials.
- Downgrading 400 to 2xx for malformed events.

Verification:

- Test sends a malformed signed append request and captures logs.
- Test checks the structured reason names and verifies forbidden sensitive
  strings are absent.

### M4 - Hash Chain And Safe-Summary Regression

Acceptance: ACC-AUDIT-002, ACC-AUDIT-006; constraints C1 and C2.

Objective:

- Prove new deny events travel through normal `EventRepo.Append`.
- Prove accepted events maintain tamper-evident chain.
- Prove denied event rows do not carry raw before/after payloads.

Verification:

- Append a successful login/admin event and a denied event; verify second
  `prevHash` equals first `eventHash`.
- Verify `eventHash` is non-empty for both.
- Verify denied row has safe summary, reason code, actor role/display, and empty
  before/after summaries.
- Re-run existing audit-history append/idempotency tests.

### M5 - Deny To Operation Log Integration / E2E

Acceptance: ACC-AUDIT-001, ACC-AUDIT-002, ACC-AUDIT-005; constraints C1, C3, C4.

Objective:

- Add an end-to-end test that triggers a real denied action and then verifies an
  Administrator can see the corresponding operation-log event.

Preferred E2E scenario:

1. Administrator creates a Sales user.
2. Sales signs in.
3. Sales attempts an administrator-only endpoint such as `/admin/users`.
4. The endpoint still returns 403 and does not expose administrator data.
5. Administrator signs back in and opens Operation Logs.
6. Operation Logs show the access-denied event with actor display, actor role,
   action, result, reason/safe summary, event id, and event hash.
7. The same page still renders no raw before/after fields and no edit/delete
   controls.

If `reasonCode` is surfaced in the UI, it must be a safe reason code from
audit-history and translated through `labels.ts` or an equivalent safe label
map. It must not expose raw request context. If the UI does not add a distinct
reason field, the safe summary must carry enough safe denial context for the
test to verify ACC-AUDIT-002.

Alternative service-level integration is acceptable only as a supplement; the
G11 package must still include one real user-facing or gateway-backed deny to
operation-log proof.

Verification:

- `frontend/e2e/oplog.spec.ts` or an equivalent new spec covers the scenario.
- Tests use existing login/session/gateway flows, not seeded audit rows alone.
- E2E continues to run with `workers:2` and `retries:1`; no `test.skip`,
  `test.only`, assertion weakening, or fake operation-log fixture may satisfy
  ACC-AUDIT-005.

### M6 - G11 Evidence And Regression Closure

Acceptance: ACC-AUDIT-001..006; constraints C1-C6.

Objective:

- Produce a concise evidence report after implementation and QA.
- Keep `BLK-PROD-AUDIT-001` open until Claude G12 independently passes the fix.

Required evidence:

- Go test commands and results for `services/identity-authz` and
  `services/audit-history`.
- Frontend E2E command/result for the deny-to-operation-log scenario and full
  relevant suite.
- Outbox drain proof for well-formed denied events.
- Audit-history row proof for actor id/role/display, action, result, reasonCode,
  safeSummary, eventHash, prevHash.
- Operation Log UI proof showing safeSummary-only rendering and zh-CN labels.
- Static diff summary proving role gates, enum/role comparison values, and
  permission policy semantics were not changed except audit-envelope fields.
- No secret values in evidence.

## ACC-AUDIT To Task Map

| ACC | Task coverage |
|---|---|
| ACC-AUDIT-001 | M1, M2, M5, M6 |
| ACC-AUDIT-002 | M1, M2, M4, M5, M6 |
| ACC-AUDIT-003 | M1, M2, M6 |
| ACC-AUDIT-004 | M3, M6 |
| ACC-AUDIT-005 | M5, M6 |
| ACC-AUDIT-006 | M1, M2, M4, M6 |

## Constraint To Task Map

| Constraint | Task coverage |
|---|---|
| C1 safe-summary only; deny events carry no raw before/after | M1, M2, M4, M5, M6 |
| C2 hash chain intact | M4, M6 |
| C3 no authz/role/enum decision change | M1, M2, M5, M6 |
| C4 zh-CN through `labels.ts` | M5, M6 |
| C5 no-downgrade existing audit | M1, M2, M4, M6 |
| C6 full gates; app code allowed only after G8 | This package; M6 evidence |

## Handoff Condition

After Claude passes the G8 handoff audit, Codex may enter G9 implementation for
M1-M4, then G10/G11 for M5-M6. Claude performs independent G12 audit before this
blocker can be resolved or any release decision can rely on the fix.
