package sqs

import (
	"encoding/json"

	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type publisherOption struct {
	encoder messaging.Encoder
}

type PublisherOption func(*publisherOption)

func newPublisherOption(opts ...PublisherOption) *publisherOption {
	options := &publisherOption{
		encoder: func(m messaging.Message) ([]byte, error) {
			return json.Marshal(m.Payload())
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
