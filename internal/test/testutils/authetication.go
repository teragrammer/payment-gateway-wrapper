package testutils

import (
	"github.com/teragrammer/payment-gateway-wrapper/internal/database/mongo"
	"github.com/teragrammer/payment-gateway-wrapper/internal/repository"
	"github.com/teragrammer/payment-gateway-wrapper/internal/services"
)

func MockJWT(username string) string {
	db, _ := mongo.DefaultMongo()
	user, _ := repository.NewUserRepository(db).GetByUsername(username)
	token, _ := services.GenerateAccessToken(user.ID.Hex(), user.Role)
	return token
}
