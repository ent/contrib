package mixin

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
	"entgo.io/ent/schema/mixin"
	"github.com/google/uuid"
)

// UUID defines the id field as a UUID string.
type UUID struct {
	mixin.Schema
}

// Fields of the UUID.
func (UUID) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			DefaultFunc(uuid.NewString).
			NotEmpty().
			Immutable().
			Unique(),
	}
}
