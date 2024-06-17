#!/bin/bash -e

readonly GENERATED_FILE=zz_generated.deepcopy.go
readonly GENERATED_FILE_NEW=zz_generated_new.deepcopy.go

go install k8s.io/code-generator/cmd/deepcopy-gen
trap "rm -f kong/${GENERATED_FILE_NEW}" EXIT

deepcopy-gen \
  -v 2 \
  --bounding-dirs kong \
  --output-file ${GENERATED_FILE_NEW} \
  --go-header-file hack/header-template.go.tmpl \
  ./kong

diff -Naur \
  kong/${GENERATED_FILE_NEW} \
  kong/${GENERATED_FILE}
