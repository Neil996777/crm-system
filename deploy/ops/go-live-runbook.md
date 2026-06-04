# CRM Go-Live Runbook (Production Deployment)

## Document Control

- Project: CRM System
- Purpose: production go-live after G12 audit PASS (2026-06-04)
- Executor: **Ops `crm-deploy` account** on `srv-volcengine-sh-01` (`118.196.44.193`).
  NOT executed by the audit platform (Claude). Claude verifies before (static smoke) and
  after (read-only live checks) only.
- Release commit: the audited HEAD (G12-passed). Confirm `git log` shows the
  "Record G12 PASS" commit or later before deploying.

## Why a fresh deploy is required

The runtime evidence recorded at G11 predates the six rounds of G12 rework (IDOR fixes,
durable identity-authz audit, optimistic concurrency, idempotency, etc.). The version that
was running on the host is therefore **stale and missing the audited security fixes**. A
read-only probe on 2026-06-05 found the production endpoint not serving (HTTPS 443 handshake
failed, HTTP empty reply), so go-live is a fresh deployment of the audited build.

## Pre-flight (already done / must be true)

- [x] Local static deploy smoke passed (`bash scripts/test_deploy_smoke.sh`).
- [ ] Audited release commit pushed to the origin the server pulls from.
- [ ] Server `/opt/crm-system/current/.env` exists with all required secrets and is NOT in
      the repo: `POSTGRES_DB`, `POSTGRES_USER`, `POSTGRES_PASSWORD`, per-service DB
      credentials, session/JWT signing secret, service-to-service (S2S) secret. (Confirm with
      Ops; never print secret values.)
- [ ] SSH as `crm-deploy` (key only; password auth disabled). Keys live under
      `/Users/neil/practice/software/.secrets/ssh-keys/` — never commit them.
- [ ] Security group `sg-366ptx1bxp9ts1e710babmc8y`: public TCP 22/80/443 only (verified).

## Deploy steps (run on the host as `crm-deploy`)

```bash
cd /opt/crm-system/current

# 1) Update the deploy dir to the audited release commit, confirm HEAD
git fetch --all
git checkout <AUDITED_RELEASE_COMMIT>          # the G12-PASS commit or later
git log --oneline -1

# 2) Back up the current production DB BEFORE migrations (rollback safety)
bash deploy/backup/backup.sh                    # encrypted, timestamped, local
bash deploy/backup/offsite-copy.sh              # copy to off-site srv-aliyun-bj-01
#    record the backup filename + checksum for rollback

# 3) Build the frontend SPA (Nginx serves /opt/crm-system/current/frontend/dist)
cd frontend && npm ci && npm run build && cd ..

# 4) Build + start the backend stack from the audited source
docker compose -f docker-compose.prod.yml up -d --build

# 5) Wait for postgres healthy, then apply migrations
docker compose -f docker-compose.prod.yml ps    # postgres + services healthy
bash scripts/migrate.sh up                       # applies all */migrations/*.up.sql

# 6) Reload Nginx (validate config first)
sudo nginx -t && sudo systemctl reload nginx

# 7) Confirm backup + cert renewal timers are active
systemctl is-active crm-backup.timer certbot-renew.timer
```

## Post-deploy verification

On the host (Ops):
```bash
docker compose -f docker-compose.prod.yml ps              # all services healthy
bash deploy/healthcheck/check_endpoint.sh                 # TEST-DEPLOY-SMOKE-001/002
```

External (Claude can run these read-only):
```bash
curl -sS https://118.196.44.193/health                    # expect HTTP 200
curl -sI  http://118.196.44.193                           # expect 301 -> https
# negative: from off-host edge (srv-aliyun-bj-01), confirm 8080 + 5432 NOT reachable publicly
```

Record ACC-017 release evidence: endpoint, TLS status, security-group rules, opened ports,
health URL, deploy timestamp, operator, smoke result
(`docs/release/acc-017-evidence-template.md`).

## Rollback

If migrations or smoke fail:
```bash
docker compose -f docker-compose.prod.yml down
# restore the DB from the step-2 backup (decrypt + psql restore per deploy/backup/restore-rehearsal.md)
git checkout <PREVIOUS_STABLE_COMMIT>
docker compose -f docker-compose.prod.yml up -d --build
```
Backup + restore rehearsal was proven during G12 (TASK-040), so rollback is exercised.

## Notes

- gateway-bff is the only published port and is bound to `127.0.0.1:8080` (loopback);
  it is reachable only through the host Nginx reverse proxy. Postgres + the other 9 Go
  services are on the internal Docker network and are never published.
- Do not take over host port 80/443 ownership beyond the CRM `server_name`; the host runs
  unrelated personal workloads (co-location constraint).
