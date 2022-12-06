package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractKeyService interface {
	// Create creates a Key in Kong.
	Create(ctx context.Context, key *Key) (*Key, error)
	// Get fetches a Key in Kong.
	Get(ctx context.Context, nameOrID *string) (*Key, error)
	// Update updates a Key in Kong
	Update(ctx context.Context, key *Key) (*Key, error)
	// Delete deletes a Key in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of Keys in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Key, *ListOpt, error)
	// ListAll fetches all Keys in Kong.
	ListAll(ctx context.Context) ([]*Key, error)
}

type KeyService service

// Create creates a Key in Kong.
// If an ID is specified, it will be used to
// create a key in Kong, otherwise an ID
// is auto-generated.
func (s *KeyService) Create(ctx context.Context,
	key *Key,
) (*Key, error) {
	queryPath := "/keys"
	method := "POST"
	if key.ID != nil {
		queryPath = queryPath + "/" + *key.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, key)
	if err != nil {
		return nil, err
	}

	var createdKey Key
	_, err = s.client.Do(ctx, req, &createdKey)
	if err != nil {
		return nil, err
	}
	return &createdKey, nil
}

// Get fetches a Consumer in Kong.
func (s *KeyService) Get(ctx context.Context,
	nameOrID *string,
) (*Key, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/keys/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var key Key
	_, err = s.client.Do(ctx, req, &key)
	if err != nil {
		return nil, err
	}
	return &key, nil
}

// Update updates a Consumer in Kong
func (s *KeyService) Update(ctx context.Context,
	key *Key,
) (*Key, error) {
	if isEmptyString(key.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/keys/%v", *key.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, key)
	if err != nil {
		return nil, err
	}

	var updatedAPI Key
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a Consumer in Kong
func (s *KeyService) Delete(ctx context.Context,
	nameOrID *string,
) error {
	if isEmptyString(nameOrID) {
		return fmt.Errorf("nameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/keys/%v", *nameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Keys in Kong.
// opt can be used to control pagination.
func (s *KeyService) List(ctx context.Context,
	opt *ListOpt,
) ([]*Key, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/keys", opt)
	if err != nil {
		return nil, nil, err
	}
	var keys []*Key

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var key Key
		err = json.Unmarshal(b, &key)
		if err != nil {
			return nil, nil, err
		}
		keys = append(keys, &key)
	}

	return keys, next, nil
}

// ListAll fetches all Keys in Kong.
// This method can take a while if there
// a lot of Keys present.
func (s *KeyService) ListAll(ctx context.Context) ([]*Key, error) {
	var keys, data []*Key
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		keys = append(keys, data...)
	}
	return keys, nil
}
