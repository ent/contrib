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
	"entgo.io/contrib/entoas"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

// Edges of the User.
func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("friends", User.Type).
			Annotations(entoas.Groups("user")),
		edge.To("following", User.Type).
			Annotations(entoas.Groups("user")).
			From("followers").
			Annotations(entoas.Groups("user")),
		edge.To("children", User.Type).
			Annotations(entoas.Groups("user")).
			From("parent").
			Unique().
			Annotations(entoas.Groups("user")),
	}
}

// Annotations of the User.
func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entoas.ReadOperation(entoas.OperationGroups("user")),
		// entoas.ListOperation(entoas.OperationGroups("user")),
	}
}
