package kong

import (
	"encoding/json"
)

// Configuration represents a config of a plugin/partial in Kong.
type Configuration map[string]interface{}

// DeepCopyInto copies the receiver, writing into out. in must be non-nil.
func (in Configuration) DeepCopyInto(out *Configuration) {
	// Resorting to JSON since interface{} cannot be DeepCopied easily.
	// This could be replaced using reflection-fu.
	// XXX Ignoring errors
	b, _ := json.Marshal(&in)
	_ = json.Unmarshal(b, out)
}

// DeepCopy copies the receiver, creating a new Configuration.
func (in Configuration) DeepCopy() Configuration {
	if in == nil {
		return nil
	}
	out := new(Configuration)
	in.DeepCopyInto(out)
	return *out
}
