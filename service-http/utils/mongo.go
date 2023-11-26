package utils

import (
	"context"
	"time"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func GetMultiWithTotal[T any](collection *mongo.Collection, filter bson.M, findOptions *options.FindOptions) ([]T, int64, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter = bson.M{
		"$and": []bson.M{
			filter,
			{"deleteDate": nil},
		},
	}

	var total int64
	countCh := make(chan int64)
	go func() {
		count, err := collection.CountDocuments(ctx, filter)
		if err != nil {
			countCh <- 0
		}
		countCh <- count
	}()
	total = <-countCh

	var result []T
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	for cur.Next(ctx) {
		var p T
		if err := cur.Decode(&p); err != nil {
			return nil, 0, err
		}
		result = append(result, p)
	}

	return result, total, nil
}

func GetMulti[T any](collection *mongo.Collection, filter bson.M, findOptions *options.FindOptions) ([]T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter = bson.M{
		"$and": []bson.M{
			filter,
			{"deleteDate": nil},
		},
	}

	var result []T
	cur, err := collection.Find(ctx, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for cur.Next(ctx) {
		var p T
		if err := cur.Decode(&p); err != nil {
			return nil, err
		}
		result = append(result, p)
	}
	return result, nil
}

func GetOne[T any](collection *mongo.Collection, filter bson.M, findOptions *options.FindOneOptions) (*T, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result T
	err := collection.FindOne(ctx, filter, findOptions).Decode(&result)
	if err != nil {
		return nil, err
	}

	return &result, nil
}
