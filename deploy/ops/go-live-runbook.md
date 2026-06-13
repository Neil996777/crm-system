# CRM Go-Live Runbook - Digest-Pinned Load And Run

## Document Control

- Project: CRM System
- Purpose: production deployment of the CI/CD migration release bundle.
- Release content commit: `66d2531`.
- Target host: `srv-volcengine-sh-01`.
- Public endpoint: `118.196.44.193` until Infrastructure Ops confirms otherwise
  in the G11 evidence.
- Executor: Ops `crm-deploy` account on the production host.
- Audit boundary: Codex does not self-approve G10/G11/G12. Claude performs G12
  after the release evidence is complete.

## Hard Preconditions

- Claude G8 handoff audit passed:
  `delivery/cicd-migration-g8-audit-decision.md`.
- The release bundle was produced by CI from release source commit `66d2531`.
- The production host does not build application source. Deployment is
  checksum-verify, image load or digest pull, run, migrate, health-check only.
- Infrastructure Ops has signed off:
  ingress/co-location, public IP, disk retention, secret path, backup/offsite
  path, monitoring, and no impact to vendor agents or Hermes/Feishu bot.
- Secret values are never printed or copied into evidence. Evidence records
  paths only.

## Required Host Paths

| Path | Purpose |
|---|---|
| `/opt/crm-system/incoming/66d2531/` | uploaded release bundle staging path |
| `/opt/crm-system/releases/66d2531/` | verified release runtime directory |
| `/opt/crm-system/secrets/prod.env` | runtime secret env file, values not printed |
| `/opt/crm-system/secrets/backup.passphrase` | backup encryption passphrase file |
| `/opt/crm-system/backups/postgres/` | local encrypted backup output |
| `/opt/crm-system/logs/` | deployment, Nginx, app, and backup logs |

`prod.env` must define the runtime PostgreSQL admin variables, service database
URLs, service-to-service secret, and the production database-role passwords used
by the release migrations. The release bundle must not contain development role
passwords. Secret values must not be printed in the transcript or copied into
evidence.

Required variable names in `prod.env`:

```text
POSTGRES_DB
POSTGRES_USER
POSTGRES_PASSWORD
SERVICE_TOKEN_SECRET
IDENTITY_AUTHZ_DATABASE_URL
LEAD_DATABASE_URL
ACCOUNT_DATABASE_URL
OPPORTUNITY_DATABASE_URL
COMMERCIAL_DATABASE_URL
WORK_DATABASE_URL
AUDIT_HISTORY_DATABASE_URL
REPORTING_DATABASE_URL
IMPORT_EXPORT_DATABASE_URL
CRM_DB_PASSWORD_IDENTITY_AUTHZ
CRM_DB_PASSWORD_LEAD
CRM_DB_PASSWORD_ACCOUNT
CRM_DB_PASSWORD_OPPORTUNITY
CRM_DB_PASSWORD_COMMERCIAL
CRM_DB_PASSWORD_WORK
CRM_DB_PASSWORD_AUDIT_HISTORY
CRM_DB_PASSWORD_REPORTING
CRM_DB_PASSWORD_IMPORT_EXPORT
```

`prod.env` is sourced by Bash during deployment. Quote values that contain shell
metacharacters, especially database URLs containing `&`, for example:

```text
LEAD_DATABASE_URL='postgres://crm_lead_user:<redacted>@postgres:5432/crm_system?sslmode=disable&search_path=lead'
```

## Forbidden Production Actions

All commands in this runbook that execute a deployment step go through
`deploy/scripts/run-release-step.sh`. The wrapper exits before execution if the
operator attempts source checkout, frontend build, Docker image build, Compose
build, or Compose up with build flags on the production host.

## Fourteen-Step Deployment

Run as `crm-deploy` unless a command explicitly uses `sudo`.

### 1. Receive the release bundle

From the build workstation or CI artifact download location, copy these files to
`/opt/crm-system/incoming/66d2531/`:

```bash
release-crm-system-66d2531.tar.gz
release-crm-system-66d2531.tar.gz.sha256
```

### 2. Verify archive checksum and unpack

```bash
set -euo pipefail
mkdir -p /opt/crm-system/incoming/66d2531 /opt/crm-system/releases
cd /opt/crm-system/incoming/66d2531
sha256sum -c release-crm-system-66d2531.tar.gz.sha256
tar -xzf release-crm-system-66d2531.tar.gz
rsync -a --delete release-crm-system-66d2531/ /opt/crm-system/releases/66d2531/
cd /opt/crm-system/releases/66d2531
```

### 3. Start the deploy transcript

```bash
export CRM_DEPLOY_TRANSCRIPT=/opt/crm-system/releases/66d2531/deploy-transcript.log
: > "$CRM_DEPLOY_TRANSCRIPT"
printf '%s release=66d2531 host=srv-volcengine-sh-01\n' "$(date -Is)" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
```

### 4. Verify bundle manifest, checksums, and static release shape

```bash
bash deploy/scripts/verify-release-bundle.sh /opt/crm-system/releases/66d2531
```

This fails if the compose file has `build:` keys, a moving image fallback, a
source checkout migration mount, a missing frontend image service, a missing
image archive, or a checksum mismatch.

### 5. Verify secret path without printing values

```bash
test -f /opt/crm-system/secrets/prod.env
test -f /opt/crm-system/secrets/backup.passphrase
printf '%s secret_path_verified path=/opt/crm-system/secrets/prod.env\n' "$(date -Is)" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
```

Do not print `prod.env`. The migration script validates that all required
variables exist and that forbidden development secret values are not used.

### 6. Record Infrastructure Ops preflight

Infrastructure Ops records these facts in the G11 evidence:

- confirmed public IP / `CRM_SERVER_NAME` value;
- public ingress remains CRM 80/443 plus SSH only;
- `gateway-bff` stays loopback `127.0.0.1:8080`;
- `frontend-web` stays loopback `127.0.0.1:8081`;
- PostgreSQL remains internal only;
- disk has room for current bundle plus one previous-good bundle;
- backup/offsite/monitoring paths are ready;
- vendor agents and Hermes/Feishu bot are untouched.

If the public IP differs from `118.196.44.193`, do not modify the verified
release bundle. Export `CRM_SERVER_NAME=<confirmed-value>` before steps 12 and
13 and record the override in evidence.

### 7. Capture rollback point and pre-migration backup

```bash
cd /opt/crm-system/releases/66d2531
if [[ -d /opt/crm-system/releases/current-good ]]; then
  readlink -f /opt/crm-system/releases/current-good | tee rollback-point.txt
else
  printf 'first-standard-deploy-no-prior-crm-runtime\n' | tee rollback-point.txt
fi

export CRM_RELEASE_DIR=/opt/crm-system/releases/66d2531
export CRM_DEPLOY_SECRET_ENV=/opt/crm-system/secrets/prod.env
if backup_file="$(bash deploy/backup/backup.sh)"; then
  bash deploy/backup/offsite-copy.sh "$backup_file"
  printf '%s backup_file=%s\n' "$(date -Is)" "$backup_file" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
else
  printf 'first-standard-deploy-no-prior-database\n' | tee backup-skipped.txt
  printf '%s backup_skipped=first-standard-deploy-no-prior-database\n' "$(date -Is)" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
fi
```

If backup is skipped because PostgreSQL is not yet running for the first clean
deployment, continue only after Infrastructure Ops records that fact in the G11
evidence.

### 8. Load app images and pull the pinned PostgreSQL runtime

```bash
bash deploy/scripts/load-images.sh /opt/crm-system/releases/66d2531
```

### 9. Verify loaded image IDs and commit labels

```bash
bash deploy/scripts/verify-loaded-images.sh /opt/crm-system/releases/66d2531
```

This fails if any app image ID differs from `.env.release` or if any app image
does not carry `org.opencontainers.image.revision=66d2531`.

### 10. Run image-only Compose

```bash
bash deploy/scripts/compose-up-release.sh /opt/crm-system/releases/66d2531
```

The wrapper runs `docker compose ... up -d` with the release env file and the
secret env loaded in the shell. It does not build.

### 11. Run migrations from release artifacts

```bash
bash deploy/scripts/migrate-release-artifacts.sh /opt/crm-system/releases/66d2531 up
```

Migration SQL is read from `/opt/crm-system/releases/66d2531/migrations/`.
Service database-role passwords are parameterized in the release migration
artifacts and injected from `prod.env` through a temporary local psql input
file. The transcript records the command and temporary path only; secret values
are not passed on the command line and are not printed. No source checkout is
used.

### 12. Apply host Nginx config and reload

```bash
bash deploy/scripts/apply-nginx-runtime-config.sh /opt/crm-system/releases/66d2531
```

The generated host Nginx config keeps 80/443 ownership scoped to the CRM
`server_name`, proxies API traffic to `127.0.0.1:8080`, and proxies the SPA to
`127.0.0.1:8081`.

### 13. Run health checks and negative port checks

```bash
bash deploy/scripts/health-check-release.sh /opt/crm-system/releases/66d2531
bash deploy/healthcheck/check_endpoint.sh
```

From an off-host network path, record both expected failures:

```bash
nc -zv 118.196.44.193 8080
nc -zv 118.196.44.193 5432
```

### 14. Freeze G11 evidence and mark current-good only after G11 review

Copy these into the G11 evidence package:

- CI test logs and Playwright artifacts from `test-results/`;
- `image-manifest.json`, `image-manifest.sha256`, `.env.release`;
- `/opt/crm-system/releases/66d2531/deploy-transcript.log`;
- `rollback-point.txt`, backup filename, backup checksum, offsite copy result;
- `docker compose ps`, health-check output, endpoint smoke output, negative
  public-port check output;
- Infrastructure Ops signoff for ingress, disk, secret, backup, monitoring, and
  public IP.

After Integration Owner confirms the evidence is complete, Infrastructure Ops
may update the previous-good pointer:

```bash
ln -sfn /opt/crm-system/releases/66d2531 /opt/crm-system/releases/current-good
```

## Rollback

Rollback also uses image-only release artifacts. Do not build on the production
host.

1. Stop the current compose project from the failed release directory.
2. Restore the database from the backup recorded in step 7, following
   `deploy/backup/restore-rehearsal.md`.
3. Load the previous-good bundle named in `rollback-point.txt`.
4. Run its `verify-release-bundle.sh`, `load-images.sh`,
   `verify-loaded-images.sh`, `compose-up-release.sh`, and health checks.
5. Record the rollback transcript and backup restore evidence for Claude G12.
