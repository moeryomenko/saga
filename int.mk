# Database related variables.
PG_USER ?= test
PG_DBS ?= "orders,payments"
PG_PASS ?= pass

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
