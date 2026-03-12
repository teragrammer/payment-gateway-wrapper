package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type XenditSetting struct {
	ID         primitive.ObjectID `json:"_id" bson:"_id"`
	ProjectID  primitive.ObjectID `json:"project_id" bson:"project_id"`
	Name       string             `json:"name" bson:"name"`
	Slug       string             `json:"slug" bson:"slug"`
	Country    string             `json:"country" bson:"country"`
	Currency   string             `json:"currency" bson:"currency"`
	SuccessURL string             `json:"success_url" bson:"success_url"`
	FailedURL  string             `json:"failed_url" bson:"failed_url"`
	PublicKey  string             `json:"public_key" bson:"public_key"`
	PrivateKey string             `json:"private_key" bson:"private_key"`
	IsActive   int                `json:"is_active" bson:"is_active"`
	UpdatedAt  time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt  time.Time          `json:"created_at" bson:"created_at"`
}
