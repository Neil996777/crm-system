# Authorization Architecture

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Strategy

Authorization is enforced by backend services. Frontend navigation, disabled
buttons, and hidden UI elements are convenience hints only and are never the
security boundary.

The CRM uses a single-team, three-role model for the committed release:

- Administrator
- Sales Manager
- Sales

Authorization combines:

- actor identity
- role
- active/disabled user status
- action
- resource type
- resource ownership or assignment
- related-parent scope where needed
- archived/default filter rules

## Service Responsibilities

| Service | Authorization Responsibility |
|---|---|
| gateway-bff | Authenticate request context, forward actor/correlation context, never make final domain authorization decisions alone. |
| identity-authz-service | Own user, role, session, permission policies, active/disabled checks, last Administrator protection, and permission decision API. |
| domain services | Enforce domain-specific authorization and business rules before mutation/query. |
| audit-history-service | Enforce record-local history visibility and admin-only global log query. |
| reporting-service | Enforce authorization before aggregates are returned. |
| import-export-service | Enforce import/export role/scope before starting runs and target service calls. |

## Permission Contract

Every protected action must be represented as:

```text
actor + action + resource + scope + condition -> allowed / denied
```

Required fields:

- actor ID
- actor role
- actor active/disabled status
- action ID
- resource type
- resource ID where available
- owner ID or assignee ID where applicable
- related parent ID where applicable
- archived/default state where applicable
- correlation ID

## Denial Contract

Permission denial must return a safe category:

- unauthenticated
- disabled user
- role denied
- scope denied
- resource unavailable
- archived/default hidden
- admin only

The API must not reveal restricted record names, amounts, contact values, or
whether an unauthorized record exists.

## Required Permission Enforcement Points

| Capability | Enforcement |
|---|---|
| Login/session | identity-authz-service rejects invalid, disabled, or expired sessions. |
| User/role management | Administrator only, with last active Administrator protection. |
| Lead/account/opportunity/commercial/work writes | Owning service checks permission and business rule before mutation. |
| List/detail/search/filter | Owning service or approved read model filters by authorized scope before returning rows. |
| Record-local history | audit-history-service checks related record visibility or approved scope summary. |
| Global operation log | Administrator only. |
| Import/export | Administrator or Sales Manager only; Sales denied. |
| Reports | Administrator all authorized scope; Sales Manager team scope; Sales denied for manager/admin reports. |
| Archive | Administrator and Sales Manager only; no hard delete. |

## Service-To-Service Authorization

Internal calls are not trusted only because they are internal.

Each internal call must include:

- caller service identity
- user actor context if user-initiated
- purpose/action
- resource reference
- correlation ID

Target services must verify whether the caller may perform the requested
purpose. This prevents a service from becoming a generic bypass around domain
authorization.

Concrete G5 mechanism:

- Every service has a stable `serviceId`.
- Internal requests include `Authorization: Bearer <service-token>`,
  `X-Service-Id`, `X-Correlation-Id`, `X-Intent`, and actor context when the
  call is user initiated.
- The service token is a signed token or equivalent HMAC/JWS credential with
  key ID, issuer service, audience service, expiry, and allowed intent. Maximum
  token lifetime is 5 minutes.
- Service credentials are stored as Docker secrets or root-readable environment
  files outside the repository. They are never committed.
- Rotation uses dual active keys by key ID. Old keys are removed only after all
  services have accepted the new key.
- Target services reject missing, expired, invalidly signed, wrong-audience, or
  disallowed-intent calls with `SERVICE_AUTH_FAILED`.
- Target services must still enforce domain authorization and business rules;
  a valid service token only proves caller identity.
- Service auth failures and sensitive denied calls emit operation log events.

## Session And Role Recheck

The architecture must support:

- active/disabled user recheck on protected requests
- role/status change invalidation or effective re-evaluation
- safe handling of stale sessions
- operation log event for sensitive access or user lifecycle changes

Concrete G5 session strategy:

- Browser sessions use an opaque session ID in an `HttpOnly`, `Secure`,
  `SameSite=Lax` cookie for production HTTPS traffic.
- identity-authz-service owns the session store, session expiry, logout, and
  revocation state.
- Default session lifetime is 12 hours absolute or 30 minutes idle, whichever
  expires first. G6 PSM may shorten these values but may not remove expiry or
  server-side revocation.
- Protected requests revalidate active user status and role/authz version
  through identity-authz-service or a bounded short-lived cache with explicit
  recheck TTL no longer than 60 seconds.
- Role or status changes increment an authz version and revoke or force
  re-evaluation of affected sessions before further protected mutations.
- Disabled users are rejected on the next protected request and all active
  sessions for that user are invalidated.
- Logout revokes the session server-side and clears the cookie.
- Session expiry, user disabled, role changed, and logout return distinct safe
  authentication errors: `SESSION_EXPIRED`, `USER_DISABLED`,
  `AUTHZ_VERSION_STALE`, and `SESSION_REVOKED`.
- User lifecycle and sensitive access changes create operation log events.

## Audit Requirements

Permission-sensitive and business-sensitive operations must produce history or
operation log events as required by the PRD and security documents.

Examples:

- login/access failure
- owner change
- stage/status change
- quote accepted
- contract status changed
- payment recorded
- archive action
- import/export run
- denied access where required

## Security Blocker Closure

This document addresses the G5 architecture design side of:

- SEC-SVC-BLK-001: service trust boundaries
- SEC-SVC-BLK-002: service-to-service authorization and audit behavior
- SEC-SVC-BLK-004: error and denial contracts

Security Compliance must review these rules before G5 can pass.
