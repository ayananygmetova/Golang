package postgres

import (
	"context"
	"final/internal/models"
	"final/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) ProductCharacteristics() store.ProductCharacteristicsRepository {
	if db.product_characteristics == nil {
		db.product_characteristics = NewProductCharacteristicsRepository(db.conn)
	}

	return db.product_characteristics
}

type ProductCharacteristicsRepository struct {
	conn *sqlx.DB
}

func NewProductCharacteristicsRepository(conn *sqlx.DB) store.ProductCharacteristicsRepository {
	return &ProductCharacteristicsRepository{conn: conn}
}

func (p ProductCharacteristicsRepository) Create(ctx context.Context, product_char *models.ProductCharacteristics) error {
	_, err := p.conn.Exec(`INSERT INTO product_characteristics(product_id, characteristics_id) VALUES ($1, $2)`,
		product_char.ProductId, product_char.CharacteristicsId)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductCharacteristicsRepository) ByID(ctx context.Context, id int) ([]*models.Characteristics, error) {
	characteristics := make([]*models.Characteristics, 0)
	if err := p.conn.Get(characteristics, "SELECT id, name FROM product_characteristics t left join characteristics t2 on t.characteristics_id=t2.id left join properties t3 on t3.id=t2.property_id where t.product_id=$1", id); err != nil {
		return nil, err
	}

	return characteristics, nil
}

func (p ProductCharacteristicsRepository) Delete(ctx context.Context, product_id int) error {
	_, err := p.conn.Exec("DELETE FROM product_characteristics WHERE product_id = $1", product_id)
	if err != nil {
		return err
	}

	return nil
}
