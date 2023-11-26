package price

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

type IProductPriceService interface {
	FindAll(c microservice.IContext) (any, error)
	FindOne(c microservice.IContext) (*ProductPrice, error)
	CreateProductPrice(c microservice.IContext, req CreateProductPrice) (string, error)
	// EventUpdateProduct(c microservice.IContext, req UpdateProductRequest) (string, error)
	DeleteProductPrice(c microservice.IContext) (string, error)
}
type productPriceService struct {
	r        IProductPriceRepository
	producer sarama.SyncProducer
}

func NewProductPriceService(r IProductPriceRepository, producer sarama.SyncProducer) IProductPriceService {
	return &productPriceService{r, producer}
}

func (svc *productPriceService) FindAll(c microservice.IContext) (any, error) {
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

	products, total, err := svc.r.FindAndTotal(filter, findOptions)
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
func (svc *productPriceService) FindOne(c microservice.IContext) (*ProductPrice, error) {
	id := c.Param("id")
	if id == "" {
		return nil, fmt.Errorf("id is required")
	}
	filter := bson.M{"id": id, "deleteDate": nil}
	findOptions := options.FindOneOptions{}
	product, err := svc.r.FindOne(filter, &findOptions)

	if err != nil {
		return nil, err

	}
	return product, nil
}
func (svc *productPriceService) CreateProductPrice(c microservice.IContext, req CreateProductPrice) (string, error) {
	id, err := utils.RandomNanoID(11)
	if err != nil {
		return "", err
	}

	if req.Status != "" {
		if req.Status != "active" && req.Status != "inActive" {
			return "", fmt.Errorf("status is required")
		}
	}

	document := CreateProductPrice{
		ID:     id,
		Name:   req.Name,
		Status: req.Status,
		Price:  req.Price,
	}

	if !req.LastUpdate.IsZero() {
		document.LastUpdate = req.LastUpdate
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   document,
	}

	eventProducer := producer.NewEventProducer(svc.producer)
	if err := eventProducer.Produce("productPrice.created", data); err != nil {
		return "", err
	}

	return id, nil
}
func (svc *productPriceService) DeleteProductPrice(c microservice.IContext) (string, error) {
	id := c.Param("id")
	if id == "" {
		return "", fmt.Errorf("id is required")
	}

	body := DeleteProductPriceRequest{
		ID:         id,
		DeleteDate: time.Now().UTC(),
	}

	_, header := c.GetHeader()
	data := Event{
		Header: header,
		Body:   body,
	}

	eventProducer := producer.NewEventProducer(svc.producer)
	if err := eventProducer.Produce("productPrice.deleted", data); err != nil {
		return "", err
	}

	return id, nil
}
