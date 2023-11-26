package price

import (
	"github.com/sing3demons/go-product-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IProductPriceRepository interface {
	FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]ProductPrice, int64, error)
	FindOne(filter bson.M, findOptions *options.FindOneOptions) (*ProductPrice, error)
}

type productPriceRepository struct {
	collection *mongo.Collection
}

func NewProductPriceRepository(collection *mongo.Collection) IProductPriceRepository {
	return &productPriceRepository{collection}
}

func (r *productPriceRepository) FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]ProductPrice, int64, error) {
	productPrices := []ProductPrice{}
	result, total, err := utils.GetMultiWithTotal[ProductPrice](r.collection, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	for _, p := range result {
		productPrices = append(productPrices, ProductPrice{
			ID:         p.ID,
			Type:       p.Type,
			Status:     p.Status,
			Href:       utils.Href(p.Type, p.ID),
			Price:      p.Price,
			Name:       p.Name,
			LastUpdate: p.LastUpdate,
		})
	}
	return productPrices, total, nil
}

func (r *productPriceRepository) FindOne(filter bson.M, findOptions *options.FindOneOptions) (*ProductPrice, error) {

	p, err := utils.GetOne[ProductPrice](r.collection, filter, findOptions)
	if err != nil {
		return nil, err
	}

	productPrice := ProductPrice{
		ID:         p.ID,
		Type:       p.Type,
		Status:     p.Status,
		Href:       utils.Href(p.Type, p.ID),
		Name:       p.Name,
		Price:      p.Price,
		LastUpdate: p.LastUpdate,
	}
	return &productPrice, nil
}
