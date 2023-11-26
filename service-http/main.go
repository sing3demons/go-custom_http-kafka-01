package main

import (
	"log"
	"net/url"
	"os"

	"github.com/IBM/sarama"
	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/sing3demons/go-product-service/category"
	"github.com/sing3demons/go-product-service/db"
	"github.com/sing3demons/go-product-service/microservice"
	"github.com/sing3demons/go-product-service/price"
	"github.com/sing3demons/go-product-service/product"
	"github.com/sing3demons/go-product-service/utils"
)

func init() {
	if os.Getenv("GIN_MODE") != gin.ReleaseMode {
		godotenv.Load(".env")
	}

	logLevel, err := logrus.ParseLevel(os.Getenv("LOG_LEVEL"))
	if err != nil {
		logLevel = logrus.InfoLevel
	}

	logrus.SetLevel(logLevel)
	logrus.SetFormatter(&logrus.JSONFormatter{})

	// setup gin
	mode := os.Getenv("GIN_MODE")
	if mode == "" {
		mode = gin.DebugMode
	}

	gin.SetMode(mode)
}

func NewSyncProducer(kafkaBrokers []string) (sarama.SyncProducer, error) {
	producer, err := sarama.NewSyncProducer(kafkaBrokers, nil)
	if err != nil {
		return nil, err
	}

	return producer, nil
}

func main() {
	_, err := os.Create("/tmp/live")
	if err != nil {
		log.Fatal(err)
	}
	defer os.Remove("/tmp/live")
	db := db.NewMongoDB()

	broker := os.Getenv("KAFKA_BROKERS")
	if broker == "" {
		broker = "localhost:9092"
	}
	kafkaBrokers := []string{broker}
	producer, err := NewSyncProducer(kafkaBrokers)
	if err != nil {
		panic(err)
	}
	defer producer.Close()

	productRepository := product.NewProductRepository(db.Collection("product"))
	productService := product.NewProductService(productRepository, producer)
	productHandler := product.NewProductHandler(productService)

	ms := microservice.NewMicroservice()
	ms.GET("", func(c microservice.IContext) {
		resp := map[string]any{
			"name":    "go-http-service",
			"version": "1.0.0",
			"start":   "OK",
		}
		c.JSON(200, resp)
	})

	ms.GET("healthz", healthCheck)

	ms.GET("/products", productHandler.FindAll)
	ms.GET("/products/:id", productHandler.FindOne)
	ms.POST("/products", productHandler.InsertProduct)
	ms.PUT("/products", productHandler.InsertProduct)
	ms.DELETE("/products/:id", productHandler.DeleteProduct)

	productPriceRepository := price.NewProductPriceRepository(db.Collection("productPrice"))
	productPriceService := price.NewProductPriceService(productPriceRepository, producer)
	productPriceHandler := price.NewProductPriceHandler(productPriceService)

	ms.GET("/productPrice", productPriceHandler.FindAll)
	ms.GET("/productPrice/:id", productPriceHandler.FindOne)
	ms.DELETE("/productPrice/:id", productPriceHandler.DeleteProductPrice)
	ms.POST("/productPrice", productPriceHandler.InsertProductPrice)

	categoryRepository := category.NewCategoryRepository(db.Collection("category"))
	categoryService := category.NewCategoryService(categoryRepository, producer)
	categoryHandler := category.NewCategoryHandler(categoryService)

	ms.POST("/category", categoryHandler.InsertProduct)
	ms.GET("/category", categoryHandler.FindCategories)
	ms.GET("/category/:id", categoryHandler.FindOne)
	ms.PATCH("/category/:id", categoryHandler.FindOne)

	ms.Start()
}

func healthCheck(c microservice.IContext) {
	type T struct {
		Name    string `json:"name"`
		Version string `json:"version"`
		Status  string `json:"start"`
	}
	url, err := url.Parse("http://localhost:8080/")
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	result, err := utils.HttpGetClient[T](url.String())
	if err != nil {
		c.Error(500, "Internal Server Error", err)
		return
	}
	c.JSON(200, result)
}
