package kong

import (
	"context"
	"fmt"
)

// AbstractSchemaService handles schemas in Kong.
type AbstractSchemaService interface {
	// Get fetches an entity schema from Kong.
	Get(ctx context.Context, entity string) (Schema, error)
}

// SchemaService handles schemas in Kong.
type SchemaService service

// Schema represents an entity schema in Kong.
type Schema map[string]interface{}

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
