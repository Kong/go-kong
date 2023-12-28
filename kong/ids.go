package kong

import (
	"fmt"
	"reflect"

	"github.com/google/uuid"
)

// FillID fills the ID of an entity. It is a no-op if the entity already has an ID.
// ID is generated in a deterministic way using UUIDv5. The UUIDv5 namespace is different for each entity type.
// The name used to generate the ID for Service is Service.Name.
func (s *Service) FillID() error {
	if s == nil {
		return fmt.Errorf("service is nil")
	}
	if s.ID != nil {
		// ID already set, do nothing.
		return nil
	}
	if s.Name == nil || *s.Name == "" {
		return fmt.Errorf("service name is required")
	}

	gen, err := idGeneratorFor(s)
	if err != nil {
		return fmt.Errorf("could not get id generator: %w", err)
	}

	s.ID = gen.buildIDFor(*s.Name)
	return nil
}

// FillID fills the ID of an entity. It is a no-op if the entity already has an ID.
// ID is generated in a deterministic way using UUIDv5. The UUIDv5 namespace is different for each entity type.
// The name used to generate the ID for Route is Route.Name.
func (r *Route) FillID() error {
	if r == nil {
		return fmt.Errorf("route is nil")
	}
	if r.ID != nil {
		// ID already set, do nothing.
		return nil
	}
	if r.Name == nil || *r.Name == "" {
		return fmt.Errorf("route name is required")
	}

	gen, err := idGeneratorFor(r)
	if err != nil {
		return fmt.Errorf("could not get id generator: %w", err)
	}

	r.ID = gen.buildIDFor(*r.Name)
	return nil
}

// FillID fills the ID of an entity. It is a no-op if the entity already has an ID.
// ID is generated in a deterministic way using UUIDv5. The UUIDv5 namespace is different for each entity type.
// The name used to generate the ID for Consumer is Consumer.Username.
func (c *Consumer) FillID() error {
	if c == nil {
		return fmt.Errorf("consumer is nil")
	}
	if c.ID != nil {
		// ID already set, do nothing.
		return nil
	}
	if c.Username == nil || *c.Username == "" {
		return fmt.Errorf("consumer username is required")
	}

	gen, err := idGeneratorFor(c)
	if err != nil {
		return fmt.Errorf("could not get id generator: %w", err)
	}

	c.ID = gen.buildIDFor(*c.Username)
	return nil
}

// FillID fills the ID of an entity. It is a no-op if the entity already has an ID.
// ID is generated in a deterministic way using UUIDv5. The UUIDv5 namespace is different for each entity type.
// The name used to generate the ID for ConsumerGroup is ConsumerGroup.Name.
func (cg *ConsumerGroup) FillID() error {
	if cg == nil {
		return fmt.Errorf("consumer group is nil")
	}
	if cg.ID != nil {
		// ID already set, do nothing.
		return nil
	}
	if cg.Name == nil || *cg.Name == "" {
		return fmt.Errorf("consumer group name is required")
	}

	gen, err := idGeneratorFor(cg)
	if err != nil {
		return fmt.Errorf("could not get id generator: %w", err)
	}

	cg.ID = gen.buildIDFor(*cg.Name)
	return nil
}

// FillID fills the ID of an entity. It is a no-op if the entity already has an ID.
// ID is generated in a deterministic way using UUIDv5. The UUIDv5 namespace is different for each entity type.
// The name used to generate the ID for Vault is Vault.Prefix.
func (v *Vault) FillID() error {
	if v == nil {
		return fmt.Errorf("vault is nil")
	}
	if v.ID != nil && *v.ID != "" {
		// ID already set, do nothing.
		return nil
	}
	if v.Prefix == nil || *v.Prefix == "" {
		return fmt.Errorf("vault prefix is required")
	}

	gen, err := idGeneratorFor(v)
	if err != nil {
		return fmt.Errorf("could not get id generator: %w", err)
	}

	v.ID = gen.buildIDFor(*v.Prefix)
	return nil
}

var (
	// _kongEntitiesNamespace is the UUIDv5 namespace used to generate IDs for Kong entities.
	_kongEntitiesNamespace = uuid.MustParse("fd02801f-0957-4a15-a55a-c8d9606f30b5")

	// _idGenerators is a map of entity type to ID generator.
	// Plural names of entities are used as names for UUIDv5 namespaces to match Kong's behavior which uses schemas
	// names for that purpose.
	// See https://github.com/Kong/kong/blob/master/kong/db/schema/others/declarative_config.lua for reference.
	_idGenerators = map[reflect.Type]idGenerator{
		reflect.TypeOf(Service{}):       newIDGeneratorFor("services"),
		reflect.TypeOf(Route{}):         newIDGeneratorFor("routes"),
		reflect.TypeOf(Consumer{}):      newIDGeneratorFor("consumers"),
		reflect.TypeOf(ConsumerGroup{}): newIDGeneratorFor("consumergroups"),
		reflect.TypeOf(Vault{}):         newIDGeneratorFor("vaults"),
	}
)

type idGenerator struct {
	namespace uuid.UUID
}

func (g idGenerator) buildIDFor(entityKey string) *string {
	id := uuid.NewSHA1(g.namespace, []byte(entityKey)).String()
	return &id
}

// newIDGeneratorFor returns a new ID generator for the given entity type. Should be used only to initialize _idGenerators.
func newIDGeneratorFor(entityPluralName string) idGenerator {
	entityTypeNamespace := uuid.NewSHA1(_kongEntitiesNamespace, []byte(entityPluralName))
	return idGenerator{namespace: entityTypeNamespace}
}

// IDFillable is a type constraint for entities that can be filled with an ID.
type IDFillable interface {
	FillID() error
}

// idGeneratorFor returns the ID generator for the given entity type.
func idGeneratorFor(entity IDFillable) (idGenerator, error) {
	generator, ok := _idGenerators[reflect.TypeOf(entity).Elem()]
	if !ok {
		// This should never happen, as the map is initialized with all supported entity types.
		// If it does happen, it is a bug in the code.
		return idGenerator{}, fmt.Errorf("unsupported entity type: '%T'", entity)
	}
	return generator, nil
}
