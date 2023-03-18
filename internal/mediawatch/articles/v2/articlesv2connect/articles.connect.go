// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/articles/v2/articles.proto

package articlesv2connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v2 "github.com/cvcio/mediawatch/internal/mediawatch/articles/v2"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect_go.IsAtLeastVersion0_1_0

const (
	// ArticlesServiceName is the fully-qualified name of the ArticlesService service.
	ArticlesServiceName = "mediawatch.articles.v2.ArticlesService"
)

// ArticlesServiceClient is a client for the mediawatch.articles.v2.ArticlesService service.
type ArticlesServiceClient interface {
	// GetArticle
	GetArticle(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.Article], error)
	// GetArticles
	GetArticles(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.ArticleList], error)
	// StreamArticles
	StreamArticles(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.ServerStreamForClient[v2.ArticleList], error)
	// StreamRelatedArticles
	StreamRelatedArticles(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.ServerStreamForClient[v2.ArticleList], error)
}

// NewArticlesServiceClient constructs a client for the mediawatch.articles.v2.ArticlesService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewArticlesServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ArticlesServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &articlesServiceClient{
		getArticle: connect_go.NewClient[v2.QueryArticle, v2.Article](
			httpClient,
			baseURL+"/mediawatch.articles.v2.ArticlesService/GetArticle",
			opts...,
		),
		getArticles: connect_go.NewClient[v2.QueryArticle, v2.ArticleList](
			httpClient,
			baseURL+"/mediawatch.articles.v2.ArticlesService/GetArticles",
			opts...,
		),
		streamArticles: connect_go.NewClient[v2.QueryArticle, v2.ArticleList](
			httpClient,
			baseURL+"/mediawatch.articles.v2.ArticlesService/StreamArticles",
			opts...,
		),
		streamRelatedArticles: connect_go.NewClient[v2.QueryArticle, v2.ArticleList](
			httpClient,
			baseURL+"/mediawatch.articles.v2.ArticlesService/StreamRelatedArticles",
			opts...,
		),
	}
}

// articlesServiceClient implements ArticlesServiceClient.
type articlesServiceClient struct {
	getArticle            *connect_go.Client[v2.QueryArticle, v2.Article]
	getArticles           *connect_go.Client[v2.QueryArticle, v2.ArticleList]
	streamArticles        *connect_go.Client[v2.QueryArticle, v2.ArticleList]
	streamRelatedArticles *connect_go.Client[v2.QueryArticle, v2.ArticleList]
}

// GetArticle calls mediawatch.articles.v2.ArticlesService.GetArticle.
func (c *articlesServiceClient) GetArticle(ctx context.Context, req *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.Article], error) {
	return c.getArticle.CallUnary(ctx, req)
}

// GetArticles calls mediawatch.articles.v2.ArticlesService.GetArticles.
func (c *articlesServiceClient) GetArticles(ctx context.Context, req *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.ArticleList], error) {
	return c.getArticles.CallUnary(ctx, req)
}

// StreamArticles calls mediawatch.articles.v2.ArticlesService.StreamArticles.
func (c *articlesServiceClient) StreamArticles(ctx context.Context, req *connect_go.Request[v2.QueryArticle]) (*connect_go.ServerStreamForClient[v2.ArticleList], error) {
	return c.streamArticles.CallServerStream(ctx, req)
}

// StreamRelatedArticles calls mediawatch.articles.v2.ArticlesService.StreamRelatedArticles.
func (c *articlesServiceClient) StreamRelatedArticles(ctx context.Context, req *connect_go.Request[v2.QueryArticle]) (*connect_go.ServerStreamForClient[v2.ArticleList], error) {
	return c.streamRelatedArticles.CallServerStream(ctx, req)
}

// ArticlesServiceHandler is an implementation of the mediawatch.articles.v2.ArticlesService
// service.
type ArticlesServiceHandler interface {
	// GetArticle
	GetArticle(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.Article], error)
	// GetArticles
	GetArticles(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.ArticleList], error)
	// StreamArticles
	StreamArticles(context.Context, *connect_go.Request[v2.QueryArticle], *connect_go.ServerStream[v2.ArticleList]) error
	// StreamRelatedArticles
	StreamRelatedArticles(context.Context, *connect_go.Request[v2.QueryArticle], *connect_go.ServerStream[v2.ArticleList]) error
}

// NewArticlesServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewArticlesServiceHandler(svc ArticlesServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/mediawatch.articles.v2.ArticlesService/GetArticle", connect_go.NewUnaryHandler(
		"/mediawatch.articles.v2.ArticlesService/GetArticle",
		svc.GetArticle,
		opts...,
	))
	mux.Handle("/mediawatch.articles.v2.ArticlesService/GetArticles", connect_go.NewUnaryHandler(
		"/mediawatch.articles.v2.ArticlesService/GetArticles",
		svc.GetArticles,
		opts...,
	))
	mux.Handle("/mediawatch.articles.v2.ArticlesService/StreamArticles", connect_go.NewServerStreamHandler(
		"/mediawatch.articles.v2.ArticlesService/StreamArticles",
		svc.StreamArticles,
		opts...,
	))
	mux.Handle("/mediawatch.articles.v2.ArticlesService/StreamRelatedArticles", connect_go.NewServerStreamHandler(
		"/mediawatch.articles.v2.ArticlesService/StreamRelatedArticles",
		svc.StreamRelatedArticles,
		opts...,
	))
	return "/mediawatch.articles.v2.ArticlesService/", mux
}

// UnimplementedArticlesServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedArticlesServiceHandler struct{}

func (UnimplementedArticlesServiceHandler) GetArticle(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.Article], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.articles.v2.ArticlesService.GetArticle is not implemented"))
}

func (UnimplementedArticlesServiceHandler) GetArticles(context.Context, *connect_go.Request[v2.QueryArticle]) (*connect_go.Response[v2.ArticleList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.articles.v2.ArticlesService.GetArticles is not implemented"))
}

func (UnimplementedArticlesServiceHandler) StreamArticles(context.Context, *connect_go.Request[v2.QueryArticle], *connect_go.ServerStream[v2.ArticleList]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.articles.v2.ArticlesService.StreamArticles is not implemented"))
}

func (UnimplementedArticlesServiceHandler) StreamRelatedArticles(context.Context, *connect_go.Request[v2.QueryArticle], *connect_go.ServerStream[v2.ArticleList]) error {
	return connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.articles.v2.ArticlesService.StreamRelatedArticles is not implemented"))
}
