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

package plugin

import (
	"fmt"
	"strings"

	"entgo.io/contrib/entgql"
	"entgo.io/ent/entc/gen"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

var (
	camel          = gen.Funcs["camel"].(func(string) string)
	annotationName = entgql.Annotation{}.Name()
)

type (
	EntGQL struct {
		debug          bool
		genTypes       []*gen.Type
		scalarMappings map[string]string
		schema         *ast.Schema
		hooks          []SchemaHook
		graph          *gen.Graph
	}

	PluginOption func(*EntGQL) error
)

var (
	_ plugin.Plugin              = &EntGQL{}
	_ plugin.ConfigMutator       = &EntGQL{}
	_ plugin.EarlySourceInjector = &EntGQL{}
)

func WithSchemaHooks(hooks ...SchemaHook) PluginOption {
	return func(e *EntGQL) error {
		e.hooks = hooks
		return nil
	}
}

func WithDebug() PluginOption {
	return func(e *EntGQL) error {
		e.debug = true
		return nil
	}
}

// SchemaHook hook to modify schema before printing
type SchemaHook func(schema *ast.Schema)

func New(graph *gen.Graph, opts ...PluginOption) (*EntGQL, error) {
	types, err := entgql.FilterNodes(graph.Nodes)
	if err != nil {
		return nil, err
	}

	// Include default mapping for time
	scalarMappings := map[string]string{
		"Time": "Time",
	}
	if graph.Annotations != nil {
		globalAnn := graph.Annotations[annotationName]
		// TODO: cleanup assertions
		if globalAnn != nil {
			if globalAnn.(entgql.Annotation).GqlScalarMappings != nil {
				scalarMappings = globalAnn.(entgql.Annotation).GqlScalarMappings
			}
		}
	}

	e := &EntGQL{
		graph:          graph,
		genTypes:       types,
		scalarMappings: scalarMappings,
		schema: &ast.Schema{
			Types:         map[string]*ast.Definition{},
			Directives:    map[string]*ast.DirectiveDefinition{},
			PossibleTypes: map[string][]*ast.Definition{},
			Implements:    map[string][]*ast.Definition{},
		},
	}
	for _, opt := range opts {
		if err = opt(e); err != nil {
			return nil, err
		}
	}

	return e, nil
}

func (e *EntGQL) Name() string {
	return "entgql"
}

func (e *EntGQL) InjectSourceEarly() *ast.Source {
	e.scalars()
	e.relayBuiltins()
	e.entBuiltins()
	err := e.entOrderBy()
	if err != nil {
		panic(err)
	}
	err = e.enums()
	if err != nil {
		panic(err)
	}
	err = e.types()
	if err != nil {
		panic(err)
	}
	for _, h := range e.hooks {
		h(e.schema)
	}
	input := e.print()
	if e.debug {
		fmt.Printf("Generated Graphql:\n%s", input)
	}

	return &ast.Source{
		Name:    "entgql.graphql",
		Input:   input,
		BuiltIn: false,
	}
}

func (e *EntGQL) print() string {
	sb := &strings.Builder{}
	printer := formatter.NewFormatter(sb)
	printer.FormatSchema(e.schema)
	return sb.String()
}
