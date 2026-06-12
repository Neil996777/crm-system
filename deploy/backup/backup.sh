#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/opt/crm-system/backups/postgres}"
LOG_FILE="${BACKUP_LOG_FILE:-/opt/crm-system/logs/backup/postgres-backup.log}"
LOCK_FILE="${BACKUP_LOCK_FILE:-/opt/crm-system/backups/postgres/.backup.lock}"
RELEASE_DIR="${CRM_RELEASE_DIR:-/opt/crm-system/releases/66d2531}"
ENV_FILE="${CRM_ENV_FILE:-${CRM_DEPLOY_SECRET_ENV:-/opt/crm-system/secrets/prod.env}}"
RELEASE_ENV_FILE="${CRM_RELEASE_ENV_FILE:-$RELEASE_DIR/.env.release}"
COMPOSE_FILE="${CRM_COMPOSE_FILE:-$RELEASE_DIR/docker-compose.prod.yml}"
BACKUP_PASSPHRASE_FILE="${BACKUP_PASSPHRASE_FILE:-/opt/crm-system/secrets/backup.passphrase}"

timestamp="$(date -u +%Y%m%dT%H%M%SZ)"
backup_file="${BACKUP_DIR}/crm-postgres-${timestamp}.sql.gz.enc"
checksum_file="${backup_file}.sha256"

log() {
  mkdir -p "$(dirname "$LOG_FILE")"
  printf '%s %s\n' "$(date -Is)" "$*" >> "$LOG_FILE"
}

fail() {
  log "backup_failed reason=$*"
  echo "backup failed: $*" >&2
  exit 1
}

mkdir -p "$BACKUP_DIR" "$(dirname "$LOG_FILE")"
chmod 700 "$BACKUP_DIR"

exec 9>"$LOCK_FILE"
flock -n 9 || fail "another backup is running"

[[ -f "$ENV_FILE" ]] || fail "missing env file $ENV_FILE"
[[ -f "$RELEASE_ENV_FILE" ]] || fail "missing release env file $RELEASE_ENV_FILE"
[[ -f "$COMPOSE_FILE" ]] || fail "missing compose file $COMPOSE_FILE"
[[ -f "$BACKUP_PASSPHRASE_FILE" ]] || fail "missing BACKUP_PASSPHRASE_FILE $BACKUP_PASSPHRASE_FILE"
[[ ! -e "$backup_file" ]] || fail "backup target already exists $backup_file"

set -a
# shellcheck disable=SC1090
source "$ENV_FILE"
# shellcheck disable=SC1090
source "$RELEASE_ENV_FILE"
set +a

[[ -n "${POSTGRES_USER:-}" ]] || fail "POSTGRES_USER is not set"
[[ -n "${POSTGRES_PASSWORD:-}" ]] || fail "POSTGRES_PASSWORD is not set"

tmp_file="${backup_file}.tmp"
trap 'rm -f "$tmp_file"' EXIT

log "backup_started file=$backup_file release_dir=$RELEASE_DIR"

docker compose --env-file "$RELEASE_ENV_FILE" -f "$COMPOSE_FILE" \
  exec -T -e "PGPASSWORD=${POSTGRES_PASSWORD}" postgres \
  pg_dumpall -U "$POSTGRES_USER" \
  | gzip -c \
  | openssl enc -aes-256-cbc -salt -pbkdf2 -pass "file:${BACKUP_PASSPHRASE_FILE}" -out "$tmp_file"

chmod 600 "$tmp_file"
mv "$tmp_file" "$backup_file"
sha256sum "$backup_file" > "$checksum_file"
chmod 600 "$checksum_file"

find "$BACKUP_DIR" -type f -name 'crm-postgres-*.sql.gz.enc' -mtime +7 -print | while read -r old_backup; do
  rm -f "$old_backup" "${old_backup}.sha256"
  log "backup_pruned file=$old_backup"
done

log "backup_succeeded file=$backup_file checksum_file=$checksum_file"
printf '%s\n' "$backup_file"
