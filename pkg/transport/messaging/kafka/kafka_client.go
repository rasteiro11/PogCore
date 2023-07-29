package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
)

func newSubscriberProvider(ctx context.Context, subscriberOptions *subscriberOptions) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(subscriberOptions.kafkaConfig)
	if err != nil {
		logger.Of(ctx).Errorf("[kafka.kafka_client] kafka.NewConsumer() returned error: %+v\n", err)
		return nil, err
	}
	return consumer, nil
}

func newPublisherProvider(ctx context.Context, subscriberOptions *publisherOptions) (*kafka.Producer, error) {
	consumer, err := kafka.NewProducer(subscriberOptions.kafkaConfig)
	if err != nil {
		logger.Of(ctx).Errorf("[kafka.kafka_client] kafka.NewConsumer() returned error: %+v\n", err)
		return nil, err
	}
	return consumer, nil
}
