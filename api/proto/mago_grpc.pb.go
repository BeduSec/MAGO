// Copyright (c) BeduSec. All rights reserved.
package proto

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

var _ = context.Background

const _ = grpc.SupportPackageIsVersion7

type MagoClient interface {
	Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error)
	ReloadRules(ctx context.Context, in *ReloadRulesRequest, opts ...grpc.CallOption) (*ReloadRulesResponse, error)
}

type magoClient struct {
	cc grpc.ClientConnInterface
}

func NewMagoClient(cc grpc.ClientConnInterface) MagoClient {
	return &magoClient{cc}
}

func (c *magoClient) Health(ctx context.Context, in *HealthRequest, opts ...grpc.CallOption) (*HealthResponse, error) {
	out := new(HealthResponse)
	err := c.cc.Invoke(ctx, "/mago.Mago/Health", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *magoClient) ReloadRules(ctx context.Context, in *ReloadRulesRequest, opts ...grpc.CallOption) (*ReloadRulesResponse, error) {
	out := new(ReloadRulesResponse)
	err := c.cc.Invoke(ctx, "/mago.Mago/ReloadRules", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

type MagoServer interface {
	Health(context.Context, *HealthRequest) (*HealthResponse, error)
	ReloadRules(context.Context, *ReloadRulesRequest) (*ReloadRulesResponse, error)
}

func RegisterMagoServer(s grpc.ServiceRegistrar, srv MagoServer) {
	s.RegisterService(&Mago_ServiceDesc, srv)
}

func _Mago_Health_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HealthRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MagoServer).Health(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mago.Mago/Health",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MagoServer).Health(ctx, req.(*HealthRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _Mago_ReloadRules_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(ReloadRulesRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(MagoServer).ReloadRules(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mago.Mago/ReloadRules",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(MagoServer).ReloadRules(ctx, req.(*ReloadRulesRequest))
	}
	return interceptor(ctx, in, info, handler)
}

var Mago_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mago.Mago",
	HandlerType: (*MagoServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "Health",
			Handler:    _Mago_Health_Handler,
		},
		{
			MethodName: "ReloadRules",
			Handler:    _Mago_ReloadRules_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mago.proto",
}