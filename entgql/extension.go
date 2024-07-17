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
	"fmt"
	"os"
	"path/filepath"
	"slices"

	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/api"
	"github.com/99designs/gqlgen/codegen/config"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	// Extension implements the entc.Extension for providing GraphQL integration.
	Extension struct {
		schemaGenerator
		entc.DefaultExtension
		outputWriter func(*ast.Schema) error
		hooks        []gen.Hook
		templates    []*gen.Template
	}

	// ExtensionOption allows for managing the Extension configuration
	// using functional options.
	ExtensionOption func(*Extension) error

	// SchemaHook is the hook that run after the GQL schema generation.
	SchemaHook func(*gen.Graph, *ast.Schema) error
)

// WithSchemaPath sets the filepath to the GraphQL schema to write the
// generated Ent types. If the file does not exist, it will generate a
// new schema. Please note that your gqlgen.yml config file should be
// updated as follows to support multiple schema files:
//
//	schema:
//	 - schema.graphql // existing schema.
//	 - ent.graphql	  // generated schema.
func WithSchemaPath(path string) ExtensionOption {
	return func(ex *Extension) error {
		ex.path = path
		return nil
	}
}

// WithOutputWriter sets the function to write the generated schema.
func WithOutputWriter(w func(*ast.Schema) error) ExtensionOption {
	return func(ex *Extension) error {
		ex.outputWriter = w
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

var (
	// WithWhereFilters configures the extension to either add or
	// remove the WhereTemplate from the code generation templates.
	//
	// Deprecated: use WithWhereInputs instead. This option is planned
	// to be removed in future versions.
	WithWhereFilters = WithWhereInputs
)

// WithWhereInputs configures the extension to either add or
// remove the WhereTemplate from the code generation templates.
//
// The WhereTemplate generates GraphQL filters to all types in the ent/schema.
func WithWhereInputs(b bool) ExtensionOption {
	return func(ex *Extension) error {
		ex.genWhereInput = b
		i, exists := ex.hasTemplate(WhereTemplate)
		if b && !exists {
			ex.templates = append(ex.templates, WhereTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithNodeDescriptor configures the extension to either add or
// remove the NodeDescriptorTemplate from the code generation templates.
//
// In case this option is enabled, EntGQL generates a `Node()` method for each
// type that returns its representation in one standard way. A common use case for
// this option is to develop an administrator tool on top of Ent as implemented in:
// https://github.com/ent/ent/issues/1000#issuecomment-735663175.
func WithNodeDescriptor(b bool) ExtensionOption {
	return func(ex *Extension) error {
		i, exists := ex.hasTemplate(NodeDescriptorTemplate)
		if b && !exists {
			ex.templates = append(ex.templates, NodeDescriptorTemplate)
		} else if !b && exists && len(ex.templates) > 0 {
			ex.templates = append(ex.templates[:i], ex.templates[i+1:]...)
		}
		return nil
	}
}

// WithRelaySpec enables or disables generating the Relay Node interface.
func WithRelaySpec(enabled bool) ExtensionOption {
	return func(e *Extension) error {
		e.relaySpec = enabled
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
func WithMapScalarFunc(scalarFunc func(*gen.Field, gen.Op) string) ExtensionOption {
	return func(ex *Extension) error {
		ex.scalarFunc = scalarFunc
		return nil
	}
}

// NewExtension creates a new extension with the given configuration.
//
//	ex, err := entgql.NewExtension(
//		entgql.WithSchemaGenerator(),
//		entgql.WithSchemaPath("../ent.graphql"),
//		entgql.WithWhereInputs(true),
//	)
func NewExtension(opts ...ExtensionOption) (*Extension, error) {
	ex := &Extension{
		templates: AllTemplates,
		schemaGenerator: schemaGenerator{
			relaySpec:    true,
			genMutations: true,
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

// Options of the extension.
func (e *Extension) Options() []entc.Option {
	return []entc.Option{
		entc.FeatureNames(gen.FeatureNamedEdges.Name),
	}
}

// genSchema returns a new hook for generating
// the GraphQL schema from the graph.
func (e *Extension) genSchemaHook() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		return gen.GenerateFunc(func(g *gen.Graph) (err error) {
			if err = next.Generate(g); err != nil {
				return err
			}
			for _, t := range g.Nodes {
				for _, f := range t.DeprecatedFields() {
					ant, err := annotation(f.Annotations)
					if err != nil {
						return err
					}
					if !slices.ContainsFunc(ant.Directives, func(d Directive) bool {
						return d.Name == "deprecated"
					}) {
						ant.Directives = append(ant.Directives, Deprecated(f.DeprecationReason()))
						if f.Annotations == nil {
							f.Annotations = make(map[string]interface{})
						}
						f.Annotations[ant.Name()] = ant
					}
				}
			}
			if !(e.genSchema || e.genWhereInput || e.genMutations) {
				return nil
			}
			schema, err := e.BuildSchema(g)
			if err != nil {
				return err
			}
			if e.outputWriter == nil {
				if e.path == "" {
					return nil
				}
				return os.WriteFile(e.path, []byte(printSchema(schema)), 0644)
			}
			return e.outputWriter(schema)
		})
	}
}

// hasTemplate reports if the template exists
// in the template list and returns its index.
func (e *Extension) hasTemplate(tem *gen.Template) (int, bool) {
	for i := range e.templates {
		if e.templates[i].Name() == tem.Templates()[1].Name() {
			return i, true
		}
	}
	return -1, false
}

var (
	_ entc.Extension = (*Extension)(nil)

	camel    = gen.Funcs["camel"].(func(string) string)
	pascal   = gen.Funcs["pascal"].(func(string) string)
	plural   = gen.Funcs["plural"].(func(string) string)
	singular = gen.Funcs["singular"].(func(string) string)
	snake    = gen.Funcs["snake"].(func(string) string)
)
