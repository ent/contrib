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

	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/require"
)

func TestRelayBuiltins(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	e.relayBuiltins()
	require.Equal(t, `"""
Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor
"""
scalar Cursor
"""
An object with an ID.
Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)
"""
interface Node {
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
`, e.print())
}

func TestRelayConnection(t *testing.T) {
	e, err := New(&gen.Graph{
		Config: &gen.Config{},
	})
	require.NoError(t, err)
	e.relayConnection(&gen.Type{
		Name: "Todo",
	})
	require.Equal(t, `"""
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
	The item at the end of the edge
	"""
	node: Todo
	"""
	A cursor for use in pagination
	"""
	cursor: Cursor!
}
`, e.print())
}
