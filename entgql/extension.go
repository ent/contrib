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
	"github.com/99designs/gqlgen/plugin"
	"github.com/99designs/gqlgen/plugin/federation"
	"github.com/99designs/gqlgen/plugin/modelgen"
	"github.com/99designs/gqlgen/plugin/resolvergen"
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

		schema *ast.Schema
		models map[string]string
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
		ex.path = path
		ex.hooks = append(ex.hooks, ex.genWhereInputs())
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
		if cfg.Schema == nil {
			// Copied from api/generate.go
			// https://github.com/99designs/gqlgen/blob/47015f12e3aa26af251fec67eab50d3388c17efe/api/generate.go#L21-L57
			var plugins []plugin.Plugin
			if cfg.Model.IsDefined() {
				plugins = append(plugins, modelgen.New())
			}
			plugins = append(plugins, resolvergen.New())
			if cfg.Federation.IsDefined() {
				plugins = append([]plugin.Plugin{federation.New()}, plugins...)
			}
			for _, opt := range gqlgenOptions {
				opt(cfg, &plugins)
			}
			for _, p := range plugins {
				if inj, ok := p.(plugin.EarlySourceInjector); ok {
					if s := inj.InjectSourceEarly(); s != nil {
						cfg.Sources = append(cfg.Sources, s)
					}
				}
			}
			if err := cfg.LoadSchema(); err != nil {
				return fmt.Errorf("failed to load schema: %w", err)
			}
			for _, p := range plugins {
				if inj, ok := p.(plugin.LateSourceInjector); ok {
					if s := inj.InjectSourceLate(cfg.Schema); s != nil {
						cfg.Sources = append(cfg.Sources, s)
					}
				}
			}
			// LoadSchema again now we have everything.
			if err := cfg.LoadSchema(); err != nil {
				return fmt.Errorf("failed to load schema: %w", err)
			}
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
		e.hooks = append(e.hooks, e.genSchema())
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
//	)
//
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{templates: AllTemplates}
	for _, opt := range opts {
		if err := opt(ex); err != nil {
			return nil, err
		}
	}
	ex.hooks = append(ex.hooks, removeOldAssets)
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
	if t, ok := typeAnnotation(f); ok {
		return t
	}
	return scalar
}

// hasMapping reports if the gqlgen.yml has custom mapping for
// the given field type and returns its GraphQL name if exists.
func (e *Extension) hasMapping(f *gen.Field) (string, bool) {
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
			if strings.HasSuffix(m, ident) && e.isInput(t) {
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

// isInput reports if the given type is an input object.
func (e *Extension) isInput(name string) bool {
	if t, ok := e.cfg.Schema.Types[name]; ok && t != nil {
		return t.IsInputType()
	}
	return false
}

// genSchema returns a new hook for generating
// the GraphQL schema from the graph.
func (e *Extension) genSchema() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) error {
			if err := next.Generate(g); err != nil {
				return err
			}

			genSchema, err := newSchemaGenerator(g)
			if err != nil {
				return err
			}
			if e.schema, err = genSchema.prepareSchema(); err != nil {
				return err
			}
			if e.models, err = genSchema.genModels(); err != nil {
				return err
			}

			return nil
		})
	}
}

// genWhereInputs returns a new hook for generating
// <T>WhereInputs in the GraphQL schema.
func (e *Extension) genWhereInputs() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		if _, exists := e.whereExists(); !exists {
			return next
		}
		inputs := make(map[string]*ast.Definition)
		return gen.GenerateFunc(func(g *gen.Graph) error {
			nodes, err := filterNodes(g.Nodes, SkipFlagWhere)
			if err != nil {
				return err
			}
			if err := next.Generate(g); err != nil {
				return err
			}
			for _, node := range nodes {
				input, err := e.whereType(node)
				if err != nil {
					return err
				}
				inputs[input.Name] = input
			}
			return e.updateSchema(inputs)
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

// updateSchema commits the changes to the GraphQL schema file.
func (e *Extension) updateSchema(inputs map[string]*ast.Definition) error {
	schema := &ast.Schema{
		Types: inputs,
	}
	return ioutil.WriteFile(e.path, []byte(printSchema(schema)), 0644)
}

// addWhereType returns the a <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *Extension) whereType(t *gen.Type) (*ast.Definition, error) {
	var (
		name    = t.Name + "WhereInput"
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

	fields, err := filterFields(append(t.Fields, t.ID), SkipFlagWhere)
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

	edges, err := filterEdges(t.Edges, SkipFlagWhere)
	if err != nil {
		return nil, err
	}
	for _, e := range edges {
		typeDef.Fields = append(typeDef.Fields,
			&ast.FieldDefinition{
				Name:        camel("has_" + e.Name),
				Type:        ast.NamedType("Boolean", nil),
				Description: e.Name + " edge predicates",
			},
			&ast.FieldDefinition{
				Name: camel("has_" + e.Name + "_with"),
				Type: ast.ListType(ast.NonNullNamedType(
					e.Type.Name+"WhereInput", nil), nil),
			},
		)
	}
	return typeDef, nil
}

func (e *Extension) fieldDefinition(f *gen.Field, op gen.Op) *ast.FieldDefinition {
	name := camel(f.Name + "_" + op.Name())
	if op == gen.EQ {
		name = camel(f.Name)
	}

	typeName := e.mapScalar(f, op)
	if f.Name == "id" {
		typeName = "ID"
	}
	def := &ast.FieldDefinition{
		Name: name,
	}
	if op.Variadic() {
		def.Type = ast.ListType(ast.NonNullNamedType(typeName, nil), nil)
	} else {
		def.Type = ast.NamedType(typeName, nil)
	}
	return def
}

var (
	_     entc.Extension = (*Extension)(nil)
	camel                = gen.Funcs["camel"].(func(string) string)
)

// typeAnnotation returns the scalar type mapping if exists (i.e. entgql.Type).
func typeAnnotation(f *gen.Field) (string, bool) {
	var ant Annotation
	if i, ok := f.Annotations[ant.Name()]; ok && ant.Decode(i) == nil && ant.Type != "" {
		return ant.Type, true
	}
	return "", false
}
