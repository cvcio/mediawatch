// Code generated by protoc-gen-go-grpc. DO NOT EDIT.

package articlesv2

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

// ArticlesServiceClient is the client API for ArticlesService service.
//
// For semantics around ctx use and closing/ending streaming RPCs, please refer to https://pkg.go.dev/google.golang.org/grpc/?tab=doc#ClientConn.NewStream.
type ArticlesServiceClient interface {
	// GetArticle
	GetArticle(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (*Article, error)
	// GetArticles
	GetArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (*ArticleList, error)
	// StreamArticles
	StreamArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (ArticlesService_StreamArticlesClient, error)
	// StreamRelatedArticles
	StreamRelatedArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (ArticlesService_StreamRelatedArticlesClient, error)
}

type articlesServiceClient struct {
	cc grpc.ClientConnInterface
}

func NewArticlesServiceClient(cc grpc.ClientConnInterface) ArticlesServiceClient {
	return &articlesServiceClient{cc}
}

func (c *articlesServiceClient) GetArticle(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (*Article, error) {
	out := new(Article)
	err := c.cc.Invoke(ctx, "/mediawatch.articles.v2.ArticlesService/GetArticle", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articlesServiceClient) GetArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (*ArticleList, error) {
	out := new(ArticleList)
	err := c.cc.Invoke(ctx, "/mediawatch.articles.v2.ArticlesService/GetArticles", in, out, opts...)
	if err != nil {
		return nil, err
	}
	return out, nil
}

func (c *articlesServiceClient) StreamArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (ArticlesService_StreamArticlesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ArticlesService_ServiceDesc.Streams[0], "/mediawatch.articles.v2.ArticlesService/StreamArticles", opts...)
	if err != nil {
		return nil, err
	}
	x := &articlesServiceStreamArticlesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ArticlesService_StreamArticlesClient interface {
	Recv() (*ArticleList, error)
	grpc.ClientStream
}

type articlesServiceStreamArticlesClient struct {
	grpc.ClientStream
}

func (x *articlesServiceStreamArticlesClient) Recv() (*ArticleList, error) {
	m := new(ArticleList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

func (c *articlesServiceClient) StreamRelatedArticles(ctx context.Context, in *QueryArticle, opts ...grpc.CallOption) (ArticlesService_StreamRelatedArticlesClient, error) {
	stream, err := c.cc.NewStream(ctx, &ArticlesService_ServiceDesc.Streams[1], "/mediawatch.articles.v2.ArticlesService/StreamRelatedArticles", opts...)
	if err != nil {
		return nil, err
	}
	x := &articlesServiceStreamRelatedArticlesClient{stream}
	if err := x.ClientStream.SendMsg(in); err != nil {
		return nil, err
	}
	if err := x.ClientStream.CloseSend(); err != nil {
		return nil, err
	}
	return x, nil
}

type ArticlesService_StreamRelatedArticlesClient interface {
	Recv() (*ArticleList, error)
	grpc.ClientStream
}

type articlesServiceStreamRelatedArticlesClient struct {
	grpc.ClientStream
}

func (x *articlesServiceStreamRelatedArticlesClient) Recv() (*ArticleList, error) {
	m := new(ArticleList)
	if err := x.ClientStream.RecvMsg(m); err != nil {
		return nil, err
	}
	return m, nil
}

// ArticlesServiceServer is the server API for ArticlesService service.
// All implementations should embed UnimplementedArticlesServiceServer
// for forward compatibility
type ArticlesServiceServer interface {
	// GetArticle
	GetArticle(context.Context, *QueryArticle) (*Article, error)
	// GetArticles
	GetArticles(context.Context, *QueryArticle) (*ArticleList, error)
	// StreamArticles
	StreamArticles(*QueryArticle, ArticlesService_StreamArticlesServer) error
	// StreamRelatedArticles
	StreamRelatedArticles(*QueryArticle, ArticlesService_StreamRelatedArticlesServer) error
}

// UnimplementedArticlesServiceServer should be embedded to have forward compatible implementations.
type UnimplementedArticlesServiceServer struct {
}

func (UnimplementedArticlesServiceServer) GetArticle(context.Context, *QueryArticle) (*Article, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticle not implemented")
}
func (UnimplementedArticlesServiceServer) GetArticles(context.Context, *QueryArticle) (*ArticleList, error) {
	return nil, status.Errorf(codes.Unimplemented, "method GetArticles not implemented")
}
func (UnimplementedArticlesServiceServer) StreamArticles(*QueryArticle, ArticlesService_StreamArticlesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamArticles not implemented")
}
func (UnimplementedArticlesServiceServer) StreamRelatedArticles(*QueryArticle, ArticlesService_StreamRelatedArticlesServer) error {
	return status.Errorf(codes.Unimplemented, "method StreamRelatedArticles not implemented")
}

// UnsafeArticlesServiceServer may be embedded to opt out of forward compatibility for this service.
// Use of this interface is not recommended, as added methods to ArticlesServiceServer will
// result in compilation errors.
type UnsafeArticlesServiceServer interface {
	mustEmbedUnimplementedArticlesServiceServer()
}

func RegisterArticlesServiceServer(s grpc.ServiceRegistrar, srv ArticlesServiceServer) {
	s.RegisterService(&ArticlesService_ServiceDesc, srv)
}

func _ArticlesService_GetArticle_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryArticle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).GetArticle(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mediawatch.articles.v2.ArticlesService/GetArticle",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).GetArticle(ctx, req.(*QueryArticle))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticlesService_GetArticles_Handler(srv interface{}, ctx context.Context, dec func(interface{}) error, interceptor grpc.UnaryServerInterceptor) (interface{}, error) {
	in := new(QueryArticle)
	if err := dec(in); err != nil {
		return nil, err
	}
	if interceptor == nil {
		return srv.(ArticlesServiceServer).GetArticles(ctx, in)
	}
	info := &grpc.UnaryServerInfo{
		Server:     srv,
		FullMethod: "/mediawatch.articles.v2.ArticlesService/GetArticles",
	}
	handler := func(ctx context.Context, req interface{}) (interface{}, error) {
		return srv.(ArticlesServiceServer).GetArticles(ctx, req.(*QueryArticle))
	}
	return interceptor(ctx, in, info, handler)
}

func _ArticlesService_StreamArticles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryArticle)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ArticlesServiceServer).StreamArticles(m, &articlesServiceStreamArticlesServer{stream})
}

type ArticlesService_StreamArticlesServer interface {
	Send(*ArticleList) error
	grpc.ServerStream
}

type articlesServiceStreamArticlesServer struct {
	grpc.ServerStream
}

func (x *articlesServiceStreamArticlesServer) Send(m *ArticleList) error {
	return x.ServerStream.SendMsg(m)
}

func _ArticlesService_StreamRelatedArticles_Handler(srv interface{}, stream grpc.ServerStream) error {
	m := new(QueryArticle)
	if err := stream.RecvMsg(m); err != nil {
		return err
	}
	return srv.(ArticlesServiceServer).StreamRelatedArticles(m, &articlesServiceStreamRelatedArticlesServer{stream})
}

type ArticlesService_StreamRelatedArticlesServer interface {
	Send(*ArticleList) error
	grpc.ServerStream
}

type articlesServiceStreamRelatedArticlesServer struct {
	grpc.ServerStream
}

func (x *articlesServiceStreamRelatedArticlesServer) Send(m *ArticleList) error {
	return x.ServerStream.SendMsg(m)
}

// ArticlesService_ServiceDesc is the grpc.ServiceDesc for ArticlesService service.
// It's only intended for direct use with grpc.RegisterService,
// and not to be introspected or modified (even as a copy)
var ArticlesService_ServiceDesc = grpc.ServiceDesc{
	ServiceName: "mediawatch.articles.v2.ArticlesService",
	HandlerType: (*ArticlesServiceServer)(nil),
	Methods: []grpc.MethodDesc{
		{
			MethodName: "GetArticle",
			Handler:    _ArticlesService_GetArticle_Handler,
		},
		{
			MethodName: "GetArticles",
			Handler:    _ArticlesService_GetArticles_Handler,
		},
	},
	Streams: []grpc.StreamDesc{
		{
			StreamName:    "StreamArticles",
			Handler:       _ArticlesService_StreamArticles_Handler,
			ServerStreams: true,
		},
		{
			StreamName:    "StreamRelatedArticles",
			Handler:       _ArticlesService_StreamRelatedArticles_Handler,
			ServerStreams: true,
		},
	},
	Metadata: "mediawatch/articles/v2/articles.proto",
}
