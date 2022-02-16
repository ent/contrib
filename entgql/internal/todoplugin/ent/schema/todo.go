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

	"entgo.io/ent/schema"
	"github.com/vektah/gqlparser/v2/ast"

	"entgo.io/contrib/entgql"
	"entgo.io/ent"
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
			),
		field.Enum("status").
			NamedValues(
				"InProgress", "IN_PROGRESS",
				"Completed", "COMPLETED",
			).
			Annotations(
				entgql.Annotation{
					OrderField: "STATUS",
					GqlDirectives: []entgql.Directive{
						{
							Name: "someDirective",
							Arguments: []entgql.DirectiveArgument{
								{
									Name:  "stringArg",
									Value: "someString",
									Kind:  ast.StringValue,
								},
								{
									Name:  "boolArg",
									Value: "FALSE",
									Kind:  ast.BooleanValue,
								},
							},
						},
					},
				},
			),
		field.Int("priority").
			Default(0).
			Annotations(
				entgql.OrderField("PRIORITY"),
			),
		field.Text("text").
			NotEmpty().
			Annotations(
				entgql.OrderField("TEXT"),
			),
	}
}

// Edges returns todo edges.
func (Todo) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("children", Todo.Type).
			//nolint SA1019 we keep this as the example.
			Annotations(entgql.Bind()).
			From("parent").
			//nolint SA1019 we keep this as the example.
			Annotations(entgql.Bind()).
			Unique(),
	}
}

// Annotations returns todo annotations.
func (Todo) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Annotation{
			RelayConnection: true,
			GqlDirectives: []entgql.Directive{
				{
					Name: "someDirective",
				},
			},
		},
	}
}
