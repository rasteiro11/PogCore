package sqs

import (
	"context"
	"errors"

	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type Publisher struct {
	options  *publisherOption
	provider *sqs.Client
}

var _ messaging.Publisher = (*Publisher)(nil)
	
func (p *Publisher) Publish(ctx context.Context, messages ...messaging.Message) error {
	if len(messages) > 10 {
		return ErrSizeLimit
	}

	req, err := mapToSendMessageInput(p.options.encoder, messages...)
	if err != nil {
		return err
	}

	result, err := p.provider.SendMessageBatch(ctx, req)
	if err != nil {
		return err
	}

	for _, r := range result.Successful {
		logger.Of(ctx).Infof("package=sqs method=Publish message=%+s\n", *r.MessageId)
	}

	var publishErr []error
	for _, r := range result.Failed {
		publishErr = append(publishErr, newPublishMessageError(r))
	}

	return errors.Join(publishErr...)
}

func NewPublisher(ctx context.Context, opts ...PublisherOption) (*Publisher, error) {
	options := newPublisherOption(opts...)

	provider, err := newProvider(ctx)
	if err != nil {
		return nil, err
	}

	return &Publisher{
		options:  options,
		provider: provider,
	}, nil
}
