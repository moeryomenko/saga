# Project specific variables
COVER_FILE ?= coverage.out

.PHONY: test
test: tools ## Run unit tests
	@set -o pipefail && go test -json -count=1 -cover -race -coverpkg=./... ./... | tparse

.PHONY: cover
cover: ## Output coverage in human readable form in html
	@go test -race ./... -coverpkg=./... -coverprofile=$(COVER_FILE)
	@go tool cover -html=$(COVER_FILE)
	@rm -f $(COVER_FILE)

.PHONY: lint
lint: tools ## Check the project with lint
	@golangci-lint run -v --fix

.PHONY: check
check: lint test ## Check project with static checks and unit tests
