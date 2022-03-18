build: ## Build
	CGO_ENABLED=0 go build -o ./coinconv cmd/main.go

lint: ## Lint
	golangci-lint run --timeout=3m

req-linter: ## Install linter
	@curl -sfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(GOPATH)/bin v1.31.0

help: ## List of commands
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'