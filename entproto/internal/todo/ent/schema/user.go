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
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// User holds the schema definition for the User entity.
type User struct {
	ent.Schema
}

func (User) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}

// Fields of the User.
func (User) Fields() []ent.Field {
	return []ent.Field{
		field.String("user_name").
			Unique().
			Annotations(entproto.Field(2)),
		field.Time("joined").
			Annotations(entproto.Field(3)),
		field.Uint("points").
			Annotations(entproto.Field(4)),
		field.Uint64("exp").
			Annotations(entproto.Field(5)),
		field.Enum("status").
			Values("pending", "active").
			Annotations(
				entproto.Field(6),
				entproto.Enum(map[string]int32{
					"pending": 1,
					"active":  2,
				}),
			),
		field.Int("external_id").
			Unique().
			Annotations(entproto.Field(8)),
		field.UUID("crm_id", uuid.New()).
			Annotations(entproto.Field(9)),
		field.Bool("banned").
			Default(false).
			Annotations(entproto.Field(10)),
	}
}

func (User) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("group", Group.Type).
			Unique().
			Annotations(
				entproto.Field(7),
			),
		edge.To("attachment", Attachment.Type).
			Unique().
			Annotations(
				entproto.Field(11),
			),
	}
}
