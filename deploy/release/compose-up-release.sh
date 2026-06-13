#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
mode="${2:-all}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
secret_env="${CRM_DEPLOY_SECRET_ENV:-/opt/crm-system/secrets/prod.env}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"
deploy_start_file="$release_dir/deploy-start.iso"

app_services=(
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
[[ -f "$secret_env" ]] || die "secret env missing: $secret_env"

case "$mode" in
  postgres|apps|all) ;;
  *) die "usage: $0 <release-dir> [postgres|apps|all]" ;;
esac

set -a
. "$env_file"
. "$secret_env"
set +a

export CRM_DEPLOY_TRANSCRIPT="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"

compose() {
  "$script_dir/run-release-step.sh" \
    docker compose \
    --env-file "$env_file" \
    -f "$compose_file" \
    "$@"
}

record_deploy_start() {
  local started_at

  started_at="${CRM_DEPLOY_STARTED_AT:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"
  printf '%s\n' "$started_at" > "$deploy_start_file"
  printf '%s deploy_started_at=%s\n' "$(date -Is)" "$started_at" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
}

wait_for_postgres() {
  local attempts="${CRM_POSTGRES_WAIT_ATTEMPTS:-60}"
  local delay_seconds="${CRM_POSTGRES_WAIT_DELAY_SECONDS:-2}"
  local attempt

  for ((attempt = 1; attempt <= attempts; attempt++)); do
    if compose exec -T postgres pg_isready \
      -U "${POSTGRES_USER:?POSTGRES_USER required}" \
      -d "${POSTGRES_DB:?POSTGRES_DB required}"; then
      printf '%s postgres_ready attempt=%s\n' "$(date -Is)" "$attempt" | tee -a "$CRM_DEPLOY_TRANSCRIPT"
      return 0
    fi
    sleep "$delay_seconds"
  done

  die "postgres did not become ready after $attempts attempts"
}

start_postgres() {
  record_deploy_start
  compose up -d postgres
  wait_for_postgres
}

start_apps() {
  [[ -f "$deploy_start_file" ]] || die "deploy start marker missing: $deploy_start_file"
  compose up -d "${app_services[@]}"
}

case "$mode" in
  postgres)
    start_postgres
    ;;
  apps)
    start_apps
    ;;
  all)
    start_postgres
    "$script_dir/migrate-release-artifacts.sh" "$release_dir" up
    start_apps
    ;;
esac
