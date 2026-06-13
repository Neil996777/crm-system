# Audit Deny Event Fix - G5 Architecture Design

Status: Ready for Claude G8 handoff audit. Implementation is not started.

## Document Control

- Change: audit access-denied event fix for `BLK-PROD-AUDIT-001`.
- Yardstick: `delivery/audit-deny-event-fix-acceptance.md`.
- Product baseline: `docs/product/acceptance-matrix.md` ACC-022.
- Security baseline: `docs/security/audit-log-spec.md`, especially
  `EVT-AUTH-LOGIN-FAILED`, `EVT-AUTH-ACCESS-DENIED`, safe-summary, and
  append-only/tamper-evident requirements.
- Producer: Codex, G5/G7/G8 production lane.
- Auditor: Claude, independent G8 handoff audit before implementation may start.

## Scope

This design fixes application-layer audit completeness for denied-access events.
It is not part of the CI/CD migration and does not inherit that migration's
"mechanism only / zero app source" constraint. After G8 passes, implementation
is expected to touch application code in:

- `services/identity-authz`
- `services/audit-history`
- tests under the affected services and `frontend/e2e`
- possibly `frontend/src/i18n/labels.ts` for display-only labels
- possibly `shared/contracts` if the implementation chooses to make the audit
  envelope a shared compile-time contract

The design forbids any authorization weakening. Role gates, role comparison
values, permission decisions, and access denial outcomes remain unchanged.

## Current-Code Root Cause

### Identity-authz produces incomplete denied-event payloads

`UserSignedIn` is the reference healthy shape. In
`services/identity-authz/internal/handler/auth.go`, successful sign-in appends
`UserSignedIn` with:

- `actorId`
- `actorRole`
- `actorDisplay`
- `role`
- `result`

By contrast, the generic denied helper `appendAccessDenied` appends
`UserAccessDenied` with only:

- `actorId`
- `reason`
- `result`

The permission-check denied branch also appends `UserAccessDenied` without
`actorRole` or `actorDisplay`, even though it has already loaded the actor.
The user-admin denied branch calls the same generic helper, so a Sales user
denied from administrator-only user management also loses role/display at the
event-construction point.

`UserSignedOut` is another identity-authz event that currently carries only
`actorId` and `result`; it may be accepted through dispatcher fallback, but it is
semantically underfilled and must be included in the audit-envelope cleanup so
the fix is not limited to one event type.

### Dispatcher fallback is not an acceptable contract

`services/identity-authz/internal/event/outbox.go` currently derives S2S actor
headers from payload keys:

- `X-Actor-User-Id` from `actorId`, fallback `system`
- `X-Actor-Role` from `actorRole`, fallback `System`
- `X-Actor-Display` from `actorDisplay`, fallback actor id

This fallback may keep some malformed events non-empty, but it masks the true
actor role with `System`, does not satisfy ACC-AUDIT-002 semantics, and does not
make the producer payload complete. The fix must make identity-authz producers
construct a complete audit envelope before enqueueing, then make the dispatcher
fail tests if an identity audit event would rely on fallback for authenticated
actor fields.

### Audit-history rejects incomplete append requests with 400

`services/audit-history/internal/handler/server.go` accepts
`POST /internal/events/append` only when the append body has:

- `eventId`
- non-zero `eventVersion`
- `action`
- `resourceType`
- `resourceId`
- `result`
- `safeSummary`
- at least one `surface`
- at least one `acceptanceId`

It then sets actor fields from S2S headers and rejects the request if any of:

- `X-Actor-User-Id`
- `X-Actor-Role`
- `X-Actor-Display`

are empty. The underlying table
`services/audit-history/migrations/0002_history_oplog.up.sql` also makes
`actor_user_id`, `actor_role`, `actor_display`, `action`, `resource_type`,
`resource_id`, `result`, `safe_summary`, `acceptance_ids`, `prev_hash`, and
`event_hash` non-null.

Therefore a malformed identity event can be rejected either because the actor
headers are incomplete or because `auditAppendBody` emits an empty `resourceId`
for unauthenticated/unknown-user denial cases where the outbox aggregate id is
empty. The design must satisfy both conditions rather than weakening
audit-history validation.

### Reason and safe-summary contract is incomplete

`auditAppendBody` maps `UserAccessDenied` to `EVT-AUTH-LOGIN-FAILED` or
`EVT-AUTH-ACCESS-DENIED`, but it does not promote payload `reason` to
audit-history `reasonCode`. It also stores the entire payload in `afterSummary`.
For denied events this is too broad. The yardstick requires the deny event to
land with a safe summary and reason code, without raw before/after payloads.

## Target Contract

All identity-authz audit events sent to audit-history must use this effective
contract.

| Field | Source | Requirement |
|---|---|---|
| `eventUid` | outbox row id | Stable producer event uid for idempotent retry. |
| `eventId` | catalog mapper | `EVT-AUTH-LOGIN-SUCCEEDED`, `EVT-AUTH-LOGIN-FAILED`, `EVT-AUTH-ACCESS-DENIED`, `EVT-USER-ROLE-CHANGED`, `EVT-USER-STATUS-CHANGED`, `EVT-LAST-ADMIN-BLOCKED`, or existing mapped value. |
| `eventVersion` | producer | `1` unless catalog version changes. |
| `surfaces` | producer contract | Includes `operation_log` for these global audit events. |
| `action` | producer/catalog | Stable display action such as `sign_in`, `sign_out`, `login_failed`, `access_denied`, `change_role`, `change_status`, `last_admin_blocked`. |
| `resourceType` | producer/catalog | `Auth` for authentication failures, `User` for user-admin denies and user events, or safe resource type for permission-check denies. |
| `resourceId` | producer/catalog | Non-empty safe id. For unknown/unauthenticated actors use a stable non-secret audit resource id such as `anonymous` or `auth`. Do not store submitted email/password. |
| `result` | producer | `success`, `denied`, `blocked`, or `failed` as applicable. |
| `reasonCode` | producer | Required for denied/failed/blocked events, derived from safe reason constants. |
| `beforeSummary` | producer | `{}` for denied auth/access events. |
| `afterSummary` | producer | `{}` for denied auth/access events. Existing successful mutation events may continue to use safe summaries. |
| `diffClassification` | producer | `Security Critical` for auth/access denial events unless a stricter existing classification applies. |
| `scopeSummary` | producer | Safe scope label such as `administrator only`, `permission denied`, or `authentication`. |
| `safeSummary` | producer | Safe text only. It may include action/result/reason labels but no credential, session token, restricted record name, or raw before/after payload. |
| `correlationId` | request/outbox | Preserve request correlation where present; fallback to event uid. |
| `causationId` | outbox row id | The producer outbox event id. |
| `acceptanceIds` | producer | At least `ACC-022`; include `ACC-AUDIT-*` in tests/evidence, not as a replacement for product ACC-022. |
| actor headers | complete producer envelope | Non-empty actor id, role, and display. For unauthenticated/unknown actor, use an audit-only actor label such as `anonymous` / `Unauthenticated` / `Anonymous`, never a real domain role comparison value. |

Audit-only actor role labels must remain display data. They must not be accepted
as authorization roles and must not be added to domain role comparisons in
`services/identity-authz/internal/domain`.

## Architecture Design

### 1. Identity-authz envelope builder

Introduce one identity-authz audit-envelope construction path used by every
identity audit event producer. The builder owns:

- actor id, role, and display resolution
- safe unauthenticated/unknown actor fallback
- resource type/id fallback for auth-denied cases
- reason code mapping
- safe summary creation
- before/after behavior for denied events
- action/catalog mapping inputs

Known actors should use the persisted `domain.User` data already loaded on the
request path. Where the path only has a user id, the implementation may look up
the user through the existing `UserRepo` inside the same durable workflow. If the
actor cannot be found, the audit event must still be well formed with an
audit-only unknown actor label and no credential details.

This is a producer-side fix. Audit-history validation should stay strict.

### 2. Denied-path producers

Update these identity-authz producers after G8 passes:

| Producer path | Current problem | Target |
|---|---|---|
| `signIn` unknown email / wrong password / inactive user | `appendAccessDenied` emits only `actorId`, `reason`, `result`; unknown actor may leave aggregate/resource id empty. | Emit complete login-failed envelope with safe unknown/known actor, `resourceType=Auth`, non-empty `resourceId`, `reasonCode=login_failed`, empty before/after, safe summary. |
| `authenticate` missing cookie / invalid session / inactive user / stale authz version | Same generic helper, no role/display, possible empty resource id. | Emit complete denied envelope; use known user role/display when safely available, otherwise audit-only unknown actor. |
| `requireAdministrator` non-admin user | Generic helper drops the already-loaded actor role/display. | Emit complete `EVT-AUTH-ACCESS-DENIED` envelope using the loaded user, `reasonCode=user_admin_denied`, access still returns 403. |
| `permissionCheck` denied decision | Denied branch drops loaded actor role/display and stores broad payload in `afterSummary`. | Emit complete denied envelope using loaded actor, safe resource type, no restricted resource name/id leak beyond accepted safe identifiers, `reasonCode=decision.DenialCategory`, access decision unchanged. |
| `actor_not_found` permission check | No actor role/display possible. | Emit well-formed audit-only unknown actor envelope without treating unknown actor as an authorized role. |
| `signOut` | Underfilled actor envelope. | Include actor role/display or safe unknown actor fallback; existing sign-out behavior unchanged. |

### 3. Dispatcher and body mapping

Update identity-authz dispatcher tests so the append body and S2S headers prove:

- no authenticated denied event depends on `System` actor fallback
- denied events set `reasonCode`
- denied events set empty `beforeSummary` and `afterSummary`
- denied events have non-empty `resourceId`
- `safeSummary` is non-empty and safe
- existing `UserSignedIn`, `UserRoleStatusChanged`, and `LastAdministratorBlocked`
  mappings still produce the same catalog ids and accepted append body shape

Dispatcher fallback can remain as a defensive fallback for legacy/unexpected
events, but tests must prevent the known identity-authz catalog events from
using it silently.

### 4. Audit-history structured rejection diagnostics

Keep audit-history's strict 400 behavior for malformed append requests. Add
structured server-side logging when append validation rejects an event:

- event uid if present
- producer service from verified token claims
- correlation id from header/body if present
- missing fields
- invalid fields
- validation stage: JSON decode, body-required-fields, actor-headers, or repo
  append

The response body should remain safe and generic. The log must not include
credentials, tokens, raw request body, password material, or raw before/after
payloads.

### 5. Safe-summary and frontend display

The Operation Logs page already renders `AuditEventCard` with `safeSummary` and
does not render raw before/after fields. This invariant must be preserved.
If new display-only labels are needed for audit-only actor labels or new reason
codes, add them through `frontend/src/i18n/labels.ts`. Do not change enum or role
comparison values. `access_denied` and `login_failed` action labels already
exist and should continue to be used.

If the implementation exposes `reasonCode` through the operation-log DTO so the
UI can show the denial reason explicitly, that field must be a safe reason-code
string only. It must not expose raw request context. If no DTO field is added,
the safe summary must still include enough safe denial context for ACC-AUDIT-002
verification.

### 6. Hash-chain preservation

The audit-history append flow computes `prevHash` from the latest accepted event
and then computes `eventHash` over the canonical event payload including actor,
safe summaries, reason code, and `prevHash`. The fix must not bypass
`EventRepo.Append`, must not update/delete accepted rows, and must not alter the
hash computation except through legitimate event field values. Tests must append
at least two events after the fix and verify the second `prevHash` equals the
first `eventHash`.

## Service And Contract Impact

| Service / package | Impact |
|---|---|
| identity-authz-service | Primary producer fix for complete audit envelopes and denied-path event construction. |
| audit-history-service | Add structured rejection logging; keep strict validation and append-only hash-chain behavior. May safely expose `reasonCode` in operation-log DTOs if needed for ACC-AUDIT-002 UI proof. |
| gateway-bff | No contract change expected. Existing `/api/operation-log` aggregate should continue to query audit-history as Administrator only. |
| frontend | Test and possibly label-only additions for zh-CN display; no operation-log raw before/after rendering. If `reasonCode` is displayed, it must go through `labels.ts` or an equivalent safe label map. |
| shared/contracts | Optional compile-time contract strengthening if implementation chooses; any change must be additive and not force unrelated services to rewrite stable event producers. |

## Non-Goals

- No change to permission rules or role comparison values.
- No widening of Sales or Sales Manager access.
- No raw password, credential, session token, submitted email, restricted record
  name, or raw before/after payload in denied audit events.
- No historical backfill of already lost denied events.
- No CI/CD pipeline redesign.

## G8 Audit Hooks

Claude's G8 audit should verify this package answers:

- Root cause is tied to `identity-authz` event construction and
  `audit-history` append validation, not to deployment mechanics.
- ACC-AUDIT-001..006 are each mapped to tasks and tests.
- C1-C6 are present as task constraints.
- The plan changes audit fields/contracts only and explicitly forbids access
  decision changes.
- Implementation is not started before the G8 handoff audit passes.
