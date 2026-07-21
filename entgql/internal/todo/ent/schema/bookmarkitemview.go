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
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/entsql"
	"entgo.io/ent/dialect/sql"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/field"
)

// BookmarkItemView is a SQL view that unions the BookmarkItem members (Todo and
// Project) to back the global `bookmarkItems` connection with DB-level pagination.
//
// It carries a "kind" discriminator that routes each row to its concrete type,
// plus the interface's shared columns (text, name) so a page whose selection
// touches only shared fields is served from the view rows alone.
type BookmarkItemView struct {
	ent.View
}

// Annotations of the BookmarkItemView.
func (BookmarkItemView) Annotations() []schema.Annotation {
	return []schema.Annotation{
		// Unions the members, tagging each row with a static "kind" discriminator
		// (matching the concrete ent type names) and projecting the shared columns
		// into a common shape. The built-in migrator skips views; tests create this
		// one to match.
		entsql.ViewFor(dialect.SQLite, func(s *sql.Selector) {
			s.Select("id").
				AppendSelectExprAs(sql.Raw("'Todo'"), "kind").
				AppendSelect("category_id", "text", "name").
				From(sql.Table("todos")).
				UnionAll(
					sql.Dialect(dialect.SQLite).
						Select("id").
						AppendSelectExprAs(sql.Raw("'Project'"), "kind").
						AppendSelectAs("category_projects", "category_id").
						AppendSelect("text", "name").
						From(sql.Table("projects")),
				)
		}),
		// Type names the interface this view backs; QueryField exposes it as the
		// top-level `bookmarkItems` connection. entgql generates a polymorphic
		// connection that paginates the view (with orderBy/where) and resolves each
		// row to its concrete node.
		entgql.Type("BookmarkItem"),
		entgql.QueryField("bookmarkItems"),
	}
}

// Fields of the BookmarkItemView. "kind" discriminates the row's concrete type;
// "text"/"name" mirror the interface's shared scalar fields.
func (BookmarkItemView) Fields() []ent.Field {
	return []ent.Field{
		field.Int("id").
			Annotations(entgql.OrderField("ID")),
		// The concrete-type discriminator: entgql.MapsTo("__typename") marks the
		// column whose value routes each row to its member type.
		field.String("kind").
			Annotations(entgql.MapsTo("__typename")),
		field.Int("category_id").
			Optional().
			Annotations(
				entgql.OrderField("CATEGORY_ID"),
				// The FK behind the shared "owner" edge. Naming that interface field
				// lets entgql verify every implementor exposes "owner".
				entgql.InterfaceField("owner"),
			),
		field.Text("text"),
		field.String("name").
			Optional(),
	}
}
