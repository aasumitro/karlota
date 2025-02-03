# Exporting bin folder to the path for makefile
export PATH   := $(PWD)/bin:$(PATH)
# Default Shell
export SHELL  := bash
# Type of OS: Linux or Darwin.
export OSTYPE := $(shell uname -s)

# --- Tooling & Variables ----------------------------------------------------------------
include ./misc/make/tools.Makefile

install-deps: gocritic gotestsum golangci-lint
deps: $(GOCRITIC) $(GOTESTSUM) $(GOLANGCI)
deps:
	@ echo "Required Tools Are Available"

.Phony: migration-table
migration-table:
	@if [ "$(name)" ]; then \
    	migrate create -ext sql -dir db/migrations $(name); \
   	fi

.Phony: migrate-up
migrate-up:
	 go run cmd/cli/main.go migration up

.Phony: migrate-up-seed
migrate-up-seed:
	 go run cmd/cli/main.go migration seed

.Phony: migrate-down
migrate-down:
	  go run cmd/cli/main.go migration down

.Phony: mock
mock: $(MOCKERY)
	mockery --all --output=mocks/genmocks --outpkg=mocks


.Phony: api-spec
api-spec: # run-tests
	@ echo "Re-generate API-Spec docs"
	@ swag init --parseDependency --parseInternal \
		--parseDepth 4 -g ./cmd/app/main.go

.Phony: run-tests
run-tests: $(GOTESTSUM) 
	@ echo "Run tests in specific directories . . ."
	@ gotestsum --format pkgname-and-test-fails \
		--hide-summary=skipped \
		-- -coverprofile=cover.out \
		$(shell find ./internal/repository/sql \
		./internal/account ./internal/chat \
		./internal/notify ./internal/utils/mailer \
		 ./internal/utils/cache ./internal/utils/security \
		 ./internal/utils/http/middleware ./internal/utils/http/wrapper \
		 ./internal/utils/ulid \
		 -type d -exec go list {} \;)
	@ rm cover.out

.Phony: run-lint
run-lint: $(GOLANGCI)
	@ echo "Code Convention Check . . ."
	@ golangci-lint cache clean
	@ golangci-lint run -c .golangci.yaml ./...

.Phony: critic-code
critic-code: $(GOCRITIC)
	@ gocritic check -enableAll ./...

.Phony: run-api
run-api:
	@echo "Run Api"
	go mod tidy -compat=1.23
	go run -race ./cmd/app/main.go

.Phony: run-app
run-app: build-fe run-api
	@echo "Run App"

.Phony: run-watch-app
run-watch-app:
	air --build.cmd "go build -o ./temps/watch/main ./cmd/app/main.go" --build.bin "./temps/watch/main"

.Phony: build-fe
build-fe:
	@ echo "Build Frontend"
	@ cd web && npm install && npm run build
