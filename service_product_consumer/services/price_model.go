package services

import "time"

type ProductPrice struct {
	ID         string    `json:"id" bson:"id"`
	Type       string    `json:"@type" bson:"@type"`
	Status     string    `json:"status" bson:"status"`
	Href       string    `json:"href"`
	Name       string    `json:"name" bson:"name"`
	Price      Price     `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate time.Time `json:"lastUpdate" bson:"lastUpdate"`
}

type DeleteProductPriceRequest struct {
	ID         string    `json:"id" bson:"id"`
	DeleteDate time.Time `json:"delete_date" bson:"deleteDate"`
}

type CreateProductPrice struct {
	ID         string    `json:"id" bson:"id"`
	Type       string    `json:"@type" bson:"@type"`
	Status     string    `json:"status" bson:"status"`
	Name       string    `json:"name" bson:"name"`
	Price      Price     `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate time.Time `json:"lastUpdate" bson:"lastUpdate"`
}

type EventCreateProductPriceRequest struct {
	Header map[string]any     `json:"header"`
	Body   CreateProductPrice `json:"body"`
}

type EventDeleProductPriceRequest struct {
	Header map[string]any            `json:"header"`
	Body   DeleteProductPriceRequest `json:"body"`
}

type Price struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}

type CreateUpdateProductPrice struct {
	ID   string `json:"id" bson:"id"`
	Name string `json:"name" bson:"name"`
}
