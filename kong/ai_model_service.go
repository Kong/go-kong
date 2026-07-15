package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractAIModelService handles AI Models in Kong.
type AbstractAIModelService interface {
	// Create creates an AI Model in Kong.
	Create(ctx context.Context, aiModel *AIModel) (*AIModel, error)
	// Get fetches an AI Model in Kong.
	Get(ctx context.Context, usernameOrID *string) (*AIModel, error)
	// Update updates an AI Model in Kong
	Update(ctx context.Context, aiModel *AIModel) (*AIModel, error)
	// Delete deletes an AI Model in Kong
	Delete(ctx context.Context, usernameOrID *string) error
	// List fetches a list of AI Models in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*AIModel, *ListOpt, error)
	// ListAll fetches all AI Models in Kong.
	ListAll(ctx context.Context) ([]*AIModel, error)
}

// AIModelService handles AI Models in Kong.
type AIModelService service

// Create creates an AI Model in Kong.
// If an ID is specified, it will be used to
// create an AI Model in Kong, otherwise an ID
// is auto-generated.
func (s *AIModelService) Create(ctx context.Context, aiModel *AIModel) (*AIModel, error) {
	queryPath := "/ai-models"
	method := "POST"
	if aiModel.ID != nil {
		queryPath = queryPath + "/" + *aiModel.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, aiModel)
	if err != nil {
		return nil, err
	}

	var createdAIModel AIModel
	_, err = s.client.Do(ctx, req, &createdAIModel)
	if err != nil {
		return nil, err
	}
	return &createdAIModel, nil
}

// Get fetches an AI Model in Kong.
func (s *AIModelService) Get(ctx context.Context,
	usernameOrID *string,
) (*AIModel, error) {
	if isEmptyString(usernameOrID) {
		return nil, fmt.Errorf(
			"usernameOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/ai-models/%v", *usernameOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var aiModel AIModel
	_, err = s.client.Do(ctx, req, &aiModel)
	if err != nil {
		return nil, err
	}
	return &aiModel, nil
}

// Update updates an AI Model in Kong
func (s *AIModelService) Update(ctx context.Context, aiModel *AIModel) (*AIModel, error) {
	if isEmptyString(aiModel.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	method := "PATCH"
	if s.client.isKonnect {
		method = "PUT"
	}

	endpoint := fmt.Sprintf("/ai-models/%v", *aiModel.ID)
	req, err := s.client.NewRequest(method, endpoint, nil, aiModel)
	if err != nil {
		return nil, err
	}

	var updatedAIModel AIModel
	_, err = s.client.Do(ctx, req, &updatedAIModel)
	if err != nil {
		return nil, err
	}
	return &updatedAIModel, nil
}

// Delete deletes an AI Model in Kong
func (s *AIModelService) Delete(ctx context.Context, usernameOrID *string) error {
	if isEmptyString(usernameOrID) {
		return fmt.Errorf("usernameOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/ai-models/%v", *usernameOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of AI Models in Kong.
// opt can be used to control pagination.
func (s *AIModelService) List(ctx context.Context,
	opt *ListOpt,
) ([]*AIModel, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/ai-models", opt)
	if err != nil {
		return nil, nil, err
	}
	var aiModels []*AIModel
	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var aiModel AIModel
		err = json.Unmarshal(b, &aiModel)
		if err != nil {
			return nil, nil, err
		}
		aiModels = append(aiModels, &aiModel)
	}

	return aiModels, next, nil
}

// ListAll fetches all AI Models in Kong.
// This method can take a while if there
// are a lot of AI Models present.
func (s *AIModelService) ListAll(ctx context.Context) ([]*AIModel, error) {
	var aiModels, data []*AIModel
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		aiModels = append(aiModels, data...)
	}
	return aiModels, nil
}
