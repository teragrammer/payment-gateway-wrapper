package main

import (
	"github.com/teragrammer/payment-gateway-wrapper/internal/app"
	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	"github.com/teragrammer/payment-gateway-wrapper/internal/routes"
)

func main() {
	cfg := config.Load()

	handler := routes.Register()
	app.Run(handler, cfg.Port)
}
