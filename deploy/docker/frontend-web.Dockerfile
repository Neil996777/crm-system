ARG NODE_BASE_IMAGE=node:24-alpine
ARG NGINX_BASE_IMAGE=nginx:1.27-alpine

FROM ${NODE_BASE_IMAGE} AS build

WORKDIR /app
COPY package*.json ./
RUN npm ci
COPY . .
RUN npm run build

FROM ${NGINX_BASE_IMAGE}

ARG RELEASE_COMMIT
ARG SERVICE_NAME=frontend-web
ARG IMAGE_SOURCE=crm-system
ARG IMAGE_CREATED

LABEL org.opencontainers.image.revision="${RELEASE_COMMIT}" \
      org.opencontainers.image.source="${IMAGE_SOURCE}" \
      org.opencontainers.image.created="${IMAGE_CREATED}" \
      com.crm.release.content="${RELEASE_COMMIT}" \
      com.crm.service="${SERVICE_NAME}"

RUN mkdir -p /usr/share/nginx/html /var/cache/nginx /var/run /tmp \
    && chown -R nginx:nginx /usr/share/nginx/html /var/cache/nginx /var/run /tmp \
    && printf '%s\n' \
      'worker_processes auto;' \
      'pid /tmp/nginx.pid;' \
      'events { worker_connections 1024; }' \
      'http {' \
      '  include /etc/nginx/mime.types;' \
      '  default_type application/octet-stream;' \
      '  access_log /dev/stdout;' \
      '  error_log /dev/stderr warn;' \
      '  sendfile on;' \
      '  server {' \
      '    listen 8080;' \
      '    server_name _;' \
      '    root /usr/share/nginx/html;' \
      '    index index.html;' \
      '    location = /health { return 200 "ok\n"; add_header Content-Type text/plain; }' \
      '    location / { try_files $uri $uri/ /index.html; }' \
      '  }' \
      '}' > /etc/nginx/nginx.conf

COPY --from=build /app/dist /usr/share/nginx/html

USER nginx
EXPOSE 8080
HEALTHCHECK --interval=10s --timeout=5s --retries=12 CMD wget -qO- http://127.0.0.1:8080/health >/dev/null 2>&1 || exit 1
ENTRYPOINT ["nginx", "-g", "daemon off;"]
