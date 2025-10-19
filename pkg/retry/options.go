package retry

import "time"

type Option func(*Options)

func WithMaxElapsedTime(d time.Duration) Option {
	return func(o *Options) { o.MaxElapsedTime = d }
}

func WithInitialInterval(d time.Duration) Option {
	return func(o *Options) { o.InitialInterval = d }
}

func WithMaxInterval(d time.Duration) Option {
	return func(o *Options) { o.MaxInterval = d }
}

func WithMultiplier(m float64) Option {
	return func(o *Options) { o.Multiplier = m }
}

func WithRandomizationFactor(f float64) Option {
	return func(o *Options) { o.RandomizationFactor = f }
}

func WithMaxRetries(n uint64) Option {
	return func(o *Options) { o.MaxRetries = n }
}

func WithOnRetry(fn func(error, time.Duration)) Option {
	return func(o *Options) { o.OnRetry = fn }
}
