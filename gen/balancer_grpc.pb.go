// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             v4.23.2
// source: balancer.proto

package gen

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

const (
	LoadBalancer_LoadBalanceRequest_FullMethodName = "/LoadBalancer/LoadBalanceRequest"
)

// LoadBalancerClient is the client API for LoadBalancer service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type LoadBalancerClient interface {
	LoadBalanceRequest(ctx context.Context, in *IncomingRequest, opts ...grpc.CallOption) (*OutgoingResponse, error)
}

type loadBalancerClient struct {
	cc grpc.ClientConnInterface
}

func NewLoadBalancerClient(cc grpc.ClientConnInterface) LoadBalancerClient {
	return &loadBalancerClient{cc}
}

func (c *loadBalancerClient) LoadBalanceRequest(ctx context.Context, in *IncomingRequest, opts ...grpc.CallOption) (*OutgoingResponse, error) {
	out := new(OutgoingResponse)
	err := c.cc.Invoke(ctx, LoadBalancer_LoadBalanceRequest_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// LoadBalancerServer is the server API for LoadBalancer service.
// All implementations must embed UnimplementedLoadBalancerServer
// for forward compatibility
type LoadBalancerServer interface {
	LoadBalanceRequest(context.Context, *IncomingRequest) (*OutgoingResponse, error)
	mustEmbedUnimplementedLoadBalancerServer()
}

// UnimplementedLoadBalancerServer must be embedded to have forward compatible implementations.
type UnimplementedLoadBalancerServer struct {
}

func (UnimplementedLoadBalancerServer) LoadBalanceRequest(context.Context, *IncomingRequest) (*OutgoingResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadBalanceRequest not implemented")
}
func (UnimplementedLoadBalancerServer) mustEmbedUnimplementedLoadBalancerServer() {}

// UnsafeLoadBalancerServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to LoadBalancerServer will
// result in compilation errors.
type UnsafeLoadBalancerServer interface {
	mustEmbedUnimplementedLoadBalancerServer()
}

func RegisterLoadBalancerServer(s grpc.ServiceRegistrar, srv LoadBalancerServer) {
	s.RegisterService(&LoadBalancer_ServiceDesc, srv)
}

func _LoadBalancer_LoadBalanceRequest_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(IncomingRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(LoadBalancerServer).LoadBalanceRequest(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: LoadBalancer_LoadBalanceRequest_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(LoadBalancerServer).LoadBalanceRequest(ctx, req.(*IncomingRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// LoadBalancer_ServiceDesc is the grpc.ServiceDesc for LoadBalancer service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var LoadBalancer_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "LoadBalancer",
	HandlerType: (*LoadBalancerServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "LoadBalanceRequest",
			Handler:    _LoadBalancer_LoadBalanceRequest_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "balancer.proto",
}
