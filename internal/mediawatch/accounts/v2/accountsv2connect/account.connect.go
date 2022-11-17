// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: mediawatch/accounts/v2/account.proto

package accountsv2connect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	v2 "github.com/cvcio/mediawatch/internal/mediawatch/accounts/v2"
	v21 "github.com/cvcio/mediawatch/internal/mediawatch/common/v2"
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
	// AccountServiceName is the fully-qualified name of the AccountService service.
	AccountServiceName = "mediawatch.accounts.v2.AccountService"
)

// AccountServiceClient is a client for the mediawatch.accounts.v2.AccountService service.
type AccountServiceClient interface {
	// create new account
	Create(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v2.Account], error)
	// get account with query
	Get(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.Account], error)
	// get list of accounts with query
	List(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.AccountList], error)
	// update an account
	Update(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// update account with fields
	UpdateWithFields(context.Context, *connect_go.Request[v2.AccountWithFields]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// delete an account
	Delete(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	Verify(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	Reset(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
}

// NewAccountServiceClient constructs a client for the mediawatch.accounts.v2.AccountService
// service. By default, it uses the Connect protocol with the binary Protobuf Codec, asks for
// gzipped responses, and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply
// the connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccountServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AccountServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accountServiceClient{
		create: connect_go.NewClient[v2.Account, v2.Account](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Create",
			opts...,
		),
		get: connect_go.NewClient[v2.QueryAccount, v2.Account](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Get",
			opts...,
		),
		list: connect_go.NewClient[v2.QueryAccount, v2.AccountList](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/List",
			opts...,
		),
		update: connect_go.NewClient[v2.Account, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Update",
			opts...,
		),
		updateWithFields: connect_go.NewClient[v2.AccountWithFields, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/UpdateWithFields",
			opts...,
		),
		delete: connect_go.NewClient[v2.Account, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Delete",
			opts...,
		),
		verify: connect_go.NewClient[v2.Account, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Verify",
			opts...,
		),
		reset: connect_go.NewClient[v2.Account, v21.ResponseWithMessage](
			httpClient,
			baseURL+"/mediawatch.accounts.v2.AccountService/Reset",
			opts...,
		),
	}
}

// accountServiceClient implements AccountServiceClient.
type accountServiceClient struct {
	create           *connect_go.Client[v2.Account, v2.Account]
	get              *connect_go.Client[v2.QueryAccount, v2.Account]
	list             *connect_go.Client[v2.QueryAccount, v2.AccountList]
	update           *connect_go.Client[v2.Account, v21.ResponseWithMessage]
	updateWithFields *connect_go.Client[v2.AccountWithFields, v21.ResponseWithMessage]
	delete           *connect_go.Client[v2.Account, v21.ResponseWithMessage]
	verify           *connect_go.Client[v2.Account, v21.ResponseWithMessage]
	reset            *connect_go.Client[v2.Account, v21.ResponseWithMessage]
}

// Create calls mediawatch.accounts.v2.AccountService.Create.
func (c *accountServiceClient) Create(ctx context.Context, req *connect_go.Request[v2.Account]) (*connect_go.Response[v2.Account], error) {
	return c.create.CallUnary(ctx, req)
}

// Get calls mediawatch.accounts.v2.AccountService.Get.
func (c *accountServiceClient) Get(ctx context.Context, req *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.Account], error) {
	return c.get.CallUnary(ctx, req)
}

// List calls mediawatch.accounts.v2.AccountService.List.
func (c *accountServiceClient) List(ctx context.Context, req *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.AccountList], error) {
	return c.list.CallUnary(ctx, req)
}

// Update calls mediawatch.accounts.v2.AccountService.Update.
func (c *accountServiceClient) Update(ctx context.Context, req *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.update.CallUnary(ctx, req)
}

// UpdateWithFields calls mediawatch.accounts.v2.AccountService.UpdateWithFields.
func (c *accountServiceClient) UpdateWithFields(ctx context.Context, req *connect_go.Request[v2.AccountWithFields]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.updateWithFields.CallUnary(ctx, req)
}

// Delete calls mediawatch.accounts.v2.AccountService.Delete.
func (c *accountServiceClient) Delete(ctx context.Context, req *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.delete.CallUnary(ctx, req)
}

// Verify calls mediawatch.accounts.v2.AccountService.Verify.
func (c *accountServiceClient) Verify(ctx context.Context, req *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.verify.CallUnary(ctx, req)
}

// Reset calls mediawatch.accounts.v2.AccountService.Reset.
func (c *accountServiceClient) Reset(ctx context.Context, req *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return c.reset.CallUnary(ctx, req)
}

// AccountServiceHandler is an implementation of the mediawatch.accounts.v2.AccountService service.
type AccountServiceHandler interface {
	// create new account
	Create(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v2.Account], error)
	// get account with query
	Get(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.Account], error)
	// get list of accounts with query
	List(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.AccountList], error)
	// update an account
	Update(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// update account with fields
	UpdateWithFields(context.Context, *connect_go.Request[v2.AccountWithFields]) (*connect_go.Response[v21.ResponseWithMessage], error)
	// delete an account
	Delete(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	Verify(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
	Reset(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error)
}

// NewAccountServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccountServiceHandler(svc AccountServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/mediawatch.accounts.v2.AccountService/Create", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Create",
		svc.Create,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/Get", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Get",
		svc.Get,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/List", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/List",
		svc.List,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/Update", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Update",
		svc.Update,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/UpdateWithFields", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/UpdateWithFields",
		svc.UpdateWithFields,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/Delete", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Delete",
		svc.Delete,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/Verify", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Verify",
		svc.Verify,
		opts...,
	))
	mux.Handle("/mediawatch.accounts.v2.AccountService/Reset", connect_go.NewUnaryHandler(
		"/mediawatch.accounts.v2.AccountService/Reset",
		svc.Reset,
		opts...,
	))
	return "/mediawatch.accounts.v2.AccountService/", mux
}

// UnimplementedAccountServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccountServiceHandler struct{}

func (UnimplementedAccountServiceHandler) Create(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v2.Account], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Create is not implemented"))
}

func (UnimplementedAccountServiceHandler) Get(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.Account], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Get is not implemented"))
}

func (UnimplementedAccountServiceHandler) List(context.Context, *connect_go.Request[v2.QueryAccount]) (*connect_go.Response[v2.AccountList], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.List is not implemented"))
}

func (UnimplementedAccountServiceHandler) Update(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Update is not implemented"))
}

func (UnimplementedAccountServiceHandler) UpdateWithFields(context.Context, *connect_go.Request[v2.AccountWithFields]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.UpdateWithFields is not implemented"))
}

func (UnimplementedAccountServiceHandler) Delete(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Delete is not implemented"))
}

func (UnimplementedAccountServiceHandler) Verify(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Verify is not implemented"))
}

func (UnimplementedAccountServiceHandler) Reset(context.Context, *connect_go.Request[v2.Account]) (*connect_go.Response[v21.ResponseWithMessage], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("mediawatch.accounts.v2.AccountService.Reset is not implemented"))
}
