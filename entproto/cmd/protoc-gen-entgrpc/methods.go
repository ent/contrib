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

func (g *serviceGenerator) generateCreateMethod() error {
	reqVar := camel(g.typeName)
	g.Tmpl(`%(reqVar) := req.Get%(typeName)()
	created, err := svc.client.%(typeName).Create().`, tmplValues{
		"reqVar":   reqVar,
		"typeName": g.typeName},
	)

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
		return toProto%(typeName)(created), nil
	case %(uniqConstraintErr)(err):
		return nil, %(grpcStatusErrorf)(%(alreadyExists), "already exists: %s", err)
	case %(constraintErr)(err):
		return nil, %(grpcStatusErrorf)(%(invalidArgument), "invalid argument: %s", err)
	default:
		return nil, %(grpcStatusErrorf)(%(internal), "internal: %s", err)
	}`, tmplValues{
		"uniqConstraintErr": protogen.GoImportPath("entgo.io/ent/dialect/sql/sqlgraph").Ident("IsUniqueConstraintError"),
		"constraintErr":     g.entPackage.Ident("IsConstraintError"),
		"grpcStatusErrorf":  status.Ident("Errorf"),
		"alreadyExists":     codes.Ident("AlreadyExists"),
		"invalidArgument":   codes.Ident("InvalidArgument"),
		"internal":          codes.Ident("Internal"),
		"typeName":          g.typeName,
	})

	return nil
}
