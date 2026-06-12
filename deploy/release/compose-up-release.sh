#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
secret_env="${CRM_DEPLOY_SECRET_ENV:-/opt/crm-system/secrets/prod.env}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"
[[ -f "$compose_file" ]] || die "compose file missing: $compose_file"
[[ -f "$secret_env" ]] || die "secret env missing: $secret_env"

set -a
. "$env_file"
. "$secret_env"
set +a

export CRM_DEPLOY_TRANSCRIPT="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"

"$script_dir/run-release-step.sh" \
  docker compose \
  --env-file "$env_file" \
  -f "$compose_file" \
  up -d
