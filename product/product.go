package product

import (
	"time"

	"github.com/jackc/pgx/v5/pgtype"
)

type Product struct {
	Id          int            `json:"id"`
	Name        pgtype.Text    `json:"name"`
	Description pgtype.Text    `json:"description"`
	Photo       pgtype.Text    `json:"photo"`
	Price       pgtype.Numeric `json:"price"`
	CreatedAt   time.Time      `json:"created_at"`
	UpdatedAt   time.Time      `json:"updated_at"`
}
