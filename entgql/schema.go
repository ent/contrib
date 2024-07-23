// Copyright 2019-present Facebook
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
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
	"path/filepath"
	"reflect"
	"strings"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/99designs/gqlgen/codegen/templates"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

const (
	// QueryType is the name of the root Query object.
	QueryType = "Query"
	// OrderDirectionEnum is the name of enum OrderDirection
	OrderDirectionEnum = "OrderDirection"
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
				{
					Name: "forceGenerate",
					Type: ast.NamedType("Boolean", nil),
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
				{
					Name: "omittable",
					Type: ast.NamedType("Boolean", nil),
				},
			},
			Locations: []ast.DirectiveLocation{
				ast.LocationFieldDefinition,
				ast.LocationInputFieldDefinition,
			},
		},
	}
	inputObjectFilter    = func(t string) bool { return strings.HasSuffix(t, "Input") }
	nonInputObjectFilter = func(t string) bool { return !inputObjectFilter(t) }
)

type schemaGenerator struct {
	path          string
	relaySpec     bool
	genSchema     bool
	genWhereInput bool
	genMutations  bool

	cfg         *config.Config
	scalarFunc  func(*gen.Field, gen.Op) string
	schemaHooks []SchemaHook
}

func (e *schemaGenerator) BuildSchema(g *gen.Graph) (s *ast.Schema, err error) {
	s = &ast.Schema{
		Directives: make(map[string]*ast.DirectiveDefinition),
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
		if node.HasCompositeID() {
			continue
		}
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
				if s.Types[def.Name] != nil {
					return fmt.Errorf("found the GQL type conflict for the node %s, please use the entgql.Type() annotation to rename the GQL type", node.Name)
				}
				s.AddTypes(def)
				e.mayAddScalars(s, def)
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
					gqlType := e.mapScalar(gqlType, f, ant, nonInputObjectFilter)
					if gqlType == "" {
						return errors.New("unable to map enum field " + f.Name)
					}
					def, err := e.buildFieldEnum(f, gqlType, fieldGoType(f, g.Package))
					if err != nil {
						return err
					}
					if def != nil {
						if s.Types[def.Name] != nil {
							continue
						}
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

					def := names.ConnectionField(name, hasOrderBy, ant.MultiOrder, hasWhereInput)
					def.Description = ant.QueryField.Description
					def.Directives = e.buildDirectives(ant.QueryField.Directives)
					queryFields = append(queryFields, def)
				}
			} else if ant.QueryField != nil {
				name := ant.QueryField.fieldName(gqlType)
				def := &ast.FieldDefinition{
					Name:        name,
					Description: ant.QueryField.Description,
					Type:        listNamedType(gqlType, false),
				}
				def.Directives = e.buildDirectives(ant.QueryField.Directives)
				queryFields = append(queryFields, def)
			}
		}

		if e.genWhereInput && !ant.Skip.Is(SkipWhereInput) {
			def, err := e.buildWhereInput(node, gqlType, names.WhereInput)
			if err != nil {
				return err
			}
			if def != nil {
				s.AddTypes(def)
			}
		}

		if e.genMutations {
			defs, err := e.buildMutationInputs(node, ant, gqlType)
			if err != nil {
				return err
			}
			if len(defs) > 0 {
				s.AddTypes(defs...)
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

func (e *schemaGenerator) mayAddScalars(s *ast.Schema, def *ast.Definition) {
	var redeclareErr bool
	// If there is a config file but the schema there was not loaded.
	if e.cfg != nil && e.cfg.Schema == nil {
		// Do not fail in case of error.
		err := e.cfg.LoadSchema()
		redeclareErr = err != nil && strings.Contains(err.Error(), "Cannot redeclare type")
	}
	for _, f := range def.Fields {
		switch name := f.Type.Name(); name {
		case "Time", "Map", "Upload", "Any", "Int32", "Int64", "Uint", "Uint32", "Uint64":
			// Skip adding it if it was added before, or it exists in other schemas.
			if s.Types[name] == nil && e.externalType(name) {
				break
			}
			// In case of a declaration error generate builtin types only no external
			// schemas were found to allow users fix these failures.
			if !redeclareErr || len(e.cfg.SchemaFilename) == 1 && filepath.Clean(e.cfg.SchemaFilename[0]) == filepath.Clean(e.path) {
				s.AddTypes(&ast.Definition{
					Name:        name,
					Kind:        ast.Scalar,
					Description: fmt.Sprintf("The builtin %s type", name),
				})
			}
		}
	}
}

// externalType indicates if the given type name exists in another schema.
func (e *schemaGenerator) externalType(name string) bool {
	if e.cfg == nil || e.cfg.Schema == nil || e.cfg.Schema.Types[name] == nil {
		return false
	}
	def := e.cfg.Schema.Types[name]
	return def.Position != nil && def.Position.Src != nil && filepath.Clean(def.Position.Src.Name) != filepath.Clean(e.path)
}

func (e *schemaGenerator) buildType(t *gen.Type, ant *Annotation, gqlType, pkg string) (*ast.Definition, error) {
	def := &ast.Definition{
		Name:       gqlType,
		Kind:       ast.Object,
		Directives: e.buildDirectives(ant.Directives),
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

	fields := allFields(t)
	for _, f := range fields {
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(SkipType) || f.Sensitive() {
			continue
		}

		fieldDefs, err := e.fieldDefinitions(gqlType, f, ant)
		if err != nil {
			return nil, err
		}
		if fieldDefs != nil {
			def.Fields = append(def.Fields, fieldDefs...)
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
			return nil, fmt.Errorf("entgql: RelayConnection cannot be defined on Unique edge: %s.%s", t.Name, edge.Name)
		}

		fields, err := e.buildEdge(t, edge, ant)
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
		list = append(list, &ast.Directive{
			Name:      d.Name,
			Arguments: d.Arguments,
		})
	}
	return list
}

func (e *schemaGenerator) enumOrderByValues(t *gen.Type, gqlType string) (*ast.Definition, error) {
	terms, err := orderFields(t)
	if err != nil {
		return nil, err
	}
	enumValues := make(ast.EnumValueList, 0, len(terms))
	for _, f := range terms {
		enumValues = append(enumValues, &ast.EnumValueDefinition{
			Name: f.GQL,
		})
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

func (e *schemaGenerator) buildEdge(node *gen.Type, edge *gen.Edge, edgeAnt *Annotation) ([]*ast.FieldDefinition, error) {
	if edge.Type.HasCompositeID() {
		return nil, nil
	}
	gqlType, ant, err := gqlTypeFromNode(edge.Type)
	if err != nil || ant.Skip.Is(SkipType) {
		return nil, err
	}
	orderFields, err := orderFields(edge.Type)
	if err != nil {
		return nil, err
	}

	var (
		fields      []*ast.FieldDefinition
		mappings    = []string{camel(edge.Name)}
		goFieldName = templates.ToGo(edge.Name)
		structField = edge.StructField()
	)
	if len(edgeAnt.Mapping) > 0 {
		mappings = edgeAnt.Mapping
	}
	for _, name := range mappings {
		fieldDef := &ast.FieldDefinition{
			Name:        name,
			Description: edge.Comment(),
		}
		switch {
		case edge.Unique:
			fieldDef.Type = namedType(gqlType, edge.Optional)
		// Avoid error in case the RelayConnection is defined on the
		// `Through` edge, but the edge-schema is not a Relay connection.
		case edgeAnt.RelayConnection && edge.Type.IsEdgeSchema() && !ant.RelayConnection:
			fieldDef.Type = listNamedType(gqlType, edge.Optional)
		case edgeAnt.RelayConnection:
			if !e.relaySpec {
				return nil, ErrRelaySpecDisabled
			}
			if !ant.RelayConnection {
				return nil, fmt.Errorf("entgql.RelayConnection() must be set on entity %q in order to define %q.%q as Relay Connection", edge.Type.Name, node.Name, edge.Name)
			}

			fieldDef = paginationNames(gqlType).
				ConnectionField(name, len(orderFields) > 0, ant.MultiOrder,
					e.genWhereInput && !edgeAnt.Skip.Is(SkipWhereInput) && !ant.Skip.Is(SkipWhereInput),
				)
		default:
			fieldDef.Type = listNamedType(gqlType, edge.Optional)
		}

		fieldDef.Directives = e.buildDirectives(edgeAnt.Directives)
		if goFieldName != templates.ToGo(name) {
			fieldDef.Directives = append(fieldDef.Directives, goField(structField))
		}
		fields = append(fields, fieldDef)
	}

	return fields, nil
}

// buildWhereInput returns the <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *schemaGenerator) buildWhereInput(t *gen.Type, nodeGQLType, gqlType string) (*ast.Definition, error) {
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

	fields := allFields(t)
	for _, f := range fields {
		if t.IsEdgeSchema() && f.IsEdgeField() || !f.Type.Comparable() || f.Sensitive() {
			continue
		}
		ant, err := annotation(f.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip.Is(SkipWhereInput) {
			continue
		}
		for i, op := range f.Ops() {
			fd := e.fieldDefinitionOp(nodeGQLType, f, ant, op)
			if i == 0 {
				fd.Description = f.Name + " field predicates"
			}
			def.Fields = append(def.Fields, fd)
		}
	}

	if t.IsEdgeSchema() {
		return def, nil
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

func (e *schemaGenerator) buildMutationInputs(t *gen.Type, ant *Annotation, gqlType string) ([]*ast.Definition, error) {
	var defs []*ast.Definition

	for _, i := range ant.MutationInputs {
		if i.IsCreate && ant.Skip.Is(SkipMutationCreateInput) {
			continue
		}
		if !i.IsCreate && ant.Skip.Is(SkipMutationUpdateInput) {
			continue
		}

		desc := MutationDescriptor{Type: t, IsCreate: i.IsCreate}
		name, err := desc.Input()
		if err != nil {
			return nil, err
		}
		fields, err := desc.InputFields()
		if err != nil {
			return nil, err
		}
		edges, err := desc.InputEdges()
		if err != nil {
			return nil, err
		}

		def := &ast.Definition{
			Name:        name,
			Kind:        ast.InputObject,
			Description: i.Description,
		}

		if def.Description == "" {
			if i.IsCreate {
				def.Description = fmt.Sprintf("%s is used for create %s object.\nInput was generated by ent.", name, gqlType)
			} else {
				def.Description = fmt.Sprintf("%s is used for update %s object.\nInput was generated by ent.", name, gqlType)
			}
		}

		for _, f := range fields {
			ant, err := annotation(f.Annotations)
			if err != nil {
				return nil, err
			}
			scalar := e.mapScalar(gqlType, f.Field, ant, inputObjectFilter)
			if scalar == "" {
				return nil, fmt.Errorf("%s is not supported as input for %s", f.Name, def.Name)
			}
			def.Fields = append(def.Fields, &ast.FieldDefinition{
				Name:        camel(f.Name),
				Type:        namedType(scalar, f.Nullable),
				Description: f.Comment(),
			})
			if f.AppendOp {
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: "append" + f.StructField(),
					Type: namedType(scalar, true),
				})
			}
			if f.ClearOp {
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: "clear" + f.StructField(),
					Type: namedType("Boolean", true),
				})
			}
		}

		for _, e := range edges {
			switch {
			case e.Unique:
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: camel(e.Name) + "ID",
					Type: namedType("ID", !i.IsCreate || e.Optional),
				})
			case i.IsCreate:
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: camel(singular(e.Name)) + "IDs",
					Type: namedType("[ID!]", e.Optional),
				})
			default:
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: "add" + pascal(singular(e.Name)) + "IDs",
					Type: namedType("[ID!]", true),
				}, &ast.FieldDefinition{
					Name: "remove" + pascal(singular(e.Name)) + "IDs",
					Type: namedType("[ID!]", true),
				})
			}
			if !i.IsCreate && e.Optional {
				def.Fields = append(def.Fields, &ast.FieldDefinition{
					Name: camel(snake(e.MutationClear())),
					Type: namedType("Boolean", true),
				})
			}
		}
		defs = append(defs, def)
	}

	return defs, nil
}

func (e *schemaGenerator) fieldDefinitions(gqlType string, f *gen.Field, ant *Annotation) ([]*ast.FieldDefinition, error) {
	ft, err := e.typeFromField(gqlType, f, ant)
	if err != nil {
		return nil, fmt.Errorf("field(%s): %w", f.Name, err)
	}
	var (
		fields      []*ast.FieldDefinition
		goFieldName = templates.ToGo(f.Name)
		structField = f.StructField()
	)
	mapping, err := fieldMapping(f)
	if err != nil {
		return nil, err
	}
	for _, name := range mapping {
		def := &ast.FieldDefinition{
			Name:        name,
			Type:        ft,
			Description: f.Comment(),
			Directives:  e.buildDirectives(ant.Directives),
		}
		// We check the field name with gqlgen's naming convention.
		// To avoid unnecessary @goField directives
		if goFieldName != templates.ToGo(name) {
			def.Directives = append(def.Directives, goField(structField))
		}
		fields = append(fields, def)
	}
	return fields, nil
}

// fieldMapping returns the GraphQL names mapping of a field.
func fieldMapping(f *gen.Field) ([]string, error) {
	ant, err := annotation(f.Annotations)
	if err != nil || ant.Skip.Is(SkipType) || f.Sensitive() {
		return nil, err
	}
	if len(ant.Mapping) > 0 {
		return ant.Mapping, nil
	}
	return []string{camel(f.Name)}, nil
}

func (e *schemaGenerator) fieldDefinitionOp(gqlType string, f *gen.Field, ant *Annotation, op gen.Op) *ast.FieldDefinition {
	def := &ast.FieldDefinition{
		Name: camel(f.Name + "_" + op.Name()),
	}
	if op == gen.EQ {
		def.Name = camel(f.Name)
	}

	if e.scalarFunc != nil {
		if t := e.scalarFunc(f, op); t != "" {
			def.Type = namedType(t, true)
			return def
		}
	}

	switch {
	case op.Niladic():
		def.Type = namedType("Boolean", true)
	case op.Variadic():
		def.Type = listNamedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	default:
		def.Type = namedType(e.mapScalar(gqlType, f, ant, inputObjectFilter), true)
	}
	return def
}

func (e *schemaGenerator) typeFromField(gqlType string, f *gen.Field, ant *Annotation) (*ast.Type, error) {
	if scalar := e.mapScalar(gqlType, f, ant, nonInputObjectFilter); scalar != "" {
		return namedType(scalar, f.Optional), nil
	}

	switch t := f.Type.Type; {
	case t == field.TypeJSON:
		return nil, fmt.Errorf("entgql: json type not implemented without setting an entgql.Type() annotation")
	case t == field.TypeOther:
		return nil, fmt.Errorf("entgql: other type must have typed defined")
	default:
		return nil, fmt.Errorf("entgql: unexpected type: %s", t.String())
	}
}

// mapScalar provides maps an ent.Schema type into GraphQL scalar type.
func (e *schemaGenerator) mapScalar(gqlType string, f *gen.Field, ant *Annotation, typeFilter func(string) bool) string {
	if ant != nil && ant.Type != "" {
		return ant.Type
	}
	scalar := f.Type.String()
	switch t := f.Type.Type; {
	case f.Name == "id":
		return "ID"
	case f.IsEdgeField():
		scalar = "ID"
	case t.Float():
		scalar = "Float"
	case t.Integer():
		scalar = "Int"
	case t == field.TypeString:
		scalar = "String"
	case t == field.TypeBool:
		scalar = "Boolean"
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		if typ, ok := e.hasMapping(f, typeFilter); ok {
			scalar = typ
		} else {
			scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		}
		if f.IsEnum() {
			// Use the GQL type as enum prefix. e.g. Todo.status
			// will generate an enum named "TodoStatus".
			scalar = gqlType + scalar
		}
		if f.Type.RType != nil && f.Type.RType.Name == "" {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				if strings.HasPrefix(f.Type.RType.Ident, "[]*") {
					scalar = "[" + scalar + "]"
				} else {
					scalar = "[" + scalar + "!]"
				}
			}
		}
	case t == field.TypeJSON:
		scalar = ""
		if f.Type.RType != nil {
			switch f.Type.RType.Kind {
			case reflect.Slice, reflect.Array:
				switch f.Type.RType.Ident {
				case "[]float64":
					scalar = "[Float!]"
				case "[]int":
					scalar = "[Int!]"
				case "[]string":
					scalar = "[String!]"
				}
			case reflect.Map:
				if f.Type.RType.Ident == "map[string]interface {}" {
					scalar = "Map"
				}
			}
		}
	}
	return scalar
}

// hasMapping reports if the gqlgen.yml has custom mapping for
// the given field type and returns its GraphQL name if exists.
func (e *schemaGenerator) hasMapping(f *gen.Field, typeFilter func(string) bool) (string, bool) {
	if e.cfg == nil {
		return "", false
	}

	// The string representation uses shortened package
	// names, and we override it for custom Go types.
	ident := f.Type.String()
	if idx := strings.IndexByte(ident, '.'); idx != -1 && f.HasGoType() && f.Type.PkgPath != "" {
		ident = f.Type.PkgPath + ident[idx:]
	}

	var gqlNames []string
	for t, v := range e.cfg.Models {
		for _, m := range v.Model {
			// A mapping was found from GraphQL name to field type.
			if strings.HasSuffix(m, ident) {
				gqlNames = append(gqlNames, t)
			}
		}
	}
	if count := len(gqlNames); count == 1 {
		return gqlNames[0], true
	} else if count > 1 && typeFilter != nil {
		for _, t := range gqlNames {
			if typeFilter(t) {
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

func allFields(t *gen.Type) []*gen.Field {
	if t.ID == nil {
		return t.Fields
	}

	// NOTE: always keep the ID field at the beginning of the list
	return append([]*gen.Field{t.ID}, t.Fields...)
}

func fieldGoType(f *gen.Field, pkg string) string {
	if f.HasGoType() {
		return entGoType(f.Type.RType.Name, f.Type.RType.PkgPath)
	}
	return fmt.Sprintf("%s/%s", pkg, f.Type.Ident)
}

func entGoType(name, pkg string) string {
	return fmt.Sprintf("%s.%s", pkg, name)
}

func builtinTypes() []*ast.Definition {
	return []*ast.Definition{
		{
			Name:        OrderDirectionEnum,
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
