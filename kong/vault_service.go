package kong

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
)

// AbstractVaultService handles Vaults in Kong.
type AbstractVaultService interface {
	// Create creates a Vault in Kong
	Create(ctx context.Context, vault *Vault) (*Vault, error)
	// Get fetches a Vault in Kong.
	Get(ctx context.Context, nameOrID *string) (*Vault, error)
	// Update updates a Vault in Kong
	Update(ctx context.Context, vault *Vault) (*Vault, error)
	// Delete deletes a Vault in Kong
	Delete(ctx context.Context, nameOrID *string) error
	// List fetches a list of Vaults in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*Vault, *ListOpt, error)
	// ListAll fetches all Vaults in Kong.
	ListAll(ctx context.Context) ([]*Vault, error)
	// Validate validates a Vault against its schema.
	Validate(ctx context.Context, vault *Vault) (bool, string, error)
}

// VaultService handles Vaults in Kong.
type VaultService service

// Create creates a Vault in Kong
// If an ID is specified, it will be used to
// create a Vault in Kong, otherwise an ID
// is auto-generated.
func (s *VaultService) Create(ctx context.Context, vault *Vault) (*Vault, error) {
	if vault == nil {
		return nil, fmt.Errorf("cannot create a nil vault")
	}

	endpoint := "/vaults"
	method := "POST"
	if vault.ID != nil {
		endpoint = endpoint + "/" + *vault.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, endpoint, nil, vault)
	if err != nil {
		return nil, err
	}

	var createdVault Vault
	_, err = s.client.Do(ctx, req, &createdVault)
	if err != nil {
		return nil, err
	}
	return &createdVault, nil
}

// Get fetches a Vault in Kong.
func (s *VaultService) Get(ctx context.Context, prefixOrID *string) (*Vault, error) {
	if isEmptyString(prefixOrID) {
		return nil, fmt.Errorf("prefixOrID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/vaults/%v", *prefixOrID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var vault Vault
	_, err = s.client.Do(ctx, req, &vault)
	if err != nil {
		return nil, err
	}
	return &vault, nil
}

// Update updates a Vault in Kong
func (s *VaultService) Update(ctx context.Context, vault *Vault) (*Vault, error) {
	if vault == nil {
		return nil, fmt.Errorf("cannot update a nil vault")
	}

	if isEmptyString(vault.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/vaults/%v", *vault.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, vault)
	if err != nil {
		return nil, err
	}

	var updatedVault Vault
	_, err = s.client.Do(ctx, req, &updatedVault)
	if err != nil {
		return nil, err
	}
	return &updatedVault, nil
}

// Delete deletes a Vault in Kong
func (s *VaultService) Delete(ctx context.Context, prefixOrID *string) error {
	if isEmptyString(prefixOrID) {
		return fmt.Errorf("prefixOrID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/vaults/%v", *prefixOrID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Vaults in Kong.
// opt can be used to control pagination.
func (s *VaultService) List(ctx context.Context,
	opt *ListOpt,
) ([]*Vault, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/vaults", opt)
	if err != nil {
		return nil, nil, err
	}
	vaults := make([]*Vault, 0, len(data))
	for _, object := range data {
		var vault Vault
		err = json.Unmarshal(object, &vault)
		if err != nil {
			return nil, nil, err
		}
		vaults = append(vaults, &vault)
	}

	return vaults, next, nil
}

// ListAll fetches all Vaults in Kong.
// This method can take a while if there
// a lot of Vaults present.
func (s *VaultService) ListAll(ctx context.Context) ([]*Vault, error) {
	var vaults, data []*Vault
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		vaults = append(vaults, data...)
	}
	return vaults, nil
}

// Validate validates a vault against its schema.
// returns validate result (passed/failed) and the message from the schema validation service if validation fails.
// returns a non-nil error if failed to call the schema validation service.
func (s *VaultService) Validate(ctx context.Context, vault *Vault) (bool, string, error) {
	const endpoint = "/schemas/vaults/validate"
	req, err := s.client.NewRequest(http.MethodPost, endpoint, nil, vault)
	if err != nil {
		return false, "", err
	}
	if _, err := s.client.Do(ctx, req, nil); err != nil {

		var apiErr *APIError
		if ok := errors.As(err, &apiErr); !ok || apiErr.Code() != http.StatusBadRequest {
			return false, "", err
		}
		return false, apiErr.message, nil
	}
	return true, "", nil
}
