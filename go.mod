module github.com/kong/go-kong

go 1.16

replace github.com/imdario/mergo v0.3.12 => github.com/Kong/mergo v0.3.13

require (
	github.com/blang/semver/v4 v4.0.0
	github.com/google/go-cmp v0.5.7
	github.com/google/go-querystring v1.1.0
	github.com/google/uuid v1.3.0
	github.com/imdario/mergo v0.3.12
	github.com/mitchellh/mapstructure v1.4.3
	github.com/stretchr/testify v1.7.0
	github.com/tidwall/gjson v1.14.0
	k8s.io/code-generator v0.23.5
)
