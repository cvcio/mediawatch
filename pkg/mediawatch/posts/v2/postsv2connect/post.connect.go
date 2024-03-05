// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/posts/v2/post.proto

package postsv2connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v2 "github.com/cvcio/mediawatch/pkg/mediawatch/posts/v2"
	http "net/http"
	strings "strings"
)

// This is a compile-time assertion to ensure that this generated file and the connect package are
// compatible. If you get a compiler error that this constant is not defined, this code was
// generated with a version of connect newer than the one compiled into your binary. You can fix the
// problem by either regenerating this code with an older version of connect or updating the connect
// version compiled into your binary.
const _ = connect.IsAtLeastVersion1_13_0

const (
	// PostServiceName is the fully-qualified name of the PostService service.
	PostServiceName = "mediawatch.posts.v2.PostService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PostServiceGetPostProcedure is the fully-qualified name of the PostService's GetPost RPC.
	PostServiceGetPostProcedure = "/mediawatch.posts.v2.PostService/GetPost"
	// PostServiceGetPostsProcedure is the fully-qualified name of the PostService's GetPosts RPC.
	PostServiceGetPostsProcedure = "/mediawatch.posts.v2.PostService/GetPosts"
	// PostServiceCreatePostProcedure is the fully-qualified name of the PostService's CreatePost RPC.
	PostServiceCreatePostProcedure = "/mediawatch.posts.v2.PostService/CreatePost"
	// PostServiceUpdatePostProcedure is the fully-qualified name of the PostService's UpdatePost RPC.
	PostServiceUpdatePostProcedure = "/mediawatch.posts.v2.PostService/UpdatePost"
	// PostServiceDeletePostProcedure is the fully-qualified name of the PostService's DeletePost RPC.
	PostServiceDeletePostProcedure = "/mediawatch.posts.v2.PostService/DeletePost"
	// PostServiceStreamPostsProcedure is the fully-qualified name of the PostService's StreamPosts RPC.
	PostServiceStreamPostsProcedure = "/mediawatch.posts.v2.PostService/StreamPosts"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	postServiceServiceDescriptor           = v2.File_mediawatch_posts_v2_post_proto.Services().ByName("PostService")
	postServiceGetPostMethodDescriptor     = postServiceServiceDescriptor.Methods().ByName("GetPost")
	postServiceGetPostsMethodDescriptor    = postServiceServiceDescriptor.Methods().ByName("GetPosts")
	postServiceCreatePostMethodDescriptor  = postServiceServiceDescriptor.Methods().ByName("CreatePost")
	postServiceUpdatePostMethodDescriptor  = postServiceServiceDescriptor.Methods().ByName("UpdatePost")
	postServiceDeletePostMethodDescriptor  = postServiceServiceDescriptor.Methods().ByName("DeletePost")
	postServiceStreamPostsMethodDescriptor = postServiceServiceDescriptor.Methods().ByName("StreamPosts")
)

// PostServiceClient is a client for the mediawatch.posts.v2.PostService service.
type PostServiceClient interface {
	GetPost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	GetPosts(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	CreatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	UpdatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	DeletePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	StreamPosts(context.Context, *connect.Request[v2.PostRequest]) (*connect.ServerStreamForClient[v2.PostResponse], error)
}

// NewPostServiceClient constructs a client for the mediawatch.posts.v2.PostService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPostServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PostServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &postServiceClient{
		getPost: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceGetPostProcedure,
			connect.WithSchema(postServiceGetPostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPosts: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceGetPostsProcedure,
			connect.WithSchema(postServiceGetPostsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		createPost: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceCreatePostProcedure,
			connect.WithSchema(postServiceCreatePostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		updatePost: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceUpdatePostProcedure,
			connect.WithSchema(postServiceUpdatePostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deletePost: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceDeletePostProcedure,
			connect.WithSchema(postServiceDeletePostMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		streamPosts: connect.NewClient[v2.PostRequest, v2.PostResponse](
			httpClient,
			baseURL+PostServiceStreamPostsProcedure,
			connect.WithSchema(postServiceStreamPostsMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// postServiceClient implements PostServiceClient.
type postServiceClient struct {
	getPost     *connect.Client[v2.PostRequest, v2.PostResponse]
	getPosts    *connect.Client[v2.PostRequest, v2.PostResponse]
	createPost  *connect.Client[v2.PostRequest, v2.PostResponse]
	updatePost  *connect.Client[v2.PostRequest, v2.PostResponse]
	deletePost  *connect.Client[v2.PostRequest, v2.PostResponse]
	streamPosts *connect.Client[v2.PostRequest, v2.PostResponse]
}

// GetPost calls mediawatch.posts.v2.PostService.GetPost.
func (c *postServiceClient) GetPost(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return c.getPost.CallUnary(ctx, req)
}

// GetPosts calls mediawatch.posts.v2.PostService.GetPosts.
func (c *postServiceClient) GetPosts(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return c.getPosts.CallUnary(ctx, req)
}

// CreatePost calls mediawatch.posts.v2.PostService.CreatePost.
func (c *postServiceClient) CreatePost(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return c.createPost.CallUnary(ctx, req)
}

// UpdatePost calls mediawatch.posts.v2.PostService.UpdatePost.
func (c *postServiceClient) UpdatePost(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return c.updatePost.CallUnary(ctx, req)
}

// DeletePost calls mediawatch.posts.v2.PostService.DeletePost.
func (c *postServiceClient) DeletePost(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return c.deletePost.CallUnary(ctx, req)
}

// StreamPosts calls mediawatch.posts.v2.PostService.StreamPosts.
func (c *postServiceClient) StreamPosts(ctx context.Context, req *connect.Request[v2.PostRequest]) (*connect.ServerStreamForClient[v2.PostResponse], error) {
	return c.streamPosts.CallServerStream(ctx, req)
}

// PostServiceHandler is an implementation of the mediawatch.posts.v2.PostService service.
type PostServiceHandler interface {
	GetPost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	GetPosts(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	CreatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	UpdatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	DeletePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error)
	StreamPosts(context.Context, *connect.Request[v2.PostRequest], *connect.ServerStream[v2.PostResponse]) error
}

// NewPostServiceHandler builds an HTTP handler from the service implementation. It returns the path
// on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPostServiceHandler(svc PostServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	postServiceGetPostHandler := connect.NewUnaryHandler(
		PostServiceGetPostProcedure,
		svc.GetPost,
		connect.WithSchema(postServiceGetPostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	postServiceGetPostsHandler := connect.NewUnaryHandler(
		PostServiceGetPostsProcedure,
		svc.GetPosts,
		connect.WithSchema(postServiceGetPostsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	postServiceCreatePostHandler := connect.NewUnaryHandler(
		PostServiceCreatePostProcedure,
		svc.CreatePost,
		connect.WithSchema(postServiceCreatePostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	postServiceUpdatePostHandler := connect.NewUnaryHandler(
		PostServiceUpdatePostProcedure,
		svc.UpdatePost,
		connect.WithSchema(postServiceUpdatePostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	postServiceDeletePostHandler := connect.NewUnaryHandler(
		PostServiceDeletePostProcedure,
		svc.DeletePost,
		connect.WithSchema(postServiceDeletePostMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	postServiceStreamPostsHandler := connect.NewServerStreamHandler(
		PostServiceStreamPostsProcedure,
		svc.StreamPosts,
		connect.WithSchema(postServiceStreamPostsMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/mediawatch.posts.v2.PostService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PostServiceGetPostProcedure:
			postServiceGetPostHandler.ServeHTTP(w, r)
		case PostServiceGetPostsProcedure:
			postServiceGetPostsHandler.ServeHTTP(w, r)
		case PostServiceCreatePostProcedure:
			postServiceCreatePostHandler.ServeHTTP(w, r)
		case PostServiceUpdatePostProcedure:
			postServiceUpdatePostHandler.ServeHTTP(w, r)
		case PostServiceDeletePostProcedure:
			postServiceDeletePostHandler.ServeHTTP(w, r)
		case PostServiceStreamPostsProcedure:
			postServiceStreamPostsHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPostServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPostServiceHandler struct{}

func (UnimplementedPostServiceHandler) GetPost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.GetPost is not implemented"))
}

func (UnimplementedPostServiceHandler) GetPosts(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.GetPosts is not implemented"))
}

func (UnimplementedPostServiceHandler) CreatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.CreatePost is not implemented"))
}

func (UnimplementedPostServiceHandler) UpdatePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.UpdatePost is not implemented"))
}

func (UnimplementedPostServiceHandler) DeletePost(context.Context, *connect.Request[v2.PostRequest]) (*connect.Response[v2.PostResponse], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.DeletePost is not implemented"))
}

func (UnimplementedPostServiceHandler) StreamPosts(context.Context, *connect.Request[v2.PostRequest], *connect.ServerStream[v2.PostResponse]) error {
	return connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.posts.v2.PostService.StreamPosts is not implemented"))
}
