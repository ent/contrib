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
	"github.com/google/uuid"
)

type MessageWithOptionals struct {
	ent.Schema
}

func (MessageWithOptionals) Fields() []ent.Field {
	return []ent.Field{
		field.String("str_optional").
			Optional().
			Annotations(entproto.Field(2)),
		field.Int8("int_optional").
			Optional().
			Annotations(entproto.Field(3)),
		field.Uint8("uint_optional").
			Optional().
			Annotations(entproto.Field(4)),
		field.Float32("float_optional").
			Optional().
			Annotations(entproto.Field(5)),
		field.Bool("bool_optional").
			Optional().
			Annotations(entproto.Field(6)),
		field.Bytes("bytes_optional").
			Optional().
			Annotations(entproto.Field(7)),
		field.UUID("uuid_optional", uuid.New()).
			Optional().
			Annotations(entproto.Field(8)),
		field.Time("time_optional").
			Optional().
			Annotations(entproto.Field(9)),
	}
}

func (MessageWithOptionals) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
	}
}
