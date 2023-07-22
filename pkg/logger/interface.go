package logger

type Logger interface {
	Debug(message string)
	Debugf(template string, args ...any)

	Info(message string)
	Infof(template string, args ...any)

	Warn(message string)
	Warnf(template string, args ...any)

	Error(message string)
	Errorf(template string, args ...any)

	Fatal(message string)
	Fatalf(template string, args ...any)
}

type Provider interface {
	Write(level Level, message string)
	WriteF(level Level, template string, args ...any)
	Enabled(level Level) bool
}
