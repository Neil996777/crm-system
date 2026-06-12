#!/usr/bin/env bash
set -euo pipefail

bundle_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$bundle_dir/.env.release"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

service_var_name() {
  printf '%s' "$1" | tr '[:lower:]' '[:upper:]' | tr '-' '_'
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"

set -a
. "$env_file"
set +a

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

for service in "${services[@]}"; do
  upper="$(service_var_name "$service")"
  image_var="CRM_IMAGE_$upper"
  id_var="CRM_IMAGE_ID_$upper"
  image="${!image_var:-}"
  expected_id="${!id_var:-}"
  [[ -n "$image" ]] || die "missing $image_var"
  [[ -n "$expected_id" ]] || die "missing $id_var"

  actual_id="$(docker image inspect "$image" --format '{{.Id}}')"
  [[ "$actual_id" == "$expected_id" ]] || die "$service image id mismatch: $actual_id != $expected_id"

  revision="$(docker image inspect "$image" --format '{{ index .Config.Labels "org.opencontainers.image.revision" }}')"
  [[ "$revision" == "$CRM_RELEASE_COMMIT" ]] || die "$service revision label mismatch: $revision"

  release_content="$(docker image inspect "$image" --format '{{ index .Config.Labels "com.crm.release.content" }}')"
  [[ "$release_content" == "$CRM_RELEASE_COMMIT" ]] || die "$service release content label mismatch: $release_content"
done

[[ "${CRM_IMAGE_POSTGRES:-}" == *@sha256:* ]] || die "postgres image is not digest-pinned"
docker image inspect "$CRM_IMAGE_POSTGRES" >/dev/null

printf 'loaded image verification complete for commit %s\n' "$CRM_RELEASE_COMMIT"
