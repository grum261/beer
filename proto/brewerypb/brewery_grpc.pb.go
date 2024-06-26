// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/brewerypb/brewery.proto

package brewerypb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.32.0 or later.
const _ = grpc.SupportPackageIsVersion7

// BreweryDeliveryServiceClient is the client API for BreweryDeliveryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type BreweryDeliveryServiceClient interface {
	CreateBrewery(ctx context.Context, in *CreateBreweryRequest, opts ...grpc.CallOption) (*CreateBreweryRespone, error)
}

type breweryDeliveryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewBreweryDeliveryServiceClient(cc grpc.ClientConnInterface) BreweryDeliveryServiceClient {
	return &breweryDeliveryServiceClient{cc}
}

func (c *breweryDeliveryServiceClient) CreateBrewery(ctx context.Context, in *CreateBreweryRequest, opts ...grpc.CallOption) (*CreateBreweryRespone, error) {
	out := new(CreateBreweryRespone)
	err := c.cc.Invoke(ctx, "/proto.brewery.BreweryDeliveryService/CreateBrewery", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// BreweryDeliveryServiceServer is the server API for BreweryDeliveryService service.
// All implementations must embed UnimplementedBreweryDeliveryServiceServer
// for forward compatibility
type BreweryDeliveryServiceServer interface {
	CreateBrewery(context.Context, *CreateBreweryRequest) (*CreateBreweryRespone, error)
	mustEmbedUnimplementedBreweryDeliveryServiceServer()
}

// UnimplementedBreweryDeliveryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedBreweryDeliveryServiceServer struct {
}

func (UnimplementedBreweryDeliveryServiceServer) CreateBrewery(context.Context, *CreateBreweryRequest) (*CreateBreweryRespone, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreateBrewery not implemented")
}
func (UnimplementedBreweryDeliveryServiceServer) mustEmbedUnimplementedBreweryDeliveryServiceServer() {
}

// UnsafeBreweryDeliveryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to BreweryDeliveryServiceServer will
// result in compilation errors.
type UnsafeBreweryDeliveryServiceServer interface {
	mustEmbedUnimplementedBreweryDeliveryServiceServer()
}

func RegisterBreweryDeliveryServiceServer(s grpc.ServiceRegistrar, srv BreweryDeliveryServiceServer) {
	s.RegisterService(&BreweryDeliveryService_ServiceDesc, srv)
}

func _BreweryDeliveryService_CreateBrewery_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(CreateBreweryRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(BreweryDeliveryServiceServer).CreateBrewery(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.brewery.BreweryDeliveryService/CreateBrewery",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(BreweryDeliveryServiceServer).CreateBrewery(ctx, req.(*CreateBreweryRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// BreweryDeliveryService_ServiceDesc is the grpc.ServiceDesc for BreweryDeliveryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var BreweryDeliveryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.brewery.BreweryDeliveryService",
	HandlerType: (*BreweryDeliveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreateBrewery",
			Handler:    _BreweryDeliveryService_CreateBrewery_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/brewerypb/brewery.proto",
}
