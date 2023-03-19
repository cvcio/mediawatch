// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/feeds/v2/feed.proto

package feedsv2connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v21 "github.com/cvcio/mediawatch/internal/mediawatch/common/v2"
	v2 "github.com/cvcio/mediawatch/internal/mediawatch/feeds/v2"
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
	// FeedServiceName is the fully-qualified name of the FeedService service.
	FeedServiceName = "mediawatch.feeds.v2.FeedService"
)

// FeedServiceClient is a client for the mediawatch.feeds.v2.FeedService service.
type FeedServiceClient interface {
	// create a new feed
	CreateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v2.Feed], error)
	// get a single feed
	GetFeed(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.Feed], error)
	// get list of feeds by query
	GetFeeds(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error)
	// update a feed
	UpdateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// delete a feed
	DeleteFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// get the stream list
	GetFeedsStreamList(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error)
}

// NewFeedServiceClient constructs a client for the mediawatch.feeds.v2.FeedService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewFeedServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) FeedServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &feedServiceClient{
		createFeed: connect_go.NewClient[v2.Feed, v2.Feed](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/CreateFeed",
			opts...,
		),
		getFeed: connect_go.NewClient[v2.QueryFeed, v2.Feed](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/GetFeed",
			opts...,
		),
		getFeeds: connect_go.NewClient[v2.QueryFeed, v2.FeedList](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/GetFeeds",
			opts...,
		),
		updateFeed: connect_go.NewClient[v2.Feed, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/UpdateFeed",
			opts...,
		),
		deleteFeed: connect_go.NewClient[v2.Feed, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/DeleteFeed",
			opts...,
		),
		getFeedsStreamList: connect_go.NewClient[v2.QueryFeed, v2.FeedList](
			httpClient,
			baseURL+"/mediawatch.feeds.v2.FeedService/GetFeedsStreamList",
			opts...,
		),
	}
}

// feedServiceClient implements FeedServiceClient.
type feedServiceClient struct {
	createFeed         *connect_go.Client[v2.Feed, v2.Feed]
	getFeed            *connect_go.Client[v2.QueryFeed, v2.Feed]
	getFeeds           *connect_go.Client[v2.QueryFeed, v2.FeedList]
	updateFeed         *connect_go.Client[v2.Feed, v21.ResponseWithMessage]
	deleteFeed         *connect_go.Client[v2.Feed, v21.ResponseWithMessage]
	getFeedsStreamList *connect_go.Client[v2.QueryFeed, v2.FeedList]
}

// CreateFeed calls mediawatch.feeds.v2.FeedService.CreateFeed.
func (c *feedServiceClient) CreateFeed(ctx context.Context, req *connect_go.Request[v2.Feed]) (*connect_go.Response[v2.Feed], error) {
	return c.createFeed.CallUnary(ctx, req)
}

// GetFeed calls mediawatch.feeds.v2.FeedService.GetFeed.
func (c *feedServiceClient) GetFeed(ctx context.Context, req *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.Feed], error) {
	return c.getFeed.CallUnary(ctx, req)
}

// GetFeeds calls mediawatch.feeds.v2.FeedService.GetFeeds.
func (c *feedServiceClient) GetFeeds(ctx context.Context, req *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error) {
	return c.getFeeds.CallUnary(ctx, req)
}

// UpdateFeed calls mediawatch.feeds.v2.FeedService.UpdateFeed.
func (c *feedServiceClient) UpdateFeed(ctx context.Context, req *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.updateFeed.CallUnary(ctx, req)
}

// DeleteFeed calls mediawatch.feeds.v2.FeedService.DeleteFeed.
func (c *feedServiceClient) DeleteFeed(ctx context.Context, req *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.deleteFeed.CallUnary(ctx, req)
}

// GetFeedsStreamList calls mediawatch.feeds.v2.FeedService.GetFeedsStreamList.
func (c *feedServiceClient) GetFeedsStreamList(ctx context.Context, req *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error) {
	return c.getFeedsStreamList.CallUnary(ctx, req)
}

// FeedServiceHandler is an implementation of the mediawatch.feeds.v2.FeedService service.
type FeedServiceHandler interface {
	// create a new feed
	CreateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v2.Feed], error)
	// get a single feed
	GetFeed(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.Feed], error)
	// get list of feeds by query
	GetFeeds(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error)
	// update a feed
	UpdateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// delete a feed
	DeleteFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// get the stream list
	GetFeedsStreamList(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error)
}

// NewFeedServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewFeedServiceHandler(svc FeedServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/mediawatch.feeds.v2.FeedService/CreateFeed", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/CreateFeed",
		svc.CreateFeed,
		opts...,
	))
	mux.Handle("/mediawatch.feeds.v2.FeedService/GetFeed", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/GetFeed",
		svc.GetFeed,
		opts...,
	))
	mux.Handle("/mediawatch.feeds.v2.FeedService/GetFeeds", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/GetFeeds",
		svc.GetFeeds,
		opts...,
	))
	mux.Handle("/mediawatch.feeds.v2.FeedService/UpdateFeed", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/UpdateFeed",
		svc.UpdateFeed,
		opts...,
	))
	mux.Handle("/mediawatch.feeds.v2.FeedService/DeleteFeed", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/DeleteFeed",
		svc.DeleteFeed,
		opts...,
	))
	mux.Handle("/mediawatch.feeds.v2.FeedService/GetFeedsStreamList", connect_go.NewUnaryHandler(
		"/mediawatch.feeds.v2.FeedService/GetFeedsStreamList",
		svc.GetFeedsStreamList,
		opts...,
	))
	return "/mediawatch.feeds.v2.FeedService/", mux
}

// UnimplementedFeedServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedFeedServiceHandler struct{}

func (UnimplementedFeedServiceHandler) CreateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v2.Feed], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.CreateFeed is not implemented"))
}

func (UnimplementedFeedServiceHandler) GetFeed(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.Feed], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.GetFeed is not implemented"))
}

func (UnimplementedFeedServiceHandler) GetFeeds(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.GetFeeds is not implemented"))
}

func (UnimplementedFeedServiceHandler) UpdateFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.UpdateFeed is not implemented"))
}

func (UnimplementedFeedServiceHandler) DeleteFeed(context.Context, *connect_go.Request[v2.Feed]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.DeleteFeed is not implemented"))
}

func (UnimplementedFeedServiceHandler) GetFeedsStreamList(context.Context, *connect_go.Request[v2.QueryFeed]) (*connect_go.Response[v2.FeedList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.feeds.v2.FeedService.GetFeedsStreamList is not implemented"))
}
