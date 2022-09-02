# Tools
.PHONY: tools
tools: ## Install all needed tools, e.g. for static checks
	@echo Installing tools from tools.txt
	@grep '@' tools.txt | xargs -tI % go install %

.PHONE: gen
gen: tools ## Generate projects files and components.
	@oapi-codegen --config api/order/config.yaml api/order/api.yaml > internal/order/infrastructure/api/http.gen.go

.PHONY: deps
deps: ## Manage go mod dependencies, beautify go.mod and go.sum files
	@go-mod-upgrade
	@go mod tidy
