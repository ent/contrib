package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// AllMethodsService holds the schema definition for the AllMethodsService entity.
type AllMethodsService struct {
	ent.Schema
}

// Fields of the AllMethodsService.
func (AllMethodsService) Fields() []ent.Field {
	return nil
}

// Edges of the AllMethodsService.
func (AllMethodsService) Edges() []ent.Edge {
	return nil
}

func (AllMethodsService) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(
			entproto.Methods(
				entproto.MethodAll,
			),
		),
	}
}
