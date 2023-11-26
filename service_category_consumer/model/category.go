package model

import "time"

type CreateCategoryReq struct {
	ID         string    `json:"id" bson:"id"`
	Name       string    `json:"name" bson:"name"`
	Type       string    `json:"@type" bson:"@type"`
	Status     string    `json:"status" bson:"status"`
	LastUpdate time.Time `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}

type UpdateCategoryReq struct {
	ID         string       `json:"id" bson:"id"`
	Products   []AddProduct `json:"products" bson:"products"`
	Name       string       `json:"name" bson:"name"`
	Type       string       `json:"@type" bson:"@type"`
	Status     string       `json:"status" bson:"status"`
	LastUpdate time.Time    `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}

type AddProduct struct {
	ID   string `json:"id,omitempty" bson:"id,omitempty"`
	Name string `json:"name,omitempty" bson:"name,omitempty"`
	Type string `json:"@type" bson:"@type"`
}

type Category struct {
	ID         string    `json:"id" bson:"id"`
	Name       string    `json:"name" bson:"name"`
	Products   []Product `json:"products" bson:"products"`
	Type       string    `json:"@type" bson:"@type"`
	Status     string    `json:"status" bson:"status"`
	Href       string    `json:"href"`
	LastUpdate time.Time `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
}

type (
	Product struct {
		ID           string         `json:"id" bson:"id"`
		Type         string         `json:"@type" bson:"@type"`
		Status       string         `json:"status" bson:"status"`
		Href         string         `json:"href"`
		Title        string         `json:"title,omitempty" bson:"title,omitempty"`
		Description  string         `json:"description,omitempty" bson:"description,omitempty"`
		Image        string         `json:"image,omitempty" bson:"image,omitempty"`
		ProductPrice []ProductPrice `json:"productPrice,omitempty" bson:"productPrice,omitempty"`
		LastUpdate   time.Time      `json:"lastUpdate,omitempty" bson:"lastUpdate,omitempty"`
	}

	ProductPrice struct {
		ID         string    `json:"id" bson:"id"`
		Type       string    `json:"@type" bson:"@type"`
		Status     string    `json:"status" bson:"status"`
		Href       string    `json:"href"`
		Name       string    `json:"name" bson:"name"`
		Price      Price     `json:"price,omitempty" bson:"price,omitempty"`
		LastUpdate time.Time `json:"lastUpdate" bson:"lastUpdate"`
	}

	Price struct {
		Unit  string  `json:"unit,omitempty" bson:"unit,omitempty"`
		Value float64 `json:"value,omitempty" bson:"value,omitempty"`
	}
)
