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

import "google.golang.org/protobuf/types/descriptorpb"

const (
	TypeBool   = descriptorpb.FieldDescriptorProto_TYPE_BOOL
	TypeString = descriptorpb.FieldDescriptorProto_TYPE_STRING
)

// Message annotates an ent.Schema to specify that protobuf message generation is required for it.
func NamedMessages(messages ...*namedMessage) MessageOption {
	return func(msg *message) {
		msg.NamedMessages = append(msg.NamedMessages, messages...)
	}
}

func NamedMessage(name string) *namedMessage {
	return &namedMessage{
		ProtoMessageOptions: protoMessageOptions{},
		Name:                name,
	}
}

type namedMessage struct {
	ProtoMessageOptions protoMessageOptions
	Name                string
	Groups              FieldGroups
	ExtraFields         []pbfield
}

func (nm *namedMessage) WithGroups(groups FieldGroups) *namedMessage {
	nm.Groups = groups
	return nm
}

func (nm *namedMessage) WithSkipID(skip bool) *namedMessage {
	nm.ProtoMessageOptions.SkipID = skip
	return nm
}

func (nm *namedMessage) WithExtraFields(fields ...*extraField) *namedMessage {
	nm.ProtoMessageOptions.ExtraFields = append(nm.ProtoMessageOptions.ExtraFields, fields...)
	return nm
}

func ExtraField(name string, number int) *extraField {
	return &extraField{
		Name: name,
		Descriptor: pbfield{
			Number: number,
		},
	}
}

func (ef *extraField) WithType(t descriptorpb.FieldDescriptorProto_Type) *extraField {
	ef.Descriptor.Type = t
	return ef
}

func (ef *extraField) WithTypeName(name string) *extraField {
	ef.Descriptor.TypeName = name
	return ef
}
