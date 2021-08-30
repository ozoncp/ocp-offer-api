package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/Shopify/sarama"
	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/database"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	"github.com/ozoncp/ocp-offer-api/internal/service"
	"github.com/ozoncp/ocp-offer-api/internal/tracer"
	"github.com/rs/zerolog/log"
)

const (
	batchSize = 2
)

func main() {
	tracer.InitTracing("ocp_offer_api-kafka_consumer")

	version, err := sarama.ParseKafkaVersion("2.8.0")
	if err != nil {
		log.Fatal().Err(err).Msg("Error parsing Kafka version")
	}

	config := sarama.NewConfig()
	config.Consumer.Group.Rebalance.Strategy = sarama.BalanceStrategyRange
	config.Consumer.Offsets.Initial = sarama.OffsetOldest
	config.Version = version

	topics := []string{cfg.Kafka.Topic}

	dsn := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db := database.NewPostgres(dsn, cfg.Database.Driver)
	r := repo.NewRepo(db, batchSize)

	consumer, ok := service.NewConsumer(r, topics, config).(*service.Consumer)
	if !ok {
		log.Fatal().Err(err).Msg("Error creating consumer")
	}

	ctx, cancel := context.WithCancel(context.Background())
	client, err := sarama.NewConsumerGroup(cfg.Kafka.Brokers, cfg.Kafka.GroupID, config)
	if err != nil {
		log.Fatal().Err(err).Msg("Error creating consumer group client")
	}

	wg := &sync.WaitGroup{}
	wg.Add(1)
	go func() {
		defer wg.Done()
		for {
			if err := client.Consume(ctx, topics, consumer); err != nil {
				log.Error().Err(err).Msg("Error from consumer")
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.Ready = make(chan bool)
		}
	}()

	<-consumer.Ready // Await till the consumer has been set up
	log.Info().Msg("Sarama consumer up and running!...")

	sigterm := make(chan os.Signal, 1)
	signal.Notify(sigterm, syscall.SIGINT, syscall.SIGTERM)

	select {
	case <-ctx.Done():
		log.Info().Msg("terminating: context cancelled")
	case <-sigterm:
		log.Info().Msg("terminating: via signal")
	}

	cancel()

	wg.Wait()

	if err = client.Close(); err != nil {
		log.Error().Msgf("Error closing client: %v", err)
	}
}
