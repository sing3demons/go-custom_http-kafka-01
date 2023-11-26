package repository

import (
	"context"
	"time"

	"github.com/sing3demons/go-category-service/model"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
)

type category struct {
	*mongo.Database
	logger *logrus.Logger
}

func NewCategory(db *mongo.Database, logger *logrus.Logger) CategoryRepository {
	return &category{db, logger}
}

type CategoryRepository interface {
	Save(doc model.CreateCategoryReq) error
	Update(req model.UpdateCategoryReq) (category *model.Category, err error)
}

func (tx *category) Save(doc model.CreateCategoryReq) error {
	dbName := "category"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	r, err := tx.Database.Collection(dbName).InsertOne(ctx, doc)
	if err != nil {
		tx.logger.WithFields(logrus.Fields{
			"dbName": dbName,
			"data":   doc,
			"error":  err,
		}).Error("insert category error")
		return err
	}

	tx.logger.WithFields(logrus.Fields{
		"dbName":   dbName,
		"data":     doc,
		"resultID": r.InsertedID,
	}).Debug("insert category success")
	return nil
}

func (tx *category) Update(req model.UpdateCategoryReq) (category *model.Category, err error) {
	dbName := "category"

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	filter := bson.M{"id": req.ID, "deleteDate": nil}
	var update model.UpdateCategoryReq
	update.ID = req.ID
	if req.Name != "" {
		update.Name = req.Name
	}
	if req.Status != "" {
		update.Status = req.Status
	}
	if req.LastUpdate.IsZero() {
		update.LastUpdate = time.Now().UTC()
	}

	if len(req.Products) > 0 {
		for _, v := range req.Products {
			product := model.AddProduct{
				ID:   v.ID,
				Type: "products",
			}
			if v.Name != "" {
				product.Name = v.Name
			}
			update.Products = append(update.Products, product)
		}
	}

	// var category category
	err = tx.Database.Collection(dbName).FindOneAndUpdate(ctx, filter, update).Decode(&category)
	if err != nil {
		tx.logger.WithFields(logrus.Fields{
			"dbName": dbName,
			"data":   update,
			"error":  err,
		}).Error("insert category error")
		return nil, err
	}

	tx.logger.WithFields(logrus.Fields{
		"dbName": dbName,
		"update": update,
		"data":   category,
	}).Debug("insert category success")
	return category, nil
}
