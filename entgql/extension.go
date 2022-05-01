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
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	// Extension implements the entc.Extension for providing GraphQL integration.
	Extension struct {
		entc.DefaultExtension
		path       string
		cfg        *config.Config
		hooks      []gen.Hook
		templates  []*gen.Template
		scalarFunc func(*gen.Field, gen.Op) string

		schema        *ast.Schema
		schemaHooks   []SchemaHook
		models        map[string]string
		genSchema     bool
		genWhereInput bool
	}

	// ExtensionOption allows for managing the Extension configuration
	// using functional options.
	ExtensionOption func(*Extension) error

	// SchemaHook is the hook that run after the GQL schema generation.
	SchemaHook func(*gen.Graph, *ast.Schema) error
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
		ex.path = path
		return nil
	}
}

// WithSchemaHook allows users to provide a list of hooks
// to run after the GQL schema generation.
func WithSchemaHook(hooks ...SchemaHook) ExtensionOption {
	return func(ex *Extension) error {
		ex.schemaHooks = append(ex.schemaHooks, hooks...)
		return nil
	}
}

// WithConfigPath sets the filepath to gqlgen.yml configuration file
// and injects its parsed version to the global annotations.
//
// Note that, enabling this option is recommended as it improves the
// GraphQL integration,
func WithConfigPath(path string, gqlgenOptions ...api.Option) ExtensionOption {
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
		ex.cfg = cfg
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

// WithWhereFilters configures the extension to either add or
// remove the WhereTemplate from the code generation templates.
//
// The WhereTemplate generates GraphQL filters to all types in the ent/schema.
func WithWhereFilters(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.genWhereInput = b
		i, exists := ex.whereExists()
		if b && !exists {
			ex.templates = append(ex.templates, WhereTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithSchemaGenerator add a hook for generate GQL schema
func WithSchemaGenerator() ExtensionOption {
	return func(e *Extension) error {
		e.genSchema = true
		return nil
	}
}

// WithMapScalarFunc allows users to provide a custom function that
// maps an ent.Field (*gen.Field) into its GraphQL scalar type. If the
// function returns an empty string, the extension fallbacks to its
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
//	)
//
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates: AllTemplates,
		schema: &ast.Schema{
			Directives: map[string]*ast.DirectiveDefinition{},
		},
	}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	ex.hooks = append(ex.hooks, ex.genSchemaHook(), removeOldAssets)
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
func (e *Extension) hasMapping(f *gen.Field) (string, bool) {
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

// genSchema returns a new hook for generating
// the GraphQL schema from the graph.
func (e *Extension) genSchemaHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			if err := next.Generate(g); err != nil {
				return err
			}

			if !e.genSchema && !e.genWhereInput {
				return nil
			}
			genSchema, err := newSchemaGenerator(g)
			if err != nil {
				return err
			}

			genSchema.genWhereInput = e.genWhereInput
			if e.genSchema {
				if err = genSchema.buildSchema(e.schema); err != nil {
					return err
				}
				if e.models, err = genSchema.genModels(); err != nil {
					return err
				}
			}

			if e.genWhereInput {
				nodes, err := filterNodes(g.Nodes, SkipWhereInput)
				if err != nil {
					return err
				}
				for _, node := range nodes {
					input, err := e.whereType(node)
					if err != nil {
						return err
					}
					e.schema.AddTypes(input)
				}
			}

			for _, h := range e.schemaHooks {
				if err = h(g, e.schema); err != nil {
					return err
				}
			}
			if e.path == "" {
				return nil
			}
			return ioutil.WriteFile(e.path, []byte(printSchema(e.schema)), 0644)
		})
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

// addWhereType returns the a <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *Extension) whereType(t *gen.Type) (*ast.Definition, error) {
	names, err := nodePaginationNames(t)
	if err != nil {
		return nil, err
	}
	var (
		name    = names.WhereInput
		typeDef = &ast.Definition{
			Name:        name,
			Kind:        ast.InputObject,
			Description: fmt.Sprintf("%s is used for filtering %s objects.\nInput was generated by ent.", name, t.Name),
			Fields: ast.FieldList{
				&ast.FieldDefinition{
					Name: "not",
					Type: ast.NamedType(name, nil),
				},
			},
		}
	)
	for _, op := range []string{"and", "or"} {
		typeDef.Fields = append(typeDef.Fields, &ast.FieldDefinition{
			Name: op,
			Type: ast.ListType(ast.NonNullNamedType(name, nil), nil),
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
			fd := e.fieldDefinition(f, op)
			if i == 0 {
				fd.Description = f.Name + " field predicates"
			}
			typeDef.Fields = append(typeDef.Fields, fd)
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

		typeDef.Fields = append(typeDef.Fields,
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
	return typeDef, nil
}

func (e *Extension) fieldDefinition(f *gen.Field, op gen.Op) *ast.FieldDefinition {
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

var (
	_ entc.Extension = (*Extension)(nil)

	camel  = gen.Funcs["camel"].(func(string) string)
	plural = gen.Funcs["plural"].(func(string) string)
)
