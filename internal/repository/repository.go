package repository

import (
	"context"
	"log"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func Cursor(collection *mongo.Collection, filter bson.M, findOptions *options.FindOptions, projection *bson.M, page int, pageSize int) ([]map[string]interface{}, error) {
	skip := (page - 1) * pageSize
	finder := findOptions.SetSkip(int64(skip)).SetLimit(int64(pageSize))

	if projection != nil {
		finder.SetProjection(*projection)
	}

	cursor, err := collection.Find(context.Background(), filter, finder)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	var results []map[string]interface{}
	for cursor.Next(context.Background()) {
		var result map[string]interface{}
		if err := cursor.Decode(&result); err != nil {
			log.Println(err.Error())
			return nil, err
		}
		results = append(results, result)
	}

	if err := cursor.Err(); err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return results, nil
}
