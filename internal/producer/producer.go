package producer

import (
	"encoding/json"
	"github.com/Shopify/sarama"
	"github.com/rs/zerolog"
	"time"
)

type Producer interface {
	Send(message Message) error
}

type producer struct {
	prod  sarama.SyncProducer
	topic string
	log   zerolog.Logger
}

func NewProducer(topic string, brokers []string, log zerolog.Logger) *producer {
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForLocal
	config.Producer.Partitioner = sarama.NewRandomPartitioner
	config.Producer.Return.Successes = true

	syncProducer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		log.Fatal().Err(err).Msgf("failed to create Sarama sync producer: topic %v, brokers %v",
			topic, brokers)
	}

	return &producer{
		prod:  syncProducer,
		topic: topic,
		log: log,
	}
}

func (p *producer) Send(message Message) error {
	bytes, err := json.Marshal(message)
	if err != nil {
		p.log.Err(err).Msg("failed marshaling message to json")
		return err
	}

	msg := &sarama.ProducerMessage{
		Topic:     p.topic,
		Key:       sarama.StringEncoder(p.topic),
		Value:     sarama.StringEncoder(bytes),
		Partition: -1,
		Timestamp: time.Time{},
	}

	_, _, err = p.prod.SendMessage(msg)
	return err
}