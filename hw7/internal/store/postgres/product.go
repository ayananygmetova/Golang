package postgres

import (
	"context"
	"hw7/internal/models"
	"hw7/internal/store"

	"github.com/jmoiron/sqlx"
)

func (db *DB) Products() store.ProductsRepository {
	if db.products == nil {
		db.products = NewProductsRepository(db.conn)
	}

	return db.products
}

type ProductsRepository struct {
	conn *sqlx.DB
}

func NewProductsRepository(conn *sqlx.DB) store.ProductsRepository {
	return &ProductsRepository{conn: conn}
}

func (p ProductsRepository) Create(ctx context.Context, product *models.Product) error {
	_, err := p.conn.Exec(`INSERT INTO products(name, manufacturer, description, price, brand, category_id) VALUES ($1, $2, $3, $4, $5, $6)`,
		product.Name, product.Manufacturer, product.Description, product.Price, product.Brand, product.CategoryId)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductsRepository) All(ctx context.Context) ([]*models.Product, error) {
	products := make([]*models.Product, 0)
	if err := p.conn.Select(&products, "SELECT * FROM products"); err != nil {
		return nil, err
	}

	return products, nil
}

func (p ProductsRepository) ByID(ctx context.Context, id int) (*models.Product, error) {
	product := new(models.Product)
	if err := p.conn.Get(product, "SELECT id, name FROM products WHERE id=$1", id); err != nil {
		return nil, err
	}

	return product, nil
}

func (p ProductsRepository) Update(ctx context.Context, product *models.Product) error {
	_, err := p.conn.Exec("UPDATE products SET name = $1 WHERE id = $2", product.Name, product.ID)
	if err != nil {
		return err
	}

	return nil
}

func (p ProductsRepository) Delete(ctx context.Context, id int) error {
	_, err := p.conn.Exec("DELETE FROM products WHERE id = $1", id)
	if err != nil {
		return err
	}

	return nil
}
