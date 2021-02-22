// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package accountpb

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
const _ = grpc.SupportPackageIsVersion7

// ServiceClient is the client API for Service service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ServiceClient interface {
	FindAccountByAPIKey(ctx context.Context, in *FindAccountByAPIKeyParams, opts ...grpc.CallOption) (*FindAccountByAPIKeyReply, error)
	FindAccountByID(ctx context.Context, in *FindAccountByIDParams, opts ...grpc.CallOption) (*FindAccountByIDReply, error)
}

type serviceClient struct {
	cc grpc.ClientConnInterface
}

func NewServiceClient(cc grpc.ClientConnInterface) ServiceClient {
	return &serviceClient{cc}
}

func (c *serviceClient) FindAccountByAPIKey(ctx context.Context, in *FindAccountByAPIKeyParams, opts ...grpc.CallOption) (*FindAccountByAPIKeyReply, error) {
	out := new(FindAccountByAPIKeyReply)
	err := c.cc.Invoke(ctx, "/accountpb.Service/FindAccountByAPIKey", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *serviceClient) FindAccountByID(ctx context.Context, in *FindAccountByIDParams, opts ...grpc.CallOption) (*FindAccountByIDReply, error) {
	out := new(FindAccountByIDReply)
	err := c.cc.Invoke(ctx, "/accountpb.Service/FindAccountByID", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// ServiceServer is the server API for Service service.
// All implementations must embed UnimplementedServiceServer
// for forward compatibility
type ServiceServer interface {
	FindAccountByAPIKey(context.Context, *FindAccountByAPIKeyParams) (*FindAccountByAPIKeyReply, error)
	FindAccountByID(context.Context, *FindAccountByIDParams) (*FindAccountByIDReply, error)
	mustEmbedUnimplementedServiceServer()
}

// UnimplementedServiceServer must be embedded to have forward compatible implementations.
type UnimplementedServiceServer struct {
}

func (UnimplementedServiceServer) FindAccountByAPIKey(context.Context, *FindAccountByAPIKeyParams) (*FindAccountByAPIKeyReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAccountByAPIKey not implemented")
}
func (UnimplementedServiceServer) FindAccountByID(context.Context, *FindAccountByIDParams) (*FindAccountByIDReply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method FindAccountByID not implemented")
}
func (UnimplementedServiceServer) mustEmbedUnimplementedServiceServer() {}

// UnsafeServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ServiceServer will
// result in compilation errors.
type UnsafeServiceServer interface {
	mustEmbedUnimplementedServiceServer()
}

func RegisterServiceServer(s grpc.ServiceRegistrar, srv ServiceServer) {
	s.RegisterService(&_Service_serviceDesc, srv)
}

func _Service_FindAccountByAPIKey_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAccountByAPIKeyParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).FindAccountByAPIKey(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accountpb.Service/FindAccountByAPIKey",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).FindAccountByAPIKey(ctx, req.(*FindAccountByAPIKeyParams))
	}
	return interceptor(ctx, in, info, handler)
}

func _Service_FindAccountByID_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(FindAccountByIDParams)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ServiceServer).FindAccountByID(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/accountpb.Service/FindAccountByID",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ServiceServer).FindAccountByID(ctx, req.(*FindAccountByIDParams))
	}
	return interceptor(ctx, in, info, handler)
}

var _Service_serviceDesc = grpc.ServiceDesc{
	ServiceName: "accountpb.Service",
	HandlerType: (*ServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "FindAccountByAPIKey",
			Handler:    _Service_FindAccountByAPIKey_Handler,
		},
		{
			MethodName: "FindAccountByID",
			Handler:    _Service_FindAccountByID_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "account/rpc/accountpb/account.proto",
}
