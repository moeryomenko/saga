include tools.mk
include check.mk
include int.mk

.PHONY: help
help: ## Prints this help message
	@echo "Commands:"
	@fgrep -h '##' $(MAKEFILE_LIST) \
		| fgrep -v fgrep \
		| sort \
		| grep -E '^[a-zA-Z_-]+:.*?## .*$$' \
		| awk 'BEGIN {FS = ":.*?## "}; {printf "  \033[36m%-30s\033[0m %s\n", $$1, $$2}'
