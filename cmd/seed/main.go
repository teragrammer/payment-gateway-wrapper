package main

import (
	"context"
	"log"

	"github.com/teragrammer/payment-gateway-wrapper/internal/config"
	mongo2 "github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo/seeder"
	"go.mongodb.org/mongo-driver/mongo"
)

func main() {
	cfg := config.Load()
	client, ctx, cancel, err := mongo2.ConnectMongo(cfg.MongoURI)
	if err != nil {
		log.Fatal(err)
	}
	defer cancel()
	defer func(client *mongo.Client, ctx context.Context) {
		err := client.Disconnect(ctx)
		if err != nil {
			log.Fatal(err)
		}
	}(client, ctx)

	db := client.Database(cfg.MongoDBName)

	for _, seed := range seeder.All {
		if err := seed(db); err != nil {
			log.Fatal("Seeding failed:", err)
		}
	}

	log.Println("Seeding completed")
}
