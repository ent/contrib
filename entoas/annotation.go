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

package entoas

import (
	"encoding/json"

	"entgo.io/contrib/entoas/serialization"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	yaml "github.com/go-faster/yaml"
	"github.com/ogen-go/ogen"
)

const XOgenOperationGroup = "x-ogen-operation-group"

type (
	// Annotation annotates fields and edges with metadata for spec generation.
	Annotation struct {
		// Groups holds the serialization groups to use on this field / edge.
		Groups serialization.Groups
		// OpenAPI Specification example value for a schema field.
		Example interface{}
		// OpenAPI Specification schema to use for a schema field.
		Schema *ogen.Schema
		// Create has meta information about a creation operation.
		Create OperationConfig
		// Read has meta information about a read operation.
		Read OperationConfig
		// Update has meta information about an update operation.
		Update OperationConfig
		// Delete has meta information about a delete operation.
		Delete OperationConfig
		// List has meta information about a list operation.
		List OperationConfig
		// ReadOnly specifies that the field/edge is read only (no create/update parameter)
		ReadOnly bool
		// Skip specifies that the field will be ignored in spec.
		Skip bool
		// Extensions has map of OpenApi extenions
		Extensions ogen.Extensions
	}
	// OperationConfig holds meta information about a REST operation.
	OperationConfig struct {
		Policy     Policy
		Groups     serialization.Groups
		Extensions ogen.Extensions
	}
	// OpenApiExtension holds meta information about OpenApi extension
	OpenApiExtension struct {
		Name  string
		Value string
	}
	// OperationConfigOption allows managing OperationConfig using functional arguments.
	OperationConfigOption func(*OperationConfig)
)

// Groups returns a OperationConfigOption that adds the given serialization groups to a OperationConfig.
func Groups(gs ...string) Annotation {
	return Annotation{Groups: gs}
}

// OperationGroups returns a OperationConfigOption that adds the given serialization groups to a OperationConfig.
func OperationGroups(gs ...string) OperationConfigOption {
	return func(c *OperationConfig) { c.Groups = gs }
}

func OperationExtentions(ext ...OpenApiExtension) OperationConfigOption {
	return func(c *OperationConfig) {
		if c.Extensions == nil {
			c.Extensions = ogen.Extensions{}
		}

		for _, e := range ext {
			c.Extensions[e.Name] = yaml.Node{Kind: yaml.ScalarNode, Value: e.Value}
		}
	}
}

func Extensions(ext ...OpenApiExtension) Annotation {
	exts := ogen.Extensions{}

	for _, e := range ext {
		exts[e.Name] = yaml.Node{Kind: yaml.ScalarNode, Value: e.Value}
	}

	return Annotation{Extensions: exts}
}

func OperationExtentionGroup(group string) OperationConfigOption {
	return OperationExtentions(OpenApiExtensionOperationGroup(group))
}

func ExtentionOperationGroup(group string) Annotation {
	return Extensions(OpenApiExtensionOperationGroup(group))
}

func OpenApiExtensionOperationGroup(group string) OpenApiExtension {
	return OpenApiExtension{
		Name:  XOgenOperationGroup,
		Value: group,
	}
}

// OperationPolicy returns a OperationConfigOption that sets the Policy of a OperationConfig to the given one.
func OperationPolicy(p Policy) OperationConfigOption {
	return func(c *OperationConfig) { c.Policy = p }
}

// Example returns an example annotation.
func Example(v interface{}) Annotation { return Annotation{Example: v} }

// Schema returns a Schema annotation.
func Schema(s *ogen.Schema) Annotation { return Annotation{Schema: s} }

// CreateOperation returns a create operation annotation.
func CreateOperation(opts ...OperationConfigOption) Annotation {
	return Annotation{Create: operationsConfig(opts)}
}

// ReadOperation returns a read operation annotation.
func ReadOperation(opts ...OperationConfigOption) Annotation {
	return Annotation{Read: operationsConfig(opts)}
}

// UpdateOperation returns an update operation annotation.
func UpdateOperation(opts ...OperationConfigOption) Annotation {
	return Annotation{Update: operationsConfig(opts)}
}

// DeleteOperation returns a delete operation annotation.
func DeleteOperation(opts ...OperationConfigOption) Annotation {
	return Annotation{Delete: operationsConfig(opts)}
}

// ListOperation returns a list operation annotation.
func ListOperation(opts ...OperationConfigOption) Annotation {
	return Annotation{List: operationsConfig(opts)}
}

// ReadOnly returns a read only field/edge annotation
func ReadOnly(readonly bool) Annotation {
	return Annotation{ReadOnly: readonly}
}

// Skip returns a skip field annotation
func Skip(skip bool) Annotation {
	return Annotation{Skip: skip}
}

func operationsConfig(opts []OperationConfigOption) OperationConfig {
	c := OperationConfig{}
	for _, opt := range opts {
		opt(&c)
	}
	return c
}

// Name implements schema.Annotation interface.
func (Annotation) Name() string { return "EntOAS" }

// Merge implements ent.Merger interface.
func (a Annotation) Merge(o schema.Annotation) schema.Annotation {
	var ant Annotation
	switch o := o.(type) {
	case Annotation:
		ant = o
	case *Annotation:
		if o != nil {
			ant = *o
		}
	default:
		return a
	}
	if ant.Example != nil {
		a.Example = ant.Example
	}
	if ant.Schema != nil {
		a.Schema = ant.Schema
	}
	a.Create.merge(ant.Create)
	a.Read.merge(ant.Read)
	a.Update.merge(ant.Update)
	a.Delete.merge(ant.Delete)
	a.List.merge(ant.List)
	if ant.ReadOnly {
		a.ReadOnly = true
	}
	if ant.Skip {
		a.Skip = true
	}
	return a
}

func (op *OperationConfig) merge(other OperationConfig) {
	if other.Policy != PolicyNone {
		op.Policy = other.Policy
	}
	if other.Groups != nil {
		op.Groups = other.Groups
	}
}

// Decode from ent.
func (a *Annotation) Decode(o interface{}) error {
	buf, err := json.Marshal(o)
	if err != nil {
		return err
	}
	return json.Unmarshal(buf, a)
}

// SchemaAnnotation returns the Annotation on the given gen.Type.
func SchemaAnnotation(n *gen.Type) (*Annotation, error) { return annotation(n.Annotations) }

// FieldAnnotation returns the Annotation on the given gen.Field.
func FieldAnnotation(f *gen.Field) (*Annotation, error) { return annotation(f.Annotations) }

// EdgeAnnotation returns the Annotation on the given gen.Edge.
func EdgeAnnotation(e *gen.Edge) (*Annotation, error) { return annotation(e.Annotations) }

// annotation decodes the Annotation from the given gen.Annotations.
func annotation(as gen.Annotations) (*Annotation, error) {
	ant := &Annotation{}
	if as != nil && as[ant.Name()] != nil {
		if err := ant.Decode(as[ant.Name()]); err != nil {
			return nil, err
		}
	}
	return ant, nil
}

var (
	_ schema.Annotation = (*Annotation)(nil)
	_ schema.Merger     = (*Annotation)(nil)
)
