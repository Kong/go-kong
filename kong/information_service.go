package kong

import (
	"context"
	"errors"
	"strings"
)

// AbstractInformationService handles Plugins in Kong.
type AbstractInformationService interface {
	// Version returns the Kong version
	Version(ctx context.Context) (string, error)
	// IsEnterprise check if Kong is enterprise edition
	IsEnterprise(ctx context.Context) (bool, error)
	// Database retrieves the database used by Kong
	Database(ctx context.Context) (string, error)
	// IsInMemory check if Kong is in memory
	IsInMemory(ctx context.Context) (bool, error)
	// IsRBACEnabled check if RBAC are enabled
	IsRBACEnabled(ctx context.Context) (bool, error)
	// IsPortalEnabled check if the portal is enabled
	IsPortalEnabled(ctx context.Context) (bool, error)
}

// InformationService handles Kong's informations.
type InformationService service

// Version returns the Kong version
func (s *InformationService) Version(ctx context.Context) (string, error) {
	informations, err := s.client.Root(ctx)
	if err != nil {
		return "", err
	}
	return VersionFromInfo(informations), nil

}

// IsEnterprise check if Kong is enterprise edition
func (s *InformationService) IsEnterprise(ctx context.Context) (bool, error) {
	version, err := s.Version(ctx)
	if err != nil {
		return false, err
	}
	return strings.Contains(version, "enterprise"), nil

}

// Database retrieves the database used by Kong
func (s *InformationService) Database(ctx context.Context) (string, error) {
	informations, err := s.client.Root(ctx)
	if err != nil {
		return "", err
	}
	configuration := configurationFromInfo(informations)
	database, ok := configuration["database"]
	if !ok {
		return "", errors.New("database could not be parsed from kong configuration")
	}
	return database.(string), nil

}

// IsInMemory check if Kong is in memory
func (s *InformationService) IsInMemory(ctx context.Context) (bool, error) {
	database, err := s.Database(ctx)
	if err != nil {
		return false, err
	}
	return "off" == database || "" == database, nil
}

// IsRBACEnabled check if RBAC are enabled
func (s *InformationService) IsRBACEnabled(ctx context.Context) (bool, error) {
	informations, err := s.client.Root(ctx)
	if err != nil {
		return false, err
	}
	configuration := configurationFromInfo(informations)
	rbac, ok := configuration["rbac"]
	if !ok {
		return false, errors.New("rbac could not be parsed from kong configuration")
	}
	return "on" == rbac.(string), nil
}

// IsPortalEnabled check if the portal is enabled
func (s *InformationService) IsPortalEnabled(ctx context.Context) (bool, error) {
	informations, err := s.client.Root(ctx)
	if err != nil {
		return false, err
	}
	configuration := configurationFromInfo(informations)
	portal, ok := configuration["portal"]
	if !ok {
		return false, errors.New("portal could not be parsed from kong configuration")
	}
	return portal.(bool), nil
}

func configurationFromInfo(info map[string]interface{}) map[string]interface{} {
	configuration, ok := info["configuration"]
	if !ok {
		return make(map[string]interface{})
	}
	return configuration.(map[string]interface{})
}
