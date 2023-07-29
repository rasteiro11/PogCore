package kafka

import (
	"encoding/json"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type PublisherOption func(*publisherOptions)

type publisherOptions struct {
	pollTimeout int
	encoder     messaging.Encoder
	kafkaConfig *kafka.ConfigMap
}

func newPublisherOption(opts ...PublisherOption) *publisherOptions {
	cfg := &publisherOptions{
		encoder: func(m messaging.Message) ([]byte, error) {
			return json.Marshal(m.Payload())
		},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithPublisherServers(bootstrapServers string) PublisherOption {
	return func(so *publisherOptions) {
		if err := so.kafkaConfig.SetKey("bootstrap.servers", bootstrapServers); err != nil {
			logger.Global().Fatalf("[kafka.publisher_option] so.kafkaConfig.SetKey() retunred error: %+v\n", err)
		}
	}
}

func WithAcks(acks string) PublisherOption {
	return func(so *publisherOptions) {
		if err := so.kafkaConfig.SetKey("acks", acks); err != nil {
			logger.Global().Fatalf("[kafka.publisher_option] so.kafkaConfig.SetKey() retunred error: %+v\n", err)
		}
	}
}

func WithEncoder(encoder messaging.Encoder) PublisherOption {
	return func(so *publisherOptions) {
		so.encoder = encoder
	}
}
