#!/usr/bin/env bash
set -euo pipefail

require_file() {
  local path="$1"
  if [[ ! -f "$path" ]]; then
    echo "missing required file: $path" >&2
    exit 1
  fi
}

require_executable() {
  local path="$1"
  require_file "$path"
  if [[ ! -x "$path" ]]; then
    echo "required file is not executable: $path" >&2
    exit 1
  fi
}

require_contains() {
  local path="$1"
  local pattern="$2"
  require_file "$path"
  if ! grep -Eq "$pattern" "$path"; then
    echo "required pattern not found in $path: $pattern" >&2
    exit 1
  fi
}

require_executable deploy/backup/backup.sh
require_executable deploy/backup/offsite-copy.sh
require_file deploy/backup/restore-rehearsal.md
require_file deploy/ops/crm-backup.service
require_file deploy/ops/crm-backup.timer
require_file docs/release/acc-017-backup-evidence-template.md

require_contains deploy/backup/backup.sh 'pg_dumpall'
require_contains deploy/backup/backup.sh 'openssl enc .*aes-256-cbc'
require_contains deploy/backup/backup.sh 'BACKUP_PASSPHRASE_FILE'
require_contains deploy/backup/backup.sh 'sha256'
require_contains deploy/backup/backup.sh 'find .*-mtime [+]7'
require_contains deploy/backup/backup.sh 'flock'

require_contains deploy/backup/offsite-copy.sh 'rsync'
require_contains deploy/backup/offsite-copy.sh 'srv-aliyun-bj-01|47\\.95\\.119\\.211'
require_contains deploy/backup/offsite-copy.sh 'sha256'

require_contains deploy/backup/restore-rehearsal.md 'users/roles'
require_contains deploy/backup/restore-rehearsal.md 'history/logs'
require_contains deploy/backup/restore-rehearsal.md 'checksum'
require_contains deploy/backup/restore-rehearsal.md 'decryption'

require_contains deploy/ops/crm-backup.timer 'OnCalendar=.*02:00'
require_contains docs/release/acc-017-backup-evidence-template.md 'Restore rehearsal'

echo "TASK-040 backup artifact static checks passed"
