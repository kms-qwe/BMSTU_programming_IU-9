APP_DIR=./app
LOCALBIN=./.localbin
GOOSE_BIN=$(LOCALBIN)/goose
SWAG_BIN=$(LOCALBIN)/swag
MIGRATIONS_DIR=./migrations
GOOSE_PKG=github.com/pressly/goose/v3/cmd/goose
SWAG_PKG=github.com/swaggo/swag/cmd/swag
GOOSE_VERSION=v3.24.1
SWAG_VERSION=v1.16.4

.PHONY: run
run:
	cd $(APP_DIR) && go run ./cmd/app

.PHONY: fmt
fmt:
	cd $(APP_DIR) && gofmt -w ./cmd ./internal ./pkg

.PHONY: tidy
tidy:
	cd $(APP_DIR) && go mod tidy

.PHONY: install-goose
install-goose:
	mkdir -p $(LOCALBIN)
	GOBIN=$(PWD)/$(LOCALBIN) go install $(GOOSE_PKG)@$(GOOSE_VERSION)

.PHONY: install-swag
install-swag:
	mkdir -p $(LOCALBIN)
	GOBIN=$(PWD)/$(LOCALBIN) go install $(SWAG_PKG)@$(SWAG_VERSION)

.PHONY: install-tools
install-tools: install-goose install-swag

.PHONY: swagger
swagger: install-swag
	cd $(APP_DIR) && ../$(SWAG_BIN) init -g cmd/app/swagger.go -o ../docs/swagger

.PHONY: migration-create
migration-create: install-goose
	@if [ -z "$(name)" ]; then echo "usage: make migration-create name=create_tables"; exit 1; fi
	$(GOOSE_BIN) -dir $(MIGRATIONS_DIR) create $(name) sql

.PHONY: up
up:
	cp build/.env.example build/.env 2>/dev/null || true
	docker compose --env-file build/.env -f build/docker-compose.yml up --build

.PHONY: down
down:
	docker compose --env-file build/.env -f build/docker-compose.yml down -v
