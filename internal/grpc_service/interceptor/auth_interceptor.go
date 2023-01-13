package interceptor

import (
	"context"

	"github.com/nei7/ntube/internal/service"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"

	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
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
		md, err := i.authorize(ctx)
		if err != nil {
			return nil, err
		}

		ctx = metadata.NewIncomingContext(ctx, md)

		return handler(ctx, req)
	}
}

func (i *AuthInterceptor) authorize(ctx context.Context) (metadata.MD, error) {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, status.Errorf(codes.Unauthenticated, "metadata is not provided")
	}

	values := md["authorization"]
	if len(values) == 0 {
		return nil, status.Errorf(codes.Unauthenticated, "authorization token is not provided")
	}

	accessToken := values[0]
	userid, err := i.tokenManager.Parse(accessToken)
	if err != nil {
		return nil, status.Errorf(codes.Unauthenticated, "invalid access token")
	}

	md.Append("id", userid)

	return md, nil
}
