package middleware

import (
	"context"
	"net/http"
	"strings"

	"github.com/teragrammer/payment-gateway-wrapper/internal/services"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func JWT() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			auth := r.Header.Get("Authorization")
			if !strings.HasPrefix(auth, "Bearer ") {
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "AUTH_JWT", "Bearer is not set correctly")
				return
			}

			token := strings.TrimPrefix(auth, "Bearer ")
			claims, err := services.ValidateAccessToken(token)
			if err != nil {
				utils.JSONErrorMessage(w, http.StatusUnauthorized, "AUTH_JWT", "Unauthorized")
				return
			}

			ctx := context.WithValue(r.Context(), services.UserContextKey, claims)
			next.ServeHTTP(w, r.WithContext(ctx))
		})
	}
}
