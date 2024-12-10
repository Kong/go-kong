package kong

// KonnectApplication represents Konnect-Application-Auth
// in Kong.
// Read https://docs.konghq.com/konnect/dev-portal/applications/application-overview/
// +k8s:deepcopy-gen=true
type KonnectApplication struct {
	ID                 *string             `json:"id"`
	CreatedAt          int64               `json:"created_at"`
	ClientID           string              `json:"client_id"`
	ConsumerGroups     []string            `json:"consumer_groups"`
	Scopes             []string            `json:"scopes"`
	AuthStrategyID     *string             `json:"auth_strategy_id"`
	ApplicationContext *ApplicationContext `json:"application_context"`
	ExhaustedScopes    []string            `json:"exhausted_scopes"`
	Tags               *[]string           `json:"tags"`
}

// ApplicationContext reprensents the application context inside the
// Konnenct-Application-Auth.
// +k8s:deepcopy-gen=true
type ApplicationContext struct {
	PortalID       *string `json:"portal_id"`
	ApplicationID  *string `json:"application_id"`
	DeveloperID    *string `json:"developer_id"`
	OrganizationID *string `json:"organization_id"`
}
