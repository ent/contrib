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
	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema/field"
	"google.golang.org/protobuf/types/descriptorpb"
)

var typeMap = map[field.Type]typeConfig{
	field.TypeBool:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_BOOL, optionalType: "google.protobuf.BoolValue"},
	field.TypeTime:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, msgTypeName: "google.protobuf.Timestamp", optionalType: "google.protobuf.Timestamp"},
	field.TypeJSON:  {unsupported: true},
	field.TypeOther: {unsupported: true},
	field.TypeUUID:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_BYTES, optionalType: "google.protobuf.BytesValue"},
	field.TypeBytes: {pbType: descriptorpb.FieldDescriptorProto_TYPE_BYTES, optionalType: "google.protobuf.BytesValue"},
	field.TypeEnum: {pbType: descriptorpb.FieldDescriptorProto_TYPE_ENUM, namer: func(fld *gen.Field) string {
		return pascal(fld.Name)
	}},
	field.TypeString:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_STRING, optionalType: "google.protobuf.StringValue"},
	field.TypeInt:     {pbType: descriptorpb.FieldDescriptorProto_TYPE_INT32, optionalType: "google.protobuf.Int32Value"},
	field.TypeInt8:    {pbType: descriptorpb.FieldDescriptorProto_TYPE_INT32, optionalType: "google.protobuf.Int32Value"},
	field.TypeInt16:   {pbType: descriptorpb.FieldDescriptorProto_TYPE_INT32, optionalType: "google.protobuf.Int32Value"},
	field.TypeInt32:   {pbType: descriptorpb.FieldDescriptorProto_TYPE_INT32, optionalType: "google.protobuf.Int32Value"},
	field.TypeInt64:   {pbType: descriptorpb.FieldDescriptorProto_TYPE_INT64, optionalType: "google.protobuf.Int64Value"},
	field.TypeUint:    {pbType: descriptorpb.FieldDescriptorProto_TYPE_UINT32, optionalType: "google.protobuf.UInt32Value"},
	field.TypeUint8:   {pbType: descriptorpb.FieldDescriptorProto_TYPE_UINT32, optionalType: "google.protobuf.UInt32Value"},
	field.TypeUint16:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_UINT32, optionalType: "google.protobuf.UInt32Value"},
	field.TypeUint32:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_UINT32, optionalType: "google.protobuf.UInt32Value"},
	field.TypeUint64:  {pbType: descriptorpb.FieldDescriptorProto_TYPE_UINT64, optionalType: "google.protobuf.UInt64Value"},
	field.TypeFloat32: {pbType: descriptorpb.FieldDescriptorProto_TYPE_FLOAT, optionalType: "google.protobuf.FloatValue"},
	field.TypeFloat64: {pbType: descriptorpb.FieldDescriptorProto_TYPE_DOUBLE, optionalType: "google.protobuf.DoubleValue"},
}

var identMap = map[string]typeConfig{
	"[]string": {pbType: descriptorpb.FieldDescriptorProto_TYPE_STRING, pbLabel: descriptorpb.FieldDescriptorProto_LABEL_REPEATED, optionalType: "google.protobuf.ListValue"},
}

type typeConfig struct {
	unsupported  bool
	pbType       descriptorpb.FieldDescriptorProto_Type
	pbLabel      descriptorpb.FieldDescriptorProto_Label
	msgTypeName  string
	optionalType string
	namer        func(fld *gen.Field) string
}
