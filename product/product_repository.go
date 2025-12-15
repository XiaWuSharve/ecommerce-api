package product

import (
	"context"
	"fmt"
	"log"

	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

type ProductRepository struct {
	pool *pgxpool.Pool
	ctx  context.Context
}

func NewProductRepository(pool *pgxpool.Pool, ctx context.Context) *ProductRepository {
	return &ProductRepository{pool: pool, ctx: ctx}
}

func (p *ProductRepository) prepare(name string, sql string) (pgx.Tx, error) {
	tx, err := p.pool.Begin(p.ctx)
	if err != nil {
		log.Fatal("failed to begin transaction: ", err)
		return nil, err
	}
	tx.Begin(p.ctx)
	_, err = tx.Prepare(p.ctx, name, sql)
	if err != nil {
		return nil, err
	}
	return tx, nil
}

func (p *ProductRepository) FindById(id int) (*ProductEntity, error) {
	tx, err := p.prepare("find_product_by_id", `
	SELECT id, name, description, photo, price, created_at, updated_at 
	FROM products 
	WHERE id = $1 
	ORDER BY created_at DESC`)
	defer tx.Rollback(p.ctx)
	if err != nil {
		return nil, fmt.Errorf("failed to prepare query: %v", err)
	}
	row := tx.QueryRow(p.ctx, "find_product_by_id", id)
	tx.Commit(p.ctx)
	var prod ProductEntity
	err = row.Scan(
		&prod.Id, &prod.Name, &prod.Description, &prod.Photo, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt,
	)
	if err != nil {
		return nil, fmt.Errorf("cannot find product with id %d: %v", id, err)
	}

	return &prod, nil
}

func (p *ProductRepository) Save(prod *ProductDto) (int, error) {
	tx, err := p.prepare("save_product", `
		INSERT INTO products (name, description, photo, price)
		VALUES ($1, $2, $3, $4)
		RETURNING id
	`)
	if err != nil {
		return 0, fmt.Errorf("failed to prepare query: %v", err)
	}

	var id int
	tx.
		QueryRow(p.ctx, "save_product", prod.Name, prod.Description, prod.Photo, prod.Price).
		Scan(&id)

	tx.Commit(p.ctx)
	return id, nil
}
