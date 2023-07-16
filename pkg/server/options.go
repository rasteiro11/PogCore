package server

func WithPrefix(prefix string) ServerOpt {
	return func(s *server) {
		s.prefix = prefix
	}
}
