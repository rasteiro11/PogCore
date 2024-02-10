package opentelemetry

import (
	"context"
	"os"

	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/trace"
)

const traceName = "testing/TraceName"

type Tracer struct {
	provider *trace.TracerProvider
}

var _ tracer.Tracer = (*Tracer)(nil)

func (t *Tracer) Start(ctx context.Context, name string) (context.Context, tracer.Span) {
	ctx, traceSpan := t.provider.Tracer(traceName).Start(ctx, name)

	traceSpan.SetAttributes(
		attribute.KeyValue{
			Key:   "service.name",
			Value: attribute.StringValue(os.Getenv("OTEL_SERVICE_NAME")),
		},
		attribute.KeyValue{
			Key:   "deployment.environment",
			Value: attribute.StringValue(os.Getenv("ENV")),
		},
	)

	span := &span{
		traceSpan: traceSpan,
		tracer:    t,
	}

	return ctx, span
}

func (t *Tracer) Inject(ctx context.Context, values map[string][]string) {
	props := otel.GetTextMapPropagator()
	props.Inject(ctx, propagation.HeaderCarrier(values))
}

func (t *Tracer) Extract(ctx context.Context, values map[string][]string) context.Context {
	props := otel.GetTextMapPropagator()
	return props.Extract(ctx, propagation.HeaderCarrier(values))
}

func (t *Tracer) Shutdown(ctx context.Context) error {
	if err := t.provider.ForceFlush(ctx); err != nil {
		return err
	}

	return t.provider.Shutdown(ctx)
}

func NewProvider(ctx context.Context, opts ...Option) tracer.Tracer {
	options, err := NewOptions(ctx, opts...)
	if err != nil {
		return nil
	}

	providerOption := options.toTraceProviderOption()
	provider := trace.NewTracerProvider(providerOption...)
	if provider == nil {
		return nil
	}

	otel.SetTracerProvider(provider)
	otel.SetTextMapPropagator(options.propagator)

	return &Tracer{provider: provider}
}
