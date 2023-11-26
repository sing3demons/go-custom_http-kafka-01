package category

import (
	"github.com/sing3demons/go-product-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type ICategoryRepository interface {
	FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]Category, int64, error)
	FindOne(filter bson.M, findOptions *options.FindOneOptions) (*Category, error)
}

type categoryRepository struct {
	collection *mongo.Collection
}

func NewCategoryRepository(collection *mongo.Collection) ICategoryRepository {
	return &categoryRepository{collection}
}

func (r *categoryRepository) FindOne(filter bson.M, findOptions *options.FindOneOptions) (category *Category, err error) {
	category, err = utils.GetOne[Category](r.collection, filter, findOptions)
	if err != nil {
		return nil, err
	}

	var products []Product
	if len(category.Products) != 0 {
		for _, v := range category.Products {
			p, err := utils.HttpGetClient[Product]("http://localhost:8080/product/" + v.ID)
			if err != nil {
				products = append(products, Product{
					ID:   p.ID,
					Type: p.Type,
				})
			}
			products = append(products, Product{
				ID:           p.ID,
				Type:         p.Type,
				Href:         utils.Href(p.Type, p.ID),
				Status:       p.Status,
				Title:        p.Title,
				Description:  p.Description,
				Image:        p.Image,
				ProductPrice: p.ProductPrice,
				LastUpdate:   p.LastUpdate,
			})
		}
	}

	result := Category{
		ID:         category.ID,
		Type:       category.Type,
		Status:     category.Status,
		Href:       utils.Href(category.Type, category.ID),
		Name:       category.Name,
		Products:   products,
		LastUpdate: category.LastUpdate,
	}

	return &result, nil
}

func (r *categoryRepository) FindAndTotal(filter bson.M, findOptions *options.FindOptions) ([]Category, int64, error) {
	categories := []Category{}
	result, total, err := utils.GetMultiWithTotal[Category](r.collection, filter, findOptions)
	if err != nil {
		return nil, 0, err
	}

	for _, category := range result {

		var products []Product
		if len(category.Products) != 0 {
			for _, v := range category.Products {
				p, err := utils.HttpGetClient[Product]("http://localhost:8080/product/" + v.ID)
				if err != nil {
					products = append(products, Product{
						ID:   p.ID,
						Type: p.Type,
					})
				}
				products = append(products, Product{
					ID:           p.ID,
					Type:         p.Type,
					Href:         utils.Href(p.Type, p.ID),
					Status:       p.Status,
					Title:        p.Title,
					Description:  p.Description,
					Image:        p.Image,
					ProductPrice: p.ProductPrice,
					LastUpdate:   p.LastUpdate,
				})
			}
		}
		categories = append(categories, Category{
			ID:         category.ID,
			Type:       category.Type,
			Status:     category.Status,
			Href:       utils.Href(category.Type, category.ID),
			Name:       category.Name,
			Products:   products,
			LastUpdate: category.LastUpdate,
		})

	}

	return categories, total, nil
}
