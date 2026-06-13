#!/usr/bin/env bash
set -euo pipefail

RELEASE_COMMIT="${RELEASE_COMMIT:-66d2531}"
PIPELINE_ROOT="${PIPELINE_ROOT:-$(cd "$(dirname "${BASH_SOURCE[0]}")/../.." && pwd)}"
SOURCE_DIR="${RELEASE_SOURCE_DIR:-$PIPELINE_ROOT}"
BUNDLE_PARENT="${RELEASE_BUNDLE_PARENT:-$PIPELINE_ROOT/dist}"
BUNDLE_DIR="${RELEASE_BUNDLE_DIR:-$BUNDLE_PARENT/release-crm-system-$RELEASE_COMMIT}"
PLATFORM="${RELEASE_PLATFORM:-linux/amd64}"
POSTGRES_TAG="${CRM_POSTGRES_TAG:-postgres:16-alpine}"
REPO_SOURCE="${REPO_SOURCE:-$(git -C "$PIPELINE_ROOT" config --get remote.origin.url || printf 'crm-system')}"
IMAGE_CREATED="${IMAGE_CREATED:-$(date -u +%Y-%m-%dT%H:%M:%SZ)}"

go_services=(
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
)

app_images=("${go_services[@]}" frontend-web)

db_role_names=(
  crm_identity_authz_user
  crm_lead_user
  crm_account_user
  crm_opportunity_user
  crm_commercial_user
  crm_work_user
  crm_audit_history_user
  crm_reporting_user
  crm_import_export_user
)

db_role_old_passwords=(
  crm_identity_authz_dev_password
  crm_lead_dev_password
  crm_account_dev_password
  crm_opportunity_dev_password
  crm_commercial_dev_password
  crm_work_dev_password
  crm_audit_history_dev_password
  crm_reporting_dev_password
  crm_import_export_dev_password
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

service_var_name() {
  printf '%s' "$1" | tr '[:lower:]' '[:upper:]' | tr '-' '_'
}

parameterize_db_role_passwords() {
  local i role old_password secret_var file

  for i in "${!db_role_names[@]}"; do
    role="${db_role_names[$i]}"
    old_password="${db_role_old_passwords[$i]}"
    secret_var="${db_role_secret_vars[$i]}"

    while IFS= read -r file; do
      ROLE="$role" OLD_PASSWORD="$old_password" SECRET_VAR="$secret_var" perl -0pi -e '
        my $role = $ENV{"ROLE"};
        my $old_password = $ENV{"OLD_PASSWORD"};
        my $secret_var = $ENV{"SECRET_VAR"};
        my $needle = "CREATE ROLE $role LOGIN PASSWORD " . chr(39) . "$old_password" . chr(39) . ";";
        my $replacement = "CREATE ROLE $role LOGIN PASSWORD :" . chr(39) . "$secret_var" . chr(39) . ";";
        s/\Q$needle\E/$replacement/g;
      ' "$file"
    done < <(find "$BUNDLE_DIR/migrations" -type f -name '*.up.sql' -print)
  done

  if grep -RIn '_dev_password' "$BUNDLE_DIR/migrations" >/dev/null; then
    grep -RIn '_dev_password' "$BUNDLE_DIR/migrations" >&2
    die "release migrations still contain development database role passwords"
  fi

  for secret_var in "${db_role_secret_vars[@]}"; do
    grep -RIn ":'$secret_var'" "$BUNDLE_DIR/migrations" >/dev/null \
      || die "release migrations do not reference $secret_var"
  done
}

expected_full="$(git -C "$SOURCE_DIR" rev-parse "$RELEASE_COMMIT^{commit}")"
actual_full="$(git -C "$SOURCE_DIR" rev-parse HEAD)"
[[ "$expected_full" == "$actual_full" ]] || die "release source is not $RELEASE_COMMIT"

if ! git -C "$SOURCE_DIR" diff --quiet || ! git -C "$SOURCE_DIR" diff --cached --quiet; then
  die "release source has tracked changes"
fi

rm -rf "$BUNDLE_DIR"
mkdir -p "$BUNDLE_DIR/images" \
  "$BUNDLE_DIR/deploy/scripts" \
  "$BUNDLE_DIR/deploy/backup" \
  "$BUNDLE_DIR/deploy/healthcheck" \
  "$BUNDLE_DIR/evidence" \
  "$BUNDLE_DIR/migrations" \
  "$BUNDLE_DIR/test-results/backend" \
  "$BUNDLE_DIR/test-results/frontend-build" \
  "$BUNDLE_DIR/test-results/e2e"

if [[ -d "$PIPELINE_ROOT/test-results/backend" ]]; then
  cp -R "$PIPELINE_ROOT/test-results/backend/." "$BUNDLE_DIR/test-results/backend/"
fi
if [[ -f "$PIPELINE_ROOT/test-results/frontend-build.log" ]]; then
  cp "$PIPELINE_ROOT/test-results/frontend-build.log" "$BUNDLE_DIR/test-results/frontend-build/build.log"
fi
if [[ -d "$SOURCE_DIR/frontend/test-results" ]]; then
  cp -R "$SOURCE_DIR/frontend/test-results/." "$BUNDLE_DIR/test-results/e2e/"
fi
if [[ -d "$SOURCE_DIR/frontend/playwright-report" ]]; then
  cp -R "$SOURCE_DIR/frontend/playwright-report" "$BUNDLE_DIR/test-results/e2e/playwright-report"
fi

cp "$PIPELINE_ROOT/docker-compose.prod.yml" "$BUNDLE_DIR/docker-compose.prod.yml"
cp "$PIPELINE_ROOT/deploy/ops/go-live-runbook.md" "$BUNDLE_DIR/deploy/runbook.md"
cp "$PIPELINE_ROOT"/deploy/release/*.sh "$BUNDLE_DIR/deploy/scripts/"
cp "$PIPELINE_ROOT"/deploy/backup/*.sh "$BUNDLE_DIR/deploy/backup/"
cp "$PIPELINE_ROOT"/deploy/backup/*.md "$BUNDLE_DIR/deploy/backup/"
cp "$PIPELINE_ROOT"/deploy/healthcheck/*.sh "$BUNDLE_DIR/deploy/healthcheck/"
cp "$PIPELINE_ROOT/delivery/cicd-release-evidence-template.md" "$BUNDLE_DIR/evidence/cicd-release-evidence-template.md"
chmod +x "$BUNDLE_DIR"/deploy/scripts/*.sh "$BUNDLE_DIR"/deploy/backup/*.sh "$BUNDLE_DIR"/deploy/healthcheck/*.sh

for migrations_dir in "$SOURCE_DIR"/services/*/migrations; do
  [[ -d "$migrations_dir" ]] || continue
  service="$(basename "$(dirname "$migrations_dir")")"
  mkdir -p "$BUNDLE_DIR/migrations/$service"
  cp "$migrations_dir"/*.sql "$BUNDLE_DIR/migrations/$service/"
done

parameterize_db_role_passwords

(
  cd "$BUNDLE_DIR/migrations"
  find . -type f -name '*.sql' | sort | while IFS= read -r file; do
    clean_file="${file#./}"
    printf '%s  %s\n' "$(sha256_file "$clean_file")" "$clean_file"
  done
) > "$BUNDLE_DIR/migrations/migration-manifest.sha256"

docker pull --platform "$PLATFORM" "$POSTGRES_TAG"
postgres_repo_digest="$(docker image inspect "$POSTGRES_TAG" --format '{{index .RepoDigests 0}}')"
[[ "$postgres_repo_digest" == *@sha256:* ]] || die "postgres digest not available"
postgres_digest="sha256:${postgres_repo_digest##*@sha256:}"
postgres_reference="${POSTGRES_TAG}@${postgres_digest}"

manifest_entries=()

record_image() {
  local service="$1"
  local image="$2"
  local tar_name="$3"
  local tar_path="$BUNDLE_DIR/images/$tar_name"
  local image_id tar_sha upper

  docker save "$image" -o "$tar_path"
  image_id="$(docker image inspect "$image" --format '{{.Id}}')"
  tar_sha="$(sha256_file "$tar_path")"
  upper="$(service_var_name "$service")"

  printf 'CRM_IMAGE_%s=%s\n' "$upper" "$image" >> "$BUNDLE_DIR/.env.release"
  printf 'CRM_IMAGE_ID_%s=%s\n' "$upper" "$image_id" >> "$BUNDLE_DIR/.env.release"
  printf 'CRM_IMAGE_TAR_SHA256_%s=%s\n' "$upper" "$tar_sha" >> "$BUNDLE_DIR/.env.release"

  manifest_entries+=("$(cat <<ENTRY
    {
      "service": "$service",
      "image": "$image",
      "imageDigest": "$image_id",
      "imageId": "$image_id",
      "archive": "images/$tar_name",
      "archiveSha256": "$tar_sha",
      "platform": "$PLATFORM",
      "sourceCommit": "$RELEASE_COMMIT",
      "labels": {
        "org.opencontainers.image.revision": "$RELEASE_COMMIT",
        "org.opencontainers.image.source": "$REPO_SOURCE",
        "org.opencontainers.image.created": "$IMAGE_CREATED",
        "com.crm.release.content": "$RELEASE_COMMIT",
        "com.crm.service": "$service"
      }
    }
ENTRY
)")
}

cat > "$BUNDLE_DIR/.env.release" <<EOF
COMPOSE_PROJECT_NAME=crm-system
CRM_RELEASE_COMMIT=$RELEASE_COMMIT
CRM_RELEASE_DIR=/opt/crm-system/releases/$RELEASE_COMMIT
CRM_DEPLOY_SECRET_ENV=/opt/crm-system/secrets/prod.env
CRM_SERVER_NAME=118.196.44.193
CRM_IMAGE_POSTGRES=$postgres_reference
EOF

for service in "${go_services[@]}"; do
  image="crm-system/$service:$RELEASE_COMMIT"
  docker build \
    --platform "$PLATFORM" \
    --build-arg RELEASE_COMMIT="$RELEASE_COMMIT" \
    --build-arg SERVICE_NAME="$service" \
    --build-arg IMAGE_SOURCE="$REPO_SOURCE" \
    --build-arg IMAGE_CREATED="$IMAGE_CREATED" \
    --label "org.opencontainers.image.revision=$RELEASE_COMMIT" \
    --label "org.opencontainers.image.source=$REPO_SOURCE" \
    --label "org.opencontainers.image.created=$IMAGE_CREATED" \
    --label "com.crm.release.content=$RELEASE_COMMIT" \
    --label "com.crm.service=$service" \
    -f "$PIPELINE_ROOT/deploy/docker/go-service.Dockerfile" \
    -t "$image" \
    "$SOURCE_DIR/services/$service"
  record_image "$service" "$image" "$service-$RELEASE_COMMIT.tar"
done

frontend_image="crm-system/frontend-web:$RELEASE_COMMIT"
docker build \
  --platform "$PLATFORM" \
  --build-arg RELEASE_COMMIT="$RELEASE_COMMIT" \
  --build-arg SERVICE_NAME="frontend-web" \
  --build-arg IMAGE_SOURCE="$REPO_SOURCE" \
  --build-arg IMAGE_CREATED="$IMAGE_CREATED" \
  --label "org.opencontainers.image.revision=$RELEASE_COMMIT" \
  --label "org.opencontainers.image.source=$REPO_SOURCE" \
  --label "org.opencontainers.image.created=$IMAGE_CREATED" \
  --label "com.crm.release.content=$RELEASE_COMMIT" \
  --label "com.crm.service=frontend-web" \
  -f "$PIPELINE_ROOT/deploy/docker/frontend-web.Dockerfile" \
  -t "$frontend_image" \
  "$SOURCE_DIR/frontend"
record_image "frontend-web" "$frontend_image" "frontend-web-$RELEASE_COMMIT.tar"

{
  printf '{\n'
  printf '  "schemaVersion": 1,\n'
  printf '  "releaseCommit": "%s",\n' "$RELEASE_COMMIT"
  printf '  "created": "%s",\n' "$IMAGE_CREATED"
  printf '  "platform": "%s",\n' "$PLATFORM"
  printf '  "distribution": "export-load",\n'
  printf '  "images": [\n'
  for i in "${!manifest_entries[@]}"; do
    if [[ "$i" -gt 0 ]]; then
      printf ',\n'
    fi
    printf '%s' "${manifest_entries[$i]}"
  done
  printf '\n  ],\n'
  printf '  "thirdPartyImages": [\n'
  printf '    {\n'
  printf '      "service": "postgres",\n'
  printf '      "image": "%s",\n' "$postgres_reference"
  printf '      "imageDigest": "%s",\n' "$postgres_digest"
  printf '      "platform": "%s",\n' "$PLATFORM"
  printf '      "sourceCommit": null\n'
  printf '    }\n'
  printf '  ]\n'
  printf '}\n'
} > "$BUNDLE_DIR/image-manifest.json"

printf '%s  image-manifest.json\n' "$(sha256_file "$BUNDLE_DIR/image-manifest.json")" > "$BUNDLE_DIR/image-manifest.sha256"

evidence_file="$BUNDLE_DIR/evidence/cicd-release-evidence-$RELEASE_COMMIT.md"
cp "$PIPELINE_ROOT/delivery/cicd-release-evidence-template.md" "$evidence_file"
{
  printf '\n## CI Artifact Pointers\n\n'
  printf '%s\n' "- Release commit: \`$RELEASE_COMMIT\`"
  printf '%s\n' '- Image manifest: `image-manifest.json`'
  printf '%s\n' '- Manifest checksum: `image-manifest.sha256`'
  printf '%s\n' '- Backend test logs: `test-results/backend/*.log`'
  printf '%s\n' '- Frontend build log: `test-results/frontend-build/build.log`'
  printf '%s\n' '- Playwright artifacts: `test-results/e2e/`'
  printf '%s\n' "- Deployment transcript: fill during G11 from \`/opt/crm-system/releases/$RELEASE_COMMIT/deploy-transcript.log\`."
  printf '%s\n' '- Rollback point: fill during G11 before running deployment.'
} >> "$evidence_file"

(
  cd "$BUNDLE_DIR"
  find . -type f ! -name "release-crm-system-$RELEASE_COMMIT.sha256" | sort | while IFS= read -r file; do
    clean_file="${file#./}"
    printf '%s  %s\n' "$(sha256_file "$clean_file")" "$clean_file"
  done
) > "$BUNDLE_DIR/release-crm-system-$RELEASE_COMMIT.sha256"

tar_path="$BUNDLE_PARENT/release-crm-system-$RELEASE_COMMIT.tar.gz"
tar -C "$BUNDLE_PARENT" -czf "$tar_path" "release-crm-system-$RELEASE_COMMIT"
printf '%s  %s\n' "$(sha256_file "$tar_path")" "$(basename "$tar_path")" > "$tar_path.sha256"

printf 'release bundle ready: %s\n' "$BUNDLE_DIR"
