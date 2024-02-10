package tracer

import (
	"context"
)

type Status struct {
	Code int64
	Msg  string
}

// Span represents a single operation within a trace.
type Span interface {
	// End ends the span.
	End()

	// SetAttribute sets the specified attribute with the given value.
	SetAttribute(key string, value interface{})

	// SetStatus sets the span's status.
	SetStatus(status Status)

	// RecordError records the given error on the span.
	RecordError(err error)

	// Tracer returns the Tracer that created this span.
	Tracer() Tracer
}

type Propagator interface {
	Inject(ctx context.Context, values map[string][]string)
	Extract(ctx context.Context, values map[string][]string) context.Context
}

// Tracer is the interface for creating and recording traces.
type Tracer interface {
	// Start creates a new span with the given context and name.
	Start(ctx context.Context, name string) (context.Context, Span)

	Shutdown(ctx context.Context) error

	Propagator
}
