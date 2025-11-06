// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package schema

import (
	"time"

	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/internal/todoglobalid/ent/schema/customstruct"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Todo defines the todo type schema.
type Todo struct {
	ent.Schema
}

// Fields returns todo fields.
func (Todo) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now).
			Immutable().
			Annotations(
				entgql.OrderField("CREATED_AT"),
				entgql.Skip(entgql.SkipMutationCreateInput),
			),
		field.Enum("status").
			NamedValues(
				"InProgress", "IN_PROGRESS",
				"Completed", "COMPLETED",
				"Pending", "PENDING",
			).
			Annotations(
				entgql.OrderField("STATUS"),
			),
		field.Int("priority").
			Default(0).
			Annotations(
				entgql.OrderField("PRIORITY_ORDER"),
				entgql.MapsTo("priorityOrder"),
			),
		field.Text("text").
			NotEmpty().
			Annotations(
				entgql.OrderField("TEXT"),
			),
		field.Bytes("blob").
			Annotations(
				entgql.Skip(),
			).
			Optional(),
		field.Int("category_id").
			Optional().
			Immutable().
			Annotations(
				entgql.MapsTo("categoryID", "category_id", "categoryX"),
			),
		field.JSON("init", map[string]any{}).
			Optional().
			Annotations(entgql.Type("Map")),
		field.JSON("custom", []customstruct.Custom{}).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput),
				entgql.Skip(entgql.SkipMutationUpdateInput),
			).
			Optional(),
		field.JSON("customp", []*customstruct.Custom{}).
			Annotations(
				entgql.Skip(entgql.SkipMutationCreateInput),
				entgql.Skip(entgql.SkipMutationUpdateInput),
			).
			Optional(),
		field.Int("value").
			Default(0),
	}
}

// Edges returns todo edges.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Todo.Type).
			Annotations(
				entgql.RelayConnection(),
				// For non-unique edges, the order field can be only on edge count.
				// The convention is "UPPER(<edge-name>)_COUNT".
				entgql.OrderField("CHILDREN_COUNT"),
			).
			From("parent").
			Annotations(
				// For unique edges, the order field can be on the edge field that is defined
				// as entgql.OrderField. The convention is "UPPER(<edge-name>)_<gql-order-field>".
				entgql.OrderField("PARENT_STATUS"),
			).
			Unique(),
		edge.From("category", Category.Type).
			Ref("todos").
			Field("category_id").
			Unique().
			Immutable().
			Annotations(
				entgql.OrderField("CATEGORY_TEXT"),
			),
		edge.To("secret", VerySecret.Type).
			Unique(),
	}
}

// Annotations returns Todo annotations.
func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField().Description("This is the todo item"),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.MultiOrder(),
	}
}
