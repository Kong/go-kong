package kong

import (
	"context"
	"encoding/json"
)

type AbstractKeySetService interface {
	// ListAll fetches all Keys in Kong.
	ListAll(ctx context.Context) ([]*KeySet, error)
}

type KeySetService service

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
