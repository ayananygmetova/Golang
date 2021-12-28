package postgres

import (
	"final/internal/store"

	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

type DB struct {
	conn *sqlx.DB

	categories      store.CategoriesRepository
	products        store.ProductsRepository
	properties      store.PropertiesRepository
	characteristics store.CharacteristicsRepository
}

func (db *DB) Connect(url string) error {
	conn, err := sqlx.Connect("pgx", url)
	if err != nil {
		return err
	}

	if err := conn.Ping(); err != nil {
		return err
	}

	db.conn = conn
	return nil
}

func (db *DB) Close() error {
	return db.conn.Close()
}

func NewDB() store.Store {
	return &DB{}
}
