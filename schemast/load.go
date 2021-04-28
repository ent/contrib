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

	"golang.org/x/tools/go/packages"
)

// Context represents an ent schema directory, parsed and loaded as ASTs, such that schema type declarations
// can be analyzed an manipulated by different programs.
type Context struct {
	SchemaPackage *packages.Package
	newTypes      map[string]*ast.File
}

// HasType reports whether typeName is already defined in the Context.
func (c *Context) HasType(typeName string) bool {
	for _, file := range c.syntax() {
		var found bool
		ast.Inspect(file, func(node ast.Node) bool {
			if decl, ok := node.(*ast.GenDecl); ok {
				if isTypeDeclFor(decl, typeName) {
					found = true
					return false
				}
			}
			return true
		})
		if found {
			return true
		}
	}
	return false
}

// lookupMethod will search the schemast.Context for the AST representing the function declaration of the requested
// methodName for type typeName.
func (c *Context) lookupMethod(typeName string, methodName string) (*ast.FuncDecl, error) {
	var found *ast.FuncDecl
	for _, file := range c.syntax() {
		ast.Inspect(file, func(node ast.Node) bool {
			if fd, ok := node.(*ast.FuncDecl); ok {
				if fd.Name.Name != methodName {
					return true
				}
				if len(fd.Recv.List) != 1 {
					return true
				}
				if id, ok := fd.Recv.List[0].Type.(*ast.Ident); ok && id.Name == typeName {
					found = fd
					return false
				}
			}
			return true
		})
		if found != nil {
			return found, nil
		}
	}
	return nil, &notFoundErr{
		typeName:   typeName,
		methodName: methodName,
	}
}

// Load loads a *schemast.Context from a path.
func Load(path string) (*Context, error) {
	pkgs, err := packages.Load(&packages.Config{
		Mode: packages.NeedName | packages.NeedTypes | packages.NeedTypesInfo | packages.NeedSyntax,
	}, path)
	if err != nil {
		return nil, fmt.Errorf("loading package: %w", err)
	}
	if len(pkgs) < 1 {
		return nil, fmt.Errorf("missing package information for: %s", path)
	}
	return &Context{
		SchemaPackage: pkgs[0],
		newTypes:      make(map[string]*ast.File),
	}, nil
}

type notFoundErr struct {
	typeName, methodName string
}

func (e *notFoundErr) Error() string {
	return fmt.Sprintf("schemast: could not find method %q for type %q", e.methodName, e.typeName)
}

func (c *Context) syntax() []*ast.File {
	var out []*ast.File
	out = append(out, c.SchemaPackage.Syntax...)
	for _, f := range c.newTypes {
		out = append(out, f)
	}
	return out
}
