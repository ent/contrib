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
		// MutationInputs defines the input types for the mutation.
		MutationInputs MutationInputType `json:"MutationInputs,omitempty"`
	}

	// Directive to apply on the field/type.
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

	// MutationInputType is a bit flag for the mutation input type.
	MutationInputType int

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
	// SkipMutationCreateInput skips generating GraphQL Create<Type>Input types.
	// If defined on a field, the type will be generated without the field.
	SkipMutationCreateInput
	// SkipMutationUpdateInput skips generating GraphQL Update<Type>Input types.
	// If defined on a field, the type will be generated without the field.
	SkipMutationUpdateInput

	// SkipAll is default mode to skip all.
	SkipAll = SkipType |
		SkipEnumField |
		SkipOrderField |
		SkipWhereInput |
		SkipMutationCreateInput |
		SkipMutationUpdateInput
)

const (
	MutationCreate MutationInputType = 1 << iota
	MutationUpdate
)

// Name implements ent.Annotation interface.
func (Annotation) Name() string {
	return "EntGQL"
}

// OrderField enables ordering in GraphQL for the annotated Ent field
// with the given name. Note that, the field type must be comparable.
//
//	field.Time("created_at").
//		Default(time.Now).
//		Immutable().
//		Annotations(
//			entgql.OrderField("CREATED_AT"),
//		)
//
// For enum fields, values must be uppercase or mapped using the NamedValues
// option:
//
//	field.Enum("status").
//		NamedValues(
//			"InProgress", "IN_PROGRESS",
//			"Completed", "COMPLETED",
//		).
//		Default("IN_PROGRESS").
//		Annotations(
//			entgql.OrderField("STATUS"),
//		)
//
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
// option in order to load edges with different name mapping.
//
//	func (Todo) Edges() []ent.Edge {
//		return []ent.Edge{
//			edge.To("parent", Todo.Type).
//			Annotations(entgql.Unbind()).
//			Unique().
//			From("children").
//			Annotations(entgql.Unbind()),
//		}
//	}
//
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
// The Type() annotation is used to map the underlying
// GraphQL type to the type.
//
// To change the GraphQL type for a type:
//
//	func (User) Annotations() []schema.Annotation {
//		return []schema.Annotation{
//			entgql.Type("MasterUser"),
//		}
//	}
//
// To change the GraphQL type for a field (rename enum type):
//
//	field.Enum("status").
//		NamedValues(
//			"InProgress", "IN_PROGRESS",
//			"Completed", "COMPLETED",
//		).
//		Default("IN_PROGRESS").
//		Annotations(
//			entgql.Type("TodoStatus"),
//		)
func Type(name string) Annotation {
	return Annotation{Type: name}
}

// Skip returns a skip annotation.
// The Skip() annotation is used to skip
// generating the type or the field from GraphQL schema.
//
// It gives you the flexibility to skip generating
// the type or the field based on the SkipMode flags.
//
// For example, if you don't want to expose a field on the <T>WhereInput type
// you can use the following:
//
//	field.String("name").
//		Annotations(
//			entgql.Skip(entgql.SkipWhereInput),
//		)
//
// Since SkipMode is a bit flag, it's possible to skip multiple modes using
// the bitwise OR operator as follows:
//
//	entgql.Skip(entgql.SkipWhereInput | entgql.SkipEnumField)
//
//
// To skip everything except the type, use the bitwise NOT operator:
//
//	entgql.Skip(^entgql.SkipType)
//
// You can also skip all modes with the `entgql.SkipAll` constant which is the default mode.
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
// The RelayConnection() annotation is used to generate
// the Relay <T>Edge, <T>Connection, and PageInfo types for a type. For example:
//
//	func (Todo) Annotations() []schema.Annotation {
//		return []schema.Annotation{
//			entgql.RelayConnection(),
//			entgql.QueryField(),
//		}
//	}
//
// The RelayConnection() annotation can also be used on
// the edge fields, to generate first, last, after, before... arguments and
// change the field type to `<T>Connection!`.
// For example to change the children field from `children: [Todo!]!` to
// `children(first: Int, last: Int, after: Cursor, before: Cursor): TodoConnection!`
//
//	func (Todo) Edges() []ent.Edge {
//		return []ent.Edge{
//				edge.To("parent", Todo.Type).
//					Unique().
//					From("children").
//					Annotation(entgql.RelayConnection()),
//		}
//	}
func RelayConnection() Annotation {
	return Annotation{RelayConnection: true}
}

// Implements returns an Implements annotation.
// The Implements() annotation is used to
// add implements interfaces to a GraphQL type.
//
// For example, to add the `Entity` to the `Todo` type:
//
//	type Todo implements Node {
//		id: ID!
//		...
//	}
//
// Add the entgql.Implements("Entity") to the Todo's annotations:
//
//	func (Todo) Annotations() []schema.Annotation {
//		return []schema.Annotation{
//			entgql.Implements("Entity"),
//		}
//	}
//
// and the GraphQL type will be generated with the implements interface.
//
//	type Todo implements Node & Entity {
//		id: ID!
//		...
//	}
func Implements(interfaces ...string) Annotation {
	return Annotation{Implements: interfaces}
}

// Directives returns a Directives annotation.
// The Directives() annotation is used to
// add directives to a GraphQL type or on the field.
//
// For example, to add the `@deprecated` directive to the `text` field:
//
//	type Todo {
//		id: ID
//		text: String
//		...
//	}
//
// Add the entgql.Directives() to the text field's annotations:
//
//	field.Text("text").
//		NotEmpty().
//		Annotations(
//			entgql.Directives("Use `description` instead."),
//		),
//
// and the GraphQL type will be generated with the directive.
//
//	type Todo @deprecated {
//		id: ID
//		text: String @deprecated(reason: "Use `description` instead.")
//		...
//	}
func Directives(directives ...Directive) Annotation {
	return Annotation{Directives: directives}
}

type queryFieldAnnotation struct {
	Annotation
}

// QueryField returns an annotation for expose the field on the Query type.
func QueryField(name ...string) queryFieldAnnotation {
	a := Annotation{QueryField: &FieldConfig{}}
	if len(name) > 0 {
		a.QueryField.Name = name[0]
	}
	return queryFieldAnnotation{Annotation: a}
}

// Directives allow you apply directives to the field.
func (a queryFieldAnnotation) Directives(directives ...Directive) queryFieldAnnotation {
	a.QueryField.Directives = directives
	return a
}

// MutationInput returns an annotation for generate input types for mutation.
func MutationInput(inputs ...MutationInputType) Annotation {
	if len(inputs) == 0 {
		return Annotation{MutationInputs: MutationCreate | MutationUpdate}
	}

	input := MutationInputType(0)
	for _, f := range inputs {
		input |= f
	}
	return Annotation{MutationInputs: input}
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
	case queryFieldAnnotation:
		ant = other.Annotation
	case *queryFieldAnnotation:
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
	if ant.MutationInputs.Any() {
		a.MutationInputs |= ant.MutationInputs
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

// Any returns true if the mutation input annotation was set.
func (f MutationInputType) Any() bool {
	return f != 0
}

// Is checks if the mutation input annotation has a specific type.
func (f MutationInputType) Has(mode MutationInputType) bool {
	return f&mode != 0
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

// NewDirective returns a GraphQL directive
// to use with the entgql.Directives annotation.
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
