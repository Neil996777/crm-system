#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
direction="${2:-up}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
migrations_dir="$release_dir/migrations"
secret_env="${CRM_DEPLOY_SECRET_ENV:-/opt/crm-system/secrets/prod.env}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

case "$direction" in
  up|down) ;;
  *) die "usage: $0 <release-dir> up|down" ;;
esac

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"
[[ -f "$compose_file" ]] || die "compose file missing: $compose_file"
[[ -d "$migrations_dir" ]] || die "migrations directory missing: $migrations_dir"
[[ -f "$secret_env" ]] || die "secret env missing: $secret_env"

set -a
. "$env_file"
. "$secret_env"
set +a

export CRM_DEPLOY_TRANSCRIPT="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"

if [[ "$direction" == "up" ]]; then
  mapfile -t files < <(find "$migrations_dir" -path '*.up.sql' -type f | sort)
else
  mapfile -t files < <(find "$migrations_dir" -path '*.down.sql' -type f | sort -r)
fi

[[ "${#files[@]}" -gt 0 ]] || die "no migration files found for direction: $direction"

for file in "${files[@]}"; do
  rel="${file#$migrations_dir/}"
  container_file="/release-migrations/$rel"
  "$script_dir/run-release-step.sh" \
    docker compose \
    --env-file "$env_file" \
    -f "$compose_file" \
    exec -T postgres \
    psql -v ON_ERROR_STOP=1 \
      -U "${POSTGRES_USER:?POSTGRES_USER required}" \
      -d "${POSTGRES_DB:?POSTGRES_DB required}" \
      -f "$container_file"
done
