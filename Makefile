# ------------------------------------------------------------------------------
# Configuration - Tooling
# ------------------------------------------------------------------------------

PROJECT_DIR := $(shell dirname $(abspath $(lastword $(MAKEFILE_LIST))))
TOOLS_VERSIONS_FILE = $(PROJECT_DIR)/.tools_versions.yaml
export MISE_DATA_DIR = $(PROJECT_DIR)/bin/

MISE := $(shell which mise)
.PHONY: mise
mise:
	@mise -V >/dev/null || (echo "mise - https://github.com/jdx/mise - not found. Please install it." && exit 1)

.PHONY: tools
tools: golangci-lint yq

# Do not store yq's version in .tools_versions.yaml as it is used to get tool versions.
# renovate: datasource=github-releases depName=mikefarah/yq
YQ_VERSION = 4.43.1
YQ = $(PROJECT_DIR)/bin/installs/yq/$(YQ_VERSION)/bin/yq
.PHONY: yq
yq: mise # Download yq locally if necessary.
	@$(MISE) plugin install --yes -q yq
	@$(MISE) install -q yq@$(YQ_VERSION)

GOLANGCI_LINT_VERSION = $(shell $(YQ) -r '.golangci-lint' < $(TOOLS_VERSIONS_FILE))
GOLANGCI_LINT = $(PROJECT_DIR)/bin/installs/golangci-lint/$(GOLANGCI_LINT_VERSION)/bin/golangci-lint
.PHONY: golangci-lint
golangci-lint: mise yq ## Download golangci-lint locally if necessary.
	@$(MISE) plugin install --yes -q golangci-lint
	@$(MISE) install -q golangci-lint@$(GOLANGCI_LINT_VERSION)

# ------------------------------------------------------------------------------
# CI
# ------------------------------------------------------------------------------


.PHONY: kong.supported-versions
kong.supported-versions:
	@curl -s https://developer.konghq.com/_api/gateway-versions.json | \
	    jq '[.[] | select(.label == null) | select(.endOfLifeDate > (now | strftime("%Y-%m-%d")))] | [.[].tag]'

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
