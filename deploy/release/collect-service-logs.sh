#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
out_dir="${2:-$release_dir/evidence/service-logs}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
deploy_start_file="$release_dir/deploy-start.iso"

services=(
  postgres
  gateway-bff
  identity-authz
  lead
  account
  opportunity
  commercial
  work
  audit-history
  reporting
  import-export
  frontend-web
)

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"
[[ -f "$compose_file" ]] || die "compose file missing: $compose_file"
[[ -f "$deploy_start_file" ]] || die "deploy start marker missing: $deploy_start_file"

set -a
. "$env_file"
set +a

export CRM_DEPLOY_TRANSCRIPT="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"
since="$(cat "$deploy_start_file")"
[[ -n "$since" ]] || die "deploy start marker is empty: $deploy_start_file"

mkdir -p "$out_dir"

for service in "${services[@]}"; do
  "$script_dir/run-release-step.sh" \
    docker compose \
    --env-file "$env_file" \
    -f "$compose_file" \
    logs --no-color --timestamps --since "$since" "$service" \
    > "$out_dir/$service.log"
done

scan_file="$out_dir/service-log-scan.txt"
if grep -HnE '28P01|password authentication failed|Bad Gateway|status=502|[[:space:]]502[[:space:]]' "$out_dir"/*.log \
  > "$scan_file"; then
  printf 'service_log_scan=FAIL since=%s evidence=%s\n' "$since" "$scan_file" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
  exit 1
fi

printf 'service_log_scan=PASS since=%s evidence=%s\n' "$since" "$out_dir" | tee "$scan_file" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
