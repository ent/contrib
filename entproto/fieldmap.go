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

package entproto

import (
	"fmt"
	"strings"
	"time"

	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/types/known/timestamppb"
)

// FieldMap returns a FieldMap containing descriptors of all of the mappings between the ent schema field
// and the protobuf message's field descriptors.
func (a *Adapter) FieldMap(schemaName string) (FieldMap, error) {
	bt, err := extractGenTypeByName(a.graph, schemaName)
	if err != nil {
		return nil, err
	}
	md, err := a.GetMessageDescriptor(schemaName)
	if err != nil {
		return nil, err
	}
	return mapFields(bt, md)
}

// FieldMap contains a mapping between the field's name in the ent schema and a FieldMappingDescriptor
type FieldMap map[string]*FieldMappingDescriptor

// Fields returns the FieldMappingDescriptor for all of the fields of the schema.
func (m FieldMap) Fields() []*FieldMappingDescriptor {
	var out []*FieldMappingDescriptor
	for _, f := range m {
		if !f.IsEdgeField {
			out = append(out, f)
		}
	}
	return out
}

// ID returns the FieldMappingDescriptor for the ID field of the schema.
func (m FieldMap) ID() *FieldMappingDescriptor {
	for _, f := range m {
		if f.IsIDField {
			return f
		}
	}
	return nil
}

// Edges returns the FieldMappingDescriptor for all of the edge fields of the schema.
func (m FieldMap) Edges() []*FieldMappingDescriptor {
	var out []*FieldMappingDescriptor
	for _, f := range m {
		if f.IsEdgeField {
			out = append(out, f)
		}
	}
	return out
}

// FieldMappingDescriptor describes the mapping from a protobuf field descriptor to an ent Schema field
type FieldMappingDescriptor struct {
	EntField          *gen.Field
	PbFieldDescriptor *desc.FieldDescriptor
	ToEntExpr         string
	IsEdgeField       bool
	IsIDField         bool
}

func (d *FieldMappingDescriptor) PbStructField() string {
	return strings.Title(camel(d.PbFieldDescriptor.GetName()))
}

func mapFields(entType *gen.Type, pbType *desc.MessageDescriptor) (FieldMap, error) {
	m := make(map[string]*FieldMappingDescriptor)
	for _, fld := range pbType.GetFields() {
		isID := pascal(fld.GetName()) == pascal(entType.ID.Name)
		var isEdge bool
		for _, edg := range entType.Edges {
			if fld.GetName() == edg.Name {
				isEdge = true
				break
			}
		}
		var ef *gen.Field
		if !isEdge {
			enf, err := extractEntFieldByName(entType, fld.GetName())
			if err != nil {
				return nil, err
			}
			ef = enf
		}
		m[fld.GetName()] = &FieldMappingDescriptor{
			PbFieldDescriptor: fld,
			IsIDField:         isID,
			IsEdgeField:       isEdge,
			EntField:          ef,
		}
	}
	return m, nil
}

func extractEntFieldByName(entType *gen.Type, name string) (*gen.Field, error) {
	for _, fld := range allFields(entType) {
		if fld.Name == name {
			return fld, nil
		}
	}
	return nil, fmt.Errorf("entproto: could not find find %q in %q", name, entType.Name)
}

func allFields(t *gen.Type) []*gen.Field {
	var all []*gen.Field
	all = append(all, t.ID)
	all = append(all, t.Fields...)
	return all
}

func ExtractTime(t *timestamppb.Timestamp) time.Time {
	return t.AsTime()
}
