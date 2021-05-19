package kong

import (
	"context"
	"encoding/json"
	"strings"
)

type Info struct {
	Version       string `json:"version,omitempty" yaml:"version,omitempty"`
	Configuration struct {
		Database string `json:"database,omitempty" yaml:"database,omitempty"`
		Portal   bool   `json:"portal,omitempty" yaml:"portal,omitempty"`
		RBAC     string `json:"rbac,omitempty" yaml:"rbac,omitempty"`
	} `json:"configuration,omitempty" yaml:"configuration,omitempty"`
}

// IsEnterprise check if Kong is enterprise edition
func (s *Info) IsEnterprise(ctx context.Context) bool {
	return strings.Contains(s.Version, "enterprise")

}

// IsInMemory check if Kong is in memory
func (s *Info) IsInMemory(ctx context.Context) bool {
	return "off" == s.Configuration.Database
}

// IsRBACEnabled check if RBAC are enabled
func (s *Info) IsRBACEnabled(ctx context.Context) bool {
	return "on" == s.Configuration.RBAC
}

// IsPortalEnabled check if the portal is enabled
func (s *Info) IsPortalEnabled(ctx context.Context) bool {
	return s.Configuration.Portal
}

func (c *Client) Info(ctx context.Context) (*Info, error) {
	information, err := c.Root(ctx)
	if err != nil {
		return nil, err
	}

	bytes, err := json.Marshal(information)
	if err != nil {
		return nil, err
	}
	var info Info
	err = json.Unmarshal(bytes, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
