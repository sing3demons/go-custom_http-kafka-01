package main

import (
	"context"
	"errors"
	"os"
	"os/signal"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	logrus "github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/mongo"
)

type IMicroservice interface {
	LogInfo(message string, fields logrus.Fields)
	LogError(message string, fields logrus.Fields)
	// Consumer Services
	Consume(servers []string, groupID string, topics []string)
}

type Microservice struct {
	logger     *logrus.Logger
	db *mongo.Database
}

func NewMicroservice() IMicroservice {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	db, err := ConnectMonoDB()
	if err != nil {
		logger.Error("Error connecting to MongoDB", err)
		panic(err)
	}

	return &Microservice{logger, db}
}

func (ms *Microservice) Consume(servers []string, groupID string, topics []string) {
	handler := NewConsumerHandler(ms)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest

	version, err := sarama.ParseKafkaVersion("1.0.0")
	if err != nil {
		ms.LogError("Error parsing Kafka version", logrus.Fields{"error": err})
	}
	config.Version = version

	client, err := sarama.NewConsumerGroup(servers, groupID, config)
	if err != nil {
		ms.LogError("Error creating consumer group client", logrus.Fields{
			"error":   err,
			"topic":   "product.created",
			"group":   groupID,
			"brokers": servers,
			"version": version,
		})
		return
	}
	defer client.Close()

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			ms.logger.Info("Kafka consumer has been started...")
			if err := client.Consume(ctx, topics, handler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				ms.LogError("Error from consumer", logrus.Fields{"error": err})
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				ms.LogError("Context has been cancelled", logrus.Fields{"error": ctx.Err()})
				return
			}

		}
	}()

	// Handle graceful shutdown
	consumptionIsPaused := false
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		ms.logger.Info("Received termination signal. Initiating shutdown...")
		cancel()
	case <-ctx.Done():
		ms.logger.Info("terminating: context cancelled")
	case <-sigusr1:
		ms.toggleConsumptionFlow(client, &consumptionIsPaused)
	}
	// Wait for the consumer to finish processing
	wg.Wait()
}

func (ms *Microservice) Start() {
	ms.logger.Info("Microservice has been started...")
}

func (ms *Microservice) toggleConsumptionFlow(client sarama.ConsumerGroup, isPaused *bool) {
	if *isPaused {
		client.ResumeAll()
		ms.logger.Info("Resuming consumption")
	} else {
		client.PauseAll()
		ms.logger.Info("Pausing consumption")
	}

	*isPaused = !*isPaused
}

// Log log message to console
func (ms *Microservice) LogInfo(message string, fields logrus.Fields) {
	ms.logger.WithFields(fields).Info(message)
}

func (ms *Microservice) LogError(message string, fields logrus.Fields) {
	ms.logger.WithFields(fields).Error(message)
}
