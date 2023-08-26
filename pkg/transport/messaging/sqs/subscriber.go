package sqs

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type Subscriber struct {
	options  *subscriberOptions
	provider *sqs.Client
	queue    string
}

var _ messaging.Subscriber = (*Subscriber)(nil)

func (s *Subscriber) Subscribe(ctx context.Context) (messaging.Message, error) {
	result, err := s.provider.ReceiveMessage(ctx, mapToReceiveMessageInput(s.queue))
	if err != nil {
		return nil, err
	}

	if len(result.Messages) == 0 {
		return nil, ErrNoMessageInQueue
	}

	return mapProviderToMessage(
		s.options.decoder, s.queue, result.Messages[0])
}

func (s *Subscriber) Commit(ctx context.Context, message messaging.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	input := mapToCommitMessageInput(s.queue, message)
	_, err := s.provider.DeleteMessage(ctx, input)
	if err != nil {
		return err
	}

	logger.Of(ctx).Infof("package=sqs method=Commit message=%s\n", "success commit")

	return nil
}

func (s *Subscriber) Close(ctx context.Context) error {
	return nil
}

func NewSubscriber(ctx context.Context, queue string, opts ...SubscriberOption) (*Subscriber, error) {
	options := newSubscriberOption(opts...)

	provider, err := newProvider(ctx)
	if err != nil {
		return nil, err
	}

	subscriber := &Subscriber{
		options:  options,
		queue:    queue,
		provider: provider,
	}

	return subscriber, nil
}
