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
	"github.com/facebook/ent"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todo/ent/schema"
	"github.com/facebookincubator/ent-contrib/entgql/internal/todopulid/ent/schema/pulid"
)

// Todo defines the todo type schema.
type Todo struct {
	ent.Schema
}

// Mixin returns todo mixed-in schema.
func (Todo) Mixin() []ent.Mixin {
	return []ent.Mixin{
		// "TD" declared once.
		pulid.MixinWithPrefix("TD"),
		// Reuse the fields and edges from base example.
		schema.Todo{},
	}
}
