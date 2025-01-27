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

// DirectiveExample holds the schema definition for the DirectiveExample entity.
type DirectiveExample struct {
	ent.Schema
}

func fieldDirective() entgql.Directive { return entgql.NewDirective("fieldDirective") }

// Fields of the DirectiveExample.
func (DirectiveExample) Fields() []ent.Field {
	return []ent.Field{
		field.Text("on_type_field").
			Optional().
			Annotations(
				entgql.Directives(fieldDirective()),
			),
		field.Text("on_mutation_fields").
			Optional().
			Annotations(
				entgql.Directives(fieldDirective().OnCreateMutationField().OnUpdateMutationField().SkipOnTypeField()),
			),
		field.Text("on_mutation_create").
			Optional().
			Annotations(
				entgql.Directives(fieldDirective().OnCreateMutationField().SkipOnTypeField()),
			),
		field.Text("on_mutation_update").
			Optional().
			Annotations(
				entgql.Directives(fieldDirective().OnUpdateMutationField().SkipOnTypeField()),
			),
		field.Text("on_all_fields").
			Optional().
			Annotations(
				entgql.Directives(fieldDirective().OnCreateMutationField().OnUpdateMutationField()),
			),
	}
}

// Annotations returns Todo annotations.
func (DirectiveExample) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.QueryField(),
		entgql.Mutations(entgql.MutationCreate(), entgql.MutationUpdate()),
	}
}
