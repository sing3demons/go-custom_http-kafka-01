package db

import (
	"context"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

type MongoDB struct {
	*mongo.Collection
}

func NewMongoDB() *mongo.Database {
	dsn := os.Getenv("MONGO_URI")
	if dsn == "" {
		dsn = "mongodb://localhost:27017/my_app?authSource=admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		panic("failed to connect database")
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		panic("failed to ping database")
	}
	// collection := client.Database("my_app").Collection("product")
	return client.Database("my_app")
}
