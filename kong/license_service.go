package kong

import (
	"context"
	"encoding/json"
	"fmt"
)

// AbstractLicenseService handles Licenses in Kong.
type AbstractLicenseService interface {
	// Create creates a License in Kong.
	Create(ctx context.Context, license *License) (*License, error)
	// Get fetches a License in Kong.
	Get(ctx context.Context, ID *string) (*License, error)
	// Update updates a License in Kong
	Update(ctx context.Context, license *License) (*License, error)
	// Delete deletes a License in Kong
	Delete(ctx context.Context, ID *string) error
	// List fetches a list of Licenses in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*License, *ListOpt, error)
	// ListAll fetches all Licenses in Kong.
	ListAll(ctx context.Context) ([]*License, error)
}

// LicenseService handles Licenses in Kong.
type LicenseService service

// Create creates a License in Kong.
// If an ID is specified, it will be used to
// create a license in Kong, otherwise an ID
// is auto-generated.
func (s *LicenseService) Create(ctx context.Context,
	license *License,
) (*License, error) {
	if license == nil {
		return nil, fmt.Errorf("cannot create a nil license")
	}

	queryPath := "/licenses"
	method := "POST"
	if license.ID != nil {
		queryPath = queryPath + "/" + *license.ID
		method = "PUT"
	}
	req, err := s.client.NewRequest(method, queryPath, nil, license)
	if err != nil {
		return nil, err
	}

	var createdLicense License
	_, err = s.client.Do(ctx, req, &createdLicense)
	if err != nil {
		return nil, err
	}
	return &createdLicense, nil
}

// Get fetches a License in Kong.
func (s *LicenseService) Get(ctx context.Context,
	ID *string,
) (*License, error) {
	if isEmptyString(ID) {
		return nil, fmt.Errorf("ID cannot be nil for Get operation")
	}

	endpoint := fmt.Sprintf("/licenses/%v", *ID)
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var license License
	_, err = s.client.Do(ctx, req, &license)
	if err != nil {
		return nil, err
	}
	return &license, nil
}

// Update updates a License in Kong
func (s *LicenseService) Update(ctx context.Context,
	license *License,
) (*License, error) {
	if license == nil {
		return nil, fmt.Errorf("cannot update a nil license")
	}

	if isEmptyString(license.ID) {
		return nil, fmt.Errorf("ID cannot be nil for Update operation")
	}

	endpoint := fmt.Sprintf("/licenses/%v", *license.ID)
	req, err := s.client.NewRequest("PATCH", endpoint, nil, license)
	if err != nil {
		return nil, err
	}

	var updatedAPI License
	_, err = s.client.Do(ctx, req, &updatedAPI)
	if err != nil {
		return nil, err
	}
	return &updatedAPI, nil
}

// Delete deletes a License in Kong
func (s *LicenseService) Delete(ctx context.Context,
	ID *string,
) error {
	if isEmptyString(ID) {
		return fmt.Errorf("ID cannot be nil for Delete operation")
	}

	endpoint := fmt.Sprintf("/licenses/%v", *ID)
	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}

// List fetches a list of Licenses in Kong.
// opt can be used to control pagination.
func (s *LicenseService) List(ctx context.Context,
	opt *ListOpt,
) ([]*License, *ListOpt, error) {
	data, next, err := s.client.list(ctx, "/licenses", opt)
	if err != nil {
		return nil, nil, err
	}
	var licenses []*License

	for _, object := range data {
		b, err := object.MarshalJSON()
		if err != nil {
			return nil, nil, err
		}
		var license License
		err = json.Unmarshal(b, &license)
		if err != nil {
			return nil, nil, err
		}
		licenses = append(licenses, &license)
	}

	return licenses, next, nil
}

// ListAll fetches all Licenses in Kong.
// This method can take a while if there
// a lot of Licenses present.
func (s *LicenseService) ListAll(ctx context.Context) ([]*License, error) {
	var licenses, data []*License
	var err error
	opt := &ListOpt{Size: pageSize}

	for opt != nil {
		data, opt, err = s.List(ctx, opt)
		if err != nil {
			return nil, err
		}
		licenses = append(licenses, data...)
	}
	return licenses, nil
}
