package rabbitmq

import (
	"context"
	"strconv"

	amqp "github.com/rabbitmq/amqp091-go"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type Subscriber struct {
	options    *subscriberOptions
	conn       *amqp.Connection
	ch         *amqp.Channel
	queue      string
	deliveries <-chan amqp.Delivery
}

var _ messaging.Subscriber = (*Subscriber)(nil)

func (s *Subscriber) Subscribe(ctx context.Context) (messaging.Message, error) {
	if s.ch == nil {
		return nil, ErrProviderNotInitialized
	}

	if s.deliveries != nil {
		select {
		case <-ctx.Done():
			return nil, ctx.Err()
		case d, ok := <-s.deliveries:
			if !ok {
				return nil, ErrNoMessageInQueue
			}
			return mapDeliveryToMessage(s.options.decoder, s.queue, d)
		}
	}

	d, ok, err := s.ch.Get(s.queue, false)
	if err != nil {
		return nil, err
	}

	if !ok {
		return nil, ErrNoMessageInQueue
	}

	return mapDeliveryToMessage(s.options.decoder, s.queue, d)
}

func (s *Subscriber) Commit(ctx context.Context, message messaging.Message) error {
	select {
	case <-ctx.Done():
		return ctx.Err()
	default:
	}

	tagStr := message.Attributes().Get(MessageDeliveryTag)
	if tagStr == "" {
		return ErrInvalidDeliveryTag
	}

	tag, err := strconv.ParseUint(tagStr, 10, 64)
	if err != nil {
		return err
	}

	if err := s.ch.Ack(tag, false); err != nil {
		return err
	}

	logger.Of(ctx).Infof("package=rabbitmq method=Commit message=%s\n", "success commit")

	return nil
}

func (s *Subscriber) Close(ctx context.Context) error {
	if s.ch != nil {
		_ = s.ch.Close()
	}
	if s.conn != nil {
		_ = s.conn.Close()
	}
	return nil
}

func NewSubscriber(ctx context.Context, queue string, opts ...SubscriberOption) (*Subscriber, error) {
	options := newSubscriberOption(opts...)

	conn, err := newProvider(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		options: options,
		conn:    conn,
		ch:      ch,
		queue:   queue,
	}, nil
}

func NewTopicSubscriber(ctx context.Context, exchange, routingKey string, opts ...SubscriberOption) (*Subscriber, error) {
	options := newSubscriberOption(opts...)

	conn, err := newProvider(ctx)
	if err != nil {
		return nil, err
	}

	ch, err := conn.Channel()
	if err != nil {
		return nil, err
	}

	if err := ch.ExchangeDeclare(exchange, "topic", true, false, false, false, nil); err != nil {
		return nil, err
	}

	q, err := ch.QueueDeclare("", false, true, true, false, nil)
	if err != nil {
		return nil, err
	}

	if err := ch.QueueBind(q.Name, routingKey, exchange, false, nil); err != nil {
		return nil, err
	}

	deliveries, err := ch.Consume(q.Name, "", false, false, false, false, nil)
	if err != nil {
		return nil, err
	}

	return &Subscriber{
		options:    options,
		conn:       conn,
		ch:         ch,
		queue:      q.Name,
		deliveries: deliveries,
	}, nil
}
