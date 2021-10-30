package inmemory

import (
	"Midterm/internal/models"
	"context"
	"fmt"
	"sync"
)

type CategoryRepo struct {
	data map[int]*models.Category
	mu   *sync.RWMutex
}

func (db *CategoryRepo) Create(ctx context.Context, category *models.Category) (*models.Category, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[category.ID] = category
	return category, nil
}

func (db *CategoryRepo) All(ctx context.Context) ([]*models.Category, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	categories := make([]*models.Category, 0, len(db.data))
	for _, category := range db.data {
		categories = append(categories, category)
	}

	return categories, nil
}

func (db *CategoryRepo) ByID(ctx context.Context, id int) (*models.Category, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	category, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no category with id %d", id)
	}

	return category, nil
}

func (db *CategoryRepo) Update(ctx context.Context, category *models.Category) (*models.Category, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[category.ID] = category
	return category, nil
}

func (db *CategoryRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
