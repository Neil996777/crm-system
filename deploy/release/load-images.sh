#!/usr/bin/env bash
set -euo pipefail

bundle_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$bundle_dir/.env.release"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"

set -a
. "$env_file"
set +a

for image_tar in "$bundle_dir"/images/*.tar; do
  [[ -f "$image_tar" ]] || die "no image archives found"
  docker load -i "$image_tar"
done

[[ "${CRM_IMAGE_POSTGRES:-}" == *@sha256:* ]] || die "CRM_IMAGE_POSTGRES is not digest-pinned"
docker pull --platform linux/amd64 "$CRM_IMAGE_POSTGRES"

printf 'image load complete for %s\n' "$bundle_dir"
