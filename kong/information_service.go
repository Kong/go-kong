package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
)

// AbstractInfoService handles Kong's Information.
type AbstractInfoService interface {
	// Get retrieves the general runtime information about the Kong gateway.
	Get(ctx context.Context) (*Info, error)
	// IsConfigReady returns true if the Kong instance has configuration and false if it does not
	IsConfigReady(ctx context.Context) (bool, error)
}

type InfoService service

// IsInMemory check if Kong is in memory
func (r *RuntimeConfiguration) IsInMemory() bool {
	return r.Database == "off"
}

// IsRBACEnabled check if RBAC are enabled
func (r *RuntimeConfiguration) IsRBACEnabled() bool {
	return r.RBAC == "on"
}

// convert convert an object to another through json marshalling
// unmarshalling
func convert(from, to interface{}) error {
	bytes, err := json.Marshal(from)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, to)
}

// Get retrieves the high-level metadata of a Kong instance.
func (s *InfoService) Get(ctx context.Context) (*Info, error) {
	information, err := s.client.Root(ctx)
	if err != nil {
		return nil, err
	}
	var info Info
	if err := convert(information, &info); err != nil {
		return nil, err
	}
	return &info, nil
}

// IsConfigReady returns true if the Kong instance has configuration and false if it does not
func (s *InfoService) IsConfigReady(ctx context.Context) (bool, error) {
	req, err := s.client.NewRequest("GET", "/config/ready/", nil, nil)
	if err != nil {
		return false, err
	}

	resp, err := s.client.Do(ctx, req, nil)
	if resp.StatusCode != http.StatusOK && resp.StatusCode != http.StatusServiceUnavailable {
		if err != nil {
			return false, fmt.Errorf("could not determine config state: %v %v", resp.StatusCode, err)
		}
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			body = []byte(err.Error())
		}
		return false, fmt.Errorf("could not determine config state: %v %v", resp.StatusCode, body)
	}
	return resp.StatusCode == http.StatusOK, nil
}
