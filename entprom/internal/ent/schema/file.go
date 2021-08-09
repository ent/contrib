package schema

import (
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// File holds the schema definition for the File entity.
type File struct {
	ent.Schema
}

// Fields of the File.
func (File) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
		field.Bool("deleted").
			Default(false),
		field.Int("parent_id").
			Optional(),
		field.Int("owner_id").
			Optional(),
	}
}

// Edges of the File.
func (File) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", File.Type).
			From("parent").
			Unique().
			Field("parent_id"),
		edge.From("owner", User.Type).
			Ref("files").
			Unique().
			Field("owner_id"),
	}
}
