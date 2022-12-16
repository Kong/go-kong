package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

type AbstractKeySetService interface {
	// Create creates a Key in Kong.
	Create(ctx context.Context, keySet *KeySet) (*KeySet, error)
	// Get fetches a Key in Kong.
	Get(ctx context.Context, nameOrID *string) (*KeySet, error)
	// Update updates a Key in Kong
	Update(ctx context.Context, keySet *KeySet) (*KeySet, error)
	// Delete deletes a Key in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of Keys in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*KeySet, *ListOpt, error)
	// ListAll fetches all Keys in Kong.
	ListAll(ctx context.Context) ([]*KeySet, error)
}

type KeySetService service

// Create creates a Key in Kong.
// If an ID is specified, it will be used to
// create a key in Kong, otherwise an ID
// is auto-generated.
func (s *KeySetService) Create(ctx context.Context,
	keySet *KeySet,
) (*KeySet, error) {
	queryPath := "/key-sets"
	method := "POST"
	if keySet.ID != nil {
		queryPath = queryPath + "/" + *keySet.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, keySet)
	if err != nil {
		return nil, err
	}

	var createdKeySet KeySet
	_, err = s.client.Do(ctx, req, &createdKeySet)
	if err != nil {
		return nil, err
	}
	return &createdKeySet, nil
}

// Get fetches a KeySet in Kong.
func (s *KeySetService) Get(ctx context.Context,
	nameOrID *string,
) (*KeySet, error) {
	if isEmptyString(nameOrID) {
		return nil, fmt.Errorf("nameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/key-sets/%v", *nameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var keySet KeySet
	_, err = s.client.Do(ctx, req, &keySet)
	if err != nil {
		return nil, err
	}
	return &keySet, nil
}

// Update updates a KeySet in Kong
func (s *KeySetService) Update(ctx context.Context,
	keySet *KeySet,
) (*KeySet, error) {
	if isEmptyString(keySet.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/key-sets/%v", *keySet.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, keySet)
	if err != nil {
		return nil, err
	}

	var updatedAPI KeySet
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a KeySet in Kong
func (s *KeySetService) Delete(ctx context.Context,
	nameOrID *string,
) error {
	if isEmptyString(nameOrID) {
		return fmt.Errorf("nameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/key-sets/%v", *nameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of KeySets in Kong.
// opt can be used to control pagination.
func (s *KeySetService) List(ctx context.Context,
	opt *ListOpt,
) ([]*KeySet, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/key-sets", opt)
	if err != nil {
		return nil, nil, err
	}
	var keysets []*KeySet

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var keyset KeySet
		err = json.Unmarshal(b, &keyset)
		if err != nil {
			return nil, nil, err
		}
		keysets = append(keysets, &keyset)
	}

	return keysets, next, nil
}

// ListAll fetches all KeySets in Kong.
// This method can take a while if there
// a lot of KeySets present.
func (s *KeySetService) ListAll(ctx context.Context) ([]*KeySet, error) {
	var keysets, data []*KeySet
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		keysets = append(keysets, data...)
	}
	return keysets, nil
}
