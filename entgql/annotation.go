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
	OrderField string `json:"order_field,omitempty"`
	// Bind implies the edge field name in graphql schema
	// is equivalent to the name used in ent schema.
	Bind bool `json:"bind,omitempty"`
	// Mapping is the edge field names as defined in graphql schema.
	Mapping []string `json:"mapping,omitempty"`
	// Skip exclude the type
	Skip bool `json:"skip,omitempty"`
	// RelayConnection expose this node as a relay connection
	RelayConnection bool `json:"relay_connection,omitempty"`
	// GqlName provide alternative name. see: https://gqlgen.com/config/#inline-config-with-directives
	GqlName string `json:"gql_name,omitempty"`
	// GqlType override type
	GqlType string `json:"gql_type,omitempty"`
	// GqlImplements extra interfaces that are implemented
	GqlImplements []string `json:"gql_implements,omitempty"`
	// GqlDirectives directives to add
	GqlDirectives []Directive `json:"gql_directives,omitempty"`
	// GqlScalarMappings defines custom scalars mappings, scalars will also be created automatically
	GqlScalarMappings map[string]string `json:"gql_scalar_mappings,omitempty"`
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
func Bind() Annotation {
	return Annotation{Bind: true}
}

// MapsTo returns a mapping annotation.
func MapsTo(names ...string) Annotation {
	return Annotation{Mapping: names}
}

// Skip returns a skip annotation.
func Skip() Annotation {
	return Annotation{Skip: true}
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
	if ant.Bind {
		a.Bind = true
	}
	if len(ant.Mapping) != 0 {
		a.Mapping = ant.Mapping
	}
	if ant.Skip {
		a.Skip = true
	}
	if ant.RelayConnection {
		a.RelayConnection = true
	}
	if ant.GqlName != "" {
		a.GqlName = ant.GqlName
	}
	if ant.GqlType != "" {
		a.GqlType = ant.GqlType
	}
	if len(ant.GqlDirectives) > 0 {
		a.GqlDirectives = ant.GqlDirectives
	}
	if len(ant.GqlImplements) > 0 {
		a.GqlImplements = ant.GqlImplements
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
