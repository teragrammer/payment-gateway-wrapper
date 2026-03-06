package migrations

import (
	"context"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func XenditSessionsCollection(db *mongo.Database) error {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	collectionName := "xendit_sessions"

	err := utils.CreateCollection(collectionName, db, ctx)
	if err != nil {
		return err
	}

	indexes := []mongo.IndexModel{
		{
			Keys: bson.D{{Key: "project_id", Value: 1}},
			Options: options.Index().
				SetName("project_id_idx"),
		},
		{
			Keys: bson.D{{Key: "reference_id", Value: -1}},
			Options: options.Index().
				SetName("reference_id_idx"),
		},
		{
			Keys: bson.D{{Key: "customer.reference_id", Value: 1}},
			Options: options.Index().
				SetName("customer_reference_id_idx"),
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
