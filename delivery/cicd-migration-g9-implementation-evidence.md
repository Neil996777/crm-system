# CRM CI/CD Migration - G9 Implementation Evidence

Status: G9 implementation return for G10/G11 execution. This document does not
claim G10, G11, or G12 pass status.

Date: 2026-06-12

Release content commit: `66d2531`

## Implemented Scope

| Task | Implementation evidence |
|---|---|
| M2 off-host CI | `.github/workflows/cicd-release.yml` checks out a clean `66d2531` worktree, runs 10 Go service tests, frontend build, builds the 10 Go runtime images before e2e, starts a CI-only Docker Compose backend stack (Postgres + migrated schemas + 10 Go services + loopback gateway), runs existing Playwright e2e with workers 2 / retries 1 against the real backend, then builds `frontend-web`, exports image tars, writes `image-manifest.json`, checksums, migrations, compose, scripts, and evidence stubs. |
| M2 image metadata | `deploy/docker/go-service.Dockerfile`, `deploy/docker/frontend-web.Dockerfile`, and `deploy/release/build-release-bundle.sh` apply `org.opencontainers.image.revision=66d2531`, `com.crm.release.content=66d2531`, `com.crm.service`, source, and created labels. |
| M3 image-only compose | `docker-compose.prod.yml` uses required `.env.release` image variables for all 10 Go services plus `frontend-web`; contains no `build:` key and no moving image fallback; Postgres uses `CRM_IMAGE_POSTGRES` digest reference; migrations mount from `/opt/crm-system/releases/66d2531/migrations`. |
| M4 load-and-run deployment | `deploy/ops/go-live-runbook.md` is a 14-step checksum/load/verify/run/migrate/nginx/health/rollback runbook. Deployment scripts live under `deploy/release/` and use `run-release-step.sh` to fail fast on source checkout or host build attempts. |
| Frontend nginx runtime | `deploy/docker/frontend-web.Dockerfile`, `docker-compose.prod.yml`, `deploy/nginx/crm.conf`, and `deploy/release/apply-nginx-runtime-config.sh` route the SPA through loopback `127.0.0.1:8081`. |
| Release evidence generation | `deploy/release/build-release-bundle.sh` copies `delivery/cicd-release-evidence-template.md` into the bundle and creates `evidence/cicd-release-evidence-66d2531.md` with CI artifact pointers. G11 must fill deploy transcript, health checks, rollback point, and Infrastructure Ops signoff. |

## BLK-CICD-G9-001 Follow-Up

| Finding | Implementation response |
|---|---|
| CI e2e had no backend stack | `.github/workflows/cicd-release.yml` now builds the 10 Go service images before Playwright, generates a temporary CI-only compose file, starts Postgres, applies all `*.up.sql` migrations from the checked-out `66d2531` release source, starts `gateway-bff` plus all 9 downstream Go services, waits for gateway `/health` and an admin `/auth/sign-in` probe, and only then runs `npm run test:e2e -- --workers=2 --retries=1`. |
| `TEST-PERSISTENCE` restarts `lead` with `docker compose restart lead` | The workflow writes `COMPOSE_FILE` and `COMPOSE_PROJECT_NAME` into `$GITHUB_ENV` before Playwright runs, so the e2e process and its child `docker compose restart lead` command target the same CI backend stack. |
| e2e diagnostics | The workflow collects `compose-ps.txt` and `compose-logs.txt` into `frontend/test-results/backend-stack/` before cleanup; the release bundle copies these under `test-results/e2e/backend-stack/`. |
| Status | Implementation returned for review. This does not claim G10/G11 pass; a real GitHub Actions run is still required to fill the CI run URL/logs and prove `61 passed` in the actual CI environment. |

## Static Verification Performed

| Check | Result |
|---|---|
| Shell syntax | PASS: `bash -n deploy/release/*.sh deploy/backup/backup.sh deploy/backup/offsite-copy.sh deploy/healthcheck/check_endpoint.sh` |
| Compose config parse | PASS: `docker compose -f docker-compose.prod.yml config --quiet` with dummy release env and dummy secrets |
| Whitespace diff | PASS: `git diff --check` |
| Application source diff | PASS: `git diff --name-only -- frontend services shared` returned empty |
| e2e weakening scan | PASS: no `test.skip`, `test.only`, `.skip(`, or `.only(` found under `frontend/e2e` |
| CI workflow parse | PASS: `.github/workflows/cicd-release.yml` parses as YAML; extracted backend-stack heredoc is unindented correctly for shell execution |
| Compose host-build scan | PASS: no `build:` key in `docker-compose.prod.yml` |
| Moving-tag/source-dependence scan | PASS: no `CRM_IMAGE_TAG`, `latest`, or `/opt/crm-system/current` in `docker-compose.prod.yml`, `deploy/ops/go-live-runbook.md`, `deploy/nginx/crm.conf`, or `deploy/backup/backup.sh` |

## Pending G10/G11 Evidence

The following evidence must be produced by CI and production deployment before
Claude G12:

| Evidence item | Producing step |
|---|---|
| 10 Go service test results | GitHub Actions job `Run 10 Go service test suites` |
| Frontend build result | GitHub Actions job `Frontend build` |
| Existing Playwright e2e result | GitHub Actions job `Run existing Playwright e2e` |
| digest/imageID/tar SHA-256 mapping | `deploy/release/build-release-bundle.sh` output `image-manifest.json` and `.env.release` |
| Deploy transcript | `/opt/crm-system/releases/66d2531/deploy-transcript.log` from the 14-step runbook |
| Health checks | `deploy/release/health-check-release.sh` and `deploy/healthcheck/check_endpoint.sh` outputs |
| Rollback point | `rollback-point.txt`, backup filename/checksum, offsite copy evidence |
| Infrastructure Ops signoff | G11 evidence entry for ingress, public IP, disk, secret path, backup/offsite, monitoring, and co-location impact |

Codex does not self-pass G10, G11, or G12.
