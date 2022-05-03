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
	"github.com/99designs/gqlgen/codegen/config"
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

type schemaGenerator struct {
	relaySpec     bool
	genSchema     bool
	genWhereInput bool

	cfg         *config.Config
	scalarFunc  func(*gen.Field, gen.Op) string
	schemaHooks []SchemaHook
}

func newSchemaGenerator() *schemaGenerator {
	return &schemaGenerator{
		relaySpec: true,
	}
}

func (e *schemaGenerator) BuildSchema(g *gen.Graph) (s *ast.Schema, err error) {
	s = &ast.Schema{
		Directives: map[string]*ast.DirectiveDefinition{},
	}
	if e.genSchema {
		s.AddTypes(builtinTypes()...)
		if e.relaySpec {
			s.AddTypes(relayBuiltinTypes(g.Package)...)
		}
		for name, d := range directives {
			s.Directives[name] = d
		}
	}
	if err := e.buildTypes(g, s); err != nil {
		return nil, err
	}

	for _, h := range e.schemaHooks {
		if err = h(g, s); err != nil {
			return nil, err
		}
	}

	return s, nil
}

func (e *schemaGenerator) buildTypes(g *gen.Graph, s *ast.Schema) error {
	var queryFields ast.FieldList
	if e.relaySpec {
		queryFields = relayBuiltinQueryFields()
	}

	for _, node := range g.Nodes {
		gqlType, ant, err := gqlTypeFromNode(node)
		if err != nil {
			return err
		}
		names := paginationNames(gqlType)

		if e.genSchema && !ant.Skip.Is(SkipType) {
			def, err := e.buildType(node, ant, gqlType, g.Package)
			if err != nil {
				return err
			}
			if def != nil {
				s.AddTypes(def)
			}
		}

		if e.genSchema && !ant.Skip.Is(SkipEnumField) {
			for _, f := range node.Fields {
				ant, err := annotation(f.Annotations)
				if err != nil {
					return err
				}
				if ant.Skip.Is(SkipEnumField) {
					continue
				}
				if f.IsEnum() {
					fieldType, err := e.typeFromField(f, false, ant.Type)
					if err != nil {
						return err
					}

					goType, ok := e.enumGoType(f, g.Package)
					if !ok {
						return fmt.Errorf("unexpected missing GoType info for enum %q", f.Name)
					}

					def, err := e.buildFieldEnum(f, fieldType.Name(), goType)
					if err != nil {
						return err
					}
					if def != nil {
						s.AddTypes(def)
					}
				}
			}
		}

		if e.genSchema && !ant.Skip.Is(SkipOrderField) {
			def, err := e.enumOrderByValues(node, names.OrderField)
			if err != nil {
				return err
			}
			if def != nil {
				def.Description = fmt.Sprintf("Properties by which %s connections can be ordered.", gqlType)
				s.AddTypes(def, names.OrderInputDef())
			}
		}

		if e.genSchema {
			if ant.RelayConnection {
				if !e.relaySpec {
					return ErrRelaySpecDisabled
				}
				s.AddTypes(names.TypeDefs()...)

				if ant.QueryField != nil {
					name := ant.QueryField.fieldName(gqlType)
					_, hasOrderBy := s.Types[names.Order]
					hasWhereInput := e.genWhereInput && !ant.Skip.Is(SkipWhereInput)

					def := names.ConnectionField(name, hasOrderBy, hasWhereInput)
					def.Directives = e.buildDirectives(ant.QueryField.Directives)
					queryFields = append(queryFields, def)
				}
			} else if ant.QueryField != nil {
				name := ant.QueryField.fieldName(gqlType)
				def := &ast.FieldDefinition{
					Name: name,
					Type: listNamedType(gqlType, false),
				}
				def.Directives = e.buildDirectives(ant.QueryField.Directives)
				queryFields = append(queryFields, def)
			}
		}

		if e.genWhereInput && !ant.Skip.Is(SkipWhereInput) {
			def, err := e.buildWhereInput(node, names.WhereInput)
			if err != nil {
				return err
			}
			if def != nil {
				s.AddTypes(def)
			}
		}
	}

	if e.genSchema && len(queryFields) > 0 {
		s.AddTypes(&ast.Definition{
			Name:   QueryType,
			Kind:   ast.Object,
			Fields: queryFields,
		})
	}

	return nil
}

func (e *schemaGenerator) buildType(t *gen.Type, ant *Annotation, gqlType, pkg string) (*ast.Definition, error) {
	def := &ast.Definition{
		Name:       gqlType,
		Kind:       ast.Object,
		Fields:     ast.FieldList{},
		Directives: e.buildDirectives(ant.Directives),
		Interfaces: []string{},
	}
	if t.Name != gqlType {
		def.Directives = append(def.Directives, goModel(entGoType(t.Name, pkg)))
	}
	if e.relaySpec {
		def.Interfaces = append(def.Interfaces, "Node")
	}
	if len(ant.Implements) > 0 {
		def.Interfaces = append(def.Interfaces, ant.Implements...)
	}

	if t.ID != nil {
		f, err := e.fieldDefinition(t.ID, true)
		if err != nil {
			return nil, err
		}
		if f != nil {
			def.Fields = append(def.Fields, f...)
		}
	}

	for _, f := range t.Fields {
		f, err := e.fieldDefinition(f, false)
		if err != nil {
			return nil, err
		}
		if f != nil {
			def.Fields = append(def.Fields, f...)
		}
	}

	for _, edge := range t.Edges {
		ant, err := annotation(edge.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(SkipType) {
			continue
		}
		if ant.RelayConnection && edge.Unique {
			return nil, fmt.Errorf("RelayConnection cannot be defined on Unique edge: %s.%s", t.Name, edge.Name)
		}

		fields, err := e.buildEdge(edge, ant)
		if err != nil {
			return nil, err
		}
		if len(fields) > 0 {
			def.Fields = append(def.Fields, fields...)
		}
	}

	return def, nil
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

func (e *schemaGenerator) enumOrderByValues(t *gen.Type, gqlType string) (*ast.Definition, error) {
	enumValues := ast.EnumValueList{}
	for _, f := range t.Fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(SkipOrderField) {
			continue
		}

		// Check if this node has an OrderBy object
		if ant.OrderField != "" {
			if ant.Skip.Is(SkipOrderField) {
				return nil, fmt.Errorf("entgql: ordered field %s.%s cannot be skipped", t.Name, f.Name)
			}
			enumValues = append(enumValues, &ast.EnumValueDefinition{
				Name: ant.OrderField,
			})
		}
	}
	if len(enumValues) == 0 {
		return nil, nil
	}

	return &ast.Definition{
		Name:       gqlType,
		Kind:       ast.Enum,
		EnumValues: enumValues,
	}, nil
}

func (e *schemaGenerator) buildFieldEnum(f *gen.Field, gqlType, goType string) (*ast.Definition, error) {
	enumValues := make(ast.EnumValueList, 0, len(f.Enums))
	for _, v := range f.Enums {
		enumValues = append(enumValues, &ast.EnumValueDefinition{
			Name: v.Value,
		})
	}
	return &ast.Definition{
		Name:        gqlType,
		Kind:        ast.Enum,
		Description: fmt.Sprintf("%s is enum for the field %s", gqlType, f.Name),
		EnumValues:  enumValues,
		Directives:  ast.DirectiveList{goModel(goType)},
	}, nil
}

func (e *schemaGenerator) buildEdge(edge *gen.Edge, edgeAnt *Annotation) ([]*ast.FieldDefinition, error) {
	gqlType, ant, err := gqlTypeFromNode(edge.Type)
	if err != nil {
		return nil, err
	}
	orderFields, err := orderFields(edge.Type)
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
				ConnectionField(name, len(orderFields) > 0,
					e.genWhereInput && !edgeAnt.Skip.Is(SkipWhereInput) && !ant.Skip.Is(SkipWhereInput),
				)
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

// buildWhereInput returns the a <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *schemaGenerator) buildWhereInput(t *gen.Type, gqlType string) (*ast.Definition, error) {
	def := &ast.Definition{
		Name:        gqlType,
		Kind:        ast.InputObject,
		Description: fmt.Sprintf("%s is used for filtering %s objects.\nInput was generated by ent.", gqlType, t.Name),
		Fields: ast.FieldList{
			&ast.FieldDefinition{
				Name: "not",
				Type: ast.NamedType(gqlType, nil),
			},
		},
	}

	for _, op := range []string{"and", "or"} {
		def.Fields = append(def.Fields, &ast.FieldDefinition{
			Name: op,
			Type: ast.ListType(ast.NonNullNamedType(gqlType, nil), nil),
		})
	}

	fields, err := filterFields(append(t.Fields, t.ID), SkipWhereInput)
	if err != nil {
		return nil, err
	}
	for _, f := range fields {
		if !f.Type.Comparable() {
			continue
		}
		for i, op := range f.Ops() {
			fd := e.fieldDefinitionOp(f, op)
			if i == 0 {
				fd.Description = f.Name + " field predicates"
			}
			def.Fields = append(def.Fields, fd)
		}
	}

	edges, err := filterEdges(t.Edges, SkipWhereInput)
	if err != nil {
		return nil, err
	}
	for _, e := range edges {
		names, err := nodePaginationNames(e.Type)
		if err != nil {
			return nil, err
		}

		def.Fields = append(def.Fields,
			&ast.FieldDefinition{
				Name:        camel("has_" + e.Name),
				Type:        namedType("Boolean", true),
				Description: e.Name + " edge predicates",
			},
			&ast.FieldDefinition{
				Name: camel("has_" + e.Name + "_with"),
				Type: listNamedType(names.WhereInput, true),
			},
		)
	}
	return def, nil
}

func (e *schemaGenerator) fieldDefinition(f *gen.Field, isID bool) ([]*ast.FieldDefinition, error) {
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
			Name:        camel(f.Name),
			Type:        ft,
			Description: f.Comment(),
			Directives:  e.buildDirectives(ant.Directives),
		},
	}, nil
}

func (e *schemaGenerator) fieldDefinitionOp(f *gen.Field, op gen.Op) *ast.FieldDefinition {
	def := &ast.FieldDefinition{
		Name: camel(f.Name + "_" + op.Name()),
	}
	if op == gen.EQ {
		def.Name = camel(f.Name)
	}

	typeName := e.mapScalar(f, op)
	if op.Variadic() {
		def.Type = listNamedType(typeName, true)
	} else {
		def.Type = namedType(typeName, true)
	}
	return def
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

// mapScalar provides maps an ent.Schema type into GraphQL scalar type.
// In order to override this function, use the WithMapScalarFunc option.
func (e *schemaGenerator) mapScalar(f *gen.Field, op gen.Op) string {
	if e.scalarFunc != nil {
		if t := e.scalarFunc(f, op); t != "" {
			return t
		}
	}
	scalar := f.Type.String()
	switch t := f.Type.Type; {
	case op.Niladic():
		return "Boolean"
	case f.Name == "id":
		return "ID"
	case t == field.TypeBool:
		scalar = "Boolean"
	case f.IsEdgeField():
		scalar = "ID"
	case t.Float():
		scalar = "Float"
	case t.Numeric():
		scalar = "Int"
	case t == field.TypeString:
		scalar = "String"
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		if typ, ok := e.hasMapping(f); ok {
			scalar = typ
		} else {
			scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		}
	}
	if ant, err := annotation(f.Annotations); err == nil && ant.Type != "" {
		return ant.Type
	}
	return scalar
}

// hasMapping reports if the gqlgen.yml has custom mapping for
// the given field type and returns its GraphQL name if exists.
func (e *schemaGenerator) hasMapping(f *gen.Field) (string, bool) {
	if e.cfg == nil {
		return "", false
	}

	var gqlNames []string
	for t, v := range e.cfg.Models {
		// The string representation uses shortened package
		// names, and we override it for custom Go types.
		ident := f.Type.String()
		if idx := strings.IndexByte(ident, '.'); idx != -1 && f.HasGoType() && f.Type.PkgPath != "" {
			ident = f.Type.PkgPath + ident[idx:]
		}
		for _, m := range v.Model {
			// A mapping was found from GraphQL name to field type.
			if strings.HasSuffix(m, ident) {
				gqlNames = append(gqlNames, t)
			}
		}
	}
	if count := len(gqlNames); count == 1 {
		return gqlNames[0], true
	} else if count > 1 {
		// If there is more than 1 mapping, we accept the one with the "Input" suffix.
		for _, t := range gqlNames {
			if strings.HasSuffix(t, "Input") {
				return t, true
			}
		}
	}
	// If no custom mapping was found, fallback to the builtin scalar
	// types as mentioned in https://gqlgen.com/reference/scalars
	switch f.Type.String() {
	case "time.Time":
		return "Time", true
	case "map[string]interface{}":
		return "Map", true
	default:
		return "", false
	}
}

func (e *schemaGenerator) enumGoType(f *gen.Field, pkg string) (string, bool) {
	if f.IsEnum() {
		if f.HasGoType() {
			return entGoType(f.Type.RType.Name, f.Type.RType.PkgPath), true
		}
		return fmt.Sprintf("%s/%s", pkg, f.Type.Ident), true
	}

	return "", false
}

func entGoType(name, pkg string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func builtinTypes() []*ast.Definition {
	return []*ast.Definition{
		{
			Name:        OrderDirection,
			Kind:        ast.Enum,
			Description: "Possible directions in which to order a list of items when provided an `orderBy` argument.",
			EnumValues: []*ast.EnumValueDefinition{
				{
					Name:        "ASC",
					Description: "Specifies an ascending order for a given `orderBy` argument.",
				},
				{
					Name:        "DESC",
					Description: "Specifies a descending order for a given `orderBy` argument.",
				},
			},
		},
	}
}

func relayBuiltinQueryFields() ast.FieldList {
	var (
		idType  = ast.NonNullNamedType("ID", nil)
		nodeDef = ast.NamedType(RelayNode, nil)
	)
	return ast.FieldList{
		{
			Name:        "node",
			Type:        nodeDef,
			Description: "Fetches an object given its ID.",
			Arguments: ast.ArgumentDefinitionList{
				{
					Name:        "id",
					Type:        idType,
					Description: "ID of the object.",
				},
			},
		},
		{
			Name:        "nodes",
			Type:        ast.NonNullListType(nodeDef, nil),
			Description: "Lookup nodes by a list of IDs.",
			Arguments: ast.ArgumentDefinitionList{
				{
					Name:        "ids",
					Type:        ast.NonNullListType(idType, nil),
					Description: "The list of node IDs.",
				},
			},
		},
	}
}

func relayBuiltinTypes(pkg string) []*ast.Definition {
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
			Directives: []*ast.Directive{
				goModel(entGoType("Noder", pkg)),
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
