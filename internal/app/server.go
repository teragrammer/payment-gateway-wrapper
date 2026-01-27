package app

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/database/redis"
)

func Run(handler http.Handler, port string) {
	srv := &http.Server{
		Addr:         ":" + port,
		Handler:      handler,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	go func() {
		err := srv.ListenAndServe()
		if err != nil {
			log.Fatal(err)
		}
	}()

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, os.Interrupt)
	<-quit

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	defer mongo.CloseMongo()
	defer redis.CloseRedis()
	err := srv.Shutdown(ctx)
	if err != nil {
		return
	}
}
