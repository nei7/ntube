// Code generated by protoc-gen-go-http. DO NOT EDIT.
// versions:
// - protoc-gen-go-http v2.7.0
// - protoc             v4.23.4
// source: auth/v1/auth.proto

package v1

import (
	context "context"
	http "github.com/go-kratos/kratos/v2/transport/http"
	binding "github.com/go-kratos/kratos/v2/transport/http/binding"
)

// This is a compile-time assertion to ensure that this generated file
// is compatible with the kratos package it is being compiled against.
var _ = new(context.Context)
var _ = binding.EncodeURL

const _ = http.SupportPackageIsVersion1

const OperationAuthSercieVerifyEmail = "/v1.AuthSercie/VerifyEmail"

type AuthSercieHTTPServer interface {
	VerifyEmail(context.Context, *VerifyEmailRequest) (*VerifyEmailResponse, error)
}

func RegisterAuthSercieHTTPServer(s *http.Server, srv AuthSercieHTTPServer) {
	r := s.Route("/")
	r.GET("v1/email/verify", _AuthSercie_VerifyEmail0_HTTP_Handler(srv))
}

func _AuthSercie_VerifyEmail0_HTTP_Handler(srv AuthSercieHTTPServer) func(ctx http.Context) error {
	return func(ctx http.Context) error {
		var in VerifyEmailRequest
		if err := ctx.BindQuery(&in); err != nil {
			return err
		}
		http.SetOperation(ctx, OperationAuthSercieVerifyEmail)
		h := ctx.Middleware(func(ctx context.Context, req interface{}) (interface{}, error) {
			return srv.VerifyEmail(ctx, req.(*VerifyEmailRequest))
		})
		out, err := h(ctx, &in)
		if err != nil {
			return err
		}
		reply := out.(*VerifyEmailResponse)
		return ctx.Result(200, reply)
	}
}

type AuthSercieHTTPClient interface {
	VerifyEmail(ctx context.Context, req *VerifyEmailRequest, opts ...http.CallOption) (rsp *VerifyEmailResponse, err error)
}

type AuthSercieHTTPClientImpl struct {
	cc *http.Client
}

func NewAuthSercieHTTPClient(client *http.Client) AuthSercieHTTPClient {
	return &AuthSercieHTTPClientImpl{client}
}

func (c *AuthSercieHTTPClientImpl) VerifyEmail(ctx context.Context, in *VerifyEmailRequest, opts ...http.CallOption) (*VerifyEmailResponse, error) {
	var out VerifyEmailResponse
	pattern := "v1/email/verify"
	path := binding.EncodeURL(pattern, in, true)
	opts = append(opts, http.Operation(OperationAuthSercieVerifyEmail))
	opts = append(opts, http.PathTemplate(pattern))
	err := c.cc.Invoke(ctx, "GET", path, nil, &out, opts...)
	if err != nil {
		return nil, err
	}
	return &out, err
}
