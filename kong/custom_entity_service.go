package kong

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/kong/go-kong/kong/custom"
)

// AbstractCustomEntityService handles custom entities in Kong.
type AbstractCustomEntityService interface {
	// Get fetches a custom entity. The primary key and all relations of the
	// entity must be populated in entity.
	Get(ctx context.Context, entity custom.Entity) (custom.Entity, error)
	// Create creates a custom entity based on entity.
	// All required fields must be present in entity.
	Create(ctx context.Context, entity custom.Entity) (custom.Entity, error)
	// Update updates a custom entity in Kong.
	Update(ctx context.Context, entity custom.Entity) (custom.Entity, error)
	// Delete deletes a custom entity in Kong.
	Delete(ctx context.Context, entity custom.Entity) error
	// List fetches all custom entities based on relations
	List(ctx context.Context, opt *ListOpt, entity custom.Entity) ([]custom.Entity, *ListOpt, error)
	// ListAll fetches all custom entities based on relations
	ListAll(ctx context.Context, entity custom.Entity) ([]custom.Entity, error)
}

// CustomEntityService handles custom entities in Kong.
type CustomEntityService service

// Get fetches a custom entity. The primary key and all relations of the
// entity must be populated in entity.
func (s *CustomEntityService) Get(ctx context.Context,
	entity custom.Entity,
) (custom.Entity, error) {
	def := s.client.Lookup(entity.Type())
	if def == nil {
		return nil, fmt.Errorf("entity '" + string(entity.Type()) +
			"' not registered")
	}

	queryPath, err := def.GetEndpoint(entity)
	if err != nil {
		return nil, err
	}

	req, err := s.client.NewRequest("GET", queryPath, nil, nil)
	if err != nil {
		return nil, err
	}

	var object custom.Object
	_, err = s.client.Do(ctx, req, &object)
	if err != nil {
		return nil, err
	}
	entity.SetObject(object)
	return entity, nil
}

// Create creates a custom entity based on entity.
// All required fields must be present in entity.
func (s *CustomEntityService) Create(ctx context.Context,
	entity custom.Entity,
) (custom.Entity, error) {
	def := s.client.Lookup(entity.Type())
	if def == nil {
		return nil, fmt.Errorf("entity '" + string(entity.Type()) +
			"' not registered")
	}

	method := "POST"
	queryPath, err := def.PostEndpoint(entity)
	if err != nil {
		return nil, err
	}
	if entity.Object() != nil {
		if _, ok := entity.Object()["id"]; ok {
			queryPath, err = def.PatchEndpoint(entity)
			if err != nil {
				return nil, err
			}
			method = "PUT"
		}
	}

	o := entity.Object()
	// Necessary to Marshal an empty map
	// as {} and not null
	if len(o) == 0 {
		o = make(map[string]interface{})
	}
	req, err := s.client.NewRequest(method, queryPath, nil, o)
	if err != nil {
		return nil, err
	}

	var object custom.Object
	_, err = s.client.Do(ctx, req, &object)
	if err != nil {
		return nil, err
	}
	entity.SetObject(object)
	return entity, nil
}

// Update updates a custom entity in Kong.
func (s *CustomEntityService) Update(ctx context.Context,
	entity custom.Entity,
) (custom.Entity, error) {
	def := s.client.Lookup(entity.Type())
	if def == nil {
		return nil, fmt.Errorf("entity '" + string(entity.Type()) +
			"' not registered")
	}

	queryPath, err := def.PatchEndpoint(entity)
	if err != nil {
		return nil, err
	}

	o := entity.Object()
	// Necessary to Marshal an empty map
	// as {} and not null
	if len(o) == 0 {
		o = make(map[string]interface{})
	}
	req, err := s.client.NewRequest("PATCH", queryPath, nil, o)
	if err != nil {
		return nil, err
	}

	var object custom.Object
	_, err = s.client.Do(ctx, req, &object)
	if err != nil {
		return nil, err
	}
	entity.SetObject(object)
	return entity, nil
}

// Delete deletes a custom entity in Kong.
func (s *CustomEntityService) Delete(ctx context.Context,
	entity custom.Entity,
) error {
	def := s.client.Lookup(entity.Type())
	if def == nil {
		return fmt.Errorf("entity '" + string(entity.Type()) +
			"' not registered")
	}

	queryPath, err := def.PatchEndpoint(entity)
	if err != nil {
		return err
	}

	req, err := s.client.NewRequest("DELETE", queryPath, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches all custom entities based on relations
func (s *CustomEntityService) List(ctx context.Context, opt *ListOpt,
	entity custom.Entity,
) ([]custom.Entity, *ListOpt, error) {
	def := s.client.Lookup(entity.Type())
	if def == nil {
		return nil, nil, fmt.Errorf("entity '" + string(entity.Type()) +
			"' not registered")
	}

	queryPath, err := def.ListEndpoint(entity)
	if err != nil {
		return nil, nil, err
	}

	data, next, err := s.client.list(ctx, queryPath, opt)
	if err != nil {
		return nil, nil, err
	}
	var entities []custom.Entity

	for _, o := range data {
		b, err := o.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var object custom.Object
		err = json.Unmarshal(b, &object)
		if err != nil {
			return nil, nil, err
		}
		e := custom.NewEntityObject(entity.Type())
		e.SetObject(object)
		for k, v := range entity.GetAllRelations() {
			e.AddRelation(k, v)
		}
		entities = append(entities, e)
	}

	return entities, next, nil
}

// ListAll fetches all custom entities based on relations
func (s *CustomEntityService) ListAll(ctx context.Context,
	entity custom.Entity,
) ([]custom.Entity, error) {
	var entities, data []custom.Entity
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt, entity)
		if err != nil {
			return nil, err
		}
		entities = append(entities, data...)
	}
	return entities, nil
}
