package store

import (
	"context"
	"hw8/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Categories() CategoriesRepository
	Products() ProductsRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type ProductsRepository interface {
	Create(ctx context.Context, product *models.Product) error
	All(ctx context.Context) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	ByCategory(ctx context.Context, id int) ([]*models.Product, error)
}
