package kong

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"reflect"
	"regexp"
	"strings"

	"github.com/imdario/mergo"
	"github.com/tidwall/gjson"
)

const (
	versionParts = 4
)

// String returns pointer to s.
func String(s string) *string {
	return &s
}

// Bool returns a pointer to b.
func Bool(b bool) *bool {
	return &b
}

// Int returns a pointer to i.
func Int(i int) *int {
	return &i
}

// Float64 returns a pointer to f.
func Float64(f float64) *float64 {
	return &f
}

func isEmptyString(s *string) bool {
	return s == nil || strings.TrimSpace(*s) == ""
}

// StringSlice converts a slice of string to a
// slice of *string
func StringSlice(elements ...string) []*string {
	var res []*string
	for _, element := range elements {
		e := element
		res = append(res, &e)
	}
	return res
}

func stringArrayToString(arr []*string) string {
	if arr == nil {
		return "nil"
	}

	var buf bytes.Buffer
	buf.WriteString("[ ")
	l := len(arr)
	for i, el := range arr {
		buf.WriteString(*el)
		if i != l-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(" ]")
	return buf.String()
}

// headerRoundTripper injects Headers into requests
// made via RT.
type headerRoundTripper struct {
	headers http.Header
	rt      http.RoundTripper
}

// RoundTrip satisfies the RoundTripper interface.
func (t headerRoundTripper) RoundTrip(req *http.Request) (*http.Response, error) {
	req = requestWithHeaders(req, t.headers)
	return t.rt.RoundTrip(req)
}

func requestWithHeaders(req *http.Request, headers http.Header) *http.Request {
	if req == nil {
		return nil
	}
	// fast-path
	if len(headers) == 0 {
		return req
	}
	newRequest := new(http.Request)
	*newRequest = *req
	newRequest.Header = req.Header.Clone()
	for k, values := range headers {
		for _, v := range values {
			newRequest.Header.Add(k, v)
		}
	}
	return newRequest
}

// HTTPClientWithHeaders returns a client which injects headers
// before sending any request.
func HTTPClientWithHeaders(client *http.Client,
	headers http.Header,
) *http.Client {
	var res *http.Client
	if client == nil {
		defaultTransport := http.DefaultTransport.(*http.Transport)
		res = &http.Client{}
		res.Transport = defaultTransport
	} else {
		res = client
		// If a client with nil transport has been provided then set it to http's
		// default transport so that the caller is still able to use it.
		if res.Transport == nil {
			res.Transport = http.DefaultTransport.(*http.Transport)
		}
	}
	res.Transport = headerRoundTripper{
		headers: headers,
		rt:      res.Transport,
	}
	return res
}

// ParseSemanticVersion creates a semantic version from the version
// returned by Kong.
func ParseSemanticVersion(v string) (Version, error) {
	re := regexp.MustCompile(`((?:\d+\.\d+\.\d+)|(?:\d+\.\d+))(?:[\.-](\d+))?(?:\-?(.+)$|$)`)
	m := re.FindStringSubmatch(v)
	if len(m) != versionParts {
		return Version{}, fmt.Errorf("unknown Kong version : '%v'", v)
	}
	if m[2] == "" {
		// Only append zero patch version if major and minor have been detected
		if strings.Count(m[1], ".") == 1 {
			m[2] = ".0"
		}
	} else if strings.Count(m[2], ".") == 0 {
		// Ensure stripped digit is prefixed with a "."
		m[2] = "." + m[2]
	}
	if m[3] != "" {
		if strings.Contains(m[3], "enterprise") {
			// Convert enterprise pre-release to build metadata
			m[3] = "+" + strings.Replace(m[3], "enterprise-edition", "enterprise", 1)
		} else {
			// Keep pre-release information intact
			m[3] = "-" + m[3]
		}
		m[3] = strings.Replace(m[3], ".", "", -1)
	}

	return NewVersion(fmt.Sprintf("%s%s%s", m[1], m[2], m[3]))
}

// VersionFromInfo retrieves the version from the response of root
// or /kong endpoints
func VersionFromInfo(info map[string]interface{}) string {
	version, ok := info["version"]
	if !ok {
		return ""
	}
	return version.(string)
}

func getDefaultProtocols(schema gjson.Result) []*string {
	var res []*string
	fields := schema.Get("fields")

	for _, field := range fields.Array() {
		protocols := field.Map()["protocols"]
		d := protocols.Get("default")
		if d.Exists() {
			for _, v := range d.Array() {
				res = append(res, String(v.String()))
			}
			return res
		}
	}
	return res
}

func getConfigSchema(schema gjson.Result) (gjson.Result, error) {
	fields := schema.Get("fields")

	for _, field := range fields.Array() {
		config := field.Map()["config"]
		if config.Exists() {
			return config, nil
		}
	}
	return schema, fmt.Errorf("no 'config' field found in schema")
}

func fillConfigRecord(schema gjson.Result, config Configuration) Configuration {
	res := config.DeepCopy()
	value := schema.Get("fields")

	value.ForEach(func(key, value gjson.Result) bool {
		// get the key name
		ms := value.Map()
		fname := ""
		for k := range ms {
			fname = k
			break
		}

		if fname == "config" {
			newConfig := fillConfigRecord(value.Get(fname), config)
			res = newConfig
			return true
		}

		// check if key is already set in the config
		if _, ok := config[fname]; ok {
			// yes, don't set it
			return true
		}
		ftype := value.Get(fname + ".type")
		if ftype.String() == "record" {
			subConfig := config[fname]
			if subConfig == nil {
				subConfig = make(map[string]interface{})
			}
			newSubConfig := fillConfigRecord(value.Get(fname), subConfig.(map[string]interface{}))
			res[fname] = map[string]interface{}(newSubConfig)
			return true
		}
		value = value.Get(fname + ".default")
		if value.Exists() {
			res[fname] = value.Value()
		} else {
			// if no default exists, set an explicit nil
			res[fname] = nil
		}
		return true
	})

	return res
}

// flattenDefaultsSchema gets an arbitrarily nested and structured entity schema
// and flattens it, turning it into a map that can be more easily unmarshalled
// into proper entity objects.
//
// This supports both Lua and JSON schemas.
//
// Sample Lua schema input:
//
//	{
//		"fields": [
//	        {
//	            "algorithm": {
//	                "default": "round-robin",
//	                "one_of": ["consistent-hashing", "least-connections", "round-robin"],
//	                "type": "string"
//	            }
//	        }, {
//	            "hash_on": {
//	                "default": "none",
//	                "one_of": ["none", "consumer", "ip", "header", "cookie"],
//	                "type": "string"
//	            }
//	        }, {
//	            "hash_fallback": {
//	                "default": "none",
//	                "one_of": ["none", "consumer", "ip", "header", "cookie"],
//	                "type": "string"
//	            }
//	        },
//	  ...
//	}
//
// Sample JSON schema input:
//
//	{
//		"properties": [
//	        {
//	            "algorithm": {
//	                "default": "round-robin",
//	                "enum": ["consistent-hashing", "least-connections", "round-robin"],
//	                "type": "string"
//	            }
//	        }, {
//	            "hash_on": {
//	                "default": "none",
//	                "enum": ["none", "consumer", "ip", "header", "cookie"],
//	                "type": "string"
//	            }
//	        }, {
//	            "hash_fallback": {
//	                "default": "none",
//	                "enum": ["none", "consumer", "ip", "header", "cookie"],
//	                "type": "string"
//	            }
//	        },
//	  ...
//	}
//
// Sample output:
//
//	{
//		"algorithm": "round-robin",
//		"hash_on": "none",
//		"hash_fallback": "none",
//	 ...
//	}
func flattenDefaultsSchema(schema gjson.Result) Schema {
	fields := schema.Get("fields")
	if fields.Exists() {
		return flattenLuaSchema(fields)
	}
	properties := schema.Get("properties")
	if properties.Exists() {
		return flattenJSONSchema(properties)
	}
	return Schema{}
}

func flattenJSONSchema(value gjson.Result) Schema {
	results := Schema{}

	value.ForEach(func(key, value gjson.Result) bool {
		name := key.String()

		ftype := value.Get("type")
		// when type==object and additionalProperties==false, the object
		// represents either a foreign relationship or a map entry.
		// In both cases, defaults don't need to be injected.
		additionalProperties := value.Get("additionalProperties")
		if ftype.String() == "object" &&
			(!additionalProperties.Exists() ||
				(additionalProperties.Exists() &&
					additionalProperties.Bool())) {
			newSubConfig := flattenDefaultsSchema(value)
			results[name] = newSubConfig
			return true
		}
		value = value.Get("default")
		if value.Exists() {
			results[name] = value.Value()
		} else {
			results[name] = nil
		}
		return true
	})

	return results
}

func flattenLuaSchema(value gjson.Result) Schema {
	results := Schema{}

	value.ForEach(func(key, value gjson.Result) bool {
		// get the key name
		ms := value.Map()
		fname := ""
		for k := range ms {
			fname = k
			break
		}

		ftype := value.Get(fname + ".type")
		if ftype.String() == "record" {
			newSubConfig := flattenDefaultsSchema(value.Get(fname))
			results[fname] = newSubConfig
			return true
		}
		value = value.Get(fname + ".default")
		if value.Exists() {
			results[fname] = value.Value()
		} else {
			results[fname] = nil
		}
		return true
	})

	return results
}

func getDefaultsObj(schema Schema) ([]byte, error) {
	jsonSchema, err := json.Marshal(&schema)
	if err != nil {
		return nil, err
	}
	gjsonSchema := gjson.ParseBytes((jsonSchema))
	defaults := flattenDefaultsSchema(gjsonSchema)
	jsonSchemaWithDefaults, err := json.Marshal(&defaults)
	if err != nil {
		return nil, err
	}
	return jsonSchemaWithDefaults, nil
}

type zeroValueTransformer struct{}

func (t zeroValueTransformer) Transformer(typ reflect.Type) func(dst, src reflect.Value) error {
	if typ == reflect.TypeOf(Bool(true)) || typ == reflect.TypeOf(Int(0)) {
		return func(dst, src reflect.Value) error {
			if dst.CanSet() {
				if src.IsNil() {
					src.Set(dst)
				}
			}
			return nil
		}
	}
	return nil
}

// FillEntityDefaults ingests entities' defaults from their schema.
func FillEntityDefaults(entity interface{}, schema Schema) error {
	if schema == nil {
		return fmt.Errorf("filling defaults for '%T': provided schema is nil", entity)
	}
	var tmpEntity interface{}
	switch entity.(type) {
	case *Target:
		tmpEntity = &Target{}
	case *Service:
		tmpEntity = &Service{}
	case *Route:
		tmpEntity = &Route{}
	case *Upstream:
		tmpEntity = &Upstream{}
	default:
		return fmt.Errorf("unsupported entity: '%T'", entity)
	}
	defaults, err := getDefaultsObj(schema)
	if err != nil {
		return fmt.Errorf("parse schema for defaults: %w", err)
	}
	if err := json.Unmarshal(defaults, &tmpEntity); err != nil {
		return fmt.Errorf("unmarshal entity with defaults: %w", err)
	}
	if err := mergo.Merge(
		entity, tmpEntity, mergo.WithTransformers(zeroValueTransformer{}),
	); err != nil {
		return fmt.Errorf("merge entity with its defaults: %w", err)
	}
	return nil
}

// FillPluginsDefaults ingests plugin's defaults from its schema.
// Takes in a plugin struct and mutate it in place.
func FillPluginsDefaults(plugin *Plugin, schema Schema) error {
	jsonb, err := json.Marshal(&schema)
	if err != nil {
		return err
	}
	gjsonSchema := gjson.ParseBytes((jsonb))
	configSchema, err := getConfigSchema(gjsonSchema)
	if err != nil {
		return err
	}
	if plugin.Config == nil {
		plugin.Config = make(Configuration)
	}
	plugin.Config = fillConfigRecord(configSchema, plugin.Config)
	if plugin.Protocols == nil {
		plugin.Protocols = getDefaultProtocols(gjsonSchema)
	}
	if plugin.Enabled == nil {
		plugin.Enabled = Bool(true)
	}
	return nil
}
