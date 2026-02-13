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

func TestEntGQL_BuildSplitSchema(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	plugin := &schemaGenerator{genSchema: true, genWhereInput: true, genMutations: true, relaySpec: true}
	split, err := plugin.BuildSplitSchema(graph)
	require.NoError(t, err)
	require.NotNil(t, split)

	// Verify shared schema contains expected types
	require.NotNil(t, split.Shared)
	require.NotNil(t, split.Shared.Types["OrderDirection"], "shared schema should contain OrderDirection")
	require.NotNil(t, split.Shared.Types["NullsDirection"], "shared schema should contain NullsDirection")
	require.NotNil(t, split.Shared.Types["Cursor"], "shared schema should contain Cursor")
	require.NotNil(t, split.Shared.Types["Node"], "shared schema should contain Node")
	require.NotNil(t, split.Shared.Types["PageInfo"], "shared schema should contain PageInfo")

	// Verify directives are in shared schema
	require.NotNil(t, split.Shared.Directives["goModel"], "shared schema should contain goModel directive")
	require.NotNil(t, split.Shared.Directives["goField"], "shared schema should contain goField directive")

	// Verify query schema
	require.NotNil(t, split.Query)
	require.NotNil(t, split.Query.Types["Query"], "query schema should contain Query type")
	queryType := split.Query.Types["Query"]
	require.True(t, len(queryType.Fields) > 0, "Query type should have fields")

	// Verify entities exist
	require.NotNil(t, split.Entities, "entities should not be nil")
	require.True(t, len(split.Entities) > 0, "should have at least one entity")

	// Verify Todo entity has expected types
	todoSchema, exists := split.Entities["Todo"]
	require.True(t, exists, "should have Todo entity")
	require.NotNil(t, todoSchema.Types["Todo"], "Todo schema should have Todo type")
	require.NotNil(t, todoSchema.Types["TodoConnection"], "Todo schema should have TodoConnection type")
	require.NotNil(t, todoSchema.Types["TodoEdge"], "Todo schema should have TodoEdge type")
	require.NotNil(t, todoSchema.Types["TodoWhereInput"], "Todo schema should have TodoWhereInput type")
}

func TestEntGQL_BuildSplitSchema_PreservesSchemaHooks(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	plugin := &schemaGenerator{
		genSchema:     true,
		genWhereInput: true,
		genMutations:  true,
		relaySpec:     true,
		schemaHooks: []SchemaHook{func(_ *gen.Graph, schema *ast.Schema) error {
			schema.AddTypes(&ast.Definition{
				Name: "HookScalar",
				Kind: ast.Scalar,
			})
			todoType := schema.Types["Todo"]
			require.NotNil(t, todoType)
			todoType.Fields = append(todoType.Fields, &ast.FieldDefinition{
				Name: "hookField",
				Type: ast.NamedType("HookScalar", nil),
			})
			return nil
		}},
	}

	split, err := plugin.BuildSplitSchema(graph)
	require.NoError(t, err)

	require.NotNil(t, split.Shared.Types["HookScalar"], "hook-added scalar should be preserved in split shared schema")
	todoSchema, ok := split.Entities["Todo"]
	require.True(t, ok, "Todo entity should exist in split schema")
	todoType := todoSchema.Types["Todo"]
	require.NotNil(t, todoType)

	var hookFieldExists bool
	for _, field := range todoType.Fields {
		if field.Name == "hookField" {
			hookFieldExists = true
			break
		}
	}
	require.True(t, hookFieldExists, "hook-added field should be preserved in split Todo type")
}

func TestEntGQL_SplitSchema_SharedTypesNotInEntities(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	plugin := &schemaGenerator{genSchema: true, genWhereInput: true, genMutations: true, relaySpec: true}
	split, err := plugin.BuildSplitSchema(graph)
	require.NoError(t, err)

	// Verify shared types are NOT in entity schemas
	for name, entitySchema := range split.Entities {
		require.Nil(t, entitySchema.Types["OrderDirection"], "entity %s should not contain OrderDirection", name)
		require.Nil(t, entitySchema.Types["NullsDirection"], "entity %s should not contain NullsDirection", name)
		require.Nil(t, entitySchema.Types["Cursor"], "entity %s should not contain Cursor", name)
		require.Nil(t, entitySchema.Types["Node"], "entity %s should not contain Node", name)
		require.Nil(t, entitySchema.Types["PageInfo"], "entity %s should not contain PageInfo", name)
		require.Nil(t, entitySchema.Types["Query"], "entity %s should not contain Query", name)
	}
}

func TestEntGQL_WithSchemaDir_MutualExclusivity(t *testing.T) {
	// Test that WithSchemaDir and WithSchemaPath are mutually exclusive
	_, err := NewExtension(
		WithSchemaPath("./test.graphql"),
		WithSchemaDir("./graphql"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")

	// Test the reverse order
	_, err = NewExtension(
		WithSchemaDir("./graphql"),
		WithSchemaPath("./test.graphql"),
	)
	require.Error(t, err)
	require.Contains(t, err.Error(), "mutually exclusive")
}

func TestEntGQL_SplitSchema_FileNaming(t *testing.T) {
	// Test that entity names are properly snake_cased for file naming
	require.Equal(t, "todo", snake("Todo"))
	require.Equal(t, "user", snake("User"))
	require.Equal(t, "bill_product", snake("BillProduct"))
	require.Equal(t, "very_secret", snake("VerySecret"))
}

func TestEntGQL_generateSplitSchema_Files(t *testing.T) {
	// Create a temporary directory for output
	tmpDir, err := os.MkdirTemp("", "entgql_split_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSchemaDir(tmpDir),
		WithWhereInputs(true),
	)
	require.NoError(t, err)

	// Manually call the generateSplitSchema method
	err = ex.generateSplitSchema(graph)
	require.NoError(t, err)

	// Verify ent_shared.graphql exists and contains expected content
	sharedContent, err := os.ReadFile(tmpDir + "/ent_shared.graphql")
	require.NoError(t, err)
	require.Contains(t, string(sharedContent), "OrderDirection")
	require.Contains(t, string(sharedContent), "NullsDirection")
	require.Contains(t, string(sharedContent), "Cursor")
	require.Contains(t, string(sharedContent), "Node")
	require.Contains(t, string(sharedContent), "PageInfo")

	// Verify ent_query.graphql exists
	queryContent, err := os.ReadFile(tmpDir + "/ent_query.graphql")
	require.NoError(t, err)
	require.Contains(t, string(queryContent), "type Query")

	// Verify ent_todo.graphql exists
	todoContent, err := os.ReadFile(tmpDir + "/ent_todo.graphql")
	require.NoError(t, err)
	require.Contains(t, string(todoContent), "type Todo")
	require.Contains(t, string(todoContent), "TodoConnection")
	require.Contains(t, string(todoContent), "TodoEdge")
	require.Contains(t, string(todoContent), "TodoWhereInput")

	// Verify shared types are not in entity files
	require.NotContains(t, string(todoContent), "type Node")
	require.NotContains(t, string(todoContent), "enum OrderDirection")
	require.NotContains(t, string(todoContent), "type PageInfo")
}

func TestEntGQL_generateSplitSchema_RemovesStaleFiles(t *testing.T) {
	tmpDir, err := os.MkdirTemp("", "entgql_split_cleanup_test")
	require.NoError(t, err)
	defer os.RemoveAll(tmpDir)

	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)

	staleSchema := tmpDir + "/ent_stale.graphql"
	staleQuery := tmpDir + "/ent_query.graphql"
	require.NoError(t, os.WriteFile(staleSchema, []byte("type Stale { id: ID! }\n"), 0644))
	require.NoError(t, os.WriteFile(staleQuery, []byte("type Query { stale: String }\n"), 0644))

	ex, err := NewExtension(
		WithSchemaGenerator(),
		WithSchemaDir(tmpDir),
		WithWhereInputs(true),
	)
	require.NoError(t, err)

	err = ex.generateSplitSchema(graph)
	require.NoError(t, err)

	_, err = os.Stat(staleSchema)
	require.True(t, os.IsNotExist(err), "stale split schema file should be removed")

	_, err = os.Stat(staleQuery)
	require.NoError(t, err, "query schema file should be rewritten when query type exists")
}
