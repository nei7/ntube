package interceptor

import (
	"context"

	"github.com/nei7/gls/internal/service"
	"google.golang.org/grpc"
)

type AuthInterceptor struct {
	tokenManager service.TokenManager
}

func NewAuthInterceptor(tokenManager service.TokenManager) *AuthInterceptor {
	return &AuthInterceptor{
		tokenManager,
	}
}

func (i *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp interface{}, err error) {

		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) Stream() grpc.StreamServerInterceptor {}
