package config

import (
	"context"
	"fmt"
	"sync"

	"github.com/rasteiro11/PogCore/pkg/logger"
)

type (
	Option func(*config)

	config struct {
		sources []Source
	}
)

var global *config
var once sync.Once

func Instance() Config {
	if global == nil {
		Init()
	}

	return global
}

func (c *config) get(key string) (string, bool) {
	for _, source := range c.sources {
		v, ok := source.Value(key)
		if ok {
			return v, ok
		}
	}

	return "", false
}

func (c *config) Int(key string) int {
	value, ok := c.get(key)
	if !ok {
		return 0
	}

	return convertToInt(value)
}

func (c *config) Bool(key string) bool {
	value, ok := c.get(key)
	if !ok {
		return ok
	}

	return convertToBool(value)
}

func (c *config) String(key string) string {
	value, ok := c.get(key)
	if ok {
		return value
	}

	return ""
}

func (c *config) ByteSlice(key string) []byte {
	value, ok := c.get(key)
	if !ok {
		return nil
	}

	return []byte(value)
}

func (c *config) RequiredInt(key string) int {
	value, ok := c.get(key)
	if !ok {
		panic(fmt.Sprintf("key not found: %s", key))
	}

	return convertToInt(value)
}

func (c *config) RequiredBool(key string) bool {
	value, ok := c.get(key)
	if !ok {
		panic(fmt.Sprintf("key not found: %s", key))
	}

	return convertToBool(value)
}

func (c *config) RequiredString(key string) string {
	value, ok := c.get(key)
	if !ok {
		panic(fmt.Sprintf("key not found: %s", key))
	}

	return value
}

func (c *config) RequiredByteSlice(key string) []byte {
	value, ok := c.get(key)
	if !ok {
		panic(fmt.Sprintf("key not found: %s", key))
	}

	return []byte(value)
}

func WithSource(source Source) Option {
	return func(c *config) {
		c.sources = append(c.sources, source)
	}
}

func Init(opts ...Option) {
	once.Do(func() {
		cfg := newConfig(opts...)
		cfg.sources = loadSources(cfg.sources)

		global = cfg
	})
}

func newConfig(opts ...Option) *config {
	c := &config{
		sources: []Source{NewEnvSource()},
	}

	for _, opt := range opts {
		opt(c)
	}

	return c
}

func loadSources(sources []Source) []Source {
	for _, source := range sources {
		switch s := source.(type) {
		case LocalSource:
			if err := s.Load(context.Background()); err != nil {
			}
		}
	}

	return sources
}
