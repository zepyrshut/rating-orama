GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")
CORE_DIR := ./core
UI_DIR := ./ui
DOCS_DIR := ./docs
LIBRARIES_DIR := ./../libraries
PG_VERSION := 16.4-alpine3.20
DB_NAME := rating
MOD_NAME := rating-orama

.PHONY: sayhello
# Print Hello World
sayhello:
	@echo "Hello World"

.PHONY: dockerize
# Creates a development database.
dockerize:
	docker run --name $(DB_NAME)-db-dev -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=developer -e POSTGRES_DB=$(DB_NAME) -p 5432:5432 -d postgres:$(PG_VERSION)

.PHONY: dockerize-test
# Creates a test database.
dockerize-test:
	docker run --name $(DB_NAME)-db-test -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=developer -e POSTGRES_DB=$(DB_NAME) -p 5433:5432 -d postgres:$(PG_VERSION)

.PHONY: undockerize
# Destroy a development database.
undockerize:
	docker rm -f $(DB_NAME)-db-dev

.PHONY: undockerize-test
# Destroy a test database.
undockerize-test:
	docker rm -f $(DB_NAME)-db-test

.PHONY: restart-db
# Restart a development database.
restart-db:
	make undockerize
	make dockerize

.PHONY: restart-db-test
# Restart a test database.
restart-db-test:
	make undockerize-test
	make dockerize-test

.PHONY: migrateup
# Migrate all schemas, triggers and data located in database/migrations.
migrateup:
	migrate -path $(CORE_DIR)/cmd/database/migrations -database "postgresql://developer:secret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose up

.PHONY: migratedown
# Migrate all schemas, triggers and data located in database/migrations.
migratedown:
	migrate -path $(CORE_DIR)/cmd/database/migrations -database "postgresql://developer:secret@localhost:5432/$(DB_NAME)?sslmode=disable" -verbose down

.PHONY: pg-dump
# Dump database to file.
pg-dump:
	docker exec -e PGPASSWORD=secret $(DB_NAME)-db-dev pg_dump -U developer --column-inserts --data-only $(DB_NAME) > $(CORE_DIR)/cmd/database/data/data.sql
	sed -i '1iSET session_replication_role = '\''replica'\'';' $(CORE_DIR)/cmd/database/data/data.sql
	sed -i '$$aSET session_replication_role = '\''origin'\'';' $(CORE_DIR)/cmd/database/data/data.sql

.PHONY: pg-restore
# Restore database from file.
pg-restore:
	docker cp $(CORE_DIR)/cmd/database/data/data.sql $(DB_NAME)-db-dev:/data.sql
	docker exec -e PGPASSWORD=secret $(DB_NAME)-db-dev psql -U developer -d $(DB_NAME) -f data.sql

.PHONY: pg-docs
# Generate docs from database.
pg-docs:
	java -jar $(LIBRARIES_DIR)/schemaspy-6.2.4.jar -t pgsql -dp $(LIBRARIES_DIR)/postgresql-42.7.4.jar -db $(DB_NAME) -host localhost -port 5432 -u developer -p secret -o $(DOCS_DIR)/database -vizjs

.PHONY: sqlc
# Generate or recreate SQLC queries.
sqlc:
	cd $(CORE_DIR) && sqlc generate

.PHONY: test
# Test all files and generate coverage file.
test:
	cd $(CORE_DIR) && $(GO) test ./... -v -covermode=count -coverprofile=./benchmark/coverage.out $(PACKAGES)

.PHONY: gomock
# Generate mock files.
gomock:
	cd $(CORE_DIR) && mockgen -package mock -destination internal/repository/mock/querier.go $(MOD_NAME)/internal/repository ExtendedQuerier

.PHONY: run
# Run project.
run:
	cd $(CORE_DIR) && $(GO) run ./cmd/.

.PHONY: bench
# Run benchmarks.
bench:
	cd $(CORE_DIR) && test -f benchmark/new_benchmark.txt && mv benchmark/new_benchmark.txt benchmark/old_benchmark.txt || true
	cd $(CORE_DIR) && $(GO) test ./... -bench=. -count=10 -benchmem > benchmark/new_benchmark.txt
	cd $(CORE_DIR) && benchstat benchmark/old_benchmark.txt benchmark/new_benchmark.txt > benchmark/benchstat.txt

.PHONY: recreate
# Destroy development DB and generate ones.
recreate:
	echo "y" | make migratedown
	make migrateup

.PHONY: tidy
# Runs a go mod tidy
tidy:
	cd $(CORE_DIR) && $(GO) mod tidy

.PHONY: build-linux
# Build and generate linux executable.
build-linux:
	cd $(CORE_DIR) && go mod tidy
	cd $(CORE_DIR) && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./tmp/$(MOD_NAME) ./cmd/.



.PHONY: pack-docker
# Run docker build for pack binary and assets to Docker container.
pack-docker:
	make test
	make build-linux
	docker build -t $(MOD_NAME):${version} -t $(MOD_NAME):latest .

.PHONY: remove-debug
# Remove all debug entries for reduce size binary.
remove-debug:
	cd $(CORE_DIR) && find . -name "*.go" -type f -exec sed -i '/slog\.Debug/d' {} +

.PHONY: fmt
# Ensure consistent code formatting.
fmt:
	cd $(CORE_DIR) && $(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
# format (check only).
fmt-check:
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

.PHONY: vet
# Examine packages and report suspicious constructs if any.
vet:
	cd $(CORE_DIR) && $(GO) vet $(VETPACKAGES)

.PHONY: tools
# Install tools (migrate and sqlc).
tools:
	@if [ $(GO_VERSION) -gt 16 ]; then \
		cd $(CORE_DIR) && $(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		cd $(CORE_DIR) && $(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	fi

.PHONY: env
# Copy .env.example to .env if .env does not already exist
env:
	cd $(CORE_DIR) && @if [ ! -f .env ]; then \
		cp .env.example .env; \
		echo ".env file created from .env.example"; \
	else \
		echo ".env file already exists"; \
	fi


.PHONY: first-run
# Runs for the first time
first-run:
	make tools
	make env
	make recreate
	make run

.PHONY: help
# Help.
help:
	@echo ''
	@echo 'Usage:'
	@echo ' make [target]'
	@echo ''
	@echo 'Targets:'
	@awk '/^[a-zA-Z\-\0-9]+:/ { \
	helpMessage = match(lastLine, /^# (.*)/); \
		if (helpMessage) { \
			helpCommand = substr($$1, 0, index($$1, ":")-1); \
			helpMessage = substr(lastLine, RSTART + 2, RLENGTH); \
			printf " - \033[36m%-20s\033[0m %s\n", helpCommand, helpMessage; \
		} \
	} \
	{ lastLine = $$0 }' $(MAKEFILE_LIST)

.DEFAULT_GOAL := help