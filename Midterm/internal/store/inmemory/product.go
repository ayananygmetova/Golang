package inmemory

import (
	"Midterm/internal/models"
	"Midterm/internal/store"
	"context"
	"fmt"
	"sync"
)

type ProductRepo struct {
	data           map[int]*models.Product
	categoriesRepo store.CategoriesRepository
	mu             *sync.RWMutex
}

func (db *ProductRepo) Create(ctx context.Context, product *models.Product) (*models.Product, error) {
	db.mu.Lock()
	defer db.mu.Unlock()
	db.data[product.ID] = product
	return product, nil
}

func (db *ProductRepo) All(ctx context.Context) ([]*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	products := make([]*models.Product, 0, len(db.data))
	for _, product := range db.data {
		products = append(products, product)
	}

	return products, nil
}

func (db *ProductRepo) ByID(ctx context.Context, id int) (*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	product, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no product with id %d", id)
	}

	return product, nil
}

func (db *ProductRepo) Update(ctx context.Context, product *models.Product) (*models.Product, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[product.ID] = product
	return product, nil
}

func (db *ProductRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
func (db *ProductRepo) ByCategory(ctx context.Context, id int) ([]*models.Product, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()
	products := make([]*models.Product, 0, len(db.data))
	for _, product := range db.data {
		if product.CategoryId == id {
			products = append(products, product)
		}
	}

	return products, nil
}
