package services

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type Service struct {
	db     *mongo.Database
	logger *logrus.Logger
}

func NewService(db *mongo.Database, logger *logrus.Logger) *Service {
	return &Service{db, logger}
}

func (svc *Service) InsertProduct(topic string, msg []byte, timestamp string) {
	var event EventCreateProductRequest
	json.Unmarshal(msg, &event)
	req := event.Body

	if req.LastUpdate.IsZero() {
		req.LastUpdate = time.Now().UTC()
	}

	document := CreateProductRequest{
		ID:           req.ID,
		Type:         "products",
		Status:       req.Status,
		Title:        req.Title,
		Description:  req.Description,
		Image:        req.Image,
		ProductPrice: req.ProductPrice,
		LastUpdate:   req.LastUpdate.UTC(),
	}

	if len(req.Category) > 0 {
		document.Category = req.Category
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	dbName := "product"
	result, err := svc.db.Collection(dbName).InsertOne(ctx, document)
	if err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    result,
			"error":     err,
			"timestamp": timestamp,
		}).Error("insert product error")
		return
	}

	var data CreateProductRequest
	if err := svc.db.Collection(dbName).FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&data); err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    result,
			"error":     err,
			"timestamp": timestamp,
		}).Error("Error decoding data")
		return
	}
	svc.logger.WithFields(logrus.Fields{
		"timestamp": timestamp,
		"result":    data,
		"headers":   event.Header,
	}).Info("Insert Product")
}

func (svc *Service) DeleteProduct(topic string, msg []byte, timestamp string) {
	dbName := "product"
	var event EventDeleteProductRequest
	json.Unmarshal(msg, &event)
	req := event.Body

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result Product
	if err := svc.db.Collection(dbName).FindOneAndUpdate(ctx, bson.M{"id": req.ID}, bson.M{
		"$set": bson.M{
			"deleteDate": req.DeleteDate,
		},
	}).Decode(&result); err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    req,
			"error":     err,
			"timestamp": timestamp,
		}).Error("Error decoding data")
		return
	}

	svc.logger.WithFields(logrus.Fields{
		"result":    result,
		"headers":   event.Header,
		"timestamp": timestamp,
	}).Info("Insert Product")
}

func (svc *Service) InsertProductPrice(topic string, msg []byte, timestamp string) {
	dbName := "productPrice"

	var event EventCreateProductPriceRequest
	if err := json.Unmarshal(msg, &event); err != nil {
		svc.logger.WithFields(logrus.Fields{
			"error":     err,
			"timestamp": timestamp,
		}).Error("Error decoding data")
		return
	}
	req := event.Body

	if req.LastUpdate.IsZero() {
		req.LastUpdate = time.Now().UTC()
	}

	document := CreateProductPrice{
		ID:         req.ID,
		Type:       "productPrice",
		Status:     req.Status,
		Name:       req.Name,
		Price:      req.Price,
		LastUpdate: req.LastUpdate.UTC(),
	}

	svc.logger.WithFields(logrus.Fields{
		"timestamp": timestamp,
		"body":      document,
		"headers":   event.Header,
	}).Debug("")

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := svc.db.Collection(dbName).InsertOne(ctx, document)
	if err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    result,
			"error":     err,
			"timestamp": timestamp,
		}).Error("insert product price error")
		return
	}

	var data ProductPrice
	if err := svc.db.Collection(dbName).FindOne(ctx, bson.M{"_id": result.InsertedID}).Decode(&data); err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    result,
			"error":     err,
			"timestamp": timestamp,
		}).Error("Error decoding data")
		return
	}
	svc.logger.WithFields(logrus.Fields{
		"result":    data,
		"headers":   event.Header,
		"timestamp": timestamp,
	}).Info("Insert Product Price")
}
func (svc *Service) DeleteProductPrice(topic string, msg []byte, timestamp string) {
	fmt.Println("DeleteProductPrice")
	dbName := "productPrice"
	var event EventDeleProductPriceRequest
	json.Unmarshal(msg, &event)
	req := event.Body

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	var result ProductPrice
	if err := svc.db.Collection(dbName).FindOneAndUpdate(ctx, bson.M{"id": req.ID}, bson.M{
		"$set": bson.M{
			"deleteDate": req.DeleteDate,
		},
	}).Decode(&result); err != nil {
		svc.logger.WithFields(logrus.Fields{
			"result":    req,
			"error":     err,
			"timestamp": timestamp,
		}).Error("Error decoding data")
		return
	}

	svc.logger.WithFields(logrus.Fields{
		"result":    result,
		"headers":   event.Header,
		"timestamp": timestamp,
	}).Info("Insert Product")
}
