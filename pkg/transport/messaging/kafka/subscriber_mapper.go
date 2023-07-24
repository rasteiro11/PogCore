package kafka

import (
	"strconv"
	"time"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

func mapProviderToMessage(decoder messaging.Decoder, queue string, m *kafka.Message) (messaging.Message, error) {
	body, err := decoder(m.Value)
	if err != nil {
		return nil, err
	}

	res := &message{
		msg:       m,
		data:      body,
		createdAt: time.Now(),
		attr:      newAttributes(),
	}

	res.Attributes().Add(MessageTopic, *m.TopicPartition.Topic)
	res.Attributes().Add(MessagePartition, strconv.Itoa(int(m.TopicPartition.Partition)))
	res.Attributes().Add(MessageTimestamp, m.Timestamp.String())

	return res, nil
}
