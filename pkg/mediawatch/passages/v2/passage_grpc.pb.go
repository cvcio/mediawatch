// Code generated by protoc-gen-go-grpc. DO NOT EDIT.
// versions:
// - protoc-gen-go-grpc v1.3.0
// - protoc             (unknown)
// source: mediawatch/passages/v2/passage.proto

package passagesv2

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
	PassageService_CreatePassage_FullMethodName = "/mediawatch.passages.v2.PassageService/CreatePassage"
	PassageService_GetPassages_FullMethodName   = "/mediawatch.passages.v2.PassageService/GetPassages"
)

// PassageServiceClient is the client API for PassageService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type PassageServiceClient interface {
	// create a new passage
	CreatePassage(ctx context.Context, in *Passage, opts ...grpc.CallOption) (*Passage, error)
	// get list of passages by query
	GetPassages(ctx context.Context, in *QueryPassage, opts ...grpc.CallOption) (*PassageList, error)
}

type passageServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewPassageServiceClient(cc grpc.ClientConnInterface) PassageServiceClient {
	return &passageServiceClient{cc}
}

func (c *passageServiceClient) CreatePassage(ctx context.Context, in *Passage, opts ...grpc.CallOption) (*Passage, error) {
	out := new(Passage)
	err := c.cc.Invoke(ctx, PassageService_CreatePassage_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *passageServiceClient) GetPassages(ctx context.Context, in *QueryPassage, opts ...grpc.CallOption) (*PassageList, error) {
	out := new(PassageList)
	err := c.cc.Invoke(ctx, PassageService_GetPassages_FullMethodName, in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

// PassageServiceServer is the server API for PassageService service.
// All implementations should embed UnimplementedPassageServiceServer
// for forward compatibility
type PassageServiceServer interface {
	// create a new passage
	CreatePassage(context.Context, *Passage) (*Passage, error)
	// get list of passages by query
	GetPassages(context.Context, *QueryPassage) (*PassageList, error)
}

// UnimplementedPassageServiceServer should be embedded to have forward compatible implementations.
type UnimplementedPassageServiceServer struct {
}

func (UnimplementedPassageServiceServer) CreatePassage(context.Context, *Passage) (*Passage, error) {
	return nil, status.Errorf(codes.Unimplemented, "method CreatePassage not implemented")
}
func (UnimplementedPassageServiceServer) GetPassages(context.Context, *QueryPassage) (*PassageList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetPassages not implemented")
}

// UnsafePassageServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to PassageServiceServer will
// result in compilation errors.
type UnsafePassageServiceServer interface {
	mustEmbedUnimplementedPassageServiceServer()
}

func RegisterPassageServiceServer(s grpc.ServiceRegistrar, srv PassageServiceServer) {
	s.RegisterService(&PassageService_ServiceDesc, srv)
}

func _PassageService_CreatePassage_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(Passage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassageServiceServer).CreatePassage(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassageService_CreatePassage_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassageServiceServer).CreatePassage(ctx, req.(*Passage))
	}
	return interceptor(ctx, in, info, handler)
}

func _PassageService_GetPassages_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryPassage)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(PassageServiceServer).GetPassages(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: PassageService_GetPassages_FullMethodName,
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(PassageServiceServer).GetPassages(ctx, req.(*QueryPassage))
	}
	return interceptor(ctx, in, info, handler)
}

// PassageService_ServiceDesc is the grpc.ServiceDesc for PassageService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var PassageService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mediawatch.passages.v2.PassageService",
	HandlerType: (*PassageServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "CreatePassage",
			Handler:    _PassageService_CreatePassage_Handler,
		},
		{
			MethodName: "GetPassages",
			Handler:    _PassageService_GetPassages_Handler,
		},
	},
	Streams:  []grpc.StreamDesc{},
	Metadata: "mediawatch/passages/v2/passage.proto",
}