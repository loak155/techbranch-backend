package interceptor

import (
	"context"
	"log/slog"

	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/auth"
	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
)

type AuthInterceptor struct {
	jwtManager jwt.JwtManager
}

func NewAuthInterceptor(jwtManager jwt.JwtManager) *AuthInterceptor {
	return &AuthInterceptor{jwtManager}
}

func (ai *AuthInterceptor) AuthFunc(ctx context.Context) (context.Context, error) {
	token, err := auth.AuthFromMD(ctx, "bearer")
	slog.Info("[Message]", "Token", token, "Error", err)
	if err != nil {
		return nil, err
	}

	claims, err := ai.jwtManager.ValidateToken(token)
	if err != nil {
		return nil, err
	}

	newCtx := myContext.SetUserID(ctx, claims.UserId)
	return newCtx, nil
}
