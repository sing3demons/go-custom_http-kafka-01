package main

import (
	"os"
	"strings"

	"github.com/joho/godotenv"
)

const (
	ProductCreatedTopic        = "product.created"
	ProductUpdatedTopic        = "product.updated"
	ProductDeletedTopic        = "product.deleted"
	ProductProductCreatedTopic = "productPrice.created"
	ProductPriceDeleteTopic    = "productPrice.deleted"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	ms := NewMicroservice()

	broker := os.Getenv("KAFKA_BROKERS")
	var kafkaBrokers []string
	if broker != "" {
		kafkaBrokers = strings.Split(broker, ",")
	} else {
		kafkaBrokers = []string{"localhost:9092"}
	}

	consumerGroupID := "product_consumer_group"
	topics := []string{
		ProductCreatedTopic,
		ProductUpdatedTopic,
		ProductDeletedTopic,
		ProductProductCreatedTopic,
		ProductPriceDeleteTopic,
	}

	ms.Consume(kafkaBrokers, consumerGroupID, topics)
}
