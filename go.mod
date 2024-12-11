module github.com/kong/go-kong

go 1.22.0

toolchain go1.22.3

replace github.com/imdario/mergo v0.3.12 => github.com/Kong/mergo v0.3.13

// Contains https://github.com/Kong/go-kong/pull/302 which introduced
// a bug with filling the config defaults:
// https://github.com/Kong/go-kong/issues/307
retract v0.39.1

require (
	github.com/google/go-cmp v0.6.0
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.6.0
	github.com/imdario/mergo v0.3.12
	github.com/kong/semver/v4 v4.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/stretchr/testify v1.10.0
	github.com/tidwall/gjson v1.18.0
	k8s.io/code-generator v0.31.4
	sigs.k8s.io/yaml v1.4.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.11.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.4 // indirect
	github.com/golang/protobuf v1.5.4 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/mod v0.17.0 // indirect
	golang.org/x/sync v0.7.0 // indirect
	golang.org/x/text v0.16.0 // indirect
	golang.org/x/tools v0.21.1-0.20240508182429-e35e4ccd0d2d // indirect
	google.golang.org/protobuf v1.34.2 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/gengo/v2 v2.0.0-20240228010128-51d4e06bde70 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20240228011516-70dd3763d340 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.4.1 // indirect
)
