package routes

import (
	"net/http"

	"github.com/go-chi/chi/v5"
	"github.com/teragrammer/payment-gateway-wrapper/internal"
	handlers2 "github.com/teragrammer/payment-gateway-wrapper/internal/core/handlers"
	middleware2 "github.com/teragrammer/payment-gateway-wrapper/internal/core/middleware"
	"github.com/teragrammer/payment-gateway-wrapper/internal/core/validations"
)

func Register() http.Handler {
	application := internal.InitializedBootstrap()

	r := chi.NewRouter()
	r.Use(middleware2.Recovery())
	r.Use(middleware2.SecurityHeaders())

	r.With(middleware2.RateLimitIP()).Get("/health", handlers2.Health)
	r.With(middleware2.RateLimitIP()).With(validations.ValidateAuthenticationInputs).Post("/api/v1/login", handlers2.Login)
	r.With(middleware2.RateLimitIP()).With(middleware2.JWT()).Get("/api/v1/logout", handlers2.Logout)

	r.Route("/api/v1/users", handlers2.UserRoutes(application.DB))
	r.Route("/api/v1/projects", handlers2.ProjectRoutes(application))

	r.With(middleware2.RateLimitIP()).Get("/api/v1/supports", handlers2.GetSupportedPayments)

	return r
}
