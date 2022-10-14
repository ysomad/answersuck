// Code generated by protoc-gen-connect-go. DO NOT EDIT.
//
// Source: user/account/service.proto

package accountconnect

import (
	context "context"
	errors "errors"
	connect_go "github.com/bufbuild/connect-go"
	account "github.com/ysomad/answersuck/api/user/account"
	emptypb "google.golang.org/protobuf/types/known/emptypb"
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
	AccountServiceName = "user.account.AccountService"
)

// AccountServiceClient is a client for the user.account.AccountService service.
type AccountServiceClient interface {
	CreateAccount(context.Context, *connect_go.Request[account.CreateAccountRequest]) (*connect_go.Response[account.CreateAccountResponse], error)
	GetAccountById(context.Context, *connect_go.Request[account.GetAccountByIdRequest]) (*connect_go.Response[account.GetAccountByIdResponse], error)
	GetAccountByEmail(context.Context, *connect_go.Request[account.GetAccountByEmailRequest]) (*connect_go.Response[account.GetAccountByEmailResponse], error)
	DeleteAccount(context.Context, *connect_go.Request[account.DeleteAccountRequest]) (*connect_go.Response[emptypb.Empty], error)
}

// NewAccountServiceClient constructs a client for the user.account.AccountService service. By
// default, it uses the Connect protocol with the binary Protobuf Codec, asks for gzipped responses,
// and sends uncompressed requests. To use the gRPC or gRPC-Web protocols, supply the
// connect.WithGRPC() or connect.WithGRPCWeb() options.
//
// The URL supplied here should be the base URL for the Connect or gRPC server (for example,
// http://api.acme.com or https://acme.com/grpc).
func NewAccountServiceClient(httpClient connect_go.HTTPClient, baseURL string, opts ...connect_go.ClientOption) AccountServiceClient {
	baseURL = strings.TrimRight(baseURL, "/")
	return &accountServiceClient{
		createAccount: connect_go.NewClient[account.CreateAccountRequest, account.CreateAccountResponse](
			httpClient,
			baseURL+"/user.account.AccountService/CreateAccount",
			opts...,
		),
		getAccountById: connect_go.NewClient[account.GetAccountByIdRequest, account.GetAccountByIdResponse](
			httpClient,
			baseURL+"/user.account.AccountService/GetAccountById",
			opts...,
		),
		getAccountByEmail: connect_go.NewClient[account.GetAccountByEmailRequest, account.GetAccountByEmailResponse](
			httpClient,
			baseURL+"/user.account.AccountService/GetAccountByEmail",
			opts...,
		),
		deleteAccount: connect_go.NewClient[account.DeleteAccountRequest, emptypb.Empty](
			httpClient,
			baseURL+"/user.account.AccountService/DeleteAccount",
			opts...,
		),
	}
}

// accountServiceClient implements AccountServiceClient.
type accountServiceClient struct {
	createAccount     *connect_go.Client[account.CreateAccountRequest, account.CreateAccountResponse]
	getAccountById    *connect_go.Client[account.GetAccountByIdRequest, account.GetAccountByIdResponse]
	getAccountByEmail *connect_go.Client[account.GetAccountByEmailRequest, account.GetAccountByEmailResponse]
	deleteAccount     *connect_go.Client[account.DeleteAccountRequest, emptypb.Empty]
}

// CreateAccount calls user.account.AccountService.CreateAccount.
func (c *accountServiceClient) CreateAccount(ctx context.Context, req *connect_go.Request[account.CreateAccountRequest]) (*connect_go.Response[account.CreateAccountResponse], error) {
	return c.createAccount.CallUnary(ctx, req)
}

// GetAccountById calls user.account.AccountService.GetAccountById.
func (c *accountServiceClient) GetAccountById(ctx context.Context, req *connect_go.Request[account.GetAccountByIdRequest]) (*connect_go.Response[account.GetAccountByIdResponse], error) {
	return c.getAccountById.CallUnary(ctx, req)
}

// GetAccountByEmail calls user.account.AccountService.GetAccountByEmail.
func (c *accountServiceClient) GetAccountByEmail(ctx context.Context, req *connect_go.Request[account.GetAccountByEmailRequest]) (*connect_go.Response[account.GetAccountByEmailResponse], error) {
	return c.getAccountByEmail.CallUnary(ctx, req)
}

// DeleteAccount calls user.account.AccountService.DeleteAccount.
func (c *accountServiceClient) DeleteAccount(ctx context.Context, req *connect_go.Request[account.DeleteAccountRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return c.deleteAccount.CallUnary(ctx, req)
}

// AccountServiceHandler is an implementation of the user.account.AccountService service.
type AccountServiceHandler interface {
	CreateAccount(context.Context, *connect_go.Request[account.CreateAccountRequest]) (*connect_go.Response[account.CreateAccountResponse], error)
	GetAccountById(context.Context, *connect_go.Request[account.GetAccountByIdRequest]) (*connect_go.Response[account.GetAccountByIdResponse], error)
	GetAccountByEmail(context.Context, *connect_go.Request[account.GetAccountByEmailRequest]) (*connect_go.Response[account.GetAccountByEmailResponse], error)
	DeleteAccount(context.Context, *connect_go.Request[account.DeleteAccountRequest]) (*connect_go.Response[emptypb.Empty], error)
}

// NewAccountServiceHandler builds an HTTP handler from the service implementation. It returns the
// path on which to mount the handler and the handler itself.
//
// By default, handlers support the Connect, gRPC, and gRPC-Web protocols with the binary Protobuf
// and JSON codecs. They also support gzip compression.
func NewAccountServiceHandler(svc AccountServiceHandler, opts ...connect_go.HandlerOption) (string, http.Handler) {
	mux := http.NewServeMux()
	mux.Handle("/user.account.AccountService/CreateAccount", connect_go.NewUnaryHandler(
		"/user.account.AccountService/CreateAccount",
		svc.CreateAccount,
		opts...,
	))
	mux.Handle("/user.account.AccountService/GetAccountById", connect_go.NewUnaryHandler(
		"/user.account.AccountService/GetAccountById",
		svc.GetAccountById,
		opts...,
	))
	mux.Handle("/user.account.AccountService/GetAccountByEmail", connect_go.NewUnaryHandler(
		"/user.account.AccountService/GetAccountByEmail",
		svc.GetAccountByEmail,
		opts...,
	))
	mux.Handle("/user.account.AccountService/DeleteAccount", connect_go.NewUnaryHandler(
		"/user.account.AccountService/DeleteAccount",
		svc.DeleteAccount,
		opts...,
	))
	return "/user.account.AccountService/", mux
}

// UnimplementedAccountServiceHandler returns CodeUnimplemented from all methods.
type UnimplementedAccountServiceHandler struct{}

func (UnimplementedAccountServiceHandler) CreateAccount(context.Context, *connect_go.Request[account.CreateAccountRequest]) (*connect_go.Response[account.CreateAccountResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("user.account.AccountService.CreateAccount is not implemented"))
}

func (UnimplementedAccountServiceHandler) GetAccountById(context.Context, *connect_go.Request[account.GetAccountByIdRequest]) (*connect_go.Response[account.GetAccountByIdResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("user.account.AccountService.GetAccountById is not implemented"))
}

func (UnimplementedAccountServiceHandler) GetAccountByEmail(context.Context, *connect_go.Request[account.GetAccountByEmailRequest]) (*connect_go.Response[account.GetAccountByEmailResponse], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("user.account.AccountService.GetAccountByEmail is not implemented"))
}

func (UnimplementedAccountServiceHandler) DeleteAccount(context.Context, *connect_go.Request[account.DeleteAccountRequest]) (*connect_go.Response[emptypb.Empty], error) {
	return nil, connect_go.NewError(connect_go.CodeUnimplemented, errors.New("user.account.AccountService.DeleteAccount is not implemented"))
}