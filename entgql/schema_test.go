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

package entgql

import (
	"reflect"
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEntGQL_buildTypes(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	disableRelayConnection(graph)
	plugin, err := newSchemaGenerator(graph)
	require.NoError(t, err)
	plugin.relaySpec = false

	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(schema)
	require.NoError(t, err)

	require.Equal(t, `type Category implements Entity {
  id: ID!
  text: String!
  uuidA: UUID
  status: CategoryStatus!
  config: CategoryConfig!
  duration: Duration!
  count: Uint64! @deprecated(reason: "We don't use this field anymore")
  strings: [String!]
  todos: [Todo!]
}
"""CategoryStatus is enum for the field status"""
enum CategoryStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status") {
  ENABLED
  DISABLED
}
type MasterUser @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent.User") {
  id: ID!
  username: String!
  age: Float!
  amount: Float!
  role: Role!
  nullableString: String
}
"""Role is enum for the field role"""
enum Role @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role") {
  ADMIN
  USER
  UNKNOWN
}
"""Status is enum for the field status"""
enum Status @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status") {
  IN_PROGRESS
  COMPLETED
}
type Todo {
  id: ID!
  createdAt: Time!
  visibilityStatus: VisibilityStatus!
  status: Status!
  priority: Int!
  text: String!
  parent: Todo
  childrenConnection: [Todo!] @goField(name: "children", forceResolver: false)
  children: [Todo!]
}
"""VisibilityStatus is enum for the field visibility_status"""
enum VisibilityStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus") {
  LISTING
  HIDDEN
}
`, printSchema(schema))
}

func TestEntGQL_buildTypes_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin, err := newSchemaGenerator(graph)

	require.NoError(t, err)
	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(schema)
	require.NoError(t, err)

	require.Equal(t, `type Category implements Node & Entity {
  id: ID!
  text: String!
  uuidA: UUID
  status: CategoryStatus!
  config: CategoryConfig!
  duration: Duration!
  count: Uint64! @deprecated(reason: "We don't use this field anymore")
  strings: [String!]
  todos: [Todo!]
}
"""A connection to a list of items."""
type CategoryConnection {
  """A list of edges."""
  edges: [CategoryEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  totalCount: Int!
}
"""An edge in a connection."""
type CategoryEdge {
  """The item at the end of the edge."""
  node: Category
  """A cursor for use in pagination."""
  cursor: Cursor!
}
input CategoryOrder {
  direction: OrderDirection! = ASC
  field: CategoryOrderField!
}
enum CategoryOrderField {
  TEXT
  DURATION
}
"""CategoryStatus is enum for the field status"""
enum CategoryStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status") {
  ENABLED
  DISABLED
}
type MasterUser implements Node @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent.User") {
  id: ID!
  username: String!
  age: Float!
  amount: Float!
  role: Role!
  nullableString: String
}
"""A connection to a list of items."""
type MasterUserConnection {
  """A list of edges."""
  edges: [MasterUserEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  totalCount: Int!
}
"""An edge in a connection."""
type MasterUserEdge {
  """The item at the end of the edge."""
  node: MasterUser
  """A cursor for use in pagination."""
  cursor: Cursor!
}
"""Role is enum for the field role"""
enum Role @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role") {
  ADMIN
  USER
  UNKNOWN
}
"""Status is enum for the field status"""
enum Status @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status") {
  IN_PROGRESS
  COMPLETED
}
type Todo implements Node {
  id: ID!
  createdAt: Time!
  visibilityStatus: VisibilityStatus!
  status: Status!
  priority: Int!
  text: String!
  parent: Todo
  childrenConnection(after: Cursor, first: Int, before: Cursor, last: Int, orderBy: TodoOrder): TodoConnection! @goField(name: "children", forceResolver: false)
  children(after: Cursor, first: Int, before: Cursor, last: Int, orderBy: TodoOrder): TodoConnection!
}
"""A connection to a list of items."""
type TodoConnection {
  """A list of edges."""
  edges: [TodoEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  totalCount: Int!
}
"""An edge in a connection."""
type TodoEdge {
  """The item at the end of the edge."""
  node: Todo
  """A cursor for use in pagination."""
  cursor: Cursor!
}
input TodoOrder {
  direction: OrderDirection! = ASC
  field: TodoOrderField!
}
enum TodoOrderField {
  CREATED_AT
  VISIBILITY_STATUS
  STATUS
  PRIORITY
  TEXT
}
"""VisibilityStatus is enum for the field visibility_status"""
enum VisibilityStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus") {
  LISTING
  HIDDEN
}
`, printSchema(schema))
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
			want: `"""A connection to a list of items."""
type TodoConnection {
  """A list of edges."""
  edges: [TodoEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  totalCount: Int!
}
"""An edge in a connection."""
type TodoEdge {
  """The item at the end of the edge."""
  node: Todo
  """A cursor for use in pagination."""
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
			want: `"""A connection to a list of items."""
type SuperTodoConnection {
  """A list of edges."""
  edges: [SuperTodoEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  totalCount: Int!
}
"""An edge in a connection."""
type SuperTodoEdge {
  """The item at the end of the edge."""
  node: SuperTodo
  """A cursor for use in pagination."""
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
interface Node {
  """The id of the object."""
  id: ID!
}
"""
Information about pagination in a connection.
https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo
"""
type PageInfo {
  """When paginating forwards, are there more items?"""
  hasNextPage: Boolean!
  """When paginating backwards, are there more items?"""
  hasPreviousPage: Boolean!
  """When paginating backwards, the cursor to continue."""
  startCursor: Cursor
  """When paginating forwards, the cursor to continue."""
  endCursor: Cursor
}
`,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got := relayBuiltinTypes()

			s := &ast.Schema{}
			s.AddTypes(got...)
			gots := printSchema(s)
			if !reflect.DeepEqual(gots, tt.want) {
				t.Errorf("relayBuiltinTypes() = %v, want %v", gots, tt.want)
			}
		})
	}
}

func TestModifyConfig_empty(t *testing.T) {
	e, err := newSchemaGenerator(&gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
		},
	})
	require.NoError(t, err)
	e.relaySpec = false

	cfg, err := e.genModels()
	require.NoError(t, err)

	expected := map[string]string{}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig(t *testing.T) {
	e, err := newSchemaGenerator(createGraph(false))
	require.NoError(t, err)

	e.relaySpec = false
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Todo":          "example.com.Todo",
		"Group":         "example.com.Group",
		"GroupWithSort": "example.com.GroupWithSort",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_relay(t *testing.T) {
	g := createGraph(true)
	e, err := newSchemaGenerator(g)
	e.relaySpec = false
	require.NoError(t, err)
	_, err = e.genModels()
	require.Error(t, err, ErrRelaySpecDisabled)

	e.relaySpec = true
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Cursor":                  "example.com.Cursor",
		"Group":                   "example.com.Group",
		"GroupConnection":         "example.com.GroupConnection",
		"GroupEdge":               "example.com.GroupEdge",
		"GroupWithSort":           "example.com.GroupWithSort",
		"GroupWithSortConnection": "example.com.GroupWithSortConnection",
		"GroupWithSortEdge":       "example.com.GroupWithSortEdge",
		"GroupWithSortOrder":      "example.com.GroupWithSortOrder",
		"GroupWithSortOrderField": "example.com.GroupWithSortOrderField",
		"Node":                    "example.com.Noder",
		"OrderDirection":          "example.com.OrderDirection",
		"PageInfo":                "example.com.PageInfo",
		"Todo":                    "example.com.Todo",
		"TodoConnection":          "example.com.TodoConnection",
		"TodoEdge":                "example.com.TodoEdge",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	disableRelayConnection(graph)

	e, err := newSchemaGenerator(graph)
	require.NoError(t, err)
	e.relaySpec = false

	cfg, err := e.genModels()
	require.NoError(t, err)

	expected := map[string]string{
		"Category":         "entgo.io/contrib/entgql/internal/todoplugin/ent.Category",
		"CategoryStatus":   "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status",
		"MasterUser":       "entgo.io/contrib/entgql/internal/todoplugin/ent.User",
		"Role":             "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role",
		"Status":           "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status",
		"Todo":             "entgo.io/contrib/entgql/internal/todoplugin/ent.Todo",
		"VisibilityStatus": "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus",
	}
	require.Equal(t, expected, cfg)
}

func TestModifyConfig_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)

	e, err := newSchemaGenerator(graph)
	require.NoError(t, err)
	cfg, err := e.genModels()
	require.NoError(t, err)
	expected := map[string]string{
		"Category":             "entgo.io/contrib/entgql/internal/todoplugin/ent.Category",
		"CategoryConnection":   "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryConnection",
		"CategoryEdge":         "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryEdge",
		"CategoryOrder":        "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrder",
		"CategoryOrderField":   "entgo.io/contrib/entgql/internal/todoplugin/ent.CategoryOrderField",
		"CategoryStatus":       "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status",
		"Cursor":               "entgo.io/contrib/entgql/internal/todoplugin/ent.Cursor",
		"MasterUser":           "entgo.io/contrib/entgql/internal/todoplugin/ent.User",
		"MasterUserConnection": "entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserConnection",
		"MasterUserEdge":       "entgo.io/contrib/entgql/internal/todoplugin/ent.MasterUserEdge",
		"Node":                 "entgo.io/contrib/entgql/internal/todoplugin/ent.Noder",
		"OrderDirection":       "entgo.io/contrib/entgql/internal/todoplugin/ent.OrderDirection",
		"PageInfo":             "entgo.io/contrib/entgql/internal/todoplugin/ent.PageInfo",
		"Role":                 "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role",
		"Status":               "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.Status",
		"Todo":                 "entgo.io/contrib/entgql/internal/todoplugin/ent.Todo",
		"TodoConnection":       "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoConnection",
		"TodoEdge":             "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoEdge",
		"TodoOrder":            "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrder",
		"TodoOrderField":       "entgo.io/contrib/entgql/internal/todoplugin/ent.TodoOrderField",
		"VisibilityStatus":     "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus",
	}
	require.Equal(t, expected, cfg)
}

func createGraph(relayConnection bool) *gen.Graph {
	return &gen.Graph{
		Config: &gen.Config{
			Package: "example.com",
			IDType: &field.TypeInfo{
				Type: field.TypeInt,
			},
		},
		Nodes: []*gen.Type{
			{
				Name: "Todo",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"RelayConnection": relayConnection,
					},
				},
			},
			{
				Name: "User",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"Skip": SkipAll,
					},
				},
			},
			{
				Name: "Group",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"RelayConnection": relayConnection,
					},
				},
			},
			{
				Name: "GroupWithSort",
				Fields: []*gen.Field{{
					Name: "Name",
					Type: &field.TypeInfo{
						Type: field.TypeString,
					},
					Annotations: map[string]interface{}{
						annotationName: map[string]interface{}{
							"OrderField": "NAME",
						},
					},
				}},
				Annotations: map[string]interface{}{
					annotationName: map[string]interface{}{
						"RelayConnection": relayConnection,
					},
				},
			},
		},
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
