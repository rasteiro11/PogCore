package config

import "context"

type (
	Config interface {
		Int(key string) int
		Bool(key string) bool
		String(key string) string
		ByteSlice(key string) []byte

		RequiredInt(key string) int
		RequiredBool(key string) bool
		RequiredString(key string) string
		RequiredByteSlice(key string) []byte
	}

	Source interface {
		Value(key string) (string, bool)
		ForEach(iterator func(key, value string))
	}

	LocalSource interface {
		Source
		Load(ctx context.Context) error
	}
)
