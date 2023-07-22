package logger

import (
	"os"
)

type (
	Option func(*logger)

	logger struct {
		provider Provider
	}
)

var global Logger

func Use(logger Logger) {
	if logger == nil {
		global = New()
		return
	}

	global = logger
}

func Global() Logger {
	if global == nil {
		return New()
	}

	return global
}

func WithProvider(provider Provider) Option {
	return func(l *logger) {
		l.provider = provider
	}
}

func New(opts ...Option) Logger {
	l := &logger{
		provider: defaultProvider(),
	}

	for _, opt := range opts {
		opt(l)
	}

	return l
}

func (l *logger) Debug(message string) {
	if l.provider.Enabled(DebugLevel) {
		l.provider.Write(DebugLevel, message)
	}
}

func (l *logger) Debugf(template string, args ...any) {
	if l.provider.Enabled(DebugLevel) {
		l.provider.WriteF(DebugLevel, template, args...)
	}
}

func (l *logger) Info(message string) {
	if l.provider.Enabled(InfoLevel) {
		l.provider.Write(InfoLevel, message)
	}
}

func (l *logger) Infof(template string, args ...any) {
	if l.provider.Enabled(InfoLevel) {
		l.provider.WriteF(InfoLevel, template, args...)
	}
}

func (l *logger) Warn(message string) {
	if l.provider.Enabled(WarnLevel) {
		l.provider.Write(WarnLevel, message)
	}
}

func (l *logger) Warnf(template string, args ...any) {
	if l.provider.Enabled(WarnLevel) {
		l.provider.WriteF(WarnLevel, template, args...)
	}
}

func (l *logger) Error(message string) {
	if l.provider.Enabled(ErrorLevel) {
		l.provider.Write(ErrorLevel, message)
	}
}

func (l *logger) Errorf(template string, args ...any) {
	if l.provider.Enabled(ErrorLevel) {
		l.provider.WriteF(ErrorLevel, template, args...)
	}
}

func (l *logger) Fatal(message string) {
	if l.provider.Enabled(FatalLevel) {
		l.provider.Write(FatalLevel, message)
		os.Exit(1)
	}
}

func (l *logger) Fatalf(template string, args ...any) {
	if l.provider.Enabled(FatalLevel) {
		l.provider.WriteF(FatalLevel, template, args...)
		os.Exit(1)
	}
}
