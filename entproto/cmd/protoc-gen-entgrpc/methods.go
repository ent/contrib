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
	camel    = gen.Funcs["camel"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
)

func (g *serviceGenerator) generateGetMethod(methodName string) error {
	idField := g.fieldMap.ID()
	convert, err := g.newConverter(idField)
	if err != nil {
		return err
	}
	if fieldNeedsValidator(idField) {
		g.generateIDFieldValidator(idField)
	}
	vars := g.withGlobals(tmplValues{
		"id":         g.renderToEnt(convert, fmt.Sprintf("req.Get%s()", idField.PbStructField())),
		"methodName": methodName,
	})
	g.Tmpl(`var (
		err error
		get *ent.%(typeName)
	)
	switch req.GetView() {
		case %(methodName)_VIEW_UNSPECIFIED, %(methodName)_BASIC:
			get, err = svc.client.%(typeName).Get(ctx, %(id))
		case %(methodName)_WITH_EDGE_IDS:
			get, err = svc.client.%(typeName).Query().
`, vars)
	for _, edg := range g.fieldMap.Edges() {
		g.Tmpl(`With%(edgeName)(func(query *ent.%(otherType)Query) {
	query.Select("id")
}).`, g.withGlobals(tmplValues{
			"edgeName":  edg.PbStructField(),
			"otherType": edg.EntEdge.Type.Name,
		}))
	}
	g.Tmpl(`
			First(ctx)
		default:
			return nil, %(statusErrf)(%(invalidArgument), "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProto%(typeName)(get), nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, vars)
	g.Tmpl(``, vars)
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
			"checkIDFlag": strconv.FormatBool(op == "update"),
		}))
	}
	switch op {
	case "create":
		g.Tmpl("m := svc.client.%(typeName).Create()", g.withGlobals())
	case "update":
		idField := g.fieldMap.ID()
		convert, err := g.newConverter(idField)
		if err != nil {
			return err
		}
		g.Tmpl(`m := svc.client.%(typeName).UpdateOneID(%(id))`, g.withGlobals(tmplValues{
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
		if fld.EntField.Optional {
			g.Tmpl("if %(reqVar).Get%(structField)() != nil {", tmplValues{
				"reqVar":      reqVar,
				"structField": fld.PbStructField(),
			})
		}
		g.Tmpl("m.Set%(entField)(%(converted))", tmplValues{
			"entField":  entField,
			"converted": g.renderToEnt(convert, fmt.Sprintf("%s.Get%s()", reqVar, fld.PbStructField())),
		})
		if fld.EntField.Optional {
			g.P("}")
		}
	}
	for _, edg := range g.fieldMap.Edges() {
		convert, err := g.newConverter(edg)
		if err != nil {
			return err
		}
		if edg.EntEdge.Unique {
			g.Tmpl("m.Set%(edgeName)ID(%(converted))", tmplValues{
				"edgeName":  edg.EntEdge.StructField(),
				"converted": g.renderToEnt(convert, fmt.Sprintf("%s.Get%s().Get%s()", reqVar, edg.PbStructField(), edg.EdgeIDPbStructField())),
			})
		} else {
			g.Tmpl(`for _, item := range %(reqVar).Get%(edgeField)() {
	m.Add%(edgeName)IDs(%(converted))
}
`, g.withGlobals(tmplValues{
				"reqVar":    reqVar,
				"edgeField": edg.PbStructField(),
				"edgeName":  singular(edg.EntEdge.StructField()),
				"converted": g.renderToEnt(convert, fmt.Sprintf("item.Get%s()", edg.EdgeIDPbStructField())),
			}))
		}
	}
	g.P("res, err := m.Save(ctx)")

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
