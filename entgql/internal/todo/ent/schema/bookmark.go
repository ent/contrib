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
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
)

// Bookmark demonstrates a 1:1 polymorphic interface field: its `item` field
// resolves to a single Todo or Project, both of which implement the BookmarkItem
// interface.
type Bookmark struct {
	ent.Schema
}

// Fields of the Bookmark.
func (Bookmark) Fields() []ent.Field {
	return []ent.Field{
		field.String("name"),
	}
}

// Edges of the Bookmark. The two to-one edges share the "item" InterfaceField, so
// entgql collapses them into a single `item: BookmarkItem` field that resolves to
// whichever target is set (at most one).
func (Bookmark) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("todo", Todo.Type).
			Unique().
			Annotations(entgql.InterfaceField("item")),
		edge.To("project", Project.Type).
			Unique().
			Annotations(entgql.InterfaceField("item")),
	}
}

// Annotations of the Bookmark.
func (Bookmark) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entgql.RelayConnection(),
		entgql.QueryField("bookmarks"),
	}
}
