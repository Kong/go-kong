package kong

import (
	"context"
	"encoding/json"
)

type AbstractKeyService interface {
	// ListAll fetches all Keys in Kong.
	ListAll(ctx context.Context) ([]*Key, error)
}

type KeyService service

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
