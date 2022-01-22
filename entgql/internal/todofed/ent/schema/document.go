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
	"entgo.io/contrib/entgql"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
	"github.com/google/uuid"
)

// Document defines the document type schema.
type Document struct {
	ent.Schema
}

// Fields returns private fields.
func (Document) Fields() []ent.Field {
	return []ent.Field{
		field.UUID("id", uuid.New()),
		field.Int("global_id"),
		field.String("name"),
	}
}

func (Document) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.NodeIDField("global_id"),
		// // nodeIDField: unable to find field field_not_exist
		// entgql.NodeIDField("field_not_exist"),
	}
}
