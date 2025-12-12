package main

import (
	"context"
	"fmt"
	"log"
	"net/http"

	"github.com/XiaWuSharve/ecommerce-api/product"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/labstack/echo/v4"
)

var pool *pgxpool.Pool
var ctx = context.Background()

var productRepository *product.ProductRepository

func init() {
	var err error
	pool, err = pgxpool.New(ctx, "postgres://postgres:admin@localhost:5432/gotodo")
	if err != nil {
		log.Fatal("unable to connect postgresql: ", err)
	}

	if err := pool.Ping(ctx); err != nil {
		log.Fatal("unable to ping postgresql: ", err)
	}

	fmt.Println("Connected to PostgreSQL database")

	productRepository = product.NewProductRepository(pool, ctx)
}

func main() {
	e := echo.New()
	// e.Use(echojwt.WithConfig(echojwt.Config{}))

	// e.POST("/product", saveProduct)
	// e.GET("/product/:id", getProduct)
	// e.PUT("/product/:id", updateProduct)
	// e.DELETE("/product/:id", deleteProduct)

	e.GET("/", func(c echo.Context) error {
		prod, err := productRepository.GetProductsById(1)
		if err != nil {
			return err
		}
		return c.JSON(http.StatusOK, prod)
		// return c.String(http.StatusOK, "Hello, World!")
	})
	e.Logger.Fatal(e.Start(":1323"))
}
