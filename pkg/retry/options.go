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
