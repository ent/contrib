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
	"entgo.io/ent/schema/field"
	"errors"
	"fmt"

	"entgo.io/ent/entc/gen"
	"entgo.io/ent/schema"
	"github.com/mitchellh/mapstructure"
	"google.golang.org/protobuf/types/descriptorpb"
	_ "google.golang.org/protobuf/types/known/emptypb"
)

const (
	ServiceAnnotation = "ProtoService"
	// MaxPageSize is the maximum page size that can be returned by a List call. Requesting page sizes larger than
	// this value will return, at most, MaxPageSize entries.
	MaxPageSize = 1000
	// MethodCreate generates a Create gRPC service method for the entproto.Service.
	MethodCreate Method = 1 << iota
	// MethodGet generates a Get gRPC service method for the entproto.Service.
	MethodGet
	// MethodUpdate generates an Update gRPC service method for the entproto.Service.
	MethodUpdate
	// MethodDelete generates a Delete gRPC service method for the entproto.Service.
	MethodDelete
	// MethodList generates a List gRPC service method for the entproto.Service.
	MethodList

	MethodListByRelatedTo

	MethodListByRelatedToPaginated
	// MethodAll generates all service methods for the entproto.Service. This is the same behavior as not including entproto.Methods.
	MethodAll = MethodCreate | MethodGet | MethodUpdate | MethodDelete | MethodList | MethodListByRelatedTo | MethodListByRelatedToPaginated
)

var (
	errNoServiceDef = errors.New("entproto: annotation entproto.Service missing")

	int32FieldType        = descriptorpb.FieldDescriptorProto_TYPE_INT32
	uint64FieldType       = descriptorpb.FieldDescriptorProto_TYPE_UINT64
	stringFieldType       = descriptorpb.FieldDescriptorProto_TYPE_STRING
	protoEnumFieldType    = descriptorpb.FieldDescriptorProto_TYPE_ENUM
	protoMessageFieldType = descriptorpb.FieldDescriptorProto_TYPE_MESSAGE

	queryByEdgeRequest = descriptorpb.DescriptorProto{
		Name: strptr("ByEntityRequest"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   strptr("id"),
				Number: int32ptr(1),
				Type:   &uint64FieldType,
			},
			{
				Name:     strptr("type"),
				Number:   int32ptr(2),
				Type:     &protoEnumFieldType,
				TypeName: strptr(NodesTypesEnumName),
			},
		},
	}

	queryByEdgeRequestPaginated = descriptorpb.DescriptorProto{
		Name: strptr("ByEntityRequestPaginated"),
		Field: []*descriptorpb.FieldDescriptorProto{
			{
				Name:   strptr("id"),
				Number: int32ptr(1),
				Type:   &uint64FieldType,
			},
			{
				Name:     strptr("type"),
				Number:   int32ptr(2),
				Type:     &protoEnumFieldType,
				TypeName: strptr(NodesTypesEnumName),
			},
			{
				Name:   strptr("page_size"),
				Number: int32ptr(3),
				Type:   &int32FieldType,
			},
			{
				Name:   strptr("page_token"),
				Number: int32ptr(4),
				Type:   &stringFieldType,
			},
			{
				Name:     strptr("view"),
				Number:   int32ptr(5),
				Type:     &protoEnumFieldType,
				TypeName: strptr("View"),
			},
		},
		EnumType: []*descriptorpb.EnumDescriptorProto{{
			Name: strptr("View"),
			Value: []*descriptorpb.EnumValueDescriptorProto{
				{Number: int32ptr(0), Name: strptr("VIEW_UNSPECIFIED")},
				{Number: int32ptr(1), Name: strptr("BASIC")},
				{Number: int32ptr(2), Name: strptr("WITH_EDGE_IDS")},
			},
		},
		}}
)

type Method uint

// Is reports whether method m matches given method n.
func (m Method) Is(n Method) bool { return m&n != 0 }

// Methods specifies the gRPC service methods to generate for the entproto.Service.
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
	// Default to generating all methods.
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

	for _, m := range []Method{MethodCreate, MethodGet, MethodUpdate, MethodDelete, MethodList, MethodListByRelatedTo, MethodListByRelatedToPaginated} {
		if !methods.Is(m) {
			continue
		}

		resources, err := a.genMethodProtos(genType, m)
		if err != nil {
			return serviceResources{}, err
		}
		out.svc.Method = append(out.svc.Method, resources.methodDescriptor)
		out.svcMessages = append(out.svcMessages, resources.messages...)
	}

	return out, nil
}

func (a *Adapter) genMethodProtos(genType *gen.Type, m Method) (methodResources, error) {
	input := &descriptorpb.DescriptorProto{}
	idField, err := toProtoFieldDescriptor(genType.ID)
	if err != nil {
		return methodResources{}, err
	}

	singleMessageField := &descriptorpb.FieldDescriptorProto{
		Name:     strptr(snake(genType.Name)),
		Number:   int32ptr(1),
		Type:     &protoMessageFieldType,
		TypeName: &genType.Name,
	}
	var (
		outputName, methodName string
		messages               []*descriptorpb.DescriptorProto
	)
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
		outputName = genType.Name
		messages = append(messages, input)
	case MethodListByRelatedTo:
		methodName = "ListByRelatedTo"

		input = &queryByEdgeRequest
		outputName = fmt.Sprintf("List%sResponse", genType.Name)
		output := &descriptorpb.DescriptorProto{
			Name: &outputName,
			Field: []*descriptorpb.FieldDescriptorProto{
				{
					Name:     strptr(snake(genType.Name) + "_list"),
					Number:   int32ptr(1),
					Label:    &repeatedFieldLabel,
					Type:     &protoMessageFieldType,
					TypeName: strptr(genType.Name),
				},
			},
		}
		messages = append(messages, output)

	case MethodListByRelatedToPaginated:
		methodName = "ListByRelatedToPaginated"

		input = &queryByEdgeRequestPaginated
		outputName = fmt.Sprintf("List%sPaginatedResponse", genType.Name)
		output := &descriptorpb.DescriptorProto{
			Name: &outputName,
			Field: []*descriptorpb.FieldDescriptorProto{
				{
					Name:     strptr(snake(genType.Name) + "_list"),
					Number:   int32ptr(1),
					Label:    &repeatedFieldLabel,
					Type:     &protoMessageFieldType,
					TypeName: strptr(genType.Name),
				},
				{
					Name:   strptr("next_page_token"),
					Number: int32ptr(2),
					Type:   &stringFieldType,
				},
			},
		}
		messages = append(messages, output)
	case MethodCreate:
		methodName = "Create"
		input.Name = strptr(fmt.Sprintf("Create%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		outputName = genType.Name
		messages = append(messages, input)
	case MethodUpdate:
		methodName = "Update"
		input.Name = strptr(fmt.Sprintf("Update%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{singleMessageField}
		outputName = genType.Name
		messages = append(messages, input)
	case MethodDelete:
		methodName = "Delete"
		input.Name = strptr(fmt.Sprintf("Delete%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{idField}
		outputName = "google.protobuf.Empty"
		messages = append(messages, input)
	case MethodList:
		if !(genType.ID.IsInt() || genType.ID.IsUUID() || genType.ID.IsString() || (genType.ID.Type != nil && genType.ID.Type.Type == field.TypeUint64)) {
			return methodResources{}, fmt.Errorf("entproto: list method does not support schema %q id type %q",
				genType.Name, genType.ID.Type.String())
		}

		methodName = "List"
		// TODO: changed the original name (List%sRequest) to ListAll%sRequest because otherwise it clashed with
		//       ListByRelatedTo response type and raises a duplication error. Should be handled if we ever use this API.
		input.Name = strptr(fmt.Sprintf("ListAll%sRequest", genType.Name))
		input.Field = []*descriptorpb.FieldDescriptorProto{
			{
				Name:   strptr("page_size"),
				Number: int32ptr(1),
				Type:   &int32FieldType,
			},
			{
				Name:   strptr("page_token"),
				Number: int32ptr(2),
				Type:   &stringFieldType,
			},
			{
				Name:     strptr("view"),
				Number:   int32ptr(3),
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
		outputName = fmt.Sprintf("List%sResponse", genType.Name)
		output := &descriptorpb.DescriptorProto{
			Name: &outputName,
			Field: []*descriptorpb.FieldDescriptorProto{
				{
					Name:     strptr(snake(genType.Name) + "_list"),
					Number:   int32ptr(1),
					Label:    &repeatedFieldLabel,
					Type:     &protoMessageFieldType,
					TypeName: strptr(genType.Name),
				},
				{
					Name:   strptr("next_page_token"),
					Number: int32ptr(2),
					Type:   &stringFieldType,
				},
			},
		}
		messages = append(messages, input, output)
	default:
		return methodResources{}, fmt.Errorf("unknown method %q", m)
	}
	return methodResources{
		methodDescriptor: &descriptorpb.MethodDescriptorProto{
			Name:       &methodName,
			InputType:  input.Name,
			OutputType: &outputName,
		},
		messages: messages,
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
