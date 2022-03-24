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
	"errors"
	"fmt"
	"reflect"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

const (
	// OrderDirection is the name of enum OrderDirection
	OrderDirection = "OrderDirection"
	// RelayCursor is the name of the cursor type
	RelayCursor = "Cursor"
	// RelayNode is the name of the interface that all nodes implement
	RelayNode = "Node"
	// RelayPageInfo is the name of the PageInfo type
	RelayPageInfo = "PageInfo"
)

var (
	// ErrRelaySpecDisabled is the error returned when the relay specification is disabled
	ErrRelaySpecDisabled = errors.New("entgql: must enable relay specification via the WithRelaySpec option")

	pos        = &ast.Position{Src: &ast.Source{BuiltIn: false}}
	directives = map[string]*ast.DirectiveDefinition{
		"goModel": {
			Name:     "goModel",
			Position: pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "model",
					Type: ast.NamedType("String", nil),
				},
				{
					Name: "models",
					Type: ast.ListType(ast.NonNullNamedType("String", nil), nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationObject,
				ast.LocationInputObject,
				ast.LocationScalar,
				ast.LocationEnum,
				ast.LocationInterface,
				ast.LocationUnion,
			},
		},
		"goField": {
			Name:     "goField",
			Position: pos,
			Arguments: ast.ArgumentDefinitionList{
				{
					Name: "forceResolver",
					Type: ast.NamedType("Boolean", nil),
				},
				{
					Name: "name",
					Type: ast.NamedType("String", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	}
)

// TODO(giautm): refactor internal APIs
type schemaGenerator struct {
	graph     *gen.Graph
	nodes     []*gen.Type
	relaySpec bool
}

func newSchemaGenerator(g *gen.Graph) (*schemaGenerator, error) {
	nodes, err := filterNodes(g.Nodes, SkipFlagType)
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
	insertDefinitions(types, builtinTypes()...)
	if e.relaySpec {
		insertDefinitions(types, relayBuiltinTypes()...)
	}
	return &ast.Schema{
		Types:      types,
		Directives: directives,
	}, nil
}

func (e *schemaGenerator) buildTypes() (map[string]*ast.Definition, error) {
	types := make(map[string]*ast.Definition)
	var defaultInterfaces []string
	if e.relaySpec {
		defaultInterfaces = append(defaultInterfaces, "Node")
	}
	for _, node := range e.nodes {
		ant, err := annotation(node.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Has(SkipFlagType) {
			continue
		}

		fields, err := e.buildTypeFields(node)
		if err != nil {
			return nil, err
		}
		typ := &ast.Definition{
			Name:       node.Name,
			Kind:       ast.Object,
			Fields:     fields,
			Directives: e.buildDirectives(ant.Directives),
			Interfaces: defaultInterfaces,
		}
		if ant.Type != "" {
			typ.Name = ant.Type
			typ.Directives = append(typ.Directives, goModel(e.entGoType(node.Name)))
		}
		if len(ant.Implements) > 0 {
			typ.Interfaces = append(typ.Interfaces, ant.Implements...)
		}
		insertDefinitions(types, typ)

		var enumOrderByValues ast.EnumValueList
		for _, f := range node.Fields {
			ant, err := annotation(f.Annotations)
			if err != nil {
				return nil, err
			}
			if ant.Skip.Has(SkipFlagType) {
				continue
			}

			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				enumOrderByValues = append(enumOrderByValues, &ast.EnumValueDefinition{
					Name: ant.OrderField,
				})
			}

			if f.IsEnum() && !ant.Skip.Has(SkipFlagFieldEnum) {
				enum, err := e.buildEnum(f, ant)
				if err != nil {
					return nil, err
				}
				insertDefinitions(types, enum)
			}
		}

		if ant.RelayConnection {
			if !e.relaySpec {
				return nil, ErrRelaySpecDisabled
			}

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

				insertDefinitions(types,
					&ast.Definition{
						Name:       pagination.OrderField,
						Kind:       ast.Enum,
						EnumValues: enumOrderByValues,
					},
					&ast.Definition{
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
					},
				)
			}
		}
	}

	return types, nil
}

func (e *schemaGenerator) buildDirectives(directives []Directive) ast.DirectiveList {
	list := make(ast.DirectiveList, 0, len(directives))
	for _, d := range directives {
		args := make(ast.ArgumentList, 0, len(d.Arguments))
		for _, a := range d.Arguments {
			args = append(args, &ast.Argument{
				Name: a.Name,
				Value: &ast.Value{
					Raw:  a.Value,
					Kind: a.Kind,
				},
			})
		}
		list = append(list, &ast.Directive{
			Name:      d.Name,
			Arguments: args,
		})
	}
	return list
}

func (e *schemaGenerator) buildEnum(f *gen.Field, ant *Annotation) (*ast.Definition, error) {
	goType, ok := e.fieldGoType(f)
	if !ok {
		return nil, fmt.Errorf("unexpected missing GoType info for enum %q", f.Name)
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
		Directives:  ast.DirectiveList{goModel(goType)},
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
	ant, err := annotation(f.Annotations)
	if err != nil {
		return nil, err
	}
	if ant.Skip.Has(SkipFlagType) {
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
			Name:       camel(f.Name),
			Type:       ft,
			Directives: e.buildDirectives(ant.Directives),
		},
	}, nil
}

func (e *schemaGenerator) typeFromField(f *gen.Field, idField bool, userDefinedType string) (*ast.Type, error) {
	nillable := f.Nillable
	typ := f.Type.Type

	// TODO(giautm): Support custom scalar types
	// TODO(giautm): Support Edge Field
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

func (e *schemaGenerator) genModels() (map[string]string, error) {
	models := make(map[string]string)

	if e.relaySpec {
		models[RelayPageInfo] = e.entGoType(RelayPageInfo)
		models[RelayNode] = e.entGoType("Noder")
		models[RelayCursor] = e.entGoType(RelayCursor)
	}
	for _, node := range e.nodes {
		ant, err := annotation(node.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Has(SkipFlagType) {
			continue
		}

		name := node.Name
		if ant.Type != "" {
			name = ant.Type
		}
		models[name] = e.entGoType(node.Name)

		var hasOrderBy bool
		for _, field := range node.Fields {
			ant, err := annotation(field.Annotations)
			if err != nil {
				return nil, err
			}
			if ant.Skip.Has(SkipFlagType) {
				continue
			}
			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				hasOrderBy = true
			}

			goType, ok := e.fieldGoType(field)
			if !ok {
				continue
			}
			// NOTE(giautm): I'm not sure this is
			// the right approach, but it passed the test
			defs, err := e.typeFromField(field, false, ant.Type)
			if err != nil {
				return nil, err
			}
			name := defs.Name()
			models[name] = goType
		}

		if ant.RelayConnection {
			if !e.relaySpec {
				return nil, ErrRelaySpecDisabled
			}
			pagination, err := nodePaginationNames(node)
			if err != nil {
				return nil, err
			}

			models[pagination.Connection] = e.entGoType(pagination.Connection)
			models[pagination.Edge] = e.entGoType(pagination.Edge)

			if hasOrderBy {
				models["OrderDirection"] = e.entGoType("OrderDirection")
				models[pagination.Order] = e.entGoType(pagination.Order)
				models[pagination.OrderField] = e.entGoType(pagination.OrderField)
			}
		}
	}

	return models, nil
}

func (e *schemaGenerator) entGoType(name string) string {
	return fmt.Sprintf("%s.%s", e.graph.Package, name)
}

func (e *schemaGenerator) fieldGoType(f *gen.Field) (string, bool) {
	switch {
	case f.IsOther() || (f.IsEnum() && f.HasGoType()):
		return fmt.Sprintf("%s.%s", f.Type.RType.PkgPath, f.Type.RType.Name), true
	case f.IsEnum():
		return fmt.Sprintf("%s/%s", e.graph.Package, f.Type.Ident), true
	default:
		return "", false
	}
}

func builtinTypes() []*ast.Definition {
	return []*ast.Definition{
		{
			Name: OrderDirection,
			Kind: ast.Enum,
			EnumValues: []*ast.EnumValueDefinition{
				{Name: "ASC"},
				{Name: "DESC"},
			},
		},
	}
}

func relayBuiltinTypes() []*ast.Definition {
	return []*ast.Definition{
		{
			Name: RelayCursor,
			Kind: ast.Scalar,
			Description: `Define a Relay Cursor type:
https://relay.dev/graphql/connections.htm#sec-Cursor`,
		},
		{
			Name: RelayNode,
			Kind: ast.Interface,
			Description: `An object with an ID.
Follows the [Relay Global Object Identification Specification](https://relay.dev/graphql/objectidentification.htm)`,
			Fields: []*ast.FieldDefinition{
				{
					Name:        "id",
					Type:        ast.NonNullNamedType("ID", nil),
					Description: "The id of the object.",
				},
			},
		},
		{
			Name: RelayPageInfo,
			Kind: ast.Object,
			Description: `Information about pagination in a connection.
https://relay.dev/graphql/connections.htm#sec-undefined.PageInfo`,
			Fields: []*ast.FieldDefinition{
				{
					Name:        "hasNextPage",
					Type:        ast.NonNullNamedType("Boolean", nil),
					Description: "When paginating forwards, are there more items?",
				},
				{
					Name:        "hasPreviousPage",
					Type:        ast.NonNullNamedType("Boolean", nil),
					Description: "When paginating backwards, are there more items?",
				},
				{
					Name:        "startCursor",
					Type:        ast.NamedType("Cursor", nil),
					Description: "When paginating backwards, the cursor to continue.",
				},
				{
					Name:        "endCursor",
					Type:        ast.NamedType("Cursor", nil),
					Description: "When paginating forwards, the cursor to continue.",
				},
			},
		},
	}
}

func relayConnectionTypes(t *gen.Type) ([]*ast.Definition, error) {
	pagination, err := nodePaginationNames(t)
	if err != nil {
		return nil, err
	}
	return []*ast.Definition{
		{
			Name:        pagination.Edge,
			Kind:        ast.Object,
			Description: "An edge in a connection.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "node",
					Type:        ast.NamedType(pagination.Node, nil),
					Description: "The item at the end of the edge.",
				},
				{
					Name:        "cursor",
					Type:        ast.NonNullNamedType("Cursor", nil),
					Description: "A cursor for use in pagination.",
				},
			},
		},
		{
			Name:        pagination.Connection,
			Kind:        ast.Object,
			Description: "A connection to a list of items.",
			Fields: []*ast.FieldDefinition{
				{
					Name:        "edges",
					Type:        ast.ListType(ast.NamedType(pagination.Edge, nil), nil),
					Description: "A list of edges.",
				},
				{
					Name:        "pageInfo",
					Type:        ast.NonNullNamedType("PageInfo", nil),
					Description: "Information to aid in pagination.",
				},
				{
					Name: "totalCount",
					Type: ast.NonNullNamedType("Int", nil),
				},
			},
		},
	}, nil
}

func insertDefinitions(types map[string]*ast.Definition, defs ...*ast.Definition) {
	for _, d := range defs {
		types[d.Name] = d
	}
}

func namedType(name string, nullable bool) *ast.Type {
	if nullable {
		return ast.NamedType(name, nil)
	}
	return ast.NonNullNamedType(name, nil)
}

func listNamedType(name string, nullable bool) *ast.Type {
	t := ast.NonNullNamedType(name, nil)
	if nullable {
		return ast.ListType(t, nil)
	}
	return ast.NonNullListType(t, nil)
}

func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb, formatter.WithIndent("  ")).
		FormatSchema(schema)
	return sb.String()
}

func goModel(ident string) *ast.Directive {
	return &ast.Directive{
		Name:     "goModel",
		Location: ast.LocationObject,
		Arguments: ast.ArgumentList{
			{
				Name: "model",
				Value: &ast.Value{
					Kind: ast.StringValue,
					Raw:  ident,
				},
			},
		},
	}
}
