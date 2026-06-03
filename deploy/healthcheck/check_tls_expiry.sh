#!/usr/bin/env bash
set -euo pipefail

CERT_PATH="${CRM_TLS_CERT_PATH:-/etc/letsencrypt/live/118.196.44.193/fullchain.pem}"
WARNING_SECONDS="${CRM_TLS_WARNING_SECONDS:-172800}"

echo "Checking TLS certificate validity for ${CERT_PATH}"
if openssl x509 -in "$CERT_PATH" -checkend "$WARNING_SECONDS" -noout >/dev/null; then
  echo "TLS certificate has at least ${WARNING_SECONDS} seconds remaining."
else
  echo "TLS certificate expires within ${WARNING_SECONDS} seconds." >&2
  exit 1
fi
