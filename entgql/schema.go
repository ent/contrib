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

package entgql

import (
	"fmt"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

type (
	// TODO(giautm): refactor internal APIs
	schemaGenerator struct {
		graph     *gen.Graph
		nodes     []*gen.Type
		relaySpec bool
	}
)

func newSchemaGenerator(g *gen.Graph) (*schemaGenerator, error) {
	nodes, err := filterNodes(g.Nodes)
	if err != nil {
		return nil, err
	}

	return &schemaGenerator{
		graph: g,
		nodes: nodes,
		// TODO(giautm): relaySpec enable by default.
		// Add an option to disable it.
		relaySpec: true,
	}, nil
}

func (e *schemaGenerator) prepareSchema() (*ast.Schema, error) {
	types, err := e.buildTypes()
	if err != nil {
		return nil, err
	}
	if e.relaySpec {
		insertDefinitions(types, relayBuiltinTypes()...)
	}

	return &ast.Schema{
		Types: types,
	}, nil
}

func (e *schemaGenerator) buildTypes() (map[string]*ast.Definition, error) {
	types := make(map[string]*ast.Definition)
	for _, node := range e.nodes {
		ant, err := decodeAnnotation(node.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip {
			continue
		}

		fields, err := e.buildTypeFields(node)
		if err != nil {
			return nil, err
		}

		name := node.Name
		if ant.Type != "" {
			name = ant.Type
		}

		var interfaces []string
		if e.relaySpec {
			interfaces = append(interfaces, "Node")
		}

		types[name] = &ast.Definition{
			Name:       name,
			Kind:       ast.Object,
			Fields:     fields,
			Interfaces: interfaces,
		}

		var enumOrderByValues ast.EnumValueList
		for _, f := range node.Fields {
			ant, err := decodeAnnotation(f.Annotations)
			if err != nil {
				return nil, err
			}
			if ant.Skip {
				continue
			}

			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				enumOrderByValues = append(enumOrderByValues, &ast.EnumValueDefinition{
					Name: ant.OrderField,
				})
			}

			enum, err := e.buildEnum(f, ant)
			if err != nil {
				return nil, err
			}
			if enum != nil {
				types[enum.Name] = enum
			}
		}

		// TODO(giautm): Added RelayConnection annotation check
		if e.relaySpec {
			defs, err := relayConnectionTypes(node)
			if err != nil {
				return nil, err
			}

			insertDefinitions(types, defs...)
			if enumOrderByValues != nil {
				pagination, err := nodePaginationNames(node)
				if err != nil {
					return nil, err
				}

				types[pagination.OrderField] = &ast.Definition{
					Name:       pagination.OrderField,
					Kind:       ast.Enum,
					EnumValues: enumOrderByValues,
				}
				types[pagination.Order] = &ast.Definition{
					Name: pagination.Order,
					Kind: ast.InputObject,
					Fields: ast.FieldList{
						{
							Name: "direction",
							Type: ast.NonNullNamedType("OrderDirection", nil),
							DefaultValue: &ast.Value{
								Raw:  "ASC",
								Kind: ast.EnumValue,
							},
						},
						{
							Name: "field",
							Type: ast.NonNullNamedType(pagination.OrderField, nil),
						},
					},
				}
			}
		}
	}

	return types, nil
}

func (e *schemaGenerator) buildEnum(f *gen.Field, ant *Annotation) (*ast.Definition, error) {
	if !f.IsEnum() {
		return nil, nil
	}

	// NOTE(giautm): I'm not sure this is
	// the right approach, but it passed the test
	defs, err := e.typeFromField(f, false, ant.Type)
	if err != nil {
		return nil, err
	}
	name := defs.Name()

	valueDefs := make(ast.EnumValueList, 0, len(f.Enums))
	for _, v := range f.Enums {
		valueDefs = append(valueDefs, &ast.EnumValueDefinition{
			Name: v.Value,
		})
	}

	return &ast.Definition{
		Name:        name,
		Kind:        ast.Enum,
		Description: fmt.Sprintf("%s is enum for the field %s", name, f.Name),
		EnumValues:  valueDefs,
	}, nil
}

func (e *schemaGenerator) buildTypeFields(t *gen.Type) (ast.FieldList, error) {
	var fields ast.FieldList
	if t.ID != nil {
		f, err := e.typeField(t.ID, true)
		if err != nil {
			return nil, err
		}
		if f != nil {
			fields = append(fields, f...)
		}
	}

	for _, f := range t.Fields {
		f, err := e.typeField(f, false)
		if err != nil {
			return nil, err
		}
		if f != nil {
			fields = append(fields, f...)
		}
	}
	return fields, nil
}

func (e *schemaGenerator) typeField(f *gen.Field, isID bool) ([]*ast.FieldDefinition, error) {
	ant, err := decodeAnnotation(f.Annotations)
	if err != nil {
		return nil, err
	}
	if ant.Skip {
		return nil, nil
	}

	ft, err := e.typeFromField(f, isID, ant.Type)
	if err != nil {
		return nil, fmt.Errorf("field(%s): %w", f.Name, err)
	}

	// TODO(giautm): support rename field
	// TODO(giautm): support mapping single field to multiple GQL fields
	return []*ast.FieldDefinition{
		{
			Name: camel(f.Name),
			Type: ft,
		},
	}, nil
}

func namedType(name string, nullable bool) *ast.Type {
	if nullable {
		return ast.NamedType(name, nil)
	}
	return ast.NonNullNamedType(name, nil)
}

func (e *schemaGenerator) typeFromField(f *gen.Field, idField bool, userDefinedType string) (*ast.Type, error) {
	nillable := f.Nillable
	typ := f.Type.Type

	// TODO(giautm): Support custom scalar types
	// TODO(giautm): Support Edge Field
	// TODO(giautm): Support some built-in JSON types: Ints(), Floats(), Strings()
	scalar := f.Type.String()
	switch {
	case userDefinedType != "":
		return namedType(userDefinedType, nillable), nil
	case idField:
		return namedType("ID", !e.relaySpec && nillable), nil
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
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		return namedType(scalar, nillable), nil
	case typ == field.TypeJSON:
		return nil, fmt.Errorf("json type not implemented")
	case typ == field.TypeOther:
		return nil, fmt.Errorf("other type must have typed defined")
	default:
		return nil, fmt.Errorf("unexpected type: %s", typ.String())
	}
}

func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb).
		FormatSchema(schema)
	return sb.String()
}
