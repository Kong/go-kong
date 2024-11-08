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

// Uint64 returns a pointer to i.
func Uint64(i uint64) *uint64 {
	return &i
}

// Float64 returns a pointer to f.
func Float64(f float64) *float64 {
	return &f
}

// JSONRawMessage returns a pointer to a json.RawMessage
func JSONRawMessage(s string) *json.RawMessage {
	j := json.RawMessage(s)
	return &j
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

// traverseConfigMap recursively traverses a plugin config
// and returns the value at the specified path.
// The path is represented as a slice of strings, where each string is a key in the map.
//
// If the path is empty, nil is returned.
//
// If the path cannot be fully traversed (e.g., a non-existent key is encountered),
// this function returns nil and an appropriate error.
//
// This function can be helpful to fetch the nested config value from a backward translation
// path provided with deprecated fields.
//
// Example usage:
//
//	configMap := map[string]interface{}{
//		"foo": map[string]interface{}{
//			"bar": 42,
//		},
//	}
//	value, err := traverseConfigMap(configMap, []string{"foo", "bar"})
//	// value comes 42 here
func traverseConfigMap(currentConfigMap map[string]interface{}, path []string) (interface{}, error) {
	if len(path) == 0 {
		return nil, nil
	}

	pathElement := path[0]
	value, ok := currentConfigMap[pathElement]
	if !ok {
		return nil, fmt.Errorf("key %q not found in map", pathElement)
	}

	switch v := value.(type) {
	case map[string]interface{}:
		// Traversing the map recursively, dissecting the path each time
		return traverseConfigMap(v, path[1:])
	default:
		return v, nil
	}
}

// readStringArray converts a gjson.Result (JSON array) into a slice of strings.
// It extracts each element of the JSON array and returns a slice containing their string values.
//
// Parameters:
//
//	gjsonArray (gjson.Result): A JSON array to be converted.
//
// Returns:
//
//	[]string: A slice of strings representing the array elements.
//
// Example:
//
//		gjsonArray := gjson.Parse(`["apple", "banana", "cherry"]`)
//		result := readStringArray(gjsonArray)
//		fmt.Println(result)
//		// Output:
//	 	// ["apple", "banana", "cherry"]
func readStringArray(gjsonArray gjson.Result) []string {
	array := gjsonArray.Array()
	result := make([]string, len(array))
	for i, pathSegment := range array {
		result[i] = pathSegment.String()
	}

	return result
}

// hasCorrespondingShorthandField checks if a given fieldName exists in any shorthand field's list of replacements.
//
// Parameters:
//
//	fieldName (string): The field name to search for.
//	allShorthandFields (map[string][][]string): A map of shorthand field names and their replacement paths.
//
// Returns:
//
//	bool: True if fieldName exists in any shorthand field's replacement paths, otherwise false.
func hasCorrespondingShorthandField(fieldName string, allShorthandFields map[string][][]string) bool {
	for _, v := range allShorthandFields {
		for _, path := range v {
			if len(path) > 0 && path[0] == fieldName {
				return true
			}
		}
	}
	return false
}

// parseReplacedWithPaths extracts replacement paths from a shorthand field's deprecation data.
//
// If "replaced_with" paths are present, they are returned; otherwise, it falls back to "translate_backwards" paths.
//
// Parameters:
//
//	shorthandField (gjson.Result): A JSON object containing deprecation data.
//
// Returns:
//
//	[][]string: A 2D slice of replacement paths.
func parseReplacedWithPaths(shorthandField gjson.Result) [][]string {
	replacements := shorthandField.Get("deprecation.replaced_with.#.path").Array()

	if len(replacements) > 0 {
		paths := make([][]string, len(replacements))
		for i, replacementPath := range replacements {
			paths[i] = readStringArray(replacementPath)
		}

		return paths
	}

	backwardTranslation := shorthandField.Get("translate_backwards")
	if !backwardTranslation.Exists() {
		return nil
	}

	return [][]string{readStringArray(backwardTranslation)}
}

// buildDeprecatedFieldWithReplacementsMap processes a list of shorthand fields and creates a map,
// linking deprecated field names to their corresponding replacement paths.
//
// The function iterates over the given `shorthandFields` JSON object, extracting the deprecated field names,
// and their respective "replaced_with" or "translate_backwards" paths.
// It builds a map where each key is a deprecated field name, and each value is,
// a slice of replacement paths (each path itself is a slice of strings).
//
// Parameters:
//
//	shorthandFields (gjson.Result): A array of JSON objects containing shorthand field data.
//																	Each shorthand field may contain deprecation information,
//																	including the paths that should be used as replacements.
//
// Returns:
//
//	map[string][][]string: A map where each key is a deprecated field name (string),
//	                        and each value is a slice of string slices representing the replacement paths for the field.
//
// Example:
//
//	shorthandFields := gjson.Parse(`[{
//	    "field1": {
//	        "deprecation": {
//	            "replaced_with": [{"path": ["new_field1_1"]}, {"path": ["new_field1_2"]}]
//	        }
//	    },
//	    "field2": {
//	        "deprecation": {
//	            "replaced_with": [{"path": ["new_field2"]}],
//	        }
//	    }
//	}]`)
//
//	result := buildDeprecatedFieldWithReplacementsMap(shorthandFields)
//	fmt.Println(result)
//	// Output:
//	// map[field1:[[new_field1_1] [new_field1_2]] field2:[[new_field2]]]
func buildDeprecatedFieldWithReplacementsMap(shorthandFields gjson.Result) map[string][][]string {
	result := make(map[string][][]string)
	shorthandFields.ForEach(func(_, value gjson.Result) bool {
		field := value.Map()
		for deprecatedFieldName, shorthandFieldConfig := range field {
			if paths := parseReplacedWithPaths(shorthandFieldConfig); paths != nil {
				result[deprecatedFieldName] = paths
			}
		}

		return true
	})

	return result
}

// processDeprecatedFields processes shorthand fields to build a map of deprecated field names
// to their replacement paths, then filters out any shorthand fields not present in the given config.
//
// Parameters:
//
//	shorthandFields (gjson.Result): JSON object containing shorthand field data, including deprecation info.
//	config (Configuration): A map that provides the valid shorthand field names to retain.
//
// Returns:
//
//	map[string][][]string: A map of deprecated field names to replacement paths, with invalid fields removed.
func processDeprecatedFields(shorthandFields gjson.Result, config Configuration) map[string][][]string {
	deprecatedFieldsWithReplacements := buildDeprecatedFieldWithReplacementsMap(shorthandFields)
	filterValidShorthands(deprecatedFieldsWithReplacements, config)

	return deprecatedFieldsWithReplacements
}

// filterValidShorthands removes any shorthand fields from the map that are not found in the given config.
//
// Parameters:
//
//	deprecatedFieldsWithReplacements (map[string][][]string): A map of deprecated field names to replacement paths.
//	config (Configuration): A map of valid shorthand field names. Any field not in this map will be removed.
//
// Returns:
//
//	void: Modifies the input map in place by deleting invalid shorthand fields.
func filterValidShorthands(deprecatedFieldsWithReplacements map[string][][]string, config Configuration) {
	for shorthandFieldName := range deprecatedFieldsWithReplacements {
		if _, ok := config[shorthandFieldName]; !ok {
			delete(deprecatedFieldsWithReplacements, shorthandFieldName)
		}
	}
}

// filterDeprecatedFieldsForGivenField filters out deprecated fields whose replacement paths do not include
// the specified fieldName. It returns a map of deprecated field names and their valid replacement paths.
//
// Parameters:
//
//	fieldName (string): The field name to filter for in the replacement paths.
//	deprecatedFieldsWithReplacements (map[string][][]string): A map of deprecated field names to replacement paths.
//
// Returns:
//
//	map[string][][]string: A filtered map of deprecated field names to replacement paths that include given fieldName.
func filterDeprecatedFieldsForGivenField(
	fieldName string,
	deprecatedFieldsWithReplacements map[string][][]string,
) map[string][][]string {
	result := make(map[string][][]string)
	for deprecatedFieldName, replacedWithPaths := range deprecatedFieldsWithReplacements {
		result[deprecatedFieldName] = make([][]string, 0)
		for _, path := range replacedWithPaths {
			if len(path) > 0 && path[0] == fieldName {
				result[deprecatedFieldName] = append(result[deprecatedFieldName], path[1:])
			}
		}

		if len(result[deprecatedFieldName]) == 0 {
			delete(result, deprecatedFieldName)
		}
	}

	return result
}

type FillRecordOptions struct {
	FillDefaults bool
	FillAuto     bool
}

// fills the config record with default values
func fillConfigRecord(
	schema gjson.Result,
	config Configuration,
	parentShorthands map[string][][]string,
	opts FillRecordOptions,
) Configuration {
	res := config.DeepCopy()
	configFields := schema.Get("fields")
	// Fetch deprecated fields
	shortHandFields := schema.Get("shorthand_fields")
	defaultRecordValue := schema.Get("default")

	// Build a { deprecatedField : [replacements]} map for easier lookup && merge it with the map received
	// from previous recursive call. This map will be used in hasCorrespondingShorthandField
	deprecatedFieldsWithReplacements := processDeprecatedFields(shortHandFields, config)
	for k, v := range parentShorthands {
		deprecatedFieldsWithReplacements[k] = v
	}
	configFields.ForEach(func(_, value gjson.Result) bool {
		// get the key name
		ms := value.Map()
		fname := ""
		for k := range ms {
			fname = k
			break
		}

		if fname == "config" {
			currentShorthands := filterDeprecatedFieldsForGivenField(fname, deprecatedFieldsWithReplacements)
			newConfig := fillConfigRecord(value.Get(fname), config, currentShorthands, opts)
			res = newConfig
			return true
		}

		// check if key is already set in the config
		if v, ok := config[fname]; ok {
			// the field is set. If it's not a map, then
			// the field is fully set. If it's a map, we
			// need to make sure that all fields are properly
			// filled.
			//
			// some fields are defined as arbitrary maps,
			// containing a 'keys' and 'values' subfields.
			// in this case, the map is already fully set.
			switch v.(type) {
			case map[string]interface{}:
				keys := value.Get(fname + ".keys")
				values := value.Get(fname + ".values")
				if keys.Exists() && values.Exists() {
					// an arbitrary map, field is already set.
					return true
				}
			case []interface{}:
				if value.Get(fname).Get("elements.type").String() != "record" &&
					config[fname] != nil {
					// Non nil array with elements which are not of type record
					// this means field is already set
					return true
				}
			default:
				// not a map, field is already set.
				return true
			}
		}

		ftype := value.Get(fname + ".type")
		frequired := value.Get(fname + ".required")
		// Recursively fill defaults only if the field is either required or a subconfig is provided
		if ftype.String() == "record" && (config[fname] != nil || (frequired.Exists() && frequired.Bool())) {
			var fieldConfig Configuration
			subConfig := config[fname]
			switch subConfig.(type) {
			case nil, []interface{}:
				// We can encounter an empty array here due to an incorrect yaml
				// setting or an empty subconfig (like acme.storage_config.kong).
				// This should be treated as nil case.
				// TODO: https://konghq.atlassian.net/browse/KOKO-1125
				fieldConfig = make(map[string]interface{})
			default:
				fieldConfig = subConfig.(map[string]interface{})
			}

			currentShorthands := filterDeprecatedFieldsForGivenField(fname, deprecatedFieldsWithReplacements)
			newSubConfig := fillConfigRecord(value.Get(fname), fieldConfig, currentShorthands, opts)
			// When we are not filling defaults, only assign the subconfig if it's not empty.
			// This is to avoid having records that are assigned empty map values when defaults
			// are not supposed to be filled.
			if opts.FillDefaults && len(newSubConfig) > 0 {
				res[fname] = map[string]interface{}(newSubConfig)
			}
			return true
		}

		// If the config already contains the corresponding shorthand (deprecated) field,
		// we don't need to process it again. This is to avoid overwriting or stubbing the value,
		// which could lead to incorrect results, especially when the new field has a default value.
		if hasCorrespondingShorthandField(fname, deprecatedFieldsWithReplacements) {
			return true
		}

		// Check if field is of type array of records (in Schema)
		// If this array is non-nil and non-empty (in Config), go through all the records in this array and add defaults
		// If the array has only primitives like string/number/boolean then the value is already set
		// If the array is empty or nil, then no defaults need to be set for its elements
		// The same logic should be applied if field is of type set of records (in Schema)
		if ftype.String() == "array" || ftype.String() == "set" {
			if value.Get(fname).Get("elements.type").String() == "record" {
				if config[fname] != nil {
					// Check sub config is of type array and it is non-empty
					if subConfigArray, ok := config[fname].([]interface{}); ok && len(subConfigArray) > 0 {
						processedSubConfigArray := make([]interface{}, len(subConfigArray))

						for i, configRecord := range subConfigArray {
							// Check if element is of type record,
							// if it is, set default values by recursively calling `fillConfigRecord`
							if configRecordMap, ok := configRecord.(map[string]interface{}); ok {
								processedConfigRecord := fillConfigRecord(
									value.Get(fname).Get("elements"),
									configRecordMap,
									deprecatedFieldsWithReplacements,
									opts)
								processedSubConfigArray[i] = processedConfigRecord
								continue
							}
							// Element not of type record, keep the value as is
							processedSubConfigArray[i] = configRecord
						}
						res[fname] = processedSubConfigArray
						return true
					}
				}
			}
		}

		auto := value.Get(fname + ".auto")
		isAuto := auto.Exists() && auto.Bool()

		// if this is an auto field and we are not filling auto fields, skip
		if !opts.FillAuto && isAuto {
			return true
		}

		// if this is a normal field and we are not filling defaults, skip
		if !opts.FillDefaults && !isAuto {
			return true
		}

		// Check if the record has a default value for the specified field.
		// If so, use it. If not, fall back to the default value of the field itself.
		if defaultRecordValue.Exists() && defaultRecordValue.Get(fname).Exists() {
			value = defaultRecordValue.Get(fname)
		} else {
			value = value.Get(fname + ".default")
		}

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

	value.ForEach(func(_, value gjson.Result) bool {
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
	case *ConsumerGroupPlugin:
		tmpEntity = &ConsumerGroupPlugin{}
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

func fillConfigRecordDefaultsAutoFields(plugin *Plugin, schema map[string]interface{}, opts FillRecordOptions) error {
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

	plugin.Config = fillConfigRecord(configSchema, plugin.Config, nil, opts)
	if plugin.Protocols == nil {
		plugin.Protocols = getDefaultProtocols(gjsonSchema)
	}
	if plugin.Enabled == nil {
		plugin.Enabled = Bool(true)
	}
	return nil
}

// FillPluginsDefaults ingests plugin's defaults from its schema.
// Takes in a plugin struct and mutate it in place.
func FillPluginsDefaults(plugin *Plugin, schema Schema) error {
	return fillConfigRecordDefaultsAutoFields(plugin, schema, FillRecordOptions{
		FillDefaults: true,
		FillAuto:     true,
	})
}

// same as FillPluginsDefaults but allows configuring whether to fill defaults and auto fields.
func FillPluginsDefaultsWithOpts(plugin *Plugin, schema map[string]interface{}, opts FillRecordOptions) error {
	return fillConfigRecordDefaultsAutoFields(plugin, schema, opts)
}

// deleteAndCollapseMap is a utility function that removes an element from a map
// based on a given path. If removing the element results in an empty map at that
// key, the key itself is also deleted from the parent map.
//
// Parameters:
//   - config: The map from which the element will be removed. The map is of type
//     `map[string]interface{}`.
//   - path: A slice of strings that represents the path to the element in the map.
//     Each element of the path corresponds to a key in the map (or a nested map).
//
// Returns:
//   - The function modifies the `config` map in place. It does not return a value.
//
// Example:
//
//	configMap := map[string]interface{}{
//	    "a": map[string]interface{}{
//	        "a1": 1,
//	    },
//	    "b": 456,
//	}
//
//	deleteAndCollapseMap(configMap, []string{"a", "a1"})
//	fmt.Println(configMap)
//	// Output:
//	// map[b:456]
func deleteAndCollapseMap(config map[string]interface{}, path []string) {
	if len(path) == 0 {
		return
	}

	key := path[0]
	if len(path) == 1 {
		delete(config, key)
		return
	}

	if nested, ok := config[key].(map[string]interface{}); ok {
		deleteAndCollapseMap(nested, path[1:])
		if len(nested) == 0 {
			delete(config, key)
		}
	}
}

// This function handles the relationship between deprecated and new plugin configuration values.
// We consider the following scenarios:
//
// - **Scenario 1**: Both old and new values are present.
//   - Action: No adjustment needed for the old plugin configuration.
//
// - **Scenario 2**: The new key is missing in the new plugin configuration.
//   - Action: Delete the deprecated value from both old and new plugin configurations.
//
// - **Scenario 3**: The new field exists but was set to `null` due to decK logic.
//   - If the deprecated value is different from `nil`:
//   - Action: Clear the deprecated value in both new and old configurations.
func clearUnmatchingDeprecationsForGivenPath(
	path []string,
	newPluginConfig Configuration,
	oldPluginConfig Configuration,
	acceptNullValue bool,
) {
	newPluginNewFieldValue, _ := traverseConfigMap(newPluginConfig, path)
	if newPluginNewFieldValue == nil {
		if !acceptNullValue {
			deleteAndCollapseMap(newPluginConfig, path)
		}
		deleteAndCollapseMap(oldPluginConfig, path)
	}
}

// clearCurrentLevelUnmatchingDeprecations compares the new and old plugin configurations to handle deprecated fields.
// For each deprecated field in the new configuration, it checks whether the corresponding field is present in
// the old configuration, and adjusts the old configuration if necessary. If the deprecated field is missing
// from the new configuration, it deletes it from the old configuration for consistency.
//
// Parameters:
//
//	newPluginConfig - The updated configuration containing potential deprecated fields.
//	oldPluginConfig - The original configuration that may contain deprecated fields.
//	schema - A JSON schema representing the structure of the configuration, including deprecated fields.
//
// This function mutates the oldPluginConfig to align with the newPluginConfig in regard to deprecated fields.
func clearCurrentLevelUnmatchingDeprecations(
	newPluginConfig Configuration,
	oldPluginConfig Configuration,
	schema *gjson.Result,
) {
	// Fetch deprecated fields
	shortHandFields := schema.Get("shorthand_fields")

	shortHandFields.ForEach(func(_, value gjson.Result) bool {
		field := value.Map()
		for deprecatedFieldName, shorthandFieldConfig := range field {
			if deprecatedFieldValue, ok := newPluginConfig[deprecatedFieldName]; ok {
				// The new plugin configuration contains deprecated field. Verify if the oldPluginConfiguration needs
				// to be adjusted in order to match newPluginConfiguration.
				//
				// Determine if we accept `null` values based on the deprecated field:
				acceptNullValue := deprecatedFieldValue == nil
				for _, path := range parseReplacedWithPaths(shorthandFieldConfig) {
					clearUnmatchingDeprecationsForGivenPath(path, newPluginConfig, oldPluginConfig, acceptNullValue)
				}
			} else {
				// The new plugin configuration does not contain deprecated fields.
				// However, for backwards compatibility, Kong sends deprecated fields in the response.
				// To ensure consistent diffs, we need to Delete deprecated fields from the old plugin configuration
				// that Kong sent us.
				delete(oldPluginConfig, deprecatedFieldName)
			}
		}

		return true
	})
}

// traverseConfigurationsAndExecute iterates over the fields in two plugin configurations (configA and configB),
// and executes the provided function 'f' on nested fields of type "record" that exist in both configurations.
//
// Parameters:
//
//	configA - The first configuration to compare.
//	configB - The second configuration to compare.
//	schema - The schema describing the fields of the configurations.
//	f - A function that will be executed on each pair of matching nested fields of type "record" in configA and configB.
func traverseConfigurationsAndExecute(
	configA Configuration,
	configB Configuration,
	schema *gjson.Result,
	f func(Configuration, Configuration, *gjson.Result),
) {
	configFields := schema.Get("fields")
	configFields.ForEach(func(_, value gjson.Result) bool {
		field := value.Map()

		for fieldName, fieldConfig := range field {
			if fieldType := fieldConfig.Get("type"); fieldType.String() == "record" {
				var nestedConfigA map[string]interface{}
				if fieldA, ok := configA[fieldName].(map[string]interface{}); ok {
					nestedConfigA = fieldA
				}

				var nestedConfigB map[string]interface{}
				if fieldB, ok := configB[fieldName].(map[string]interface{}); ok {
					nestedConfigB = fieldB
				}

				if nestedConfigA != nil && nestedConfigB != nil {
					f(nestedConfigA, nestedConfigB, &fieldConfig)
				}
			}
		}

		return true
	})
}

// same as ClearUnmatchingDeprecations but this function below is adjusted for recursive use (which is required
// when the schema contains nested records and they can have nested shorthands).
func clearUnmatchingDeprecationsHelper(
	newPluginConfig Configuration,
	oldPluginConfig Configuration,
	schema *gjson.Result,
) {
	clearCurrentLevelUnmatchingDeprecations(newPluginConfig, oldPluginConfig, schema)

	// Recursively walk through configuration to clear any nested unmatching deprecations.
	traverseConfigurationsAndExecute(
		newPluginConfig,
		oldPluginConfig,
		schema,
		clearUnmatchingDeprecationsHelper,
	)
}

// ClearUnmatchingDeprecations is a function that go through a pair of
// configurations: newPlugin and oldPlugin, and by using schema it makes sure those two configurations
// are aligned.
// It does so by removing deprecated or new fields in "oldPlugin" that were not defined in "oldPlugin".
// Furthermore it'll remove new field from newPlugin when it's value is nil and the corresponding deprecated
// value is not nil (in that case we can be sure that "oldPlugin" contains
// specific value for both new and old fields).
func ClearUnmatchingDeprecations(newPlugin *Plugin, oldPlugin *Plugin, schema map[string]interface{}) error {
	jsonb, err := json.Marshal(&schema)
	if err != nil {
		return err
	}
	gjsonSchema := gjson.ParseBytes(jsonb)
	configSchema, err := getConfigSchema(gjsonSchema)
	if err != nil {
		return err
	}

	if newPlugin != nil && oldPlugin != nil {
		clearUnmatchingDeprecationsHelper(newPlugin.Config, oldPlugin.Config, &configSchema)
	}

	return nil
}
