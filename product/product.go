package product

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type ProductDto struct {
	Name        string  `json:"name"`
	Description string  `json:"description"`
	Photo       string  `json:"photo"`
	Price       float64 `json:"price"`
}

type ProductEntity struct {
	Id          int            `json:"id"`
	Name        string         `json:"name"`
	Description pgtype.Text    `json:"description"`
	Photo       pgtype.Text    `json:"photo"`
	Price       pgtype.Numeric `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
