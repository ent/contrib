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
	todoschema "entgo.io/contrib/entgql/internal/todo/ent/schema"
	"entgo.io/contrib/entgql/internal/todopulid/ent/schema/pulid"
	"entgo.io/ent"
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
		// "PR" declared once.
		pulid.MixinWithPrefix("PR"),
		// Reuse the fields and edges from base example.
		todoschema.Project{},
	}
}
