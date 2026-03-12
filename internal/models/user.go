package models

import (
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Role int

const (
	Admin   Role = 0
	Manager      = 1
)

type User struct {
	ID        *primitive.ObjectID `json:"_id" bson:"_id"`
	FullName  string              `json:"full_name" bson:"full_name"`
	Username  string              `json:"username" bson:"username"`
	Password  string              `json:"password" bson:"password"`
	Role      Role                `json:"role" bson:"role"`
	UpdatedAt *time.Time          `json:"updated_at" bson:"updated_at"`
	CreatedAt time.Time           `json:"created_at" bson:"created_at"`
}

func (u *User) ToBson() bson.M {
	bs := bson.M{
		"full_name":  u.FullName,
		"username":   u.Username,
		"password":   u.Password,
		"role":       u.Role,
		"created_at": u.CreatedAt,
	}
	if u.ID != nil {
		bs["_id"] = u.ID
	}
	if u.UpdatedAt != nil {
		bs["updated_at"] = u.UpdatedAt
	}
	return bs
}
