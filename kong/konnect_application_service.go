package kong

import (
	"context"
	"encoding/json"
)

var _ AbstractKonnectApplicationService = &KonnectApplicationService{}

// AbstractACLService handles consumer ACL groups in Kong.
type AbstractKonnectApplicationService interface {
	// Create creates a Konnect Application in Kong.
	Create(ctx context.Context, key *KonnectApplication) (*KonnectApplication, error)
	// List fetches list of Konnect Applications in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*KonnectApplication, *ListOpt, error)
	// ListAll fetches all Konnect Applications in Kong.
	ListAll(ctx context.Context) ([]*KonnectApplication, error)
}

type KonnectApplicationService service

func (k *KonnectApplicationService) Create(ctx context.Context, key *KonnectApplication) (*KonnectApplication, error) {
	queryPath := "/konnect_applications"
	method := "POST"
	if key.ID != nil {
		queryPath = queryPath + "/" + *key.ID
		method = "PUT"
	}
	req, err := k.client.NewRequest(method, queryPath, nil, key)
	if err != nil {
		return nil, err
	}

	var createdKey KonnectApplication
	_, err = k.client.Do(ctx, req, &createdKey)
	if err != nil {
		return nil, err
	}
	return &createdKey, nil
}

func (k *KonnectApplicationService) List(ctx context.Context, opt *ListOpt) ([]*KonnectApplication, *ListOpt, error) {
	data, next, err := k.client.list(ctx, "/konnect_applications", opt)
	if err != nil {
		return nil, nil, err
	}
	var kaas []*KonnectApplication

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var kaa KonnectApplication
		err = json.Unmarshal(b, &kaa)
		if err != nil {
			return nil, nil, err
		}
		kaas = append(kaas, &kaa)
	}

	return kaas, next, nil
}

func (k *KonnectApplicationService) ListAll(ctx context.Context) ([]*KonnectApplication, error) {
	var kaa, data []*KonnectApplication
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = k.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		kaa = append(kaa, data...)
	}
	return kaa, nil
}
