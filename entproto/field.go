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

package entproto

import (
	"fmt"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"github.com/go-viper/mapstructure/v2"
	"google.golang.org/protobuf/types/descriptorpb"
)

const FieldAnnotation = "ProtoField"

type FieldOption func(*pbfield)

func Field(num int, options ...FieldOption) schema.Annotation {
	f := pbfield{Number: num}
	for _, apply := range options {
		apply(&f)
	}
	return f
}

type pbfield struct {
	Number   int
	Type     descriptorpb.FieldDescriptorProto_Type
	TypeName string
}

func (f pbfield) Name() string {
	return FieldAnnotation
}

// Type overrides the default mapping between ent types and protobuf types.
// Example:
//
//	field.Uint8("custom_pb").
//		Annotations(
//			entproto.Field(2,
//				entproto.Type(descriptorpb.FieldDescriptorProto_TYPE_UINT64),
//			),
//		)
func Type(typ descriptorpb.FieldDescriptorProto_Type) FieldOption {
	return func(p *pbfield) {
		p.Type = typ
	}
}

// TypeName sets the pb descriptors type name, needed if the Type attribute is TYPE_ENUM or TYPE_MESSAGE.
func TypeName(n string) FieldOption {
	return func(p *pbfield) {
		p.TypeName = n
	}
}

func extractFieldAnnotation(fld *gen.Field) (*pbfield, error) {
	annot, ok := fld.Annotations[FieldAnnotation]
	if !ok {
		return nil, fmt.Errorf("entproto: field %q does not have an entproto.Field annnoation", fld.Name)
	}

	var out pbfield
	err := mapstructure.Decode(annot, &out)
	if err != nil {
		return nil, fmt.Errorf("entproto: unable to decode entproto.Field annotation for field %q: %w",
			fld.Name, err)
	}

	return &out, nil
}

func extractEdgeAnnotation(edge *gen.Edge) (*pbfield, error) {
	annot, ok := edge.Annotations[FieldAnnotation]
	if !ok {
		return nil, fmt.Errorf("entproto: edge %q does not have an entproto.Field annotation", edge.Name)
	}

	var out pbfield
	err := mapstructure.Decode(annot, &out)
	if err != nil {
		return nil, fmt.Errorf("entproto: unable to decode entproto.Field annotation for field %q: %w",
			edge.Name, err)
	}

	return &out, nil
}
