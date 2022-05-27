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

	"entgo.io/contrib/entgql/internal/todopulid/ent/schema/pulid"
	"entgo.io/ent"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Friendship holds the edge schema definition of the Friendship relationship.
type Friendship struct {
	ent.Schema
}

// func (Friendship) Annotations() []schema.Annotation {
// 	return []schema.Annotation{
// 		field.ID("user_id", "friend_id"),
// 	}
// }

// Fields of the Friendship.
func (Friendship) Fields() []ent.Field {
	return []ent.Field{
		field.Time("created_at").
			Default(time.Now),
		field.String("user_id").GoType(pulid.ID("")),
		field.String("friend_id").GoType(pulid.ID("")),
	}
}

// Edges of the Friendship.
func (Friendship) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("user", User.Type).
			Required().
			Unique().
			Field("user_id"),
		edge.To("friend", User.Type).
			Required().
			Unique().
			Field("friend_id"),
	}
}
