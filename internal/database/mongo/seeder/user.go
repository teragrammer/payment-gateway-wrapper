package seeder

import (
	"context"
	"time"

	"github.com/teragrammer/payment-gateway-wrapper/internal/models"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/utils"
	"go.mongodb.org/mongo-driver/mongo"
)

func UserSeeder(db *mongo.Database) error {
	password, err := utils.GeneratePasswordHash("12345678")
	if err != nil {
		return err
	}

	users := []models.User{
		{FullName: "Administrator", Username: "admin", Password: password, Role: models.Admin, CreatedAt: time.Now()},
	}

	for _, user := range users {
		_, err := db.Collection(repository.UserCollection).InsertOne(context.Background(), user.ToBson())
		if err != nil {
			return err
		}
	}

	return nil
}
