package rabbitmq

import (
	"strconv"

	"github.com/rabbitmq/amqp091-go"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

func mapToAMQPPublishing(encoder messaging.Encoder, message messaging.Message) (string, amqp091.Publishing, error) {
	payload, err := encoder(message)
	if err != nil {
		return "", amqp091.Publishing{}, err
	}

	headers := amqp091.Table{}
	for k, v := range message.Attributes().Values() {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        payload,
		Headers:     headers,
	}

	queue := message.Attributes().Get(MessageQueue)

	return queue, pub, nil
}

// mapPublishingOnly returns an amqp publishing with headers set from message attributes.
func mapPublishingOnly(encoder messaging.Encoder, message messaging.Message) (amqp091.Publishing, error) {
	payload, err := encoder(message)
	if err != nil {
		return amqp091.Publishing{}, err
	}

	headers := amqp091.Table{}
	for k, v := range message.Attributes().Values() {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	pub := amqp091.Publishing{
		ContentType: "application/json",
		Body:        payload,
		Headers:     headers,
	}

	return pub, nil
}

func mapDeliveryToMessage(decoder messaging.Decoder, queue string, d amqp091.Delivery) (messaging.Message, error) {
	body, err := decoder(d.Body)
	if err != nil {
		return nil, err
	}

	res := &message{
		data:      body,
		createdAt: d.Timestamp,
		attr:      newAttributes(),
	}

	if d.MessageId != "" {
		res.Attributes().Add(MessageID, d.MessageId)
	}

	res.Attributes().Add(MessageDeliveryTag, strconv.FormatUint(d.DeliveryTag, 10))
	res.Attributes().Add(MessageQueue, queue)

	for k, v := range d.Headers {
		switch val := v.(type) {
		case string:
			res.Attributes().Add(k, val)
		}
	}

	return res, nil
}
