package store

import (
	"Midterm/internal/models"
	"context"
)

type Store interface {
	Categories() CategoriesRepository
	Products() ProductsRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, laptop *models.Category) (*models.Category, error)
	All(ctx context.Context) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, laptop *models.Category) (*models.Category, error)
	Delete(ctx context.Context, id int) error
}

type ProductsRepository interface {
	Create(ctx context.Context, phone *models.Product) (*models.Product, error)
	All(ctx context.Context) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, laptop *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id int) error
	ByCategory(ctx context.Context, id int) ([]*models.Product, error)
}
