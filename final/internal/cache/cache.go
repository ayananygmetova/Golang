package cache

import (
	"context"
	"final/internal/models"
)

type Cache interface {
	Close() error

	Categories() CategoriesCacheRepo
	Products() ProductsCacheRepo

	DeleteAll(ctx context.Context) error
}

type CategoriesCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Category) error
	Get(ctx context.Context, key string) ([]*models.Category, error)
}
type ProductsCacheRepo interface {
	Set(ctx context.Context, key string, value []*models.Product) error
	Get(ctx context.Context, key string) ([]*models.Product, error)
}
