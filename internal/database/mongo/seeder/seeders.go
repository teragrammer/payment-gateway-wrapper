package seeder

import "go.mongodb.org/mongo-driver/mongo"

type Seed func(db *mongo.Database) error

var All = []Seed{
	UserSeeder,
}
