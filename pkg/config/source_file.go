package config

import (
	"bufio"
	"context"
	"errors"
	"os"
	"strings"
)

type FileSourceOption func(*fileSource)

type fileSource struct {
	path string
	env  map[string]string
}

var (
	ErrEnvFile         = errors.New("error env file contains bad key value pairs")
	ErrFileNotProvided = errors.New("error file not provided")
)

func WithFileSource(path string) FileSourceOption {
	return func(fs *fileSource) {
		fs.path = path
	}
}

var _ Source = (*fileSource)(nil)

func (s *fileSource) Value(key string) (value string, exists bool) {
	v, ok := s.env[key]
	if ok {
		return v, ok
	}

	return "", false
}

func (s *fileSource) ForEach(iterator func(key, valeu string)) {
	for k, v := range s.env {
		iterator(k, v)
	}
}

func (s *fileSource) Load(ctx context.Context) error {
	if s.path == "" {
		return ErrFileNotProvided
	}

	file, err := os.Open(s.path)
	if err != nil {
		return err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		pair := strings.SplitN(scanner.Text(), "=", 2)

		if len(pair) != 2 {
			return ErrEnvFile
		}

		s.env[pair[0]] = pair[1]
	}

	return nil
}

func NewFileSource(opts ...FileSourceOption) Source {
	s := &fileSource{
		env: map[string]string{},
	}

	for _, opt := range opts {
		opt(s)
	}

	return s
}
