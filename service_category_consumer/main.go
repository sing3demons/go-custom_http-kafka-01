package main

import (
	"context"
	"errors"
	"fmt"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"syscall"

	"github.com/IBM/sarama"
	"github.com/joho/godotenv"
	"github.com/sirupsen/logrus"

	"github.com/sing3demons/go-category-service/repository"
	"github.com/sing3demons/go-category-service/service"
)

func init() {
	godotenv.Load(".env")
}

func main() {
	logger := logrus.New()
	logger.SetFormatter(&logrus.JSONFormatter{})
	logger.SetOutput(os.Stdout)
	logger.SetLevel(logrus.InfoLevel)

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.NewBalanceStrategyRange()
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	version, _ := sarama.ParseKafkaVersion("1.0.0")
	config.Version = version

	broker := os.Getenv("KAFKA_BROKERS")
	var servers []string
	if broker != "" {
		servers = strings.Split(broker, ",")
	} else {
		servers = []string{"localhost:9092"}
	}

	groupID := "category-service"
	consumer, err := sarama.NewConsumerGroup(servers, groupID, config)
	if err != nil {
		panic(err)
	}
	defer consumer.Close()

	db, err := ConnectMonoDB()
	if err != nil {
		panic(err)
	}

	topics := []string{
		"category.created",
		"category.deleted",
		"category.updated",
	}

	repo := repository.NewCategory(db, logger)
	serviceCategory := service.NewCategoryEventHandler(repo, logger)
	consumerHandler := service.NewConsumerHandler(serviceCategory)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	var wg sync.WaitGroup
	wg.Add(1)

	logger.Info("Category consumer started...")
	go func() {
		defer wg.Done()
		for {
			if err := consumer.Consume(ctx, topics, consumerHandler); err != nil {
				if errors.Is(err, sarama.ErrClosedConsumerGroup) {
					return
				}
				log.Panicf("Error from consumer: %v", err)
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
		}
	}()

	// Handle graceful shutdown
	sigusr1 := make(chan os.Signal, 1)
	signal.Notify(sigusr1, syscall.SIGUSR1)

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)
	select {
	case <-sigterm:
		fmt.Println("Received termination signal. Initiating shutdown...")
		cancel()
	case <-ctx.Done():
		fmt.Println("terminating: context cancelled")

	}
	// Wait for the consumer to finish processing
	wg.Wait()
}
