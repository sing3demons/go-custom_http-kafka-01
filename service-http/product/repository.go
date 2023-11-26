package product

import (
	"context"
	"time"

	"github.com/sing3demons/go-product-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IProductRepository interface {
	FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]Product, int64, error)
	InsertOne(document interface{}) (interface{}, error)
	FindProduct(filter bson.M, findOptions *options.FindOneOptions) (*Product, error)
}

type productRepository struct {
	collection *mongo.Collection
}

func NewProductRepository(collection *mongo.Collection) IProductRepository {
	return &productRepository{collection}
}

func (r *productRepository) FindAll(filter bson.M, findOptions *options.FindOptions) ([]Product, error) {
	products := []Product{}
	result, err := utils.GetMulti[Product](r.collection, filter, findOptions)
	if err != nil {
		return nil, err
	}

	for _, p := range result {

		var productPrice []ProductPrice
		if len(p.ProductPrice) != 0 {
			for _, v := range p.ProductPrice {
				p, err := utils.HttpGetClient[ProductPrice]("http://localhost:8080/productPrice/" + v.ID)
				if err != nil {
					productPrice = append(productPrice, ProductPrice{
						ID:   p.ID,
						Type: p.Type,
						Name: p.Name,
					})
				}
				productPrice = append(productPrice, ProductPrice{
					ID:         p.ID,
					Type:       p.Type,
					Href:       utils.Href(p.Type, p.ID),
					Price:      p.Price,
					Name:       p.Name,
					Status:     p.Status,
					LastUpdate: p.LastUpdate,
				})
			}
		} else {
			productPrice = p.ProductPrice
		}
		products = append(products, Product{
			ID:           p.ID,
			Type:         p.Type,
			Href:         utils.Href(p.Type, p.ID),
			Title:        p.Title,
			ProductPrice: productPrice,
			Description:  p.Description,
			Image:        p.Image,
			Status:       p.Status,
			LastUpdate:   p.LastUpdate,
		})
	}

	return products, nil
}

func (r *productRepository) FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]Product, int64, error) {
	products := []Product{}
	result, total, err := utils.GetMultiWithTotal[Product](r.collection, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	for _, p := range result {
		var categories []Category
		if len(p.Category) != 0 {
			for _, v := range p.Category {
				c, err := utils.HttpGetClient[Category]("http://localhost:8080/category/" + v.ID)
				if err != nil {
					categories = append(categories, Category{
						ID:   c.ID,
						Type: c.Type,
						Name: c.Name,
					})
				}
				categories = append(categories, Category{
					ID:   c.ID,
					Type: c.Type,
					Href: utils.Href(c.Type, c.ID),
					Name: c.Name,
				})
			}
		} else {
			categories = p.Category
		}
		var productPrice []ProductPrice
		if len(p.ProductPrice) != 0 {
			for _, v := range p.ProductPrice {
				price, err := utils.HttpGetClient[ProductPrice]("http://localhost:8080/productPrice/" + v.ID)
				if err != nil {
					productPrice = append(productPrice, ProductPrice{
						ID:   price.ID,
						Type: price.Type,
					})
					continue
				}
				if price.ID != "" {
					productPrice = append(productPrice, ProductPrice{
						ID:         price.ID,
						Type:       price.Type,
						Href:       utils.Href(price.Type, price.ID),
						Status:     price.Status,
						Price:      price.Price,
						Name:       price.Name,
						LastUpdate: price.LastUpdate,
					})
				} else {
					r := ProductPrice{
						ID:   v.ID,
						Type: "productPrice",
					}
					productPrice = append(productPrice, r)

				}
			}
		} else {
			productPrice = append(productPrice, ProductPrice{
				ID:   p.ID,
				Type: p.Type,
			})
		}
		products = append(products, Product{
			MID:          p.MID,
			ID:           p.ID,
			Type:         p.Type,
			Category:     categories,
			Status:       p.Status,
			Href:         utils.Href(p.Type, p.ID),
			Title:        p.Title,
			ProductPrice: productPrice,
			Description:  p.Description,
			Image:        p.Image,
			LastUpdate:   p.LastUpdate,
		})
	}

	return products, total, nil
}

func (r *productRepository) FindProduct(filter bson.M, findOptions *options.FindOneOptions) (*Product, error) {
	p, err := utils.GetOne[Product](r.collection, filter, findOptions)
	if err != nil {
		if err == mongo.ErrNoDocuments {
			return nil, mongo.ErrNoDocuments
		}
		return nil, err
	}

	var categories []Category
	if len(p.Category) != 0 {
		for _, v := range p.Category {
			c, err := utils.HttpGetClient[Category]("http://localhost:8080/category/" + v.ID)
			if err != nil {
				categories = append(categories, Category{
					ID:   c.ID,
					Type: c.Type,
				})
			}
			categories = append(categories, Category{
				ID:   c.ID,
				Type: c.Type,
				Href: utils.Href(c.Type, c.ID),
				Name: c.Name,
			})
		}
	} else {
		categories = p.Category
	}

	var productPrice []ProductPrice
	if len(p.ProductPrice) != 0 {
		for _, v := range p.ProductPrice {
			price, err := utils.HttpGetClient[ProductPrice]("http://localhost:8080/productPrice/" + v.ID)
			if err != nil {
				productPrice = append(productPrice, ProductPrice{
					ID:   price.ID,
					Type: price.Type,
				})
				continue
			}
			if price.ID != "" {
				productPrice = append(productPrice, ProductPrice{
					ID:         price.ID,
					Type:       price.Type,
					Href:       utils.Href(price.Type, price.ID),
					Status:     price.Status,
					Price:      price.Price,
					Name:       price.Name,
					LastUpdate: price.LastUpdate,
				})
			} else {
				r := ProductPrice{
					ID:   v.ID,
					Type: "productPrice",
				}
				productPrice = append(productPrice, r)

			}
		}
	} else {
		productPrice = append(productPrice, ProductPrice{
			ID:   p.ID,
			Type: p.Type,
		})
	}

	product := &Product{
		MID:          p.MID,
		ID:           p.ID,
		Type:         p.Type,
		Status:       p.Status,
		Category:     categories,
		Href:         utils.Href(p.Type, p.ID),
		Title:        p.Title,
		ProductPrice: productPrice,
		Description:  p.Description,
		Image:        p.Image,
		LastUpdate:   p.LastUpdate,
	}

	return product, nil
}

func (r *productRepository) InsertOne(document interface{}) (interface{}, error) {
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	result, err := r.collection.InsertOne(ctx, document)
	if err != nil {
		return nil, err
	}

	return result.InsertedID, nil
}
