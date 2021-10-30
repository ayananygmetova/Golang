package store

import (
	"Midterm/internal/models"
	"context"
)

type Store interface {
	Categories() CategoriesRepository
	Products() ProductsRepository
	Characteristics() CharacteristicsRepository
	Properties() PropertiesRepository
}

type CategoriesRepository interface {
	Create(ctx context.Context, category *models.Category) (*models.Category, error)
	All(ctx context.Context) ([]*models.Category, error)
	ByID(ctx context.Context, id int) (*models.Category, error)
	Update(ctx context.Context, category *models.Category) (*models.Category, error)
	Delete(ctx context.Context, id int) error
}

type ProductsRepository interface {
	Create(ctx context.Context, product *models.Product) (*models.Product, error)
	All(ctx context.Context) ([]*models.Product, error)
	ByID(ctx context.Context, id int) (*models.Product, error)
	Update(ctx context.Context, product *models.Product) (*models.Product, error)
	Delete(ctx context.Context, id int) error
	ByCategory(ctx context.Context, id int) ([]*models.Product, error)
}

type CharacteristicsRepository interface {
	Create(ctx context.Context, characteristics *models.Characteristics) (*models.Characteristics, error)
	All(ctx context.Context) ([]*models.Characteristics, error)
	ByID(ctx context.Context, id int) (*models.Characteristics, error)
	Update(ctx context.Context, characteristics *models.Characteristics) (*models.Characteristics, error)
	Delete(ctx context.Context, id int) error
}
type PropertiesRepository interface {
	Create(ctx context.Context, property *models.Property) (*models.Property, error)
	All(ctx context.Context) ([]*models.Property, error)
	ByID(ctx context.Context, id int) (*models.Property, error)
	Update(ctx context.Context, property *models.Property) (*models.Property, error)
	Delete(ctx context.Context, id int) error
}
