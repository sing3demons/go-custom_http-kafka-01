package main

import (
	"encoding/json"
	"strings"
	"time"

	"github.com/IBM/sarama"
	"github.com/sing3demons/go-consumer-service/middleware"
	"github.com/sing3demons/go-consumer-service/services"
	logrus "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type consumerHandler struct {
	db     *mongo.Database
	logger *logrus.Logger
}

func NewConsumerHandler(ms *Microservice) sarama.ConsumerGroupHandler {
	return consumerHandler{db: ms.db, logger: ms.logger}
}

func (obj consumerHandler) Setup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerHandler) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

func (obj consumerHandler) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	ev := services.NewService(obj.db, obj.logger)
	for msg := range claim.Messages() {
		timestamp := time.Now().Format("20060102150405")
		obj.logger.WithFields(logrus.Fields{
			"timestamp": timestamp,
			"topic":     msg.Topic,
			"partition": msg.Partition,
			"offset":    msg.Offset,
			"key":       string(msg.Key),
			"value":     string(msg.Value),
		}).Info("Consume message")

		validate := obj.validateHeader(msg.Value)
		if validate {
			switch msg.Topic {
			case "product.created":
				ev.InsertProduct(msg.Topic, msg.Value, timestamp)
			case "product.deleted":
				ev.DeleteProduct(msg.Topic, msg.Value, timestamp)
			case "productPrice.created":
				ev.InsertProductPrice(msg.Topic, msg.Value, timestamp)
			case "productPrice.deleted":
				ev.DeleteProductPrice(msg.Topic, msg.Value, timestamp)
			}
		} else {
			obj.logger.WithFields(logrus.Fields{
				"timestamp": timestamp,
				"topic":     msg.Topic,
				"partition": msg.Partition,
				"offset":    msg.Offset,
				"key":       string(msg.Key),
				"value":     string(msg.Value),
			}).Error("Consume message")
		}
		session.MarkMessage(msg, "")
	}

	return nil
}

type Event struct {
	Header map[string]any `json:"header"`
	Body   any            `json:"body"`
}

func (obj consumerHandler) validateHeader(msg []byte) bool {
	var event Event
	var header map[string]any
	if err := json.Unmarshal(msg, &event); err != nil {
		obj.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("Error decoding data")
		return false
	}
	header = event.Header
	Authorization, ok := header["Authorization"]
	if !ok {
		obj.logger.WithFields(logrus.Fields{
			"error": "authorization is required",
		}).Error("Error decoding data")
		return false
	}

	token := strings.Split(Authorization.(string), "Bearer ")[1]
	if token == "" {
		obj.logger.WithFields(logrus.Fields{
			"error": "authorization is required",
		}).Error("Error decoding data")
		return false
	}

	mapClaims, err := middleware.ValidateToken(token)
	if err != nil {
		obj.logger.WithFields(logrus.Fields{
			"error": err,
		}).Error("validate token error")
		return false
	}

	obj.logger.WithFields(logrus.Fields{
		"claims": mapClaims,
	}).Info("validate token success")

	return true
}
