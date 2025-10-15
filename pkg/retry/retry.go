package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var defaultOptions = Options{
	MaxElapsedTime:  10 * time.Second,
	InitialInterval: 200 * time.Millisecond,
	MaxInterval:     2 * time.Second,
}

type Options struct {
	MaxElapsedTime  time.Duration
	InitialInterval time.Duration
	MaxInterval     time.Duration
}

func Run(ctx context.Context, fn func() error, optFns ...Option) error {
	opts := defaultOptions
	for _, o := range optFns {
		o(&opts)
	}

	b := backoff.NewExponentialBackOff()
	b.InitialInterval = opts.InitialInterval
	b.MaxInterval = opts.MaxInterval
	b.MaxElapsedTime = opts.MaxElapsedTime

	operation := func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
			return fn()
		}
	}

	return backoff.Retry(operation, backoff.WithContext(b, ctx))
}

func RunPermanent(ctx context.Context, fn func() error, isRetryable func(error) bool, optFns ...Option) error {
	return Run(ctx, func() error {
		err := fn()
		if err != nil && !isRetryable(err) {
			return backoff.Permanent(err)
		}
		return err
	}, optFns...)
}
