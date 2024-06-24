package middlewares

import (
	"Book_market_api/response"
	"Book_market_api/utils"
	"context"
	"net/http"
	"strings"
)

type ContextKey string

const (
	ContextUserID_v1 ContextKey = "userId"
	ContextUserID_v2 ContextKey = "userId"
)

func getTokenFromHeader(r *http.Request) string {
	token := r.Header.Get("Authorization")
	if token != "" && strings.HasPrefix(token, "Bearer ") {
		return strings.TrimPrefix(token, "Bearer ")
	}
	return ""
}

func Authenticate_v1(roles ...string) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token := getTokenFromHeader(r)
			if token == "" {
				response.ErrorResponse(w, 305)
				return
			}
			userId, err := utils.VerifyToken_v1(token, roles)
			if err != nil {
				response.ErrorResponse(w, 305)
				return
			}
			ctx := context.WithValue(r.Context(), ContextUserID_v1, userId)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}

func Authenticate_v2(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		token := getTokenFromHeader(r)
		if token == "" {
			response.ErrorResponse(w, 305)
			return
		}
		userId, err := utils.VerifyToken_v2(token)
		if err != nil {
			response.ErrorResponse(w, 305)
			return
		}
		ctx := context.WithValue(r.Context(), ContextUserID_v2, userId)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}
