package product

import (
	"time"

	"go.mongodb.org/mongo-driver/bson/primitive"
)

type Product struct {
	MID          primitive.ObjectID `json:"_id" bson:"_id"`
	Type         string             `json:"@type" bson:"@type"`
	Status       string             `json:"status" bson:"status"`
	Href         string             `json:"href"`
	ID           string             `json:"id" bson:"id"`
	Title        string             `json:"title,omitempty" bson:"title,omitempty"`
	Description  string             `json:"description,omitempty" bson:"description,omitempty"`
	Image        string             `json:"image,omitempty" bson:"image,omitempty"`
	ProductPrice []ProductPrice     `json:"productPrice,omitempty" bson:"productPrice,omitempty"`
	LastUpdate   time.Time          `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Category     []Category         `json:"category,omitempty" bson:"category,omitempty"`
}

type Category struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Type string `json:"@type" bson:"@type"`
	Href string `json:"href,omitempty"`
}

type CreateProductRequest struct {
	ID           string                     `json:"id" bson:"id" form:"id"`
	Type         string                     `json:"@type" bson:"@type"`
	Status       string                     `json:"status" bson:"status"`
	Title        string                     `json:"title,omitempty" bson:"title,omitempty" form:"title" binding:"required"`
	Description  string                     `json:"description,omitempty" bson:"description,omitempty" form:"description,omitempty"`
	Image        string                     `json:"image,omitempty" bson:"image,omitempty" form:"image,omitempty"`
	ProductPrice []CreateUpdateProductPrice `json:"productPrice,omitempty" bson:"productPrice,omitempty" form:"productPrice,omitempty"`
	LastUpdate   time.Time                  `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Category     []Category                 `json:"category,omitempty" bson:"category,omitempty"`
}

type CreateUpdateProductPrice struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}

type UpdateProductRequest struct {
	Status       string                     `json:"status" bson:"status"`
	Title        string                     `json:"title,omitempty" bson:"title,omitempty" form:"title,omitempty"`
	Description  string                     `json:"description,omitempty" bson:"description,omitempty" form:"description,omitempty"`
	Image        string                     `json:"image,omitempty" bson:"image,omitempty" form:"image,omitempty"`
	ProductPrice []CreateUpdateProductPrice `json:"productPrice,omitempty" bson:"productPrice,omitempty" form:"productPrice,omitempty"`
	LastUpdate   time.Time                  `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	Category     []Category                 `json:"category,omitempty" bson:"category,omitempty"`
}

type DeleteProductRequest struct {
	ID         string    `json:"id" bson:"id"`
	DeleteDate time.Time `json:"delete_date" bson:"deleteDate"`
}

type DeleteProductPriceRequest struct {
	ID         string    `json:"id" bson:"id"`
	DeleteDate time.Time `json:"delete_date" bson:"deleteDate"`
}

type Event struct {
	Header any `json:"header"`
	Body   any `json:"body"`
}
type Price struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}

type ProductPrice struct {
	ID         string    `json:"id,omitempty" bson:"id,omitempty"`
	Type       string    `json:"@type,omitempty" bson:"@type,omitempty"`
	Status     string    `json:"status,omitempty" bson:"status,omitempty"`
	Href       string    `json:"href,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Price      *Price     `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate *time.Time `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}
