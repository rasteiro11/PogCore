package logger

import (
	"fmt"
	"log"
)

type std struct {
	level        Level
	colorPallete map[Level]string
}

var (
	colorReset  = "\033[0m"
	colorRed    = "\033[91m"
	colorGreen  = "\033[92m"
	colorYellow = "\033[93m"
	colorBlue   = "\033[94m"
	colorPurple = "\033[95m"
	colorCyan   = "\033[96m"
)

type StdProviderLoggerOption func(*std)

var _ Provider = (*std)(nil)

func defaultProvider() Provider {
	return &std{
		level:        DefaultLevel(),
		colorPallete: defaultColorPalette(),
	}
}

func defaultColorPalette() map[Level]string {
	return map[Level]string{
		DebugLevel: colorGreen,
		InfoLevel:  colorCyan,
		WarnLevel:  colorYellow,
		ErrorLevel: colorRed,
		FatalLevel: colorPurple,
	}
}

func (s *std) color(l Level) string {
	color, ok := s.colorPallete[l]
	if !ok {
		return colorReset
	}

	return color
}

func WithLevelColoring(level Level, color string) StdProviderLoggerOption {
	return func(s *std) {
		s.colorPallete[level] = color
	}
}

func (s *std) Write(level Level, message string) {
	log.Printf("[%s%s%s] - %s", s.color(level), level.String(), colorReset, message)
}

func (s *std) WriteF(level Level, template string, args ...any) {
	msg := fmt.Sprintf("[%s%s%s] - %s", s.color(level), level.String(), colorReset, template)
	if len(args) > 0 {
		log.Printf(msg, args...)
		return
	}

	log.Print(msg)
}

func (s *std) Enabled(level Level) bool {
	return level >= s.level
}

func NewStdLogger(opts ...StdProviderLoggerOption) Provider {
	p := &std{
		level:        DefaultLevel(),
		colorPallete: defaultColorPalette(),
	}

	for _, opt := range opts {
		opt(p)
	}

	return p
}
