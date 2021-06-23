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
	"entgo.io/contrib/entproto"
	"google.golang.org/protobuf/compiler/protogen"
)

func typeNeedsValidator(d entproto.FieldMap) bool {
	for _, fld := range d {
		if fieldNeedsValidator(fld) {
			return true
		}
	}
	return false
}

func fieldNeedsValidator(d *entproto.FieldMappingDescriptor) bool {
	f := d.EntField
	if d.IsEdgeField {
		f = d.EntEdge.Type.ID
	}
	return f.IsUUID()
}

// generateValidator generates a validation function for the service entity, to verify that
// the gRPC input is safe to pass to ent. Ent has already rich validation functionality and
// this layer should *only* assert invariants that are expected by ent but cannot be guaranteed
// by gRPC. For instance, TypeUUID is serialized as a proto bytes field, must be 16-bytes long.
func (g *serviceGenerator) generateValidator() {
	g.Tmpl(`
	// validate%(typeName) validates that all fields are encoded properly and are safe to pass
	// to the ent entity builder.
	func validate%(typeName)(x *%(typeName), checkId bool) error {`, g.withGlobals())
	for _, fld := range g.fieldMap.Fields() {
		if fieldNeedsValidator(fld) {
			var idCheckSuffix string
			if fld.IsIDField {
				idCheckSuffix = "&& checkId"
			}

			if fld.EntField.IsUUID() {
				g.Tmpl(`if err := %(validateUUID)(x.Get%(pbField)()); err != nil %(suffix) {
					return err
				}`, g.withGlobals(tmplValues{
					"pbField":      fld.PbStructField(),
					"validateUUID": protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ValidateUUID"),
					"suffix":       idCheckSuffix,
				}))
			}
		}
	}
	for _, edg := range g.fieldMap.Edges() {
		if fieldNeedsValidator(edg) {
			f := edg.EntEdge.Type.ID
			if f.IsUUID() {
				vars := g.withGlobals(tmplValues{
					"pbField":      edg.PbStructField(),
					"edgeIdField":  edg.EdgeIDPbStructField(),
					"validateUUID": protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ValidateUUID"),
				})
				if !edg.EntEdge.Unique {
					if !edg.EntEdge.Unique {
						g.Tmpl(`for _, item := range x.Get%(pbField)() {
	if err := %(validateUUID)(item.Get%(edgeIdField)()); err != nil {
		return err
	}
}
`, vars)
					}
				} else {
					g.Tmpl(`if err := %(validateUUID)(x.Get%(pbField)().Get%(edgeIdField)()); err != nil {
					return err
				}`, vars)
				}
			}
		}
	}
	g.P("return nil")
	g.P("}")
}

func (g *serviceGenerator) generateIDFieldValidator(idField *entproto.FieldMappingDescriptor) {
	if idField.EntField.IsUUID() {
		g.Tmpl(`if err := %(validateUUID)(req.Get%(pbField)()); err != nil {
					return nil, %(statusErrf)(%(invalidArgument), "invalid argument: %s", err)
				}`, g.withGlobals(tmplValues{
			"pbField":      idField.PbStructField(),
			"validateUUID": protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ValidateUUID"),
		}))
		return
	}
	panic("entproto: id field validation not implemented for " + idField.EntField.Type.String())
}
