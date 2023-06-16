// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"go/parser"
	"go/token"

	"github.com/go-openapi/inflect"
)

// RemoveType removes the type definition as well as any method receivers or associated comment groups from the context.
func (c *Context) RemoveType(typeName string) error {
	_, found := c.newTypes[typeName]
	if found {
		delete(c.newTypes, typeName)
	}
	for _, file := range c.syntax() {
		toRemDecl := make(map[ast.Decl]struct{})
		toRemComments := make(map[*ast.CommentGroup]struct{})
		ast.Inspect(file, func(node ast.Node) bool {
			switch n := node.(type) {
			case *ast.FuncDecl:
				if len(n.Recv.List) > 0 {
					if id, ok := n.Recv.List[0].Type.(*ast.Ident); ok && id.Name == typeName {
						toRemDecl[n] = struct{}{}
						toRemComments[n.Doc] = struct{}{}
						found = true
					}
				}
			case *ast.GenDecl:
				if isTypeDeclFor(n, typeName) {
					toRemDecl[n] = struct{}{}
					toRemComments[n.Doc] = struct{}{}
				}
			}
			return true
		})
		var newComments []*ast.CommentGroup
		for _, comm := range file.Comments {
			if _, ok := toRemComments[comm]; !ok {
				newComments = append(newComments, comm)
			}
		}
		file.Comments = newComments
		var newDecls []ast.Decl
		for _, decl := range file.Decls {
			if _, ok := toRemDecl[decl]; !ok {
				newDecls = append(newDecls, decl)
			}
		}
		file.Decls = newDecls
	}
	if !found {
		return fmt.Errorf("schemast: type %q not found", typeName)
	}
	return nil
}

func (c *Context) AddType(typeName string) error {
	body := fmt.Sprintf(`package schema
import (
	"entgo.io/ent"
	"entgo.io/ent/schema"
)
type %s struct {
	ent.Schema
}
func (%s) Fields() []ent.Field {
	return nil
}
func (%s) Edges() []ent.Edge {
	return nil
}
func (%s) Annotations() []schema.Annotation {
	return nil
}
`, typeName, typeName, typeName, typeName)
	fn := inflect.Underscore(typeName) + ".go"
	f, err := parser.ParseFile(c.SchemaPackage.Fset, fn, body, 0)
	if err != nil {
		return err
	}
	c.newTypes[typeName] = f
	return nil
}

func isTypeDeclFor(n *ast.GenDecl, typeName string) bool {
	if n.Tok == token.TYPE && len(n.Specs) > 0 {
		if ts, ok := n.Specs[0].(*ast.TypeSpec); ok {
			return ts.Name.Name == typeName
		}
	}
	return false
}
