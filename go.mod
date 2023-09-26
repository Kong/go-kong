module github.com/kong/go-kong

go 1.19

replace github.com/imdario/mergo v0.3.12 => github.com/Kong/mergo v0.3.13

// Contains https://github.com/Kong/go-kong/pull/302 which introduced
// a bug with filling the config defaults:
// https://github.com/Kong/go-kong/issues/307
retract v0.39.1

require (
	github.com/google/go-cmp v0.5.9
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.3.1
	github.com/imdario/mergo v0.3.12
	github.com/kong/semver/v4 v4.0.1
	github.com/mitchellh/mapstructure v1.5.0
	github.com/stretchr/testify v1.8.4
	github.com/tidwall/gjson v1.17.0
	k8s.io/code-generator v0.28.2
	sigs.k8s.io/yaml v1.3.0
)

require (
	github.com/Masterminds/goutils v1.1.1 // indirect
	github.com/Masterminds/semver v1.5.0 // indirect
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/emicklei/go-restful/v3 v3.9.0 // indirect
	github.com/go-logr/logr v1.2.4 // indirect
	github.com/go-openapi/jsonpointer v0.19.6 // indirect
	github.com/go-openapi/jsonreference v0.20.2 // indirect
	github.com/go-openapi/swag v0.22.3 // indirect
	github.com/go-task/slim-sprig v2.20.0+incompatible // indirect
	github.com/golang/protobuf v1.5.3 // indirect
	github.com/google/gnostic-models v0.6.8 // indirect
	github.com/google/gofuzz v1.2.0 // indirect
	github.com/huandu/xstrings v1.4.0 // indirect
	github.com/josharian/intern v1.0.0 // indirect
	github.com/json-iterator/go v1.1.12 // indirect
	github.com/mailru/easyjson v0.7.7 // indirect
	github.com/modern-go/concurrent v0.0.0-20180306012644-bacd9c7ef1dd // indirect
	github.com/modern-go/reflect2 v1.0.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	github.com/spf13/pflag v1.0.5 // indirect
	github.com/tidwall/match v1.1.1 // indirect
	github.com/tidwall/pretty v1.2.0 // indirect
	golang.org/x/mod v0.10.0 // indirect
	golang.org/x/sys v0.10.0 // indirect
	golang.org/x/tools v0.8.0 // indirect
	google.golang.org/protobuf v1.30.0 // indirect
	gopkg.in/yaml.v2 v2.4.0 // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
	k8s.io/gengo v0.0.0-20220902162205-c0856e24416d // indirect
	k8s.io/klog/v2 v2.100.1 // indirect
	k8s.io/kube-openapi v0.0.0-20230717233707-2695361300d9 // indirect
	sigs.k8s.io/structured-merge-diff/v4 v4.2.3 // indirect
)
