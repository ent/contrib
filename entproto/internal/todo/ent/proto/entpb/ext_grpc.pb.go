// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.4.0
// source: entpb/ext.proto

package entpb

import (
	context "context"
	empty "github.com/golang/protobuf/ptypes/empty"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// NoopServiceClient is the client API for NoopService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type NoopServiceClient interface {
	Crickets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error)
}

type noopServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewNoopServiceClient(cc grpc.ClientConnInterface) NoopServiceClient {
	return &noopServiceClient{cc}
}

func (c *noopServiceClient) Crickets(ctx context.Context, in *empty.Empty, opts ...grpc.CallOption) (*empty.Empty, error) {
	out := new(empty.Empty)
	err := c.cc.Invoke(ctx, "/entpb.NoopService/Crickets", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// NoopServiceServer is the server API for NoopService service.
// All implementations must embed UnimplementedNoopServiceServer
// for forward compatibility
type NoopServiceServer interface {
	Crickets(context.Context, *empty.Empty) (*empty.Empty, error)
	mustEmbedUnimplementedNoopServiceServer()
}

// UnimplementedNoopServiceServer must be embedded to have forward compatible implementations.
type UnimplementedNoopServiceServer struct {
}

func (UnimplementedNoopServiceServer) Crickets(context.Context, *empty.Empty) (*empty.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Crickets not implemented")
}
func (UnimplementedNoopServiceServer) mustEmbedUnimplementedNoopServiceServer() {}

// UnsafeNoopServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to NoopServiceServer will
// result in compilation errors.
type UnsafeNoopServiceServer interface {
	mustEmbedUnimplementedNoopServiceServer()
}

func RegisterNoopServiceServer(s grpc.ServiceRegistrar, srv NoopServiceServer) {
	s.RegisterService(&NoopService_ServiceDesc, srv)
}

func _NoopService_Crickets_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(empty.Empty)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(NoopServiceServer).Crickets(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/entpb.NoopService/Crickets",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(NoopServiceServer).Crickets(ctx, req.(*empty.Empty))
	}
	return interceptor(ctx, in, info, handler)
}

// NoopService_ServiceDesc is the grpc.ServiceDesc for NoopService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var NoopService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "entpb.NoopService",
	HandlerType: (*NoopServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Crickets",
			Handler:    _NoopService_Crickets_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "entpb/ext.proto",
}
