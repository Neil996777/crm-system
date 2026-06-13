#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
direction="${2:-up}"
env_file="$release_dir/.env.release"
compose_file="$release_dir/docker-compose.prod.yml"
migrations_dir="$release_dir/migrations"
secret_env="${CRM_DEPLOY_SECRET_ENV:-/opt/crm-system/secrets/prod.env}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

db_role_secret_vars=(
  CRM_DB_PASSWORD_IDENTITY_AUTHZ
  CRM_DB_PASSWORD_LEAD
  CRM_DB_PASSWORD_ACCOUNT
  CRM_DB_PASSWORD_OPPORTUNITY
  CRM_DB_PASSWORD_COMMERCIAL
  CRM_DB_PASSWORD_WORK
  CRM_DB_PASSWORD_AUDIT_HISTORY
  CRM_DB_PASSWORD_REPORTING
  CRM_DB_PASSWORD_IMPORT_EXPORT
)

runtime_secret_vars=(
  POSTGRES_PASSWORD
  SERVICE_TOKEN_SECRET
  IDENTITY_AUTHZ_DATABASE_URL
  LEAD_DATABASE_URL
  ACCOUNT_DATABASE_URL
  OPPORTUNITY_DATABASE_URL
  COMMERCIAL_DATABASE_URL
  WORK_DATABASE_URL
  AUDIT_HISTORY_DATABASE_URL
  REPORTING_DATABASE_URL
  IMPORT_EXPORT_DATABASE_URL
)

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

require_var() {
  local var_name="$1"
  local value="${!var_name:-}"

  [[ -n "$value" ]] || die "$var_name is not set in $secret_env"
  case "$value" in
    *$'\n'*|*$'\r'*)
      die "$var_name contains a newline or carriage return"
      ;;
  esac
}

reject_development_secret() {
  local var_name="$1"
  local value="${!var_name:-}"

  case "$value" in
    crm_admin_dev_password|crm_dev_service_token_secret|*_dev_password*|*crm_*_dev_password*)
      die "$var_name uses a forbidden development secret value"
      ;;
  esac
}

psql_meta_quote() {
  local value="$1"

  value="${value//\\/\\\\}"
  value="${value//\'/\\\'}"
  printf '%s' "$value"
}

write_psql_secret_vars() {
  local var_name value

  for var_name in "${db_role_secret_vars[@]}"; do
    value="${!var_name}"
    printf "\\set %s '%s'\n" "$var_name" "$(psql_meta_quote "$value")"
  done
}

require_var POSTGRES_DB
require_var POSTGRES_USER
for var_name in "${runtime_secret_vars[@]}" "${db_role_secret_vars[@]}"; do
  require_var "$var_name"
  reject_development_secret "$var_name"
done

files=()
if [[ "$direction" == "up" ]]; then
  while IFS= read -r file; do
    files+=("$file")
  done < <(find "$migrations_dir" -path '*.up.sql' -type f | sort)
else
  while IFS= read -r file; do
    files+=("$file")
  done < <(find "$migrations_dir" -path '*.down.sql' -type f | sort -r)
fi

[[ "${#files[@]}" -gt 0 ]] || die "no migration files found for direction: $direction"

tmp_dir="$(mktemp -d)"
chmod 700 "$tmp_dir"
trap 'rm -rf "$tmp_dir"' EXIT

for file in "${files[@]}"; do
  rel="${file#$migrations_dir/}"
  safe_rel="$(printf '%s' "$rel" | tr '/.' '__')"
  stdin_file="$tmp_dir/$safe_rel.sql"

  {
    write_psql_secret_vars
    printf '\\echo running release migration %s\n' "$rel"
    cat "$file"
  } > "$stdin_file"
  chmod 600 "$stdin_file"

  CRM_RELEASE_STDIN_FILE="$stdin_file" "$script_dir/run-release-step.sh" \
    docker compose \
    --env-file "$env_file" \
    -f "$compose_file" \
    exec -T postgres \
    psql -v ON_ERROR_STOP=1 \
      -U "${POSTGRES_USER:?POSTGRES_USER required}" \
      -d "${POSTGRES_DB:?POSTGRES_DB required}"
done
