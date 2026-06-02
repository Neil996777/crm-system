#!/usr/bin/env bash
set -euo pipefail

services=(gateway-bff identity-authz lead account opportunity commercial work audit-history reporting import-export)

for svc in "${services[@]}"; do
  if rg -n 'crm-system/services/[^"]+/internal' "services/${svc}" | rg -v "crm-system/services/${svc}/internal"; then
    echo "cross-service internal imports are forbidden" >&2
    exit 1
  fi
done

if rg -n 'DATABASE_URL' services/gateway-bff | rg -v 'should-not-be-used'; then
  echo "gateway-bff must not have database credentials" >&2
  exit 1
fi

gateway_compose_block="$(awk '/^  gateway-bff:/{flag=1; next} /^  [a-z-]+:/{flag=0} flag {print}' docker-compose.yml)"
if printf '%s\n' "$gateway_compose_block" | rg -n 'DATABASE_URL'; then
  echo "gateway-bff must not have database credentials" >&2
  exit 1
fi

for svc in "${services[@]}"; do
  test -f "services/${svc}/go.mod"
  test -f "services/${svc}/cmd/server/main.go"
done

for svc in identity-authz lead account opportunity commercial work audit-history reporting import-export; do
  test -f "services/${svc}/migrations/0001_init_schema.up.sql"
  test -f "services/${svc}/migrations/0001_init_schema.down.sql"
done

echo "TASK-001 static checks passed"
