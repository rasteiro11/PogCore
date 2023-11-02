package gorm

import (
	"context"
	"fmt"

	"github.com/rasteiro11/PogCore/pkg/logger"
	"github.com/rasteiro11/PogCore/pkg/repository"
	"gorm.io/gorm"
)

type Option[T any] func(*repo[T])

type repo[T any] struct {
	db *gorm.DB
}

func WithDatabase[T any](db *gorm.DB) Option[T] {
	return func(r *repo[T]) {
		r.db = db
	}
}

func NewGormRepo[T any](otps ...Option[T]) repository.Repository[T] {
	r := &repo[T]{}

	for _, opt := range otps {
		opt(r)
	}

	return r
}

func (r *repo[T]) Create(ctx context.Context, model *T) (*T, error) {
	if err := r.db.Create(model).Error; err != nil {
		logger.Of(ctx).Errorf("[repository.gorm.Cerate] db.Create() returned error: %+v\n", err)
		return nil, err
	}

	return model, nil
}

func (r *repo[T]) FindOne(ctx context.Context, model *T) (*T, error) {
	var m T
	if err := r.db.Where(model).Take(&m).Error; err != nil {
		logger.Of(ctx).Errorf("[repository.gorm.FindOne] db.Take() returned error: %+v\n", err)
		return nil, err
	}

	return &m, nil
}

func (r *repo[T]) Update(ctx context.Context, model *T) (*T, error) {
	if err := r.db.Debug().Model(model).Updates(&model).Error; err != nil {
		logger.Of(ctx).Errorf("[repository.gorm.Update] db.Take() returned error: %+v\n", err)
		return nil, err
	}

	return model, nil
}

func (r *repo[T]) FindAllWithPages(ctx context.Context, query *repository.PagesQuery[T]) (*repository.PagesResponse[T], error) {
	var data []*T
	var modelQuery T
	var count int64
	dbInstance := r.db.Debug().Model(&modelQuery).Where(&query.Query)

	if query.OrderBy != nil {
		if query.OrderBy.Field != "" && query.OrderBy.Direction != "" {
			dbInstance = dbInstance.Order(fmt.Sprintf("%s %s", query.OrderBy.Field, query.OrderBy.Direction))
		}
	}

	dbInstance.Count(&count)
	offset := (query.Page - 1) * query.PerPage
	pages := count / int64(query.PerPage)
	if query.PerPage > int(count) {
		pages = 1
	}

	if err := dbInstance.Order("id desc").Limit(query.PerPage).Offset(offset).Find(&data).Error; err != nil {
		return nil, err
	}

	return &repository.PagesResponse[T]{
		Data:  data,
		Total: int(count),
		Page:  query.Page,
		Pages: int(pages),
	}, nil
}

// var _ Repository[models.Transaction] = (*repository)(nil)
