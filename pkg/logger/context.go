package logger

import "context"

type loggerContextKey struct{}

var loggerContext = &loggerContextKey{}

func Of(ctx context.Context) Logger {
	if l := fromContext(ctx); l != nil {
		return l
	}

	return Global()
}

func fromContext(ctx context.Context) Logger {
	if l, ok := ctx.Value(loggerContext).(Logger); ok {
		return l
	}

	return nil
}
