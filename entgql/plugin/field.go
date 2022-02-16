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

package plugin

import (
	"fmt"
	"reflect"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/vektah/gqlparser/v2/ast"
)

func (e *EntGQL) typeFields(t *gen.Type) (ast.FieldList, error) {
	var fields ast.FieldList
	if t.ID != nil {
		f, err := e.typeField(t.ID, true)
		if err != nil {
			return nil, err
		}
		if f != nil {
			fields = append(fields, f)
		}
	}
	for _, f := range t.Fields {
		f, err := e.typeField(f, false)
		if err != nil {
			return nil, err
		}
		if f != nil {
			fields = append(fields, f)
		}
	}
	return fields, nil
}

func (e *EntGQL) typeField(f *gen.Field, idField bool) (*ast.FieldDefinition, error) {
	ann := &entgql.Annotation{}
	err := ann.Decode(f.Annotations[ann.Name()])
	if err != nil {
		return nil, err
	}
	if ann.Skip {
		return nil, nil
	}
	ft, err := e.entTypToGqlType(f, idField, ann.Type)
	if err != nil {
		return nil, fmt.Errorf("field(%s): %w", f.Name, err)
	}
	return &ast.FieldDefinition{
		Name:       camel(f.Name),
		Type:       ft,
		Directives: e.directives(ann.GQLDirectives),
	}, nil
}

func namedType(name string, nillable bool) *ast.Type {
	if !nillable {
		return ast.NonNullNamedType(name, nil)
	}
	return ast.NamedType(name, nil)
}

func listNamedType(name string, nillable bool) *ast.Type {
	if !nillable {
		return ast.NonNullListType(namedType(name, false), nil)
	}
	return ast.ListType(namedType(name, false), nil)
}

func (e *EntGQL) entTypToGqlType(f *gen.Field, idField bool, userDefinedType string) (*ast.Type, error) {
	nillable := f.Nillable
	typ := f.Type.Type
	typeName := strings.TrimPrefix(typ.ConstName(), "Type")
	switch {
	case userDefinedType != "":
		return namedType(userDefinedType, nillable), nil
	case e.scalarMappings[typeName] != "":
		return namedType(e.scalarMappings[typeName], nillable), nil
	case idField:
		// Id cannot be null for node interface
		return namedType("ID", false), nil
	case f.IsEnum():
		// Guess enum type
		return namedType(strings.Title(f.Name), nillable), nil
	case typ.Float():
		return namedType("Float", nillable), nil
	case typ.Integer():
		return namedType("Int", nillable), nil
	case typ == field.TypeString:
		return namedType("String", nillable), nil
	case typ == field.TypeBool:
		return namedType("Boolean", nillable), nil
	case typ == field.TypeBytes:
		return nil, fmt.Errorf("bytes type not implemented")
	case typ == field.TypeJSON:
		if f.Type.RType != nil {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				switch f.Type.RType.Ident {
				case "[]float64":
					return listNamedType("Float", f.Optional), nil
				case "[]int":
					return listNamedType("Int", f.Optional), nil
				case "[]string":
					return listNamedType("String", f.Optional), nil
				}
			}
		}
		return nil, fmt.Errorf("json type not implemented")
	case typ == field.TypeOther:
		return nil, fmt.Errorf("other type must have typed defined")
	default:
		return nil, fmt.Errorf("unexpected type: %s", typ.String())
	}
}
