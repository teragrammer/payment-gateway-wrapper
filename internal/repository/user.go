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

var UserCollection = "users"

type UserRepository struct {
	collection *mongo.Collection
}

func NewUserRepository(db *mongo.Database) *UserRepository {
	collection := db.Collection(UserCollection)
	return &UserRepository{collection: collection}
}

func (r *UserRepository) GetByID(id string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.M{"_id": utils.StringToObjectID(id)}).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUsername(username string) (*models.User, error) {
	var user models.User
	err := r.collection.FindOne(context.Background(), bson.D{{"username", username}}).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) GetByUniqueUsername(username string, id string) (*models.User, error) {
	filter := bson.M{
		"username": username,
		"_id":      bson.M{"$ne": utils.StringToObjectID(id)},
	}

	var user models.User
	err := r.collection.FindOne(context.Background(), filter).Decode(&user)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}
	return &user, nil
}

func (r *UserRepository) Browse(search string, page int, pageSize int) ([]map[string]interface{}, error) {
	filter := bson.M{}

	if search != "" {
		filter["username"] = bson.M{
			"$regex":   search,
			"$options": "i",
		}
	}

	// exclude sensitive information by default for browsing
	projection := bson.M{
		"password": 0,
	}

	findOptions := options.Find().SetSort(bson.D{{Key: "created_at", Value: 1}})
	results, err := Cursor(r.collection, filter, findOptions, &projection, page, pageSize)
	if err != nil {
		return nil, err
	}

	return results, nil
}

func (r *UserRepository) Create(inputs bson.D) (interface{}, error) {
	result, err := r.collection.InsertOne(context.Background(), inputs)
	if err != nil {
		log.Println(err.Error())
		return nil, err
	}

	return result.InsertedID, nil
}

func (r *UserRepository) Update(id string, inputs bson.D) error {
	filter := bson.M{"_id": utils.StringToObjectID(id)}

	_, err := r.collection.UpdateOne(context.Background(), filter, bson.D{{"$set", inputs}})
	if err != nil {
		log.Println(err.Error())
		return err
	}

	return nil
}

func (r *UserRepository) Delete(id string) (int64, error) {
	filter := bson.M{"_id": utils.StringToObjectID(id)}

	result, err := r.collection.DeleteOne(context.Background(), filter)
	if err != nil {
		log.Println(err.Error())
		return 0, err
	}

	log.Println("Number of deleted users: ", result.DeletedCount)
	return result.DeletedCount, nil
}
