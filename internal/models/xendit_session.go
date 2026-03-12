package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Status int

const (
	Pending Status = iota
	Processing
	Success
	Failed
	Expired
)

type XenditSession struct {
	ID          primitive.ObjectID `json:"_id" bson:"_id"`
	ProjectID   primitive.ObjectID `json:"project_id" bson:"project_id"`
	ReferenceID string             `json:"reference_id" bson:"reference_id"`
	SessionType string             `json:"session_type" bson:"session_type"`
	Mode        string             `json:"mode" bson:"mode"`
	Amount      float64            `json:"amount" bson:"amount"`
	Currency    string             `json:"currency" bson:"currency"`
	Country     string             `json:"country" bson:"country"`
	Customer    Customer           `json:"customer" bson:"customer"`
	Status      Status             `json:"status" bson:"status"`
	UpdatedAt   time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt   time.Time          `json:"created_at" bson:"created_at"`
}

type Customer struct {
	ReferenceID      string           `json:"reference_id" bson:"reference_id"`
	Type             string           `json:"type" bson:"type"`
	Email            string           `json:"email" bson:"email"`
	MobileNumber     string           `json:"mobile_number" bson:"mobile_number"`
	IndividualDetail IndividualDetail `json:"individual_detail" bson:"individual_detail"`
}

type IndividualDetail struct {
	GivenName string `json:"given_name" bson:"given_name"`
	Surname   string `json:"surname" bson:"surname"`
}
