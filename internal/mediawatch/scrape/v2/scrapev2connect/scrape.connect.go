// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/scrape/v2/scrape.proto

package scrapev2connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v2 "github.com/cvcio/mediawatch/pkg/mediawatch/scrape/v2"
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
	// ScrapeServiceName is the fully-qualified name of the ScrapeService service.
	ScrapeServiceName = "mediawatch.scrape.v2.ScrapeService"
)

// ScrapeServiceClient is a client for the mediawatch.scrape.v2.ScrapeService service.
type ScrapeServiceClient interface {
	// Endpoint Scrape
	Scrape(context.Context, *connect_go.Request[v2.ScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error)
	// Endpoint SimpleScrape
	SimpleScrape(context.Context, *connect_go.Request[v2.SimpleScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error)
	// Endpoint ReloadPassages
	ReloadPassages(context.Context, *connect_go.Request[v2.Empty]) (*connect_go.Response[v2.ReloadPassagesResponse], error)
}

// NewScrapeServiceClient constructs a client for the mediawatch.scrape.v2.ScrapeService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewScrapeServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) ScrapeServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &scrapeServiceClient{
		scrape: connect_go.NewClient[v2.ScrapeRequest, v2.ScrapeResponse](
			httpClient,
			baseURL+"/mediawatch.scrape.v2.ScrapeService/Scrape",
			opts...,
		),
		simpleScrape: connect_go.NewClient[v2.SimpleScrapeRequest, v2.ScrapeResponse](
			httpClient,
			baseURL+"/mediawatch.scrape.v2.ScrapeService/SimpleScrape",
			opts...,
		),
		reloadPassages: connect_go.NewClient[v2.Empty, v2.ReloadPassagesResponse](
			httpClient,
			baseURL+"/mediawatch.scrape.v2.ScrapeService/ReloadPassages",
			opts...,
		),
	}
}

// scrapeServiceClient implements ScrapeServiceClient.
type scrapeServiceClient struct {
	scrape         *connect_go.Client[v2.ScrapeRequest, v2.ScrapeResponse]
	simpleScrape   *connect_go.Client[v2.SimpleScrapeRequest, v2.ScrapeResponse]
	reloadPassages *connect_go.Client[v2.Empty, v2.ReloadPassagesResponse]
}

// Scrape calls mediawatch.scrape.v2.ScrapeService.Scrape.
func (c *scrapeServiceClient) Scrape(ctx context.Context, req *connect_go.Request[v2.ScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error) {
	return c.scrape.CallUnary(ctx, req)
}

// SimpleScrape calls mediawatch.scrape.v2.ScrapeService.SimpleScrape.
func (c *scrapeServiceClient) SimpleScrape(ctx context.Context, req *connect_go.Request[v2.SimpleScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error) {
	return c.simpleScrape.CallUnary(ctx, req)
}

// ReloadPassages calls mediawatch.scrape.v2.ScrapeService.ReloadPassages.
func (c *scrapeServiceClient) ReloadPassages(ctx context.Context, req *connect_go.Request[v2.Empty]) (*connect_go.Response[v2.ReloadPassagesResponse], error) {
	return c.reloadPassages.CallUnary(ctx, req)
}

// ScrapeServiceHandler is an implementation of the mediawatch.scrape.v2.ScrapeService service.
type ScrapeServiceHandler interface {
	// Endpoint Scrape
	Scrape(context.Context, *connect_go.Request[v2.ScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error)
	// Endpoint SimpleScrape
	SimpleScrape(context.Context, *connect_go.Request[v2.SimpleScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error)
	// Endpoint ReloadPassages
	ReloadPassages(context.Context, *connect_go.Request[v2.Empty]) (*connect_go.Response[v2.ReloadPassagesResponse], error)
}

// NewScrapeServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewScrapeServiceHandler(svc ScrapeServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/mediawatch.scrape.v2.ScrapeService/Scrape", connect_go.NewUnaryHandler(
		"/mediawatch.scrape.v2.ScrapeService/Scrape",
		svc.Scrape,
		opts...,
	))
	mux.Handle("/mediawatch.scrape.v2.ScrapeService/SimpleScrape", connect_go.NewUnaryHandler(
		"/mediawatch.scrape.v2.ScrapeService/SimpleScrape",
		svc.SimpleScrape,
		opts...,
	))
	mux.Handle("/mediawatch.scrape.v2.ScrapeService/ReloadPassages", connect_go.NewUnaryHandler(
		"/mediawatch.scrape.v2.ScrapeService/ReloadPassages",
		svc.ReloadPassages,
		opts...,
	))
	return "/mediawatch.scrape.v2.ScrapeService/", mux
}

// UnimplementedScrapeServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedScrapeServiceHandler struct{}

func (UnimplementedScrapeServiceHandler) Scrape(context.Context, *connect_go.Request[v2.ScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.scrape.v2.ScrapeService.Scrape is not implemented"))
}

func (UnimplementedScrapeServiceHandler) SimpleScrape(context.Context, *connect_go.Request[v2.SimpleScrapeRequest]) (*connect_go.Response[v2.ScrapeResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.scrape.v2.ScrapeService.SimpleScrape is not implemented"))
}

func (UnimplementedScrapeServiceHandler) ReloadPassages(context.Context, *connect_go.Request[v2.Empty]) (*connect_go.Response[v2.ReloadPassagesResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.scrape.v2.ScrapeService.ReloadPassages is not implemented"))
}
