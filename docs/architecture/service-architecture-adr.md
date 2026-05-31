# Service Architecture ADR

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## ADR-ARCH-001: Physical Multi-Service Deployment On One ECS

### Status

Accepted for G5 Re-review

### Context

The project is a production CRM for internal ToB sales team use. Company
standards default to microservice-oriented delivery and require
service-boundary-first governance. The user selected direct multi-service Go
deployment with Docker isolation.

### Decision

Use multiple Go backend services, each packaged and deployed as an independent
Docker container on one runtime host (`srv-volcengine-sh-01`, Volcengine ECS,
4 vCPU / 8 GiB) with Docker Compose. A second host (`srv-aliyun-bj-01`, Alibaba
Cloud ECS) serves as the off-server backup target only, not a runtime host.

### Consequences

- Service runtime boundaries exist from the committed release.
- Deployment and debugging are more complex than a modular monolith.
- Docker Compose is sufficient for the committed scale and current infrastructure.
- Kubernetes is not required for the committed scope.
- The runtime host is co-located with existing personal workloads, so CRM must
  honor disk/memory headroom and reverse-proxy/port-80 co-location constraints
  recorded in `deployment-notes.md`.

## ADR-ARCH-002: Service-Owned Data In One PostgreSQL Instance

### Status

Accepted for G5 Re-review

### Context

The user selected self-hosted PostgreSQL on the ECS server. Full database per
service infrastructure would add operational overhead. Shared database without
ownership controls would violate service governance.

### Decision

Use one PostgreSQL instance with independent service-owned database or schema
boundaries and independent database users.

Each service may access only its own database/schema. Cross-service database
access is forbidden.

### Consequences

- Operational cost is lower than running many PostgreSQL instances.
- Service ownership remains enforceable through database permissions.
- Cross-service joins are prohibited.
- Reporting and import/export must use events, public APIs, or owned read
  models.

## ADR-ARCH-003: No Unified Database CRUD Service

### Status

Accepted for G5 Re-review

### Context

A unified database-write service was discussed. It would expose database
operations through API and make all business services depend on one central data
operation service.

### Decision

Do not create a unified database CRUD service.

Each business service owns its data writes, validation, lifecycle rules, and
business invariants. Other services may request business capabilities through
public Command APIs, Query APIs, or event contracts.

### Consequences

- Preserves high cohesion and low coupling.
- Avoids a generic database API becoming the real business core.
- Requires stronger service contracts.
- Requires explicit integration and contract tests.

## ADR-ARCH-004: Local Automatic Backup With 7-Day Retention

### Status

Accepted for pre-release baseline; off-server target defined; production release
blocked until off-server backup copy and restore rehearsal evidence exist

### Context

The available infrastructure is two personal ECS servers: a Volcengine runtime
host and an Alibaba Cloud host. The user accepted local backup with 7-day
cleanup and confirmed the Alibaba host as the off-server backup target.

### Decision

Use automatic encrypted local PostgreSQL backup on the runtime host
(`srv-volcengine-sh-01`). Each backup creates a new timestamped file. Backups
older than 7 days are automatically deleted.

Before production release, copy encrypted backups off-server to `srv-aliyun-bj-01`
(Alibaba Cloud, Beijing — different provider and region) and prove restore
rehearsal evidence. Same-host-only backup cannot satisfy P0/P1 release
completion.

### Consequences

- Meets the current local backup baseline.
- Requires restore rehearsal before release.
- Local-only backup does not provide disaster recovery while it shares the
  runtime host; the off-server copy to `srv-aliyun-bj-01` (different provider and
  region) provides that separation once copy evidence exists.
- Production release requires off-server backup copy evidence and restore
  rehearsal, or a formal scope change by the user.

## ADR-ARCH-005: HTTPS-Only Production Ingress And Service Authentication

### Status

Accepted for G5 Re-review

### Context

The CRM handles restricted customer, contact, quote, contract, payment, user,
and audit data. Docker internal networking does not prove trusted caller
identity. Production login/session traffic cannot use plaintext HTTP.

### Decision

Use HTTPS-only production ingress through the reverse proxy. HTTP may redirect
to HTTPS. Browser sessions use secure cookies. Internal service calls require
authenticated service identity, actor context when user initiated, allowed
intent, correlation ID, and target-service verification.

### Consequences

- A valid production HTTPS endpoint and TLS evidence are required before
  production release.
- Service-to-service credentials, rotation, and rejection behavior must be
  represented in PSM and tasks.
- Security Compliance must review this before G5 can pass.
