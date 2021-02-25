OK_COLOR=\033[32;01m
NO_COLOR=\033[0m
MAKE_COLOR=\033[33;01m%-20s\033[0m

all: tools lint test ## install tools, lint and test

deps: ## install dependencies
	@echo "$(OK_COLOR)--> Download go.mod dependencies$(NO_COLOR)"
	go mod download
	go mod vendor

tools: ## install dev tools, linters, code generators, etc..
	@echo "$(OK_COLOR)--> Installing tools from tools/tools.go$(NO_COLOR)"
	@export GOBIN=$$PWD/tools/bin; export PATH=$$GOBIN:$$PATH; cat tools/tools.go | grep _ | awk -F'"' '{print $$2}' | xargs -tI % go install %

lint: ## run linters
	@echo "$(OK_COLOR)--> Running linters$(NO_COLOR)"
	tools/bin/golangci-lint run

test: test-unit ## run all tests

test-unit: ## run unit tests
	@echo "$(OK_COLOR)--> Running unit tests$(NO_COLOR)"
	go test --race --count=1 ./...

test-coverage: ## run all tests with coverage
	@echo "$(OK_COLOR)--> Generating code coverage$(NO_COLOR)"
	tools/coverage.sh

fmt: ## format go files
	@echo "$(OK_COLOR)--> Formatting go files$(NO_COLOR)"
	gofumpt -l -w .

clean: ## remove tools
	@echo "$(OK_COLOR)--> Clean up$(NO_COLOR)"
	rm -rf $(PWD)/tools/bin
	rm coverage.txt c.out

help: ## show this help screen
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m\033[0m\n"} /^[a-zA-Z_-]+:.*?##/ { printf "  $(MAKE_COLOR) %s\n", $$1, $$2 } /^##@/ { printf "\n$(MAKE_COLOR)\n", substr($$0, 5) } ' $(MAKEFILE_LIST)

# To avoid unintended conflicts with file names, always add to .PHONY
# unless there is a reason not to.
# https://www.gnu.org/software/make/manual/html_node/Phony-Targets.html
.PHONY: all deps tools lint test test-unit test-coverage fmt clean help