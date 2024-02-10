package tracer

import "context"

// noopTracer is a no-operation implementation of the Tracer interface.
type noopTracer struct {
}

// Start creates a new span with the given context and name, but does nothing.
func (n *noopTracer) Start(ctx context.Context, name string) (context.Context, Span) {
	return ctx, &noopSpan{}
}

// WithSpan creates a new context with the given span, but does nothing.
func (n *noopTracer) WithSpan(ctx context.Context, span Span) context.Context {
	return ctx
}

func (n *noopTracer) Inject(ctx context.Context, values map[string][]string) {
}

func (n *noopTracer) Extract(ctx context.Context, values map[string][]string) context.Context {
	return ctx
}

func (n *noopTracer) Shutdown(ctx context.Context) error { return nil }

// noopSpan is a no-operation implementation of the Span interface.
type noopSpan struct{}

// End ends the span, but does nothing.
func (n *noopSpan) End() {}

// SetAttribute sets the specified attribute with the given value, but does nothing.
func (n *noopSpan) SetAttribute(key string, value interface{}) {}

// SetStatus sets the span's status, but does nothing.
func (n *noopSpan) SetStatus(status Status) {}

// RecordError records the given error on the span, but does nothing.
func (n *noopSpan) RecordError(err error) {}

// Tracer returns the Tracer that created this span, but returns a NoOpTracer.
func (n *noopSpan) Tracer() Tracer {
	return &noopTracer{}
}
