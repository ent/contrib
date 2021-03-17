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

package main

import (
	"entgo.io/ent/entc/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	camel = gen.Funcs["camel"].(func(string) string)
)

func (g *serviceGenerator) generateGetMethod() error {
	idField := g.fieldMap.ID()
	cast, err := g.castToEntFunc(idField)
	if err != nil {
		return err
	}
	g.Tmpl(`get, err := svc.client.%(typeName).Get(ctx, %(cast)(req.Get%(pbIdField)()))
	switch {
	case err == nil:
		return toProto%(typeName)(get), nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, g.withGlobals(tmplValues{
		"cast":      cast,
		"pbIdField": idField.PbStructField(),
	}))
	return nil
}

func (g *serviceGenerator) generateDeleteMethod() error {
	idField := g.fieldMap.ID()
	cast, err := g.castToEntFunc(idField)
	if err != nil {
		return err
	}
	g.Tmpl(`err := svc.client.%(typeName).DeleteOneID(%(cast)(req.Get%(pbIdField)())).Exec(ctx)
	switch {
	case err == nil:
		return &%(empty){}, nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, g.withGlobals(tmplValues{
		"cast":      cast,
		"pbIdField": idField.PbStructField(),
		"empty":     protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb").Ident("Empty"),
	}))
	return nil
}

func (g *serviceGenerator) generateUpdateMethod() error {
	return g.generateMutationMethod("update")
}

func (g *serviceGenerator) generateCreateMethod() error {
	return g.generateMutationMethod("create")
}

func (g *serviceGenerator) generateMutationMethod(op string) error {
	reqVar := camel(g.typeName)
	g.Tmpl("%(reqVar) := req.Get%(typeName)()", g.withGlobals(tmplValues{
		"reqVar": reqVar,
	}))
	switch op {
	case "create":
		g.Tmpl("res, err := svc.client.%(typeName).Create().", g.withGlobals())
	case "update":
		idField := g.fieldMap.ID()
		cast, err := g.castToEntFunc(idField)
		if err != nil {
			return err
		}
		g.Tmpl(`res, err := svc.client.%(typeName).UpdateOneID(%(cast)(%(reqVar).Get%(pbIdField)())).`, g.withGlobals(tmplValues{
			"pbIdField": idField.PbStructField(),
			"cast":      cast,
			"reqVar":    reqVar,
		}))
	}

	for _, fld := range g.fieldMap.Fields() {
		if fld.IsIDField {
			continue
		}
		castFn, err := g.castToEntFunc(fld)
		if err != nil {
			return err
		}
		entField := fld.EntField.StructField()
		g.Tmpl("Set%(entField)( %(castFn)(%(reqVar).Get%(pbField)())).", tmplValues{
			"reqVar":   reqVar,
			"entField": entField,
			"castFn":   castFn,
			"pbField":  fld.PbStructField(),
		})
	}
	for _, edg := range g.fieldMap.Edges() {
		if edg.EntEdge.Unique {
			cast, err := g.castToEntFunc(edg)
			if err != nil {
				return err
			}
			g.Tmpl("Set%(edgeName)ID(%(cast)(%(reqVar).Get%(pbField)().Get%(edgeIdField)())).", tmplValues{
				"edgeName":    edg.EntEdge.StructField(),
				"pbField":     edg.PbStructField(),
				"reqVar":      reqVar,
				"edgeIdField": edg.EdgeIDPbStructField(),
				"cast":        cast,
			})
		}
	}
	g.P("Save(ctx)")

	g.Tmpl(`
	switch {
	case err == nil:
		return toProto%(typeName)(res), nil
	case %(uniqConstraintErr)(err):
		return nil, %(statusErrf)(%(alreadyExists), "already exists: %s", err)
	case %(constraintErr)(err):
		return nil, %(statusErrf)(%(invalidArgument), "invalid argument: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal: %s", err)
	}`, g.withGlobals())

	return nil
}

func (g *serviceGenerator) withGlobals(additionals ...tmplValues) tmplValues {
	m := tmplValues{
		"uniqConstraintErr": protogen.GoImportPath("entgo.io/ent/dialect/sql/sqlgraph").Ident("IsUniqueConstraintError"),
		"constraintErr":     g.entPackage.Ident("IsConstraintError"),
		"isNotFound":        g.entPackage.Ident("IsNotFound"),
		"statusErrf":        status.Ident("Errorf"),
		"alreadyExists":     codes.Ident("AlreadyExists"),
		"invalidArgument":   codes.Ident("InvalidArgument"),
		"notFound":          codes.Ident("NotFound"),
		"internal":          codes.Ident("Internal"),
		"typeName":          g.typeName,
	}
	for _, additional := range additionals {
		for k, v := range additional {
			m[k] = v
		}
	}
	return m
}
