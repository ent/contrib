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
	// QueryType is the name of the root Query object.
	QueryType = "Query"
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
	nodes, err := filterNodes(g.Nodes, SkipType)
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

func (e *schemaGenerator) buildSchema(s *ast.Schema) error {
	err := e.buildTypes(s)
	if err != nil {
		return err
	}
	s.AddTypes(builtinTypes()...)
	if e.relaySpec {
		s.AddTypes(relayBuiltinTypes()...)
	}

	for name, d := range directives {
		s.Directives[name] = d
	}
	return nil
}

func (e *schemaGenerator) buildTypes(s *ast.Schema) error {
	var (
		defaultInterfaces []string
		queryFields       ast.FieldList
	)
	if e.relaySpec {
		defaultInterfaces = append(defaultInterfaces, "Node")

		var (
			idType  = ast.NonNullNamedType("ID", nil)
			nodeDef = ast.NamedType(RelayNode, nil)
		)
		queryFields = append(queryFields,
			&ast.FieldDefinition{
				Name: "node",
				Type: nodeDef,
				Arguments: ast.ArgumentDefinitionList{
					{Name: "id", Type: idType},
				},
			},
			&ast.FieldDefinition{
				Name: "nodes",
				Arguments: ast.ArgumentDefinitionList{
					{Name: "ids", Type: ast.NonNullListType(idType, nil)},
				},
				Type: ast.NonNullListType(nodeDef, nil),
			},
		)
	}

	for _, node := range e.nodes {
		gqlType, ant, err := gqlTypeFromNode(node)
		if err != nil {
			return err
		}
		if ant.Skip.Is(SkipType) {
			continue
		}

		fields, err := e.buildTypeFields(node)
		if err != nil {
			return err
		}
		typ := &ast.Definition{
			Name:       gqlType,
			Kind:       ast.Object,
			Fields:     fields,
			Directives: e.buildDirectives(ant.Directives),
			Interfaces: defaultInterfaces,
		}
		if node.Name != gqlType {
			typ.Directives = append(typ.Directives, goModel(e.entGoType(node.Name)))
		}
		if len(ant.Implements) > 0 {
			typ.Interfaces = append(typ.Interfaces, ant.Implements...)
		}
		s.AddTypes(typ)

		var enumOrderByValues []string
		for _, f := range node.Fields {
			ant, err := annotation(f.Annotations)
			if err != nil {
				return err
			}
			if ant.Skip.Is(SkipType) {
				continue
			}

			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				if ant.Skip.Is(SkipOrderField) {
					return fmt.Errorf("entgql: ordered field %s.%s cannot be skipped", node.Name, f.Name)
				}
				enumOrderByValues = append(enumOrderByValues, ant.OrderField)
			}

			if f.IsEnum() && !ant.Skip.Is(SkipEnumField) {
				enum, err := e.buildEnum(f, ant)
				if err != nil {
					return err
				}
				s.AddTypes(enum)
			}
		}

		for _, edge := range node.Edges {
			ant, err := annotation(edge.Annotations)
			if err != nil {
				return err
			}
			if ant.Skip.Is(SkipType) {
				continue
			}
			if ant.RelayConnection && edge.Unique {
				return fmt.Errorf("RelayConnection cannot be defined on Unique edge: %s.%s", node.Name, edge.Name)
			}
			fields, err := e.buildEdge(edge, ant)
			if err != nil {
				return err
			}
			if len(fields) > 0 {
				typ.Fields = append(typ.Fields, fields...)
			}
		}

		if ant.RelayConnection {
			if !e.relaySpec {
				return ErrRelaySpecDisabled
			}

			pagination := paginationNames(gqlType)

			s.AddTypes(pagination.TypeDefs()...)
			if len(enumOrderByValues) > 0 && !ant.Skip.Is(SkipOrderField) {
				s.AddTypes(pagination.OrderByTypeDefs(enumOrderByValues)...)
			}
		}

		if ant.QueryField != "" {
			if ant.RelayConnection {
				pagination := paginationNames(gqlType)
				queryFields = append(queryFields, pagination.ConnectionField(ant.QueryField))
			} else {
				queryFields = append(queryFields, &ast.FieldDefinition{
					Name: ant.QueryField,
					Type: listNamedType(gqlType, false),
				})
			}
		}
	}

	if len(queryFields) > 0 {
		s.AddTypes(&ast.Definition{
			Name:   QueryType,
			Kind:   ast.Object,
			Fields: queryFields,
		})
	}

	return nil
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
	goType, ok := e.enumGoType(f)
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

func (e *schemaGenerator) buildEdge(edge *gen.Edge, edgeAnt *Annotation) ([]*ast.FieldDefinition, error) {
	gqlType, ant, err := gqlTypeFromNode(edge.Type)
	if err != nil {
		return nil, err
	}

	var (
		edgeField = camel(edge.Name)
		mappings  = []string{edgeField}
	)
	if len(edgeAnt.Mapping) > 0 {
		mappings = edgeAnt.Mapping
	}

	var fields []*ast.FieldDefinition
	for _, name := range mappings {
		fieldDef := &ast.FieldDefinition{Name: name}
		switch {
		case edge.Unique:
			fieldDef.Type = namedType(gqlType, edge.Optional)
		case edgeAnt.RelayConnection:
			if !e.relaySpec {
				return nil, ErrRelaySpecDisabled
			}
			if !ant.RelayConnection {
				return nil, fmt.Errorf("entgql: must enable Relay Connection via the entgql.RelayConnection annotation on the %s entity", edge.Type.Name)
			}

			fieldDef = paginationNames(gqlType).
				ConnectionField(name)
		default:
			fieldDef.Type = listNamedType(gqlType, edge.Optional)
		}

		fieldDef.Directives = e.buildDirectives(edgeAnt.Directives)
		if name != edgeField {
			fieldDef.Directives = append(fieldDef.Directives, goField(edgeField))
		}
		fields = append(fields, fieldDef)
	}

	return fields, nil
}

func (e *schemaGenerator) typeField(f *gen.Field, isID bool) ([]*ast.FieldDefinition, error) {
	ant, err := annotation(f.Annotations)
	if err != nil {
		return nil, err
	}
	if ant.Skip.Is(SkipType) {
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
		return namedType("ID", false), nil
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
		gqlType, ant, err := gqlTypeFromNode(node)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(SkipType) {
			continue
		}

		models[gqlType] = e.entGoType(node.Name)

		var hasOrderBy bool
		for _, field := range node.Fields {
			ant, err := annotation(field.Annotations)
			if err != nil {
				return nil, err
			}
			if ant.Skip.Is(SkipType) {
				continue
			}
			// Check if this node has an OrderBy object
			if ant.OrderField != "" {
				hasOrderBy = true
			}

			// We only map the Go types generated by Ent, for example: enum.
			goType, ok := e.enumGoType(field)
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

func (e *schemaGenerator) enumGoType(f *gen.Field) (string, bool) {
	if f.IsEnum() {
		if f.HasGoType() {
			return fmt.Sprintf("%s.%s", f.Type.RType.PkgPath, f.Type.RType.Name), true
		}
		return fmt.Sprintf("%s/%s", e.graph.Package, f.Type.Ident), true
	}

	return "", false
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
	return pagination.TypeDefs(), nil
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

func goField(name string) *ast.Directive {
	return &ast.Directive{
		Name:     "goField",
		Location: ast.LocationFieldDefinition,
		Arguments: ast.ArgumentList{
			{
				Name: "name",
				Value: &ast.Value{
					Kind: ast.StringValue,
					Raw:  name,
				},
			},
			{
				Name: "forceResolver",
				Value: &ast.Value{
					Kind: ast.BooleanValue,
					Raw:  "false",
				},
			},
		},
	}
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
