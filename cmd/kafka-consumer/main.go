package main

import (
	"encoding/json"
	"fmt"

	"github.com/Shopify/sarama"
	cfg "github.com/ozoncp/ocp-offer-api/internal/config"
	"github.com/ozoncp/ocp-offer-api/internal/producer"
	"github.com/rs/zerolog/log"
)

func subscribe(topic string, consumer sarama.Consumer) error {
	// get all partitions on the given topic
	partitionList, err := consumer.Partitions(topic)

	if err != nil {
		return err
	}

	// get offset for the oldest message on the topic.
	initialOffset := sarama.OffsetOldest

	for _, partition := range partitionList {
		pc, err := consumer.ConsumePartition(topic, partition, initialOffset)

		if err != nil {
			return err
		}

		for message := range pc.Messages() {
			messageReceived(message)
		}
	}

	return nil
}

func messageReceived(message *sarama.ConsumerMessage) {
	var msg producer.Message
	err := json.Unmarshal(message.Value, &msg)

	if err != nil {
		fmt.Printf("Error unmarshalling message: %s\n", err)
	}

	log.Info().Msgf("Message: %v", msg.Body)
}

func main() {
	consumer, err := sarama.NewConsumer(cfg.Kafka.Brokers, nil)
	if err != nil {
		log.Fatal().Msgf("NewConsumer error: %v", err)
	}

	log.Info().Msgf("Waiting messages from Kafka ...")
	err = subscribe(cfg.Kafka.Topic, consumer)

	if err != nil {
		log.Fatal().Err(err).Msgf("Subscribe failed")
	}
}
