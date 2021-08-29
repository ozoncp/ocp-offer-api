package main

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"os/signal"
	"sync"
	"syscall"

	_ "github.com/jackc/pgx/v4"
	_ "github.com/jackc/pgx/v4/stdlib"
	_ "github.com/lib/pq"

	"github.com/Shopify/sarama"
	"github.com/jmoiron/sqlx"
	"github.com/mitchellh/mapstructure"
	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/producer"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	"github.com/ozoncp/ocp-offer-api/internal/tracer"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
	"github.com/rs/zerolog/log"
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

	consumer := Consumer{
		ready: make(chan bool),
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
			if err := client.Consume(ctx, topics, &consumer); err != nil {
				log.Error().Err(err).Msg("Error from consumer")
			}
			// check if context was cancelled, signaling that the consumer should stop
			if ctx.Err() != nil {
				return
			}
			consumer.ready = make(chan bool)
		}
	}()

	<-consumer.ready // Await till the consumer has been set up
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

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	ready chan bool
	repo  repo.IRepository
}

var (
	batchSize uint = 1
)

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	c.repo = repo.NewRepo(createDB(), batchSize)

	// Mark the consumer as ready
	close(c.ready)

	return nil
}

// Cleanup is run at the end of a session, once all ConsumeClaim goroutines have exited.
func (c *Consumer) Cleanup(sarama.ConsumerGroupSession) error {
	return nil
}

// ConsumeClaim must start a consumer loop of ConsumerGroupClaim's Messages().
func (c *Consumer) ConsumeClaim(session sarama.ConsumerGroupSession, claim sarama.ConsumerGroupClaim) error {
	// NOTE:
	// Do not move the code below to a goroutine.
	// The `ConsumeClaim` itself is called within a goroutine, see:
	// https://github.com/Shopify/sarama/blob/master/consumer_group.go#L27-L29
	for message := range claim.Messages() {
		messageReceived(message, c.repo)
		session.MarkMessage(message, "")
	}

	return nil
}

func messageReceived(m *sarama.ConsumerMessage, r repo.IRepository) {
	var msg producer.Message

	if err := json.Unmarshal(m.Value, &msg); err != nil {
		log.Error().Err(err).Msg("Message unmarshal error")

		return
	}

	ctx := context.Background()

	if msg.Type == producer.TypeMultiCreateOffers {
		var mapOffers map[string]models.Offer

		if err := mapstructure.Decode(msg.Value, &mapOffers); err != nil {
			log.Error().Err(err).Msg("Message unmarshal error")
		}

		offers := utils.ConvertOffersMapStringToSlice(mapOffers)

		log.Info().
			Uint16("__type", uint16(msg.Type)).
			Interface("offers", offers).
			Msg("Message received")

		if _, err := r.MultiCreateOffer(ctx, offers); err != nil {
			log.Error().Err(err).Send()
		}

		return
	}

	var offer models.Offer
	if err := mapstructure.Decode(msg.Value, &offer); err != nil {
		log.Error().Err(err).Msg("Message unmarshal error")
	}

	log.Info().
		Uint16("__type", uint16(msg.Type)).
		Interface("offer", offer).
		Msg("Message received")

	switch msg.Type {
	case producer.TypeCreateOffer:
		if _, err := r.CreateOffer(ctx, offer); err != nil {
			log.Error().Err(err).Send()
		}

	case producer.TypeUpdateOffer:
		if err := r.UpdateOffer(ctx, offer); err != nil {
			log.Error().Err(err).Send()
		}

	case producer.TypeDeleteOffer:
		if err := r.RemoveOffer(ctx, offer.ID); err != nil {
			log.Error().Err(err).Send()
		}

	default:
		log.Warn().Msgf("Ignore message: %v", msg)
	}
}

func createDB() *sqlx.DB {
	dataSourceName := fmt.Sprintf("host=%v port=%v user=%v password=%v dbname=%v sslmode=%v",
		cfg.Database.Host,
		cfg.Database.Port,
		cfg.Database.User,
		cfg.Database.Password,
		cfg.Database.Name,
		cfg.Database.SSLMode,
	)

	db, err := sqlx.Open(cfg.Database.Driver, dataSourceName)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create database connection")

		return nil
	}

	if err = db.Ping(); err != nil {
		log.Fatal().Err(err).Msgf("failed ping the database")

		return nil
	}

	return db
}
