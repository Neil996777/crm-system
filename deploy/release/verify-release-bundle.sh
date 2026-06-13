#!/usr/bin/env bash
set -euo pipefail

bundle_dir="${1:-/opt/crm-system/releases/66d2531}"
release_commit="${RELEASE_COMMIT:-66d2531}"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -d "$bundle_dir" ]] || die "bundle directory not found: $bundle_dir"
cd "$bundle_dir"

[[ -f image-manifest.json ]] || die "image-manifest.json missing"
[[ -f image-manifest.sha256 ]] || die "image-manifest.sha256 missing"
[[ -f "release-crm-system-$release_commit.sha256" ]] || die "release checksum missing"
[[ -f .env.release ]] || die ".env.release missing"
[[ -f docker-compose.prod.yml ]] || die "docker-compose.prod.yml missing"

sha256sum -c image-manifest.sha256
sha256sum -c "release-crm-system-$release_commit.sha256"

grep -q "\"releaseCommit\": \"$release_commit\"" image-manifest.json || die "manifest release commit mismatch"
grep -q "org.opencontainers.image.revision" image-manifest.json || die "manifest lacks OCI revision labels"
grep -q "CRM_RELEASE_COMMIT=$release_commit" .env.release || die ".env.release commit mismatch"
grep -q '^CRM_IMAGE_POSTGRES=.*@sha256:' .env.release || die "postgres image is not digest-pinned"

if grep -nE '^[[:space:]]+build:' docker-compose.prod.yml; then
  die "compose contains build keys"
fi
if grep -n 'CRM_IMAGE_TAG:-latest' docker-compose.prod.yml; then
  die "compose contains latest fallback"
fi
if grep -n '/opt/crm-system/current/services' docker-compose.prod.yml; then
  die "compose still depends on source checkout migrations"
fi
grep -q 'frontend-web:' docker-compose.prod.yml || die "frontend-web service missing"
grep -q "CRM_RELEASE_COMMIT" docker-compose.prod.yml || die "compose does not use release commit path"

services=(
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

db_role_secret_tokens=(
  __CRM_DB_PASSWORD_IDENTITY_AUTHZ__
  __CRM_DB_PASSWORD_LEAD__
  __CRM_DB_PASSWORD_ACCOUNT__
  __CRM_DB_PASSWORD_OPPORTUNITY__
  __CRM_DB_PASSWORD_COMMERCIAL__
  __CRM_DB_PASSWORD_WORK__
  __CRM_DB_PASSWORD_AUDIT_HISTORY__
  __CRM_DB_PASSWORD_REPORTING__
  __CRM_DB_PASSWORD_IMPORT_EXPORT__
)

is_expected_secret_token() {
  local candidate="$1"
  local token

  for token in "${db_role_secret_tokens[@]}"; do
    [[ "$candidate" == "$token" ]] && return 0
  done

  return 1
}

for service in "${services[@]}"; do
  tar_file="images/$service-$release_commit.tar"
  [[ -f "$tar_file" ]] || die "image archive missing: $tar_file"
  grep -q "crm-system/$service:$release_commit" image-manifest.json || die "manifest missing $service"
done

[[ -f migrations/migration-manifest.sha256 ]] || die "migration manifest missing"
(cd migrations && sha256sum -c migration-manifest.sha256)

if grep -RIn '_dev_password' migrations >/dev/null; then
  grep -RIl '_dev_password' migrations >&2
  die "release migrations contain development database role passwords"
fi

if grep -RIl ":'CRM_DB_PASSWORD_" migrations >/dev/null; then
  grep -RIl ":'CRM_DB_PASSWORD_" migrations >&2
  die "release migrations contain psql secret variables; DO-block migrations require render tokens"
fi

for secret_token in "${db_role_secret_tokens[@]}"; do
  count="$( (grep -RohF "$secret_token" migrations || true) | wc -l | tr -d '[:space:]')"
  [[ "$count" == "1" ]] \
    || die "release migrations contain $count occurrences of $secret_token; expected 1"
done

bad_password_files=()
while IFS= read -r sql_file; do
  while IFS= read -r password_literal; do
    if ! is_expected_secret_token "$password_literal"; then
      bad_password_files+=("$sql_file")
      break
    fi
  done < <(perl -nle 'while (/PASSWORD\s+'\''([^'\'']*)'\''/g) { print $1 }' "$sql_file")
done < <(find migrations -type f -name '*.sql' -print)

if [[ "${#bad_password_files[@]}" -gt 0 ]]; then
  printf '%s\n' "${bad_password_files[@]}" | sort -u >&2
  die "release migrations contain unexpected password literals"
fi

printf 'release bundle verified: %s\n' "$bundle_dir"
