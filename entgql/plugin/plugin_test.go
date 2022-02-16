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

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestEmpty(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.Equal(t, ``, e.print())
	require.NoError(t, err)
}

func TestInjectSourceEarlyEmpty(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	s := e.InjectSourceEarly()
	require.False(t, s.BuiltIn)
	require.Equal(t, `scalar Cursor
interface Node {
	id: ID!
}
enum OrderDirection {
	ASC
	DESC
}
type PageInfo {
	hasNextPage: Boolean!
	hasPreviousPage: Boolean!
	startCursor: Cursor
	endCursor: Cursor
}
scalar Time
`, s.Input)
}

func TestInjectSourceEarly(t *testing.T) {
	ann := entgql.Annotation{GqlScalarMappings: map[string]string{
		"Time": "Time",
	}}
	graph, err := entc.LoadGraph("../internal/todoplugin/ent/schema", &gen.Config{
		Annotations: map[string]interface{}{
			ann.Name(): ann,
		},
	})
	require.NoError(t, err)
	plugin, err := New(graph)
	require.NoError(t, err)
	s := plugin.InjectSourceEarly()
	require.Equal(t, expected, s.Input)
}

var expected = `type Category implements Node {
	id: ID!
	text: String!
	status: CategoryStatus!
	config: CategoryConfig!
	duration: Duration!
	count: Uint64!
	strings: [String!]
}
input CategoryOrder {
	direction: OrderDirection!
	field: CategoryOrderField!
}
enum CategoryOrderField {
	TEXT
	DURATION
}
enum CategoryStatus {
	ENABLED
	DISABLED
}
scalar Cursor
interface Node {
	id: ID!
}
enum OrderDirection {
	ASC
	DESC
}
type PageInfo {
	hasNextPage: Boolean!
	hasPreviousPage: Boolean!
	startCursor: Cursor
	endCursor: Cursor
}
enum Role {
	ADMIN
	USER
	UNKNOWN
}
enum Status {
	IN_PROGRESS
	COMPLETED
}
scalar Time
type Todo implements Node @someDirective {
	id: ID!
	createdAt: Time!
	status: Status! @someDirective(stringArg: "someString", boolArg: FALSE)
	priority: Int!
	text: String!
}
type TodoConnection {
	edges: [TodoEdge]
	pageInfo: PageInfo!
	totalCount: Int!
}
type TodoEdge {
	node: Todo
	cursor: Cursor
}
input TodoOrder {
	direction: OrderDirection!
	field: TodoOrderField!
}
enum TodoOrderField {
	CREATED_AT
	STATUS
	PRIORITY
	TEXT
}
type User implements Node {
	id: ID!
	username: String!
	age: Float!
	amount: Float!
	role: Role!
}
`
