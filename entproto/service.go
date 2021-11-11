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

const (
	ServiceAnnotation        = "ProtoService"
	MethodCreate      Method = 1 << iota
	MethodGet
	MethodUpdate
	MethodDelete
	MethodAll = MethodCreate | MethodGet | MethodUpdate | MethodDelete
)

var (
	errNoServiceDef = errors.New("entproto: annotation entproto.Service missing")
)

type Method uint

// Is reports whether method m matches given method n.
func (m Method) Is(n Method) bool { return m&n != 0 }

func Methods(methods Method) ServiceOption {
	return func(s *service) {
		s.Methods = methods
	}
}

type service struct {
	Generate bool
	Methods  Method
}

func (service) Name() string {
	return ServiceAnnotation
}

// ServiceOption configures the entproto.Service annotation.
type ServiceOption func(svc *service)

// Service annotates an ent.Schema to specify that protobuf service generation is required for it.
func Service(opts ...ServiceOption) schema.Annotation {
	s := service{
		Generate: true,
	}
	for _, apply := range opts {
		apply(&s)
	}
	// Default to generating all methods
	if s.Methods == 0 {
		s.Methods = MethodAll
	}
	return s
}

func (a *Adapter) createServiceResources(genType *gen.Type, methods Method) (serviceResources, error) {
	name := genType.Name
	serviceFqn := fmt.Sprintf("%sService", name)

	out := serviceResources{
		svc: &descriptorpb.ServiceDescriptorProto{
			Name: &serviceFqn,
		},
	}

	for _, m := range []Method{MethodCreate, MethodGet, MethodUpdate, MethodDelete} {
		if !methods.Is(m) {
			continue
		}

		resources, err := a.genMethodProtos(genType, m)
		if err != nil {
			return serviceResources{}, err
		}
		out.svc.Method = append(out.svc.Method, resources.methodDescriptor)
		out.svcMessages = append(out.svcMessages, resources.input)
	}

	return out, nil
}

func (a *Adapter) genMethodProtos(genType *gen.Type, m Method) (methodResources, error) {
	input := &descriptorpb.DescriptorProto{}
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
	var output, methodName string
	switch m {
	case MethodGet:
		methodName = "Get"
		input.Name = strptr(fmt.Sprintf("Get%sRequest", genType.Name))
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
	case MethodCreate:
		methodName = "Create"
		input.Name = strptr(fmt.Sprintf("Create%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		output = genType.Name
	case MethodUpdate:
		methodName = "Update"
		input.Name = strptr(fmt.Sprintf("Update%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		output = genType.Name
	case MethodDelete:
		methodName = "Delete"
		input.Name = strptr(fmt.Sprintf("Delete%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{idField}
		output = "google.protobuf.Empty"
	default:
		return methodResources{}, fmt.Errorf("unknown method %q", m)
	}
	return methodResources{
		methodDescriptor: &descriptorpb.MethodDescriptorProto{
			Name:       &methodName,
			InputType:  input.Name,
			OutputType: &output,
		},
		input: input,
	}, nil
}

type methodResources struct {
	methodDescriptor *descriptorpb.MethodDescriptorProto
	input            *descriptorpb.DescriptorProto
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
