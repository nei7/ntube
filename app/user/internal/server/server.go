package server

import (
	"context"

	"github.com/go-kratos/kratos/v2/middleware"
	"github.com/go-kratos/kratos/v2/middleware/auth/jwt"
	"github.com/go-kratos/kratos/v2/middleware/selector"
	jwtv4 "github.com/golang-jwt/jwt/v4"
	"github.com/google/wire"
	v1 "github.com/nei7/ntube/api/user/v1"
)

var ProviderSet = wire.NewSet(NewGRPCServer, NewHTTPServer)

func selectorFn(secret string) middleware.Middleware {
	return selector.Server(jwt.Server(func(t *jwtv4.Token) (interface{}, error) {
		return []byte(secret), nil
	})).Match(func(ctx context.Context, operation string) bool {
		whiteList := map[string]bool{
			v1.OperationUserServiceVerifyPassword: false,
			v1.OperationUserServiceCreateUser:     false,
		}
		if _, ok := whiteList[operation]; ok {
			return false
		}
		return true
	}).Build()
}
