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
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// ClashingText defines the ClashingText type schema.
// Its name clashes with something we may not be able to use in GraphQL.
type ClashingText struct {
	ent.Schema
}

// Fields returns private fields.
func (ClashingText) Fields() []ent.Field {
	return []ent.Field{
		field.String("content"),
	}
}

func (ClashingText) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.Type("NotClashingText"),
	}
}
