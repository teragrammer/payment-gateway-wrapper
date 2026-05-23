package utils

import (
	"context"
	"errors"

	"go.mongodb.org/mongo-driver/mongo"
)

func CreateCollection(collectionName string, db *mongo.Database, ctx context.Context) error {
	err := db.CreateCollection(ctx, collectionName)
	if err != nil {
		var cmdErr mongo.CommandError
		if errors.As(err, &cmdErr) && cmdErr.Code == 48 {
			// Collection already exists → OK
			return nil
		}

		return err
	}
	return nil
}
