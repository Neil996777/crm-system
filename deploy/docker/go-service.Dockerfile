ARG GO_BASE_IMAGE=golang:1.26-alpine
ARG RUNTIME_BASE_IMAGE=alpine:3.21

FROM ${GO_BASE_IMAGE} AS build

WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM ${RUNTIME_BASE_IMAGE}

ARG RELEASE_COMMIT
ARG SERVICE_NAME
ARG IMAGE_SOURCE=crm-system
ARG IMAGE_CREATED

LABEL org.opencontainers.image.revision="${RELEASE_COMMIT}" \
      org.opencontainers.image.source="${IMAGE_SOURCE}" \
      org.opencontainers.image.created="${IMAGE_CREATED}" \
      com.crm.release.content="${RELEASE_COMMIT}" \
      com.crm.service="${SERVICE_NAME}"

RUN addgroup -S crm && adduser -S -G crm crm
WORKDIR /app
COPY --from=build /out/server /app/server
USER crm

EXPOSE 8080
ENTRYPOINT ["/app/server"]
