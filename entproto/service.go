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
	"errors"
	"fmt"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/types/descriptorpb"
	_ "google.golang.org/protobuf/types/known/emptypb"
)

type method string

const (
	ServiceAnnotation        = "ProtoService"
	get               method = "Get"
	create                   = "Create"
	update                   = "Update"
	delete_                  = "Delete"
)

var (
	errNoServiceDef = errors.New("entproto: annotation entproto.Service missing")
)

type ServiceOptionsEnum int

const (
	// OptionGenGetAll generates a GetAll RPC for retrieving all ids of existing entities.
	OptionGenGetAll ServiceOptionsEnum = iota
)

type serviceOptions struct {
	GenerateGetAll bool
}

type service struct {
	Generate bool
	Options  *serviceOptions
}

func (service) Name() string {
	return ServiceAnnotation
}

func Service(options ...ServiceOptionsEnum) schema.Annotation {
	newService := &service{
		Generate: true,
		Options:  &serviceOptions{},
	}

	for _, option := range options {
		switch option {
		case OptionGenGetAll:
			newService.Options.GenerateGetAll = true
		}
	}

	return newService
}

func (a *Adapter) createServiceResources(genType *gen.Type, options *serviceOptions) (serviceResources, error) {
	name := genType.Name
	serviceFqn := fmt.Sprintf("%sService", name)

	out := serviceResources{
		svc: &descriptorpb.ServiceDescriptorProto{
			Name: &serviceFqn,
		},
	}

	for _, m := range []method{create, get, update, delete_} {
		resources, err := a.genMethodProtos(genType, m)
		if err != nil {
			return serviceResources{}, err
		}
		out.svc.Method = append(out.svc.Method, resources.methodDescriptor)
		out.svcMessages = append(out.svcMessages, resources.messages...)
	}

	if options != nil {
		resources, err := a.genOptionalMethodProtos(genType, options)
		if err != nil {
			return serviceResources{}, err
		}

		for _, resource := range resources {
			out.svc.Method = append(out.svc.Method, resource.methodDescriptor)
			out.svcMessages = append(out.svcMessages, resource.messages...)
		}
	}

	return out, nil
}

func (a *Adapter) genOptionalMethodProtos(genType *gen.Type, options *serviceOptions) ([]*methodResources, error) {
	var methods []*methodResources

	if options.GenerateGetAll {
		protoMessageFieldType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
		singleMessageField := &descriptorpb.FieldDescriptorProto{
			Name:     strptr(snake(genType.Name)),
			Number:   int32ptr(1),
			Type:     &protoMessageFieldType,
			TypeName: strptr(fmt.Sprintf("%sId", genType.Name)),
		}
		singleMessageField.Name = strptr(*singleMessageField.Name + "_ids")

		repeatedFieldLabel = descriptorpb.FieldDescriptorProto_LABEL_REPEATED
		singleMessageField.Label = &repeatedFieldLabel

		output := &descriptorpb.DescriptorProto{
			Name:  strptr(fmt.Sprintf("%sIds", genType.Name)),
			Field: []*descriptorpb.FieldDescriptorProto{singleMessageField},
		}

		methods = append(methods, &methodResources{
			methodDescriptor: &descriptorpb.MethodDescriptorProto{
				Name:       strptr("GetAll"),
				InputType:  strptr("google.protobuf.Empty"),
				OutputType: output.Name,
			},
			messages: []*descriptorpb.DescriptorProto{
				output,
			},
		})
	}

	return methods, nil
}

func (a *Adapter) genMethodProtos(genType *gen.Type, m method) (methodResources, error) {
	input := &descriptorpb.DescriptorProto{
		Name: strptr(fmt.Sprintf("%s%sRequest", m, genType.Name)),
	}
	idField, err := toProtoFieldDescriptor(genType.ID)
	if err != nil {
		return methodResources{}, err
	}
	protoMessageFieldType := descriptorpb.FieldDescriptorProto_TYPE_MESSAGE
	protoEnumFieldType := descriptorpb.FieldDescriptorProto_TYPE_ENUM
	singleMessageField := &descriptorpb.FieldDescriptorProto{
		Name:     strptr(snake(genType.Name)),
		Number:   int32ptr(1),
		Type:     &protoMessageFieldType,
		TypeName: &genType.Name,
	}
	var output string
	switch m {
	case get:
		input.Field = []*descriptorpb.FieldDescriptorProto{
			idField,
			{
				Name:     strptr("view"),
				Number:   int32ptr(2),
				Type:     &protoEnumFieldType,
				TypeName: strptr("View"),
			},
		}
		input.EnumType = append(input.EnumType, &descriptorpb.EnumDescriptorProto{
			Name: strptr("View"),
			Value: []*descriptorpb.EnumValueDescriptorProto{
				{Number: int32ptr(0), Name: strptr("VIEW_UNSPECIFIED")},
				{Number: int32ptr(1), Name: strptr("BASIC")},
				{Number: int32ptr(2), Name: strptr("WITH_EDGE_IDS")},
			},
		})
		output = genType.Name
	case create:
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		output = genType.Name
	case update:
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		output = genType.Name
	case delete_:
		input.Field = []*descriptorpb.FieldDescriptorProto{idField}
		output = "google.protobuf.Empty"
	default:
		return methodResources{}, fmt.Errorf("unknown method %q", m)
	}
	return methodResources{
		methodDescriptor: &descriptorpb.MethodDescriptorProto{
			Name:       strptr(string(m)),
			InputType:  input.Name,
			OutputType: &output,
		},
		messages: []*descriptorpb.DescriptorProto{
			input,
		},
	}, nil
}

type methodResources struct {
	methodDescriptor *descriptorpb.MethodDescriptorProto
	messages         []*descriptorpb.DescriptorProto
}

type serviceResources struct {
	svc         *descriptorpb.ServiceDescriptorProto
	svcMessages []*descriptorpb.DescriptorProto
}

func extractServiceAnnotation(sch *gen.Type) (*service, error) {
	annot, ok := sch.Annotations[ServiceAnnotation]
	if !ok {
		return nil, fmt.Errorf("%w: entproto: schema %q does not have an entproto.Service annotation",
			errNoServiceDef, sch.Name)
	}

	var out service
	err := mapstructure.Decode(annot, &out)
	if err != nil {
		return nil, fmt.Errorf("entproto: unable to decode entproto.Service annotation for schema %q: %w",
			sch.Name, err)
	}

	return &out, nil
}
