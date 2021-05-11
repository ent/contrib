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
	"entgo.io/ent/entc/gen"
	"fmt"
	"github.com/vektah/gqlparser/v2/ast"
)

var (
	RelayCursor   = "Cursor"
	RelayNode     = "Node"
	RelayPageInfo = "PageInfo"
)

func (e *Entgqlgen) relayBuiltins() {
	e.insertDefinitions([]*ast.Definition{
		{
			Name: RelayCursor,
			Kind: ast.Scalar,
		},
		{
			Kind: ast.Interface,
			Name: RelayNode,
			Fields: []*ast.FieldDefinition{
				{
					Name: "id",
					Type: ast.NonNullNamedType("ID", nil),
				},
			},
		},
		{
			Name: RelayPageInfo,
			Kind: ast.Object,
			Fields: []*ast.FieldDefinition{
				{
					Name: "hasNextPage",
					Type: ast.NonNullNamedType("Boolean", nil),
				},
				{
					Name: "hasPreviousPage",
					Type: ast.NonNullNamedType("Boolean", nil),
				},
				{
					Name: "startCursor",
					Type: ast.NamedType("Cursor", nil),
				},
				{
					Name: "endCursor",
					Type: ast.NamedType("Cursor", nil),
				},
			},
		},
	})
}

func (e *Entgqlgen) relayConnection(t *gen.Type) {
	e.insertDefinition(&ast.Definition{
		Name: fmt.Sprintf("%sEdge", t.Name),
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "node",
				Type: ast.NamedType(t.Name, nil),
			},
			{
				Name: "cursor",
				Type: ast.NamedType("Cursor", nil),
			},
		},
	})
	e.insertDefinition(&ast.Definition{
		Name: fmt.Sprintf("%sConnection", t.Name),
		Kind: ast.Object,
		Fields: []*ast.FieldDefinition{
			{
				Name: "edges",
				Type: ast.ListType(ast.NamedType(fmt.Sprintf("%sEdge", t.Name), nil), nil),
			},
			{
				Name: "pageInfo",
				Type: ast.NonNullNamedType("PageInfo", nil),
			},
			{
				Name: "totalCount",
				Type: ast.NonNullNamedType("Int", nil),
			},
		},
	})
}
