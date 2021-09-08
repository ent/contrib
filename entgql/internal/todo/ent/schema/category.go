// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
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
		field.Text("text").
			NotEmpty().
			Annotations(
				entgql.OrderField("TEXT"),
			),
		field.Enum("status").
			NamedValues(
				"Enabled", "ENABLED",
				"Disabled", "DISABLED",
			),
		field.Other("config", &schematype.CategoryConfig{}).
			SchemaType(map[string]string{
				dialect.SQLite: "json",
			}).
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
				entgql.Type("Uint64"),
			),
	}
}

// Edges of the Category.
func (Category) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("todos", Todo.Type),
	}
}
