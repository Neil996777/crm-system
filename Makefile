SERVICES := gateway-bff identity-authz lead account opportunity commercial work audit-history reporting import-export
export GOCACHE := $(CURDIR)/.cache/go-build
export GOMODCACHE := $(CURDIR)/.cache/go-mod

.PHONY: test-task001 test-services test-contracts tidy-services static-check migrate-up migrate-down

test-task001: static-check test-contracts test-services

test-services:
	@set -e; for svc in $(SERVICES); do \
		echo "==> go test services/$$svc"; \
		(cd services/$$svc && go test ./...); \
	done

test-contracts:
	@cd shared/contracts && go test ./...

tidy-services:
	@set -e; for svc in $(SERVICES); do \
		echo "==> go mod tidy services/$$svc"; \
		(cd services/$$svc && go mod tidy); \
	done

static-check:
	@bash scripts/task001_static_check.sh

migrate-up:
	@bash scripts/migrate.sh up

migrate-down:
	@bash scripts/migrate.sh down
