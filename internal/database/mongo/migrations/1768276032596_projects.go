package migrations

import (
	"context"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ProjectsCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionName := "projects"

	err := utils.CreateCollection(collectionName, db, ctx)
	if err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "user_id", Value: 1}},
			Options: options.Index().
				SetName("user_id_idx"),
		},
		{
			Keys: bson.D{{Key: "private_key", Value: -1}},
			Options: options.Index().
				SetUnique(true).
				SetName("uniq_private_key"),
		},
		{
			Keys: bson.D{{Key: "created_at", Value: -1}},
			Options: options.Index().
				SetName("created_at_idx"),
		},
	}

	_, err = db.Collection(collectionName).Indexes().CreateMany(ctx, indexes)
	return err
}
