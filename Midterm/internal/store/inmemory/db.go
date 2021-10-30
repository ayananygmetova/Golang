package inmemory

import (
	"Midterm/internal/models"
	"Midterm/internal/store"
	"sync"
)

type DB struct {
	categoriesRepo store.CategoriesRepository
	productsRepo   store.ProductsRepository

	mu *sync.RWMutex
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
