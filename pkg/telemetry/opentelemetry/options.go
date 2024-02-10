package opentelemetry

import (
	"context"

	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	"go.opentelemetry.io/otel/sdk/trace"
)

type options struct {
	spanProcessor trace.SpanProcessor
	traceSampler  trace.Sampler
	propagator    propagation.TextMapPropagator
	resource      *resource.Resource
}

func (o *options) toTraceProviderOption() []trace.TracerProviderOption {
	opts := make([]trace.TracerProviderOption, 0)

	if o.resource != nil {
		opts = append(opts, trace.WithResource(o.resource))
	}

	if o.traceSampler != nil {
		opts = append(opts, trace.WithSampler(o.traceSampler))
	}

	if o.spanProcessor != nil {
		opts = append(opts, trace.WithSpanProcessor(o.spanProcessor))
	}

	return opts
}

type Option func(*options)

func NewOptions(ctx context.Context, opts ...Option) (*options, error) {
	options := &options{
		traceSampler:  trace.AlwaysSample(),
		propagator:    propagation.TraceContext{},
		spanProcessor: NewProcessorFromEnv(),
		resource:      NewResourceFromEnv(),
	}

	for _, opt := range opts {
		opt(options)
	}

	return options, nil
}
