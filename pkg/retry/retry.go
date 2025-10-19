package retry

import (
	"context"
	"time"

	"github.com/cenkalti/backoff/v4"
)

var defaultOptions = Options{
	MaxElapsedTime:      10 * time.Second,
	InitialInterval:     200 * time.Millisecond,
	MaxInterval:         2 * time.Second,
	Multiplier:          2.0,
	RandomizationFactor: 0.5,
	MaxRetries:          0,
	OnRetry:             nil,
}

type Options struct {
	MaxElapsedTime      time.Duration
	InitialInterval     time.Duration
	MaxInterval         time.Duration
	Multiplier          float64
	RandomizationFactor float64
	MaxRetries          uint64
	OnRetry             func(err error, next time.Duration)
}

func Run(ctx context.Context, fn func() error, optFns ...Option) error {
	opts := applyOptions(optFns...)
	b := newBackoff(ctx, opts)

	var attempts uint64
	return backoff.Retry(func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
			if err := fn(); err != nil {
				attempts++
				handleRetryCallback(opts, err, b.NextBackOff())
				if opts.MaxRetries > 0 && attempts >= opts.MaxRetries {
					return backoff.Permanent(err)
				}
				return err
			}
			return nil
		}
	}, b)
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

func RunWithResult[T any](ctx context.Context, fn func(ctx context.Context) (T, error), isRetryable func(error) bool, optFns ...Option) (*T, error) {
	opts := applyOptions(optFns...)
	b := newBackoff(ctx, opts)

	var (
		result   T
		attempts uint64
	)

	err := backoff.Retry(func() error {
		select {
		case <-ctx.Done():
			return backoff.Permanent(ctx.Err())
		default:
			value, err := fn(ctx)
			if err != nil {
				attempts++
				handleRetryCallback(opts, err, b.NextBackOff())
				if !isRetryable(err) {
					return backoff.Permanent(err)
				}
				if opts.MaxRetries > 0 && attempts >= opts.MaxRetries {
					return backoff.Permanent(err)
				}
				return err
			}
			result = value
			return nil
		}
	}, b)

	if err != nil {
		return nil, err
	}
	return &result, nil
}

func applyOptions(optFns ...Option) Options {
	opts := defaultOptions
	for _, fn := range optFns {
		fn(&opts)
	}
	return opts
}

func newBackoff(ctx context.Context, opts Options) backoff.BackOffContext {
	b := backoff.NewExponentialBackOff()
	b.InitialInterval = opts.InitialInterval
	b.MaxInterval = opts.MaxInterval
	b.Multiplier = opts.Multiplier
	b.RandomizationFactor = opts.RandomizationFactor
	b.MaxElapsedTime = opts.MaxElapsedTime

	if opts.MaxElapsedTime == 0 {
		b.MaxElapsedTime = backoff.Stop
	}

	b.Reset()
	return backoff.WithContext(b, ctx)
}

func handleRetryCallback(opts Options, err error, next time.Duration) {
	if opts.OnRetry != nil {
		opts.OnRetry(err, next)
	}
}
