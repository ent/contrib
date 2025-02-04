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
	"entgo.io/ent"
	"entgo.io/ent/dialect"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"

	"entgo.io/contrib/entgql/internal/todo/ent/schema/schematype"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Fields of the Category.
func (Category) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Annotations(
				// Setting the OrderField explicitly on the "ID"
				// field, adds it to the generated GraphQL schema.
				entgql.OrderField("ID"),
			),
		field.Text("text").
			NotEmpty().
			Annotations(
				entgql.OrderField("TEXT"),
			),
		field.Enum("status").
			NamedValues(
				"Enabled", "ENABLED",
				"Disabled", "DISABLED",
			).
			Annotations(
				entgql.Type("CategoryStatus"),
				entgql.OrderField("STATUS"),
			),
		field.Other("config", &schematype.CategoryConfig{}).
			SchemaType(map[string]string{
				dialect.SQLite: "json",
			}).
			Optional(),
		field.JSON("types", &schematype.CategoryTypes{}).
			Optional(),
		field.Int64("duration").
			GoType(time.Duration(0)).
			Optional().
			Annotations(
				entgql.OrderField("DURATION"),
				entgql.Type("Duration"),
			),
		field.Uint64("count").
			Optional().
			Annotations(
				entgql.OrderField("COUNT"),
				entgql.Type("Uint64"),
			),
		field.Strings("strings").
			Optional().
			Deprecated("use `string` instead"),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("todos", Todo.Type).
			Annotations(
				entgql.RelayConnection(),
				entgql.OrderField("TODOS_COUNT"),
			),
		edge.To("sub_categories", Category.Type).
			Annotations(entgql.RelayConnection()),
	}
}

// Annotations returns Todo annotations.
func (Category) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.RelayConnection(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
		entgql.MultiOrder(),
	}
}
