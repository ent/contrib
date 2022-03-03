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

package plugin

import (
	"testing"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
	"github.com/vektah/gqlparser/v2/ast"
)

func TestEntGQL_buildTypes(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin, err := NewEntGQLPlugin(graph)

	require.NoError(t, err)
	types, err := plugin.buildTypes()
	require.NoError(t, err)

	require.Equal(t, `type Category {
	id: ID!
	text: String!
	uuidA: UUID
	status: CategoryStatus!
	config: CategoryConfig!
	duration: Duration!
	count: Uint64!
}
"""
CategoryStatus is enum for the field status
"""
enum CategoryStatus {
	ENABLED
	DISABLED
}
type MasterUser {
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
enum Role {
	ADMIN
	USER
	UNKNOWN
}
"""
Status is enum for the field status
"""
enum Status {
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
enum VisibilityStatus {
	LISTING
	HIDDEN
}
`, printSchema(&ast.Schema{
		Types: types,
	}))
}

func TestEntGQL_buildTypes_todoplugin_relay(t *testing.T) {
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{})
	require.NoError(t, err)
	plugin, err := NewEntGQLPlugin(graph, WithRelaySpecification(true))

	require.NoError(t, err)
	types, err := plugin.buildTypes()
	require.NoError(t, err)

	require.Equal(t, `type Category implements Node {
	id: ID!
	text: String!
	uuidA: UUID
	status: CategoryStatus!
	config: CategoryConfig!
	duration: Duration!
	count: Uint64!
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
"""
CategoryStatus is enum for the field status
"""
enum CategoryStatus {
	ENABLED
	DISABLED
}
type MasterUser implements Node {
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
enum Role {
	ADMIN
	USER
	UNKNOWN
}
"""
Status is enum for the field status
"""
enum Status {
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
"""
VisibilityStatus is enum for the field visibility_status
"""
enum VisibilityStatus {
	LISTING
	HIDDEN
}
`, printSchema(&ast.Schema{
		Types: types,
	}))
}
