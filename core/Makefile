GO ?= go
GOFMT ?= gofmt "-s"
GO_VERSION=$(shell $(GO) version | cut -c 14- | cut -d' ' -f1 | cut -d'.' -f2)
PACKAGES ?= $(shell $(GO) list ./...)
VETPACKAGES ?= $(shell $(GO) list ./...)
GOFILES := $(shell find . -name "*.go")

.PHONY: sayhello
# Print Hello World
sayhello:
	@echo "Hello World"

.PHONY: dockerize
# Creates a development database.
dockerize:
	docker run --name rating-db-dev -e POSTGRES_PASSWORD=secret -e POSTGRES_USER=developer -e POSTGRES_DB=rating -p 5432:5432 -d postgres:16.3-alpine3.20

.PHONY: undockerize
# Destroy a development database.
undockerize:
	docker rm -f rating-db-dev

.PHONY: migrateup
# Migrate all schemas, triggers and data located in cmd/database/migrations.
migrateup:
	migrate -path cmd/database/migrations -database "postgresql://developer:secret@localhost:5432/rating?sslmode=disable" -verbose up

.PHONY: sqlc
# Generate or recreate SQLC queries.
sqlc:
	sqlc generate

.PHONY: test
# Test all files and generate coverage file.
test:
	$(GO) test -v -covermode=count -coverprofile=coverage.out $(PACKAGES)

.PHONY: gomock
# Generate mock files.
gomock:
	mockgen -package mock -destination internal/repository/mock/querier.go arena-coins/internal/repository ExtendedQuerier

.PHONY: run
# Run project.
run:
	$(GO) run ./cmd/.

.PHONY: recreate
# Destroy development DB and generate ones.
recreate:
	make undockerize
	make dockerize
	sleep 2
	make migrateup

.PHONY: tidy
# Runs a go mod tidy
tidy:
	$(GO) mod tidy

.PHONY: build-linux
# Build and generate linux executable.
build-linux:
	make tidy
	make remove-debug
	CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o ./tmp/rating ./cmd/.

.PHONY: pack-docker
# Run docker build for pack binary and assets to Docker container.
pack-docker:
	make test
	make build-linux
	docker build -t arena-coins:${version} -t arena-coins:latest .

.PHONY: remove-debug
# Remove all debug entries for reduce size binary.
remove-debug:
	find . -name "*.go" -type f -exec sed -i '/slog\.Debug/d' {} +

.PHONY: fmt
# Ensure consistent code formatting.
fmt:
	$(GOFMT) -w $(GOFILES)

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
	$(GO) vet $(VETPACKAGES)

.PHONY: tools
# Install tools (migrate and sqlc).
tools:
	@if [ $(GO_VERSION) -gt 16 ]; then \
		$(GO) install -tags 'postgres' github.com/golang-migrate/migrate/v4/cmd/migrate@latest; \
		$(GO) install github.com/sqlc-dev/sqlc/cmd/sqlc@latest; \
	fi

.PHONY: env
# Copy .env.example to .env if .env does not already exist
env:
	@if [ ! -f .env ]; then \
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