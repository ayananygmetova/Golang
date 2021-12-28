package store

import (
	"context"
	"final/internal/models"
)

type Store interface {
	Connect(url string) error
	Close() error

	Categories() CategoriesRepository
	Products() ProductsRepository
	Properties() PropertiesRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) error
	All(ctx context.Context, filter *models.CategoriesFilter) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) error
	Delete(ctx context.Context, id int) error
}

type ProductsRepository interface {
	Create(ctx context.Context, product *models.Product) error
	All(ctx context.Context, filter *models.ProductsFilter) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) error
	Delete(ctx context.Context, id int) error
	// ByCategory(ctx context.Context, id int) ([]*models.Product, error)
}

type PropertiesRepository interface {
	Create(ctx context.Context, property *models.Property) error
	All(ctx context.Context) ([]*models.Property, error)
	ByID(ctx context.Context, id int) (*models.Property, error)
	Update(ctx context.Context, property *models.Property) error
	Delete(ctx context.Context, id int) error
}

type CharacteristicsRepository interface {
	Create(ctx context.Context, characteristic *models.Characteristics) error
	All(ctx context.Context) ([]*models.Characteristics, error)
	ByID(ctx context.Context, id int) (*models.Characteristics, error)
	Update(ctx context.Context, characteristic *models.Characteristics) error
	Delete(ctx context.Context, id int) error
}
