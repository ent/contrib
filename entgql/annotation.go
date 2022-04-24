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
	"encoding/json"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"github.com/vektah/gqlparser/v2/ast"
)

type (
	// Annotation annotates fields and edges with metadata for templates.
	Annotation struct {
		// OrderField is the ordering field as defined in graphql schema.
		OrderField string `json:"OrderField,omitempty"`
		// Unbind implies the edge field name in GraphQL schema is not equivalent
		// to the name used in ent schema. That means, by default, edges with this
		// annotation will not be eager-loaded on Paginate calls. See the `MapsTo`
		// option in order to load edges be different name mapping.
		Unbind bool `json:"Unbind,omitempty"`
		// Mapping is the edge field names as defined in graphql schema.
		Mapping []string `json:"Mapping,omitempty"`
		// Type is the underlying GraphQL type name (e.g. Boolean).
		Type string `json:"Type,omitempty"`
		// Skip exclude the type
		Skip SkipMode `json:"Skip,omitempty"`
		// RelayConnection enables the Relay Connection specification for the entity.
		// It's also can apply on an edge to create the Relay-style filter.
		RelayConnection bool `json:"RelayConnection,omitempty"`
		// Implements defines a list of interfaces implemented by the type.
		Implements []string `json:"Implements,omitempty"`
		// Directives to add on the field/type.
		Directives []Directive `json:"Directives,omitempty"`
		// QueryField exposes the generated type with the given string under the Query object.
		QueryField *FieldConfig `json:"QueryField,omitempty"`
	}

	// Directive to apply on the field/type
	Directive struct {
		Name      string              `json:"name,omitempty"`
		Arguments []DirectiveArgument `json:"arguments,omitempty"`
	}

	// DirectiveArgument return a GraphQL directive argument
	DirectiveArgument struct {
		Name  string        `json:"name,omitempty"`
		Value string        `json:"value,omitempty"`
		Kind  ast.ValueKind `json:"kind,omitempty"`
	}

	// SkipMode is a bit flag for the Skip annotation.
	SkipMode int

	FieldConfig struct {
		// Name is the name of the field in the Query object.
		Name string `json:"Name,omitempty"`
		// Directives to add on the field
		Directives []Directive `json:"Directives,omitempty"`
	}
)

const (
	// SkipType skips generating GraphQL types or fields in the schema.
	SkipType SkipMode = 1 << iota
	// SkipEnumField skips generating GraphQL enums for enum fields in the schema.
	SkipEnumField
	// SkipOrderField skips generating GraphQL order inputs and enums for ordered-fields in the schema.
	SkipOrderField
	// SkipWhereInput skips generating GraphQL WhereInput types.
	// If defined on a field, the type will be generated without the field.
	SkipWhereInput

	// SkipAll is default mode to skip all.
	SkipAll = SkipType |
		SkipEnumField |
		SkipOrderField |
		SkipWhereInput
)

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "EntGQL"
}

// OrderField returns an order field annotation.
func OrderField(name string) Annotation {
	return Annotation{OrderField: name}
}

// Bind returns a binding annotation.
//
// No-op function to avoid breaking the existing schema.
// You can safely remove this function from your scheme.
//
// Deprecated: the Bind option predates the Unbind option, and it is planned
// to be removed in future versions. Users should not use this annotation as it
// is a no-op call, or use `Unbind` in order to disable `Bind`.
func Bind() Annotation {
	return Annotation{}
}

// Unbind implies the edge field name in GraphQL schema is not equivalent
// to the name used in ent schema. That means, by default, edges with this
// annotation will not be eager-loaded on Paginate calls. See the `MapsTo`
// option in order to load edges be different name mapping.
func Unbind() Annotation {
	return Annotation{Unbind: true}
}

// MapsTo returns a mapping annotation.
func MapsTo(names ...string) Annotation {
	return Annotation{
		Mapping: names,
		Unbind:  true, // Unbind because it cant be used with mapping names.
	}
}

// Type returns a type mapping annotation.
func Type(name string) Annotation {
	return Annotation{Type: name}
}

// Skip returns a skip annotation.
func Skip(flags ...SkipMode) Annotation {
	if len(flags) == 0 {
		return Annotation{Skip: SkipAll}
	}

	skip := SkipMode(0)
	for _, f := range flags {
		skip |= f
	}
	return Annotation{Skip: skip}
}

// RelayConnection returns a relay connection annotation.
func RelayConnection() Annotation {
	return Annotation{RelayConnection: true}
}

// Implements returns an Implements annotation.
func Implements(interfaces ...string) Annotation {
	return Annotation{Implements: interfaces}
}

// Directives returns a Directives annotation.
func Directives(directives ...Directive) Annotation {
	return Annotation{Directives: directives}
}

type QueryFieldAnnotation struct {
	Annotation
}

// QueryField returns an annotation for expose the field on the Query type.
func QueryField(name ...string) QueryFieldAnnotation {
	a := Annotation{QueryField: &FieldConfig{}}
	if len(name) > 0 {
		a.QueryField.Name = name[0]
	}
	return QueryFieldAnnotation{Annotation: a}
}

func (a QueryFieldAnnotation) Directives(directives ...Directive) QueryFieldAnnotation {
	a.QueryField.Directives = directives
	return a
}

// Merge implements the schema.Merger interface.
func (a Annotation) Merge(other schema.Annotation) schema.Annotation {
	var ant Annotation
	switch other := other.(type) {
	case Annotation:
		ant = other
	case *Annotation:
		if other != nil {
			ant = *other
		}
	case QueryFieldAnnotation:
		ant = other.Annotation
	case *QueryFieldAnnotation:
		if other != nil {
			ant = other.Annotation
		}
	default:
		return a
	}
	if ant.OrderField != "" {
		a.OrderField = ant.OrderField
	}
	if ant.Unbind {
		a.Unbind = true
	}
	if len(ant.Mapping) != 0 {
		a.Mapping = ant.Mapping
	}
	if ant.Type != "" {
		a.Type = ant.Type
	}
	if ant.Skip.Any() {
		a.Skip |= ant.Skip
	}
	if ant.RelayConnection {
		a.RelayConnection = true
	}
	if len(ant.Implements) > 0 {
		a.Implements = append(a.Implements, ant.Implements...)
	}
	if len(ant.Directives) > 0 {
		a.Directives = append(a.Directives, ant.Directives...)
	}
	if ant.QueryField != nil {
		if a.QueryField == nil {
			a.QueryField = &FieldConfig{}
		}
		a.QueryField.merge(ant.QueryField)
	}
	return a
}

// Decode unmarshalls the annotation.
func (a *Annotation) Decode(annotation interface{}) error {
	buf, err := json.Marshal(annotation)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

// Any returns true if the skip annotation was set.
func (f SkipMode) Any() bool {
	return f != 0
}

// Is checks if the skip annotation has a specific flag.
func (f SkipMode) Is(mode SkipMode) bool {
	return f&mode != 0
}

func (c FieldConfig) fieldName(gqlType string) string {
	if c.Name != "" {
		return c.Name
	}
	return camel(plural(gqlType))
}

func (c *FieldConfig) merge(ant *FieldConfig) {
	if ant == nil {
		return
	}
	if ant.Name != "" {
		c.Name = ant.Name
	}
	c.Directives = append(c.Directives, ant.Directives...)
}

// annotation extracts the entgql.Annotation or returns its empty value.
func annotation(ants gen.Annotations) (*Annotation, error) {
	ant := &Annotation{}
	if ants != nil && ants[ant.Name()] != nil {
		if err := ant.Decode(ants[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Merger     = (*Annotation)(nil)
)

// NewDirective return a GraphQL directive
func NewDirective(name string, args ...DirectiveArgument) Directive {
	return Directive{
		Name:      name,
		Arguments: args,
	}
}

// Deprecated create `@deprecated` directive to apply on the field/type
func Deprecated(reason string) Directive {
	var args []DirectiveArgument
	if reason != "" {
		args = append(args, DirectiveArgument{
			Name:  "reason",
			Kind:  ast.StringValue,
			Value: reason,
		})
	}
	return NewDirective("deprecated", args...)
}
