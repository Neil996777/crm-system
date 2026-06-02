package handler

import (
	"net/http"

	"crm-system/services/gateway-bff/internal/authz"
)

type Config struct {
	IdentityBaseURL string
	Routes          map[string]string
	HTTPClient      *http.Client
}

type Gateway struct {
	config Config
	authz  authz.Client
}

func NewGatewayServer(config Config) http.Handler {
	gateway := &Gateway{
		config: config,
		authz: authz.Client{
			BaseURL:    config.IdentityBaseURL,
			HTTPClient: config.HTTPClient,
		},
	}
	mux := http.NewServeMux()
	mux.HandleFunc("GET /auth/current", gateway.authProxy)
	mux.HandleFunc("POST /auth/sign-in", gateway.authProxy)
	mux.HandleFunc("POST /auth/sign-out", gateway.authProxy)
	mux.HandleFunc("GET /admin/users", gateway.authProxy)
	mux.HandleFunc("POST /admin/users", gateway.authProxy)
	mux.HandleFunc("PATCH /admin/users/{id}/role", gateway.authProxy)
	mux.HandleFunc("PATCH /admin/users/{id}/status", gateway.authProxy)
	mux.HandleFunc("GET /api/operation-log", gateway.operationLog)
	mux.HandleFunc("GET /api/{resource}/{id}/history", gateway.recordHistory)
	mux.HandleFunc("GET /api/{resource}", gateway.proxy)
	mux.HandleFunc("GET /api/{resource}/{id}", gateway.proxy)
	mux.HandleFunc("GET /api/{resource}/{id}/{child}", gateway.proxy)
	mux.HandleFunc("POST /api/{resource}", gateway.proxy)
	mux.HandleFunc("POST /api/{resource}/duplicate-checks", gateway.proxy)
	mux.HandleFunc("POST /api/{resource}/{id}/{action}", gateway.proxy)
	mux.HandleFunc("PATCH /api/{resource}/{id}", gateway.proxy)
	return mux
}
