package kong

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/mitchellh/mapstructure"
)

// -----------------------------------------------------------------------------
// Kong Listeners - Public Types
// -----------------------------------------------------------------------------

// ProxyListener is a configured listener on the Kong Gateway for L7 routing.
type ProxyListener struct {
	SSL           bool   `json:"ssl"             mapstructure:"ssl"`
	Listener      string `json:"listener"        mapstructure:"listener"`
	Port          int    `json:"port"            mapstructure:"port"`
	Bind          bool   `json:"bind"            mapstructure:"bind"`
	IP            string `json:"ip"              mapstructure:"ip"`
	HTTP2         bool   `json:"http2"           mapstructure:"http2"`
	ProxyProtocol bool   `json:"proxy_protocol"  mapstructure:"proxy_protocol"`
	Deferred      bool   `json:"deferred"        mapstructure:"deferred"`
	ReusePort     bool   `json:"reuseport"       mapstructure:"reuseport"`
	Backlog       bool   `json:"backlog=%d+"     mapstructure:"backlog=%d+"`
}

// StreamListener is a configured listener on the Kong Gateway for L4 routing.
type StreamListener struct {
	UDP           bool   `json:"udp"             mapstructure:"udp"`
	SSL           bool   `json:"ssl"             mapstructure:"ssl"`
	ProxyProtocol bool   `json:"proxy_protocol"  mapstructure:"proxy_protocol"`
	IP            string `json:"ip"              mapstructure:"ip"`
	Listener      string `json:"listener"        mapstructure:"listener"`
	Port          int    `json:"port"            mapstructure:"port"`
	Bind          bool   `json:"bind"            mapstructure:"bind"`
	ReusePort     bool   `json:"reuseport"       mapstructure:"reuseport"`
	Backlog       bool   `json:"backlog=%d+"     mapstructure:"backlog=%d+"`
}

// -----------------------------------------------------------------------------
// Kong Listeners - Client Methods
// -----------------------------------------------------------------------------

// Listeners returns the proxy_listeners and stream_listeners that are currently configured in the
// Kong root as convenient native types rather than JSON or unstructured.
func (c *Client) Listeners(ctx context.Context) ([]ProxyListener, []StreamListener, error) {
	rootJSON, err := c.RootJSON(ctx)
	if err != nil {
		return nil, nil, fmt.Errorf("couldn't get root JSON when trying to determine listeners: %w", err)
	}

	root := dataPlaneConfigWrapper{}
	if err := json.Unmarshal(rootJSON, &root); err != nil {
		return nil, nil, fmt.Errorf("couldn't decode root JSON when trying to determine listeners: %w", err)
	}

	return root.Config.ProxyListeners, root.Config.StreamListeners, nil
}

// -----------------------------------------------------------------------------
// Kong Listeners - Private Wrapper Types
// -----------------------------------------------------------------------------

type dataPlaneConfigWrapper struct {
	Config dataPlaneConfig `json:"configuration"`
}

type dataPlaneConfig struct {
	ProxyListeners  []ProxyListener  `json:"proxy_listeners"`
	StreamListeners []StreamListener `json:"stream_listeners"`
}

// UnmarshalJSON implements custom JSON unmarshaling for this type which must
// be done because the Kong Admin API will return empty objects when a list
// is empty which will confuse the default unmarshaler.
func (d *dataPlaneConfig) UnmarshalJSON(data []byte) error {
	wrapper := make(map[string]interface{})
	if err := json.Unmarshal(data, &wrapper); err != nil {
		return err
	}

	proxyListenersRaw, ok := wrapper["proxy_listeners"]
	if ok {
		listeners, ok := proxyListenersRaw.([]interface{})
		if ok {
			d.ProxyListeners = make([]ProxyListener, 0, len(listeners))
			for _, listener := range listeners {
				proxyListener := ProxyListener{}
				if err := mapstructure.Decode(listener, &proxyListener); err != nil {
					return err
				}
				d.ProxyListeners = append(d.ProxyListeners, proxyListener)
			}
		}
	}

	streamListenersRaw, ok := wrapper["stream_listeners"]
	if ok {
		listeners, ok := streamListenersRaw.([]interface{})
		if ok {
			d.StreamListeners = make([]StreamListener, 0, len(listeners))
			for _, listener := range listeners {
				streamListener := StreamListener{}
				if err := mapstructure.Decode(listener, &streamListener); err != nil {
					return err
				}
				d.StreamListeners = append(d.StreamListeners, streamListener)
			}
		}
	}

	return nil
}
