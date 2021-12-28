package postgres

import (
	"context"
	"final/internal/models"
	"final/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Properties() store.PropertiesRepository {
	if db.properties == nil {
		db.properties = NewPropertiesRepository(db.conn)
	}

	return db.properties
}

type PropertiesRepository struct {
	conn *sqlx.DB
}

func NewPropertiesRepository(conn *sqlx.DB) store.PropertiesRepository {
	return &PropertiesRepository{conn: conn}
}

func (p PropertiesRepository) Create(ctx context.Context, property *models.Property) error {
	_, err := p.conn.Exec(`INSERT INTO properties(name) VALUES ($1)`, property.Name)
	if err != nil {
		return err
	}

	return nil
}

func (p PropertiesRepository) All(ctx context.Context) ([]*models.Property, error) {
	basicQuery := "SELECT * FROM properties"
	properties := make([]*models.Property, 0)
	if err := p.conn.SelectContext(ctx, &properties, basicQuery); err != nil {
		return nil, err
	}

	return properties, nil
}

func (p PropertiesRepository) ByID(ctx context.Context, id int) (*models.Property, error) {
	property := new(models.Property)
	if err := p.conn.Get(property, "SELECT id, name FROM properties WHERE id=$1", id); err != nil {
		return nil, err
	}

	return property, nil
}

func (p PropertiesRepository) Update(ctx context.Context, property *models.Property) error {
	_, err := p.conn.Exec("UPDATE properties SET name = $1 WHERE id = $2", property.Name, property.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p PropertiesRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM properties WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
