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
	plugin := newSchemaGenerator()
	plugin.genSchema = true
	plugin.relaySpec = false

	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(graph, schema)
	require.NoError(t, err)

	require.Equal(t, `type Category {
  id: ID!
  text: String!
  status: CategoryStatus!
  config: CategoryConfig
  duration: Duration
  count: Uint64
  strings: [String!]
  todos: [Todo!]
}
"""Ordering options for Category connections"""
input CategoryOrder {
  """The ordering direction."""
  direction: OrderDirection! = ASC
  """The field by which to order Categories."""
  field: CategoryOrderField!
}
"""Properties by which Category connections can be ordered."""
enum CategoryOrderField {
  TEXT
  DURATION
}
"""CategoryStatus is enum for the field status"""
enum CategoryStatus @goModel(model: "entgo.io/contrib/entgql/internal/todo/ent/category.Status") {
  ENABLED
  DISABLED
}
"""
CreateTodoInput is used for create Todo object.
Input was generated by ent.
"""
input CreateTodoInput {
  status: TodoStatus!
  priority: Int
  text: String!
  parentID: ID
  childIDs: [ID!]
  categoryID: ID
  secretID: ID
}
type Friendship {
  id: ID!
  createdAt: Time!
  user: User!
  friend: User!
}
type Group {
  id: ID!
  name: String!
  users: [User!]
}
type Query {
  groups: [Group!]!
  todos: [Todo!]!
  users: [User!]!
}
type Todo {
  id: ID!
  createdAt: Time!
  status: TodoStatus!
  priority: Int!
  text: String!
  categoryID: ID
  parent: Todo
  children: [Todo!]
  category: Category
}
"""Ordering options for Todo connections"""
input TodoOrder {
  """The ordering direction."""
  direction: OrderDirection! = ASC
  """The field by which to order Todos."""
  field: TodoOrderField!
}
"""Properties by which Todo connections can be ordered."""
enum TodoOrderField {
  CREATED_AT
  STATUS
  PRIORITY
  TEXT
}
"""TodoStatus is enum for the field status"""
enum TodoStatus @goModel(model: "entgo.io/contrib/entgql/internal/todo/ent/todo.Status") {
  IN_PROGRESS
  COMPLETED
}
type User {
  id: ID!
  name: String!
  groups: [Group!]
  friends: [User!]
  friendships: [Friendship!]
}
`, printSchema(schema))
}

func TestEntGQL_buildTypes_todoplugin_relay(t *testing.T) {
	s, err := gen.NewStorage("sql")
	require.NoError(t, err)

	graph, err := entc.LoadGraph("./internal/todo/ent/schema", &gen.Config{
		Storage: s,
	})
	require.NoError(t, err)
	plugin := newSchemaGenerator()
	plugin.genSchema = true
	plugin.genWhereInput = true
	schema := &ast.Schema{
		Types: make(map[string]*ast.Definition),
	}
	err = plugin.buildTypes(graph, schema)
	require.NoError(t, err)

	require.Equal(t, `type Category implements Node {
  id: ID!
  text: String!
  status: CategoryStatus!
  config: CategoryConfig
  duration: Duration
  count: Uint64
  strings: [String!]
  todos(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Ordering options for Todos returned from the connection."""
    orderBy: TodoOrder

    """Filtering options for Todos returned from the connection."""
    where: TodoWhereInput
  ): TodoConnection!
}
"""Ordering options for Category connections"""
input CategoryOrder {
  """The ordering direction."""
  direction: OrderDirection! = ASC
  """The field by which to order Categories."""
  field: CategoryOrderField!
}
"""Properties by which Category connections can be ordered."""
enum CategoryOrderField {
  TEXT
  DURATION
}
"""CategoryStatus is enum for the field status"""
enum CategoryStatus @goModel(model: "entgo.io/contrib/entgql/internal/todo/ent/category.Status") {
  ENABLED
  DISABLED
}
"""
CategoryWhereInput is used for filtering Category objects.
Input was generated by ent.
"""
input CategoryWhereInput {
  not: CategoryWhereInput
  and: [CategoryWhereInput!]
  or: [CategoryWhereInput!]
  """id field predicates"""
  id: ID
  idNEQ: ID
  idIn: [ID!]
  idNotIn: [ID!]
  idGT: ID
  idGTE: ID
  idLT: ID
  idLTE: ID
  """text field predicates"""
  text: String
  textNEQ: String
  textIn: [String!]
  textNotIn: [String!]
  textGT: String
  textGTE: String
  textLT: String
  textLTE: String
  textContains: String
  textHasPrefix: String
  textHasSuffix: String
  textEqualFold: String
  textContainsFold: String
  """status field predicates"""
  status: CategoryStatus
  statusNEQ: CategoryStatus
  statusIn: [CategoryStatus!]
  statusNotIn: [CategoryStatus!]
  """config field predicates"""
  config: CategoryConfig
  configNEQ: CategoryConfig
  configIn: [CategoryConfig!]
  configNotIn: [CategoryConfig!]
  configGT: CategoryConfig
  configGTE: CategoryConfig
  configLT: CategoryConfig
  configLTE: CategoryConfig
  configIsNil: Boolean
  configNotNil: Boolean
  """duration field predicates"""
  duration: Duration
  durationNEQ: Duration
  durationIn: [Duration!]
  durationNotIn: [Duration!]
  durationGT: Duration
  durationGTE: Duration
  durationLT: Duration
  durationLTE: Duration
  durationIsNil: Boolean
  durationNotNil: Boolean
  """count field predicates"""
  count: Uint64
  countNEQ: Uint64
  countIn: [Uint64!]
  countNotIn: [Uint64!]
  countGT: Uint64
  countGTE: Uint64
  countLT: Uint64
  countLTE: Uint64
  countIsNil: Boolean
  countNotNil: Boolean
  """todos edge predicates"""
  hasTodos: Boolean
  hasTodosWith: [TodoWhereInput!]
}
"""
CreateTodoInput is used for create Todo object.
Input was generated by ent.
"""
input CreateTodoInput {
  status: TodoStatus!
  priority: Int
  text: String!
  parentID: ID
  childIDs: [ID!]
  categoryID: ID
  secretID: ID
}
type Friendship implements Node {
  id: ID!
  createdAt: Time!
  user: User!
  friend: User!
}
"""
FriendshipWhereInput is used for filtering Friendship objects.
Input was generated by ent.
"""
input FriendshipWhereInput {
  not: FriendshipWhereInput
  and: [FriendshipWhereInput!]
  or: [FriendshipWhereInput!]
  """id field predicates"""
  id: ID
  idNEQ: ID
  idIn: [ID!]
  idNotIn: [ID!]
  idGT: ID
  idGTE: ID
  idLT: ID
  idLTE: ID
  """created_at field predicates"""
  createdAt: Time
  createdAtNEQ: Time
  createdAtIn: [Time!]
  createdAtNotIn: [Time!]
  createdAtGT: Time
  createdAtGTE: Time
  createdAtLT: Time
  createdAtLTE: Time
}
type Group implements Node {
  id: ID!
  name: String!
  users(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Filtering options for Users returned from the connection."""
    where: UserWhereInput
  ): UserConnection!
}
"""A connection to a list of items."""
type GroupConnection {
  """A list of edges."""
  edges: [GroupEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  """Identifies the total count of items in the connection."""
  totalCount: Int!
}
"""An edge in a connection."""
type GroupEdge {
  """The item at the end of the edge."""
  node: Group
  """A cursor for use in pagination."""
  cursor: Cursor!
}
"""
GroupWhereInput is used for filtering Group objects.
Input was generated by ent.
"""
input GroupWhereInput {
  not: GroupWhereInput
  and: [GroupWhereInput!]
  or: [GroupWhereInput!]
  """id field predicates"""
  id: ID
  idNEQ: ID
  idIn: [ID!]
  idNotIn: [ID!]
  idGT: ID
  idGTE: ID
  idLT: ID
  idLTE: ID
  """name field predicates"""
  name: String
  nameNEQ: String
  nameIn: [String!]
  nameNotIn: [String!]
  nameGT: String
  nameGTE: String
  nameLT: String
  nameLTE: String
  nameContains: String
  nameHasPrefix: String
  nameHasSuffix: String
  nameEqualFold: String
  nameContainsFold: String
  """users edge predicates"""
  hasUsers: Boolean
  hasUsersWith: [UserWhereInput!]
}
type Query {
  """Fetches an object given its ID."""
  node(
    """ID of the object."""
    id: ID!
  ): Node
  """Lookup nodes by a list of IDs."""
  nodes(
    """The list of node IDs."""
    ids: [ID!]!
  ): [Node]!
  groups(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Filtering options for Groups returned from the connection."""
    where: GroupWhereInput
  ): GroupConnection!
  todos(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Ordering options for Todos returned from the connection."""
    orderBy: TodoOrder

    """Filtering options for Todos returned from the connection."""
    where: TodoWhereInput
  ): TodoConnection!
  users(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Filtering options for Users returned from the connection."""
    where: UserWhereInput
  ): UserConnection!
}
type Todo implements Node {
  id: ID!
  createdAt: Time!
  status: TodoStatus!
  priority: Int!
  text: String!
  categoryID: ID
  parent: Todo
  children(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Ordering options for Todos returned from the connection."""
    orderBy: TodoOrder

    """Filtering options for Todos returned from the connection."""
    where: TodoWhereInput
  ): TodoConnection!
  category: Category
}
"""A connection to a list of items."""
type TodoConnection {
  """A list of edges."""
  edges: [TodoEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  """Identifies the total count of items in the connection."""
  totalCount: Int!
}
"""An edge in a connection."""
type TodoEdge {
  """The item at the end of the edge."""
  node: Todo
  """A cursor for use in pagination."""
  cursor: Cursor!
}
"""Ordering options for Todo connections"""
input TodoOrder {
  """The ordering direction."""
  direction: OrderDirection! = ASC
  """The field by which to order Todos."""
  field: TodoOrderField!
}
"""Properties by which Todo connections can be ordered."""
enum TodoOrderField {
  CREATED_AT
  STATUS
  PRIORITY
  TEXT
}
"""TodoStatus is enum for the field status"""
enum TodoStatus @goModel(model: "entgo.io/contrib/entgql/internal/todo/ent/todo.Status") {
  IN_PROGRESS
  COMPLETED
}
"""
TodoWhereInput is used for filtering Todo objects.
Input was generated by ent.
"""
input TodoWhereInput {
  not: TodoWhereInput
  and: [TodoWhereInput!]
  or: [TodoWhereInput!]
  """id field predicates"""
  id: ID
  idNEQ: ID
  idIn: [ID!]
  idNotIn: [ID!]
  idGT: ID
  idGTE: ID
  idLT: ID
  idLTE: ID
  """created_at field predicates"""
  createdAt: Time
  createdAtNEQ: Time
  createdAtIn: [Time!]
  createdAtNotIn: [Time!]
  createdAtGT: Time
  createdAtGTE: Time
  createdAtLT: Time
  createdAtLTE: Time
  """status field predicates"""
  status: TodoStatus
  statusNEQ: TodoStatus
  statusIn: [TodoStatus!]
  statusNotIn: [TodoStatus!]
  """priority field predicates"""
  priority: Int
  priorityNEQ: Int
  priorityIn: [Int!]
  priorityNotIn: [Int!]
  priorityGT: Int
  priorityGTE: Int
  priorityLT: Int
  priorityLTE: Int
  """text field predicates"""
  text: String
  textNEQ: String
  textIn: [String!]
  textNotIn: [String!]
  textGT: String
  textGTE: String
  textLT: String
  textLTE: String
  textContains: String
  textHasPrefix: String
  textHasSuffix: String
  textEqualFold: String
  textContainsFold: String
  """category_id field predicates"""
  categoryID: ID
  categoryIDNEQ: ID
  categoryIDIn: [ID!]
  categoryIDNotIn: [ID!]
  categoryIDIsNil: Boolean
  categoryIDNotNil: Boolean
  """parent edge predicates"""
  hasParent: Boolean
  hasParentWith: [TodoWhereInput!]
  """children edge predicates"""
  hasChildren: Boolean
  hasChildrenWith: [TodoWhereInput!]
  """category edge predicates"""
  hasCategory: Boolean
  hasCategoryWith: [CategoryWhereInput!]
}
type User implements Node {
  id: ID!
  name: String!
  groups(
    """Returns the elements in the list that come after the specified cursor."""
    after: Cursor

    """Returns the first _n_ elements from the list."""
    first: Int

    """Returns the elements in the list that come before the specified cursor."""
    before: Cursor

    """Returns the last _n_ elements from the list."""
    last: Int

    """Filtering options for Groups returned from the connection."""
    where: GroupWhereInput
  ): GroupConnection!
  friends: [User!]
  friendships: [Friendship!]
}
"""A connection to a list of items."""
type UserConnection {
  """A list of edges."""
  edges: [UserEdge]
  """Information to aid in pagination."""
  pageInfo: PageInfo!
  """Identifies the total count of items in the connection."""
  totalCount: Int!
}
"""An edge in a connection."""
type UserEdge {
  """The item at the end of the edge."""
  node: User
  """A cursor for use in pagination."""
  cursor: Cursor!
}
"""
UserWhereInput is used for filtering User objects.
Input was generated by ent.
"""
input UserWhereInput {
  not: UserWhereInput
  and: [UserWhereInput!]
  or: [UserWhereInput!]
  """id field predicates"""
  id: ID
  idNEQ: ID
  idIn: [ID!]
  idNotIn: [ID!]
  idGT: ID
  idGTE: ID
  idLT: ID
  idLTE: ID
  """name field predicates"""
  name: String
  nameNEQ: String
  nameIn: [String!]
  nameNotIn: [String!]
  nameGT: String
  nameGTE: String
  nameLT: String
  nameLTE: String
  nameContains: String
  nameHasPrefix: String
  nameHasSuffix: String
  nameEqualFold: String
  nameContainsFold: String
  """groups edge predicates"""
  hasGroups: Boolean
  hasGroupsWith: [GroupWhereInput!]
  """friends edge predicates"""
  hasFriends: Boolean
  hasFriendsWith: [UserWhereInput!]
  """friendships edge predicates"""
  hasFriendships: Boolean
  hasFriendshipsWith: [FriendshipWhereInput!]
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
  """Identifies the total count of items in the connection."""
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
  """Identifies the total count of items in the connection."""
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
interface Node @goModel(model: "todo/ent.Noder") {
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
