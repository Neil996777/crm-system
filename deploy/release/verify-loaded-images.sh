#!/usr/bin/env bash
set -euo pipefail

bundle_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$bundle_dir/.env.release"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

sha256_file() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum "$1" | awk '{print $1}'
  else
    shasum -a 256 "$1" | awk '{print $1}'
  fi
}

sha256_stream() {
  if command -v sha256sum >/dev/null 2>&1; then
    sha256sum | awk '{print $1}'
  else
    shasum -a 256 | awk '{print $1}'
  fi
}

service_var_name() {
  printf '%s' "$1" | tr '[:lower:]' '[:upper:]' | tr '-' '_'
}

archive_config_digest() {
  local archive="$1"
  local manifest_json config_path digest

  manifest_json="$(tar -xOf "$archive" manifest.json)"
  config_path="$(printf '%s' "$manifest_json" | sed -nE 's/.*"Config":"([^"]+)".*/\1/p' | head -n 1)"
  [[ -n "$config_path" ]] || die "image archive manifest lacks Config entry: $archive"

  digest="$(tar -xOf "$archive" "$config_path" | sha256_stream)"
  printf 'sha256:%s' "$digest"
}

cleanup_tmp_dir() {
  [[ -n "${tmp_dir:-}" && -d "$tmp_dir" ]] || return 0
  rm -rf "$tmp_dir"
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"
[[ -d "$bundle_dir/images" ]] || die "images directory missing: $bundle_dir/images"

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

tmp_dir="$(mktemp -d)"
trap cleanup_tmp_dir EXIT

printf 'service\timage\tarchive_sha256\tarchive_config_digest\tloaded_config_digest\trevision_label\trelease_content_label\n'

for service in "${services[@]}"; do
  upper="$(service_var_name "$service")"
  image_var="CRM_IMAGE_$upper"
  id_var="CRM_IMAGE_ID_$upper"
  tar_var="CRM_IMAGE_TAR_SHA256_$upper"
  image="${!image_var:-}"
  expected_config_digest="${!id_var:-}"
  expected_archive_sha="${!tar_var:-}"
  archive="$bundle_dir/images/$service-$CRM_RELEASE_COMMIT.tar"
  loaded_archive="$tmp_dir/$service-loaded.tar"

  [[ -n "$image" ]] || die "missing $image_var"
  [[ -n "$expected_config_digest" ]] || die "missing $id_var"
  [[ -n "$expected_archive_sha" ]] || die "missing $tar_var"
  [[ -f "$archive" ]] || die "image archive missing: $archive"

  actual_archive_sha="$(sha256_file "$archive")"
  [[ "$actual_archive_sha" == "$expected_archive_sha" ]] \
    || die "$service archive sha mismatch: $actual_archive_sha != $expected_archive_sha"

  archive_config="$(archive_config_digest "$archive")"
  [[ "$archive_config" == "$expected_config_digest" ]] \
    || die "$service archive config digest mismatch: $archive_config != $expected_config_digest"

  docker save "$image" -o "$loaded_archive"
  loaded_config="$(archive_config_digest "$loaded_archive")"
  rm -f "$loaded_archive"
  [[ "$loaded_config" == "$expected_config_digest" ]] \
    || die "$service loaded config digest mismatch: $loaded_config != $expected_config_digest"

  revision="$(docker image inspect "$image" --format '{{ index .Config.Labels "org.opencontainers.image.revision" }}')"
  [[ "$revision" == "$CRM_RELEASE_COMMIT" ]] || die "$service revision label mismatch: $revision"

  release_content="$(docker image inspect "$image" --format '{{ index .Config.Labels "com.crm.release.content" }}')"
  [[ "$release_content" == "$CRM_RELEASE_COMMIT" ]] || die "$service release content label mismatch: $release_content"

  printf '%s\t%s\tsha256:%s\t%s\t%s\t%s\t%s\n' \
    "$service" \
    "$image" \
    "$actual_archive_sha" \
    "$archive_config" \
    "$loaded_config" \
    "$revision" \
    "$release_content"
done

[[ "${CRM_IMAGE_POSTGRES:-}" == *@sha256:* ]] || die "postgres image is not digest-pinned"
docker image inspect "$CRM_IMAGE_POSTGRES" >/dev/null

printf 'loaded image verification complete for commit %s\n' "$CRM_RELEASE_COMMIT"
