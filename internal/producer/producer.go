package producer

import (
	"context"
	"encoding/json"
	"time"

	"github.com/Shopify/sarama"
	"github.com/fatih/structs"
	"github.com/opentracing/opentracing-go"
	"github.com/ozoncp/ocp-offer-api/internal/models"
	"github.com/rs/zerolog/log"
)

type IProducer interface {
	CreateOffer(offer models.Offer) error
	UpdateOffer(offer models.Offer) error
	DeleteOffer(offerID uint64) error
}

type Producer struct {
	producer    sarama.SyncProducer
	topicName   string
	messageChan chan *sarama.ProducerMessage
}

type MessageType uint16

const (
	TypeCreateOffer MessageType = iota
	TypeUpdateOffer
	TypeDeleteOffer
)

type Message struct {
	Type  MessageType
	Value map[string]interface{}
}

func New(ctx context.Context, brokers []string, topicName string, capacity uint64) (IProducer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		return nil, err
	}

	p := &Producer{
		producer:    producer,
		topicName:   topicName,
		messageChan: make(chan *sarama.ProducerMessage, capacity),
	}

	go p.listener(ctx)

	return p, nil
}

func (p *Producer) CreateOffer(offer models.Offer) error {
	return p.publish("Producer.CreateOffer", TypeCreateOffer, structs.Map(offer))
}

func (p *Producer) UpdateOffer(offer models.Offer) error {
	return p.publish("Producer.UpdateOffer", TypeUpdateOffer, structs.Map(offer))
}

func (p *Producer) DeleteOffer(offerID uint64) error {
	return p.publish("Producer.DeleteOffer", TypeDeleteOffer, structs.Map(models.Offer{ID: offerID}))
}

// ---

func (p *Producer) listener(ctx context.Context) {
	for {
		select {
		case msg := <-p.messageChan:
			partition, offset, err := p.producer.SendMessage(msg)

			if err != nil {
				log.Error().Msgf("failed to send message to kafka: %v", err)
				log.Error().Msgf("retry ...")

				p.messageChan <- msg
			}

			log.Info().
				Int32("partition", partition).
				Str("topic", msg.Topic).
				Msgf("Delivered message to topic %s [%d] at offset %v", msg.Topic, partition, offset)

		case <-ctx.Done():
			close(p.messageChan)
			p.producer.Close()

			return
		}
	}
}

func (p *Producer) publish(spanName string, msgType MessageType, value map[string]interface{}) error {
	span := opentracing.GlobalTracer().StartSpan(spanName)
	defer span.Finish()

	b, err := json.Marshal(
		Message{
			Type:  msgType,
			Value: value,
		})
	if err != nil {
		return err
	}

	p.messageChan <- &sarama.ProducerMessage{
		Topic:     p.topicName,
		Key:       sarama.StringEncoder(p.topicName),
		Value:     sarama.StringEncoder(b),
		Partition: -1,
		Timestamp: time.Now(),
	}

	return nil
}
