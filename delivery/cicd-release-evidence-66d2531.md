# CRM CI/CD Release Evidence - 66d2531

Status: G11 rework evidence return for Claude G12 audit. Codex does not
self-pass G10, G11, or G12.

## Document Control

| Field | Evidence |
|---|---|
| Release content commit | `66d2531` |
| Runbook/script commit used | `f4aee22` |
| CI run | `27461747030`, success, `https://github.com/Neil996777/crm-system/actions/runs/27461747030` |
| Artifact | `release-crm-system-66d2531` |
| Local artifact path | `/private/tmp/crm-g11-redeploy-27461747030-new/release-crm-system-66d2531.tar.gz` |
| Production release path | `/opt/crm-system/releases/66d2531` |
| Target host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Deploy window | `2026-06-13T17:32:53+08:00` to `2026-06-13T18:19:22+08:00` |
| Executor | root SSH for privileged host operations; Docker release steps run as `crm-deploy` where possible |
| Reproducibility claim | New CI bundle deployed on `srv-volcengine-sh-01`; no host script, compose, SQL, image, or migration edits; deployment used submitted bundle scripts/runbook commands and transcripted runtime data reset |
| G10/G11 status | Evidence produced and returned for audit; no Codex gate pass claimed |
| G12 auditor | Claude independent audit |

## CI Test Results

| Suite | Result | Evidence |
|---|---|---|
| 10 Go service tests | PASS | bundle `test-results/backend/*.log` |
| Frontend build | PASS | bundle `test-results/frontend-build/build.log` |
| Playwright backend stack | PASS | bundle `test-results/e2e/backend-stack/compose-ps.txt`; 11 backend containers healthy in CI |
| Playwright e2e | PASS | CI run `27461747030`; workflow echoed `61 total (57 stable + 4 isolated)` and uploaded artifact `7609238700` |
| Release static checks | PASS | `verify-release-bundle.sh`; migrations contain 9 `__CRM_DB_PASSWORD_*__` tokens, 0 `_dev_password`, 0 `:'CRM_DB_PASSWORD_'` |

## Bundle And Secret-Safe Scan

| Check | Result | Evidence |
|---|---|---|
| Tarball checksum | PASS | local `sha256sum -c release-crm-system-66d2531.tar.gz.sha256` |
| Bundle verification on host | PASS | `delivery/cicd-g11-deploy-transcript-66d2531.log` |
| Migration token count | PASS | 9 expected `__CRM_DB_PASSWORD_*__` password literals in release migrations |
| Development password scan | PASS | 0 `_dev_password`; 0 psql `:'CRM_DB_PASSWORD_'` variables |
| Secret path | PASS | `/opt/crm-system/secrets/prod.env` and `/opt/crm-system/secrets/backup.passphrase`; values not printed |
| Transcript secret scan | PASS | 20 secret values checked in shell memory only; values not printed |

## Digest And Running Image Evidence

`srv-volcengine-sh-01` uses Docker with containerd snapshotter:
`driver-type=io.containerd.snapshotter.v1`. The submitted
`verify-loaded-images.sh` from commit `f4aee22` was run directly on the host and
passed using the store-independent archive/config/loaded-config digest method.

| Evidence | Path |
|---|---|
| Bundle manifest | `/opt/crm-system/releases/66d2531/image-manifest.json` |
| Containerd-safe archive/config verification | `delivery/cicd-g11-containerd-safe-image-verify-66d2531.tsv`; remote `evidence/containerd-safe-image-verify-redeploy.tsv` |
| Running container image/label table | `delivery/cicd-g11-running-image-inspect-66d2531.tsv`; remote `evidence/running-image-inspect-redeploy.tsv` |

All 11 app image archives matched manifest `archiveSha256`; all archive config
digests matched manifest `imageId`; all loaded app images have
`org.opencontainers.image.revision=66d2531` and
`com.crm.release.content=66d2531`. PostgreSQL is pinned to
`postgres:16-alpine@sha256:16bc17c64a573ef34162af9298258d1aec548232985b33ed7b1eac33ba35c229`.

## Deploy Transcript

| Item | Evidence |
|---|---|
| Full transcript | `delivery/cicd-g11-deploy-transcript-66d2531.log`; remote `/opt/crm-system/releases/66d2531/deploy-transcript.log` |
| Bundle transfer | `rsync` of CI tar + sha to `/opt/crm-system/incoming/66d2531/` |
| Bundle verification | `verify-release-bundle.sh` PASS in transcript |
| Image load | `load-images.sh` loaded 11 image tar files and pulled pinned PostgreSQL |
| Image verification | submitted `verify-loaded-images.sh` PASS on containerd host |
| Compose | submitted `compose-up-release.sh`; image-only `docker compose up -d`, no `--build` |
| Runtime data reset | `run-release-step.sh` stopped stack; old bind data moved to `/opt/crm-system/volumes/postgres.redeploy-20260613T1737-g11`; fresh `/opt/crm-system/volumes/postgres` created |
| Migrations | submitted `migrate-release-artifacts.sh up`; release SQL rendered to 0600 temp stdin; fresh run created schemas/tables/roles |
| Rendered SQL cleanup | PASS: `/tmp/tmp.Gm8LE49goc` absent after migration |
| DB role login | PASS: 9 service roles logged in using secret-safe `PGPASSFILE` stdin; see `delivery/cicd-g11-db-role-login-66d2531.tsv` |
| Nginx | submitted `apply-nginx-runtime-config.sh`; `nginx -t` PASS; reload PASS |
| Endpoint smoke | submitted `deploy/healthcheck/check_endpoint.sh` PASS |
| Forbidden command scan | PASS: no host `npm run build`, `docker build`, compose build, `up --build`, or source checkout |

## Runtime Notes Disclosed For Audit

| Note | Evidence |
|---|---|
| `down -v` did not clear Postgres data because compose uses bind mount `/opt/crm-system/volumes/postgres` | transcript shows first migration saw existing objects, then a guarded runtime data reset moved the bind directory and reran compose/migrate |
| Runtime data reset was not a script/compose/SQL adaptation | `run-release-step.sh` logged `mv` and `install` only; old data retained at `/opt/crm-system/volumes/postgres.redeploy-20260613T1737-g11` |
| `crm-deploy` has Docker group but no passwordless sudo | transcript shows sudo refusal for privileged data move; root then ran the same submitted guard script for host-owned data directory operations |
| Offsite backup key is missing on the Volcengine host | local encrypted backup succeeded; `offsite-copy.sh` failed with missing `/opt/crm-system/secrets/aliyun-bj-ecs-root.pem`; this is an infrastructure follow-up, not an app or bundle adaptation |

## Health, TLS, Headers, Ports

| Check | Result | Evidence |
|---|---|---|
| Compose health | PASS: 12/12 containers healthy | transcript `docker compose ps` |
| Local gateway health | PASS | `curl -fsS http://127.0.0.1:8080/health` in transcript |
| Local frontend health | PASS | `curl -fsS http://127.0.0.1:8081/health` in transcript |
| HTTPS `/health` | PASS, HTTP/2 200 | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| HTTP redirect | PASS | `deploy/healthcheck/check_endpoint.sh` in transcript |
| Security headers | PASS: HSTS, X-Content-Type-Options, Referrer-Policy, CSP | off-host evidence |
| TLS certificate | PASS; issuer `Let's Encrypt YE2`; valid `2026-06-13T03:02:25Z` to `2026-06-19T19:02:24Z` | `delivery/cicd-g11-tls-listeners-66d2531.txt` |
| Public 8080 | CLOSED from `srv-aliyun-bj-01` off-host path (`/dev/tcp` timeout status 124) | off-host evidence |
| Public 5432 | CLOSED from `srv-aliyun-bj-01` off-host path (`/dev/tcp` timeout status 124) | off-host evidence |
| Host listeners | Public `80/443/22`; CRM `8080/8081` loopback only; PostgreSQL internal only | `delivery/cicd-g11-tls-listeners-66d2531.txt` |
| Frontend tmpfs | PASS: `/var/cache/nginx` mounted `uid=101,gid=101,mode=0770`; `client_temp` exists; container healthy | `delivery/cicd-g11-frontend-tmpfs-66d2531.txt` |

## Rollback And Backup

| Field | Evidence |
|---|---|
| Previous-good image manifest | `first-standard-deploy-no-prior-crm-runtime` |
| Previous-good bundle checksum | n/a for first standard deploy pointer |
| Pre-reset encrypted backup | `/opt/crm-system/backups/postgres/crm-postgres-20260613T093420Z.sql.gz.enc` |
| Backup SHA256 | `6d7c4429f9a2936d0ad53a781ea8a4b509f0745df059ef5e7a99dbfdf7b609a1` |
| Offsite copy | FAILED: missing `/opt/crm-system/secrets/aliyun-bj-ecs-root.pem` on target host |
| Runtime data rollback | Old Postgres bind data retained at `/opt/crm-system/volumes/postgres.redeploy-20260613T1737-g11` |
| Restore path | bundle `deploy/backup/restore-rehearsal.md`; no live restore rehearsal claimed in this evidence |

`/opt/crm-system/releases/current-good` was not updated; Codex is not self-passing
G11.

## Infrastructure Ops Signoff Evidence

| Item | Evidence |
|---|---|
| Public IP | Confirmed `118.196.44.193`; SSH used current workstation bind address `192.168.0.100` because infra register value `192.168.43.246` was stale |
| Ingress/co-location | Host Nginx listens on `80/443` for CRM `server_name`; CRM gateway/frontend are loopback only |
| Security group / public exposure | Off-host confirms `80/443` reachable, `8080/5432` closed |
| Disk | Root disk `40G`, used `13G`, available `26G`, `33%` after deploy |
| Secret path | `/opt/crm-system/secrets/prod.env`, `/opt/crm-system/secrets/backup.passphrase`, both `0600`; values not recorded |
| Backup/offsite | Local encrypted backup generated; offsite key missing and recorded as follow-up |
| Monitoring | Health endpoints and TLS/header checks recorded; no alerting system change made |
| Vendor agents | OpenClaw remains loopback (`18789`, `34499`); cloud/vendor agents not stopped |
| Hermes/Feishu bot | Root user `hermes-gateway.service` remained active during preflight; no deployment step touched it |

## ACC-CICD Evidence Map

| ACC | Evidence reference | Status for Claude G12 audit |
|---|---|---|
| ACC-CICD-001 | CI build + transcript no-build scan | Evidence returned; no Codex pass |
| ACC-CICD-002 | run `27461747030` + bundle test-results | Evidence returned; no Codex pass |
| ACC-CICD-003 | transcript load/run/migrate only, no host build | Evidence returned; no Codex pass |
| ACC-CICD-004 | manifest + containerd-safe TSV + running image labels | Evidence returned; no Codex pass |
| ACC-CICD-005 | bundle verification + compose runtime + frontend tmpfs evidence | Evidence returned; no Codex pass |
| ACC-CICD-006 | transcript, health, backup, rollback point; offsite copy gap disclosed | Evidence returned; no Codex pass |
| ACC-CICD-007 | runtime uses images/config/secrets; no source checkout dependency | Evidence returned; no Codex pass |
| ACC-CICD-008 | no app source changes during deployment; runtime data reset disclosed above | Evidence returned; no Codex pass |
