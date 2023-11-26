package producer

import (
	"encoding/json"

	"github.com/IBM/sarama"
	logger "github.com/sirupsen/logrus"
)

type eventProducer struct {
	producer   sarama.SyncProducer
}

func NewEventProducer(producer sarama.SyncProducer) eventProducer {
	return eventProducer{producer}
}

func (e *eventProducer) Produce(topic string, event any) (err error) {
	value, err := json.Marshal(event)
	if err != nil {
		return err
	}

	msg := sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(value),
	}
	partition, offset, err := e.producer.SendMessage(&msg)
	if err != nil {
		return err
	}

	logger.WithFields(logger.Fields{
		"topic":     topic,
		"partition": partition,
		"offset":    offset,
		"event":     event,
	}).Info("send message to kafka")

	return nil
}
