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
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/graphql-go/graphql"
	"github.com/graphql-go/graphql/language/ast"
	"github.com/graphql-go/graphql/language/kinds"
	"github.com/graphql-go/graphql/language/parser"
	"github.com/graphql-go/graphql/language/printer"
	"github.com/graphql-go/graphql/language/source"
	"github.com/graphql-go/graphql/language/visitor"
)

type (
	// Extension implements the entc.Extension for providing GraphQL integration.
	Extension struct {
		entc.DefaultExtension
		path         string
		doc          *ast.Document
		cfg          *config.Config
		hooks        []gen.Hook
		templates    []*gen.Template
		scalarFunc   func(*gen.Field, gen.Op) string
		naming       func(string) string
		orderBy      bool
		whereFilters bool
		relaySpec    bool
		relayQuery   func(string) string
	}

	// ExtensionOption allows for managing the Extension configuration
	// using functional options.
	ExtensionOption func(*Extension) error
)

// WithSchemaPath sets the filepath to the GraphQL schema to write the
// generated Ent types. If the file does not exist, it will generate a
// new schema. Please note, that your gqlgen.yml config file should be
// updated as follows to support multiple schema files:
//
//	schema:
//	 - schema.graphql // existing schema.
//	 - ent.graphql	  // generated schema.
//
func WithSchemaPath(path string) ExtensionOption {
	return func(ex *Extension) error {
		buf, err := ioutil.ReadFile(path)
		if err != nil && !os.IsNotExist(err) {
			return fmt.Errorf("reading graphql schema %q: %w", path, err)
		}
		ex.doc, err = parser.Parse(parser.ParseParams{
			Source: &source.Source{
				Body: buf,
				Name: filepath.Base(path),
			},
		})
		if err != nil {
			return fmt.Errorf("parsing graphql schema %q: %w", path, err)
		}
		ex.path = path
		ex.hooks = append(ex.hooks, ex.genSchema())
		return nil
	}
}

// GQLConfigAnnotation is the annotation key/name that holds gqlgen
// configuration if it was provided by the `WithConfigPath` option.
const GQLConfigAnnotation = "GQLConfig"

// WithConfigPath sets the filepath to gqlgen.yml configuration file
// and injects its parsed version to the global annotations.
//
// Note that, enabling this option is recommended as it improves the
// GraphQL integration,
func WithConfigPath(path string) ExtensionOption {
	return func(ex *Extension) (err error) {
		cwd, err := os.Getwd()
		if err != nil {
			return fmt.Errorf("unable to get working directory: %w", err)
		}
		if err := os.Chdir(filepath.Dir(path)); err != nil {
			return fmt.Errorf("unable to enter config dir: %w", err)
		}
		defer func() {
			if cerr := os.Chdir(cwd); cerr != nil {
				err = fmt.Errorf("unable to restore working directory: %w", cerr)
			}
		}()
		cfg, err := config.LoadConfig(filepath.Base(path))
		if err != nil {
			return err
		}
		if cfg.Schema == nil {
			if err := cfg.LoadSchema(); err != nil {
				return err
			}
		}
		ex.cfg = cfg
		ex.hooks = append(ex.hooks, func(next gen.Generator) gen.Generator {
			return gen.GenerateFunc(func(g *gen.Graph) error {
				if g.Annotations == nil {
					g.Annotations = gen.Annotations{}
				}
				g.Annotations[GQLConfigAnnotation] = cfg
				return next.Generate(g)
			})
		})
		return nil
	}
}

// WithTemplates overrides the default templates (entgql.AllTemplates)
// with specific templates.
func WithTemplates(templates ...*gen.Template) ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = templates
		return nil
	}
}

// WithQueryFilters adds the template to generate a basic Filter function to each
// query type in the ent schema.
// This is for GraphQL queries not using pagination/relay
func WithQueryFilters() ExtensionOption {
	return func(ex *Extension) error {
		ex.templates = append(ex.templates, FilterTemplate)
		return nil
	}
}

// WithNaming takes either "snake" or "camel" and sets that as the naming
// convention when generating the GraphQL schema
func WithNaming(s string) ExtensionOption {
	return func(ex *Extension) error {
		switch s {
		case "snake":
			ex.naming = snake
		case "camel":
			ex.naming = camel
		default:
			return fmt.Errorf("unknown naming convention \"%s\"", s)
		}
		return nil
	}
}

// WithCustomRelaySpec can be used to customize the query fields that return the
// relay specificaiton connections.
//
// A function taking the node name as input should return a string for the custom
// query name that you want to provide. Note that the naming convention will be
// applied to this as well to make sure it follows your naming convention
func WithCustomRelaySpec(b bool, queryName func(string) string) ExtensionOption {
	return func(ex *Extension) error {
		ex.relaySpec = b
		ex.relayQuery = queryName
		return nil
	}
}

func WithRelaySpec(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.relaySpec = b
		return nil
	}
}

func WithOrderBy(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.orderBy = b
		return nil
	}
}

// WithWhereFilters configures the extension to either add or
// remove the WhereTemplate from the code generation templates.
//
// The WhereTemplate generates GraphQL filters to all types in the ent/schema.
func WithWhereFilters(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.whereFilters = b
		i, exists := ex.whereExists()
		if b && !exists {
			ex.templates = append(ex.templates, WhereTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithMapScalarFunc allows users to provides a custom function that
// maps an ent.Field (*gen.Field) into its GraphQL scalar type. If the
// function returns an empty string, the extension fallbacks to the its
// default mapping.
//
//	ex, err := entgql.NewExtension(
//		entgql.WithMapScalarFunc(func(f *gen.Field, op gen.Op) string {
//			if t, ok := knowType(f, op); ok {
//				return t
//			}
//			// Fallback to the default mapping.
//			return ""
//		}),
//	)
//
func WithMapScalarFunc(scalarFunc func(*gen.Field, gen.Op) string) ExtensionOption {
	return func(ex *Extension) error {
		ex.scalarFunc = scalarFunc
		return nil
	}
}

// NewExtension creates a new extension with the given configuration.
//
//	ex, err := entgql.NewExtension(
//		entgql.WithWhereFilters(true),
//		entgql.WithSchemaPath("../schema.graphql"),
//		entgql.WithNaming("camel"),
//	)
//
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates: AllTemplates,
		naming:    camel,
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	return ex, nil
}

// Templates of the extension.
func (e *Extension) Templates() []*gen.Template {
	return e.templates
}

// Hooks of the extension.
func (e *Extension) Hooks() []gen.Hook {
	return e.hooks
}

// mapScalar provides maps an ent.Schema type into GraphQL scalar type.
// In order to override this function, use the WithMapScalarFunc option.
func (e *Extension) mapScalar(f *gen.Field, op gen.Op) string {
	if e.scalarFunc != nil {
		if t := e.scalarFunc(f, op); t != "" {
			return t
		}
	}
	scalar := f.Type.String()
	switch t := f.Type.Type; {
	case op.Niladic() || t == field.TypeBool:
		scalar = graphql.Boolean.Name()
	case f.IsEdgeField():
		scalar = graphql.ID.Name()
	case t.Integer():
		scalar = graphql.Int.Name()
	case t.Float():
		scalar = graphql.Float.Name()
	case t == field.TypeString:
		scalar = graphql.String.Name()
	case strings.ContainsRune(scalar, '.'): // Time, Enum or Other.
		if typ, ok := e.hasMapping(f); ok {
			scalar = typ
		} else {
			scalar = scalar[strings.LastIndexByte(scalar, '.')+1:]
		}
	}
	return scalar
}

// genSchema generates the graphql schema
func (e *Extension) genSchema() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		var (
			queries    = make(map[string]*ast.FieldDefinition)
			interfaces = make(map[string]*ast.InterfaceDefinition)
			scalars    = make(map[string]*ast.ScalarDefinition)
			objects    = make(map[string]*ast.ObjectDefinition)
			enums      = make(map[string]*ast.EnumDefinition)
			inputs     = make(map[string]*ast.InputObjectDefinition)
		)
		return gen.GenerateFunc(func(g *gen.Graph) error {
			nodes, err := filterNodes(g.Nodes)
			if err != nil {
				return err
			}
			if err := next.Generate(g); err != nil {
				return err
			}

			addTimeScalarDefinition(scalars)
			addNodeDefinition(interfaces)
			// Iterate through the nodes and create the base objects, enums and queries
			for _, node := range nodes {
				e.addEnumDefinitions(enums, node)
				e.addObjectDefinition(objects, node)

				var hasOrderByType bool
				// If orderBy has been enabled, create the types for ordering
				if e.orderBy {
					hasOrderByType = e.addOrderByDefinitions(enums, inputs, node)
				}

				if e.whereFilters {
					// Create the where input types
					err := e.addWhereType(inputs, node)
					if err != nil {
						return err
					}
				}

				if e.relaySpec {
					e.addRelayObjectDefinitions(objects, node)
					e.addQueryDefinition(queries, node, true, hasOrderByType)
				}
				// If either the relay spec is not enabled, or it is enabled but
				// with a custom name, then generate also the "default" query
				// for this node
				if !e.relaySpec || e.relayQuery != nil {
					e.addQueryDefinition(queries, node, false, hasOrderByType)
				}
			}
			if e.orderBy {
				// Create the global OrderDirection enum
				e.addGlobalOrderByDefinition(enums)
			}
			if e.relaySpec {
				// Create the global relay spec objects
				addCursorScalarDefinition(scalars)
				addPageInfoDefinition(objects)
			}
			return e.updateSchema(scalars, interfaces, queries, objects, enums, inputs)
		})
	}
}

// addObjectDefinition creates a GraphQL type for the given node, and adds it to the map of objects
func (e *Extension) addObjectDefinition(objects map[string]*ast.ObjectDefinition, node *gen.Type) {
	obj := ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name: newASTName(node.Name),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: fmt.Sprintf("%s represents the node %s in the ent schema.\nGenerated by ent.", node.Name, node.Name),
		}),
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("id"),
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(graphql.ID.String()),
				}),
			}),
		},
		Interfaces: []*ast.Named{newASTNamed("Node")},
	})

	for _, f := range node.Fields {
		obj.Fields = append(obj.Fields, ast.NewFieldDefinition(
			&ast.FieldDefinition{
				Name: newASTName(f.Name),
				Type: newASTNamed(e.mapScalar(f, gen.EQ)),
			}),
		)
	}

	for _, edge := range node.Edges {
		var edgeType ast.Type
		edgeType = newASTNamed(edge.Type.Name)
		if !edge.Optional {
			edgeType = ast.NewNonNull(&ast.NonNull{Type: edgeType})
		}
		if !edge.Unique {
			edgeType = ast.NewList(&ast.List{Type: edgeType})
		}
		obj.Fields = append(obj.Fields, ast.NewFieldDefinition(
			&ast.FieldDefinition{
				Name: newASTName(edge.Name),
				Type: edgeType,
			},
		))
	}

	objects[node.Name] = obj
}

// addEnumDefinitions creates GraphQL enums defined by the given node, and adds them to the map of enums
func (e *Extension) addEnumDefinitions(enums map[string]*ast.EnumDefinition, node *gen.Type) {
	for _, ef := range node.EnumFields() {
		enum := ast.NewEnumDefinition(&ast.EnumDefinition{
			Name: newASTName(e.mapScalar(ef, gen.EQ)),
		})
		for _, ev := range ef.Enums {
			enum.Values = append(enum.Values, ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Name: newASTName(ev.Value),
			}))
		}
		enums[e.mapScalar(ef, gen.EQ)] = enum
	}
}

// addGlobalOrderByDefinition creates a GraphQL enum for the OrderDirection
func (e *Extension) addGlobalOrderByDefinition(enums map[string]*ast.EnumDefinition) {
	enums["OrderDirection"] = ast.NewEnumDefinition(&ast.EnumDefinition{
		Name: newASTName("OrderDirection"),
		Values: []*ast.EnumValueDefinition{
			ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Name: newASTName("ASC"),
			}),
			ast.NewEnumValueDefinition(&ast.EnumValueDefinition{
				Name: newASTName("DESC"),
			}),
		},
	})
}

// addOrderByDefinitions creates GraphQL inputs for ordering for the given node
func (e *Extension) addOrderByDefinitions(enums map[string]*ast.EnumDefinition, inputs map[string]*ast.InputObjectDefinition, node *gen.Type) bool {
	var hasOrderByType bool
	orderByField := node.Name + "OrderField"
	// Create the order by enum which we will populate with those fields that have
	// been annotated with entgql.OrderField("...")
	orderBy := ast.NewEnumDefinition(&ast.EnumDefinition{
		Name: newASTName(orderByField),
	})
	for _, f := range node.Fields {
		// Check if graphql ordering has been applied to this field through an annotation.
		// If it has, add it to the order by enum for this node.
		var ant Annotation
		if i, ok := f.Annotations[ant.Name()]; ok && ant.Decode(i) == nil && ant.OrderField != "" {
			orderBy.Values = append(orderBy.Values, ast.NewEnumValueDefinition(
				&ast.EnumValueDefinition{
					Name: newASTName(f.Name),
				},
			))
		}
	}
	// Add the order by enum if there is at least one value to
	// order by, otherwise we skip it from the schema
	if len(orderBy.Values) > 0 {
		hasOrderByType = true
		enums[orderByField] = orderBy
		inputs[orderType(node.Name)] = ast.NewInputObjectDefinition(
			&ast.InputObjectDefinition{
				Name: newASTName(orderType(node.Name)),
				Fields: []*ast.InputValueDefinition{
					ast.NewInputValueDefinition(&ast.InputValueDefinition{
						Name: newASTName("direction"),
						Type: newASTNamed("OrderDirection!"),
					}),
					ast.NewInputValueDefinition(&ast.InputValueDefinition{
						Name: newASTName("field"),
						Type: newASTNamed(orderByField),
					}),
				},
			},
		)
	}
	return hasOrderByType
}

// addOrderByDefinitions creates a GraphQL query field for the given node, adding
// the relevanant arguments based on whether the relay spec should be followed,
// ordering is enabled, and/or the where filters are enabled
func (e *Extension) addQueryDefinition(queries map[string]*ast.FieldDefinition, node *gen.Type, relaySpec bool, orderBy bool) {
	queryName := node.Table()
	if relaySpec && e.relayQuery != nil {
		queryName = e.naming(e.relayQuery(queryName))
	}
	query := ast.NewFieldDefinition(&ast.FieldDefinition{
		Name: newASTName(queryName),
		Type: ast.NewList(&ast.List{
			Type: newASTNamed(node.Name),
		}),
		Arguments: []*ast.InputValueDefinition{
			ast.NewInputValueDefinition(&ast.InputValueDefinition{
				Name: newASTName("first"),
				Type: newASTNamed(graphql.Int.String()),
			}),
			ast.NewInputValueDefinition(&ast.InputValueDefinition{
				Name: newASTName("last"),
				Type: newASTNamed(graphql.Int.String()),
			}),
		},
	})

	if relaySpec {
		query.Type = ast.NewList(&ast.List{
			Type: newASTNamed(node.Name + "Connection"),
		})
		query.Arguments = append(query.Arguments,
			ast.NewInputValueDefinition(&ast.InputValueDefinition{
				Name: newASTName("before"),
				Type: newASTNamed("Cursor"),
			}),
			ast.NewInputValueDefinition(&ast.InputValueDefinition{
				Name: newASTName("after"),
				Type: newASTNamed("Cursor"),
			}))
	}

	if orderBy {
		query.Arguments = append(query.Arguments, ast.NewInputValueDefinition(
			&ast.InputValueDefinition{
				Name: newASTName(e.naming("order_by")),
				Type: newASTNamed(orderType(node.Name)),
			},
		))
	}

	if e.whereFilters {
		query.Arguments = append(query.Arguments, ast.NewInputValueDefinition(
			&ast.InputValueDefinition{
				Name: newASTName("where"),
				Type: newASTNamed(node.Name + "WhereInput"),
			},
		))
	}

	queries[queryName] = query
}

// addRelayObjectDefinitions creates the types (Connection and Edge) for the given node
// to follow the relay spec
func (e *Extension) addRelayObjectDefinitions(objects map[string]*ast.ObjectDefinition, node *gen.Type) {
	conn := ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name: newASTName(node.Name + "Connection"),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: fmt.Sprintf("%s supports the relay connection specification for node %s in the ent schema.\nGenerated by ent.", node.Name+"Connection", node.Name),
		}),
	})
	conn.Fields = []*ast.FieldDefinition{
		ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: newASTName(e.naming("total_count")),
			Type: newASTNamed(graphql.Int.String() + "!"),
		}),
		ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: newASTName(e.naming("page_info")),
			Type: newASTNamed("PageInfo!"),
		}),
		ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: newASTName("edges"),
			Type: ast.NewList(&ast.List{
				Type: newASTNamed(node.Name + "Edge"),
			}),
		}),
	}
	objects[node.Name+"Connection"] = conn

	edge := ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name: newASTName(node.Name + "Edge"),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: fmt.Sprintf("%s supports the relay edge specification for node %s in the ent schema.\nGenerated by ent.", node.Name+"Connection", node.Name),
		}),
	})
	edge.Fields = []*ast.FieldDefinition{
		ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: newASTName("node"),
			Type: newASTNamed(node.Name),
		}),
		ast.NewFieldDefinition(&ast.FieldDefinition{
			Name: newASTName("cursor"),
			Type: newASTNamed("Cursor!"),
		}),
	}
	objects[node.Name+"Edge"] = edge
}

// hasMapping reports if the gqlgen.yml has custom mapping for
// the given field type and returns its GraphQL name if exists.
func (e *Extension) hasMapping(f *gen.Field) (string, bool) {
	var ant Annotation
	// If the field was defined with a `entgql.Type` option (e.g. `entgql.Type("Boolean")`).
	if i, ok := f.Annotations[ant.Name()]; ok && ant.Decode(i) == nil && ant.Type != "" {
		return ant.Type, true
	}
	if e.cfg == nil {
		return "", false
	}
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

// whereExists reports if the WhereTemplate exists
// in the template list and returns its index.
func (e *Extension) whereExists() (int, bool) {
	for i := range e.templates {
		if e.templates[i] == WhereTemplate {
			return i, true
		}
	}
	return -1, false
}

// updateSchema commits the changes to the GraphQL schema file.
func (e *Extension) updateSchema(
	scalars map[string]*ast.ScalarDefinition, interfaces map[string]*ast.InterfaceDefinition,
	queries map[string]*ast.FieldDefinition, objects map[string]*ast.ObjectDefinition,
	enums map[string]*ast.EnumDefinition, inputs map[string]*ast.InputObjectDefinition) error {
	var queryUpdated bool
	visitor.Visit(e.doc, &visitor.VisitorOptions{
		LeaveKindMap: map[string]visitor.VisitFunc{
			kinds.ScalarDefinition: func(p visitor.VisitFuncParams) (string, interface{}) {
				if node, ok := p.Node.(*ast.ScalarDefinition); ok && scalars[node.Name.Value] != nil {
					scalar := scalars[node.Name.Value]
					delete(scalars, node.Name.Value)
					return visitor.ActionUpdate, scalar
				}
				return visitor.ActionNoChange, nil
			},
			kinds.InterfaceDefinition: func(p visitor.VisitFuncParams) (string, interface{}) {
				if node, ok := p.Node.(*ast.InterfaceDefinition); ok && interfaces[node.Name.Value] != nil {
					inter := interfaces[node.Name.Value]
					delete(interfaces, node.Name.Value)
					return visitor.ActionUpdate, inter
				}
				return visitor.ActionNoChange, nil
			},
			kinds.EnumDefinition: func(p visitor.VisitFuncParams) (string, interface{}) {
				if node, ok := p.Node.(*ast.EnumDefinition); ok && enums[node.Name.Value] != nil {
					enum := enums[node.Name.Value]
					delete(enums, node.Name.Value)
					return visitor.ActionUpdate, enum
				}
				return visitor.ActionNoChange, nil
			},
			kinds.ObjectDefinition: func(p visitor.VisitFuncParams) (string, interface{}) {
				if node, ok := p.Node.(*ast.ObjectDefinition); ok {
					if node.Name.Value == "Query" {
						// Mark that we have updated the query
						queryUpdated = true
						// Handle the query object
						updateQueryFields(node, queries)
						return visitor.ActionUpdate, node
					} else if objects[node.Name.Value] != nil {
						field := objects[node.Name.Value]
						delete(objects, node.Name.Value)
						return visitor.ActionUpdate, field
					}
				}
				return visitor.ActionNoChange, nil
			},
			kinds.InputObjectDefinition: func(p visitor.VisitFuncParams) (string, interface{}) {
				// If the input object was found in the schema, we update its definition.
				if node, ok := p.Node.(*ast.InputObjectDefinition); ok && inputs[node.Name.Value] != nil {
					input := inputs[node.Name.Value]
					delete(inputs, node.Name.Value)
					return visitor.ActionUpdate, input
				}
				return visitor.ActionNoChange, nil
			},
		},
	}, nil)
	if !queryUpdated {
		query := ast.NewObjectDefinition(&ast.ObjectDefinition{
			Name:        newASTName("Query"),
			Description: ast.NewStringValue(&ast.StringValue{Value: "Query generated by ent."}),
		})
		updateQueryFields(query, queries)
		e.doc.Definitions = append(e.doc.Definitions, query)
	}
	// Sorting the input types is not needed, because in the next iteration
	// the hook updates the generated types without changing their position.
	for _, scalar := range scalars {
		e.doc.Definitions = append(e.doc.Definitions, scalar)
	}
	for _, inter := range interfaces {
		e.doc.Definitions = append(e.doc.Definitions, inter)
	}
	for _, enum := range enums {
		e.doc.Definitions = append(e.doc.Definitions, enum)
	}
	for _, field := range objects {
		e.doc.Definitions = append(e.doc.Definitions, field)
	}
	for _, input := range inputs {
		e.doc.Definitions = append(e.doc.Definitions, input)
	}
	return ioutil.WriteFile(e.path, []byte(printer.Print(e.doc).(string)), 0644)
}

// addWhereType returns the a <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *Extension) addWhereType(inputs map[string]*ast.InputObjectDefinition, t *gen.Type) error {
	var (
		name  = t.Name + "WhereInput"
		input = ast.NewInputObjectDefinition(&ast.InputObjectDefinition{
			Name: newASTName(name),
			Description: ast.NewStringValue(&ast.StringValue{
				Value: fmt.Sprintf("%s is used for filtering %s objects.\nInput was generated by ent.", name, t.Name),
			}),
			Fields: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name: newASTName("not"),
					Type: newASTNamed(name),
				}),
			},
		})
	)
	for _, op := range []string{"and", "or"} {
		input.Fields = append(input.Fields, ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: newASTName(op),
			Type: ast.NewList(&ast.List{
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(name),
				}),
			}),
		}))
	}
	fields, err := filterFields(t.Fields)
	if err != nil {
		return err
	}
	for _, f := range fields {
		if !f.Type.Comparable() {
			continue
		}
		for i, op := range f.Ops() {
			fd := e.fieldDefinition(f, op)
			if i == 0 {
				fd.Description = ast.NewStringValue(&ast.StringValue{
					Value: f.Name + " field predicates",
				})
			}
			input.Fields = append(input.Fields, fd)
		}
	}
	edges, err := filterEdges(t.Edges)
	if err != nil {
		return err
	}
	for _, edge := range edges {
		input.Fields = append(input.Fields, ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: newASTName(e.naming("has_" + edge.Name)),
			Type: newASTNamed("Boolean"),
			Description: ast.NewStringValue(&ast.StringValue{
				Value: edge.Name + " edge predicates",
			}),
		}), ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: newASTName(e.naming("has_" + edge.Name + "_with")),
			Type: ast.NewList(&ast.List{
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(edge.Type.Name + "WhereInput"),
				}),
			}),
		}))
	}
	inputs[name] = input
	return nil
}

func (e *Extension) fieldDefinition(f *gen.Field, op gen.Op) *ast.InputValueDefinition {
	name := e.naming(f.Name + "_" + op.Name())
	if op == gen.EQ {
		name = e.naming(f.Name)
	}
	fieldType := e.mapScalar(f, op)
	def := ast.NewInputValueDefinition(&ast.InputValueDefinition{
		Name: newASTName(name),
		Type: newASTNamed(fieldType),
	})
	if op.Variadic() {
		def.Type = ast.NewList(&ast.List{
			Type: ast.NewNonNull(&ast.NonNull{
				Type: def.Type,
			}),
		})
	}
	return def
}

// updateQueryFields takes the GraphQL Query object (type) definition and a list
// of queries (fields) to add to the Query object.
// It updates those that exists and adds the remainder to the Query object
func updateQueryFields(query *ast.ObjectDefinition, queries map[string]*ast.FieldDefinition) {
	for idx := range query.Fields {
		field := query.Fields[idx]
		if qf, ok := queries[field.Name.Value]; ok {
			// Overwrite the query field
			query.Fields[idx] = qf
			delete(queries, field.Name.Value)
		}
	}
	for _, field := range queries {
		query.Fields = append(query.Fields, field)
	}
}

func addTimeScalarDefinition(scalars map[string]*ast.ScalarDefinition) {
	scalars["Time"] = ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: newASTName("Time"),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: "Maps a Time GraphQL scalar to a Go time.Time struct.\nGenerated by ent.",
		}),
	})
}

func addCursorScalarDefinition(scalars map[string]*ast.ScalarDefinition) {
	scalars["Cursor"] = ast.NewScalarDefinition(&ast.ScalarDefinition{
		Name: newASTName("Cursor"),
		Description: ast.NewStringValue(&ast.StringValue{
			Value: "Define a Relay Cursor type:\nhttps://relay.dev/graphql/connections.htm#sec-Cursor\nGenerated by ent.",
		}),
	})
}

func addNodeDefinition(interfaces map[string]*ast.InterfaceDefinition) {
	interfaces["Node"] = ast.NewInterfaceDefinition(&ast.InterfaceDefinition{
		Name: newASTName("Node"),
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("id"),
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(graphql.ID.String()),
				}),
			}),
		},
	})
}

func addPageInfoDefinition(objects map[string]*ast.ObjectDefinition) {
	objects["PageInfo"] = ast.NewObjectDefinition(&ast.ObjectDefinition{
		Name: newASTName("PageInfo"),
		Fields: []*ast.FieldDefinition{
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("hasNextPage"),
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(graphql.Boolean.String()),
				}),
			}),
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("hasPreviousPage"),
				Type: ast.NewNonNull(&ast.NonNull{
					Type: newASTNamed(graphql.Boolean.String()),
				}),
			}),
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("startCursor"),
				Type: newASTNamed("Cursor"),
			}),
			ast.NewFieldDefinition(&ast.FieldDefinition{
				Name: newASTName("endCursor"),
				Type: newASTNamed("Cursor"),
			}),
		},
	})
}

func orderType(name string) string {
	return name + "Order"
}

func newASTName(name string) *ast.Name {
	return ast.NewName(&ast.Name{Value: name})
}

func newASTNamed(name string) *ast.Named {
	return ast.NewNamed(&ast.Named{
		Name: newASTName(name),
	})
}

var (
	_     entc.Extension = (*Extension)(nil)
	camel                = gen.Funcs["camel"].(func(string) string)
	snake                = gen.Funcs["snake"].(func(string) string)
)
