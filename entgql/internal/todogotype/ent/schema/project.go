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
	"math/big"
	"math/rand"

	todoschema "entgo.io/contrib/entgql/internal/todo/ent/schema"
	"entgo.io/contrib/entgql/internal/todogotype/ent/schema/bigintgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Project defines the project type schema. It reuses the base example's edges
// and annotations so the Category.projects edge (and the BookmarkItem interface it
// participates in) resolves in this example too.
type Project struct {
	ent.Schema
}

// Mixin returns project mixed-in schema.
func (Project) Mixin() []ent.Mixin {
	return []ent.Mixin{
		// Reuse the fields and edges from base example.
		todoschema.Project{},
	}
}

// Fields returns project fields, overriding the base int id with the custom
// BigInt id type used by the other types in this example.
func (Project) Fields() []ent.Field {
	return []ent.Field{
		field.String("id").
			GoType(bigintgql.BigInt{}).
			Unique().
			Immutable().
			DefaultFunc(func() bigintgql.BigInt {
				//nolint:gosec
				return bigintgql.BigInt{Int: big.NewInt(int64(rand.Float64()))}
			}),
	}
}
