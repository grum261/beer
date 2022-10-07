// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.2.0
// - protoc             v3.6.1
// source: proto/userpb/user.proto

package userpb

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

// UserDeliveryServiceClient is the client API for UserDeliveryService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type UserDeliveryServiceClient interface {
	AuthUserHandler(ctx context.Context, in *AuthUserRequest, opts ...grpc.CallOption) (*AuthUserResponse, error)
	SendFriendRequestHandler(ctx context.Context, in *FriendRequest, opts ...grpc.CallOption) (*FriendResponse, error)
	UpdateFriendRequestHandler(ctx context.Context, in *UpdateFriendRequest, opts ...grpc.CallOption) (*emptypb.Empty, error)
	GetUserFriends(ctx context.Context, in *UserFriendsRequest, opts ...grpc.CallOption) (*UserFriendsResponse, error)
}

type userDeliveryServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewUserDeliveryServiceClient(cc grpc.ClientConnInterface) UserDeliveryServiceClient {
	return &userDeliveryServiceClient{cc}
}

func (c *userDeliveryServiceClient) AuthUserHandler(ctx context.Context, in *AuthUserRequest, opts ...grpc.CallOption) (*AuthUserResponse, error) {
	out := new(AuthUserResponse)
	err := c.cc.Invoke(ctx, "/proto.user.UserDeliveryService/AuthUserHandler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDeliveryServiceClient) SendFriendRequestHandler(ctx context.Context, in *FriendRequest, opts ...grpc.CallOption) (*FriendResponse, error) {
	out := new(FriendResponse)
	err := c.cc.Invoke(ctx, "/proto.user.UserDeliveryService/SendFriendRequestHandler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDeliveryServiceClient) UpdateFriendRequestHandler(ctx context.Context, in *UpdateFriendRequest, opts ...grpc.CallOption) (*emptypb.Empty, error) {
	out := new(emptypb.Empty)
	err := c.cc.Invoke(ctx, "/proto.user.UserDeliveryService/UpdateFriendRequestHandler", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *userDeliveryServiceClient) GetUserFriends(ctx context.Context, in *UserFriendsRequest, opts ...grpc.CallOption) (*UserFriendsResponse, error) {
	out := new(UserFriendsResponse)
	err := c.cc.Invoke(ctx, "/proto.user.UserDeliveryService/GetUserFriends", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// UserDeliveryServiceServer is the server API for UserDeliveryService service.
// All implementations must embed UnimplementedUserDeliveryServiceServer
// for forward compatibility
type UserDeliveryServiceServer interface {
	AuthUserHandler(context.Context, *AuthUserRequest) (*AuthUserResponse, error)
	SendFriendRequestHandler(context.Context, *FriendRequest) (*FriendResponse, error)
	UpdateFriendRequestHandler(context.Context, *UpdateFriendRequest) (*emptypb.Empty, error)
	GetUserFriends(context.Context, *UserFriendsRequest) (*UserFriendsResponse, error)
	mustEmbedUnimplementedUserDeliveryServiceServer()
}

// UnimplementedUserDeliveryServiceServer must be embedded to have forward compatible implementations.
type UnimplementedUserDeliveryServiceServer struct {
}

func (UnimplementedUserDeliveryServiceServer) AuthUserHandler(context.Context, *AuthUserRequest) (*AuthUserResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method AuthUserHandler not implemented")
}
func (UnimplementedUserDeliveryServiceServer) SendFriendRequestHandler(context.Context, *FriendRequest) (*FriendResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SendFriendRequestHandler not implemented")
}
func (UnimplementedUserDeliveryServiceServer) UpdateFriendRequestHandler(context.Context, *UpdateFriendRequest) (*emptypb.Empty, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateFriendRequestHandler not implemented")
}
func (UnimplementedUserDeliveryServiceServer) GetUserFriends(context.Context, *UserFriendsRequest) (*UserFriendsResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetUserFriends not implemented")
}
func (UnimplementedUserDeliveryServiceServer) mustEmbedUnimplementedUserDeliveryServiceServer() {}

// UnsafeUserDeliveryServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to UserDeliveryServiceServer will
// result in compilation errors.
type UnsafeUserDeliveryServiceServer interface {
	mustEmbedUnimplementedUserDeliveryServiceServer()
}

func RegisterUserDeliveryServiceServer(s grpc.ServiceRegistrar, srv UserDeliveryServiceServer) {
	s.RegisterService(&UserDeliveryService_ServiceDesc, srv)
}

func _UserDeliveryService_AuthUserHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(AuthUserRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDeliveryServiceServer).AuthUserHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.user.UserDeliveryService/AuthUserHandler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDeliveryServiceServer).AuthUserHandler(ctx, req.(*AuthUserRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDeliveryService_SendFriendRequestHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDeliveryServiceServer).SendFriendRequestHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.user.UserDeliveryService/SendFriendRequestHandler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDeliveryServiceServer).SendFriendRequestHandler(ctx, req.(*FriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDeliveryService_UpdateFriendRequestHandler_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateFriendRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDeliveryServiceServer).UpdateFriendRequestHandler(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.user.UserDeliveryService/UpdateFriendRequestHandler",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDeliveryServiceServer).UpdateFriendRequestHandler(ctx, req.(*UpdateFriendRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _UserDeliveryService_GetUserFriends_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UserFriendsRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(UserDeliveryServiceServer).GetUserFriends(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/proto.user.UserDeliveryService/GetUserFriends",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(UserDeliveryServiceServer).GetUserFriends(ctx, req.(*UserFriendsRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// UserDeliveryService_ServiceDesc is the grpc.ServiceDesc for UserDeliveryService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var UserDeliveryService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "proto.user.UserDeliveryService",
	HandlerType: (*UserDeliveryServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "AuthUserHandler",
			Handler:    _UserDeliveryService_AuthUserHandler_Handler,
		},
		{
			MethodName: "SendFriendRequestHandler",
			Handler:    _UserDeliveryService_SendFriendRequestHandler_Handler,
		},
		{
			MethodName: "UpdateFriendRequestHandler",
			Handler:    _UserDeliveryService_UpdateFriendRequestHandler_Handler,
		},
		{
			MethodName: "GetUserFriends",
			Handler:    _UserDeliveryService_GetUserFriends_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "proto/userpb/user.proto",
}
