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

	"entgo.io/contrib/entgql/internal/todo/ent/schema"
	"entgo.io/contrib/entgql/internal/todogotype/ent/schema/bigintgql"
	"entgo.io/ent"
	"entgo.io/ent/schema/field"
)

// Category holds the schema definition for the Category entity.
type Category struct {
	ent.Schema
}

// Mixin returns category mixed-in schema.
func (Category) Mixin() []ent.Mixin {
	return []ent.Mixin{
		// Reuse the fields and edges from base example.
		schema.Category{},
	}
}

// Fields returns category fields.
func (Category) Fields() []ent.Field {
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
