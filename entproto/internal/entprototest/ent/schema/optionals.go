package schema

import (
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

type MessageWithOptionals struct {
	ent.Schema
}

func (MessageWithOptionals) Fields() []ent.Field {
	return []ent.Field{
		field.String("str_field").
			Optional().
			Annotations(entproto.Field(2)),
		field.Int8("int_field").
			Optional().
			Annotations(entproto.Field(3)),
		field.Uint8("uint_field").
			Optional().
			Annotations(entproto.Field(4)),
		field.Float32("float_field").
			Optional().
			Annotations(entproto.Field(5)),
		field.Bool("bool_field").
			Optional().
			Annotations(entproto.Field(6)),
		field.Bytes("bytes_field").
			Optional().
			Annotations(entproto.Field(7)),
		field.UUID("uuid_field", uuid.New()).
			Optional().
			Annotations(entproto.Field(8)),
		field.Time("time_field").
			Optional().
			Annotations(entproto.Field(9)),
	}
}

func (MessageWithOptionals) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}
