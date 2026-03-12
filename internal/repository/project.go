package repository

import (
	"context"
	"log"

	"github.com/teragrammer/payment-gateway-wrapper/internal/models"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

var ProjectCollection = "projects"

type ProjectRepository struct {
	collection *mongo.Collection
}

func NewProjectRepository(db *mongo.Database) *ProjectRepository {
	collection := db.Collection(ProjectCollection)
	return &ProjectRepository{collection: collection}
}

func (r *ProjectRepository) Create(inputs bson.D) (interface{}, error) {
	result, err := r.collection.InsertOne(context.Background(), inputs)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *ProjectRepository) Browse(search string, page int, pageSize int) ([]map[string]interface{}, error) {
	filter := bson.M{}

	if search != "" {
		filter["name"] = bson.M{
			"$regex":   search,
			"$options": "i",
		}
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	results, err := Cursor(r.collection, filter, findOptions, nil, page, pageSize)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *ProjectRepository) GetByID(id string) (*models.Project, error) {
	var project models.Project
	err := r.collection.FindOne(context.Background(), bson.M{"_id": utils.StringToObjectID(id)}).Decode(&project)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &project, nil
}

func (r *ProjectRepository) Update(id string, inputs bson.D) error {
	filter := bson.M{"_id": utils.StringToObjectID(id)}

	_, err := r.collection.UpdateOne(context.Background(), filter, bson.D{{"$set", inputs}})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (r *ProjectRepository) Delete(id string) (int64, error) {
	filter := bson.M{"_id": utils.StringToObjectID(id)}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	log.Println("Number of deleted projects: ", result.DeletedCount)
	return result.DeletedCount, nil
}
