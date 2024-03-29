// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/enrich/v2/enrich.proto

package enrichv2connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v2 "github.com/cvcio/mediawatch/pkg/mediawatch/enrich/v2"
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
	// EnrichServiceName is the fully-qualified name of the EnrichService service.
	EnrichServiceName = "mediawatch.enrich.v2.EnrichService"
)

// EnrichServiceClient is a client for the mediawatch.enrich.v2.EnrichService service.
type EnrichServiceClient interface {
	NLP(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	StopWords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Keywords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Entities(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Summary(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Topics(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Quotes(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Claims(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
}

// NewEnrichServiceClient constructs a client for the mediawatch.enrich.v2.EnrichService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewEnrichServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) EnrichServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &enrichServiceClient{
		nLP: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/NLP",
			opts...,
		),
		stopWords: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/StopWords",
			opts...,
		),
		keywords: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Keywords",
			opts...,
		),
		entities: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Entities",
			opts...,
		),
		summary: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Summary",
			opts...,
		),
		topics: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Topics",
			opts...,
		),
		quotes: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Quotes",
			opts...,
		),
		claims: connect_go.NewClient[v2.EnrichRequest, v2.EnrichResponse](
			httpClient,
			baseURL+"/mediawatch.enrich.v2.EnrichService/Claims",
			opts...,
		),
	}
}

// enrichServiceClient implements EnrichServiceClient.
type enrichServiceClient struct {
	nLP       *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	stopWords *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	keywords  *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	entities  *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	summary   *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	topics    *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	quotes    *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
	claims    *connect_go.Client[v2.EnrichRequest, v2.EnrichResponse]
}

// NLP calls mediawatch.enrich.v2.EnrichService.NLP.
func (c *enrichServiceClient) NLP(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.nLP.CallUnary(ctx, req)
}

// StopWords calls mediawatch.enrich.v2.EnrichService.StopWords.
func (c *enrichServiceClient) StopWords(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.stopWords.CallUnary(ctx, req)
}

// Keywords calls mediawatch.enrich.v2.EnrichService.Keywords.
func (c *enrichServiceClient) Keywords(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.keywords.CallUnary(ctx, req)
}

// Entities calls mediawatch.enrich.v2.EnrichService.Entities.
func (c *enrichServiceClient) Entities(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.entities.CallUnary(ctx, req)
}

// Summary calls mediawatch.enrich.v2.EnrichService.Summary.
func (c *enrichServiceClient) Summary(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.summary.CallUnary(ctx, req)
}

// Topics calls mediawatch.enrich.v2.EnrichService.Topics.
func (c *enrichServiceClient) Topics(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.topics.CallUnary(ctx, req)
}

// Quotes calls mediawatch.enrich.v2.EnrichService.Quotes.
func (c *enrichServiceClient) Quotes(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.quotes.CallUnary(ctx, req)
}

// Claims calls mediawatch.enrich.v2.EnrichService.Claims.
func (c *enrichServiceClient) Claims(ctx context.Context, req *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return c.claims.CallUnary(ctx, req)
}

// EnrichServiceHandler is an implementation of the mediawatch.enrich.v2.EnrichService service.
type EnrichServiceHandler interface {
	NLP(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	StopWords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Keywords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Entities(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Summary(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Topics(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Quotes(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
	Claims(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error)
}

// NewEnrichServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewEnrichServiceHandler(svc EnrichServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/mediawatch.enrich.v2.EnrichService/NLP", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/NLP",
		svc.NLP,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/StopWords", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/StopWords",
		svc.StopWords,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Keywords", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Keywords",
		svc.Keywords,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Entities", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Entities",
		svc.Entities,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Summary", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Summary",
		svc.Summary,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Topics", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Topics",
		svc.Topics,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Quotes", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Quotes",
		svc.Quotes,
		opts...,
	))
	mux.Handle("/mediawatch.enrich.v2.EnrichService/Claims", connect_go.NewUnaryHandler(
		"/mediawatch.enrich.v2.EnrichService/Claims",
		svc.Claims,
		opts...,
	))
	return "/mediawatch.enrich.v2.EnrichService/", mux
}

// UnimplementedEnrichServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedEnrichServiceHandler struct{}

func (UnimplementedEnrichServiceHandler) NLP(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.NLP is not implemented"))
}

func (UnimplementedEnrichServiceHandler) StopWords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.StopWords is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Keywords(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Keywords is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Entities(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Entities is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Summary(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Summary is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Topics(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Topics is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Quotes(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Quotes is not implemented"))
}

func (UnimplementedEnrichServiceHandler) Claims(context.Context, *connect_go.Request[v2.EnrichRequest]) (*connect_go.Response[v2.EnrichResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.enrich.v2.EnrichService.Claims is not implemented"))
}
