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
	"reflect"
	"testing"

	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

func Test_relayConnectionTypes(t *testing.T) {
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

			s := &ast.Schema{
				Types: map[string]*ast.Definition{},
			}
			insertDefinitions(s.Types, got...)
			gots := printSchema(s)
			if !reflect.DeepEqual(gots, tt.want) {
				t.Errorf("relayConnection() = %v, want %v", gots, tt.want)
			}
		})
	}
}

func Test_relayBuiltinTypes(t *testing.T) {
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
			got := relayBuiltinTypes()

			s := &ast.Schema{
				Types: map[string]*ast.Definition{},
			}
			insertDefinitions(s.Types, got...)
			gots := printSchema(s)
			if !reflect.DeepEqual(gots, tt.want) {
				t.Errorf("relayBuiltinTypes() = %v, want %v", gots, tt.want)
			}
		})
	}
}
