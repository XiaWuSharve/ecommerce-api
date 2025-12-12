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

func (p *ProductRepository) GetProductsById(id int) (Product, error) {
	tx, err := p.prepare("get_product_by_id", `
	SELECT id, name, description, photo, price, created_at, updated_at 
	FROM products 
	WHERE id = $1 
	ORDER BY created_at DESC`)
	if err != nil {
		log.Fatal("failed to prepare query: ", err)
	}
	rows, err := tx.Query(p.ctx, "get_product_by_id", id)
	if err != nil {
		log.Fatal("failed to execute query: ", err)
	}
	defer rows.Close()
	// err = tx.Commit(p.ctx)
	// if err != nil {
	// 	log.Fatal("failed to commit transaction: ", err)
	// }
	var prod Product
	if rows.Next() {
		err = rows.Scan(
			&prod.Id, &prod.Name, &prod.Description, &prod.Photo, &prod.Price, &prod.CreatedAt, &prod.UpdatedAt,
		)
		if err != nil {
			log.Fatal("failed to read a row from query: ", err)
		}
	} else {
		return prod, fmt.Errorf("product not found")
	}

	if err := rows.Err(); err != nil {
		log.Fatal("failed while reading rows from query: ", err)
	}

	return prod, nil
}
