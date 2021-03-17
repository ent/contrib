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

package main

import (
	"fmt"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc/gen"
	"google.golang.org/protobuf/compiler/protogen"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

func (g *serviceGenerator) castToProtoFunc(fld *entproto.FieldMappingDescriptor) (interface{}, error) {
	// TODO(rotemtam): don't wrap if the ent type == the pb type
	pbd := fld.PbFieldDescriptor
	switch pbd.GetType() {
	case dpb.FieldDescriptorProto_TYPE_INT32:
		return "int32", nil
	case dpb.FieldDescriptorProto_TYPE_INT64:
		return "int64", nil
	case dpb.FieldDescriptorProto_TYPE_STRING:
		return "string", nil
	case dpb.FieldDescriptorProto_TYPE_UINT32:
		return "uint32", nil
	case dpb.FieldDescriptorProto_TYPE_UINT64:
		return "uint64", nil
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		ident := g.pbEnumIdent(fld)
		methodName := "toProto" + ident.GoName
		return methodName, nil
	case dpb.FieldDescriptorProto_TYPE_MESSAGE:
		if name := pbd.GetMessageType().GetFullyQualifiedName(); name != "google.protobuf.Timestamp" {
			return nil, fmt.Errorf("entproto: no mapping for pb message type %q", pbd.GetFullyQualifiedName())
		}
		newTS := protogen.GoImportPath("google.golang.org/protobuf/types/known/timestamppb").Ident("New")
		return newTS, nil
	default:
		return nil, fmt.Errorf("entproto: no mapping for pb field type %q", pbd.GetType())
	}
}

func (g *serviceGenerator) castToEntFunc(fd *entproto.FieldMappingDescriptor) (interface{}, error) {
	var fld *gen.Field
	if fd.IsEdgeField {
		fld = fd.EntEdge.Type.ID
	} else {
		fld = fd.EntField
	}
	switch {
	case fld.IsBool(), fld.IsBytes(), fld.IsString(), fld.Type.Numeric():
		return fld.Type.String(), nil
	case fld.IsTime():
		return protogen.GoImportPath("entgo.io/contrib/entproto").Ident("ExtractTime"), nil
	case fld.IsEnum():
		ident := g.pbEnumIdent(fd)
		methodName := "toEnt" + ident.GoName
		return methodName, nil
	// case field.TypeJSON:
	// case field.TypeUUID:
	// case field.TypeOther:
	default:
		return nil, fmt.Errorf("entproto: no mapping to ent field type %q", fld.Type.ConstName())
	}
}
