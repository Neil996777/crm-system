#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
EVIDENCE_DIR="${ROOT_DIR}/docs/release/evidence"

require_file() {
  local path="$1"
  if [[ ! -s "${path}" ]]; then
    echo "TEST-RELEASE-EVIDENCE-001 failed: missing or empty ${path#${ROOT_DIR}/}" >&2
    exit 1
  fi
}

require_text() {
  local path="$1"
  local pattern="$2"
  require_file "${path}"
  if ! grep -Eiq "${pattern}" "${path}"; then
    echo "TEST-RELEASE-EVIDENCE-001 failed: ${path#${ROOT_DIR}/} missing pattern ${pattern}" >&2
    exit 1
  fi
}

TLS_CURL="${EVIDENCE_DIR}/tls-curl-https-2026-06-03.txt"
TLS_OPENSSL="${EVIDENCE_DIR}/tls-openssl-2026-06-03.txt"
TLS_REDIRECT="${EVIDENCE_DIR}/tls-curl-redirect-2026-06-03.txt"
CERTBOT_RENEW="${EVIDENCE_DIR}/certbot-renew-dry-run-2026-06-03.txt"
NEGATIVE_PROBES="${EVIDENCE_DIR}/external-negative-probes-2026-06-03.txt"
RESTORE_TRANSCRIPT="${EVIDENCE_DIR}/backup-restore-transcript-2026-06-03.txt"
OPERATOR_ACCESS="${EVIDENCE_DIR}/operator-access-transcript-2026-06-03.txt"

require_text "${TLS_CURL}" 'HTTP/[0-9.]+ 200|HTTP/2 200'
require_text "${TLS_CURL}" 'Strict-Transport-Security'
require_text "${TLS_OPENSSL}" 'Verify return code: 0|subject=.*118\.196\.44\.193|IP Address:118\.196\.44\.193'
require_text "${TLS_REDIRECT}" 'HTTP/[0-9.]+ 301|HTTP/[0-9.]+ 308'
require_text "${TLS_REDIRECT}" 'Location: https://118\.196\.44\.193'
require_text "${CERTBOT_RENEW}" 'Congratulations, all simulated renewals succeeded|dry run'
require_text "${NEGATIVE_PROBES}" '118\.196\.44\.193.*8080'
require_text "${NEGATIVE_PROBES}" '118\.196\.44\.193.*5432'
require_text "${NEGATIVE_PROBES}" 'refused|timed out|succeeded: no|failed'
require_text "${RESTORE_TRANSCRIPT}" 'sha256sum -c|OK'
require_text "${RESTORE_TRANSCRIPT}" 'psql|\\du|schema|audit'
require_text "${OPERATOR_ACCESS}" 'getent passwd crm-deploy crm-ops'
require_text "${OPERATOR_ACCESS}" 'passwordauthentication no'
require_text "${OPERATOR_ACCESS}" 'permitrootlogin prohibit-password|PermitRootLogin prohibit-password'
require_text "${OPERATOR_ACCESS}" 'drwx------|-rw-------'

echo "TEST-RELEASE-EVIDENCE-001 passed"
