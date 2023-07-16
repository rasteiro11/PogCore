package server

import "github.com/gofiber/fiber/v2"

type (
	Server interface {
		AddHandler(
			path, group, method string,
			handler fiber.Handler,
			middlewares ...fiber.Handler,
		)
		Use(group string, middlewares ...fiber.Handler)
		Start(port string) error
		PrintRouter()
	}
)
