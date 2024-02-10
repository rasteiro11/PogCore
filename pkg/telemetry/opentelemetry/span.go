package opentelemetry

import (
	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/telemetry/tracer"
	"go.opentelemetry.io/otel/codes"
	"go.opentelemetry.io/otel/trace"
)

type span struct {
	traceSpan trace.Span
	tracer    *Tracer
}

var _ tracer.Span = (*span)(nil)

func (s *span) End() {
	s.traceSpan.End()
}

func (s *span) SetAttribute(key string, value interface{}) {
	if !s.traceSpan.IsRecording() {
		return
	}

	attr, err := parseAttributeKeyValue(key, value)
	if err != nil {
		logger.Global().Warnf("package=opentelemetry operation=ExtractAttributeValue error=%+v\n", err)
		return
	}

	s.traceSpan.SetAttributes(attr)
}

func (s *span) SetStatus(status tracer.Status) {
	s.traceSpan.SetStatus(codes.Code(status.Code), status.Msg)
}

func (s *span) RecordError(err error) {
	s.traceSpan.SetStatus(codes.Error, err.Error())
	s.traceSpan.RecordError(err, trace.WithStackTrace(true))
}

func (s *span) Tracer() tracer.Tracer { return s.tracer }
