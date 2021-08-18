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
		path       string
		doc        *ast.Document
		cfg        *config.Config
		hooks      []gen.Hook
		templates  []*gen.Template
		scalarFunc func(*gen.Field, gen.Op) string
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
		ex.hooks = append(ex.hooks, ex.genWhereInputs())
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
	ex.hooks = append(ex.hooks, removeOldGeneratedFiles())
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
	case t.Float():
		scalar = graphql.Float.Name()
	case t.Numeric():
		scalar = graphql.Int.Name()
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

// genWhereInputs returns a new hook for generating
// <T>WhereInputs in the GraphQL schema.
func (e *Extension) genWhereInputs() gen.Hook {
	return func(next gen.Generator) gen.Generator {
		if _, exists := e.whereExists(); !exists {
			return next
		}
		inputs := make(map[string]*ast.InputObjectDefinition)
		return gen.GenerateFunc(func(g *gen.Graph) error {
			nodes, err := filterNodes(g.Nodes)
			if err != nil {
				return err
			}
			if err := next.Generate(g); err != nil {
				return err
			}
			for _, node := range nodes {
				name, input, err := e.whereType(node)
				if err != nil {
					return err
				}
				inputs[name] = input
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
func (e *Extension) updateSchema(inputs map[string]*ast.InputObjectDefinition) error {
	visitor.Visit(e.doc, &visitor.VisitorOptions{
		LeaveKindMap: map[string]visitor.VisitFunc{
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
	// Sorting the input types is not needed, because in the next iteration
	// the hook updates the generated types without changing their position.
	for _, input := range inputs {
		e.doc.Definitions = append(e.doc.Definitions, input)
	}
	return ioutil.WriteFile(e.path, []byte(printer.Print(e.doc).(string)), 0644)
}

// addWhereType returns the a <T>WhereInput to the given schema type (e.g. User -> UserWhereInput).
func (e *Extension) whereType(t *gen.Type) (string, *ast.InputObjectDefinition, error) {
	var (
		name  = t.Name + "WhereInput"
		input = ast.NewInputObjectDefinition(&ast.InputObjectDefinition{
			Name: ast.NewName(&ast.Name{
				Value: name,
			}),
			Description: ast.NewStringValue(&ast.StringValue{
				Value: fmt.Sprintf("%s is used for filtering %s objects.\nInput was generated by ent.", name, t.Name),
			}),
			Fields: []*ast.InputValueDefinition{
				ast.NewInputValueDefinition(&ast.InputValueDefinition{
					Name: ast.NewName(&ast.Name{
						Value: "not",
					}),
					Type: ast.NewNamed(&ast.Named{
						Name: ast.NewName(&ast.Name{
							Value: name,
						}),
					}),
				}),
			},
		})
	)
	for _, op := range []string{"and", "or"} {
		input.Fields = append(input.Fields, ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: ast.NewName(&ast.Name{
				Value: op,
			}),
			Type: ast.NewList(&ast.List{
				Type: ast.NewNonNull(&ast.NonNull{
					Type: ast.NewNamed(&ast.Named{
						Name: ast.NewName(&ast.Name{
							Value: name,
						}),
					}),
				}),
			}),
		}))
	}
	fields, err := filterFields(append(t.Fields, t.ID))
	if err != nil {
		return "", nil, err
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
		return "", nil, err
	}
	for _, e := range edges {
		input.Fields = append(input.Fields, ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: ast.NewName(&ast.Name{
				Value: camel("has_" + e.Name),
			}),
			Type: ast.NewNamed(&ast.Named{
				Name: ast.NewName(&ast.Name{
					Value: "Boolean",
				}),
			}),
			Description: ast.NewStringValue(&ast.StringValue{
				Value: e.Name + " edge predicates",
			}),
		}), ast.NewInputValueDefinition(&ast.InputValueDefinition{
			Name: ast.NewName(&ast.Name{
				Value: camel("has_" + e.Name + "_with"),
			}),
			Type: ast.NewList(&ast.List{
				Type: ast.NewNonNull(&ast.NonNull{
					Type: ast.NewNamed(&ast.Named{
						Name: ast.NewName(&ast.Name{
							Value: e.Type.Name + "WhereInput",
						}),
					}),
				}),
			}),
		}))
	}
	return name, input, nil
}

func (e *Extension) fieldDefinition(f *gen.Field, op gen.Op) *ast.InputValueDefinition {
	name := camel(f.Name + "_" + op.Name())
	if op == gen.EQ {
		name = camel(f.Name)
	}
	def := ast.NewInputValueDefinition(&ast.InputValueDefinition{
		Name: ast.NewName(&ast.Name{
			Value: name,
		}),
		Type: ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{
				Value: e.mapScalar(f, op),
			}),
		}),
	})
	if f.Name == "id" {
		def.Type = ast.NewNamed(&ast.Named{
			Name: ast.NewName(&ast.Name{
				Value: graphql.ID.Name(),
			}),
		})
	}
	if op.Variadic() {
		def.Type = ast.NewList(&ast.List{
			Type: ast.NewNonNull(&ast.NonNull{
				Type: def.Type,
			}),
		})
	}
	return def
}

var (
	_     entc.Extension = (*Extension)(nil)
	camel                = gen.Funcs["camel"].(func(string) string)
)