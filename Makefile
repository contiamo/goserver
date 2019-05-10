# Setup name variables for the package/tool
.PROJECT_ROOT=$(shell pwd)
.BIN_DIR=$(.PROJECT_ROOT)/bin

################################
################################
.PHONY: help
help:
	@grep -E '^[a-zA-Z_-]+:.*?## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'


.PHONY: env
env: ## Print debug information about your local environment
	@echo git: $(shell git version)
	@echo go: $(shell go version)
	@echo golint: $(shell which golint)
	@echo gofmt: $(shell which gofmt)
	@echo staticcheck: $(shell which staticcheck)


.PHONY: changelog
changelog: ## Print git hitstory based changelog
	@git --no-pager log --no-merges --pretty=format:"%h : %s (by %an)" $(shell git describe --tags --abbrev=0)...HEAD
	@echo ""


################################
################################
.PHONY: lint
lint: ## Verifies `golint` passes
	@echo "+ $@"
	@golint -set_exit_status $(shell go list ./... | grep -v "example" | grep -v "grpc/test")

.PHONY: fmt
fmt: $(shell find ./grpc ./http ./aes) ## Verifies all files have been `gofmt`ed
	@echo "+ $@"
	@gofmt -s -l . | tee /dev/stderr

.PHONY: test
test: lint
	go test -cover ./...
