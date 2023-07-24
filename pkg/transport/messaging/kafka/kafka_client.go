package kafka

import (
	"context"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
)

func newProvider(ctx context.Context, subscriberOptions *subscriberOptions) (*kafka.Consumer, error) {
	consumer, err := kafka.NewConsumer(subscriberOptions.kafkaConfig)
	if err != nil {
		logger.Of(ctx).Errorf("[kafka.kafka_client] kafka.NewConsumer() returned error: %+v\n", err)
		return nil, err
	}
	return consumer, nil
}
