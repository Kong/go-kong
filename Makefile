# ------------------------------------------------------------------------------
# Configuration - Tooling
# ------------------------------------------------------------------------------

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))

.PHONY: _download_tool
_download_tool:
	(cd third_party && go mod tidy && \
		GOBIN=$(PROJECT_DIR)/bin go generate -tags=third_party ./$(TOOL).go )

.PHONY: tools
tools: golangci-lint

GOLANGCI_LINT = $(PROJECT_DIR)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: ## Download golangci-lint locally if necessary.
	@$(MAKE) _download_tool TOOL=golangci-lint

# ------------------------------------------------------------------------------
# Testing
# ------------------------------------------------------------------------------

.PHONY: test
test:
	go test -v ./...

.PHONY: test-enterprise
test-enterprise:
	go test -tags=enterprise -v ./...

.PHONY: lint
lint: golangci-lint
	$(GOLANGCI_LINT) run -v ./...

.PHONY: verify-codegen
verify-codegen:
	./hack/verify-deepcopy-gen.sh

.PHONY: update-codegen
update-codegen:
	./hack/update-deepcopy-gen.sh

.PHONY: setup-kong-dbless
setup-kong-dbless:
	bash .ci/setup_kong.sh --dbless

.PHONY: setup-kong-postgres
setup-kong-postgres:
	bash .ci/setup_kong.sh --postgres

.PHONY: setup-kong-ee
setup-kong-ee:
	bash .ci/setup_kong_ee.sh

.PHONY: teardown
teardown:
	bash .ci/teardown.sh

.PHONY: test-coverage-enterprise
test-coverage-enterprise:
	go test -tags=enterprise -race -v -count=1 -coverprofile=coverage.out.tmp ./...
	# ignoring generated code for coverage
	grep -E -v 'generated.deepcopy.go' coverage.out.tmp > coverage.out
	rm -f coverage.out.tmp

.PHONY: test-coverage
test-coverage:
	go test -race -v -count=1 -coverprofile=coverage.out.tmp ./...
	# ignoring generated code for coverage
	grep -E -v 'generated.deepcopy.go' coverage.out.tmp > coverage.out
	rm -f coverage.out.tmp
