package inmemory

import (
	"Midterm/internal/models"
	"context"
	"fmt"
	"sync"
)

type PropertiesRepo struct {
	data map[int]*models.Property
	mu   *sync.RWMutex
}

func (db *PropertiesRepo) Create(ctx context.Context, property *models.Property) (*models.Property, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[property.ID] = property
	return property, nil
}

func (db *PropertiesRepo) All(ctx context.Context) ([]*models.Property, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	properties := make([]*models.Property, 0, len(db.data))
	for _, property := range db.data {
		properties = append(properties, property)
	}

	return properties, nil
}

func (db *PropertiesRepo) ByID(ctx context.Context, id int) (*models.Property, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	property, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no property with id %d", id)
	}

	return property, nil
}

func (db *PropertiesRepo) Update(ctx context.Context, property *models.Property) (*models.Property, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[property.ID] = property
	return property, nil
}

func (db *PropertiesRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
