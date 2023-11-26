package category

import (
	"fmt"
	"math"
	"strconv"

	"github.com/IBM/sarama"
	"github.com/sing3demons/go-product-service/microservice"
	"github.com/sing3demons/go-product-service/producer"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/sing3demons/go-product-service/utils"
)

type ICategoryService interface {
	FindAll(c microservice.IContext) (any, error)
	FindOne(c microservice.IContext) (*Category, error)
	// EventCreateProduct(c microservice.IContext, req CreateProductRequest) (string, error)
	// EventUpdateProduct(c microservice.IContext, req UpdateProductRequest) (string, error)
	// EventDeleteProduct(c microservice.IContext) (string, error)
	CreateCategory(c microservice.IContext, req CreateCategoryReq) (string, error)
}
type categoryService struct {
	r        ICategoryRepository
	producer sarama.SyncProducer
}

func NewCategoryService(r ICategoryRepository, producer sarama.SyncProducer) ICategoryService {
	return &categoryService{r, producer}
}

func (s *categoryService) FindAll(c microservice.IContext) (any, error) {
	filter := bson.M{}
	findOptions := options.Find()
	if s := c.QueryString("s"); s != "" {
		filter = bson.M{
			"$or": []bson.M{
				{
					"title": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
				{
					"description": bson.M{
						"$regex": primitive.Regex{
							Pattern: s,
							Options: "i",
						},
					},
				},
			},
		}
	}

	if sort := c.QueryString("sort"); sort != "" {
		if sort == "asc" {
			findOptions.SetSort(bson.D{{Key: "price", Value: 1}})
		} else if sort == "desc" {
			findOptions.SetSort(bson.D{{Key: "price", Value: -1}})
		}
	}

	p := c.QueryString("page")
	if p == "" {
		p = "1"
	}

	page, _ := strconv.Atoi(p)
	if page == 0 {
		page = 1
	}
	var perPage int64 = 9

	findOptions.SetSkip((int64(page) - 1) * perPage)
	findOptions.SetLimit(perPage)

	categories, total, err := s.r.FindAndTotal(filter, findOptions)
	if err != nil {
		return nil, err
	}

	response := map[string]any{
		"data":      categories,
		"total":     total,
		"page":      page,
		"last_page": int64(math.Ceil(float64(total) / float64(perPage))),
	}

	return response, nil
}

func (s *categoryService) FindOne(c microservice.IContext) (*Category, error) {
	id := c.Param("id")
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	filter := bson.M{"id": id, "deleteDate": nil}
	findOptions := options.FindOneOptions{}
	category, err := s.r.FindOne(filter, &findOptions)

	if err != nil {
		return nil, err

	}
	return category, nil
}

func (s *categoryService) CreateCategory(c microservice.IContext, req CreateCategoryReq) (string, error) {
	id, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}

	document := CreateCategoryReq{
		ID:     id,
		Name:   req.Name,
		Type:   "category",
		Status: req.Status,
	}

	if !req.LastUpdate.IsZero() {
		document.LastUpdate = req.LastUpdate
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   document,
	}

	eventProducer := producer.NewEventProducer(s.producer)
	if err := eventProducer.Produce("category.created", data); err != nil {
		return "", err
	}

	return id, nil
}

func (s *categoryService) UpdateCategory(c microservice.IContext, req UpdateCategoryReq) (string, error) {
	id, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}

	var document UpdateCategoryReq
	if !req.LastUpdate.IsZero() {
		document.LastUpdate = req.LastUpdate
	}

	if len(req.Products) != 0 {
		for _, v := range req.Products {
			product := Product{
				ID:   v.ID,
				Type: "product",
			}
			if v.Title != "" {
				product.Title = v.Title
			}
			document.Products = append(document.Products, product)
		}
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   document,
	}

	eventProducer := producer.NewEventProducer(s.producer)
	if err := eventProducer.Produce("category.updated", data); err != nil {
		return "", err
	}

	return id, nil
}
