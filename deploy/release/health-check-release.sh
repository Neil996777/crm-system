#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"
[[ -f "$compose_file" ]] || die "compose file missing: $compose_file"

server_name_override="${CRM_SERVER_NAME:-}"
set -a
. "$env_file"
set +a
if [[ -n "$server_name_override" ]]; then
  CRM_SERVER_NAME="$server_name_override"
fi

export CRM_DEPLOY_TRANSCRIPT="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"

"$script_dir/run-release-step.sh" docker compose --env-file "$env_file" -f "$compose_file" ps
"$script_dir/run-release-step.sh" curl -fsS http://127.0.0.1:8080/health
"$script_dir/run-release-step.sh" curl -fsS http://127.0.0.1:8081/health
"$script_dir/run-release-step.sh" curl -kfsS "https://${CRM_SERVER_NAME:?CRM_SERVER_NAME required}/health"

printf '\nNegative public-port checks must be run from an off-host network path and recorded in G11 evidence:\n'
printf '  nc -zv %s 8080  # expected failure\n' "$CRM_SERVER_NAME"
printf '  nc -zv %s 5432  # expected failure\n' "$CRM_SERVER_NAME"
