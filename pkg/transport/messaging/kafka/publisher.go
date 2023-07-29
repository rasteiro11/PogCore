package kafka

import (
	"context"
	"errors"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type publisher struct {
	options  *publisherOptions
	provider *kafka.Producer
}

func NewPublisher(ctx context.Context, opts ...PublisherOption) (messaging.Publisher, error) {
	options := newPublisherOption(opts...)

	provider, err := newPublisherProvider(ctx, options)
	if err != nil {
		return nil, err
	}

	publisher := &publisher{
		options:  options,
		provider: provider,
	}

	return publisher, nil
}

func (p *publisher) Publish(ctx context.Context, messages ...messaging.Message) error {
	var publishErr []error
	delivery_chan := make(chan kafka.Event, 10000)

	for _, msg := range messages {
		kafkaMsg, err := kafkaMessageMapper(p.options.encoder, msg)
		if err != nil {
			return err
		}
		err = p.provider.Produce(
			kafkaMsg,
			delivery_chan,
		)
		if err != nil {
			logger.Of(ctx).Errorf("[messaging.kafka.publisher.Publish] p.provider.Produce() returned error: %+v\n", err)
			return err
		}

		e := <-delivery_chan
		m := e.(*kafka.Message)

		if m.TopicPartition.Error != nil {
			logger.Of(ctx).Errorf("[messaging.kafka.publisher.Publish] m.TopicPartition.Error returned error: %+v\n", err)
			publishErr = append(publishErr, m.TopicPartition.Error)
		}
	}

	close(delivery_chan)

	return errors.Join(publishErr...)
}
