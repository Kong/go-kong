package kong

// Key represents a JWK or PEM key in Kong.
// +k8s:deepcopy-gen=true
type Key struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Set       *KeySet   `json:"set,omitempty" yaml:"set,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	KID       *string   `json:"kid,omitempty" yaml:"kid,omitempty"`
	JWK       *string   `json:"jwk,omitempty" yaml:"jwk,omitempty"`
	PEM       *PEM      `json:"pem,omitempty" yaml:"pem,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// KeySet represents a set of keys in Kong.
// +k8s:deepcopy-gen=true
type KeySet struct {
	ID        *string   `json:"id,omitempty" yaml:"id,omitempty"`
	CreatedAt *int64    `json:"created_at,omitempty" yaml:"created_at,omitempty"`
	UpdatedAt *int64    `json:"updated_at,omitempty" yaml:"updated_at,omitempty"`
	Name      *string   `json:"name,omitempty" yaml:"name,omitempty"`
	Tags      []*string `json:"tags,omitempty" yaml:"tags,omitempty"`
}

// PEM represents a PEM formatted key in Kong.
// +k8s:deepcopy-gen=true
type PEM struct {
	PublicKey  *string `json:"public_key,omitempty" yaml:"public_key,omitempty"`
	PrivateKey *string `json:"private_key,omitempty" yaml:"private_key,omitempty"`
}
