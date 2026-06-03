#!/usr/bin/env bash
set -euo pipefail

BACKUP_DIR="${BACKUP_DIR:-/opt/crm-system/backups/postgres}"
LOG_FILE="${BACKUP_LOG_FILE:-/opt/crm-system/logs/backup/postgres-backup.log}"
OFFSITE_HOST="${OFFSITE_HOST:-47.95.119.211}" # srv-aliyun-bj-01
OFFSITE_USER="${OFFSITE_USER:-root}"
OFFSITE_DIR="${OFFSITE_DIR:-/opt/crm-system/backups/postgres}"
OFFSITE_SSH_KEY="${OFFSITE_SSH_KEY:-/opt/crm-system/secrets/aliyun-bj-ecs-root.pem}"
BACKUP_FILE="${1:-}"

log() {
  mkdir -p "$(dirname "$LOG_FILE")"
  printf '%s %s\n' "$(date -Is)" "$*" >> "$LOG_FILE"
}

fail() {
  log "offsite_copy_failed reason=$*"
  echo "offsite copy failed: $*" >&2
  exit 1
}

if [[ -z "$BACKUP_FILE" ]]; then
  BACKUP_FILE="$(find "$BACKUP_DIR" -type f -name 'crm-postgres-*.sql.gz.enc' -print | sort | tail -n 1)"
fi

[[ -n "$BACKUP_FILE" ]] || fail "no encrypted backup found"
[[ -f "$BACKUP_FILE" ]] || fail "backup file not found $BACKUP_FILE"
[[ -f "${BACKUP_FILE}.sha256" ]] || fail "sha256 file not found ${BACKUP_FILE}.sha256"
[[ -f "$OFFSITE_SSH_KEY" ]] || fail "offsite SSH key not found $OFFSITE_SSH_KEY"

log "offsite_copy_started file=$BACKUP_FILE target=${OFFSITE_USER}@${OFFSITE_HOST}:${OFFSITE_DIR}"

ssh -i "$OFFSITE_SSH_KEY" -o StrictHostKeyChecking=accept-new "${OFFSITE_USER}@${OFFSITE_HOST}" \
  "mkdir -p '$OFFSITE_DIR' && chmod 700 '$OFFSITE_DIR'"

if ssh -i "$OFFSITE_SSH_KEY" -o StrictHostKeyChecking=accept-new "${OFFSITE_USER}@${OFFSITE_HOST}" \
  "command -v rsync >/dev/null 2>&1"; then
  rsync -av --chmod=F600,D700 -e "ssh -i $OFFSITE_SSH_KEY -o StrictHostKeyChecking=accept-new" \
    "$BACKUP_FILE" "${BACKUP_FILE}.sha256" "${OFFSITE_USER}@${OFFSITE_HOST}:${OFFSITE_DIR}/"
else
  scp -i "$OFFSITE_SSH_KEY" -o StrictHostKeyChecking=accept-new \
    "$BACKUP_FILE" "${BACKUP_FILE}.sha256" "${OFFSITE_USER}@${OFFSITE_HOST}:${OFFSITE_DIR}/"
  ssh -i "$OFFSITE_SSH_KEY" -o StrictHostKeyChecking=accept-new "${OFFSITE_USER}@${OFFSITE_HOST}" \
    "chmod 600 '$OFFSITE_DIR/$(basename "$BACKUP_FILE")' '$OFFSITE_DIR/$(basename "${BACKUP_FILE}.sha256")'"
fi

ssh -i "$OFFSITE_SSH_KEY" -o StrictHostKeyChecking=accept-new "${OFFSITE_USER}@${OFFSITE_HOST}" \
  "cd '$OFFSITE_DIR' && sha256sum -c '$(basename "${BACKUP_FILE}.sha256")'"

log "offsite_copy_succeeded file=$BACKUP_FILE target=${OFFSITE_USER}@${OFFSITE_HOST}:${OFFSITE_DIR}"
