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
	"fmt"
	"go/ast"
	"go/token"
	"strconv"

	"entgo.io/ent/schema/edge"
)

// Edge converts a *edge.Descriptor back into an *ast.CallExpr of the ent edge package that can be used
// to construct it.
func Edge(desc *edge.Descriptor) (*ast.CallExpr, error) {
	if len(desc.Annotations) > 0 {
		return nil, fmt.Errorf("schemast: unsupported feature: Annotations")
	}
	builder := newEdgeCall(desc)
	if desc.RefName != "" {
		builder.method("Ref", strLit(desc.RefName))
	}
	if desc.Required {
		builder.method("Required")
	}
	if desc.Unique {
		builder.method("Unique")
	}
	if desc.Field != "" {
		builder.method("Field", strLit(desc.Field))
	}
	if desc.StorageKey != nil {
		tbl := fnCall(selectorLit("edge", "Table"), strLit(desc.StorageKey.Table))
		col := fnCall(selectorLit("edge", "Column"), strLit(desc.StorageKey.Columns[0]))
		if len(desc.StorageKey.Columns) == 2 {
			to, from := strLit(desc.StorageKey.Columns[0]), strLit(desc.StorageKey.Columns[1])
			col = fnCall(selectorLit("edge", "Columns"), to, from)
		}
		builder.method("StorageKey", tbl, col)
	}
	if desc.Tag != "" {
		builder.method("StructTag", strLit(desc.Tag))
	}
	return builder.curr, nil
}

// AppendEdge adds an edge to the returned values of the Edges method of type typeName.
func (c *Context) AppendEdge(typeName string, desc *edge.Descriptor) error {
	stmt, err := c.edgesReturnStmt(typeName)
	if err != nil {
		return err
	}
	newEdge, err := Edge(desc)
	if err != nil {
		return err
	}
	returned := stmt.Results[0]
	switch r := returned.(type) {
	case *ast.Ident:
		if r.Name != "nil" {
			return fmt.Errorf("schemast: unexpected ident. expected nil got %s", r.Name)
		}
		stmt.Results = []ast.Expr{newEdgeSliceWith(newEdge)}
	case *ast.CompositeLit:
		r.Elts = append(r.Elts, newEdge)
	default:
		return fmt.Errorf("schemast: unexpected AST component type %T", r)
	}
	return nil
}

// RemoveEdge removes an edge from the returned values of the Edges method of type typeName.
func (c *Context) RemoveEdge(typeName string, edgeName string) error {
	stmt, err := c.edgesReturnStmt(typeName)
	if err != nil {
		return err
	}
	returned, ok := stmt.Results[0].(*ast.CompositeLit)
	if !ok {
		return fmt.Errorf("schemast: unexpected AST component type %T", stmt.Results[0])
	}
	for i, item := range returned.Elts {
		call, ok := item.(*ast.CallExpr)
		if !ok {
			return fmt.Errorf("schemast: expected return statement elements to be call expressions")
		}
		name, err := extractEdgeName(call)
		if err != nil {
			return err
		}
		if name == edgeName {
			returned.Elts = append(returned.Elts[:i], returned.Elts[i+1:]...)
			return nil
		}
	}
	return fmt.Errorf("schemast: could not find edge %q in type %q", edgeName, typeName)
}

func newEdgeCall(desc *edge.Descriptor) *builderCall {
	constructor := "To"
	if desc.Inverse {
		constructor = "From"
	}
	return &builderCall{
		curr: &ast.CallExpr{
			Fun: &ast.SelectorExpr{
				X:   ast.NewIdent("edge"),
				Sel: ast.NewIdent(constructor),
			},
			Args: []ast.Expr{
				strLit(desc.Name),
				selectorLit(desc.Type, "Type"),
			},
		},
	}
}

func (c *Context) edgesReturnStmt(typeName string) (*ast.ReturnStmt, error) {
	fd, err := c.lookupMethod(typeName, "Edges")
	if err != nil {
		return nil, err
	}
	if len(fd.Body.List) != 1 {
		return nil, fmt.Errorf("schmeast: Edges() func body must have a single element")
	}
	if _, ok := fd.Body.List[0].(*ast.ReturnStmt); !ok {
		return nil, fmt.Errorf("schmeast: Edges() func body must contain a return statement")
	}
	return fd.Body.List[0].(*ast.ReturnStmt), err
}

func newEdgeSliceWith(f *ast.CallExpr) *ast.CompositeLit {
	return &ast.CompositeLit{
		Type: &ast.ArrayType{
			Elt: &ast.SelectorExpr{
				X:   ast.NewIdent("ent"),
				Sel: ast.NewIdent("Edge"),
			},
		},
		Elts: []ast.Expr{
			f,
		},
	}
}

func extractEdgeName(fd *ast.CallExpr) (string, error) {
	sel, ok := fd.Fun.(*ast.SelectorExpr)
	if !ok {
		return "", fmt.Errorf("schemast: unexpected type %T", fd.Fun)
	}
	if inner, ok := sel.X.(*ast.CallExpr); ok {
		return extractEdgeName(inner)
	}
	if final, ok := sel.X.(*ast.Ident); ok && final.Name != "edge" {
		return "", fmt.Errorf(`schemast: expected edge AST to be of form edge.<To/From>("name")`)
	}
	if len(fd.Args) == 0 {
		return "", fmt.Errorf("schemast: expected edge constructor to have at least name arg")
	}
	name, ok := fd.Args[0].(*ast.BasicLit)
	if !ok && name.Kind == token.STRING {
		return "", fmt.Errorf("schemast: expected edge name to be a string literal")
	}
	return strconv.Unquote(name.Value)
}
