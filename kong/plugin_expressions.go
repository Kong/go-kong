package kong

import (
	"encoding/json"
)

type PluginExpressions map[string]interface{}

func (in PluginExpressions) DeepCopyInto(out *PluginExpressions) {
	b, _ := json.Marshal(&in)
	_ = json.Unmarshal(b, out)
}

func (in PluginExpressions) DeepCopy() PluginExpressions {
	if in == nil {
		return nil
	}

	out := new(PluginExpressions)
	in.DeepCopyInto(out)
	return *out
}
