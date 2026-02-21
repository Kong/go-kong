module github.com/kong/go-kong

go 1.24.6

replace github.com/imdario/mergo v0.3.12 => github.com/Kong/mergo v0.3.13

// Contains https://github.com/Kong/go-kong/pull/302 which introduced
// a bug with filling the config defaults:
// https://github.com/Kong/go-kong/issues/307
retract v0.39.1

require (
	github.com/google/go-cmp v0.7.0
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.6.0
	github.com/imdario/mergo v0.3.16
	github.com/kong/semver/v4 v4.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/samber/lo v1.50.0
	github.com/stretchr/testify v1.10.0
	github.com/tidwall/gjson v1.18.0
	k8s.io/code-generator v0.33.3
	sigs.k8s.io/yaml v1.5.0
)

require (
	github.com/davecgh/go-spew v1.1.2-0.20180830191138-d8f796af33cc // indirect
	github.com/emicklei/go-restful/v3 v3.12.0 // indirect
	github.com/go-logr/logr v1.4.2 // indirect
	github.com/go-openapi/jsonpointer v0.21.0 // indirect
	github.com/go-openapi/jsonreference v0.21.0 // indirect
	github.com/go-openapi/swag v0.23.0 // indirect
	github.com/google/gnostic-models v0.6.9 // indirect
	github.com/google/pprof v0.0.0-20241210010833-40e02aabc2ad // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/pmezard/go-difflib v1.0.1-0.20181226105442-5d4384ee4fb2 // indirect
	github.com/spf13/pflag v1.0.6 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	go.yaml.in/yaml/v2 v2.4.2 // indirect
	golang.org/x/mod v0.24.0 // indirect
	golang.org/x/sync v0.15.0 // indirect
	golang.org/x/sys v0.34.0 // indirect
	golang.org/x/text v0.24.0 // indirect
	golang.org/x/tools v0.32.0 // indirect
	google.golang.org/protobuf v1.36.6 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/gengo/v2 v2.0.0-20250207200755-1244d31929d7 // indirect
	k8s.io/klog/v2 v2.130.1 // indirect
	k8s.io/kube-openapi v0.0.0-20250318190949-c8a335a9a2ff // indirect
	sigs.k8s.io/randfill v1.0.0 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.7.0 // indirect
)
