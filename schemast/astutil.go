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
	"sort"
	"strconv"

	"go.uber.org/multierr"
)

type builderCall struct {
	curr *ast.CallExpr
}

func (f *builderCall) method(name string, args ...ast.Expr) {
	f.curr = &ast.CallExpr{
		Fun: &ast.SelectorExpr{
			X:   f.curr,
			Sel: ast.NewIdent(name),
		},
		Args: args,
	}
}

func (f *builderCall) annotate(annots ...ast.Expr) {
	if len(annots) > 0 {
		f.method("Annotations", annots...)
	}
}

func combineUnsupported(err error, feature string) error {
	return multierr.Combine(err, fmt.Errorf("schemast: unsupported feature %s", feature))
}

func strMapLit(m map[string]string) ast.Expr {
	c := &ast.CompositeLit{
		Type: &ast.MapType{
			Key:   ast.NewIdent("string"),
			Value: ast.NewIdent("string"),
		},
	}
	keys := make([]string, 0, len(m))
	for k := range m {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		c.Elts = append(c.Elts, &ast.KeyValueExpr{
			Key:   strLit(k),
			Value: strLit(m[k]),
		})
	}
	return c
}

func strLit(lit string) ast.Expr {
	return &ast.BasicLit{
		Kind:  token.STRING,
		Value: strconv.Quote(lit),
	}
}

func structAttr(name string, val ast.Expr) ast.Expr {
	return &ast.KeyValueExpr{
		Key: &ast.BasicLit{
			Kind:  token.STRING,
			Value: name,
		},
		Value: val,
	}
}

func intLit(lit int) ast.Expr {
	return &ast.BasicLit{
		Kind:  token.INT,
		Value: strconv.Itoa(lit),
	}
}

func selectorLit(x, sel string) *ast.SelectorExpr {
	return &ast.SelectorExpr{
		X:   ast.NewIdent(x),
		Sel: ast.NewIdent(sel),
	}
}

func fnCall(sel *ast.SelectorExpr, args ...ast.Expr) *ast.CallExpr {
	return &ast.CallExpr{
		Fun:  sel,
		Args: args,
	}
}
