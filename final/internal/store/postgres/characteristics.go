package postgres

import (
	"context"
	"final/internal/models"
	"final/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Characteristics() store.CharacteristicsRepository {
	if db.characteristics == nil {
		db.characteristics = NewCharacteristicsRepository(db.conn)
	}

	return db.characteristics
}

type CharacteristicsRepository struct {
	conn *sqlx.DB
}

func NewCharacteristicsRepository(conn *sqlx.DB) store.CharacteristicsRepository {
	return &CharacteristicsRepository{conn: conn}
}

func (p CharacteristicsRepository) Create(ctx context.Context, characteristic *models.Characteristics) error {
	_, err := p.conn.Exec(`INSERT INTO characteristics(value, property_id) VALUES ($1, $2)`,
		characteristic.Value, characteristic.PropertyId)
	if err != nil {
		return err
	}

	return nil
}

func (p CharacteristicsRepository) All(ctx context.Context) ([]*models.Characteristics, error) {
	basicQuery := "SELECT * FROM characteristics"
	characteristics := make([]*models.Characteristics, 0)
	if err := p.conn.SelectContext(ctx, &characteristics, basicQuery); err != nil {
		return nil, err
	}

	return characteristics, nil
}

func (p CharacteristicsRepository) ByID(ctx context.Context, id int) (*models.Characteristics, error) {
	characteristic := new(models.Characteristics)
	if err := p.conn.Get(characteristic, "SELECT id, value, property_id FROM characteristics WHERE id=$1", id); err != nil {
		return nil, err
	}

	return characteristic, nil
}

func (p CharacteristicsRepository) Update(ctx context.Context, characteristic *models.Characteristics) error {
	_, err := p.conn.Exec("UPDATE characteristics SET value = $1 WHERE id = $2", characteristic.Value, characteristic.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p CharacteristicsRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM characteristics WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
