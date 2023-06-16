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

package entprototest

import (
	"path/filepath"
	"testing"

	"entgo.io/contrib/entproto"
	"entgo.io/ent/entc"
	"entgo.io/ent/entc/gen"
	"github.com/stretchr/testify/suite"
	"google.golang.org/protobuf/types/descriptorpb"
)

type AdapterTestSuite struct {
	suite.Suite
	adapter *entproto.Adapter
}

// This will run before before the tests in the suite are run
func (suite *AdapterTestSuite) SetupSuite() {
	graph, err := entc.LoadGraph("./ent/schema", &gen.Config{})
	if err != nil {
		suite.FailNowf("test suite init failed", "%v", err)
	}
	adapter, err := entproto.LoadAdapter(graph)
	if err != nil {
		suite.FailNowf("test suite failed", "%v", err)
	}
	suite.adapter = adapter
}

func TestSuite(t *testing.T) {
	suite.Run(t, new(AdapterTestSuite))
}

func (suite *AdapterTestSuite) TestValidMessage() {
	fd, err := suite.adapter.GetFileDescriptor("ValidMessage")
	suite.NoError(err)

	suite.Equal(filepath.Join("entpb", "entpb.proto"), fd.GetName())

	suite.Equal("entgo.io/contrib/entproto/internal/entprototest/ent/proto/entpb",
		fd.GetFileOptions().GetGoPackage())
	message := fd.FindMessage("entpb.ValidMessage")
	suite.Len(message.GetFields(), 6)

	idField := message.FindFieldByName("id")
	suite.NotNil(idField)
	suite.EqualValues(1, idField.GetNumber())

	nameField := message.FindFieldByName("name")
	suite.NotNil(nameField)
	suite.EqualValues(2, nameField.GetNumber())

	tsField := message.FindFieldByName("ts")
	suite.NotNil(nameField)
	suite.EqualValues(3, tsField.GetNumber())
	suite.EqualValues("google.protobuf.Timestamp", tsField.GetMessageType().GetFullyQualifiedName())

	uuField := message.FindFieldByName("uuid")
	suite.NotNil(uuField)

	u8Field := message.FindFieldByName("u8")
	suite.NotNil(u8Field)
	suite.EqualValues(u8Field.GetType(), descriptorpb.FieldDescriptorProto_TYPE_UINT64)

	opti8Field := message.FindFieldByName("opti8")
	suite.NotNil(opti8Field)
	suite.EqualValues("google.protobuf.Int32Value", opti8Field.GetMessageType().GetFullyQualifiedName())
}

func (suite *AdapterTestSuite) TestWktProtosDropped() {
	all := suite.adapter.AllFileDescriptors()
	_, present := all["google/protobuf/timestamp.proto"]
	suite.False(present, "wkt timestamp proto should not be included in the output descriptors")
}

func (suite *AdapterTestSuite) TestImplicitSkippedMessage() {
	_, err := suite.adapter.GetFileDescriptor("ImplicitSkippedMessage")
	suite.EqualError(err, entproto.ErrSchemaSkipped.Error())
}

func (suite *AdapterTestSuite) TestMessageWithStrings() {
	message, err := suite.adapter.GetMessageDescriptor("MessageWithStrings")
	suite.NoError(err)
	field := message.FindFieldByName("strings")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_STRING, field.GetType(), "expected repeated")
	suite.Require().True(field.IsRepeated(), "expected repeated")
}

func (suite *AdapterTestSuite) TestMessageWithInts() {
	message, err := suite.adapter.GetMessageDescriptor("MessageWithInts")
	suite.NoError(err)
	field := message.FindFieldByName("int32s")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_INT32, field.GetType(), "expected repeated")
	suite.Require().True(field.IsRepeated(), "expected repeated")
	field = message.FindFieldByName("int64s")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_INT64, field.GetType(), "expected repeated")
	suite.Require().True(field.IsRepeated(), "expected repeated")
	field = message.FindFieldByName("uint32s")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_UINT32, field.GetType(), "expected repeated")
	suite.Require().True(field.IsRepeated(), "expected repeated")
	field = message.FindFieldByName("uint64s")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_UINT64, field.GetType(), "expected repeated")
	suite.Require().True(field.IsRepeated(), "expected repeated")
}

func (suite *AdapterTestSuite) TestExplicitSkippedMessage() {
	_, err := suite.adapter.GetFileDescriptor("ExplicitSkippedMessage")
	suite.EqualError(err, entproto.ErrSchemaSkipped.Error())
}

func (suite *AdapterTestSuite) TestSkippedFieldAndEdge() {
	message, err := suite.adapter.GetMessageDescriptor("User")
	suite.Require().NoError(err)

	postsField := message.FindFieldByName("unnecessary")
	suite.Require().Nil(postsField)

	edgeField := message.FindFieldByName("skip_edge")
	suite.Require().Nil(edgeField)
}

func (suite *AdapterTestSuite) TestInvalidField() {
	_, err := suite.adapter.GetFileDescriptor("InvalidFieldMessage")
	suite.EqualError(err, "unsupported field type \"TypeJSON\"")
}

func (suite *AdapterTestSuite) TestEnumWithConflictingValue() {
	_, err := suite.adapter.GetFileDescriptor("EnumWithConflictingValue")
	suite.EqualError(err, "entproto: Enum option \"EnumJpegAlt\" produces conflicting pbfield name \"IMAGE_JPEG\" after normalization")
}

func (suite *AdapterTestSuite) TestDuplicateNumber() {
	_, err := suite.adapter.GetFileDescriptor("DuplicateNumberMessage")
	suite.EqualError(err, "entproto: field 2 already defined on message \"DuplicateNumberMessage\"")
}

func (suite *AdapterTestSuite) TestDependsOnSkippedMessage() {
	_, err := suite.adapter.GetFileDescriptor("DependsOnSkipped")
	suite.EqualError(err, "entproto: message \"ImplicitSkippedMessage\" is not generated")
}

func (suite *AdapterTestSuite) TestMessageWithPackageName() {
	fd, err := suite.adapter.GetFileDescriptor("MessageWithPackageName")
	suite.NoError(err)
	suite.Equal(filepath.Join("io", "entgo", "apps", "todo", "todo.proto"), fd.GetName())
	suite.Equal("entgo.io/contrib/entproto/internal/entprototest/ent/proto/io/entgo/apps/todo",
		fd.GetFileOptions().GetGoPackage())
}

func (suite *AdapterTestSuite) TestManyToOne() {
	message, err := suite.adapter.GetMessageDescriptor("BlogPost")
	suite.NoError(err)

	authorField := message.FindFieldByName("author")

	suite.EqualValues(authorField.GetNumber(), 4)
}

func (suite *AdapterTestSuite) TestOneToMany() {
	message, err := suite.adapter.GetMessageDescriptor("User")
	suite.Require().NoError(err)

	postsField := message.FindFieldByName("blog_posts")
	suite.Require().NotNil(postsField)

	suite.EqualValues(postsField.GetNumber(), 3)
	suite.EqualValues(descriptorpb.FieldDescriptorProto_LABEL_REPEATED, postsField.GetLabel())
}

func (suite *AdapterTestSuite) TestManyToMany() {
	postMessage, err := suite.adapter.GetMessageDescriptor("BlogPost")
	suite.Require().NoError(err)
	categoryField := postMessage.FindFieldByName("categories")
	suite.Require().NotNil(categoryField)
	suite.EqualValues(descriptorpb.FieldDescriptorProto_LABEL_REPEATED, categoryField.GetLabel())

	categoryMessage, err := suite.adapter.GetMessageDescriptor("BlogPost")
	suite.Require().NoError(err)
	postsField := categoryMessage.FindFieldByName("categories")
	suite.Require().NotNil(postsField)
	suite.EqualValues(descriptorpb.FieldDescriptorProto_LABEL_REPEATED, postsField.GetLabel())
}

func (suite *AdapterTestSuite) TestEnumMessage() {
	fd, err := suite.adapter.GetFileDescriptor("MessageWithEnum")
	suite.NoError(err)

	message := fd.FindMessage("entpb.MessageWithEnum")
	suite.Len(message.GetFields(), 4)

	// an enum field with defaults
	enumField := message.FindFieldByName("enum_type")
	suite.EqualValues(2, enumField.GetNumber())
	suite.EqualValues(descriptorpb.FieldDescriptorProto_TYPE_ENUM, enumField.GetType())
	enumDesc := enumField.GetEnumType()
	suite.EqualValues("entpb.MessageWithEnum.EnumType", enumDesc.GetFullyQualifiedName())
	suite.EqualValues(0, enumDesc.FindValueByName("ENUM_TYPE_PENDING").GetNumber())
	suite.EqualValues(1, enumDesc.FindValueByName("ENUM_TYPE_ACTIVE").GetNumber())
	suite.EqualValues(2, enumDesc.FindValueByName("ENUM_TYPE_SUSPENDED").GetNumber())
	suite.EqualValues(3, enumDesc.FindValueByName("ENUM_TYPE_DELETED").GetNumber())

	// an enum field without defaults
	enumField = message.FindFieldByName("enum_without_default")
	suite.EqualValues(3, enumField.GetNumber())
	suite.EqualValues(descriptorpb.FieldDescriptorProto_TYPE_ENUM, enumField.GetType())
	enumDesc = enumField.GetEnumType()
	suite.EqualValues("entpb.MessageWithEnum.EnumWithoutDefault", enumDesc.GetFullyQualifiedName())
	suite.EqualValues(0, enumDesc.FindValueByName("ENUM_WITHOUT_DEFAULT_UNSPECIFIED").GetNumber())
	suite.EqualValues(1, enumDesc.FindValueByName("ENUM_WITHOUT_DEFAULT_FIRST").GetNumber())
	suite.EqualValues(2, enumDesc.FindValueByName("ENUM_WITHOUT_DEFAULT_SECOND").GetNumber())

	// an enum field with special characters
	enumField = message.FindFieldByName("enum_with_special_characters")
	suite.EqualValues(4, enumField.GetNumber())
	suite.EqualValues(descriptorpb.FieldDescriptorProto_TYPE_ENUM, enumField.GetType())
	enumDesc = enumField.GetEnumType()
	suite.EqualValues("entpb.MessageWithEnum.EnumWithSpecialCharacters", enumDesc.GetFullyQualifiedName())
	suite.EqualValues(0, enumDesc.FindValueByName("ENUM_WITH_SPECIAL_CHARACTERS_UNSPECIFIED").GetNumber())
	suite.EqualValues(1, enumDesc.FindValueByName("ENUM_WITH_SPECIAL_CHARACTERS_IMAGE_JPEG").GetNumber())
	suite.EqualValues(2, enumDesc.FindValueByName("ENUM_WITH_SPECIAL_CHARACTERS_IMAGE_PNG").GetNumber())
}

func (suite *AdapterTestSuite) TestMessageWithId() {
	message, err := suite.adapter.GetMessageDescriptor("MessageWithID")
	suite.NoError(err)

	suite.Len(message.GetFields(), 1)

	idField := message.FindFieldByName("id")
	suite.Require().NotNil(idField)
	suite.EqualValues(10, idField.GetNumber())
}

func (suite *AdapterTestSuite) TestMessageWithFieldOne() {
	_, err := suite.adapter.GetMessageDescriptor("MessageWithFieldOne")
	suite.EqualError(err, "entproto: field \"field_one\" has number 1 which is reserved for id")
}

func (suite *AdapterTestSuite) TestInterpackageDep() {
	message, err := suite.adapter.GetMessageDescriptor("Portal")
	suite.Require().NoError(err)
	suite.Len(message.GetFields(), 4)
}

func (suite *AdapterTestSuite) TestOptionals() {
	message, err := suite.adapter.GetMessageDescriptor("MessageWithOptionals")
	suite.Require().NoError(err)

	intField := message.FindFieldByName("int_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, intField.GetType())
	suite.Require().EqualValues("Int32Value", intField.GetMessageType().GetName())

	uintField := message.FindFieldByName("uint_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, uintField.GetType())
	suite.Require().EqualValues("UInt32Value", uintField.GetMessageType().GetName())

	floatField := message.FindFieldByName("float_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, floatField.GetType())
	suite.Require().EqualValues("FloatValue", floatField.GetMessageType().GetName())

	strField := message.FindFieldByName("str_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, strField.GetType())
	suite.Require().EqualValues("StringValue", strField.GetMessageType().GetName())

	boolField := message.FindFieldByName("bool_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, boolField.GetType())
	suite.Require().EqualValues("BoolValue", boolField.GetMessageType().GetName())

	bytesField := message.FindFieldByName("bytes_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, bytesField.GetType())
	suite.Require().EqualValues("BytesValue", bytesField.GetMessageType().GetName())

	uuidField := message.FindFieldByName("uuid_optional")
	suite.Require().EqualValues(descriptorpb.FieldDescriptorProto_TYPE_MESSAGE, bytesField.GetType())
	suite.Require().EqualValues("BytesValue", uuidField.GetMessageType().GetName())
}
