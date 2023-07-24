package kafka

import (
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

const (
	MessageTopic     = "X-Message-Topic"
	MessagePartition = "X-Message-Partition"
	MessageTimestamp = "X-Message-Timestamp"
)

type message struct {
	msg       *kafka.Message
	data      any
	attr      *attributes
	createdAt time.Time
}

type MessageOption func(*message)

func (m *message) Attributes() messaging.Attributes { return m.attr }

func (m *message) Payload() any {
	return m.data
}

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

	message.attr.Add(MessageTopic, queue)

	return message
}
