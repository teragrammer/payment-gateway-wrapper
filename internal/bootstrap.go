package internal

import (
	"log"

	"github.com/redis/go-redis/v9"
	mongo2 "github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	redis2 "github.com/teragrammer/payment-gateway-wrapper/internal/database/redis"
	"go.mongodb.org/mongo-driver/mongo"
)

type Bootstrap struct {
	DB    *mongo.Database
	Redis redis.UniversalClient
}

func InitializedBootstrap() Bootstrap {
	db, err := mongo2.DefaultMongo()
	if err != nil {
		log.Fatalf("Error connecting to Mongo: %v", err.Error())
	}

	return Bootstrap{DB: db, Redis: redis2.ConnectRedis()}
}
