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

type converter struct {
	toEntConversion    string
	toEntConstructor   protogen.GoIdent
	toProtoConversion  string
	toProtoConstructor protogen.GoIdent
}

func (g *serviceGenerator) newConverter(fld *entproto.FieldMappingDescriptor) (*converter, error) {
	out := &converter{}
	pbd := fld.PbFieldDescriptor
	switch pbd.GetType() {
	case dpb.FieldDescriptorProto_TYPE_BOOL, dpb.FieldDescriptorProto_TYPE_STRING:
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
		enumName := fld.PbFieldDescriptor.GetEnumType().GetName()
		method := fmt.Sprintf("toProto%s_%s", g.typeName, enumName)
		out.toProtoConstructor = g.file.GoImportPath.Ident(method)
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
	case efld.Type.Numeric():
		out.toEntConversion = efld.Type.String()
	case efld.IsTime():
		out.toEntConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ExtractTime")
	case efld.IsEnum():
		enumName := fld.PbFieldDescriptor.GetEnumType().GetName()
		method := fmt.Sprintf("toEnt%s_%s", g.typeName, enumName)
		out.toEntConstructor = g.file.GoImportPath.Ident(method)
	case efld.IsUUID():
		out.toEntConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("MustBytesToUUID")
	default:
		return nil, fmt.Errorf("entproto: no mapping to ent field type %q", efld.Type.ConstName())
	}
	return out, nil
}

func (g *serviceGenerator) renderToProto(fc *converter, ident string) string {
	var left, right string
	if fc.toProtoConstructor.GoName != "" {
		left += g.QualifiedGoIdent(fc.toProtoConstructor) + "("
		right += ")"
	}
	if fc.toProtoConversion != "" {
		left += fc.toProtoConversion + "("
		right += ")"
	}
	return left + ident + right
}

func (g *serviceGenerator) renderToEnt(fc *converter, ident string) string {
	var left, right string
	if fc.toEntConstructor.GoName != "" {
		left += g.QualifiedGoIdent(fc.toEntConstructor) + "("
		right += ")"
	}
	if fc.toEntConversion != "" {
		left += fc.toEntConversion + "("
		right += ")"
	}
	return left + ident + right
}
