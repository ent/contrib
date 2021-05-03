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

package schemast

import (
	"testing"

	"entgo.io/contrib/entproto"
	"entgo.io/ent"
	"entgo.io/ent/schema"
	"entgo.io/ent/schema/edge"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
)

func TestUpsert(t *testing.T) {
	tt, err := newPrintTest(t)
	require.NoError(t, err)
	mutations := []Mutator{
		&UpsertSchema{
			Name: "User", // An existing schema
			Fields: []ent.Field{
				field.String("name"),
			},
			Edges: []ent.Edge{
				WithType(edge.From("administered", placeholder.Type).Ref("admin"), "Team"),
			},
			Annotations: []schema.Annotation{
				entproto.Message(),
			},
		},
		&UpsertSchema{
			Name: "Team", // A new schema
			Fields: []ent.Field{
				field.String("name"),
			},
			Edges: []ent.Edge{
				WithType(edge.To("admin", placeholder.Type), "User"),
			},
			Annotations: []schema.Annotation{
				entproto.Message(),
			},
		},
	}
	err = Mutate(tt.ctx, mutations...)
	require.NoError(t, err)
	require.NoError(t, tt.print())
	require.NoError(t, tt.load())

	team := tt.getType("Team")
	require.NotNil(t, team)
	require.Len(t, team.Edges, 1)
	require.Len(t, team.Annotations, 1)
	user := tt.getType("User")
	require.NotNil(t, user)
	require.Len(t, user.Fields, 1)
	require.Len(t, user.Edges, 1)
	require.Len(t, user.Annotations, 1)
}

func WithType(e ent.Edge, typeName string) ent.Edge {
	e.Descriptor().Type = typeName
	return e
}

type placeholder struct {
}

func (placeholder) Type() {

}
