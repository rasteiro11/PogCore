package rabbitmq

import (
	"time"

	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

const (
	MessageIdempotencyKey = "X-Message-Idempotency-Key"
	MessageDeliveryTag    = "X-Message-Delivery-Tag"
	MessageID             = "X-Message-ID"
	MessageQueue          = "X-Message-Queue"
)

type message struct {
	data      any
	attr      *attributes
	createdAt time.Time
}

type MessageOption func(*message)

func (m *message) Attributes() messaging.Attributes { return m.attr }

func (m *message) Payload() any { return m.data }

func (m *message) ToHash() string { return "" }

func NewMessage(payload any, queue string, opts ...MessageOption) messaging.Message {
	message := &message{
		data:      payload,
		attr:      newAttributes(),
		createdAt: time.Now(),
	}

	for _, opt := range opts {
		opt(message)
	}

	message.attr.Add(MessageQueue, queue)

	return message
}
