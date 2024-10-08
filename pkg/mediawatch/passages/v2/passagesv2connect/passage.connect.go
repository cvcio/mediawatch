// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/passages/v2/passage.proto

package passagesv2connect

import (
	connect "connectrpc.com/connect"
	context "context"
	errors "errors"
	v21 "github.com/cvcio/mediawatch/pkg/mediawatch/common/v2"
	v2 "github.com/cvcio/mediawatch/pkg/mediawatch/passages/v2"
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
	// PassageServiceName is the fully-qualified name of the PassageService service.
	PassageServiceName = "mediawatch.passages.v2.PassageService"
)

// These constants are the fully-qualified names of the RPCs defined in this package. They're
// exposed at runtime as Spec.Procedure and as the final two segments of the HTTP route.
//
// Note that these are different from the fully-qualified method names used by
// google.golang.org/protobuf/reflect/protoreflect. To convert from these constants to
// reflection-formatted method names, remove the leading slash and convert the remaining slash to a
// period.
const (
	// PassageServiceCreatePassageProcedure is the fully-qualified name of the PassageService's
	// CreatePassage RPC.
	PassageServiceCreatePassageProcedure = "/mediawatch.passages.v2.PassageService/CreatePassage"
	// PassageServiceGetPassagesProcedure is the fully-qualified name of the PassageService's
	// GetPassages RPC.
	PassageServiceGetPassagesProcedure = "/mediawatch.passages.v2.PassageService/GetPassages"
	// PassageServiceDeletePassageProcedure is the fully-qualified name of the PassageService's
	// DeletePassage RPC.
	PassageServiceDeletePassageProcedure = "/mediawatch.passages.v2.PassageService/DeletePassage"
)

// These variables are the protoreflect.Descriptor objects for the RPCs defined in this package.
var (
	passageServiceServiceDescriptor             = v2.File_mediawatch_passages_v2_passage_proto.Services().ByName("PassageService")
	passageServiceCreatePassageMethodDescriptor = passageServiceServiceDescriptor.Methods().ByName("CreatePassage")
	passageServiceGetPassagesMethodDescriptor   = passageServiceServiceDescriptor.Methods().ByName("GetPassages")
	passageServiceDeletePassageMethodDescriptor = passageServiceServiceDescriptor.Methods().ByName("DeletePassage")
)

// PassageServiceClient is a client for the mediawatch.passages.v2.PassageService service.
type PassageServiceClient interface {
	// create a new passage
	CreatePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v2.Passage], error)
	// get list of passages by query
	GetPassages(context.Context, *connect.Request[v2.QueryPassage]) (*connect.Response[v2.PassageList], error)
	// delete a passage by id
	DeletePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v21.ResponseWithMessage], error)
}

// NewPassageServiceClient constructs a client for the mediawatch.passages.v2.PassageService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewPassageServiceClient(httpClient connect.HTTPClient, baseURL string, opts ...connect.ClientOption) PassageServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &passageServiceClient{
		createPassage: connect.NewClient[v2.Passage, v2.Passage](
			httpClient,
			baseURL+PassageServiceCreatePassageProcedure,
			connect.WithSchema(passageServiceCreatePassageMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		getPassages: connect.NewClient[v2.QueryPassage, v2.PassageList](
			httpClient,
			baseURL+PassageServiceGetPassagesProcedure,
			connect.WithSchema(passageServiceGetPassagesMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
		deletePassage: connect.NewClient[v2.Passage, v21.ResponseWithMessage](
			httpClient,
			baseURL+PassageServiceDeletePassageProcedure,
			connect.WithSchema(passageServiceDeletePassageMethodDescriptor),
			connect.WithClientOptions(opts...),
		),
	}
}

// passageServiceClient implements PassageServiceClient.
type passageServiceClient struct {
	createPassage *connect.Client[v2.Passage, v2.Passage]
	getPassages   *connect.Client[v2.QueryPassage, v2.PassageList]
	deletePassage *connect.Client[v2.Passage, v21.ResponseWithMessage]
}

// CreatePassage calls mediawatch.passages.v2.PassageService.CreatePassage.
func (c *passageServiceClient) CreatePassage(ctx context.Context, req *connect.Request[v2.Passage]) (*connect.Response[v2.Passage], error) {
	return c.createPassage.CallUnary(ctx, req)
}

// GetPassages calls mediawatch.passages.v2.PassageService.GetPassages.
func (c *passageServiceClient) GetPassages(ctx context.Context, req *connect.Request[v2.QueryPassage]) (*connect.Response[v2.PassageList], error) {
	return c.getPassages.CallUnary(ctx, req)
}

// DeletePassage calls mediawatch.passages.v2.PassageService.DeletePassage.
func (c *passageServiceClient) DeletePassage(ctx context.Context, req *connect.Request[v2.Passage]) (*connect.Response[v21.ResponseWithMessage], error) {
	return c.deletePassage.CallUnary(ctx, req)
}

// PassageServiceHandler is an implementation of the mediawatch.passages.v2.PassageService service.
type PassageServiceHandler interface {
	// create a new passage
	CreatePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v2.Passage], error)
	// get list of passages by query
	GetPassages(context.Context, *connect.Request[v2.QueryPassage]) (*connect.Response[v2.PassageList], error)
	// delete a passage by id
	DeletePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v21.ResponseWithMessage], error)
}

// NewPassageServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewPassageServiceHandler(svc PassageServiceHandler, opts ...connect.HandlerOption) (string, http.Handler) {
	passageServiceCreatePassageHandler := connect.NewUnaryHandler(
		PassageServiceCreatePassageProcedure,
		svc.CreatePassage,
		connect.WithSchema(passageServiceCreatePassageMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	passageServiceGetPassagesHandler := connect.NewUnaryHandler(
		PassageServiceGetPassagesProcedure,
		svc.GetPassages,
		connect.WithSchema(passageServiceGetPassagesMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	passageServiceDeletePassageHandler := connect.NewUnaryHandler(
		PassageServiceDeletePassageProcedure,
		svc.DeletePassage,
		connect.WithSchema(passageServiceDeletePassageMethodDescriptor),
		connect.WithHandlerOptions(opts...),
	)
	return "/mediawatch.passages.v2.PassageService/", http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case PassageServiceCreatePassageProcedure:
			passageServiceCreatePassageHandler.ServeHTTP(w, r)
		case PassageServiceGetPassagesProcedure:
			passageServiceGetPassagesHandler.ServeHTTP(w, r)
		case PassageServiceDeletePassageProcedure:
			passageServiceDeletePassageHandler.ServeHTTP(w, r)
		default:
			http.NotFound(w, r)
		}
	})
}

// UnimplementedPassageServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedPassageServiceHandler struct{}

func (UnimplementedPassageServiceHandler) CreatePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v2.Passage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.passages.v2.PassageService.CreatePassage is not implemented"))
}

func (UnimplementedPassageServiceHandler) GetPassages(context.Context, *connect.Request[v2.QueryPassage]) (*connect.Response[v2.PassageList], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.passages.v2.PassageService.GetPassages is not implemented"))
}

func (UnimplementedPassageServiceHandler) DeletePassage(context.Context, *connect.Request[v2.Passage]) (*connect.Response[v21.ResponseWithMessage], error) {
	return nil, connect.NewError(connect.CodeUnimplemented, errors.New("mediawatch.passages.v2.PassageService.DeletePassage is not implemented"))
}
