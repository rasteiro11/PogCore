package grpcserver

import (
	"github.com/rasteiro11/PogCore/pkg/config"
	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
	"net"
)

type Server interface {
	Run() error
	Register(desc grpc.ServiceDesc, service any)
	Close()
}

type (
	Option func(*server)

	server struct {
		port       string
		useReflect bool
		srv        *grpc.Server
		opts       []grpc.ServerOption
	}
)

func WithReflectionEnabled() Option {
	return func(s *server) {
		s.useReflect = true
	}
}

func WithServerPort(port string) Option {
	return func(s *server) {
		s.port = port
	}
}

func WithServerOption(opts ...grpc.ServerOption) Option {
	return func(s *server) {
		s.opts = append(s.opts, opts...)
	}
}

func NewServer(opts ...Option) Server {
	s := &server{
		port: config.Instance().String("SERVICE_PORT"),
	}

	for _, opt := range opts {
		opt(s)
	}

	server := grpc.NewServer()
	if s.useReflect {
		reflection.Register(server)
	}

	s.srv = server

	return s
}

func (s *server) Run() error {
	lis, err := net.Listen("tcp", ":"+s.port)
	if err != nil {
		return err
	}

	return s.srv.Serve(lis)
}

func (s *server) Close() {
	s.srv.GracefulStop()
}

func (s *server) Register(desc grpc.ServiceDesc, service any) {
	s.srv.RegisterService(&desc, service)
}
