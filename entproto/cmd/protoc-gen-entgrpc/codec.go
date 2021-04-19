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
	"google.golang.org/protobuf/compiler/protogen"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

type fieldCodec struct {
	toEntConversion    string
	toEntConstructor   protogen.GoIdent
	toProtoConversion  string
	toProtoConstructor protogen.GoIdent
}

func (g *serviceGenerator) newFieldCodec(fld *entproto.FieldMappingDescriptor) (*fieldCodec, error) {
	out := &fieldCodec{}
	pbd := fld.PbFieldDescriptor
	switch pbd.GetType() {
	case dpb.FieldDescriptorProto_TYPE_BOOL, dpb.FieldDescriptorProto_TYPE_STRING:
		break
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		if fld.EntField != nil && fld.EntField.IsUUID() {
			out.toProtoConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("MustExtractUUIDBytes")
		}
	case dpb.FieldDescriptorProto_TYPE_INT32:
		if fld.EntField.Type.String() != "int32" {
			out.toProtoConversion = "int32"
		}
	case dpb.FieldDescriptorProto_TYPE_INT64:
		if fld.EntField.Type.String() != "int64" {
			out.toProtoConversion = "int64"
		}
	case dpb.FieldDescriptorProto_TYPE_UINT32:
		if fld.EntField.Type.String() != "uint32" {
			out.toProtoConversion = "uint32"
		}
	case dpb.FieldDescriptorProto_TYPE_UINT64:
		if fld.EntField.Type.String() != "uint64" {
			out.toProtoConversion = "uint64"
		}
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		enumTypeName := fld.PbFieldDescriptor.GetEnumType().GetName()
		enumMethod := fmt.Sprintf("toProto%s_%s", g.typeName, enumTypeName)
		out.toProtoConstructor = g.file.GoImportPath.Ident(enumMethod)
	case dpb.FieldDescriptorProto_TYPE_MESSAGE:
		if fld.IsEdgeField {
			break
		}
		if name := pbd.GetMessageType().GetFullyQualifiedName(); name != "google.protobuf.Timestamp" {
			return nil, fmt.Errorf("entproto: no mapping for pb message type %q", pbd.GetFullyQualifiedName())
		}
		out.toProtoConstructor = protogen.GoImportPath("google.golang.org/protobuf/types/known/timestamppb").Ident("New")
	default:
		return nil, fmt.Errorf("entproto: no mapping for pb field type %q", pbd.GetType())
	}
	efld := fld.EntField
	if fld.IsEdgeField {
		efld = fld.EntEdge.Type.ID
	}
	switch {
	case efld.IsBool(), efld.IsBytes(), efld.IsString():
		break
	case efld.Type.Numeric():
		out.toEntConversion = efld.Type.String()
	case efld.IsTime():
		out.toEntConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ExtractTime")
	case efld.IsEnum():
		enumTypeName := fld.PbFieldDescriptor.GetEnumType().GetName()
		enumMethod := fmt.Sprintf("toEnt%s_%s", g.typeName, enumTypeName)
		out.toEntConstructor = g.file.GoImportPath.Ident(enumMethod)
	case efld.IsUUID():
		out.toEntConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("MustBytesToUUID")
	// case field.TypeJSON:
	// case field.TypeOther:
	default:
		return nil, fmt.Errorf("entproto: no mapping to ent field type %q", efld.Type.ConstName())
	}
	return out, nil
}

func (g *serviceGenerator) renderToProto(fc *fieldCodec, ident string) string {
	var out, closing string
	if fc.toProtoConstructor.GoName != "" {
		out += g.QualifiedGoIdent(fc.toProtoConstructor) + "("
		closing += ")"
	}
	if fc.toProtoConversion != "" {
		out += fc.toProtoConversion + "("
		closing += ")"
	}
	return out + ident + closing
}

func (g *serviceGenerator) renderToEnt(fc *fieldCodec, ident string) string {
	var out, closing string
	if fc.toEntConstructor.GoName != "" {
		out += g.QualifiedGoIdent(fc.toEntConstructor) + "("
		closing += ")"
	}
	if fc.toEntConversion != "" {
		out += fc.toEntConversion + "("
		closing += ")"
	}
	return out + ident + closing
}
