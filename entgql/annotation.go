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

	"entgo.io/ent/schema"
	"github.com/vektah/gqlparser/v2/ast"
)

// Annotation annotates fields and edges with metadata for templates.
type Annotation struct {
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
	Skip bool `json:"Skip,omitempty"`
	// RelayConnection expose this node as a relay connection
	RelayConnection bool `json:"RelayConnection,omitempty"`
	// GQLName provide alternative name. see: https://gqlgen.com/config/#inline-config-with-directives
	GQLName string `json:"GQLName,omitempty"`
	// GQLImplements extra interfaces that are implemented
	GQLImplements []string `json:"GQLImplements,omitempty"`
	// GQLDirectives directives to add
	GQLDirectives []Directive `json:"GQLDirectives,omitempty"`
}

type Directive struct {
	Name      string              `json:"name,omitempty"`
	Arguments []DirectiveArgument `json:"arguments,omitempty"`
}

type DirectiveArgument struct {
	Name  string        `json:"name,omitempty"`
	Value string        `json:"value,omitempty"`
	Kind  ast.ValueKind `json:"kind,omitempty"`
}

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
func Skip() Annotation {
	return Annotation{Skip: true}
}

// Skip returns a relay connection annotation.
func RelayConnection() Annotation {
	return Annotation{RelayConnection: true}
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
	if ant.Skip {
		a.Skip = true
	}
	if ant.RelayConnection {
		a.RelayConnection = true
	}
	if ant.GQLName != "" {
		a.GQLName = ant.GQLName
	}
	if len(ant.GQLDirectives) > 0 {
		a.GQLDirectives = append(a.GQLDirectives, ant.GQLDirectives...)
	}
	if len(ant.GQLImplements) > 0 {
		a.GQLImplements = append(a.GQLImplements, ant.GQLImplements...)
	}
	return a
}

// Decode unmarshal annotation
func (a *Annotation) Decode(annotation interface{}) error {
	buf, err := json.Marshal(annotation)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Merger     = (*Annotation)(nil)
)
