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
	"encoding"
	"fmt"
	"reflect"
	"strings"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"github.com/jhump/protoreflect/desc"
	"google.golang.org/protobuf/compiler/protogen"
	dpb "google.golang.org/protobuf/types/descriptorpb"
)

var (
	binaryMarshallerUnmarshallerType = reflect.TypeOf((*BinaryMarshallerUnmarshaller)(nil)).Elem()
)

type BinaryMarshallerUnmarshaller interface {
	encoding.BinaryMarshaler
	encoding.BinaryUnmarshaler
}

type converter struct {
	toEntConversion              string
	toEntScannerConversion       string
	toEntConstructor             protogen.GoIdent
	toEntMarshallerConstructor   protogen.GoIdent
	toEntScannerConstructor      protogen.GoIdent
	toEntModifier                string
	toProtoConversion            string
	toProtoConstructor           protogen.GoIdent
	toProtoMarshallerConstructor protogen.GoIdent
	toProtoValuer                string
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
		} else if err := convertPbMessageType(pbd.GetMessageType(), fld.EntField, out); err != nil {
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
	case implements(efld.Type.RType, binaryMarshallerUnmarshallerType) && efld.HasGoType():
		// Ident returned from ent already has the packagename prefixed. Strip it since `g.QualifiedGoIdent`
		// adds it back.
		split := strings.Split(efld.Type.Ident, ".")
		out.toEntMarshallerConstructor = protogen.GoImportPath(efld.Type.PkgPath).Ident(split[1])
	case efld.Type.ValueScanner():
		switch {
		case efld.HasGoType():
			// Ident returned from ent already has the packagename prefixed. Strip it since `g.QualifiedGoIdent`
			// adds it back.
			split := strings.Split(efld.Type.Ident, ".")
			out.toEntScannerConstructor = protogen.GoImportPath(efld.Type.PkgPath).Ident(split[1])
		case efld.IsBool():
			out.toEntScannerConversion = "bool"
		case efld.IsBytes():
			out.toEntScannerConversion = "[]byte"
		case efld.IsString():
			out.toEntScannerConversion = "string"
		}
	case efld.IsBool(), efld.IsBytes(), efld.IsString():
	case efld.Type.Numeric():
		out.toEntConversion = efld.Type.String()
	case efld.IsTime():
		out.toEntConstructor = protogen.GoImportPath("entgo.io/contrib/entproto/runtime").Ident("ExtractTime")
	case efld.IsEnum():
		enumName := fld.PbFieldDescriptor.GetEnumType().GetName()
		method := fmt.Sprintf("toEnt%s_%s", g.entType.Name, enumName)
		out.toEntConstructor = g.file.GoImportPath.Ident(method)
	default:
		return nil, fmt.Errorf("entproto: no mapping to ent field type %q", efld.Type.ConstName())
	}
	return out, nil
}

// Supported value scanner types (https://golang.org/pkg/database/sql/driver/#Value): [int64, float64, bool, []byte, string, time.Time]
func basicTypeConversion(md *desc.FieldDescriptor, entField *gen.Field, conv *converter) error {
	switch md.GetType() {
	case dpb.FieldDescriptorProto_TYPE_BOOL:
		if entField.Type.Valuer() {
			conv.toProtoValuer = "bool"
		}
	case dpb.FieldDescriptorProto_TYPE_STRING:
		if entField.Type.Valuer() {
			conv.toProtoValuer = "string"
		}
	case dpb.FieldDescriptorProto_TYPE_BYTES:
		if implements(entField.Type.RType, binaryMarshallerUnmarshallerType) {
			// Ident returned from ent already has the packagename prefixed. Strip it since `g.QualifiedGoIdent`
			// adds it back.
			split := strings.Split(entField.Type.Ident, ".")
			conv.toProtoMarshallerConstructor = protogen.GoImportPath(entField.Type.PkgPath).Ident(split[1])
		} else if entField.Type.Valuer() {
			conv.toProtoValuer = "[]byte"
		}
	case dpb.FieldDescriptorProto_TYPE_INT32:
		if entField.Type.String() != "int32" {
			conv.toProtoConversion = "int32"
		}
	case dpb.FieldDescriptorProto_TYPE_INT64:
		if entField.Type.Valuer() {
			conv.toProtoValuer = "int64"
		} else if entField.Type.String() != "int64" {
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

func convertPbMessageType(md *desc.MessageDescriptor, entField *gen.Field, conv *converter) error {
	switch {
	case md.GetFullyQualifiedName() == "google.protobuf.Timestamp":
		conv.toProtoConstructor = protogen.GoImportPath("google.golang.org/protobuf/types/known/timestamppb").Ident("New")
	case isWrapperType(md):
		fqn := md.GetFullyQualifiedName()
		typ := strings.Split(fqn, ".")[2]
		constructor := strings.TrimSuffix(typ, "Value")
		conv.toProtoConstructor = protogen.GoImportPath("google.golang.org/protobuf/types/known/wrapperspb").Ident(constructor)

		goType := wrapperPrimitives[fqn]
		if entField.Type.Valuer() {
			conv.toProtoValuer = goType
		} else if entField.Type.String() != goType {
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
	case conv.toEntMarshallerConstructor.GoName != "":
		return fmt.Sprintf(`%s, err := %s.MarshalBinary()
		if err != nil {
			return nil, err
		}
		`, varName, ident)
	case conv.toProtoValuer != "" && conv.toProtoConstructor.GoName != "":
		// Returns logic to extract a valuers value, cast it to the appropriate go type, and pass it to the proto constructor.
		return fmt.Sprintf(`%sValue, err := %s.Value()
		if err != nil {
			return nil, err
		}

		%sTyped, ok := %sValue.(%s)
		if !ok {
			return nil, %s("casting value to %s")
		}

		%s := %s(%sTyped)
		`, varName, ident, varName, varName, conv.toProtoValuer, "%(newError)", conv.toProtoValuer, varName, g.QualifiedGoIdent(conv.toProtoConstructor), varName)
	case conv.toProtoValuer != "":
		// Returns logic to extract a valuers value and cast it to the appropriate go type
		return fmt.Sprintf(`%sValue, err := %s.Value()
		if err != nil {
			return nil, err
		}

		%s, ok := %sValue.(%s)
		if !ok {
			return nil, %s("casting value to %s")
		}
		`, varName, ident, varName, varName, conv.toProtoValuer, "%(newError)", conv.toProtoValuer)
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
	case conv.toEntMarshallerConstructor.GoName != "":
		return fmt.Sprintf(`var %s %s
		if err := (&%s).UnmarshalBinary(%s); err != nil {
			return nil, %s
		}
		`, varName, g.QualifiedGoIdent(conv.toEntMarshallerConstructor), varName, ident, `%(statusErrf)(%(invalidArgument), "invalid argument: %s", err)`)
	case conv.toEntScannerConstructor.GoName != "":
		// Returns: varName, err := EntConstructor(protoMember) with error handler
		return fmt.Sprintf(`%s := %s{}
		if err := (&%s).Scan(%s); err != nil {
			return nil, %s
		}
		`, varName, g.QualifiedGoIdent(conv.toEntScannerConstructor), varName, ident, `%(statusErrf)(%(invalidArgument), "invalid argument: %s", err)`)
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

func implements(r *field.RType, typ reflect.Type) bool {
	if r == nil {
		return false
	}
	n := typ.NumMethod()
	for i := 0; i < n; i++ {
		m0 := typ.Method(i)
		m1, ok := r.Methods[m0.Name]
		if !ok || len(m1.In) != m0.Type.NumIn() || len(m1.Out) != m0.Type.NumOut() {
			return false
		}
		in := m0.Type.NumIn()
		for j := 0; j < in; j++ {
			if !m1.In[j].TypeEqual(m0.Type.In(j)) {
				return false
			}
		}
		out := m0.Type.NumOut()
		for j := 0; j < out; j++ {
			if !m1.Out[j].TypeEqual(m0.Type.Out(j)) {
				return false
			}
		}
	}
	return true
}
