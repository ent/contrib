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

package entgql

import (
	"os"
	"reflect"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEntGQL_buildTypes(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)
	disableRelayConnection(graph)
	plugin := &schemaGenerator{genSchema: true, genMutations: true}
	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(graph, schema)
	require.NoError(t, err)
	schemaExpect, err := os.ReadFile("./testdata/schema.graphql")
	require.NoError(t, err)
	output := printSchema(schema)
	if string(schemaExpect) != output {
		require.NoError(t, os.WriteFile("./testdata/schema_output.graphql", []byte(output), 0644))
	}
	require.Equal(t, string(schemaExpect), output)
}

func TestEntGQL_buildTypes_todoplugin_relay(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)
	plugin := &schemaGenerator{genSchema: true, genWhereInput: true, genMutations: true, relaySpec: true}
	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(graph, schema)
	require.NoError(t, err)
	schemaExpect, err := os.ReadFile("./testdata/schema_relay.graphql")
	require.NoError(t, err)
	output := printSchema(schema)
	if string(schemaExpect) != output {
		require.NoError(t, os.WriteFile("./testdata/schema_relay_output.graphql", []byte(output), 0644))
	}
	require.Equal(t, string(schemaExpect), output)
}

func TestSchema_relayConnectionTypes(t *testing.T) {
	type args struct {
		t *gen.Type
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Todo",
			args: args{
				t: &gen.Type{
					Name: "Todo",
				},
			},
			want: `"""
A connection to a list of items.
"""
type TodoConnection {
  """
  A list of edges.
  """
  edges: [TodoEdge]
  """
  Information to aid in pagination.
  """
  pageInfo: PageInfo!
  """
  Identifies the total count of items in the connection.
  """
  totalCount: Int!
}
"""
An edge in a connection.
"""
type TodoEdge {
  """
  The item at the end of the edge.
  """
  node: Todo
  """
  A cursor for use in pagination.
  """
  cursor: Cursor!
}
`,
		},
		{
			name: "Todo_with_type",
			args: args{
				t: &gen.Type{
					Name: "Todo",
					Annotations: map[string]interface{}{
						annotationName: map[string]interface{}{
							"Type": "SuperTodo",
						},
					},
				},
			},
			want: `"""
A connection to a list of items.
"""
type SuperTodoConnection {
  """
  A list of edges.
  """
  edges: [SuperTodoEdge]
  """
  Information to aid in pagination.
  """
  pageInfo: PageInfo!
  """
  Identifies the total count of items in the connection.
  """
  totalCount: Int!
}
"""
An edge in a connection.
"""
type SuperTodoEdge {
  """
  The item at the end of the edge.
  """
  node: SuperTodo
  """
  A cursor for use in pagination.
  """
  cursor: Cursor!
}
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := relayConnectionTypes(tt.args.t)
			if (err != nil) != tt.wantErr {
				t.Errorf("relayConnection() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

			s := &ast.Schema{}
			s.AddTypes(got...)
			gots := printSchema(s)
			if !reflect.DeepEqual(gots, tt.want) {
				t.Errorf("relayConnection() = %v, want %v", gots, tt.want)
			}
		})
	}
}

func TestSchema_relayBuiltinTypes(t *testing.T) {
	tests := []struct {
		name string
		want string
	}{
		{
			name: "relayBuiltinTypes",
			want: `"""
Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor
"""
scalar Cursor
"""
An object with an ID.
Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)
"""
interface Node @goModel(model: "todo/ent.Noder") {
  """
  The id of the object.
  """
  id: ID!
}
"""
Information about pagination in a connection.
https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo
"""
type PageInfo {
  """
  When paginating forwards, are there more items?
  """
  hasNextPage: Boolean!
  """
  When paginating backwards, are there more items?
  """
  hasPreviousPage: Boolean!
  """
  When paginating backwards, the cursor to continue.
  """
  startCursor: Cursor
  """
  When paginating forwards, the cursor to continue.
  """
  endCursor: Cursor
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relayBuiltinTypes("todo/ent")

			s := &ast.Schema{}
			s.AddTypes(got...)
			gots := printSchema(s)
			if !reflect.DeepEqual(gots, tt.want) {
				t.Errorf("relayBuiltinTypes() = %v, want %v", gots, tt.want)
			}
		})
	}
}

func disableRelayConnection(g *gen.Graph) {
	disable := func(a gen.Annotations) {
		if ant, ok := a[annotationName]; ok {
			if m, ok := ant.(map[string]interface{}); ok {
				m["RelayConnection"] = false
			}
		}
	}

	for _, n := range g.Nodes {
		disable(n.Annotations)
		for _, f := range n.Fields {
			disable(f.Annotations)
		}
		for _, e := range n.Edges {
			disable(e.Annotations)
		}
	}
}
