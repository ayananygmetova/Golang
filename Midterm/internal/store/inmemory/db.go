package inmemory

import (
	"Midterm/internal/models"
	"Midterm/internal/store"
	"sync"
)

type DB struct {
	categoriesRepo      store.CategoriesRepository
	productsRepo        store.ProductsRepository
	characteristicsRepo store.CharacteristicsRepository
	propertiesRepo      store.PropertiesRepository
	mu                  *sync.RWMutex
}

func NewDB() store.Store {
	return &DB{
		mu: new(sync.RWMutex),
	}
}

func (db *DB) Categories() store.CategoriesRepository {
	if db.categoriesRepo == nil {
		db.categoriesRepo = &CategoryRepo{
			data: make(map[int]*models.Category),
			mu:   new(sync.RWMutex),
		}
	}

	return db.categoriesRepo
}

func (db *DB) Products() store.ProductsRepository {
	if db.productsRepo == nil {
		db.productsRepo = &ProductRepo{
			data:           make(map[int]*models.Product),
			categoriesRepo: db.Categories(),
			mu:             new(sync.RWMutex),
		}
	}

	return db.productsRepo
}

func (db *DB) Properties() store.PropertiesRepository {
	if db.propertiesRepo == nil {
		db.propertiesRepo = &PropertiesRepo{
			data: make(map[int]*models.Property),
			mu:   new(sync.RWMutex),
		}
	}

	return db.propertiesRepo
}
func (db *DB) Characteristics() store.CharacteristicsRepository {
	if db.characteristicsRepo == nil {
		db.characteristicsRepo = &CharacteristicsRepo{
			data: make(map[int]*models.Characteristics),
			mu:   new(sync.RWMutex),
		}
	}

	return db.characteristicsRepo
}
