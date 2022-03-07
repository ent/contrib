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
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
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

	types, err := plugin.buildTypes()
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
}
"""
CategoryStatus is enum for the field status
"""
enum CategoryStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/category.Status") {
	ENABLED
	DISABLED
}
type MasterUser @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent.User") {
	id: ID
	username: String!
	age: Float!
	amount: Float!
	role: Role!
	nullableString: String
}
"""
Role is enum for the field role
"""
enum Role @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role") {
	ADMIN
	USER
	UNKNOWN
}
"""
Status is enum for the field status
"""
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
}
"""
VisibilityStatus is enum for the field visibility_status
"""
enum VisibilityStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus") {
	LISTING
	HIDDEN
}
`, printSchema(&ast.Schema{
		Types: types,
	}))
}

func TestEntGQL_buildTypes_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("./internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin, err := newSchemaGenerator(graph)

	require.NoError(t, err)
	types, err := plugin.buildTypes()
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
}
"""
A connection to a list of items.
"""
type CategoryConnection {
	"""
	A list of edges.
	"""
	edges: [CategoryEdge]
	"""
	Information to aid in pagination.
	"""
	pageInfo: PageInfo!
	totalCount: Int!
}
"""
An edge in a connection.
"""
type CategoryEdge {
	"""
	The item at the end of the edge.
	"""
	node: Category
	"""
	A cursor for use in pagination.
	"""
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
"""
CategoryStatus is enum for the field status
"""
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
"""
A connection to a list of items.
"""
type MasterUserConnection {
	"""
	A list of edges.
	"""
	edges: [MasterUserEdge]
	"""
	Information to aid in pagination.
	"""
	pageInfo: PageInfo!
	totalCount: Int!
}
"""
An edge in a connection.
"""
type MasterUserEdge {
	"""
	The item at the end of the edge.
	"""
	node: MasterUser
	"""
	A cursor for use in pagination.
	"""
	cursor: Cursor!
}
"""
Role is enum for the field role
"""
enum Role @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/role.Role") {
	ADMIN
	USER
	UNKNOWN
}
"""
Status is enum for the field status
"""
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
}
"""
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
"""
VisibilityStatus is enum for the field visibility_status
"""
enum VisibilityStatus @goModel(model: "entgo.io/contrib/entgql/internal/todoplugin/ent/todo.VisibilityStatus") {
	LISTING
	HIDDEN
}
`, printSchema(&ast.Schema{
		Types: types,
	}))
}

func disableRelayConnection(g *gen.Graph) {
	for _, n := range g.Nodes {
		if ant, ok := n.Annotations[annotationName]; ok {
			if m, ok := ant.(map[string]interface{}); ok {
				m["RelayConnection"] = false
			}
		}
	}
}
