package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// OneMethodService holds the schema definition for the OneMethodService entity.
type OneMethodService struct {
	ent.Schema
}

// Fields of the OneMethodService.
func (OneMethodService) Fields() []ent.Field {
	return nil
}

// Edges of the OneMethodService.
func (OneMethodService) Edges() []ent.Edge {
	return nil
}

func (OneMethodService) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(
			entproto.Methods(
				entproto.MethodGet,
			),
		),
	}
}
