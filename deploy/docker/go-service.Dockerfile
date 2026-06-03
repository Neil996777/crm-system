FROM golang:1.26-alpine AS build

WORKDIR /src
ENV GOPROXY=https://goproxy.cn,direct
ENV GOSUMDB=sum.golang.google.cn
COPY . .
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux go build -o /out/server ./cmd/server

FROM alpine:3.21

RUN addgroup -S crm && adduser -S -G crm crm
WORKDIR /app
COPY --from=build /out/server /app/server
USER crm

EXPOSE 8080
ENTRYPOINT ["/app/server"]
