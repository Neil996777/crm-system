# CRM CI/CD Release Evidence - 66d2531

Status: G11 rework evidence return for Claude G12 audit. Codex does not
self-pass G10, G11, or G12.

## Document Control

| Field | Evidence |
|---|---|
| Release content commit | `66d2531` |
| Runbook/script commit used | `1ee5dda` |
| CI run | `27464903082`, success, `https://github.com/Neil996777/crm-system/actions/runs/27464903082` |
| Artifact | `release-crm-system-66d2531` |
| Local artifact path | `/private/tmp/crm-g11-redeploy-27464903082/release-crm-system-66d2531.tar.gz` |
| Local artifact SHA-256 | `43f3f414684bad389a8b0c07b7f4c33c1167ed28d547e92fedbea3dded2c2c02` |
| Production release path | `/opt/crm-system/releases/66d2531` |
| Target host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Deploy window | `2026-06-13T19:23:23+08:00` to `2026-06-13T19:31:41+08:00` |
| Startup-log scan window | From `deploy-start.iso` = `2026-06-13T11:24:46Z` |
| Executor | `crm-deploy` SSH for routine Docker release steps; root SSH only for Nginx reload and host-owned Postgres bind-data reset |
| G10/G11 status | Evidence produced and returned for audit; no Codex gate pass claimed |
| G12 auditor | Claude independent audit |

## CI Test Results

| Suite | Result | Evidence |
|---|---|---|
| 10 Go service tests | PASS | CI run `27464903082`; bundle `test-results/backend/*.log` |
| Frontend build | PASS | CI run `27464903082`; bundle `test-results/frontend-build/build.log` |
| Playwright backend stack | PASS | CI run `27464903082`; bundle `test-results/e2e/backend-stack/compose-ps.txt` |
| Playwright e2e | PASS | CI run `27464903082` |
| Release static checks | PASS | CI and host `verify-release-bundle.sh`; no `build:`, no `latest`, no `_dev_password`, no psql `:'CRM_DB_PASSWORD_'` |

## Bundle And Secret-Safe Scan

| Check | Result | Evidence |
|---|---|---|
| Tarball checksum | PASS | Local `shasum -a 256` matched `release-crm-system-66d2531.tar.gz.sha256` |
| Bundle verification on host | PASS | `delivery/cicd-g11-deploy-transcript-66d2531.log` and host command output |
| Migration token count | PASS | 9 expected `__CRM_DB_PASSWORD_*__` password tokens in release migrations |
| Development password scan | PASS | 0 `_dev_password`; 0 psql `:'CRM_DB_PASSWORD_'` variables |
| Secret path | PASS | `/opt/crm-system/secrets/prod.env`, `backup.passphrase`, and `aliyun-bj-ecs-root.pem`; values not printed |
| Transcript/service-log secret scan | PASS | `delivery/cicd-g11-secret-scan-66d2531.txt` checked 20 secret values without printing them |

## Digest And Running Image Evidence

| Evidence | Path |
|---|---|
| Containerd-safe archive/config verification | `delivery/cicd-g11-containerd-safe-image-verify-66d2531.tsv` |
| Running container image/label table | `delivery/cicd-g11-running-image-inspect-66d2531.tsv` |

All 11 app image archives matched manifest `archiveSha256`; all archive config
digests matched manifest `imageId`; all loaded app images have
`org.opencontainers.image.revision=66d2531` and
`com.crm.release.content=66d2531`. PostgreSQL is pinned to
`postgres:16-alpine@sha256:16bc17c64a573ef34162af9298258d1aec548232985b33ed7b1eac33ba35c229`.

## Deploy Transcript And Ordering

| Item | Evidence |
|---|---|
| Full transcript | `delivery/cicd-g11-deploy-transcript-66d2531.log` |
| Bundle transfer and checksum | `rsync` to `/opt/crm-system/incoming/66d2531/`; host `sha256sum -c` PASS |
| Bundle verification | `verify-release-bundle.sh` PASS |
| Pre-reset backup | `delivery/cicd-g11-backup-offsite-66d2531.txt` |
| Offsite copy | PASS to `srv-aliyun-bj-01`; remote sha256 check OK |
| Image load | `load-images.sh` loaded 11 image tar files and pulled pinned PostgreSQL |
| Image verification | Submitted `verify-loaded-images.sh` PASS on containerd host |
| Runtime data reset | Old bind data retained at `/opt/crm-system/volumes/postgres.redeploy-20260613T112438Z-g11-004`; fresh `/opt/crm-system/volumes/postgres` created |
| Startup ordering | Submitted `compose-up-release.sh ... postgres` -> `migrate-release-artifacts.sh up` -> `compose-up-release.sh ... apps` |
| Migrations | Release SQL rendered through 0600 temp stdin and cleaned; rendered tmp dir cleanup PASS |
| DB role login | PASS: 9 service roles logged in using secret-safe temporary PGPASSFILE stdin; `delivery/cicd-g11-db-role-login-66d2531.tsv` |
| Nginx | Submitted `apply-nginx-runtime-config.sh`; `nginx -t` PASS; reload PASS |
| Endpoint smoke | Submitted `deploy/healthcheck/check_endpoint.sh` PASS |

## Startup Race Evidence

| Check | Result | Evidence |
|---|---|---|
| Service container logs from deploy start | Captured for Postgres, 10 Go services, and frontend | `delivery/cicd-g11-service-logs-66d2531/` |
| Service-log race scan | PASS: 0 `28P01`, 0 password-authentication failures, 0 `502` evidence | `delivery/cicd-g11-service-logs-66d2531/service-log-scan.txt` |
| Local grep recheck | PASS: no matches for `28P01`, password auth, Bad Gateway, or 502 in committed service logs | local grep run before commit |
| Nginx 502 scan since deploy start | PASS: no access/error 502 entries since `2026-06-13T19:24:46+08:00` | `delivery/cicd-g11-nginx-502-scan-66d2531.txt` |

## Health, TLS, Headers, Ports

| Check | Result | Evidence |
|---|---|---|
| Compose health | PASS: 12/12 containers healthy | transcript `docker compose ps` |
| Local gateway health | PASS | `curl -fsS http://127.0.0.1:8080/health` in transcript |
| Local frontend health | PASS | `curl -fsS http://127.0.0.1:8081/health` in transcript |
| HTTPS `/health` from off-host | PASS, HTTP/2 200 | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| HTTP redirect | PASS, HTTP 301 to HTTPS | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| Security headers | PASS: HSTS, X-Content-Type-Options, Referrer-Policy, CSP | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| TLS certificate | PASS; Let's Encrypt YE2; valid to `2026-06-19T19:02:24Z` | `delivery/cicd-g11-tls-listeners-66d2531.txt` |
| Public 8080 | CLOSED from `srv-aliyun-bj-01` off-host path, timeout status 124 | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| Public 5432 | CLOSED from `srv-aliyun-bj-01` off-host path, timeout status 124 | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| Host listeners | Public `80/443/22`; CRM `8080/8081` loopback only; PostgreSQL internal only | `delivery/cicd-g11-tls-listeners-66d2531.txt` |
| Frontend tmpfs | PASS: nginx-owned cache/client_temp writable; container healthy | `delivery/cicd-g11-frontend-tmpfs-66d2531.txt` |

## Rollback And Backup

| Field | Evidence |
|---|---|
| Previous-good image manifest | `first-standard-deploy-no-prior-g12-approved-runtime` |
| Previous-good bundle checksum | n/a for first standard deploy pointer |
| Pre-reset encrypted backup | `/opt/crm-system/backups/postgres/crm-postgres-20260613T112341Z.sql.gz.enc` |
| Backup SHA256 | `cf960acefabe379a974789aa52d0e00120fbd42eae2c03af4c99ad596a88a06b` |
| Offsite copy | PASS to `srv-aliyun-bj-01`; checksum OK |
| Runtime data rollback | Old Postgres bind data retained at `/opt/crm-system/volumes/postgres.redeploy-20260613T112438Z-g11-004` |
| Restore path | bundle `deploy/backup/restore-rehearsal.md`; no live restore rehearsal claimed in this evidence |

`/opt/crm-system/releases/current-good` was not updated; Codex is not self-passing
G11 or G12.

## Infrastructure Ops Signoff Evidence

| Item | Evidence |
|---|---|
| Public IP | Confirmed `118.196.44.193`; SSH used current workstation bind address `192.168.0.100` |
| Ingress/co-location | Host Nginx listens on `80/443` for CRM `server_name`; CRM gateway/frontend are loopback only |
| Security group / public exposure | Off-host confirms `80/443` reachable, `8080/5432` closed |
| Disk | Root disk `40G`, used `13G`, available `25G`, `34%` after deploy |
| Secret path | `prod.env`, `backup.passphrase`, and offsite SSH key exist under `/opt/crm-system/secrets`; values not recorded |
| Backup/offsite | Local encrypted backup generated and copied offsite successfully |
| Monitoring | Health endpoints and TLS/header checks recorded; no alerting system change made |
| Vendor agents | OpenClaw remains loopback (`18789`, `34499`); cloud/vendor agents not stopped |
| Hermes/Feishu bot | Root user `hermes-gateway.service` remained active; deployment did not touch it |

## Scope Boundary

No `services/`, `frontend/`, `shared/`, or `.github` tracked diff was part of
the runtime fix commit. `BLK-PROD-AUDIT-001` was not modified; it remains an
application bug outside this CI/CD mechanism-only scope.

## ACC-CICD Evidence Map

| ACC | Evidence reference | Status for Claude G12 audit |
|---|---|---|
| ACC-CICD-001 | CI build + transcript no-build scan + script guard | Evidence returned; no Codex pass |
| ACC-CICD-002 | CI run `27464903082` + bundle test-results | Evidence returned; no Codex pass |
| ACC-CICD-003 | transcript load/run/migrate only; postgres-before-migration-before-app ordering | Evidence returned; no Codex pass |
| ACC-CICD-004 | manifest + containerd-safe TSV + running image labels | Evidence returned; no Codex pass |
| ACC-CICD-005 | bundle verification + image-only compose + frontend tmpfs evidence | Evidence returned; no Codex pass |
| ACC-CICD-006 | transcript, health, backup, offsite copy, rollback point | Evidence returned; no Codex pass |
| ACC-CICD-007 | runtime uses images/config/secrets; no source checkout dependency | Evidence returned; no Codex pass |
| ACC-CICD-008 | no application source changes; `BLK-PROD-AUDIT-001` untouched | Evidence returned; no Codex pass |
