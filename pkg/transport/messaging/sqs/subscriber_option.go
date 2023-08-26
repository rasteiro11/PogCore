package sqs

import (
	"encoding/json"
	"reflect"

	"github.com/rasteiro11/PogCore/pkg/transport/messaging"
)

type subscriberOptions struct {
	decoder messaging.Decoder
}

type SubscriberOption func(*subscriberOptions)

func WithDecoder(decoder messaging.Decoder) SubscriberOption {
	return func(so *subscriberOptions) {
		so.decoder = decoder
	}
}

func WithDecoderTarget(typeof any) SubscriberOption {
	return func(s *subscriberOptions) {
		v := reflect.TypeOf(typeof)

		s.decoder = func(b []byte) (any, error) {
			target := reflect.New(v).Interface()
			if err := json.Unmarshal(b, target); err != nil {
				return nil, err
			}

			return target, nil
		}
	}
}

func newSubscriberOption(opts ...SubscriberOption) *subscriberOptions {
	options := &subscriberOptions{
		decoder: func(b []byte) (any, error) {
			var target any

			if err := json.Unmarshal(b, &target); err != nil {
				return nil, err
			}

			return target, nil
		},
	}

	for _, opt := range opts {
		opt(options)
	}

	return options
}
