# ADR-CICD-001: Image Distribution Channel And Frontend Runtime

Architecture Decision Record for the CRM CI/CD migration follow-on change.

| Field | Value |
|---|---|
| ADR ID | ADR-CICD-001 |
| Title | Image distribution channel and frontend runtime |
| Status | Accepted for G5 handoff audit |
| Date | 2026-06-12 |
| Owner | Architecture |
| Required reviewer | Infrastructure Ops |
| Trigger | CI/CD and release-pipeline migration re-entering G5 |
| Related Acceptance ID(s) | ACC-CICD-001..008 |
| Related Services / Gate | SVC-001..010, frontend-web, G5/G8/G11/G12 |

## Context

The CRM production host `srv-volcengine-sh-01` is a single Volcengine ECS host.
Company infrastructure registers state that `/opt/crm-system` was removed on
2026-06-06, Docker has no CRM images or containers, and there is no approved
image registry. The new company CI/CD standard requires off-host build,
CI-owned build/test/image production, CD pull/load-and-run only, digest
traceability, and release evidence.

The release owner confirmed the yardstick on 2026-06-12:

- Release content commit is `66d2531`.
- Image distribution defaults to export/load (`docker save` -> `scp` ->
  `docker load`) unless a registry is explicitly introduced later.
- Frontend runtime should be an nginx container image, not loose `dist` files
  copied to the host.

This ADR decides D1 and D2 for the G5 package. It does not approve release and
does not bypass G8/G10/G11/G12.

## Decision

### D1: Use Export/Load As The Approved Channel

The CRM migration will use the export/load form:

1. Off-host CI/build workstation checks out release commit `66d2531`.
2. CI runs the required QA suite and backend tests.
3. CI builds all application images for `linux/amd64`.
4. CI exports each image as a digest-recorded image artifact.
5. CI transfers the release bundle to `srv-volcengine-sh-01` with `scp`.
6. CD verifies artifact checksums and image digests, then runs `docker load`.
7. CD starts the already-loaded images with image-only Compose.

No registry is introduced in this change. A future switch to a registry requires
a new ADR or an explicit supersession of this ADR.

### D2: Package The Frontend As An nginx Image

The frontend will be built off-host into an nginx runtime image named
`crm-system/frontend-web:<commit>`. The frontend image is part of the same
digest manifest and deployment path as the Go service images.

The production host remains a co-located host with host-level Nginx terminating
80/443 for the CRM `server_name`. The frontend nginx container is bound only to a
loopback port such as `127.0.0.1:8081`; host Nginx proxies the SPA route to that
loopback container and proxies `/api`, `/auth`, `/admin`, and `/health` to the
gateway BFF at `127.0.0.1:8080`.

### Image Set

The release image set contains the 10 Go services plus the nginx frontend image.
The official PostgreSQL image is also digest-pinned as a third-party runtime
dependency in Compose and release evidence, but it is not an application image.

| Runtime | Image name | Source commit mapping |
|---|---|---|
| Go service | `crm-system/gateway-bff:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/identity-authz:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/lead:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/account:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/opportunity:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/commercial:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/work:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/audit-history:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/reporting:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Go service | `crm-system/import-export:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Frontend | `crm-system/frontend-web:66d2531` | `org.opencontainers.image.revision=66d2531` + digest manifest |
| Third-party runtime | `postgres:16-alpine@sha256:<captured>` | Upstream digest recorded; source commit is not applicable |

### Tag And Digest Strategy

- Local tag: `crm-system/<image>:66d2531`.
- Production release identity: image digest recorded in
  `image-manifest.json`, not the tag.
- Transport checksum: SHA-256 of every exported image tar and of the full release
  bundle.
- Labels on every application image:
  - `org.opencontainers.image.revision=66d2531`
  - `org.opencontainers.image.source=<CRM repo URL or local repo identifier>`
  - `org.opencontainers.image.created=<CI timestamp>`
  - `com.crm.release.content=66d2531`
  - `com.crm.service=<service name>`

For export/load, Compose uses immutable local tags as handles because there is
no registry namespace to resolve `image@sha256:<digest>` from. The deploy
runbook still treats the digest as the release identity: before `docker compose
up -d`, it verifies that every loaded local image tag resolves to the exact image
digest recorded in `image-manifest.json`. If any digest or archive checksum
mismatches, deployment stops before any container is started.

### Digest To Commit Mapping

CI produces and retains:

- `image-manifest.json`: service, local tag, image digest/image ID, tar checksum,
  labels, build platform, source commit.
- `image-manifest.sha256`: checksum for the manifest.
- `release-crm-system-66d2531.sha256`: checksum for the full release bundle.
- CI test-result artifacts tied to the same commit and manifest.

The G11 evidence package copies this mapping into the release evidence template.
G12 verifies that the running images on `srv-volcengine-sh-01` match the
manifest, digest, and commit.

## Alternatives Considered

- Registry form. It is standard-compliant and better for multiple hosts, but
  company infrastructure currently has no registry and this release targets one
  host. Creating and operating a new registry would add credentials, retention,
  backup, access control, monitoring, and failure modes not needed for the
  single-host migration.
- Dist plus `scp` frontend. It avoids a frontend container but weakens the
  release unit because the frontend becomes loose files rather than a
  digest-recorded image. It also creates a second deployment path separate from
  the Go services.
- Continue host source checkout and build. Rejected. It violates
  `standards/cicd-and-release-standard.md` §1 and ACC-CICD-001/003/005/007.

## Consequences

- The production host no longer needs Git source checkout, Node, Go build tools,
  Docker buildx, or build cache for normal deployment.
- CI must build `linux/amd64` images off-host.
- The release bundle must include image artifacts, image manifest, migration SQL
  release artifacts, image-only Compose, deploy scripts/runbook, checksums, and
  evidence stubs.
- Disk retention on `srv-volcengine-sh-01` must keep only the current release
  bundle and one previous-good rollback bundle unless Infrastructure Ops approves
  more retention.
- Secrets are never written to the image manifest or evidence. Evidence records
  only secret file locations and recovery methods.

## Infrastructure Ops Review Points

Infrastructure Ops must review before G8/G11:

- No public port beyond CRM 80/443 and SSH is required.
- `gateway-bff` remains loopback-only at `127.0.0.1:8080`.
- `frontend-web` is loopback-only, for example `127.0.0.1:8081`.
- PostgreSQL has no public host port and remains on the internal Compose
  network.
- Vendor agents and the Hermes/Feishu bot are not changed.
- Secret files are restored on the host without recording values in docs.
- Backup, restore rehearsal, monitoring, and TLS evidence are captured at G11.

## Traceability

| Acceptance | ADR coverage |
|---|---|
| ACC-CICD-001 | Selects off-host build and rejects host build. |
| ACC-CICD-002 | Defines CI-built image set and export artifacts. |
| ACC-CICD-003 | Defines CD as checksum verify + `docker load` + run only. |
| ACC-CICD-004 | Defines digest as release identity and labels commit mapping. |
| ACC-CICD-005 | Requires image-only Compose for all app runtimes. |
| ACC-CICD-006 | Defines manifest and release evidence inputs. |
| ACC-CICD-007 | Removes host source checkout and build-cache dependence. |
| ACC-CICD-008 | Pins release content to `66d2531`; mechanism only. |

## Supersession

None.
