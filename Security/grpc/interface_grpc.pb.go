// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package secops

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

// SecOpsClient is the client API for SecOps service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type SecOpsClient interface {
	Ping(ctx context.Context, in *Share, opts ...grpc.CallOption) (*Reply, error)
}

type secOpsClient struct {
	cc grpc.ClientConnInterface
}

func NewSecOpsClient(cc grpc.ClientConnInterface) SecOpsClient {
	return &secOpsClient{cc}
}

func (c *secOpsClient) Ping(ctx context.Context, in *Share, opts ...grpc.CallOption) (*Reply, error) {
	out := new(Reply)
	err := c.cc.Invoke(ctx, "/SecOps.SecOps/ping", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// SecOpsServer is the server API for SecOps service.
// All implementations must embed UnimplementedSecOpsServer
// for forward compatibility
type SecOpsServer interface {
	Ping(context.Context, *Share) (*Reply, error)
	mustEmbedUnimplementedSecOpsServer()
}

// UnimplementedSecOpsServer must be embedded to have forward compatible implementations.
type UnimplementedSecOpsServer struct {
}

func (UnimplementedSecOpsServer) Ping(context.Context, *Share) (*Reply, error) {
	return nil, status.Errorf(codes.Unimplemented, "method Ping not implemented")
}
func (UnimplementedSecOpsServer) mustEmbedUnimplementedSecOpsServer() {}

// UnsafeSecOpsServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to SecOpsServer will
// result in compilation errors.
type UnsafeSecOpsServer interface {
	mustEmbedUnimplementedSecOpsServer()
}

func RegisterSecOpsServer(s grpc.ServiceRegistrar, srv SecOpsServer) {
	s.RegisterService(&SecOps_ServiceDesc, srv)
}

func _SecOps_Ping_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Share)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(SecOpsServer).Ping(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/SecOps.SecOps/ping",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(SecOpsServer).Ping(ctx, req.(*Share))
	}
	return interceptor(ctx, in, info, handler)
}

// SecOps_ServiceDesc is the grpc.ServiceDesc for SecOps service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var SecOps_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "SecOps.SecOps",
	HandlerType: (*SecOpsServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "ping",
			Handler:    _SecOps_Ping_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "interface.proto",
}
