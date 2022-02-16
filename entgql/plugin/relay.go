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
	"fmt"

	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	RelayCursor   = "Cursor"
	RelayNode     = "Node"
	RelayPageInfo = "PageInfo"
)

func (e *EntGQL) relayBuiltins() {
	e.insertDefinitions([]*ast.Definition{
		{
			Name: RelayCursor,
			Kind: ast.Scalar,
			Description: `Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor`,
		},
		{
			Kind: ast.Interface,
			Name: RelayNode,
			Description: `An object with an ID.
Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)`,
			Fields: []*ast.FieldDefinition{
				{
					Name:        "id",
					Type:        ast.NonNullNamedType("ID", nil),
					Description: "The id of the object.",
				},
			},
		},
		{
			Name: RelayPageInfo,
			Kind: ast.Object,
			Description: `Information about pagination in a connection.
https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo`,
			Fields: []*ast.FieldDefinition{
				{
					Name:        "hasNextPage",
					Type:        ast.NonNullNamedType("Boolean", nil),
					Description: "When paginating forwards, are there more items?",
				},
				{
					Name:        "hasPreviousPage",
					Type:        ast.NonNullNamedType("Boolean", nil),
					Description: "When paginating backwards, are there more items?",
				},
				{
					Name:        "startCursor",
					Type:        ast.NamedType("Cursor", nil),
					Description: "When paginating backwards, the cursor to continue.",
				},
				{
					Name:        "endCursor",
					Type:        ast.NamedType("Cursor", nil),
					Description: "When paginating forwards, the cursor to continue.",
				},
			},
		},
	})
}

func (e *EntGQL) relayConnection(t *gen.Type) {
	e.insertDefinition(&ast.Definition{
		Name:        fmt.Sprintf("%sEdge", t.Name),
		Kind:        ast.Object,
		Description: "An edge in a connection.",
		Fields: []*ast.FieldDefinition{
			{
				Name:        "node",
				Type:        ast.NamedType(t.Name, nil),
				Description: "The item at the end of the edge",
			},
			{
				Name:        "cursor",
				Type:        ast.NamedType("Cursor", nil),
				Description: "A cursor for use in pagination",
			},
		},
	})
	e.insertDefinition(&ast.Definition{
		Name:        fmt.Sprintf("%sConnection", t.Name),
		Kind:        ast.Object,
		Description: "A connection to a list of items.",
		Fields: []*ast.FieldDefinition{
			{
				Name:        "edges",
				Type:        ast.ListType(ast.NamedType(fmt.Sprintf("%sEdge", t.Name), nil), nil),
				Description: "A list of edges.",
			},
			{
				Name:        "pageInfo",
				Type:        ast.NonNullNamedType("PageInfo", nil),
				Description: "Information to aid in pagination.",
			},
			{
				Name: "totalCount",
				Type: ast.NonNullNamedType("Int", nil),
			},
		},
	})
}
