package rest

import (
	"net/http"

	"github.com/gofiber/fiber/v2"
)

type (
	response[T any] struct {
		statusCode int
		body       T
		headers    http.Header
	}
)

func configureResponse[T any](c *fiber.Ctx, statusCode int, opts ...ResponseOpt[T]) *response[T] {
	res := &response[T]{
		statusCode: statusCode,
		headers:    make(http.Header),
	}

	for _, opt := range opts {
		opt(res)
	}

	return res
}

type ResponseOpt[T any] func(r *response[T])

func NewResponse[T any](c *fiber.Ctx, statusCode int, opts ...ResponseOpt[T]) *response[T] {
	response := configureResponse(c, statusCode, opts...)

	if response.statusCode != 0 {
		c.Response().SetStatusCode(response.statusCode)
	}

	return response
}

func (res *response[T]) JSON(c *fiber.Ctx) error {
	return c.JSON(res.body)
}

func NewStatusOk[T any](c *fiber.Ctx, opts ...ResponseOpt[T]) error {
	return NewResponse(c, http.StatusOK, opts...).JSON(c)
}

func NewStatusCreated[T any](c *fiber.Ctx, opts ...ResponseOpt[T]) error {
	return NewResponse(c, http.StatusCreated, opts...).JSON(c)
}

func NewStatusBadRequest(c *fiber.Ctx, err error) error {
	return NewResponse(c, http.StatusBadRequest, WithBody(err.Error())).JSON(c)
}

func NewStatusInternalServerError(c *fiber.Ctx, err error) error {
	return NewResponse(c, http.StatusInternalServerError, WithBody(err.Error())).JSON(c)
}

func NewStatusUnauthorized(c *fiber.Ctx, err error) error {
	return NewResponse(c, http.StatusUnauthorized, WithBody(err.Error())).JSON(c)
}

func NewStatusNotFound(c *fiber.Ctx, err error) error {
	return NewResponse(c, http.StatusNotFound, WithBody(err.Error())).JSON(c)
}
