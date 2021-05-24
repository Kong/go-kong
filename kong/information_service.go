package kong

import (
	"context"
	"encoding/json"
)

// AbstractInfoService handles Kong's Information.
type AbstractInfoService interface {
	// Get retrieves the information concerning Kong
	Get(ctx context.Context) (*Info, error)
}

type InfoService service

// IsInMemory check if Kong is in memory
func (r *RuntimeConfiguration) IsInMemory() bool {
	return "off" == r.Database
}

// IsRBACEnabled check if RBAC are enabled
func (r *RuntimeConfiguration) IsRBACEnabled() bool {
	return "on" == r.RBAC
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

// Get retrieves the information concerning Kong
func (s *InfoService) Get(ctx context.Context) (*Info, error) {
	information, err := s.client.Root(ctx)
	if err != nil {
		return nil, err
	}
	var info Info
	err = convert(information, &info)
	if err != nil {
		return nil, err
	}
	return &info, nil
}
