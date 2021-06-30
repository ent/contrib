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
	"strings"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc/gen"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/compiler/protogen"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

type converter struct {
	toEntConversion       string
	toEntConstructor      protogen.GoIdent
	toEntTryConstructor   protogen.GoIdent
	toEntModifier         string
	toProtoConversion     string
	toProtoConstructor    protogen.GoIdent
	toProtoTryConstructor protogen.GoIdent
}

func (g *serviceGenerator) newConverter(fld *entproto.FieldMappingDescriptor) (*converter, error) {
	out := &converter{}
	pbd := fld.PbFieldDescriptor
	switch pbd.GetType() {
	case dpb.FieldDescriptorProto_TYPE_BOOL, dpb.FieldDescriptorProto_TYPE_STRING,
		dpb.FieldDescriptorProto_TYPE_BYTES, dpb.FieldDescriptorProto_TYPE_INT32,
		dpb.FieldDescriptorProto_TYPE_INT64, dpb.FieldDescriptorProto_TYPE_UINT32,
		dpb.FieldDescriptorProto_TYPE_UINT64:
		if err := basicTypeConversion(fld.PbFieldDescriptor, fld.EntField, out); err != nil {
			return nil, err
		}
	case dpb.FieldDescriptorProto_TYPE_ENUM:
		enumName := fld.PbFieldDescriptor.GetEnumType().GetName()
		method := fmt.Sprintf("toProto%s_%s", g.entType.Name, enumName)
		out.toProtoConstructor = g.file.GoImportPath.Ident(method)
	case dpb.FieldDescriptorProto_TYPE_MESSAGE:
		if fld.IsEdgeField {
			if err := basicTypeConversion(fld.EdgeIDPbStructFieldDesc(), fld.EntEdge.Type.ID, out); err != nil {
				return nil, err
			}
		} else if err := convertPbMessageType(pbd.GetMessageType(), fld.EntField.Type.String(), out); err != nil {
			return nil, err
		}
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
		method := fmt.Sprintf("toEnt%s_%s", g.entType.Name, enumName)
		out.toEntConstructor = g.file.GoImportPath.Ident(method)
	case efld.IsUUID():
		out.toEntTryConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("BytesToUUID")
	default:
		return nil, fmt.Errorf("entproto: no mapping to ent field type %q", efld.Type.ConstName())
	}
	return out, nil
}

func basicTypeConversion(md *desc.FieldDescriptor, entField *gen.Field, conv *converter) error {
	switch md.GetType() {
	case dpb.FieldDescriptorProto_TYPE_BOOL, dpb.FieldDescriptorProto_TYPE_STRING:
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		if entField.IsUUID() {
			conv.toProtoTryConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("UUIDToBytes")
		}
	case dpb.FieldDescriptorProto_TYPE_INT32:
		if entField.Type.String() != "int32" {
			conv.toProtoConversion = "int32"
		}
	case dpb.FieldDescriptorProto_TYPE_INT64:
		if entField.Type.String() != "int64" {
			conv.toProtoConversion = "int64"
		}
	case dpb.FieldDescriptorProto_TYPE_UINT32:
		if entField.Type.String() != "uint32" {
			conv.toProtoConversion = "uint32"
		}
	case dpb.FieldDescriptorProto_TYPE_UINT64:
		if entField.Type.String() != "uint64" {
			conv.toProtoConversion = "uint64"
		}
	}
	return nil
}

func convertPbMessageType(md *desc.MessageDescriptor, entFieldType string, conv *converter) error {
	switch {
	case md.GetFullyQualifiedName() == "google.protobuf.Timestamp":
		conv.toProtoConstructor = protogen.GoImportPath("google.golang.org/protobuf/types/known/timestamppb").Ident("New")
	case isWrapperType(md):
		fqn := md.GetFullyQualifiedName()
		typ := strings.Split(fqn, ".")[2]
		constructor := strings.TrimSuffix(typ, "Value")
		conv.toProtoConstructor = protogen.GoImportPath("google.golang.org/protobuf/types/known/wrapperspb").Ident(constructor)
		if goType := wrapperPrimitives[fqn]; entFieldType != goType {
			conv.toProtoConversion = goType
		}
		conv.toEntModifier = ".GetValue()"
	default:
		return fmt.Errorf("entproto: no mapping for pb field type %q", md.GetFullyQualifiedName())
	}
	return nil
}

func (g *serviceGenerator) renderToProto(conv *converter, varName, ident string) string {
	// If we are going to cast to a GoType, wrap it first. i.e. int64(entMember)
	if conv.toProtoConversion != "" {
		ident = fmt.Sprintf("%s(%s)", conv.toProtoConversion, ident)
	}

	switch {
	case conv.toProtoTryConstructor.GoName != "":
		// Returns: varName, err := ProtcoConstructor(entMember) with error handler
		return fmt.Sprintf(`%s, err := %s(%s)
		if err != nil {
			return nil, err
		}
		`, varName, g.QualifiedGoIdent(conv.toProtoTryConstructor), ident)
	case conv.toProtoConstructor.GoName != "":
		// Returns: varName := ProtcoConstructor(entMember)
		return fmt.Sprintf("%s := %s(%s)", varName, g.QualifiedGoIdent(conv.toProtoConstructor), ident)
	default:
		// Returns: varName := entMember
		return fmt.Sprintf("%s := %s", varName, ident)
	}
}

func (g *serviceGenerator) renderToEnt(conv *converter, varName, ident string) string {
	// Attach a modifer to ident (i.e. .GetValue())
	if conv.toEntModifier != "" {
		ident += conv.toEntModifier
	}

	switch {
	case conv.toEntTryConstructor.GoName != "":
		// Returns: varName, err := EntConstructor(protoMember) with error handler
		return fmt.Sprintf(`%s, err := %s(%s)
		if err != nil {
			return nil, %s
		}
		`, varName, g.QualifiedGoIdent(conv.toEntTryConstructor), ident, `%(statusErrf)(%(invalidArgument), "invalid argument: %s", err)`)
	case conv.toEntConstructor.GoName != "":
		// Returns: varName := EntConstructor(protoMember)
		return fmt.Sprintf("%s := %s(%s)", varName, g.QualifiedGoIdent(conv.toEntConstructor), ident)
	case conv.toEntConversion != "":
		// Returns: varName := EntConversion(protoMember)
		return fmt.Sprintf("%s := %s(%s)", varName, conv.toEntConversion, ident)
	default:
		// Returns: varName := protoMember
		return fmt.Sprintf("%s := %s", varName, ident)
	}
}

func isWrapperType(md *desc.MessageDescriptor) bool {
	_, ok := wrapperPrimitives[md.GetFullyQualifiedName()]
	return ok
}

var wrapperPrimitives = map[string]string{
	"google.protobuf.DoubleValue": "float64",
	"google.protobuf.FloatValue":  "float32",
	"google.protobuf.Int64Value":  "int64",
	"google.protobuf.UInt64Value": "uint64",
	"google.protobuf.Int32Value":  "int32",
	"google.protobuf.UInt32Value": "uint32",
	"google.protobuf.BoolValue":   "bool",
	"google.protobuf.StringValue": "string",
	"google.protobuf.BytesValue":  "[]byte",
}
