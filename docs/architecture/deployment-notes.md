# Deployment Notes

## Document Control

- Project: CRM System
- Phase: G5 Architecture Design
- Owner Agent: Architecture
- Status: Revised for G5 Re-review
- Date: 2026-05-30

## Deployment Target

The v1 production runtime host is one Volcengine ECS server, with a second
Alibaba Cloud ECS server as the off-server backup target. Both assets are
recorded in `company/infrastructure/server-inventory.md`.

| Role | Registered Asset | Provider / Region | Public IP | Spec |
|---|---|---|---:|---|
| Production runtime host | `srv-volcengine-sh-01` | Volcengine / East China 2, Shanghai | `118.196.44.193` | 4 vCPU / 8 GiB / 40 GiB |
| Off-server backup target | `srv-aliyun-bj-01` | Alibaba Cloud / North China 2, Beijing | `47.95.119.211` | 2 vCPU / 2 GiB / 40 GiB |

Cross-provider and cross-region separation between the runtime host (Volcengine,
Shanghai) and the backup target (Alibaba Cloud, Beijing) provides disaster-
recovery isolation for the off-server backup requirement.

Runtime components (on the runtime host):

- reverse proxy: CRM ingress served behind the existing host Nginx (see
  co-location constraints); CRM must not claim exclusive ownership of port 80
- web app / gateway-bff
- Go backend service containers
- PostgreSQL container
- backup job container or host cron job
- optional log/metrics sidecar or lightweight collector

### Runtime Host Co-location Constraints

The runtime host already runs unrelated personal workloads (Nginx on port 80,
Hermes on `8642`, OpenClaw gateway). CRM is co-located and must not assume an
empty host.

- Port 80 is already used by the existing host Nginx. CRM HTTPS ingress must be
  routed through the existing reverse proxy (a dedicated `server_name` or
  subpath), not by taking over port 80.
- Runtime host disk is 40 GiB at about 61% used (~15 GiB free, verified
  2026-05-29). CRM PostgreSQL data, Docker images, and 7-day local backups must
  fit the remaining headroom and keep disk free at or above the 20% warning /
  10% critical thresholds; sustained pressure triggers cleanup or disk
  expansion.
- The 8 GiB memory is shared with the existing Hermes/OpenClaw workloads. The
  CRM stack (multiple Go services + PostgreSQL) must fit the remaining headroom;
  sustained memory pressure, swap use, or OOM is an upgrade trigger.
- Pre-existing public ports `8642` (Hermes) and the OpenClaw gateway are NOT
  owned by CRM. CRM does not manage, expose, or depend on them.

## Container Deployment

Docker Compose manages the v1 runtime. Each service has:

- independent image
- independent environment configuration
- independent health check
- internal Docker network name
- service-specific database credentials
- service-specific logs

The PostgreSQL container is not exposed publicly.

## Network Exposure

Allowed public exposure:

- HTTPS entry through reverse proxy for production
- HTTP only for redirecting to HTTPS in production, or for restricted local
  pre-release diagnostics when explicitly documented
- SSH access controlled by operator policy

Forbidden public exposure:

- PostgreSQL port
- internal service ports
- backup directory
- secrets/config files
- admin/debug endpoints

## Configuration And Secrets

Secrets must not be committed to the repository.

Required secret/config classes:

- database passwords per service
- session/JWT signing secret or equivalent
- service-to-service authentication secret/certificate/key
- backup destination/path and encryption setting where used
- reverse proxy TLS configuration

## Production Endpoint Strategy

Architecture selects this endpoint strategy for ACC-017:

- Pre-release internal validation may use the runtime host public IP
  `118.196.44.193` (`srv-volcengine-sh-01`) with restricted security group rules
  and documented test evidence.
- Production release requires an HTTPS endpoint with a valid TLS certificate.
  If no domain/certificate is available, production release is blocked until a
  domain or otherwise approved HTTPS endpoint is configured.
- HTTP production requests must redirect to HTTPS.
- Login/session cookies must be `Secure`, `HttpOnly`, and `SameSite=Lax`.
- Reverse proxy must set security headers for product traffic:
  `Strict-Transport-Security`, `X-Content-Type-Options`, `Referrer-Policy`, and
  a project-approved `Content-Security-Policy`.
- Admin/debug endpoints must not be publicly exposed.
- ACC-017 release evidence must record endpoint, TLS certificate status,
  security group inbound rules, opened ports, health check URL, deployment
  timestamp, operator, and smoke-test result.

## PostgreSQL Backup

The v1 backup strategy is:

- automatic PostgreSQL backup
- local ECS backup directory
- timestamped backup file per run
- no overwrite of prior backup file
- retain latest 7 days
- automatically delete backups older than 7 days
- log backup success/failure
- encrypt backup files before storage
- restrict backup directory permissions to the ops account/root and backup job
- copy encrypted backup to the off-server backup target `srv-aliyun-bj-01`
  (Alibaba Cloud, Beijing) before production release

The backup schedule is expected to be daily, for example 02:00 local time.

Same-host local backup is a baseline only. The off-server backup target is
defined as `srv-aliyun-bj-01` (a different provider and region from the runtime
host). Production release remains blocked until Infrastructure Ops records
successful encrypted off-server backup copy evidence to that target plus restore
rehearsal evidence, or the user records a formal P0/P1 scope change.

## Restore Requirement

Before production release, Integration Owner and Infrastructure Ops must prove
at least one restore rehearsal:

- restore backup into a test restore target or controlled restore procedure
- verify users/roles, CRM records, history/logs, and service DB permissions
- document command, timestamp, backup file, result, and operator
- document backup checksum, encryption/decryption step, and access subject
- ensure restored restricted data is handled only in the controlled restore
  procedure and is not exposed in product UI, public logs, or repository files

## Production Backup Release Rule

Backups stored only on the same ECS host as the database do not satisfy
production release expectations. ECS disk loss, host compromise, or accidental
deletion could affect database and backups together.

Local backup remains part of v1 operations, but same-host-only backup is a
release blocker. Acceptable production closure uses the defined encrypted
off-server backup target `srv-aliyun-bj-01` (Alibaba Cloud, Beijing), plus
restore rehearsal evidence. The Alibaba host already runs the Alibaba Cloud
backup client (`hbrclient`); whether CRM uses a direct encrypted file copy to
that host or its cloud backup service is an Infrastructure Ops implementation
choice recorded before release.

## OQ-001 Resolution For Architecture

| Item | Decision |
|---|---|
| Production runtime host | Volcengine ECS `srv-volcengine-sh-01` (East China 2, Shanghai, 4 vCPU / 8 GiB), co-located with existing personal workloads |
| Off-server backup target | Alibaba Cloud ECS `srv-aliyun-bj-01` (North China 2, Beijing, 2 vCPU / 2 GiB) |
| Orchestration | Docker Compose |
| Database | Self-hosted PostgreSQL Docker container on the runtime host |
| Backup location | Local backup directory on the runtime host, copied off-server to `srv-aliyun-bj-01` |
| Backup retention | 7 days local; off-server copy before production release |
| Production backup evidence | Encrypted off-server backup copy to `srv-aliyun-bj-01` plus restore rehearsal required before production release |
| Environment ownership | infrastructure-ops owns runtime environment and both registered assets; architecture owns design constraints |
| Endpoint | Pre-release may use runtime host IP `118.196.44.193`; production requires HTTPS endpoint and TLS evidence |
| Domain | Not specified yet; absence of a valid HTTPS endpoint blocks production release |

Domain name remains a deployment configuration item. It does not block recording
the architecture decision, but ACC-017 production evidence must record the
actual HTTPS endpoint and TLS evidence used during release validation.

## Health And Observability

Each service must expose an internal health endpoint covering:

- process up
- database connectivity to owned database/schema
- dependency readiness where appropriate

Logs must include:

- timestamp
- service name
- correlation ID
- actor ID when user-initiated and safe
- error code/category
- dependency call status

G8 tasks must include implementation of health checks and log correlation.

Minimum operational evidence before release:

- ECS CPU, memory, disk, and container health monitoring target.
- Alert target or documented operator notification path.
- Disk free threshold at or above 20% for warning and 10% for critical action.
- Memory pressure and container restart threshold.
- Backup success/failure monitoring.
- Reverse proxy access/error log rotation.
- Service logs with correlation ID.
- Runtime path `/opt/crm-system`.
- Data volume path `/opt/crm-system/volumes/postgres`.
- Backup path `/opt/crm-system/backups/postgres`.
- Logs path `/opt/crm-system/logs`.

The runtime host (`srv-volcengine-sh-01`) provides 4 vCPU / 8 GiB, shared with
existing personal workloads (Hermes, OpenClaw). Production release requires
Infrastructure Ops to record CRM memory and disk usage evidence under load
against the remaining shared headroom, plus an upgrade trigger, for example
sustained memory pressure, swap use, repeated OOM kills, container restart
instability, or disk free dropping below the 10% critical threshold. The
off-server backup target `srv-aliyun-bj-01` (2 vCPU / 2 GiB) is sized for backup
storage only and is not a runtime profile.

## Operator Access

Long-term production operation must not rely on routine root SSH use.

- Root may be used only for initial provisioning or emergency recovery.
- A named deploy/ops user with least required privileges must be created before
  production release.
- SSH access, key ownership, and sudo boundary must be reviewed by Security
  Compliance before G8 implementation tasks are approved.
