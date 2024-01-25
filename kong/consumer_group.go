package kong

// ConsumerGroupObject represents a ConsumerGroup in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupObject struct {
	ConsumerGroup *ConsumerGroup         `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	Consumers     []*Consumer            `json:"consumers,omitempty" yaml:"consumers,omitempty"`
	Plugins       []*ConsumerGroupPlugin `json:"plugins,omitempty" yaml:"plugins,omitempty"`
}

// ConsumerGroup represents a ConsumerGroup in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroup struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// ConsumerGroupConsumer represents a ConsumerGroupConsumer in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupConsumer struct {
	Consumer      *Consumer      `json:"consumer,omitempty" yaml:"consumer,omitempty"`
	ConsumerGroup *ConsumerGroup `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	CreatedAt     *int64         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
}

// ConsumerGroupRLA represents a ConsumerGroupRLA in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupRLA struct {
	ConsumerGroup *string       `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	Config        Configuration `json:"config,omitempty" yaml:"config,omitempty"`
	Plugin        *string       `json:"plugin,omitempty" yaml:"plugin,omitempty"`
}

// ConsumerGroupPlugin represents a ConsumerGroupPlugin in Kong.
// +k8s:deepcopy-gen=true
type ConsumerGroupPlugin struct {
	ID            *string        `json:"id,omitempty" yaml:"id,omitempty"`
	Name          *string        `json:"name,omitempty" yaml:"name,omitempty"`
	CreatedAt     *int64         `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	Config        Configuration  `json:"config,omitempty" yaml:"config,omitempty"`
	ConsumerGroup *ConsumerGroup `json:"consumer_group,omitempty" yaml:"consumer_group,omitempty"`
	ConfigSource  *string        `json:"_config,omitempty" yaml:"_config,omitempty"`
}

// FriendlyName returns the endpoint key name or ID.
func (cg *ConsumerGroup) FriendlyName() string {
	if cg.Name != nil {
		return *cg.Name
	}
	if cg.ID != nil {
		return *cg.ID
	}
	return ""
}
