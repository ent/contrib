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
	"entgo.io/ent/schema"
)

// Mutator changes a Context.
type Mutator interface {
	Mutate(ctx *Context) error
}

// Mutate applies a sequence of mutations to a Context
func Mutate(ctx *Context, mutations ...Mutator) error {
	for _, mut := range mutations {
		if err := mut.Mutate(ctx); err != nil {
			return err
		}
	}
	return nil
}

// UpsertSchema implements Mutator. UpsertSchema will add to the Context the type named Name if not present and rewrite
// the type's Fields and Edges methods to return the desired fields and edges.
type UpsertSchema struct {
	Name        string
	Fields      []ent.Field
	Edges       []ent.Edge
	Indexes     []ent.Index
	Annotations []schema.Annotation
}

// Mutate applies the UpsertSchema mutation to the Context.
func (u *UpsertSchema) Mutate(ctx *Context) error {
	if !ctx.HasType(u.Name) {
		if err := ctx.AddType(u.Name); err != nil {
			return err
		}
	}
	if err := resetMethods(ctx, u.Name); err != nil {
		return err
	}
	for _, fld := range u.Fields {
		if err := ctx.AppendField(u.Name, fld.Descriptor()); err != nil {
			return err
		}
	}
	for _, edg := range u.Edges {
		if err := ctx.AppendEdge(u.Name, edg.Descriptor()); err != nil {
			return err
		}
	}
	for _, annot := range u.Annotations {
		if err := ctx.AppendTypeAnnotation(u.Name, annot); err != nil {
			return err
		}
	}
	for _, idx := range u.Indexes {
		if err := ctx.AppendIndex(u.Name, idx); err != nil {
			return err
		}
	}
	return nil
}

func resetMethods(ctx *Context, typeName string) error {
	for _, m := range []string{"Fields", "Edges", "Annotations", "Indexes"} {
		if _, ok := ctx.lookupMethod(typeName, m); !ok {
			continue
		}
		stmt, err := ctx.returnStmt(typeName, m)
		if err != nil {
			return err
		}
		stmt.Results = []ast.Expr{ast.NewIdent("nil")}
	}
	return nil
}

func (c *Context) appendReturnItem(k kind, typeName string, item ast.Expr) error {
	if _, ok := c.lookupMethod(typeName, k.methodName); !ok {
		if err := c.appendMethod(typeName, k.methodName, k.ifaceSelector); err != nil {
			return err
		}
	}
	stmt, err := c.returnStmt(typeName, k.methodName)
	if err != nil {
		return err
	}
	return appendToReturn(stmt, k.ifaceSelector, item)
}
