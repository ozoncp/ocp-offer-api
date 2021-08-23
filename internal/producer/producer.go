package producer

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/Shopify/sarama"
	"github.com/rs/zerolog/log"
)

type Producer interface {
	Send(msg Message) error
}

type producer struct {
	dataProducer sarama.SyncProducer
	topic        string
	msgChan      chan *sarama.ProducerMessage
}

type MessageType int

const (
	Create MessageType = iota
	Update
	Remove
)

type Message struct {
	Type MessageType
	Body map[string]interface{}
}

func New(ctx context.Context, addrs []string, topic string, capacity uint64) (Producer, error) {
	config := sarama.NewConfig()
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.RequiredAcks = sarama.WaitForAll
	config.Producer.Return.Successes = true

	p, err := sarama.NewSyncProducer(addrs, config)

	if err != nil {
		return nil, err
	}

	newProducer := &producer{
		dataProducer: p,
		topic:        topic,
		msgChan:      make(chan *sarama.ProducerMessage, capacity),
	}

	go newProducer.handleMessage(ctx)

	return newProducer, nil
}

func (dProducer *producer) handleMessage(ctx context.Context) {
	for {
		select {
		case msg := <-dProducer.msgChan:
			partition, offset, err := dProducer.dataProducer.SendMessage(msg)

			if err != nil {
				log.Error().Err(err).Msg("Failed to send message to kafka, retry ...")
				dProducer.msgChan <- msg
			} else {
				log.Info().
					Int32("partition", partition).
					Int64("offset", offset).
					Msg("Message sent successfully to kafka")
			}

		case <-ctx.Done():
			close(dProducer.msgChan)
			dProducer.dataProducer.Close()
			return
		}
	}
}

func (dProducer *producer) Send(msg Message) error {

	dataBytes, err := json.Marshal(msg)

	if err != nil {
		log.Error().Err(err).Msg("Failed to marshal message to json")
		return err
	}

	dProducer.msgChan <- &sarama.ProducerMessage{
		Topic:     dProducer.topic,
		Key:       sarama.StringEncoder(dProducer.topic),
		Value:     sarama.StringEncoder(dataBytes),
		Partition: -1,
		Timestamp: time.Time{},
	}

	return nil
}

func CreateMessage(msgType MessageType, id uint64, timestamp time.Time) Message {
	return Message{
		Type: msgType,
		Body: map[string]interface{}{
			"Id":        id,
			"Operation": fmt.Sprintf("%s offer", convertMessageTypeToString(msgType)),
			"Timestamp": timestamp,
		},
	}
}

func convertMessageTypeToString(msgType MessageType) string {
	switch msgType {
	case Create:
		return "Created"
	case Update:
		return "Updated"
	case Remove:
		return "Removed"
	}

	return "Unknown message type"
}
