//nolint:lll
package kong

import (
	"context"
	"fmt"
	"io"
)

// AbstractConfigService handles Config in Kong.
type AbstractConfigService interface {
	// ReloadDeclarativeRawConfig sends out the specified config to configured Admin
	// API endpoint using the provided reader which should contain the JSON
	// serialized body that adheres to the configuration format specified at:
	// https://docs.konghq.com/gateway/latest/production/deployment-topologies/db-less-and-declarative-config/#declarative-configuration-format
	ReloadDeclarativeRawConfig(ctx context.Context, config io.Reader, checkHash bool) error
}

// ConfigService handles Config in Kong.
type ConfigService service

// ReloadDeclarativeRawConfig sends out the specified config to configured Admin
// API endpoint using the provided reader which should contain the JSON
// serialized body that adheres to the configuration format specified at:
// https://docs.konghq.com/gateway/latest/production/deployment-topologies/db-less-and-declarative-config/#declarative-configuration-format
func (c *ConfigService) ReloadDeclarativeRawConfig(
	ctx context.Context,
	config io.Reader,
	checkHash bool,
) error {
	type sendConfigParams struct {
		CheckHash int `url:"check_hash"`
	}
	var checkHashI int
	if checkHash {
		checkHashI = 1
	}
	req, err := c.client.NewRequest("POST", "/config", sendConfigParams{CheckHash: checkHashI}, config)
	if err != nil {
		return fmt.Errorf("creating new HTTP request for /config: %w", err)
	}

	resp, err := c.client.DoRAW(ctx, req)
	if err != nil {
		return fmt.Errorf("failed posting new config to /config: %w", err)
	}
	defer resp.Body.Close()

	if resp.StatusCode < 200 || resp.StatusCode >= 400 {
		b, err := io.ReadAll(resp.Body)
		if err != nil {
			return fmt.Errorf(
				"failed posting new config to /config: got status code %d (and failed to read the response body): %w",
				resp.StatusCode, err,
			)
		}

		return fmt.Errorf(
			"failed posting new config to /config: got status code %d, body: %s",
			resp.StatusCode, b,
		)
	}

	return nil
}
