package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
)

// TwoMethodService holds the schema definition for the TwoMethodService entity.
type TwoMethodService struct {
	ent.Schema
}

// Fields of the TwoMethodService.
func (TwoMethodService) Fields() []ent.Field {
	return nil
}

// Edges of the TwoMethodService.
func (TwoMethodService) Edges() []ent.Edge {
	return nil
}

func (TwoMethodService) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(
			entproto.Methods(
				entproto.MethodCreate |
					entproto.MethodGet,
			),
		),
	}
}
