# Project specific variables
COVER_FILE ?= coverage.out

# Database related variables.
PG_USER ?= test
PG_DBS ?= orders 
PG_PASS ?= pass

# Tools
.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@echo Installing tools from tools.txt
	@grep '@' tools.txt | xargs -tI % go install %

# Main targets
.PHONY: test
test: ## Run unit tests
	@go test -race ./... -coverprofile=$(COVER_FILE)
	@go tool cover -func=$(COVER_FILE) | grep ^total

$(COVER_FILE):
	$(MAKE) test

.PHONY: cover
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	@go tool cover -html=$(COVER_FILE)
	@rm -f $(COVER_FILE)

.PHONY: lint
lint: tools ## Check the project with lint
	@golangci-lint run -v --fix

.PHONY: check
check: lint test ## Check project with static checks and unit tests

.PHONE: gen
gen: tools ## Generate projects files and components.
	@oapi-codegen --config api/order/config.yaml api/order/api.yaml > internal/order/infrastructure/api/http.gen.go

.PHONY: deps
deps: ## Manage go mod dependencies, beautify go.mod and go.sum files
	@go mod tidy

.PHONY: clean
clean: ## Clean the project from built files
	@rm -f $(COVER_FILE)

.PHONY: up
up: ## Up local development environments, see hack/docker-compose.yml
	@PG_USER=$(PG_USER) PG_DBS=$(PG_DBS) PG_PASS=$(PG_PASS) docker compose -f hack/docker-compose.yml --project-directory hack up \
		--build --force-recreate --renew-anon-volumes -d

.PHONY: down
down: ## Stop and down local development environments.
	@PG_USER=$(PG_USER) PG_DBS=$(PG_DBS) PG_PASS=$(PG_PASS) docker compose -f hack/docker-compose.yml --project-directory hack down

.PHONY: run
run: ## Run given `service` on local environment.
	@./hack/run.sh $(service)

.PHONY: integrations
integrations: ## Run integrations test
	@./hack/run-in-docker-compose.sh

.PHONY: help
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'
