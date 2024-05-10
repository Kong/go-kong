package kong

import (
	"context"
	"errors"
	"fmt"
	"net/http"
)

// EntityType describes a type of Kong entity.
type EntityType string

// Defition of constants for standard types of Kong entities.
const (
	EntityTypeServices       EntityType = "services"
	EntityTypeRoutes         EntityType = "routes"
	EntityTypeUpstreams      EntityType = "upstreams"
	EntityTypeTargets        EntityType = "targets"
	EntityTypePlugins        EntityType = "plugins"
	EntityTypeConsumers      EntityType = "consumers"
	EntityTypeConsumerGroups EntityType = "consumer_groups"
	EntityTypeSNIs           EntityType = "snis"
	EntityTypeCertificates   EntityType = "certificates"
	EntityTypeCACertificates EntityType = "ca_certificates"
	EntityTypeTags           EntityType = "tags"
)

// AbstractSchemaService handles schemas in Kong.
type AbstractSchemaService interface {
	// Get fetches an entity schema from Kong.
	Get(ctx context.Context, entity string) (Schema, error)
	// Validate validates an arbitrary entity in Kong.
	// Used for custom entities, or entities that do not have a Validate method in the corresponding service.
	Validate(ctx context.Context, entityType EntityType, entity interface{}) (bool, string, error)
}

// SchemaService handles schemas in Kong.
type SchemaService service

// Schema represents an entity schema in Kong.
type Schema map[string]interface{}

var _ AbstractSchemaService = &SchemaService{}

// Get retrieves the full schema of kong entities.
func (s *SchemaService) Get(ctx context.Context, entity string) (Schema, error) {
	req, err := s.client.NewRequest("GET", fmt.Sprintf("/schemas/%s", entity), nil, nil)
	if err != nil {
		return nil, err
	}
	var schema Schema
	_, err = s.client.Do(ctx, req, &schema)
	if err != nil {
		return nil, err
	}
	return schema, nil
}

// Validate validates an arbitrary entity in Kong.
func (s *SchemaService) Validate(ctx context.Context, entityType EntityType, entity interface{}) (bool, string, error) {
	endpoint := fmt.Sprintf("/schemas/%s/validate", entityType)
	req, err := s.client.NewRequest("POST", endpoint, nil, entity)
	if err != nil {
		return false, "", err
	}
	resp, err := s.client.Do(ctx, req, nil)
	if err != nil {
		if resp == nil {
			return false, "", err
		}
		var apiErr *APIError
		if ok := errors.As(err, &apiErr); !ok || apiErr.Code() != http.StatusBadRequest {
			return false, "", err
		}
		return false, apiErr.message, nil
	}
	return true, "", nil
}
