package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Project struct {
	ID         *primitive.ObjectID `json:"_id" bson:"_id"`
	UserID     primitive.ObjectID  `json:"user_id" bson:"user_id"`
	Name       string              `json:"name" bson:"name"`
	PublicKey  string              `json:"public_key" bson:"public_key"`
	PrivateKey string              `json:"private_key" bson:"private_key"`
	IsActive   int                 `json:"is_active" bson:"is_active"`
	UpdatedAt  *time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt  time.Time           `json:"created_at" bson:"created_at"`
}

func (p *Project) ToBson() bson.M {
	bn := bson.M{
		"name":        p.Name,
		"user_id":     p.UserID,
		"public_key":  p.PublicKey,
		"private_key": p.PrivateKey,
		"is_active":   p.IsActive,
		"created_at":  p.CreatedAt,
	}
	if p.ID != nil {
		bn["_id"] = p.ID
	}
	if p.UpdatedAt != nil {
		bn["updated_at"] = p.UpdatedAt
	}
	return bn
}
