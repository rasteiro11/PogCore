package server

import (
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
)

type (
	ServerOpt func(*server)

	route struct {
		method, fullPath string
	}

	server struct {
		engine *fiber.App
		prefix string
		routes map[string]route
	}
)

func (s *server) Use(group string, middlewares ...fiber.Handler) {
	fullPath := s.prefix + group

	for _, middleware := range middlewares {
		s.engine.Use(fullPath, middleware)
	}
}

func (s *server) AddHandler(path, group, method string, handler fiber.Handler, middlewares ...fiber.Handler) {
	fullPath := s.prefix + group + path

	_, ok := s.routes[fullPath+method]
	if ok {
		log.Println("[server.AddHandler] there is already a handler for this path")
		return
	}

	s.routes[fullPath+method] = route{
		method:   method,
		fullPath: fullPath,
	}

	for _, middleware := range middlewares {
		s.engine.Use(fullPath, middleware)
	}
	s.engine.Add(method, fullPath, handler)
}

func (s *server) PrintRouter() {
	for _, route := range s.routes {
		log.Printf("[METHOD] %s - [PATH] %s", route.method, route.fullPath)
	}
}

func (s *server) Start(port string) error {
	return s.engine.Listen(port)
}

func NewServer(opts ...ServerOpt) Server {
	s := &server{
		routes: make(map[string]route),
	}

	s.engine = fiber.New()

	s.engine.Use(cors.New(cors.Config{
		AllowHeaders:     "*",
		AllowOrigins:     "*",
		AllowCredentials: true,
		AllowMethods:     "*",
	}))

	for _, opt := range opts {
		opt(s)
	}

	return s
}
