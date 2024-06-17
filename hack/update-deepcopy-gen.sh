#!/bin/bash -e

readonly GENERATED_FILE=zz_generated.deepcopy.go

go install k8s.io/code-generator/cmd/deepcopy-gen

deepcopy-gen \
  -v 2 \
  --bounding-dirs kong \
  --output-file ${GENERATED_FILE} \
  --go-header-file hack/header-template.go.tmpl \
  ./kong
