#!/usr/bin/env bash
set -euo pipefail

ENDPOINT="${CRM_ENDPOINT:-https://118.196.44.193}"
HTTP_ENDPOINT="${CRM_HTTP_ENDPOINT:-http://118.196.44.193}"

tmp_headers="$(mktemp)"
trap 'rm -f "$tmp_headers"' EXIT

echo "TEST-DEPLOY-SMOKE-001 checking HTTPS endpoint ${ENDPOINT}"
curl --fail --silent --show-error --location --max-time 15 --output /dev/null --dump-header "$tmp_headers" "${ENDPOINT}/health"

grep -Eiq '^Strict-Transport-Security:' "$tmp_headers"
grep -Eiq '^X-Content-Type-Options:[[:space:]]*nosniff' "$tmp_headers"
grep -Eiq '^Referrer-Policy:' "$tmp_headers"
grep -Eiq '^Content-Security-Policy:' "$tmp_headers"

echo "TEST-DEPLOY-SMOKE-001 passed"

echo "TEST-DEPLOY-SMOKE-002 checking HTTP redirects to HTTPS"
redirect="$(curl --silent --show-error --max-time 15 --output /dev/null --write-out '%{http_code} %{redirect_url}' "${HTTP_ENDPOINT}/health")"
case "$redirect" in
  "301 ${ENDPOINT}/health"|"308 ${ENDPOINT}/health")
    echo "TEST-DEPLOY-SMOKE-002 passed"
    ;;
  *)
    echo "TEST-DEPLOY-SMOKE-002 failed: expected HTTP redirect to ${ENDPOINT}/health, got ${redirect}" >&2
    exit 1
    ;;
esac
