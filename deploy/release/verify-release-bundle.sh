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

for service in "${services[@]}"; do
  tar_file="images/$service-$release_commit.tar"
  [[ -f "$tar_file" ]] || die "image archive missing: $tar_file"
  grep -q "crm-system/$service:$release_commit" image-manifest.json || die "manifest missing $service"
done

[[ -f migrations/migration-manifest.sha256 ]] || die "migration manifest missing"
(cd migrations && sha256sum -c migration-manifest.sha256)

printf 'release bundle verified: %s\n' "$bundle_dir"
