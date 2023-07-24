package kafka

import (
	"encoding/json"
	"reflect"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type SubscriberOption func(*subscriberOptions)

type subscriberOptions struct {
	pollTimeout int
	decoder     messaging.Decoder
	kafkaConfig *kafka.ConfigMap
}

func newSubscriberOption(opts ...SubscriberOption) *subscriberOptions {
	cfg := &subscriberOptions{
		decoder: func(b []byte) (any, error) {
			return b, nil
		},
		kafkaConfig: &kafka.ConfigMap{},
	}

	for _, opt := range opts {
		opt(cfg)
	}

	return cfg
}

func WithPollTimeout(pollTimeout int) SubscriberOption {
	return func(so *subscriberOptions) {
		so.pollTimeout = pollTimeout
	}
}

func WithOffsetReset(offsetReset string) SubscriberOption {
	return func(so *subscriberOptions) {
		if err := so.kafkaConfig.SetKey("auto.offset.reset", offsetReset); err != nil {
			logger.Global().Fatalf("[kafka.subscriber_option] so.kafkaConfig.SetKey() retunred error: %+v\n", err)
		}
	}
}

func WithBootstrapServers(bootstrapServers string) SubscriberOption {
	return func(so *subscriberOptions) {
		if err := so.kafkaConfig.SetKey("bootstrap.servers", bootstrapServers); err != nil {
			logger.Global().Fatalf("[kafka.subscriber_option] so.kafkaConfig.SetKey() retunred error: %+v\n", err)
		}
	}
}

func WithGroupID(groupId string) SubscriberOption {
	return func(so *subscriberOptions) {
		if err := so.kafkaConfig.SetKey("group.id", groupId); err != nil {
			logger.Global().Fatalf("[kafka.subscriber_option] so.kafkaConfig.SetKey() retunred error: %+v\n", err)
		}
	}
}

func WithDecoder(decoder messaging.Decoder) SubscriberOption {
	return func(so *subscriberOptions) {
		so.decoder = decoder
	}
}

func WithDecoderTarget(typeof any) SubscriberOption {
	return func(s *subscriberOptions) {
		v := reflect.TypeOf(typeof)

		s.decoder = func(b []byte) (any, error) {
			target := reflect.New(v).Interface()
			if err := json.Unmarshal(b, target); err != nil {
				return nil, err
			}

			return target, nil
		}
	}
}
