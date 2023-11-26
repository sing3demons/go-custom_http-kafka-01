package main

import (
	"context"
	"log"
	"os"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.mongodb.org/mongo-driver/mongo/readpref"
)

func ConnectMonoDB() (*mongo.Database, error) {
	dsn := os.Getenv("MONGO_URI")
	if dsn == "" {
		dsn = "mongodb://localhost:27017/my_app?authSource=admin"
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	client, err := mongo.Connect(ctx, options.Client().ApplyURI(dsn))
	if err != nil {
		return nil, err
	}

	if err := client.Ping(ctx, readpref.Primary()); err != nil {
		return nil, err
	}

	indexModel := mongo.IndexModel{
		Keys: bson.D{{Key: "id", Value: 1}},
	}
	client.Database("my_app").Collection("product").Indexes().CreateOne(context.TODO(), indexModel)
	client.Database("my_app").Collection("productPrice").Indexes().CreateOne(context.TODO(), indexModel)

	return client.Database("my_app"), nil

	// return collection, nil
}

func DisconnectMongo(client *mongo.Client) {
	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()

	err := client.Disconnect(ctx)
	if err != nil {
		log.Fatal(err)
	}
}
