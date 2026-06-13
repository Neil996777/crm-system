# CRM CI/CD Release Evidence - 66d2531

Status: G11 evidence return for Claude G12 audit. Codex does not self-pass G10,
G11, or G12.

## Document Control

| Field | Evidence |
|---|---|
| Release content commit | `66d2531` |
| CI run | `27453655699`, success, `https://github.com/Neil996777/crm-system/actions/runs/27453655699` |
| Artifact | `release-crm-system-66d2531` |
| Local artifact path | `/private/tmp/crm-g11-run27453655699/release-crm-system-66d2531.tar.gz` |
| Production release path | `/opt/crm-system/releases/66d2531` |
| Target host | `srv-volcengine-sh-01` / `118.196.44.193` |
| Deploy window | `2026-06-13T11:42:07+08:00` to `2026-06-13T12:10:19+08:00` |
| Executor | root SSH with non-privileged Docker steps run via `sudo -u crm-deploy`; `crm-deploy` has no SSH key or passwordless sudo |
| G10/G11 status | Evidence produced and returned for audit; no Codex gate pass claimed |
| G12 auditor | Claude independent audit |

## CI Test Results

| Suite | Result | Evidence |
|---|---|---|
| 10 Go service tests | PASS | bundle `test-results/backend/*.log` |
| Frontend build | PASS | bundle `test-results/frontend-build/build.log` (`built in 2.20s`) |
| Playwright backend stack | PASS | bundle `test-results/e2e/backend-stack/compose-ps.txt`; 11 backend containers healthy in CI |
| Playwright e2e | PASS | bundle `test-results/e2e/.last-run.json` = `{"status":"passed","failedTests":[]}`; run 27453655699 recorded full 61 by prior audit |
| Release static checks | PASS | transcript bundle verification; no `build:`, no `latest`, 0 `_dev_password`, 9 `CRM_DB_PASSWORD_*` placeholders |

## Digest And Running Image Evidence

Because `srv-volcengine-sh-01` uses Docker 29.1.3 with `Storage Driver:
overlayfs` / `driver-type: io.containerd.snapshotter.v1`, Docker store IDs differ
from the bundle manifest config digests. The original `verify-loaded-images.sh`
failed on `.Id` as expected for this store type; the deployment recorded a
containerd-safe verification using archive SHA256 + config blob SHA256 + loaded
image labels.

| Evidence | Path |
|---|---|
| Bundle manifest | `/opt/crm-system/releases/66d2531/image-manifest.json` |
| Containerd-safe archive/config verification | `delivery/cicd-g11-containerd-safe-image-verify-66d2531.tsv` and remote `evidence/containerd-safe-image-verify.tsv` |
| Running container image/label table | `delivery/cicd-g11-running-image-inspect-66d2531.tsv` and remote `evidence/running-image-inspect.tsv` |

All 11 app image archives matched manifest `archiveSha256`; all 11 archive config
digests matched manifest `imageId`; all 11 loaded images have
`org.opencontainers.image.revision=66d2531` and
`com.crm.release.content=66d2531`. PostgreSQL is pinned to
`postgres:16-alpine@sha256:16bc17c64a573ef34162af9298258d1aec548232985b33ed7b1eac33ba35c229`.

## Deploy Transcript

| Item | Evidence |
|---|---|
| Full transcript | `delivery/cicd-g11-deploy-transcript-66d2531.log`; remote `/opt/crm-system/releases/66d2531/deploy-transcript.log` |
| Bundle transfer | `rsync` of tar + sha to `/opt/crm-system/incoming/66d2531/` |
| Bundle verification | `verify-release-bundle.sh` PASS in transcript |
| Secret path check | `/opt/crm-system/secrets/prod.env` and `/opt/crm-system/secrets/backup.passphrase`; values not printed |
| Pre-deploy rollback | `first-standard-deploy-no-prior-crm-runtime`; pre-DB backup skipped because no prior database existed |
| Image load | 11 image tar files loaded; pinned PostgreSQL pulled |
| Compose | image-only `docker compose up -d`, no `--build` |
| Migrations | release SQL from `/opt/crm-system/releases/66d2531/migrations/`; secrets injected through 0600 temp stdin |
| Nginx | `/etc/nginx/conf.d/crm.conf`, `nginx -t` PASS, reload PASS |
| Post-deploy backup | encrypted backup `/opt/crm-system/backups/postgres/crm-postgres-20260613T040202Z.sql.gz.enc` |
| Offsite backup | encrypted file copied to `srv-aliyun-bj-01:/opt/crm-system/backups/postgres/` and sha256 verified |

Transcript scans:

| Scan | Result |
|---|---|
| Secret value scan | PASS, 20 secret values checked against transcript |
| Forbidden command scan | PASS, no `npm run build`, `docker build`, `docker compose build`, `docker compose up --build`, or `git checkout` command |

## Runtime Reworks Recorded During Deploy

These are explicitly disclosed for Claude audit:

| Rework | Reason | Scope |
|---|---|---|
| Containerd-safe image verification | Host uses containerd snapshotter; `.Id` comparison mismatched manifest config digest | Verification method only; images unchanged |
| Rendered 0600 migration stdin | `psql` variables are not substituted inside dollar-quoted `DO $$` blocks | Migration runner behavior only; release SQL source unchanged; secret values not in command line/transcript |
| Frontend tmpfs permission options | `frontend-web` runs as UID/GID 101 and could not create `/var/cache/nginx/client_temp` on root-owned tmpfs | Production compose runtime options only; image unchanged; no build |
| TLS renewal | Existing short-lived cert had expired on 2026-06-09; renewal initially timed out due certbot random sleep | Renewed with `--no-random-sleep-on-renew`; Nginx reloaded |

## Health, TLS, Headers, Ports

| Check | Result | Evidence |
|---|---|---|
| Compose health | 12/12 containers healthy | transcript final `docker compose ps` |
| Local gateway health | PASS | `curl -fsS http://127.0.0.1:8080/health` in transcript |
| Local frontend health | PASS | `curl -fsS http://127.0.0.1:8081/health` in transcript |
| HTTPS `/health` | PASS, HTTP/2 200 | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| HTTP redirect | PASS, HTTP 301 to `https://118.196.44.193/health` | same |
| Security headers | PASS: HSTS, X-Content-Type-Options, Referrer-Policy, CSP | same |
| TLS certificate | PASS after renewal; issuer `Let's Encrypt YE2`; valid `2026-06-13T03:02:25Z` to `2026-06-19T19:02:24Z` | transcript and `openssl x509` output |
| Public 8080 | CLOSED from `srv-aliyun-bj-01` off-host path | `delivery/cicd-g11-offhost-health-and-negative-ports-66d2531.txt` |
| Public 5432 | CLOSED from `srv-aliyun-bj-01` off-host path | same |
| Host listeners | Public `80/443/22`; CRM `8080/8081` loopback only; PostgreSQL internal only | transcript final `ss -lntp` |

## Rollback Point

| Field | Evidence |
|---|---|
| Previous-good image manifest | `first-standard-deploy-no-prior-crm-runtime` |
| Previous-good bundle checksum | n/a for first clean standard deploy |
| Pre-deploy database backup | `first-standard-deploy-no-prior-database` because `/opt/crm-system` and Postgres runtime were absent |
| Post-deploy encrypted backup | `/opt/crm-system/backups/postgres/crm-postgres-20260613T040202Z.sql.gz.enc` |
| Backup SHA256 | `0e214854fe9138bd94f9cecc3f34a4bdf44b0d8d64d653c406fa6e6b8c988f28` |
| Offsite copy | `srv-aliyun-bj-01:/opt/crm-system/backups/postgres/crm-postgres-20260613T040202Z.sql.gz.enc` |
| Restore path | bundle `deploy/backup/restore-rehearsal.md`; no live restore rehearsal claimed in this evidence |

`/opt/crm-system/releases/current-good` was not updated; Codex is not self-passing
G11.

## Infrastructure Ops Signoff Evidence

| Item | Evidence |
|---|---|
| Public IP | Confirmed `118.196.44.193` |
| Ingress/co-location | Host Nginx listens on `80/443` for CRM `server_name`; CRM gateway/frontend are loopback only |
| Security group / public exposure | Off-host confirms `80/443` reachable, `8080/5432` closed |
| Disk | Root disk `40G`, used `12G`, available `26G`, `32%` after deploy |
| Secret path | `/opt/crm-system/secrets/prod.env`, `/opt/crm-system/secrets/backup.passphrase`, both `0600`; values not recorded |
| Backup/offsite | Encrypted backup generated and copied off-server to `srv-aliyun-bj-01` |
| Monitoring | Health endpoints and TLS/header checks recorded; no alerting system change made |
| Vendor agents | OpenClaw remains loopback (`18789`, `34499`); cloud/vendor agents not stopped |
| Hermes/Feishu bot | Root user `hermes-gateway.service` remained active during preflight; no deployment step touched it |

## ACC-CICD Evidence Map

| ACC | Evidence reference | Status for Claude G12 audit |
|---|---|---|
| ACC-CICD-001 | CI build + transcript no-build scan | Evidence returned; no Codex pass |
| ACC-CICD-002 | run 27453655699 + bundle test-results | Evidence returned; no Codex pass |
| ACC-CICD-003 | transcript load/run/migrate only | Evidence returned; no Codex pass |
| ACC-CICD-004 | manifest + containerd-safe TSV + running image labels | Evidence returned; no Codex pass |
| ACC-CICD-005 | bundle verification + compose runtime | Evidence returned; no Codex pass |
| ACC-CICD-006 | transcript, health, backup/offsite, rollback point | Evidence returned; no Codex pass |
| ACC-CICD-007 | runtime uses images/config/secrets; no source checkout dependency | Evidence returned; no Codex pass |
| ACC-CICD-008 | no app source changes during deployment; runtime reworks disclosed above | Evidence returned; no Codex pass |
