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
	"entgo.io/ent/schema/field"
	"github.com/99designs/gqlgen/plugin"
	"github.com/vektah/gqlparser/v2/ast"
	"github.com/vektah/gqlparser/v2/formatter"
)

type (
	// EntGQL is a plugin that generates GQL schema from the Ent's Graph
	EntGQL struct {
		graph  *gen.Graph
		nodes  []*gen.Type
		schema *ast.Schema

		relaySpec bool
	}

	// EntGQLOption is a option for the EntGQL plugin
	EntGQLOption func(*EntGQL) error
)

var (
	annotationName = entgql.Annotation{}.Name()
	camel          = gen.Funcs["camel"].(func(string) string)

	_ plugin.Plugin              = (*EntGQL)(nil)
	_ plugin.EarlySourceInjector = (*EntGQL)(nil)
	_ plugin.ConfigMutator       = (*EntGQL)(nil)
)

// WithRelaySpecification adds the Relay specification to the schema
func WithRelaySpecification(relaySpec bool) EntGQLOption {
	return func(e *EntGQL) error {
		e.relaySpec = relaySpec
		return nil
	}
}

// NewEntGQLPlugin creates a new EntGQL plugin
func NewEntGQLPlugin(graph *gen.Graph, opts ...EntGQLOption) (*EntGQL, error) {
	nodes, err := entgql.FilterNodes(graph.Nodes)
	if err != nil {
		return nil, err
	}

	e := &EntGQL{
		graph: graph,
		nodes: nodes,
	}
	for _, opt := range opts {
		if err = opt(e); err != nil {
			return nil, err
		}
	}

	types, err := e.buildTypes()
	if err != nil {
		return nil, err
	}

	e.schema = &ast.Schema{
		Types: types,
	}
	return e, nil
}

// Name implements the Plugin interface.
func (*EntGQL) Name() string {
	return "entgql"
}

// InjectSourceEarly implements the EarlySourceInjector interface.
func (e *EntGQL) InjectSourceEarly() *ast.Source {
	return nil
	// return &ast.Source{
	// 	Name:    "entgql.graphql",
	// 	Input:   printSchema(e.schema),
	// 	BuiltIn: false,
	// }
}

func (e *EntGQL) buildTypes() (map[string]*ast.Definition, error) {
	types := map[string]*ast.Definition{}
	for _, node := range e.nodes {
		ant, err := entgql.DecodeAnnotation(node.Annotations)
		if err != nil {
			return nil, err
		}
		if ant.Skip {
			continue
		}

		fields, err := e.buildTypeFields(node)
		if err != nil {
			return nil, err
		}

		name := node.Name
		if ant.Type != "" {
			name = ant.Type
		}

		var interfaces []string
		if e.relaySpec {
			interfaces = append(interfaces, "Node")
		}

		types[name] = &ast.Definition{
			Name:       name,
			Kind:       ast.Object,
			Fields:     fields,
			Interfaces: interfaces,
		}

		// TODO(giautm): Added RelayConnection annotation check
		if e.relaySpec {
			defs, err := relayConnectionTypes(node)
			if err != nil {
				return nil, err
			}

			insertDefinitions(types, defs...)
		}
	}

	return types, nil
}

func (e *EntGQL) buildTypeFields(t *gen.Type) (ast.FieldList, error) {
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

func (e *EntGQL) typeField(f *gen.Field, isID bool) ([]*ast.FieldDefinition, error) {
	ant, err := entgql.DecodeAnnotation(f.Annotations)
	if err != nil {
		return nil, err
	}
	if ant.Skip {
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
			Name: camel(f.Name),
			Type: ft,
		},
	}, nil
}

func namedType(name string, nullable bool) *ast.Type {
	if nullable {
		return ast.NamedType(name, nil)
	}
	return ast.NonNullNamedType(name, nil)
}

func (e *EntGQL) typeFromField(f *gen.Field, idField bool, userDefinedType string) (*ast.Type, error) {
	nillable := f.Nillable
	typ := f.Type.Type

	// TODO(giautm): Support custom scalar types
	// TODO(giautm): Support Edge Field
	// TODO(giautm): Support some built-in JSON types: Ints(), Floats(), Strings()
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
		return nil, fmt.Errorf("json type not implemented")
	case typ == field.TypeOther:
		return nil, fmt.Errorf("other type must have typed defined")
	default:
		return nil, fmt.Errorf("unexpected type: %s", typ.String())
	}
}

func printSchema(schema *ast.Schema) string {
	sb := &strings.Builder{}
	formatter.
		NewFormatter(sb).
		FormatSchema(schema)
	return sb.String()
}
