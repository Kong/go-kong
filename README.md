# go-kong

Go bindings for Kong's Admin API

[![GoDoc](https://godoc.org/github.com/kong/go-kong?status.svg)](https://godoc.org/github.com/kong/go-kong/kong)
[![codecov](https://codecov.io/gh/Kong/go-kong/branch/main/graph/badge.svg?token=OLN3HEOIVP)](https://codecov.io/gh/Kong/go-kong)
[![Go Report Card](https://goreportcard.com/badge/github.com/kong/go-kong)](https://goreportcard.com/report/github.com/kong/go-kong)

[![Build Status](https://github.com/kong/go-kong/actions/workflows/integration-test.yaml/badge.svg)](https://github.com/Kong/go-kong/actions/workflows/integration-test.yaml)
[![Build Status](https://github.com/kong/go-kong/actions/workflows/integration-test-nightly.yaml/badge.svg)](https://github.com/Kong/go-kong/actions/workflows/integration-test-nightly.yaml)
[![Build Status](https://github.com/kong/go-kong/actions/workflows/integration-test-enterprise.yaml/badge.svg)](https://github.com/Kong/go-kong/actions/workflows/integration-test-enterprise.yaml)
[![Build Status](https://github.com/kong/go-kong/actions/workflows/integration-test-enterprise-nightly.yaml/badge.svg)](https://github.com/Kong/go-kong/actions/workflows/integration-test-enterprise-nightly.yaml)

## Importing

```shell
go get github.com/kong/go-kong/kong
```

## Compatibility

`go-kong` is compatible with Kong 2.x and 3.x.
Semantic versioning is followed for versioning `go-kong`.

## Generators

Some code in this repo such as `kong/zz_generated.deepcopy.go` is generated
from API types (see `kong/types.go`).

After making a change to an API type you can run the generators with:

```shell
./hack/update-deepcopy-gen.sh
```

## License

go-kong is licensed with Apache License Version 2.0.
Please read the LICENSE file for more details.
