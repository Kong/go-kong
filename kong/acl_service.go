package kong

import (
	"context"
	"encoding/json"
)

// AbstractACLService handles consumer ACL groups in Kong.
type AbstractACLService interface {
	// Create adds a consumer to an ACL group in Kong
	Create(ctx context.Context, consumerUsernameOrID *string, aclGroup *ACLGroup) (*ACLGroup, error)
	// Get fetches an ACL group for a consumer in Kong.
	Get(ctx context.Context, consumerUsernameOrID, groupOrID *string) (*ACLGroup, error)
	// Update updates an ACL group for a consumer in Kong
	Update(ctx context.Context, consumerUsernameOrID *string, aclGroup *ACLGroup) (*ACLGroup, error)
	// Delete deletes an ACL group association for a consumer in Kong
	Delete(ctx context.Context, consumerUsernameOrID, groupOrID *string) error
	// List fetches a list of ACL group and consumer associations in Kong.
	List(ctx context.Context, opt *ListOpt) ([]*ACLGroup, *ListOpt, error)
	// ListAll fetches all ACL group associations in Kong.
	ListAll(ctx context.Context) ([]*ACLGroup, error)
	// ListAllByOpt fetches all ACL group associations filtered by opt in Kong.
	ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*ACLGroup, error)
	// ListAllByTags fetches all ACL group associations filtered by tags in Kong.
	ListAllByTags(ctx context.Context, tags []string) ([]*ACLGroup, error)
	// ListForConsumer fetches a list of ACL groups
	// in Kong associated with a specific consumer.
	ListForConsumer(ctx context.Context, consumerUsernameOrID *string, opt *ListOpt) ([]*ACLGroup, *ListOpt, error)
}

// ACLService handles consumer ACL groups in Kong.
type ACLService service

// Create adds a consumer to an ACL group in Kong
// If an ID is specified, it will be used to
// create the group association in Kong, otherwise an ID
// is auto-generated.
func (s *ACLService) Create(ctx context.Context,
	consumerUsernameOrID *string, aclGroup *ACLGroup) (*ACLGroup, error) {

	cred, err := s.client.credentials.Create(ctx, "acl",
		consumerUsernameOrID, aclGroup)
	if err != nil {
		return nil, err
	}

	var createdACLGroup ACLGroup
	err = json.Unmarshal(cred, &createdACLGroup)
	if err != nil {
		return nil, err
	}

	return &createdACLGroup, nil
}

// Get fetches an ACL group for a consumer in Kong.
func (s *ACLService) Get(ctx context.Context,
	consumerUsernameOrID, groupOrID *string) (*ACLGroup, error) {

	cred, err := s.client.credentials.Get(ctx, "acl",
		consumerUsernameOrID, groupOrID)
	if err != nil {
		return nil, err
	}

	var aclGroup ACLGroup
	err = json.Unmarshal(cred, &aclGroup)
	if err != nil {
		return nil, err
	}

	return &aclGroup, nil
}

// Update updates an ACL group for a consumer in Kong
func (s *ACLService) Update(ctx context.Context,
	consumerUsernameOrID *string, aclGroup *ACLGroup) (*ACLGroup, error) {

	cred, err := s.client.credentials.Update(ctx, "acl",
		consumerUsernameOrID, aclGroup)
	if err != nil {
		return nil, err
	}

	var updatedACLGroup ACLGroup
	err = json.Unmarshal(cred, &updatedACLGroup)
	if err != nil {
		return nil, err
	}

	return &updatedACLGroup, nil
}

// Delete deletes an ACL group association for a consumer in Kong
func (s *ACLService) Delete(ctx context.Context,
	consumerUsernameOrID, groupOrID *string) error {
	return s.client.credentials.Delete(ctx, "acl",
		consumerUsernameOrID, groupOrID)
}

// List fetches a list of ACL group and consumer associations in Kong.
// opt can be used to control pagination.
func (s *ACLService) List(ctx context.Context,
	opt *ListOpt) ([]*ACLGroup, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/acls", opt)
}

// ListAll fetches all ACL group associations in Kong.
// This method can take a while if there
// a lot of ACLGroup associations are present.
func (s *ACLService) ListAll(ctx context.Context) ([]*ACLGroup, error) {
	return s.ListAllByTags(ctx, nil)
}

// ListAllByTags fetches all ACL group associations filtered by tags in Kong.
// This method can take a while if there
// a lot of ACLGroup associations are present.
func (s *ACLService) ListAllByTags(ctx context.Context, tags []string) ([]*ACLGroup, error) {
	return s.ListAllByOpt(ctx, newOpt(tags))
}

// ListAllByOpt fetches all ACL group associations filtered by opt in Kong.
// This method can take a while if there
// a lot of ACLGroup associations are present.
func (s *ACLService) ListAllByOpt(ctx context.Context, opt *ListOpt) ([]*ACLGroup, error) {
	return s.listAllByEndpointAndOpt(ctx, "/acls", opt)
}

// ListForConsumer fetches a list of ACL groups
// in Kong associated with a specific consumer.
// opt can be used to control pagination.
func (s *ACLService) ListForConsumer(ctx context.Context,
	consumerUsernameOrID *string, opt *ListOpt) ([]*ACLGroup, *ListOpt, error) {
	return s.listByEndpointAndOpt(ctx, "/consumers/"+*consumerUsernameOrID+"/acls", opt)
}

func (s *ACLService) listAllByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*ACLGroup, error) {
	data, err := s.client.listAll(ctx, endpoint, opt, true)
	if err != nil {
		return nil, err
	}
	var aclGroups []*ACLGroup
	for _, object := range data {
		var aclGroup ACLGroup
		err = json.Unmarshal(object, &aclGroup)
		if err != nil {
			return nil, err
		}
		aclGroups = append(aclGroups, &aclGroup)
	}

	return aclGroups, nil
}

func (s *ACLService) listByEndpointAndOpt(ctx context.Context,
	endpoint string, opt *ListOpt) ([]*ACLGroup, *ListOpt, error) {
	data, next, err := s.client.list(ctx, endpoint, opt)
	if err != nil {
		return nil, nil, err
	}
	var aclGroups []*ACLGroup
	for _, object := range data {
		var aclGroup ACLGroup
		err = json.Unmarshal(object, &aclGroup)
		if err != nil {
			return nil, nil, err
		}
		aclGroups = append(aclGroups, &aclGroup)
	}

	return aclGroups, next, nil
}
