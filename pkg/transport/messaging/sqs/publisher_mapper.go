package sqs

import (
	"strconv"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/sqs"
	"github.com/aws/aws-sdk-go-v2/service/sqs/types"
	"github.com/google/uuid"
	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

func mapToSendMessageInput(encoder messaging.Encoder, messages ...messaging.Message) (*sqs.SendMessageBatchInput, error) {
	entries, err := mapToSendMessageEntries(encoder, messages...)
	if err != nil {
		return nil, err
	}

	queue := messages[0].Attributes().Get(MessageQueue)

	return &sqs.SendMessageBatchInput{
		Entries:  entries,
		QueueUrl: aws.String(queue),
	}, nil
}

func mapToSendMessageEntries(encoder messaging.Encoder, messages ...messaging.Message) ([]types.SendMessageBatchRequestEntry, error) {
	entries := make([]types.SendMessageBatchRequestEntry, 0, len(messages))
	for _, message := range messages {
		payload, err := encoder(message)
		if err != nil {
			return nil, err
		}

		entry, err := mapToSendMessageEntry(string(payload), message.Attributes())
		if err != nil {
			return nil, err
		}

		entries = append(entries, entry)
	}

	return entries, nil
}

func mapToSendMessageEntry(body string, attributes messaging.Attributes) (types.SendMessageBatchRequestEntry, error) {
	res := types.SendMessageBatchRequestEntry{
		Id:          aws.String(uuid.NewString()),
		MessageBody: aws.String(body),
	}

	if delay, ok := attributes.Lookup(MessageDelaySecond); ok {
		val, err := strconv.Atoi(delay)
		if err != nil {
			return res, err
		}

		res.DelaySeconds = int32(val)
	}

	return res, nil
}
