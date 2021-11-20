// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package entpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// AttachmentServiceClient is the client API for AttachmentService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type AttachmentServiceClient interface {
	Create(ctx context.Context, in *CreateAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error)
	Get(ctx context.Context, in *GetAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error)
	Update(ctx context.Context, in *UpdateAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error)
	Delete(ctx context.Context, in *DeleteAttachmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *ListAttachmentRequest, opts ...grpc.CallOption) (*ListAttachmentResponse, error)
}

type attachmentServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewAttachmentServiceClient(cc grpc.ClientConnInterface) AttachmentServiceClient {
	return &attachmentServiceClient{cc}
}

func (c *attachmentServiceClient) Create(ctx context.Context, in *CreateAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error) {
	out := new(Attachment)
	err := c.cc.Invoke(ctx, "/entpb.AttachmentService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *attachmentServiceClient) Get(ctx context.Context, in *GetAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error) {
	out := new(Attachment)
	err := c.cc.Invoke(ctx, "/entpb.AttachmentService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *attachmentServiceClient) Update(ctx context.Context, in *UpdateAttachmentRequest, opts ...grpc.CallOption) (*Attachment, error) {
	out := new(Attachment)
	err := c.cc.Invoke(ctx, "/entpb.AttachmentService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *attachmentServiceClient) Delete(ctx context.Context, in *DeleteAttachmentRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/entpb.AttachmentService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *attachmentServiceClient) List(ctx context.Context, in *ListAttachmentRequest, opts ...grpc.CallOption) (*ListAttachmentResponse, error) {
	out := new(ListAttachmentResponse)
	err := c.cc.Invoke(ctx, "/entpb.AttachmentService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// AttachmentServiceServer is the server API for AttachmentService service.
// All implementations must embed UnimplementedAttachmentServiceServer
// for forward compatibility
type AttachmentServiceServer interface {
	Create(context.Context, *CreateAttachmentRequest) (*Attachment, error)
	Get(context.Context, *GetAttachmentRequest) (*Attachment, error)
	Update(context.Context, *UpdateAttachmentRequest) (*Attachment, error)
	Delete(context.Context, *DeleteAttachmentRequest) (*emptypb.Empty, error)
	List(context.Context, *ListAttachmentRequest) (*ListAttachmentResponse, error)
	mustEmbedUnimplementedAttachmentServiceServer()
}

// UnimplementedAttachmentServiceServer must be embedded to have forward compatible implementations.
type UnimplementedAttachmentServiceServer struct {
}

func (UnimplementedAttachmentServiceServer) Create(context.Context, *CreateAttachmentRequest) (*Attachment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedAttachmentServiceServer) Get(context.Context, *GetAttachmentRequest) (*Attachment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedAttachmentServiceServer) Update(context.Context, *UpdateAttachmentRequest) (*Attachment, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedAttachmentServiceServer) Delete(context.Context, *DeleteAttachmentRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedAttachmentServiceServer) List(context.Context, *ListAttachmentRequest) (*ListAttachmentResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedAttachmentServiceServer) mustEmbedUnimplementedAttachmentServiceServer() {}

// UnsafeAttachmentServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to AttachmentServiceServer will
// result in compilation errors.
type UnsafeAttachmentServiceServer interface {
	mustEmbedUnimplementedAttachmentServiceServer()
}

func RegisterAttachmentServiceServer(s grpc.ServiceRegistrar, srv AttachmentServiceServer) {
	s.RegisterService(&AttachmentService_ServiceDesc, srv)
}

func _AttachmentService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttachmentServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.AttachmentService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttachmentServiceServer).Create(ctx, req.(*CreateAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AttachmentService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttachmentServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.AttachmentService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttachmentServiceServer).Get(ctx, req.(*GetAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AttachmentService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttachmentServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.AttachmentService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttachmentServiceServer).Update(ctx, req.(*UpdateAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AttachmentService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttachmentServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.AttachmentService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttachmentServiceServer).Delete(ctx, req.(*DeleteAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _AttachmentService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListAttachmentRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(AttachmentServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.AttachmentService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(AttachmentServiceServer).List(ctx, req.(*ListAttachmentRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// AttachmentService_ServiceDesc is the grpc.ServiceDesc for AttachmentService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var AttachmentService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entpb.AttachmentService",
	HandlerType: (*AttachmentServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _AttachmentService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _AttachmentService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _AttachmentService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _AttachmentService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _AttachmentService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entpb/entpb.proto",
}

// MultiWordSchemaServiceClient is the client API for MultiWordSchemaService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type MultiWordSchemaServiceClient interface {
	Create(ctx context.Context, in *CreateMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error)
	Get(ctx context.Context, in *GetMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error)
	Update(ctx context.Context, in *UpdateMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error)
	Delete(ctx context.Context, in *DeleteMultiWordSchemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *ListMultiWordSchemaRequest, opts ...grpc.CallOption) (*ListMultiWordSchemaResponse, error)
}

type multiWordSchemaServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewMultiWordSchemaServiceClient(cc grpc.ClientConnInterface) MultiWordSchemaServiceClient {
	return &multiWordSchemaServiceClient{cc}
}

func (c *multiWordSchemaServiceClient) Create(ctx context.Context, in *CreateMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error) {
	out := new(MultiWordSchema)
	err := c.cc.Invoke(ctx, "/entpb.MultiWordSchemaService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiWordSchemaServiceClient) Get(ctx context.Context, in *GetMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error) {
	out := new(MultiWordSchema)
	err := c.cc.Invoke(ctx, "/entpb.MultiWordSchemaService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiWordSchemaServiceClient) Update(ctx context.Context, in *UpdateMultiWordSchemaRequest, opts ...grpc.CallOption) (*MultiWordSchema, error) {
	out := new(MultiWordSchema)
	err := c.cc.Invoke(ctx, "/entpb.MultiWordSchemaService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiWordSchemaServiceClient) Delete(ctx context.Context, in *DeleteMultiWordSchemaRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/entpb.MultiWordSchemaService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *multiWordSchemaServiceClient) List(ctx context.Context, in *ListMultiWordSchemaRequest, opts ...grpc.CallOption) (*ListMultiWordSchemaResponse, error) {
	out := new(ListMultiWordSchemaResponse)
	err := c.cc.Invoke(ctx, "/entpb.MultiWordSchemaService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// MultiWordSchemaServiceServer is the server API for MultiWordSchemaService service.
// All implementations must embed UnimplementedMultiWordSchemaServiceServer
// for forward compatibility
type MultiWordSchemaServiceServer interface {
	Create(context.Context, *CreateMultiWordSchemaRequest) (*MultiWordSchema, error)
	Get(context.Context, *GetMultiWordSchemaRequest) (*MultiWordSchema, error)
	Update(context.Context, *UpdateMultiWordSchemaRequest) (*MultiWordSchema, error)
	Delete(context.Context, *DeleteMultiWordSchemaRequest) (*emptypb.Empty, error)
	List(context.Context, *ListMultiWordSchemaRequest) (*ListMultiWordSchemaResponse, error)
	mustEmbedUnimplementedMultiWordSchemaServiceServer()
}

// UnimplementedMultiWordSchemaServiceServer must be embedded to have forward compatible implementations.
type UnimplementedMultiWordSchemaServiceServer struct {
}

func (UnimplementedMultiWordSchemaServiceServer) Create(context.Context, *CreateMultiWordSchemaRequest) (*MultiWordSchema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedMultiWordSchemaServiceServer) Get(context.Context, *GetMultiWordSchemaRequest) (*MultiWordSchema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedMultiWordSchemaServiceServer) Update(context.Context, *UpdateMultiWordSchemaRequest) (*MultiWordSchema, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedMultiWordSchemaServiceServer) Delete(context.Context, *DeleteMultiWordSchemaRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedMultiWordSchemaServiceServer) List(context.Context, *ListMultiWordSchemaRequest) (*ListMultiWordSchemaResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedMultiWordSchemaServiceServer) mustEmbedUnimplementedMultiWordSchemaServiceServer() {
}

// UnsafeMultiWordSchemaServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to MultiWordSchemaServiceServer will
// result in compilation errors.
type UnsafeMultiWordSchemaServiceServer interface {
	mustEmbedUnimplementedMultiWordSchemaServiceServer()
}

func RegisterMultiWordSchemaServiceServer(s grpc.ServiceRegistrar, srv MultiWordSchemaServiceServer) {
	s.RegisterService(&MultiWordSchemaService_ServiceDesc, srv)
}

func _MultiWordSchemaService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateMultiWordSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiWordSchemaServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.MultiWordSchemaService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiWordSchemaServiceServer).Create(ctx, req.(*CreateMultiWordSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiWordSchemaService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetMultiWordSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiWordSchemaServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.MultiWordSchemaService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiWordSchemaServiceServer).Get(ctx, req.(*GetMultiWordSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiWordSchemaService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateMultiWordSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiWordSchemaServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.MultiWordSchemaService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiWordSchemaServiceServer).Update(ctx, req.(*UpdateMultiWordSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiWordSchemaService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteMultiWordSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiWordSchemaServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.MultiWordSchemaService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiWordSchemaServiceServer).Delete(ctx, req.(*DeleteMultiWordSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _MultiWordSchemaService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListMultiWordSchemaRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MultiWordSchemaServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.MultiWordSchemaService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MultiWordSchemaServiceServer).List(ctx, req.(*ListMultiWordSchemaRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// MultiWordSchemaService_ServiceDesc is the grpc.ServiceDesc for MultiWordSchemaService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var MultiWordSchemaService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entpb.MultiWordSchemaService",
	HandlerType: (*MultiWordSchemaServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _MultiWordSchemaService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _MultiWordSchemaService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _MultiWordSchemaService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _MultiWordSchemaService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _MultiWordSchemaService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entpb/entpb.proto",
}

// NilExampleServiceClient is the client API for NilExampleService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NilExampleServiceClient interface {
	Create(ctx context.Context, in *CreateNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error)
	Get(ctx context.Context, in *GetNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error)
	Update(ctx context.Context, in *UpdateNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error)
	Delete(ctx context.Context, in *DeleteNilExampleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *ListNilExampleRequest, opts ...grpc.CallOption) (*ListNilExampleResponse, error)
}

type nilExampleServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNilExampleServiceClient(cc grpc.ClientConnInterface) NilExampleServiceClient {
	return &nilExampleServiceClient{cc}
}

func (c *nilExampleServiceClient) Create(ctx context.Context, in *CreateNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error) {
	out := new(NilExample)
	err := c.cc.Invoke(ctx, "/entpb.NilExampleService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nilExampleServiceClient) Get(ctx context.Context, in *GetNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error) {
	out := new(NilExample)
	err := c.cc.Invoke(ctx, "/entpb.NilExampleService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nilExampleServiceClient) Update(ctx context.Context, in *UpdateNilExampleRequest, opts ...grpc.CallOption) (*NilExample, error) {
	out := new(NilExample)
	err := c.cc.Invoke(ctx, "/entpb.NilExampleService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nilExampleServiceClient) Delete(ctx context.Context, in *DeleteNilExampleRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/entpb.NilExampleService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *nilExampleServiceClient) List(ctx context.Context, in *ListNilExampleRequest, opts ...grpc.CallOption) (*ListNilExampleResponse, error) {
	out := new(ListNilExampleResponse)
	err := c.cc.Invoke(ctx, "/entpb.NilExampleService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NilExampleServiceServer is the server API for NilExampleService service.
// All implementations must embed UnimplementedNilExampleServiceServer
// for forward compatibility
type NilExampleServiceServer interface {
	Create(context.Context, *CreateNilExampleRequest) (*NilExample, error)
	Get(context.Context, *GetNilExampleRequest) (*NilExample, error)
	Update(context.Context, *UpdateNilExampleRequest) (*NilExample, error)
	Delete(context.Context, *DeleteNilExampleRequest) (*emptypb.Empty, error)
	List(context.Context, *ListNilExampleRequest) (*ListNilExampleResponse, error)
	mustEmbedUnimplementedNilExampleServiceServer()
}

// UnimplementedNilExampleServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNilExampleServiceServer struct {
}

func (UnimplementedNilExampleServiceServer) Create(context.Context, *CreateNilExampleRequest) (*NilExample, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedNilExampleServiceServer) Get(context.Context, *GetNilExampleRequest) (*NilExample, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedNilExampleServiceServer) Update(context.Context, *UpdateNilExampleRequest) (*NilExample, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedNilExampleServiceServer) Delete(context.Context, *DeleteNilExampleRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedNilExampleServiceServer) List(context.Context, *ListNilExampleRequest) (*ListNilExampleResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedNilExampleServiceServer) mustEmbedUnimplementedNilExampleServiceServer() {}

// UnsafeNilExampleServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NilExampleServiceServer will
// result in compilation errors.
type UnsafeNilExampleServiceServer interface {
	mustEmbedUnimplementedNilExampleServiceServer()
}

func RegisterNilExampleServiceServer(s grpc.ServiceRegistrar, srv NilExampleServiceServer) {
	s.RegisterService(&NilExampleService_ServiceDesc, srv)
}

func _NilExampleService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateNilExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NilExampleServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NilExampleService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NilExampleServiceServer).Create(ctx, req.(*CreateNilExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NilExampleService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNilExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NilExampleServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NilExampleService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NilExampleServiceServer).Get(ctx, req.(*GetNilExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NilExampleService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNilExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NilExampleServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NilExampleService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NilExampleServiceServer).Update(ctx, req.(*UpdateNilExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NilExampleService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteNilExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NilExampleServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NilExampleService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NilExampleServiceServer).Delete(ctx, req.(*DeleteNilExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _NilExampleService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListNilExampleRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NilExampleServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NilExampleService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NilExampleServiceServer).List(ctx, req.(*ListNilExampleRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// NilExampleService_ServiceDesc is the grpc.ServiceDesc for NilExampleService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NilExampleService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entpb.NilExampleService",
	HandlerType: (*NilExampleServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _NilExampleService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _NilExampleService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _NilExampleService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _NilExampleService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _NilExampleService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entpb/entpb.proto",
}

// UserServiceClient is the client API for UserService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserServiceClient interface {
	Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error)
	Get(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error)
	Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error)
	Delete(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	List(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserResponse, error)
}

type userServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserServiceClient(cc grpc.ClientConnInterface) UserServiceClient {
	return &userServiceClient{cc}
}

func (c *userServiceClient) Create(ctx context.Context, in *CreateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/entpb.UserService/Create", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Get(ctx context.Context, in *GetUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/entpb.UserService/Get", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Update(ctx context.Context, in *UpdateUserRequest, opts ...grpc.CallOption) (*User, error) {
	out := new(User)
	err := c.cc.Invoke(ctx, "/entpb.UserService/Update", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) Delete(ctx context.Context, in *DeleteUserRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/entpb.UserService/Delete", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userServiceClient) List(ctx context.Context, in *ListUserRequest, opts ...grpc.CallOption) (*ListUserResponse, error) {
	out := new(ListUserResponse)
	err := c.cc.Invoke(ctx, "/entpb.UserService/List", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserServiceServer is the server API for UserService service.
// All implementations must embed UnimplementedUserServiceServer
// for forward compatibility
type UserServiceServer interface {
	Create(context.Context, *CreateUserRequest) (*User, error)
	Get(context.Context, *GetUserRequest) (*User, error)
	Update(context.Context, *UpdateUserRequest) (*User, error)
	Delete(context.Context, *DeleteUserRequest) (*emptypb.Empty, error)
	List(context.Context, *ListUserRequest) (*ListUserResponse, error)
	mustEmbedUnimplementedUserServiceServer()
}

// UnimplementedUserServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserServiceServer struct {
}

func (UnimplementedUserServiceServer) Create(context.Context, *CreateUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Create not implemented")
}
func (UnimplementedUserServiceServer) Get(context.Context, *GetUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Get not implemented")
}
func (UnimplementedUserServiceServer) Update(context.Context, *UpdateUserRequest) (*User, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Update not implemented")
}
func (UnimplementedUserServiceServer) Delete(context.Context, *DeleteUserRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Delete not implemented")
}
func (UnimplementedUserServiceServer) List(context.Context, *ListUserRequest) (*ListUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method List not implemented")
}
func (UnimplementedUserServiceServer) mustEmbedUnimplementedUserServiceServer() {}

// UnsafeUserServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserServiceServer will
// result in compilation errors.
type UnsafeUserServiceServer interface {
	mustEmbedUnimplementedUserServiceServer()
}

func RegisterUserServiceServer(s grpc.ServiceRegistrar, srv UserServiceServer) {
	s.RegisterService(&UserService_ServiceDesc, srv)
}

func _UserService_Create_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Create(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.UserService/Create",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Create(ctx, req.(*CreateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Get_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Get(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.UserService/Get",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Get(ctx, req.(*GetUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Update_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Update(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.UserService/Update",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Update(ctx, req.(*UpdateUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_Delete_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(DeleteUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).Delete(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.UserService/Delete",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).Delete(ctx, req.(*DeleteUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserService_List_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ListUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserServiceServer).List(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.UserService/List",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserServiceServer).List(ctx, req.(*ListUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserService_ServiceDesc is the grpc.ServiceDesc for UserService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entpb.UserService",
	HandlerType: (*UserServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Create",
			Handler:    _UserService_Create_Handler,
		},
		{
			MethodName: "Get",
			Handler:    _UserService_Get_Handler,
		},
		{
			MethodName: "Update",
			Handler:    _UserService_Update_Handler,
		},
		{
			MethodName: "Delete",
			Handler:    _UserService_Delete_Handler,
		},
		{
			MethodName: "List",
			Handler:    _UserService_List_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entpb/entpb.proto",
}
