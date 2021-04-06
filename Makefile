.DEFAULT_GOAL := test-all

.PHONY: test-all
test-all: lint test

.PHONY: test
test:
	go test -v ./...

.PHONY: test-enterprise
test:
	go test -tags=enterprise -v ./...

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: verify-codegen
verify-codegen:
	./hack/verify-deepcopy-gen.sh

.PHONY: update-codegen
update-codegen:
	./hack/update-deepcopy-gen.sh

.PHONY: setup-kong
setup-kong:
	bash .ci/setup_kong.sh

.PHONY: setup-lint
setup-lint:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $(go env GOPATH)/bin v1.25.0

.PHONY: coverage
coverage:
	go test -tags=enterprise -race -v -count=1 -coverprofile=coverage.out.tmp ./...
	# ignoring generated code for coverage
	grep -E -v 'generated.deepcopy.go' coverage.out.tmp > coverage.out
	rm -f coverage.out.tmp

.PHONY: upload-coverage
upload-coverage:
	bash <(curl -s https://codecov.io/bash)