package server_test

import (
	"flashcards/core/server"
	"fmt"
	"net/http"
	"testing"

	"github.com/gofiber/fiber/v2"
)

func createServer(t *testing.T) server.Server {
	s := server.NewServer(server.WithPrefix("/gamer"))
	if s == nil {
		t.Fatalf("[createServer] server was not created")
	}
	return s
}

func TestNewServer(t *testing.T) {
	createServer(t)
}

func TestAddSimpleHandler(t *testing.T) {
	s := createServer(t)
	s.AddHandler("/list_gamers", "/gamers", http.MethodPost, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"simple": "return"})
	})
}

func TestAddHandlerWithMIddleware(t *testing.T) {
	s := createServer(t)
	s.AddHandler("/list_gamers", "/gamers", http.MethodPost, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"simple": "return"})
	},
		func(c *fiber.Ctx) error {
			fmt.Println("PASSING MIDDLEWARE 2")
			return c.Next()
		},
	)
}

func TestAddHandlerAndMiddleware(t *testing.T) {
	s := createServer(t)
	s.Use("/gamers", func(c *fiber.Ctx) error {
		fmt.Println("PASSING MIDDLEWARE")
		return c.Next()
	},
		func(c *fiber.Ctx) error {
			fmt.Println("PASSING MIDDLEWARE 2")
			return c.Next()
		},
	)
	s.AddHandler("/list_gamers", "/gamers", http.MethodPost, func(c *fiber.Ctx) error {
		return c.JSON(fiber.Map{"simple": "return"})
	},
		func(c *fiber.Ctx) error {
			fmt.Println("PASSING MIDDLEWARE 3")
			return c.Next()
		},
	)
}

func TestStartServer(t *testing.T) {
	s := createServer(t)
	go func() {
		s.Start(":8080")
	}()
}
