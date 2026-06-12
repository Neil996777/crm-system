#!/usr/bin/env bash
set -euo pipefail

release_dir="${1:-/opt/crm-system/releases/66d2531}"
env_file="$release_dir/.env.release"
nginx_target="${CRM_NGINX_TARGET:-/etc/nginx/conf.d/crm.conf}"
script_dir="$(cd "$(dirname "${BASH_SOURCE[0]}")" && pwd)"

die() {
  printf 'ERROR: %s\n' "$*" >&2
  exit 1
}

[[ -f "$env_file" ]] || die ".env.release missing: $env_file"

server_name_override="${CRM_SERVER_NAME:-}"
set -a
. "$env_file"
set +a
if [[ -n "$server_name_override" ]]; then
  CRM_SERVER_NAME="$server_name_override"
fi

server_name="${CRM_SERVER_NAME:?CRM_SERVER_NAME required}"
transcript="${CRM_DEPLOY_TRANSCRIPT:-$release_dir/deploy-transcript.log}"
export CRM_DEPLOY_TRANSCRIPT="$transcript"

tmp_conf="$(mktemp)"
trap 'rm -f "$tmp_conf"' EXIT

cat > "$tmp_conf" <<EOF
server {
    listen 80;
    server_name $server_name;

    location /.well-known/acme-challenge/ {
        root /var/www/html;
    }

    location / {
        return 301 https://\$host\$request_uri;
    }
}

server {
    listen 443 ssl http2;
    server_name $server_name;

    ssl_certificate /etc/letsencrypt/live/$server_name/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/$server_name/privkey.pem;
    ssl_session_cache shared:CRMSSL:10m;
    ssl_session_timeout 10m;
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_prefer_server_ciphers off;

    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;
    add_header X-Content-Type-Options "nosniff" always;
    add_header Referrer-Policy "strict-origin-when-cross-origin" always;
    add_header Content-Security-Policy "default-src 'self'; script-src 'self'; style-src 'self' 'unsafe-inline'; img-src 'self' data:; connect-src 'self'; frame-ancestors 'none'; base-uri 'self'; form-action 'self'" always;

    access_log /opt/crm-system/logs/nginx/access.log;
    error_log /opt/crm-system/logs/nginx/error.log warn;

    location = /health {
        proxy_pass http://127.0.0.1:8080/health;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Request-Id \$request_id;
    }

    location /auth/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Request-Id \$request_id;
    }

    location /admin/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Request-Id \$request_id;
    }

    location /api/ {
        proxy_pass http://127.0.0.1:8080;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Request-Id \$request_id;
    }

    location / {
        proxy_pass http://127.0.0.1:8081;
        proxy_http_version 1.1;
        proxy_set_header Host \$host;
        proxy_set_header X-Forwarded-Proto https;
        proxy_set_header X-Forwarded-For \$proxy_add_x_forwarded_for;
        proxy_set_header X-Request-Id \$request_id;
    }
}
EOF

"$script_dir/run-release-step.sh" sudo install -m 0644 "$tmp_conf" "$nginx_target"
"$script_dir/run-release-step.sh" sudo nginx -t
"$script_dir/run-release-step.sh" sudo systemctl reload nginx
