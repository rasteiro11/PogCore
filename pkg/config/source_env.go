package config

import (
	"os"
	"strings"
)

type envSource struct{}

var _ Source = (*envSource)(nil)

func (s *envSource) Value(key string) (value string, exists bool) {
	return os.LookupEnv(key)
}

func (s *envSource) ForEach(iterator func(key, valeu string)) {
	for _, e := range os.Environ() {
		pair := strings.SplitN(e, "=", 2)
		iterator(pair[0], pair[1])
	}
}

func (s *envSource) Load() error {
	return nil
}

func NewEnvSource() Source {
	return &envSource{}
}
