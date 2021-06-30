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

	"entgo.io/ent"
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

	vars := g.withGlobals(tmplValues{
		"varName":     idField.EntField.Name,
		"extract":     g.renderToEnt(convert, idField.EntField.Name, fmt.Sprintf("req.Get%s()", idField.PbStructField())),
		"methodName":  methodName,
		"idPredicate": protogen.GoImportPath(string(g.entPackage) + "/" + g.entType.Package()).Ident("ID"),
	})
	g.Tmpl(`var (
		err error
		get *ent.%(typeName)
	)
	%(extract)
	switch req.GetView() {
		case %(methodName)_VIEW_UNSPECIFIED, %(methodName)_BASIC:
			get, err = svc.client.%(typeName).Get(ctx, %(varName))
		case %(methodName)_WITH_EDGE_IDS:
			get, err = svc.client.%(typeName).Query().
				Where(%(idPredicate)(%(varName))).
`, vars)
	for _, edg := range g.fieldMap.Edges() {
		et := edg.EntEdge.Type
		g.Tmpl(`With%(edgeName)(func(query *ent.%(otherType)Query) {
	query.Select(%(fieldID))
}).`, g.withGlobals(tmplValues{
			"edgeName":  edg.PbStructField(),
			"otherType": et.Name,
			"fieldID":   protogen.GoImportPath(string(g.entPackage) + "/" + et.Package()).Ident(et.ID.Constant()),
		}))
	}
	g.Tmpl(`
			Only(ctx)
		default:
			return nil, %(statusErrf)(%(invalidArgument), "invalid argument: unknown view")
	}
	switch {
	case err == nil:
		return toProto%(typeName)(get)
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, vars)

	return nil
}

func (g *serviceGenerator) generateDeleteMethod() error {
	idField := g.fieldMap.ID()
	convert, err := g.newConverter(idField)
	if err != nil {
		return err
	}

	g.Tmpl(`var err error
	%(extract)
	err = svc.client.%(typeName).DeleteOneID(%(varName)).Exec(ctx)
	switch {
	case err == nil:
		return &%(empty){}, nil
	case %(isNotFound)(err):
		return nil, %(statusErrf)(%(notFound), "not found: %s", err)
	default:
		return nil, %(statusErrf)(%(internal), "internal error: %s", err)
	}`, g.withGlobals(tmplValues{
		"varName": idField.EntField.Name,
		"extract": g.renderToEnt(convert, idField.EntField.Name, fmt.Sprintf("req.Get%s()", idField.PbStructField())),
		"empty":   protogen.GoImportPath("google.golang.org/protobuf/types/known/emptypb").Ident("Empty"),
	}))
	return nil
}

func (g *serviceGenerator) generateUpdateMethod() error {
	return g.generateMutationMethod(ent.OpUpdateOne)
}

func (g *serviceGenerator) generateCreateMethod() error {
	return g.generateMutationMethod(ent.OpCreate)
}

func (g *serviceGenerator) generateMutationMethod(op ent.Op) error {
	reqVar := camel(g.entType.Name)
	g.Tmpl("%(reqVar) := req.Get%(typeName)()", g.withGlobals(tmplValues{
		"reqVar": reqVar,
	}))

	switch op {
	case ent.OpCreate:
		g.Tmpl("m := svc.client.%(typeName).Create()", g.withGlobals())
	case ent.OpUpdateOne:
		idField := g.fieldMap.ID()
		varName := camel(reqVar + "_" + idField.EntField.Name)
		convert, err := g.newConverter(idField)
		if err != nil {
			return err
		}
		g.Tmpl(`%(extract)
		m := svc.client.%(typeName).UpdateOneID(%(varName))`, g.withGlobals(tmplValues{
			"varName": varName,
			"extract": g.renderToEnt(convert, varName, fmt.Sprintf("%s.Get%s()", reqVar, idField.PbStructField())),
			"reqVar":  reqVar,
		}))
	}

	for _, fld := range g.fieldMap.Fields() {
		if fld.IsIDField || (op.Is(ent.OpUpdateOne) && fld.EntField.Immutable) {
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
		varName := camel(reqVar + "_" + fld.EntField.Name)
		g.Tmpl(`%(extract)
		m.Set%(entField)(%(varName))`, g.withGlobals(tmplValues{
			"entField": entField,
			"varName":  varName,
			"extract":  g.renderToEnt(convert, varName, fmt.Sprintf("%s.Get%s()", reqVar, fld.PbStructField())),
		}))
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
			varName := camel(reqVar + "_" + edg.EntEdge.StructField())
			g.Tmpl(`%(extract)
			m.Set%(edgeName)ID(%(varName))`, g.withGlobals(tmplValues{
				"edgeName": edg.EntEdge.StructField(),
				"varName":  varName,
				"extract":  g.renderToEnt(convert, varName, fmt.Sprintf("%s.Get%s().Get%s()", reqVar, edg.PbStructField(), edg.EdgeIDPbStructField())),
			}))
		} else {
			varName := camel(edg.EntEdge.StructField())
			g.Tmpl(`for _, item := range %(reqVar).Get%(edgeField)() {
	%(extract)
	m.Add%(edgeName)IDs(%(varName))
}
`, g.withGlobals(tmplValues{
				"reqVar":    reqVar,
				"edgeField": edg.PbStructField(),
				"edgeName":  singular(edg.EntEdge.StructField()),
				"varName":   varName,
				"extract":   g.renderToEnt(convert, varName, fmt.Sprintf("item.Get%s()", edg.EdgeIDPbStructField())),
			}))
		}
	}
	g.P("res, err := m.Save(ctx)")

	g.Tmpl(`
	switch {
	case err == nil:
		proto, err := toProto%(typeName)(res)
		if err != nil {
			return nil, %(statusErrf)(%(internal), "internal: %s", err)
		}
		return proto, nil
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
		"typeName":          g.entType.Name,
		"fmtErr":            protogen.GoImportPath("fmt").Ident("Errorf"),
	}
	for _, additional := range additionals {
		for k, v := range additional {
			m[k] = v
		}
	}
	return m
}
