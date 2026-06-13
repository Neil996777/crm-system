# CRM CI/CD Migration - G8 Task Package

Status: Ready for Claude G8 handoff audit. Implementation is not started.

This task package converts the G5 design into executable G9-G11 work after
Claude passes the G8 handoff audit. Codex does not self-pass G8 or G12.

## Global Constraints

- Release content commit: `66d2531`.
- Scope: mechanism only.
- Allowed implementation file classes after G8: CI config,
  `docker-compose.prod.yml`, deployment runbooks/scripts, Dockerfiles, release
  evidence templates.
- Forbidden implementation changes: frontend/backend application source,
  `shared`, business logic, API behavior, data-model semantics, enum values, role
  comparison values, zh-CN display behavior, e2e skip/only/slow weakening.
- Production host deployment must not run `git checkout`, `npm run build`,
  `docker build`, `docker compose build`, or `docker compose up --build`.
- Infrastructure Ops is a required reviewer for G8/G11 and supplies the
  ingress, disk, secret, backup, and monitoring review.

## Task DAG

```text
M1
├── M2
├── M3
└── M4
    ├── M5
    └── M6
```

M2, M3, and M4 can be implemented in parallel after M1 is accepted. M5 and M6
consume their outputs during G10/G11.

## Tasks

| ID | Owner | Objective | Primary files after G8 | Verification |
|---|---|---|---|---|
| M1 | Architecture + Infrastructure Ops | Record export/load channel and nginx frontend image decisions. | `docs/architecture/adr/ADR-CICD-001-image-channel-and-frontend-runtime.md`, `docs/product/decision-log.md` | Claude G8 verifies D1/D2/D3 and ACC-CICD traceability. |
| M2 | Infrastructure Ops + QA Execution | Implement off-host CI that runs backend tests, frontend build, e2e, builds 10 Go images + frontend image for `linux/amd64`, labels them, and exports image artifacts. | CI config, Dockerfiles, release bundle scripts | CI record includes tests, image manifest, checksums, labels, digest to commit map. |
| M3 | Infrastructure Ops + Backend Engineer | Convert production Compose to image-only. | `docker-compose.prod.yml` | `grep build:` returns 0; no `latest` fallback; no frontend host build; no source checkout mount for migrations. |
| M4 | Infrastructure Ops | Replace go-live runbook with digest-verified export/load deployment. | deployment runbook/scripts | Static scan contains no host build or `git checkout` path; runbook steps are checksum/load/verify/run/health/rollback. |
| M5 | Integration Owner + QA Execution | Produce release evidence for tests, manifest, deploy transcript, health checks, and rollback point. | `delivery/cicd-release-evidence-template.md`, G11 evidence artifact | Evidence has all five standard §4 items and no secret values. |
| M6 | Integration Owner + Infrastructure Ops | Verify running images map to commit `66d2531` by digest and labels. | deploy verification script, evidence artifact | Running container image IDs/digests and labels match `image-manifest.json`; postgres digest pinned. |

## Detailed Task Contracts

### M1 - ADR And Decision Log

- Acceptance: ACC-CICD-001..008.
- Reference: ADR-CICD-001, `delivery/cicd-migration-acceptance.md`.
- Done when:
  - D1 export/load is selected.
  - D2 frontend nginx image is selected.
  - D3 release content commit `66d2531` is recorded.
  - Registry is rejected for this single-host release with cost/ops rationale.
  - Future registry switch requires ADR supersession.
- No implementation before G8.

### M2 - Off-Host CI Build/Test/Image/Export

- Acceptance: ACC-CICD-001, 002, 004, 006, 008.
- Required tests:
  - 10 Go service test suites.
  - Frontend build.
  - Existing Playwright e2e suite required by the yardstick.
- Required image outputs:
  - 10 Go service images.
  - `frontend-web` nginx image.
  - Upstream postgres digest captured in the manifest.
- Required metadata:
  - Commit `66d2531`.
  - `linux/amd64`.
  - OCI revision label.
  - Image digest/image ID.
  - Image tar SHA-256.
- Failure conditions:
  - Any CI job builds on `srv-volcengine-sh-01`.
  - Any image uses `latest` as release identity.
  - Any image lacks commit label or manifest mapping.

### M3 - Image-Only Compose

- Acceptance: ACC-CICD-001, 003, 005, 007, 008.
- Required service shape:
  - `image:` for all app services and frontend.
  - no `build:`.
  - no `CRM_IMAGE_TAG:-latest`.
  - `gateway-bff` stays loopback-only at `127.0.0.1:8080`.
  - `frontend-web` loopback-only, for example `127.0.0.1:8081`.
  - PostgreSQL pinned to an upstream digest.
  - Migrations read from release artifact path, not Git checkout.
- Preserve:
  - internal Compose network.
  - `read_only`, `tmpfs`, `cap_drop`, `no-new-privileges`.
  - health checks.
  - log mounts.

### M4 - Digest-Pinned Deploy Runbook

- Acceptance: ACC-CICD-001, 003, 004, 006, 007, 008.
- Required deployment steps:
  - receive release bundle.
  - verify bundle and manifest checksums.
  - verify secret path exists without printing values.
  - create pre-migration backup.
  - `docker load` app images.
  - verify loaded image IDs/digests and labels.
  - `docker compose up -d` without `--build`.
  - run migrations from release artifact SQL with database-role passwords
    injected from `prod.env` without printing or command-line exposure.
  - validate Nginx config and reload.
  - run health and negative public-port checks.
  - record full transcript.
- Forbidden commands in production runbook:
  - `git checkout`
  - `npm run build`
  - `docker build`
  - `docker compose build`
  - `docker compose up --build`

### M5 - Release Evidence

- Acceptance: ACC-CICD-002, 004, 006, 008.
- Done when evidence records:
  - CI test result references.
  - digest to commit table for every app image.
  - postgres upstream digest.
  - deploy transcript.
  - health checks.
  - rollback point.
  - no secret values.
  - release migrations contain no fixed development database-role passwords.
- Evidence template: `delivery/cicd-release-evidence-template.md`.

### M6 - Running Digest Verification

- Acceptance: ACC-CICD-004, 006, 007, 008.
- Done when G11 evidence proves:
  - `docker compose ps` expected containers are running.
  - `docker inspect` image IDs match `image-manifest.json`.
  - every app image label has `org.opencontainers.image.revision=66d2531`.
  - no source tree is required to run or rollback.
  - rollback bundle or first-clean-deploy rollback record exists.

## ACC-CICD To Task Map

| ACC | Task coverage |
|---|---|
| ACC-CICD-001 | M2, M3, M4 |
| ACC-CICD-002 | M2, M5 |
| ACC-CICD-003 | M3, M4 |
| ACC-CICD-004 | M1, M2, M4, M5, M6 |
| ACC-CICD-005 | M3 |
| ACC-CICD-006 | M4, M5, M6 |
| ACC-CICD-007 | M3, M4, M6 |
| ACC-CICD-008 | M1, M2, M3, M4, M5, M6 |

## Handoff Condition

After Claude G8 handoff audit passes, Codex may enter G9 implementation for M2,
M3, and M4. G10/G11 must then produce M5/M6 evidence. Claude performs G12
independent audit before any release decision.
