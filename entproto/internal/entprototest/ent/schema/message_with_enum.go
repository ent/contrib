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
	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// MessageWithEnum holds the schema definition for the MessageWithEnum entity.
type MessageWithEnum struct {
	ent.Schema
}

// Fields of the MessageWithEnum.
func (MessageWithEnum) Fields() []ent.Field {
	return []ent.Field{
		field.Enum("enum_type").
			Values("pending", "active", "suspended", "deleted").
			Default("pending").
			Annotations(
				entproto.Field(2),
				entproto.Enum(map[string]int32{
					"pending":   0,
					"active":    1,
					"suspended": 2,
					"deleted":   3,
				}),
			),
		field.Enum("enum_without_default").
			Values("first", "second").
			Annotations(
				entproto.Field(3),
				entproto.Enum(map[string]int32{
					"first":  1,
					"second": 2,
				}),
			),
		field.Enum("enum_with_special_characters").
			NamedValues(
				"jpeg", "image/jpeg",
				"png", "image/png").
			Annotations(
				entproto.Field(4),
				entproto.Enum(map[string]int32{
					"image/jpeg": 1,
					"image/png":  2,
				}),
			),
	}
}

func (MessageWithEnum) Annotations() []schema.Annotation {
	return []schema.Annotation{entproto.Message()}
}
