package auth

import (
	"context"
	"fmt"
	"strconv"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/redis"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type AuthInterceptor struct {
	jwtManager   jwt.JwtManager
	redisManager redis.RedisManager
	authMethods  map[string]bool
}

func NewAuthInterceptor(jwtManager jwt.JwtManager, redisManager redis.RedisManager, authMethods map[string]bool) *AuthInterceptor {
	return &AuthInterceptor{jwtManager, redisManager, authMethods}
}

func (ai *AuthInterceptor) Unary() grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
		if val, ok := ai.authMethods[info.FullMethod]; !ok {
			return nil, status.Errorf(codes.Unauthenticated, "invalid url")
		} else if val {
			newCtx, err := ai.AuthFunc(ctx)
			if err != nil {
				return nil, status.Errorf(codes.Unauthenticated, "invalid token")
			}
			return handler(newCtx, req)
		}
		return handler(ctx, req)
	}
}

func (ai *AuthInterceptor) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "Bearer")
	if err != nil {
		return nil, err
	}

	claims, err := ai.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, err
	}

	jti, err := ai.redisManager.Get(ctx, claims.Subject)
	if err != nil {
		return nil, err
	}

	if jti != claims.Id {
		return nil, fmt.Errorf("jti is not valid")
	}

	newCtx := myContext.SetUserID(ctx, userID)
	return newCtx, nil
}
