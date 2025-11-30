package rabbitmq

import (
	"encoding/json"

	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type publisherOptions struct {
	encoder messaging.Encoder
}

type PublisherOption func(*publisherOptions)

func WithEncoder(encoder messaging.Encoder) PublisherOption {
	return func(po *publisherOptions) {
		po.encoder = encoder
	}
}

func WithEncoderDefaultJSON() PublisherOption {
	return func(po *publisherOptions) {
		po.encoder = func(m messaging.Message) ([]byte, error) {
			return json.Marshal(m.Payload())
		}
	}
}

func newPublisherOption(opts ...PublisherOption) *publisherOptions {
	options := &publisherOptions{}

	options.encoder = func(m messaging.Message) ([]byte, error) {
		return json.Marshal(m.Payload())
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
