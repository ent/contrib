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
	"entgo.io/ent/entc/gen"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	// RelayCursor is the name of the cursor type
	RelayCursor = "Cursor"
	// RelayNode is the name of the interface that all nodes implement
	RelayNode = "Node"
	// RelayPageInfo is the name of the PageInfo type
	RelayPageInfo = "PageInfo"
)

func relayBuiltinTypes() []*ast.Definition {
	return []*ast.Definition{
		{
			Name: RelayCursor,
			Kind: ast.Scalar,
			Description: `Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor`,
		},
		{
			Name: RelayNode,
			Kind: ast.Interface,
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
	}
}

func relayConnectionTypes(t *gen.Type) ([]*ast.Definition, error) {
	pagination, err := nodePaginationNames(t)
	if err != nil {
		return nil, err
	}

	return []*ast.Definition{
		{
			Name:        pagination.Edge,
			Kind:        ast.Object,
			Description: "An edge in a connection.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "node",
					Type:        ast.NamedType(pagination.Node, nil),
					Description: "The item at the end of the edge.",
				},
				{
					Name:        "cursor",
					Type:        ast.NonNullNamedType("Cursor", nil),
					Description: "A cursor for use in pagination.",
				},
			},
		},
		{
			Name:        pagination.Connection,
			Kind:        ast.Object,
			Description: "A connection to a list of items.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "edges",
					Type:        ast.ListType(ast.NamedType(pagination.Edge, nil), nil),
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
		},
	}, nil
}

func insertDefinitions(types map[string]*ast.Definition, defs ...*ast.Definition) {
	for _, d := range defs {
		types[d.Name] = d
	}
}
