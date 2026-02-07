package middleware

import (
	"log"
	"net/http"

	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
)

func Recovery() Middleware {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			defer func() {
				if r := recover(); r != nil {
					log.Printf("Recovered from panic: %v", r)
					utils.JSONErrorMessage(w, http.StatusInternalServerError, "CRITICAL", "Whoops something went wrong")
				}
			}()
			next.ServeHTTP(w, r)
		})
	}
}
