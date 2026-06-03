#!/usr/bin/env bash
set -euo pipefail

ROOT_DIR="$(cd "$(dirname "${BASH_SOURCE[0]}")/.." && pwd)"
ENDPOINT="${CRM_ENDPOINT:-https://118.196.44.193}"

require_file() {
  local path="$1"
  if [[ ! -f "${ROOT_DIR}/${path}" ]]; then
    echo "TEST-DEPLOY-SMOKE failed: missing ${path}" >&2
    exit 1
  fi
}

require_text() {
  local path="$1"
  local pattern="$2"
  if ! grep -Eq "${pattern}" "${ROOT_DIR}/${path}"; then
    echo "TEST-DEPLOY-SMOKE failed: ${path} missing pattern ${pattern}" >&2
    exit 1
  fi
}

require_file "docker-compose.prod.yml"
require_file "deploy/nginx/crm.conf"
require_file "deploy/healthcheck/check_endpoint.sh"
require_file "deploy/monitoring/thresholds.md"
require_file "deploy/ops/operator-access.md"
require_file "docs/release/acc-017-evidence-template.md"

require_text "docker-compose.prod.yml" "/opt/crm-system/volumes/postgres"
require_text "docker-compose.prod.yml" "/opt/crm-system/logs"
require_text "docker-compose.prod.yml" "restart: unless-stopped"
require_text "docker-compose.prod.yml" "127\\.0\\.0\\.1:8080:8080"
require_text "docker-compose.prod.yml" "internal: true"

require_text "deploy/nginx/crm.conf" "listen 80"
require_text "deploy/nginx/crm.conf" 'return 301 https://\$host\$request_uri'
require_text "deploy/nginx/crm.conf" "listen 443 ssl"
require_text "deploy/nginx/crm.conf" "ssl_certificate /etc/letsencrypt/live/118\\.196\\.44\\.193/fullchain\\.pem"
require_text "deploy/nginx/crm.conf" "Strict-Transport-Security"
require_text "deploy/nginx/crm.conf" "X-Content-Type-Options"
require_text "deploy/nginx/crm.conf" "Referrer-Policy"
require_text "deploy/nginx/crm.conf" "Content-Security-Policy"
require_text "deploy/nginx/crm.conf" "proxy_pass http://127\\.0\\.0\\.1:8080"

require_text "deploy/healthcheck/check_endpoint.sh" "TEST-DEPLOY-SMOKE-001"
require_text "deploy/healthcheck/check_endpoint.sh" "TEST-DEPLOY-SMOKE-002"
require_text "deploy/monitoring/thresholds.md" "20%"
require_text "deploy/monitoring/thresholds.md" "10%"
require_text "deploy/ops/operator-access.md" "root"
require_text "docs/release/acc-017-evidence-template.md" "${ENDPOINT}"

echo "TEST-DEPLOY-SMOKE static checks passed for ${ENDPOINT}"
