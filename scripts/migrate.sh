#!/usr/bin/env bash
set -euo pipefail

direction="${1:-up}"
case "$direction" in
  up|down) ;;
  *) echo "usage: bash scripts/migrate.sh up|down" >&2; exit 2 ;;
esac

if [[ "$direction" == "up" ]]; then
  mapfile -t files < <(find services -path '*/migrations/*.up.sql' -type f | sort)
else
  mapfile -t files < <(find services -path '*/migrations/*.down.sql' -type f | sort -r)
fi

for file in "${files[@]}"; do
  echo "applying $file"
  docker compose exec -T postgres psql -U crm_admin -d crm_system < "$file"
done
