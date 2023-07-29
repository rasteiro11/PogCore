package kafka

import (
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

func kafkaMessageMapper(encoder messaging.Encoder, msg messaging.Message) (*kafka.Message, error) {
	topic := msg.Attributes().Get(MessageTopic)

	value, err := encoder(msg)
	if err != nil {
		logger.Global().Errorf("[messaging.kafka.publisher_mapper.kafkaMessageMapper] encoder() returned error: %+v\n", err)
		return nil, err
	}

	kafkaMessage := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
		Value:          value,
	}

	return kafkaMessage, nil
}
