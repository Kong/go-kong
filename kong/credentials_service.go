package kong

import (
	"context"
	"encoding/json"
	"fmt"
	"reflect"
)

// abstractCredentialService handles credentials in Kong.
type abstractCredentialService interface {
	// Create creates a credential in Kong of type credType.
	Create(ctx context.Context, credType string, consumerUsernameOrID *string,
		credential interface{}, skipHash bool) (json.RawMessage, error)
	// Get fetches a credential of credType with credIdentifier from Kong.
	Get(ctx context.Context, credType string, consumerUsernameOrID *string,
		credIdentifier *string) (json.RawMessage, error)
	// GetByID fetches a credential of credType with ID from Kong.
	GetByID(ctx context.Context, credType string, credIdentifier *string) (json.RawMessage, error)
	// Update updates credential in Kong
	Update(ctx context.Context, credType string, consumerUsernameOrID *string,
		credential interface{}, skipHash bool) (json.RawMessage, error)
	// Delete deletes a credential in Kong
	Delete(ctx context.Context, credType string, consumerUsernameOrID, credIdentifier *string) error
}

// credentialService handles credentials in Kong.
type credentialService service

var credPath = map[string]string{
	"key-auth":   "key-auth",
	"basic-auth": "basic-auth",
	"hmac-auth":  "hmac-auth",
	"jwt-auth":   "jwt",
	"acl":        "acls",
	"oauth2":     "oauth2",
	"mtls-auth":  "mtls-auth",
}

var credPathAsParentResource = map[string]string{
	"key-auth":   "key-auths",
	"basic-auth": "basic-auths",
	"hmac-auth":  "hmac-auths",
	"jwt-auth":   "jwts",
	"acl":        "acls",
	"oauth2":     "oauth2",
	"mtls-auth":  "mtls-auths",
}

type queryStruct struct {
	// Skip hashing of password when creating/updating basic-auth credentials.
	SkipHash string `url:"skip_hash,omitempty"`
}

// Create creates a credential in Kong of type credType.
// If an ID is specified in the credential, it will be used to
// create a credential in Kong, otherwise an ID
// is auto-generated.
func (s *credentialService) Create(ctx context.Context, credType string,
	consumerUsernameOrID *string,
	credential interface{},
	skipHash bool,
) (json.RawMessage, error) {
	if isEmptyString(consumerUsernameOrID) {
		return nil, fmt.Errorf("consumerUsernameOrID cannot be nil")
	}

	subPath, ok := credPath[credType]
	if !ok {
		return nil, fmt.Errorf("unknown credential type: %v", credType)
	}
	endpoint := "/consumers/" + *consumerUsernameOrID + "/" + subPath
	method := "POST"
	if credential != nil {
		if id, ok := credential.(id); ok {
			if !reflect.ValueOf(id).IsNil() {
				uuid := id.id()
				if !isEmptyString(uuid) {
					endpoint = endpoint + "/" + *uuid
					method = "PUT"
				}
			}
		}
	}

	var q *queryStruct

	if credType == "basic-auth" && skipHash {
		q = &queryStruct{
			SkipHash: "true",
		}
	}

	req, err := s.client.NewRequest(method, endpoint, q, credential)
	if err != nil {
		return nil, err
	}

	var createdCredential json.RawMessage
	_, err = s.client.Do(ctx, req, &createdCredential)
	if err != nil {
		return nil, err
	}
	return createdCredential, nil
}

// Get fetches a credential of credType with credIdentifier and consumerIdentifier from Kong.
func (s *credentialService) Get(ctx context.Context, credType string,
	consumerUsernameOrID *string,
	credIdentifier *string,
) (json.RawMessage, error) {
	if isEmptyString(credIdentifier) {
		return nil, fmt.Errorf("credIdentifier cannot be nil for Get operation")
	}
	if isEmptyString(consumerUsernameOrID) {
		return nil, fmt.Errorf("consumerUsernameOrID cannot be nil")
	}

	subPath, ok := credPath[credType]
	if !ok {
		return nil, fmt.Errorf("unknown credential type: %v", credType)
	}
	endpoint := "/consumers/" + *consumerUsernameOrID + "/" +
		subPath + "/" + *credIdentifier
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var cred json.RawMessage
	_, err = s.client.Do(ctx, req, &cred)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// GetByID fetches a credential of credType with ID from Kong.
func (s *credentialService) GetByID(ctx context.Context, credType string, credIdentifier *string,
) (json.RawMessage, error) {
	if isEmptyString(credIdentifier) {
		return nil, fmt.Errorf("credIdentifier cannot be nil for GetByID operation")
	}

	subPath, ok := credPathAsParentResource[credType]
	if !ok {
		return nil, fmt.Errorf("unknown credential type: %v", credType)
	}
	endpoint := "/" +
		subPath + "/" + *credIdentifier
	req, err := s.client.NewRequest("GET", endpoint, nil, nil)
	if err != nil {
		return nil, err
	}

	var cred json.RawMessage
	_, err = s.client.Do(ctx, req, &cred)
	if err != nil {
		return nil, err
	}
	return cred, nil
}

// Update updates credential in Kong
func (s *credentialService) Update(ctx context.Context, credType string,
	consumerUsernameOrID *string,
	credential interface{},
	skipHash bool,
) (json.RawMessage, error) {
	if isEmptyString(consumerUsernameOrID) {
		return nil, fmt.Errorf("consumerUsernameOrID cannot be nil")
	}

	subPath, ok := credPath[credType]
	if !ok {
		return nil, fmt.Errorf("unknown credential type: %v", credType)
	}

	endpoint := "/consumers/" + *consumerUsernameOrID + "/" + subPath + "/"

	credID := ""
	if credential != nil {
		if id, ok := credential.(id); ok {
			if !reflect.ValueOf(id).IsNil() {
				uuid := id.id()
				if !isEmptyString(uuid) {
					credID = *uuid
				}
			}
		}
	}
	if credID == "" {
		return nil, fmt.Errorf("cannot update a credential without an ID")
	}

	endpoint = endpoint + credID

	var q *queryStruct

	if credType == "basic-auth" && skipHash {
		q = &queryStruct{
			SkipHash: "true",
		}
	}
	req, err := s.client.NewRequest("PATCH", endpoint, q, credential)
	if err != nil {
		return nil, err
	}

	var updatedCred json.RawMessage
	_, err = s.client.Do(ctx, req, &updatedCred)
	if err != nil {
		return nil, err
	}
	return updatedCred, nil
}

// Delete deletes a credential in Kong
func (s *credentialService) Delete(ctx context.Context, credType string,
	consumerUsernameOrID, credIdentifier *string,
) error {
	if isEmptyString(credIdentifier) {
		return fmt.Errorf("credIdentifier cannot be nil for Delete operation")
	}

	subPath, ok := credPath[credType]
	if !ok {
		return fmt.Errorf("unknown credential type: %v", credType)
	}
	endpoint := "/consumers/" + *consumerUsernameOrID + "/" + subPath + "/" +
		*credIdentifier

	req, err := s.client.NewRequest("DELETE", endpoint, nil, nil)
	if err != nil {
		return err
	}

	_, err = s.client.Do(ctx, req, nil)
	return err
}
