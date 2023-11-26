package services

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type EventCreateProductRequest struct {
	Header map[string]any       `json:"header"`
	Body   CreateProductRequest `json:"body"`
}

type DeleteProductRequest struct {
	ID         string    `json:"id" bson:"id"`
	DeleteDate time.Time `json:"delete_date" bson:"deleteDate"`
}
type Product struct {
	MID          primitive.ObjectID         `json:"_id" bson:"_id"`
	Type         string                     `json:"@type" bson:"@type"`
	Status       string                     `json:"status" bson:"status"`
	Category     []Category                 `json:"category,omitempty" bson:"category,omitempty"`
	Href         string                     `json:"href"`
	ID           string                     `json:"id" bson:"id"`
	Title        string                     `json:"title,omitempty" bson:"title,omitempty"`
	Description  string                     `json:"description,omitempty" bson:"description,omitempty"`
	Image        string                     `json:"image,omitempty" bson:"image,omitempty"`
	ProductPrice []CreateUpdateProductPrice `json:"productPrice,omitempty" bson:"productPrice,omitempty" form:"productPrice,omitempty"`
}

type Category struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Type string `json:"@type" bson:"@type"`
}

type EventDeleteProductRequest struct {
	Header map[string]any       `json:"header"`
	Body   DeleteProductRequest `json:"body"`
}

type CreateProductRequest struct {
	ID           string                     `json:"id" bson:"id" form:"id"`
	Type         string                     `json:"@type" bson:"@type"`
	Category     []Category                 `json:"category,omitempty" bson:"category,omitempty"`
	Status       string                     `json:"status" bson:"status"`
	Title        string                     `json:"title,omitempty" bson:"title,omitempty" form:"title" binding:"required"`
	Description  string                     `json:"description,omitempty" bson:"description,omitempty" form:"description,omitempty"`
	Image        string                     `json:"image,omitempty" bson:"image,omitempty" form:"image,omitempty"`
	ProductPrice []CreateUpdateProductPrice `json:"productPrice,omitempty" bson:"productPrice,omitempty" form:"productPrice,omitempty"`
	LastUpdate   time.Time                  `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}
