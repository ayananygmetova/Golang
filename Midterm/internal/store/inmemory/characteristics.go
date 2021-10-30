package inmemory

import (
	"Midterm/internal/models"
	"context"
	"fmt"
	"sync"
)

type CharacteristicsRepo struct {
	data map[int]*models.Characteristics
	mu   *sync.RWMutex
}

func (db *CharacteristicsRepo) Create(ctx context.Context, characteristics *models.Characteristics) (*models.Characteristics, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[characteristics.ID] = characteristics
	return characteristics, nil
}

func (db *CharacteristicsRepo) All(ctx context.Context) ([]*models.Characteristics, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	characteristics := make([]*models.Characteristics, 0, len(db.data))
	for _, characteristic := range db.data {

		characteristics = append(characteristics, characteristic)
	}

	return characteristics, nil
}

func (db *CharacteristicsRepo) ByID(ctx context.Context, id int) (*models.Characteristics, error) {
	db.mu.RLock()
	defer db.mu.RUnlock()

	characteristic, ok := db.data[id]
	if !ok {
		return nil, fmt.Errorf("no characteristic with id %d", id)
	}

	return characteristic, nil
}

func (db *CharacteristicsRepo) Update(ctx context.Context, characteristic *models.Characteristics) (*models.Characteristics, error) {
	db.mu.Lock()
	defer db.mu.Unlock()

	db.data[characteristic.ID] = characteristic
	return characteristic, nil
}

func (db *CharacteristicsRepo) Delete(ctx context.Context, id int) error {
	db.mu.Lock()
	defer db.mu.Unlock()

	delete(db.data, id)
	return nil
}
