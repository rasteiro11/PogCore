package repository

import (
	"context"
)

type Repository[T any] interface {
	FindOne(ctx context.Context, model *T) (*T, error)
	Create(ctx context.Context, model *T) (*T, error)
	Update(ctx context.Context, model *T) (*T, error)
	FindAllWithPages(ctx context.Context, q *PagesQuery[T]) (*PagesResponse[T], error)
}

type PagesQuery[T any] struct {
	PerPage int      `json:"per_page"`
	Page    int      `json:"page"`
	Query   *T       `json:"query"`
	OrderBy *OrderBy `json:"order_by"`
}

type OrderBy struct {
	Field     string `json:"field"`
	Direction string `json:"direction"`
}

type PagesResponse[T any] struct {
	Data  []*T `json:"data"`
	Page  int  `json:"page"`
	Pages int  `json:"pages"`
	Total int  `json:"total"`
}
