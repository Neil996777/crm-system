# CRM CI/CD Release Architecture - G5 Design Package

Status: G5 design and G8 task-package input. No implementation before Claude G8
handoff audit passes.

Authority:

- `delivery/cicd-migration-acceptance.md`
- `standards/cicd-and-release-standard.md`
- `company/operating-model.md`
- `company/infrastructure/*` registers
- `docs/architecture/adr/ADR-CICD-001-image-channel-and-frontend-runtime.md`

## Scope Boundary

This package changes the release mechanism only. The G9 implementation may edit
only CI configuration, `docker-compose.prod.yml`, deployment runbooks/scripts,
Dockerfiles, and release evidence templates. Application source for the frontend
and backend, shared contracts, business logic, APIs, data model semantics,
zh-CN display rules, enum values, and role comparison values remain unchanged.

Release content commit is `66d2531`.

## M2 CI Design

CI runs off the production host and owns build, test, image creation, and export.

### Required Stages

| Stage | Design |
|---|---|
| Checkout | Fetch exactly commit `66d2531`; fail if the working tree is dirty or the commit is not the release content commit. |
| Backend test | Run all 10 Go modules with backend tests, using the existing test strategy and real PostgreSQL where the suite requires it. |
| Frontend test/build | Run `npm ci`, `npm run build`, and the QA-defined Playwright suite from `frontend/` with the current workers/retries policy. Yardstick expects the existing 61 e2e tests plus backend tests. |
| Image build | Build 10 Go service images plus `frontend-web` for `linux/amd64`; no production host build. |
| Label | Apply the OCI labels listed in ADR-CICD-001, including `org.opencontainers.image.revision=66d2531`. |
| Export | Save every app image as an image artifact and record image digest/image ID plus tar SHA-256. |
| Release bundle | Create a release bundle with image tars, `image-manifest.json`, image-only Compose, deployment scripts/runbook, migration SQL release artifacts, checksums, and the evidence template. |
| Retention | Retain CI logs, test output, image manifest, and release bundle checksums as G11/G12 evidence. |

### Required Image Build Inputs

- Existing Go service Dockerfile pattern stays under `deploy/docker/`.
- A frontend nginx Dockerfile is added during G9 only after G8 passes.
- Base images must be pinned by digest where practical during G9:
  - Go builder base.
  - Alpine runtime base.
  - Node builder base for frontend.
  - nginx runtime base for frontend.
  - PostgreSQL official runtime image in Compose.

### CI Outputs

```text
release-crm-system-66d2531/
  images/
    gateway-bff-66d2531.tar
    identity-authz-66d2531.tar
    lead-66d2531.tar
    account-66d2531.tar
    opportunity-66d2531.tar
    commercial-66d2531.tar
    work-66d2531.tar
    audit-history-66d2531.tar
    reporting-66d2531.tar
    import-export-66d2531.tar
    frontend-web-66d2531.tar
  image-manifest.json
  image-manifest.sha256
  docker-compose.prod.yml
  deploy/
    runbook.md
    scripts/
      verify-release-bundle.sh
      load-images.sh
      verify-loaded-images.sh
      migrate-release-artifacts.sh
      apply-nginx-runtime-config.sh
  migrations/
    <service>/*.up.sql
    <service>/*.down.sql
    migration-manifest.sha256
  evidence/
    cicd-release-evidence-template.md
  test-results/
    backend/
    frontend-build/
    e2e/
```

The migration SQL in this bundle is a release artifact derived from commit
`66d2531`; it is not a Git checkout and is not used for building on the
production host.

## M3 Image-Only `docker-compose.prod.yml` Design

The production Compose file becomes image-only. It contains no `build:` key,
does not build the frontend, and does not mount a source checkout.

### Compose Image Variables

Compose reads immutable local tag handles from a release `.env.release` file:

```text
CRM_RELEASE_COMMIT=66d2531
CRM_IMAGE_GATEWAY_BFF=crm-system/gateway-bff:66d2531
CRM_IMAGE_ID_GATEWAY_BFF=sha256:<recorded-image-id>
CRM_IMAGE_FRONTEND_WEB=crm-system/frontend-web:66d2531
CRM_IMAGE_ID_FRONTEND_WEB=sha256:<recorded-image-id>
CRM_IMAGE_POSTGRES=postgres:16-alpine@sha256:<recorded-upstream-digest>
```

`docker-compose.prod.yml` uses `image: ${CRM_IMAGE_<SERVICE>:?}` for every
application runtime. The deploy script verifies each local tag resolves to the
matching `CRM_IMAGE_ID_*` before running Compose.

### Service Shape

| Service | Runtime change |
|---|---|
| `gateway-bff` | `image: ${CRM_IMAGE_GATEWAY_BFF:?}`; keep `127.0.0.1:8080:8080`; no `build:`. |
| `identity-authz` | `image: ${CRM_IMAGE_IDENTITY_AUTHZ:?}`; no `build:`. |
| `lead` | `image: ${CRM_IMAGE_LEAD:?}`; no `build:`. |
| `account` | `image: ${CRM_IMAGE_ACCOUNT:?}`; no `build:`. |
| `opportunity` | `image: ${CRM_IMAGE_OPPORTUNITY:?}`; no `build:`. |
| `commercial` | `image: ${CRM_IMAGE_COMMERCIAL:?}`; no `build:`. |
| `work` | `image: ${CRM_IMAGE_WORK:?}`; no `build:`. |
| `audit-history` | `image: ${CRM_IMAGE_AUDIT_HISTORY:?}`; no `build:`. |
| `reporting` | `image: ${CRM_IMAGE_REPORTING:?}`; no `build:`. |
| `import-export` | `image: ${CRM_IMAGE_IMPORT_EXPORT:?}`; no `build:`. |
| `frontend-web` | New nginx runtime image, `image: ${CRM_IMAGE_FRONTEND_WEB:?}`, loopback-only publish such as `127.0.0.1:8081:80`; no source mount. |
| `postgres` | `image: ${CRM_IMAGE_POSTGRES:?}` pinned to upstream digest; internal network only. |

### Required Differences From Current Compose

| Current `docker-compose.prod.yml` behavior | New design | ACC mapping |
|---|---|---|
| 10 Go services each contain `build:`. | Remove every `build:` key. | ACC-CICD-001, ACC-CICD-003, ACC-CICD-005 |
| Images use `${CRM_IMAGE_TAG:-latest}`. | Images use required immutable commit tags from `.env.release`; release identity is digest in manifest. No `latest` fallback. | ACC-CICD-004, ACC-CICD-005 |
| Frontend is not a Compose image; host Nginx serves `/opt/crm-system/current/frontend/dist`. | Add `frontend-web` nginx image built by CI and loopback-proxied by host Nginx. | ACC-CICD-001, ACC-CICD-002, ACC-CICD-005 |
| PostgreSQL init mounts SQL from `/opt/crm-system/current/services/...`. | Mount SQL from `/opt/crm-system/releases/66d2531/migrations/...`, a checksum-verified release artifact, not a Git checkout. | ACC-CICD-003, ACC-CICD-007 |
| Rollback requires `git checkout <prev>` and `up --build`. | Rollback loads previous-good image bundle and runs image-only Compose. | ACC-CICD-003, ACC-CICD-006, ACC-CICD-007 |
| Gateway loopback bind, internal network, `read_only`, `cap_drop`, `no-new-privileges`, health checks, logs. | Preserve these settings. | ACC-CICD-005, C5 |

## M4 Digest-Pinned Deploy Runbook Design

The G9 runbook replaces host checkout/build with load-and-run:

1. Copy the release bundle to `/opt/crm-system/incoming/66d2531/`.
2. Verify bundle checksum and manifest checksum.
3. Verify no secret values are present in bundle artifacts.
4. Confirm host runtime secrets exist at the recorded secret path without
   printing values.
5. Back up current production data before migration.
6. Load app image tars with `docker load`.
7. Verify each loaded image tag resolves to the digest/image ID in
   `image-manifest.json`.
8. Confirm image labels include `org.opencontainers.image.revision=66d2531`.
9. Install or refresh the runtime release directory:
   `/opt/crm-system/releases/66d2531/`.
10. Run image-only Compose:
    `docker compose -f docker-compose.prod.yml --env-file .env.release up -d`
    with no `--build`.
11. Run migrations from release artifact SQL only; do not read from Git.
12. Apply host Nginx runtime config through a deployment script after `nginx -t`;
    keep public 80/443 scoped to the CRM `server_name`.
13. Run health checks and negative public-port checks.
14. Capture deploy transcript, health output, image inspect output, and rollback
    point evidence.

The runbook must fail fast on any use of `git checkout`, `npm run build`,
`docker build`, `docker compose build`, or `docker compose up --build` on the
production host.

## M5 Release Evidence Design

Release evidence must contain the five required standard §4 items:

| Evidence item | Required source |
|---|---|
| Test results | CI backend tests, frontend build, and e2e run tied to commit `66d2531`. |
| Digest to commit | `image-manifest.json`, image inspect output, image labels. |
| Deploy transcript | Full `scp`/checksum/load/verify/compose/migrate/nginx/health output. |
| Health checks | Compose health, HTTPS `/health`, HTTP redirect, security headers, TLS, negative `8080`/`5432` public checks. |
| Rollback point | Previous-good image manifest and backup/restore path. For first clean deployment, the rollback point is "no prior CRM runtime" plus the pre-deploy database backup state; after the first standard deploy, keep the previous-good bundle. |

The template is `delivery/cicd-release-evidence-template.md`.

## M6 Commit Traceability Verification

G11 must prove:

- Every running app container image has the tag from the manifest.
- Every running app container image ID/digest matches the manifest.
- Every app image label includes `org.opencontainers.image.revision=66d2531`.
- The release evidence digest table maps each app image to `66d2531`.
- `postgres` is pinned to an upstream digest and is not a moving tag.
- No running container was built on the production host during deployment.
- No release step depends on `/opt/crm-system/current/.git` or source checkout.

## Infrastructure Ops Review View

| Review area | Design requirement |
|---|---|
| Ingress / co-location | Only CRM 80/443 through host Nginx for the CRM `server_name`; no takeover of unrelated host ingress. |
| Public ports | SSH 22 remains infrastructure-owned; CRM gateway and frontend container remain loopback-only; PostgreSQL remains internal. |
| Disk | Keep current bundle + one previous-good bundle; prune older loaded images only after rollback evidence is captured. |
| Secrets | Secret values never enter Git, release bundle, image labels, or evidence. Evidence records only path and recovery procedure. |
| Backup / restore | Pre-deploy backup before migration; off-server copy and restore rehearsal are G11/G12 evidence. |
| Monitoring | Compose health, endpoint smoke, TLS expiry, disk, Nginx logs, app logs, and container health checks recorded. |
| Existing services | Do not stop, remove, or reconfigure vendor agents or the Hermes/Feishu bot. |

## ACC-CICD Traceability

| ACC | Delivery coverage |
|---|---|
| ACC-CICD-001 | M2 off-host CI build; M3 no `build:`; M4 no host build commands. |
| ACC-CICD-002 | M2 CI test/build/image/export outputs; M5 test evidence. |
| ACC-CICD-003 | M4 checksum verify + `docker load` + `docker compose up -d`; no `git checkout` for build. |
| ACC-CICD-004 | ADR tag/digest scheme; M2 labels; M6 running-image digest verification. |
| ACC-CICD-005 | M3 image-only Compose diff table and service shape. |
| ACC-CICD-006 | M5 evidence template and runbook transcript requirements. |
| ACC-CICD-007 | M3/M4 release artifact migrations; no source checkout or build cache dependency. |
| ACC-CICD-008 | C1 scope boundary; release content commit `66d2531`; no app source edits. |
