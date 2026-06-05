# Go-Live Kickback #1 — CRM System (Claude → Codex)

## Document Control

- Source: post-deploy verification 2026-06-05
  (`archive/reviews/g12-audit/post-deploy-verification-2026-06-05.md`)
- Decision: **go-live NOT VERIFIED.** `BLK-GOLIVE-004` Open. Do not declare release successful.
- Executor: Codex (Ops `crm-deploy`); then return to Claude for re-verification.

## What's wrong

The production stack is running image `crm-system/<svc>:prod-20260603` (built 2026-06-03) — the
**pre-G12-rework build**. The go-live updated the source tree (`git checkout da9d63c`) and captured
evidence, but never ran `docker compose up -d --build`, `scripts/migrate.sh up`, or the frontend
build. So the audited security fixes (IDOR BLK-G12-015/016, durable audit 017/026, optimistic
concurrency 018, idempotency 019, …) are NOT in production, and the new migrations were never
applied. Smoke passed but only checks online/HTTPS/redirect, which cannot detect stale code.

## Hard constraints

- **Do NOT modify application code.** This is a deploy, not a rework. Deploy `da9d63c` as-is.
- **Back up the production DB before migrating** (you already have
  `crm-postgres-20260604T170022Z.sql.gz.enc`; take a fresh one if any data changed since).
- If `up -d --build` or `migrate.sh up` fails, STOP, roll back per the runbook, and raise a blocker
  — do not work around it or fabricate evidence.
- Do not change security-group / Nginx ownership scope (co-location constraint).

## Step 1 — First, establish the truth (cheap, do this before redeploying)

Run and capture:
```bash
for s in gateway-bff identity-authz lead account opportunity commercial work reporting audit-history import-export; do
  docker image inspect crm-system/$s:prod-20260603 --format "$s {{.Created}}" 2>/dev/null
done
```
- If `Created` is 2026-06-03 → confirmed stale; proceed to Step 2 (real deploy).
- If `Created` is 2026-06-05 → the tag was reused on a rebuild; then prove the code identity another
  way (Step 3 behavioral check) and provide the `up -d --build` log that produced them.

## Step 2 — Actually deploy the audited build

Per `deploy/ops/go-live-runbook.md`, as `crm-deploy` in `/opt/crm-system/current` (already at
`da9d63c`):
```bash
# fresh pre-migration backup if data changed since 20260604T170022Z
bash deploy/backup/backup.sh && bash deploy/backup/offsite-copy.sh

# frontend
cd frontend && npm ci && npm run build && cd ..

# rebuild + restart backend from da9d63c source
docker compose -f docker-compose.prod.yml up -d --build

# apply new migrations (must show 0004_lead_conversion_idempotency etc.)
docker compose -f docker-compose.prod.yml ps          # postgres healthy first
bash scripts/migrate.sh up

# reload nginx
sudo nginx -t && sudo systemctl reload nginx
```
Capture the FULL output of `up -d --build` and `migrate.sh up` into the evidence transcript — these
were missing last time and are the proof the audited code is deployed.

## Step 3 — Evidence to return (all of these)

1. `docker image inspect ... {{.Created}}` for all 9 services (post-rebuild: must be 2026-06-05).
2. Full `up -d --build` + `migrate.sh up` output (migrations applying the new files).
3. Post-rebuild `docker compose ps` (containers freshly recreated, image identity clear).
4. A behavioral proof of a rework-only fix in prod — e.g. an authenticated by-id read of a
   non-owned commercial contract returns safe `404` (BLK-G12-015), or the work `VERSION_CONFLICT`
   path (BLK-G12-018). Pick one and capture request/response.
5. `bash deploy/healthcheck/check_endpoint.sh` (TEST-DEPLOY-SMOKE-001/002) re-run.
6. Confirm `certbot-renew.timer` will renew before the cert expires 2026-06-09 (e.g.
   `certbot renew --dry-run` and timer next-run).

## Definition of Done

- Running production images provably built from `da9d63c`; new migrations applied; frontend
  rebuilt; smoke re-passed; one behavioral rework-fix proof captured; cert renewal confirmed.
- `planning/blockers.md` BLK-GOLIVE-004 set to Resolved with the cited evidence; ACC-017 evidence
  updated; commit made.
- Return to Claude for independent re-verification. Do NOT declare go-live successful yourself.
