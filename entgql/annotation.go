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
)

// Annotation annotates fields and edges with metadata for templates.
type Annotation struct {
	// OrderField is the ordering field as defined in graphql schema.
	OrderField string `json:"OrderField,omitempty"`
	// Bind implies the edge field name in graphql schema
	// is equivalent to the name used in ent schema.
	Bind bool `json:"Bind,omitempty"`
	// Mapping is the edge field names as defined in graphql schema.
	Mapping []string `json:"Mapping,omitempty"`
	// Type is the underlying GraphQL type name (e.g. Boolean).
	Type string `json:"Type,omitempty"`
	// Skip exclude the type
	Skip bool `json:"Skip,omitempty"`
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

// Type returns a type mapping annotation.
func Type(name string) Annotation {
	return Annotation{Type: name}
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
	if ant.Type != "" {
		a.Type = ant.Type
	}
	if ant.Skip {
		a.Skip = true
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
