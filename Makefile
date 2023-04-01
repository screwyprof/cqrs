# This repo's root import path.
PKG := gitlab.com/screwyprof/cqrs

## DO NOT EDIT BELLOW THIS LINE
SHELL := bash

GO_FILES = $(shell find . -name "*.go" | uniq)
GO_PACKAGES = $(shell go list ./... | tr '\n', ',')
LOCAL_PACKAGES="github.com/screwyprof/"

IGNORE_COVERAGE_FOR= -e .*test

OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
MAKE_COLOR=\033[33;01m%-20s\033[0m

all: tools lint test ## install tools, lint and test

deps: ## install dependencies
	@echo -e "$(OK_COLOR)--> Downloading go.mod dependencies$(NO_COLOR)"
	@go mod download

tools: ## install dev tools, linters, code generators, etc..
	@echo -e "$(OK_COLOR)--> Installing tools from tools/tools.go$(NO_COLOR)"
	@export GOBIN=$$PWD/tools/bin; export PATH=$$GOBIN:$$PATH; cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

lint: ## run linters
	@echo -e "$(OK_COLOR)--> Running linters$(NO_COLOR)"
	@tools/bin/golangci-lint run

test: ## run  tests
	@echo -e "$(OK_COLOR)--> Running unit tests$(NO_COLOR)"
	go test -v --race --count=1 -coverprofile=coverage.tmp ./...
	@set -euo pipefail && cat coverage.tmp | grep -v $(IGNORE_COVERAGE_FOR) > coverage.out && rm coverage.tmp

coverage: test ## show test coverage report
	@echo -e "$(OK_COLOR)--> Showing test coverage$(NO_COLOR)"
	@go tool cover -func=coverage.out

fmt: ## format go files
	@echo -e "$(OK_COLOR)--> Formatting go files$(NO_COLOR)"
	@go mod tidy
	@go fmt ./...
	@tools/bin/gofumpt -l -w .
	@tools/bin/gci write $(GO_FILES) -s standard  -s default -s "prefix($(LOCAL_PACKAGES))"

clean: ## remove tools
	@echo -e "$(OK_COLOR)--> Clean up$(NO_COLOR)"
	rm -rf $(PWD)/tools/bin
	rm -rf coverage.txt *.out *.tmp

help: ## show this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(MAKE_COLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(MAKE_COLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all deps tools lint test coverage fmt clean help