package service

import (
	"context"
	"encoding/json"

	"github.com/Shopify/sarama"
	"github.com/mitchellh/mapstructure"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/ozoncp/ocp-offer-api/internal/repo"
	utils "github.com/ozoncp/ocp-offer-api/internal/utils/models"
	"github.com/rs/zerolog/log"
)

type IConsumer interface {
	sarama.ConsumerGroupHandler
	MessageReceived(*sarama.ConsumerMessage)
}

// Consumer represents a Sarama consumer group consumer.
type Consumer struct {
	IConsumer
	Ready  chan bool
	repo   repo.IRepository
	topics []string
	cfg    *sarama.Config
}

func NewConsumer(r repo.IRepository, topics []string, cfg *sarama.Config) IConsumer {
	return &Consumer{
		repo:   r,
		topics: topics,
		cfg:    cfg,
		Ready:  make(chan bool),
	}
}

// Setup is run at the beginning of a new session, before ConsumeClaim.
func (c *Consumer) Setup(sarama.ConsumerGroupSession) error {
	// Mark the consumer as ready
	close(c.Ready)

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
		c.MessageReceived(message)
		session.MarkMessage(message, "")
	}

	return nil
}

func (c *Consumer) MessageReceived(m *sarama.ConsumerMessage) {
	var msg Message

	if err := json.Unmarshal(m.Value, &msg); err != nil {
		log.Error().Err(err).Msg("Message unmarshal error")

		return
	}

	ctx := context.Background()

	if msg.Type == TypeMultiCreateOffers {
		var mapOffers map[string]models.Offer

		if err := mapstructure.Decode(msg.Value, &mapOffers); err != nil {
			log.Error().Err(err).Msg("Message unmarshal error")
		}

		offers := utils.ConvertOffersMapStringToSlice(mapOffers)

		log.Info().
			Uint16("__type", uint16(msg.Type)).
			Interface("offers", offers).
			Msg("Message received")

		if _, err := c.repo.MultiCreateOffer(ctx, offers); err != nil {
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
	case TypeCreateOffer:
		if _, err := c.repo.CreateOffer(ctx, offer); err != nil {
			log.Error().Err(err).Send()
		}

	case TypeUpdateOffer:
		if err := c.repo.UpdateOffer(ctx, offer); err != nil {
			log.Error().Err(err).Send()
		}

	case TypeDeleteOffer:
		if err := c.repo.RemoveOffer(ctx, offer.ID); err != nil {
			log.Error().Err(err).Send()
		}

	default:
		log.Warn().Msgf("Ignore message: %v", msg)
	}
}
