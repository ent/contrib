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
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/bionicstork/contrib/entproto"
)

type BlogPost struct {
	ent.Schema
}

func (BlogPost) Edges() []ent.Edge {
	return []ent.Edge{
		edge.To("author", User.Type).
			Unique().
			Annotations(entproto.Field(4)),
		edge.From("categories", Category.Type).
			Ref("blog_posts").
			Annotations(entproto.Field(5)),
	}
}

func (BlogPost) Fields() []ent.Field {
	return []ent.Field{
		field.String("title").
			Annotations(entproto.Field(2)),
		field.String("body").
			Annotations(entproto.Field(3)),
		field.Int("external_id").
			Unique().
			Annotations(entproto.Field(7)),
	}
}

func (BlogPost) Annotations() []schema.Annotation {
	return []schema.Annotation{
		entproto.Message(),
		entproto.Service(),
	}
}
