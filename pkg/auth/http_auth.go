package auth

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"strings"

	myContext "github.com/loak155/techbranch-backend/pkg/context"
	"github.com/loak155/techbranch-backend/pkg/jwt"
	"github.com/loak155/techbranch-backend/pkg/redis"
)

type AuthRequest struct {
	Mehtod string
	URL    *regexp.Regexp
	Auth   bool
}

type AuthHandler struct {
	jwtManager   jwt.JwtManager
	redisManager redis.RedisManager
	authRequests []AuthRequest
}

func NewAuthHandler(jwtManager jwt.JwtManager, redisManager redis.RedisManager, AuthRequests []AuthRequest) *AuthHandler {
	return &AuthHandler{jwtManager, redisManager, AuthRequests}
}

func (ah *AuthHandler) HttpAuth(handler http.Handler) http.Handler {
	return http.HandlerFunc(func(res http.ResponseWriter, req *http.Request) {
		for _, authRequest := range AuthRequests {
			if authRequest.Mehtod == req.Method && authRequest.URL.MatchString(req.URL.Path) {
				if authRequest.Auth {
					req, err := ah.AuthFunc(req)
					if err != nil {
						http.Error(res, "invalid token", http.StatusUnauthorized)
						return
					}
					handler.ServeHTTP(res, req)
					return
				}
				handler.ServeHTTP(res, req)
				return
			}
		}
		http.Error(res, "invalid url", http.StatusUnauthorized)
	})
}

func (ah *AuthHandler) AuthFunc(req *http.Request) (*http.Request, error) {
	authorizationHeader := req.Header.Get("Authorization")
	if authorizationHeader == "" {
		return nil, fmt.Errorf("authorization token is required")
	}

	ary := strings.Split(authorizationHeader, " ")
	if len(ary) != 2 || ary[0] != "Bearer" {
		return nil, fmt.Errorf("authorization token is required")
	}

	claims, err := ah.jwtManager.ValidateToken(ary[1])
	if err != nil {
		return nil, err
	}

	userID, err := strconv.Atoi(claims.Subject)
	if err != nil {
		return nil, err
	}

	jti, err := ah.redisManager.Get(req.Context(), claims.Subject)
	if err != nil {
		return nil, err
	}

	if jti != claims.Id {
		return nil, fmt.Errorf("jti is not valid")
	}

	newCtx := myContext.SetUserID(req.Context(), userID)
	return req.WithContext(newCtx), nil
}
