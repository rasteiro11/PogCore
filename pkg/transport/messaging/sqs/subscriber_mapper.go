package sqs

import (
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

func mapToReceiveMessageInput(queue string) *sqs.ReceiveMessageInput {
	return &sqs.ReceiveMessageInput{
		QueueUrl:            aws.String(queue),
		MaxNumberOfMessages: 1,
		WaitTimeSeconds:     20,
	}
}

func mapToCommitMessageInput(queue string, message messaging.Message) *sqs.DeleteMessageInput {
	return &sqs.DeleteMessageInput{
		QueueUrl: aws.String(queue),
		ReceiptHandle: aws.String(
			message.Attributes().Get(MessageReceiptHandle),
		),
	}
}

func mapProviderToMessage(decoder messaging.Decoder, queue string, m types.Message) (messaging.Message, error) {
	body, err := decoder([]byte(*m.Body))
	if err != nil {
		return nil, err
	}

	res := &message{
		data:      body,
		createdAt: time.Now(),
		attr:      newAttributes(),
	}

	res.Attributes().Add(MessageID, *m.MessageId)
	res.Attributes().Add(MessageIdempotencyKey, m.Attributes["MessageDeduplicationId"])
	res.Attributes().Add(MessageReceiptHandle, *m.ReceiptHandle)
	res.Attributes().Add(MessageQueue, queue)

	return res, nil
}
