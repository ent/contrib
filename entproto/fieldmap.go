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
	"sort"
	"strings"

	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/desc"
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
	return a.mapFields(bt, md)
}

// FieldMap contains a mapping between the field's name in the ent schema and a FieldMappingDescriptor.
type FieldMap map[string]*FieldMappingDescriptor

// Fields returns the FieldMappingDescriptor for all of the fields of the schema. Items are sorted alphabetically
// on pb field name.
func (m FieldMap) Fields() []*FieldMappingDescriptor {
	var out []*FieldMappingDescriptor
	for _, f := range m {
		if !f.IsEdgeField {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].PbStructField() < out[j].PbStructField()
	})
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

// Edges returns the FieldMappingDescriptor for all of the edge fields of the schema. Items are sorted alphabetically
// on pb field name.
func (m FieldMap) Edges() []*FieldMappingDescriptor {
	var out []*FieldMappingDescriptor
	for _, f := range m {
		if f.IsEdgeField {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].PbStructField() < out[j].PbStructField()
	})

	return out
}

func (m FieldMap) Enums() []*FieldMappingDescriptor {
	var out []*FieldMappingDescriptor
	for _, f := range m {
		if f.IsEnumFIeld {
			out = append(out, f)
		}
	}
	sort.Slice(out, func(i, j int) bool {
		return out[i].PbStructField() < out[j].PbStructField()
	})
	return out
}

// FieldMappingDescriptor describes the mapping from a protobuf field descriptor to an ent Schema field
type FieldMappingDescriptor struct {
	EntField          *gen.Field
	EntEdge           *gen.Edge
	PbFieldDescriptor *desc.FieldDescriptor
	IsEdgeField       bool
	IsIDField         bool
	IsEnumFIeld       bool
	ReferencedPbType  *desc.MessageDescriptor
}

func (d *FieldMappingDescriptor) PbStructField() string {
	return camelCase(d.PbFieldDescriptor.GetName())
}

func (d *FieldMappingDescriptor) EdgeIDPbStructField() string {
	return camelCase(d.EntEdge.Ref.Type.ID.Name)
}

func (d *FieldMappingDescriptor) EdgeIDPbStructFieldDesc() *desc.FieldDescriptor {
	field := strings.Title(camel(d.EntEdge.Ref.Type.ID.Name))
	return d.ReferencedPbType.FindFieldByName(snake(field))
}

func (a *Adapter) mapFields(entType *gen.Type, pbType *desc.MessageDescriptor) (FieldMap, error) {
	m := make(map[string]*FieldMappingDescriptor)
	for _, fld := range pbType.GetFields() {
		fd := &FieldMappingDescriptor{
			PbFieldDescriptor: fld,
			IsIDField:         pascal(fld.GetName()) == pascal(entType.ID.Name),
			IsEnumFIeld:       fld.GetEnumType() != nil,
		}
		for _, edg := range entType.Edges {
			if fld.GetName() == edg.Name {
				fd.IsEdgeField = true
				break
			}
		}
		if fd.IsEdgeField {
			edg, err := extractEntEdgeByName(entType, fld.GetName())
			if err != nil {
				return nil, err
			}
			fd.EntEdge = edg
			referenced, err := a.GetMessageDescriptor(edg.Type.Name)
			if err != nil {
				return nil, err
			}
			fd.ReferencedPbType = referenced
		} else {
			enf, err := extractEntFieldByName(entType, fld.GetName())
			if err != nil {
				return nil, err
			}
			fd.EntField = enf
		}
		m[fld.GetName()] = fd
	}
	return m, nil
}

func extractEntFieldByName(entType *gen.Type, name string) (*gen.Field, error) {
	if name == entType.ID.Name {
		return entType.ID, nil
	}
	for _, fld := range entType.Fields {
		if fld.Name == name {
			return fld, nil
		}
	}
	return nil, fmt.Errorf("entproto: could not find field %q in %q", name, entType.Name)
}

func extractEntEdgeByName(entType *gen.Type, name string) (*gen.Edge, error) {
	for _, edg := range entType.Edges {
		if edg.Name == name {
			return edg, nil
		}
	}
	return nil, fmt.Errorf("entproto: could not find find edge %q in %q", name, entType.Name)
}

// Is c an ASCII lower-case letter?
func isASCIILower(c byte) bool {
	return 'a' <= c && c <= 'z'
}

// Is c an ASCII digit?
func isASCIIDigit(c byte) bool {
	return '0' <= c && c <= '9'
}

// camelCase was copied from https://github.com/golang/protobuf/blob/v1.5.2/protoc-gen-go/generator/generator.go#L2648
// camelCase returns the CamelCased name.
// If there is an interior underscore followed by a lower case letter,
// drop the underscore and convert the letter to upper case.
// There is a remote possibility of this rewrite causing a name collision,
// but it's so remote we're prepared to pretend it's nonexistent - since the
// C++ generator lowercases names, it's extremely unlikely to have two fields
// with different capitalizations.
// In short, _my_field_name_2 becomes XMyFieldName_2.
func camelCase(s string) string {
	if s == "" {
		return ""
	}
	t := make([]byte, 0, 32)
	i := 0
	if s[0] == '_' {
		// Need a capital letter; drop the '_'.
		t = append(t, 'X')
		i++
	}
	// Invariant: if the next letter is lower case, it must be converted
	// to upper case.
	// That is, we process a word at a time, where words are marked by _ or
	// upper case letter. Digits are treated as words.
	for ; i < len(s); i++ {
		c := s[i]
		if c == '_' && i+1 < len(s) && isASCIILower(s[i+1]) {
			continue // Skip the underscore in s.
		}
		if isASCIIDigit(c) {
			t = append(t, c)
			continue
		}
		// Assume we have a letter now - if not, it's a bogus identifier.
		// The next word is a sequence of characters that must start upper case.
		if isASCIILower(c) {
			c ^= ' ' // Make it a capital letter.
		}
		t = append(t, c) // Guaranteed not lower case.
		// Accept lower case sequence that follows.
		for i+1 < len(s) && isASCIILower(s[i+1]) {
			i++
			t = append(t, s[i])
		}
	}
	return string(t)
}
