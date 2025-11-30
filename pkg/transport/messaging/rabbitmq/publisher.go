package rabbitmq

import (
	"context"
	"errors"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type Publisher struct {
	options *publisherOptions
	conn    *amqp.Connection
	ch      *amqp.Channel
}

var _ messaging.Publisher = (*Publisher)(nil)

func (p *Publisher) Publish(ctx context.Context, messages ...messaging.Message) error {
	if len(messages) == 0 {
		return nil
	}

	if p.ch == nil {
		return errors.New("rabbitmq channel not initialized")
	}

	for _, m := range messages {
		queue, pub, err := mapToAMQPPublishing(p.options.encoder, m)
		if err != nil {
			return err
		}

		_, err = p.ch.QueueDeclare(
			queue,
			true,
			false,
			false,
			false,
			nil,
		)
		if err != nil {
			return err
		}

		if err := p.ch.PublishWithContext(ctx, "", queue, false, false, pub); err != nil {
			return err
		}

		logger.Of(ctx).Infof("package=rabbitmq method=Publish message=sent queue=%s", queue)
	}

	return nil
}

func (p *Publisher) PublishToQueue(ctx context.Context, queue string, messages ...messaging.Message) error {
	if len(messages) == 0 {
		return nil
	}

	if p.ch == nil {
		return errors.New("rabbitmq channel not initialized")
	}

	_, err := p.ch.QueueDeclare(
		queue,
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	for _, m := range messages {
		pub, err := mapPublishingOnly(p.options.encoder, m)
		if err != nil {
			return err
		}

		if err := p.ch.PublishWithContext(ctx, "", queue, false, false, pub); err != nil {
			return err
		}

		logger.Of(ctx).Infof("package=rabbitmq method=PublishToQueue message=sent queue=%s", queue)
	}

	return nil
}

func (p *Publisher) PublishToTopic(ctx context.Context, exchange, routingKey string, messages ...messaging.Message) error {
	if len(messages) == 0 {
		return nil
	}

	if p.ch == nil {
		return errors.New("rabbitmq channel not initialized")
	}

	if err := p.ch.ExchangeDeclare(
		exchange,
		"topic",
		true,
		false,
		false,
		false,
		nil,
	); err != nil {
		return err
	}

	for _, m := range messages {
		pub, err := mapPublishingOnly(p.options.encoder, m)
		if err != nil {
			return err
		}

		if err := p.ch.PublishWithContext(ctx, exchange, routingKey, false, false, pub); err != nil {
			return err
		}

		logger.Of(ctx).Infof("package=rabbitmq method=PublishToTopic message=sent exchange=%s routingKey=%s", exchange, routingKey)
	}

	return nil
}

func NewPublisher(ctx context.Context, opts ...PublisherOption) (*Publisher, error) {
	options := newPublisherOption(opts...)

	conn, err := newProvider(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Publisher{
		options: options,
		conn:    conn,
		ch:      ch,
	}, nil
}
