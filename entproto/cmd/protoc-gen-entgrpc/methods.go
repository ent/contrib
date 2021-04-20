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
	"fmt"
	"strconv"

	"entgo.io/ent/entc/gen"
	"google.golang.org/protobuf/compiler/protogen"
)

var (
	camel = gen.Funcs["camel"].(func(string) string)
)

func (g *serviceGenerator) generateGetMethod() error {
	idField := g.fieldMap.ID()
	convert, err := g.newConverter(idField)
	if err != nil {
		return err
	}
	if fieldNeedsValidator(idField) {
		g.generateIDFieldValidator(idField)
	}
	g.Tmpl(`get, err := svc.client.%(typeName).Get(ctx, %(id))
	switch {
	case err == nil:
		return toProto%(typeName)(get), nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, g.withGlobals(tmplValues{
		"id": g.renderToEnt(convert, fmt.Sprintf("req.Get%s()", idField.PbStructField())),
	}))
	return nil
}

func (g *serviceGenerator) generateDeleteMethod() error {
	idField := g.fieldMap.ID()
	convert, err := g.newConverter(idField)
	if err != nil {
		return err
	}
	if fieldNeedsValidator(idField) {
		g.generateIDFieldValidator(idField)
	}
	g.Tmpl(`err := svc.client.%(typeName).DeleteOneID(%(id)).Exec(ctx)
	switch {
	case err == nil:
		return &%(empty){}, nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, g.withGlobals(tmplValues{
		"id":    g.renderToEnt(convert, fmt.Sprintf("req.Get%s()", idField.PbStructField())),
		"empty": protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb").Ident("Empty"),
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
	if typeNeedsValidator(g.fieldMap) {
		g.Tmpl(`if err := validate%(typeName)(%(reqVar), %(checkIDFlag)); err != nil {
			return nil, %(statusErrf)(%(invalidArgument), "invalid argument: %s", err)
		}`, g.withGlobals(tmplValues{
			"reqVar":      reqVar,
			"checkIDFlag": strconv.FormatBool(op == "create"),
		}))
	}
	switch op {
	case "create":
		g.Tmpl("res, err := svc.client.%(typeName).Create().", g.withGlobals())
	case "update":
		idField := g.fieldMap.ID()
		convert, err := g.newConverter(idField)
		if err != nil {
			return err
		}
		g.Tmpl(`res, err := svc.client.%(typeName).UpdateOneID(%(id)).`, g.withGlobals(tmplValues{
			"id":     g.renderToEnt(convert, fmt.Sprintf("%s.Get%s()", reqVar, idField.PbStructField())),
			"reqVar": reqVar,
		}))
	}

	for _, fld := range g.fieldMap.Fields() {
		if fld.IsIDField || (op == "update" && fld.EntField.Immutable) {
			continue
		}
		convert, err := g.newConverter(fld)
		if err != nil {
			return err
		}
		entField := fld.EntField.StructField()
		g.Tmpl("Set%(entField)(%(converted)).", tmplValues{
			"entField":  entField,
			"converted": g.renderToEnt(convert, fmt.Sprintf("%s.Get%s()", reqVar, fld.PbStructField())),
		})
	}
	for _, edg := range g.fieldMap.Edges() {
		convert, err := g.newConverter(edg)
		if err != nil {
			return err
		}
		if edg.EntEdge.Unique {
			g.Tmpl("Set%(edgeName)ID(%(converted)).", tmplValues{
				"edgeName":  edg.EntEdge.StructField(),
				"converted": g.renderToEnt(convert, fmt.Sprintf("%s.Get%s().Get%s()", reqVar, edg.PbStructField(), edg.EdgeIDPbStructField())),
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
		"fmtErr":            protogen.GoImportPath("fmt").Ident("Errorf"),
	}
	for _, additional := range additionals {
		for k, v := range additional {
			m[k] = v
		}
	}
	return m
}
