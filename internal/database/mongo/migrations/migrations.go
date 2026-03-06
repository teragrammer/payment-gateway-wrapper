package migrations

import "go.mongodb.org/mongo-driver/mongo"

type Migration func(db *mongo.Database) error

var All = []Migration{
	UsersCollection,
	ProjectsCollection,
	XenditSettingsCollection,
	XenditSessionsCollection,
}
