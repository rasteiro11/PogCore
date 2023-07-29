package kafka

import (
	"context"
	"errors"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type subscriber struct {
	options  *subscriberOptions
	provider *kafka.Consumer
	queue    string
}

var (
	ErrCastingMessage = errors.New("error casting message")
	ErrEventIgnored   = errors.New("error event ignored")
)

var _ messaging.Subscriber = (*subscriber)(nil)

func (s *subscriber) Subscribe(ctx context.Context) (messaging.Message, error) {
	ev := s.provider.Poll(s.options.pollTimeout)
	switch e := ev.(type) {
	case *kafka.Message:
		return mapProviderToMessage(
			s.options.decoder, s.queue, e)
	case kafka.Error:
		logger.Of(ctx).Errorf("[kafka.subscriber.Subscribe] e.provider.Poll() returned error: %+v\n", e)
		return nil, e
	default:
		return nil, ErrEventIgnored
	}

}

func (s *subscriber) Commit(ctx context.Context, m messaging.Message) error {
	cast, ok := m.(*message)
	if !ok {
		return ErrCastingMessage
	}

	if _, err := s.provider.CommitMessage(cast.msg); err != nil {
		logger.Of(ctx).Errorf("[kafka.subscriber.Commit] s.provider.CommitMessage() returned error: %+v\n", err)
	}

	logger.Of(ctx).Infof("package=kafka method=Commit message=%s\n", "success commit")

	return nil
}

func (s *subscriber) Close(ctx context.Context) error {
	return s.provider.Close()
}

func NewSubscriber(ctx context.Context, queue string, opts ...SubscriberOption) (messaging.Subscriber, error) {
	options := newSubscriberOption(opts...)

	provider, err := newSubscriberProvider(ctx, options)
	if err != nil {
		return nil, err
	}

	subscriber := &subscriber{
		options:  options,
		queue:    queue,
		provider: provider,
	}

	if err := subscriber.provider.SubscribeTopics([]string{subscriber.queue}, nil); err != nil {
		logger.Of(ctx).Errorf("[kafka.subscriber.NewSubscriber] s.provider.SubscribeTopics() returned error: %+v\n", err)
		return nil, err
	}

	return subscriber, nil
}
