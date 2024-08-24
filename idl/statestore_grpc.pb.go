// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.5.1
// - protoc             v5.27.3
// source: statestore.proto

package idl

import (
	context "context"
	grpc "google.golang.org/grpc"
	codes "google.golang.org/grpc/codes"
	status "google.golang.org/grpc/status"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the grpc package it is being compiled against.
// Requires gRPC-Go v1.64.0 or later.
const _ = grpc.SupportPackageIsVersion9

const (
	StateStoreService_SayHello_FullMethodName          = "/statestore.StateStoreService/SayHello"
	StateStoreService_LoadNodes_FullMethodName         = "/statestore.StateStoreService/LoadNodes"
	StateStoreService_PrintNodes_FullMethodName        = "/statestore.StateStoreService/PrintNodes"
	StateStoreService_GetNodes_FullMethodName          = "/statestore.StateStoreService/GetNodes"
	StateStoreService_LoadPods_FullMethodName          = "/statestore.StateStoreService/LoadPods"
	StateStoreService_GetPods_FullMethodName           = "/statestore.StateStoreService/GetPods"
	StateStoreService_UpdatePods_FullMethodName        = "/statestore.StateStoreService/UpdatePods"
	StateStoreService_UpdateNodeTainted_FullMethodName = "/statestore.StateStoreService/UpdateNodeTainted"
)

// StateStoreServiceClient is the client API for StateStoreService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type StateStoreServiceClient interface {
	SayHello(ctx context.Context, in *HelloWorldRequest, opts ...grpc.CallOption) (*HelloWorldResponse, error)
	LoadNodes(ctx context.Context, in *LoadNodeRequest, opts ...grpc.CallOption) (*LoadNodeResponse, error)
	PrintNodes(ctx context.Context, in *PrintNodeRequest, opts ...grpc.CallOption) (*PrintNodeResponse, error)
	GetNodes(ctx context.Context, in *GetNodeRequest, opts ...grpc.CallOption) (*GetNodeResponse, error)
	LoadPods(ctx context.Context, in *LoadPodRequest, opts ...grpc.CallOption) (*LoadPodResponse, error)
	GetPods(ctx context.Context, in *GetPodRequest, opts ...grpc.CallOption) (*GetPodResponse, error)
	UpdatePods(ctx context.Context, in *UpdatePodRequest, opts ...grpc.CallOption) (*UpdatePodResponse, error)
	UpdateNodeTainted(ctx context.Context, in *UpdateNodeTaintRequest, opts ...grpc.CallOption) (*UpdateNodeTaintResponse, error)
}

type stateStoreServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewStateStoreServiceClient(cc grpc.ClientConnInterface) StateStoreServiceClient {
	return &stateStoreServiceClient{cc}
}

func (c *stateStoreServiceClient) SayHello(ctx context.Context, in *HelloWorldRequest, opts ...grpc.CallOption) (*HelloWorldResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(HelloWorldResponse)
	err := c.cc.Invoke(ctx, StateStoreService_SayHello_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) LoadNodes(ctx context.Context, in *LoadNodeRequest, opts ...grpc.CallOption) (*LoadNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoadNodeResponse)
	err := c.cc.Invoke(ctx, StateStoreService_LoadNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) PrintNodes(ctx context.Context, in *PrintNodeRequest, opts ...grpc.CallOption) (*PrintNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(PrintNodeResponse)
	err := c.cc.Invoke(ctx, StateStoreService_PrintNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) GetNodes(ctx context.Context, in *GetNodeRequest, opts ...grpc.CallOption) (*GetNodeResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetNodeResponse)
	err := c.cc.Invoke(ctx, StateStoreService_GetNodes_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) LoadPods(ctx context.Context, in *LoadPodRequest, opts ...grpc.CallOption) (*LoadPodResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(LoadPodResponse)
	err := c.cc.Invoke(ctx, StateStoreService_LoadPods_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) GetPods(ctx context.Context, in *GetPodRequest, opts ...grpc.CallOption) (*GetPodResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(GetPodResponse)
	err := c.cc.Invoke(ctx, StateStoreService_GetPods_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) UpdatePods(ctx context.Context, in *UpdatePodRequest, opts ...grpc.CallOption) (*UpdatePodResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdatePodResponse)
	err := c.cc.Invoke(ctx, StateStoreService_UpdatePods_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *stateStoreServiceClient) UpdateNodeTainted(ctx context.Context, in *UpdateNodeTaintRequest, opts ...grpc.CallOption) (*UpdateNodeTaintResponse, error) {
	cOpts := append([]grpc.CallOption{grpc.StaticMethod()}, opts...)
	out := new(UpdateNodeTaintResponse)
	err := c.cc.Invoke(ctx, StateStoreService_UpdateNodeTainted_FullMethodName, in, out, cOpts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// StateStoreServiceServer is the server API for StateStoreService service.
// All implementations must embed UnimplementedStateStoreServiceServer
// for forward compatibility.
type StateStoreServiceServer interface {
	SayHello(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error)
	LoadNodes(context.Context, *LoadNodeRequest) (*LoadNodeResponse, error)
	PrintNodes(context.Context, *PrintNodeRequest) (*PrintNodeResponse, error)
	GetNodes(context.Context, *GetNodeRequest) (*GetNodeResponse, error)
	LoadPods(context.Context, *LoadPodRequest) (*LoadPodResponse, error)
	GetPods(context.Context, *GetPodRequest) (*GetPodResponse, error)
	UpdatePods(context.Context, *UpdatePodRequest) (*UpdatePodResponse, error)
	UpdateNodeTainted(context.Context, *UpdateNodeTaintRequest) (*UpdateNodeTaintResponse, error)
	mustEmbedUnimplementedStateStoreServiceServer()
}

// UnimplementedStateStoreServiceServer must be embedded to have
// forward compatible implementations.
//
// NOTE: this should be embedded by value instead of pointer to avoid a nil
// pointer dereference when methods are called.
type UnimplementedStateStoreServiceServer struct{}

func (UnimplementedStateStoreServiceServer) SayHello(context.Context, *HelloWorldRequest) (*HelloWorldResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method SayHello not implemented")
}
func (UnimplementedStateStoreServiceServer) LoadNodes(context.Context, *LoadNodeRequest) (*LoadNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadNodes not implemented")
}
func (UnimplementedStateStoreServiceServer) PrintNodes(context.Context, *PrintNodeRequest) (*PrintNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method PrintNodes not implemented")
}
func (UnimplementedStateStoreServiceServer) GetNodes(context.Context, *GetNodeRequest) (*GetNodeResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetNodes not implemented")
}
func (UnimplementedStateStoreServiceServer) LoadPods(context.Context, *LoadPodRequest) (*LoadPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method LoadPods not implemented")
}
func (UnimplementedStateStoreServiceServer) GetPods(context.Context, *GetPodRequest) (*GetPodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPods not implemented")
}
func (UnimplementedStateStoreServiceServer) UpdatePods(context.Context, *UpdatePodRequest) (*UpdatePodResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdatePods not implemented")
}
func (UnimplementedStateStoreServiceServer) UpdateNodeTainted(context.Context, *UpdateNodeTaintRequest) (*UpdateNodeTaintResponse, error) {
	return nil, status.Errorf(codes.Unimplemented, "method UpdateNodeTainted not implemented")
}
func (UnimplementedStateStoreServiceServer) mustEmbedUnimplementedStateStoreServiceServer() {}
func (UnimplementedStateStoreServiceServer) testEmbeddedByValue()                           {}

// UnsafeStateStoreServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to StateStoreServiceServer will
// result in compilation errors.
type UnsafeStateStoreServiceServer interface {
	mustEmbedUnimplementedStateStoreServiceServer()
}

func RegisterStateStoreServiceServer(s grpc.ServiceRegistrar, srv StateStoreServiceServer) {
	// If the following call pancis, it indicates UnimplementedStateStoreServiceServer was
	// embedded by pointer and is nil.  This will cause panics if an
	// unimplemented method is ever invoked, so we test this at initialization
	// time to prevent it from happening at runtime later due to I/O.
	if t, ok := srv.(interface{ testEmbeddedByValue() }); ok {
		t.testEmbeddedByValue()
	}
	s.RegisterService(&StateStoreService_ServiceDesc, srv)
}

func _StateStoreService_SayHello_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(HelloWorldRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).SayHello(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_SayHello_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).SayHello(ctx, req.(*HelloWorldRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_LoadNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).LoadNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_LoadNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).LoadNodes(ctx, req.(*LoadNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_PrintNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(PrintNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).PrintNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_PrintNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).PrintNodes(ctx, req.(*PrintNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_GetNodes_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetNodeRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).GetNodes(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_GetNodes_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).GetNodes(ctx, req.(*GetNodeRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_LoadPods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(LoadPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).LoadPods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_LoadPods_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).LoadPods(ctx, req.(*LoadPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_GetPods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(GetPodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).GetPods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_GetPods_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).GetPods(ctx, req.(*GetPodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_UpdatePods_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdatePodRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).UpdatePods(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_UpdatePods_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).UpdatePods(ctx, req.(*UpdatePodRequest))
	}
	return interceptor(ctx, in, info, handler)
}

func _StateStoreService_UpdateNodeTainted_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(UpdateNodeTaintRequest)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(StateStoreServiceServer).UpdateNodeTainted(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: StateStoreService_UpdateNodeTainted_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(StateStoreServiceServer).UpdateNodeTainted(ctx, req.(*UpdateNodeTaintRequest))
	}
	return interceptor(ctx, in, info, handler)
}

// StateStoreService_ServiceDesc is the grpc.ServiceDesc for StateStoreService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var StateStoreService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "statestore.StateStoreService",
	HandlerType: (*StateStoreServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "SayHello",
			Handler:    _StateStoreService_SayHello_Handler,
		},
		{
			MethodName: "LoadNodes",
			Handler:    _StateStoreService_LoadNodes_Handler,
		},
		{
			MethodName: "PrintNodes",
			Handler:    _StateStoreService_PrintNodes_Handler,
		},
		{
			MethodName: "GetNodes",
			Handler:    _StateStoreService_GetNodes_Handler,
		},
		{
			MethodName: "LoadPods",
			Handler:    _StateStoreService_LoadPods_Handler,
		},
		{
			MethodName: "GetPods",
			Handler:    _StateStoreService_GetPods_Handler,
		},
		{
			MethodName: "UpdatePods",
			Handler:    _StateStoreService_UpdatePods_Handler,
		},
		{
			MethodName: "UpdateNodeTainted",
			Handler:    _StateStoreService_UpdateNodeTainted_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "statestore.proto",
}
