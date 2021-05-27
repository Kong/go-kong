package kong

import "fmt"

// Identifier returns the endpoint key name or ID.
func (s1 *Service) Identifier() string {
	if s1.Name != nil {
		return *s1.Name
	}
	return *s1.ID
}

// Identifier returns the endpoint key name or ID.
func (r1 *Route) Identifier() string {
	if r1.Name != nil {
		return *r1.Name
	}
	return *r1.ID
}

// Identifier returns the endpoint key name or ID.
func (u1 *Upstream) Identifier() string {
	if u1.Name != nil {
		return *u1.Name
	}
	return *u1.ID
}

// Identifier returns the endpoint key name or ID.
func (t1 *Target) Identifier() string {
	if t1.Target != nil {
		return *t1.Target
	}
	return *t1.ID
}

// Identifier returns the endpoint key name or ID.
func (c1 *Certificate) Identifier() string {
	if c1.ID != nil {
		return *c1.ID
	}
	return *c1.Cert
}

// Identifier returns the endpoint key name or ID.
func (s1 *SNI) Identifier() string {
	if s1.Name != nil {
		return *s1.Name
	}
	return *s1.ID
}

// Identifier returns the endpoint key name or ID.
func (p1 *Plugin) Identifier() string {
	if p1.Name != nil {
		return *p1.Name
	}
	return *p1.ID
}

// Identifier returns the endpoint key name or ID.
func (c1 *Consumer) Identifier() string {
	if c1.Username != nil {
		return *c1.Username
	}
	return *c1.ID
}

// Identifier returns the endpoint key name or ID.
func (c1 *CACertificate) Identifier() string {
	if c1.ID != nil {
		return *c1.ID
	}
	return *c1.Cert
}

// Identifier returns the endpoint key name or ID.
func (r1 *RBACRole) Identifier() string {
	if r1.Name != nil {
		return *r1.Name
	}
	return *r1.ID
}

// Identifier returns a composite ID base on Role ID, workspace, and endpoint
func (r1 *RBACEndpointPermission) Identifier() string {
	if r1.Endpoint != nil {
		return fmt.Sprintf("%s-%s-%s", *r1.Role.ID, *r1.Workspace, *r1.Endpoint)
	}
	return *r1.Endpoint
}
