package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type Attachment struct {
	ent.Schema
}

func (Attachment) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()).
			Annotations(entproto.Field(1)),
	}
}

func (Attachment) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}