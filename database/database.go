package database

import (
	"context"

	"github.com/albertopformoso/graphql-todo-list/config"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func ConnectDB(ctx context.Context, config config.MongoConfiguration) *mongo.Database {
	connection := options.Client().ApplyURI(config.Server)
	client, err := mongo.Connect(ctx, connection)
	if err != nil {
		panic(err)
	}

	return client.Database(config.Database)
}
