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
	"entgo.io/contrib/entgql"
	"entgo.io/contrib/entgql/internal/todo/ent/schema/annotation"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

// Project holds the schema definition for the GroupTodo entity.
type Project struct {
	ent.Schema
}

// NOTE: This schema intentionally does not have any fields.

// Edges of the Project.
func (Project) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("todo", Todo.Type).
			Annotations(entgql.RelayConnection()),
	}
}

// Annotations returns Group annotations.
func (Project) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField(),
		entgql.Directives(
			annotation.HasPermissions([]string{"ADMIN", "MODERATOR"}),
		),
		entgql.MultiOrder(),
	}
}
