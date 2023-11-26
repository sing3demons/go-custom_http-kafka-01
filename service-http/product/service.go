package product

import (
	"fmt"
	"math"
	"strconv"
	"time"

	"github.com/IBM/sarama"
	"github.com/sing3demons/go-product-service/microservice"
	"github.com/sing3demons/go-product-service/producer"

	"github.com/sing3demons/go-product-service/utils"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IProductService interface {
	FindAll(c microservice.IContext) (any, error)
	FindOne(c microservice.IContext) (*Product, error)
	EventCreateProduct(c microservice.IContext, req CreateProductRequest) (string, error)
	EventUpdateProduct(c microservice.IContext, req UpdateProductRequest) (string, error)
	EventDeleteProduct(c microservice.IContext) (string, error)
}
type productService struct {
	r        IProductRepository
	producer sarama.SyncProducer
}

func NewProductService(r IProductRepository, producer sarama.SyncProducer) IProductService {
	return &productService{r, producer}
}

func (s *productService) FindAll(c microservice.IContext) (any, error) {
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

	products, total, err := s.r.FindAndTotal(filter, findOptions)
	if err != nil {
		return nil, err
	}

	response := map[string]any{
		"data":      products,
		"total":     total,
		"page":      page,
		"last_page": int64(math.Ceil(float64(total) / float64(perPage))),
	}

	return response, nil
}

func (s *productService) InsertProduct(req CreateProductRequest) (interface{}, error) {
	id, err := utils.RandomNanoID(11)
	if err != nil {
		return nil, err
	}

	document := CreateProductRequest{
		ID:           id,
		Type:         "products",
		Status:       req.Status,
		Title:        req.Title,
		ProductPrice: req.ProductPrice,
		Description:  req.Description,
		Image:        req.Image,
	}

	if len(req.Category) != 0 {
		for _, v := range req.Category {
			category := Category{
				ID:   v.ID,
				Type: "category",
			}
			if v.Name != "" {
				category.Name = v.Name
			}
			document.Category = append(document.Category, category)
		}
	}

	result, err := s.r.InsertOne(document)
	if err != nil {
		return nil, err
	}

	filter := bson.M{"_id": result}
	findOptions := options.FindOneOptions{}
	product, err := s.r.FindProduct(filter, &findOptions)
	if err != nil {
		return nil, err
	}

	return product, nil
}

func (s *productService) EventCreateProduct(c microservice.IContext, req CreateProductRequest) (string, error) {
	id, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}

	sub, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateToken(sub)
	if err != nil {
		return "", err
	}
	c.SetAuthorization(token)

	document := CreateProductRequest{
		ID:           id,
		Type:         "products",
		Status:       req.Status,
		Title:        req.Title,
		ProductPrice: req.ProductPrice,
		Description:  req.Description,
		Image:        req.Image,
	}
	if len(req.Category) != 0 {
		for _, v := range req.Category {
			category := Category{
				ID:   v.ID,
				Type: "category",
			}
			if v.Name != "" {
				category.Name = v.Name
			}
			document.Category = append(document.Category, category)
		}
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   document,
	}

	eventProducer := producer.NewEventProducer(s.producer)
	if err := eventProducer.Produce("product.created", data); err != nil {
		return "", err
	}

	return id, nil
}

func (s *productService) FindOne(c microservice.IContext) (*Product, error) {
	id := c.Param("id")
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	filter := bson.M{"id": id, "deleteDate": nil}
	findOptions := options.FindOneOptions{}
	product, err := s.r.FindProduct(filter, &findOptions)

	if err != nil {
		return nil, err

	}
	return product, nil
}

func (s *productService) EventUpdateProduct(c microservice.IContext, req UpdateProductRequest) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("id is required")
	}
	sub, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateToken(sub)
	if err != nil {
		return "", err
	}
	c.SetAuthorization(token)

	document := UpdateProductRequest{
		Title:        req.Title,
		ProductPrice: req.ProductPrice,
		Description:  req.Description,
		Image:        req.Image,
	}

	if len(req.Category) > 0 {
		for _, v := range req.Category {
			category := Category{
				ID:   v.ID,
				Type: "category",
			}
			if v.Name != "" {
				category.Name = v.Name
			}
			document.Category = append(document.Category, category)
		}
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   document,
	}

	eventProducer := producer.NewEventProducer(s.producer)
	if err := eventProducer.Produce("product.updated", data); err != nil {
		return "", err
	}

	return id, nil
}

func (s *productService) EventDeleteProduct(c microservice.IContext) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("id is required")
	}

	sub, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}
	token, err := utils.GenerateToken(sub)
	if err != nil {
		return "", err
	}
	c.SetAuthorization(token)

	product, err := s.r.FindProduct(bson.M{"id": id, "deleteDate": nil}, &options.FindOneOptions{})
	if err != nil {
		return "", err
	}

	document := DeleteProductRequest{
		ID:         product.ID,
		DeleteDate: time.Now().UTC(),
	}

	_, header := c.GetHeader()
	productBody := Event{
		Header: header,
		Body:   document,
	}
	eventProducer := producer.NewEventProducer(s.producer)
	if err := eventProducer.Produce("product.deleted", productBody); err != nil {
		return "", err
	}

	if len(product.ProductPrice) != 0 {
		for _, v := range product.ProductPrice {
			price := Event{
				Header: header,
				Body: DeleteProductPriceRequest{
					ID:         v.ID,
					DeleteDate: time.Now().UTC(),
				}}
			if err := eventProducer.Produce("productPrice.deleted", price); err != nil {
				return "", err
			}
		}
	}

	return id, nil
}
