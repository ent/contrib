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

package schemast

import (
	"go/ast"

	"entgo.io/ent"
	"entgo.io/ent/schema/index"
)

// Index converts a *index.Descriptor back into an *ast.CallExpr of the ent Index package that can be used
// to construct it.
func Index(desc *index.Descriptor) (*ast.CallExpr, error) {
	idx := newIndexCall(desc)
	if desc.Unique {
		idx.method("Unique")
	}
	if desc.StorageKey != "" {
		idx.method("StorageKey", strLit(desc.StorageKey))
	}
	if len(desc.Edges) > 0 {
		var edges []ast.Expr
		for _, edg := range desc.Edges {
			edges = append(edges, strLit(edg))
		}
		idx.method("Edges", edges...)
	}
	return idx.curr, nil
}

// AppendIndex adds an index to the returned values of the Indexes method of type typeName.
func (c *Context) AppendIndex(typeName string, idx ent.Index) error {
	newIdx, err := Index(idx.Descriptor())
	if err != nil {
		return err
	}
	return c.appendReturnItem(indexKind, typeName, newIdx)
}

func newIndexCall(desc *index.Descriptor) *builderCall {
	var fields []ast.Expr
	for _, fld := range desc.Fields {
		fields = append(fields, strLit(fld))
	}
	return &builderCall{
		curr: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("index"),
				Sel: ast.NewIdent("Fields"),
			},
			Args: fields,
		},
	}
}
