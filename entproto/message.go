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
)

const MessageAnnotation = "ProtoMessage"

// MessageOption configures the entproto.Message annotation
type MessageOption func(msg *message)

// Message annotates an ent.Schema to specify that protobuf message generation is required for it.
func Message(opts ...MessageOption) schema.Annotation {
	m := message{
		Generate: true,
		Package:  "entpb",
	}
	for _, apply := range opts {
		apply(&m)
	}
	return m
}

// SkipGen annotates an ent.Schema to specify that protobuf message generation is not required for it.
// This is useful in cases where a schema ent.Mixin sets Generated to true and you want to specifically set it
// to false for this schema.
func SkipGen() schema.Annotation {
	return message{
		Generate: false,
	}
}

// PackageName modifies the generated message's protobuf package name
func PackageName(pkg string) MessageOption {
	return func(msg *message) {
		msg.Package = pkg
	}
}

type message struct {
	Generate bool
	Package  string
}

func (m message) Name() string {
	return MessageAnnotation
}

func (message) Merge(other schema.Annotation) schema.Annotation {
	return other
}

func extractMessageAnnotation(sch *gen.Type) (*message, error) {
	annot, ok := sch.Annotations[MessageAnnotation]
	if !ok {
		return nil, fmt.Errorf("entproto: schema %q does not have an entproto.Message annotation", sch.Name)
	}

	var out message
	err := mapstructure.Decode(annot, &out)
	if err != nil {
		return nil, fmt.Errorf("entproto: unable to decode entproto.Message annotation for schema %q: %w",
			sch.Name, err)
	}

	return &out, nil
}
