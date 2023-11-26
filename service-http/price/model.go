package price

import "time"

type Event struct {
	Header any `json:"header"`
	Body   any `json:"body"`
}

type CreateProductPrice struct {
	ID         string    `json:"id,omitempty" bson:"id,omitempty"`
	Name       string    `json:"name,omitempty" bson:"name,omitempty"`
	Status     string    `json:"status,omitempty" bson:"status,omitempty"`
	Price      Price     `json:"price,omitempty" bson:"price,omitempty"`
	LastUpdate time.Time `json:"lastUpdate" bson:"lastUpdate"`
}

type Price struct {
	Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
	Value float64 `json:"value,omitempty" bson:"value,omitempty"`
}

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