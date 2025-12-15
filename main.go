package main

import (
	"context"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"

	"github.com/XiaWuSharve/ecommerce-api/connection"
	"github.com/XiaWuSharve/ecommerce-api/my_jwt"
	"github.com/XiaWuSharve/ecommerce-api/product"
	"github.com/XiaWuSharve/ecommerce-api/stripe"
	"github.com/golang-jwt/jwt/v5"
	echojwt "github.com/labstack/echo-jwt/v4"
	"github.com/labstack/echo/v4"
	"github.com/labstack/echo/v4/middleware"
)

var productRepository *product.ProductRepository
var signedKey = "sharve"

func init() {
	ctx := context.Background()
	pool, err := connection.NewPostgresConnection("postgres://postgres:admin@localhost:5432/gotodo", ctx)
	if err != nil {
		panic(err)
	}
	fmt.Println("Connected to PostgreSQL database")

	productRepository = product.NewProductRepository(pool, ctx)
}

func main() {
	global := echo.New()
	// auth := global.Group("")
	admin := global.Group("")

	global.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{
			"http://127.0.0.1:3000",
			"http://localhost:3000",
		},
		AllowMethods: []string{
			http.MethodGet,
			http.MethodPost,
			http.MethodOptions,
		},
		AllowHeaders: []string{
			echo.HeaderOrigin,
			echo.HeaderContentType,
			echo.HeaderAccept,
			echo.HeaderAuthorization,
		},
	}))
	admin.Use(echojwt.WithConfig(echojwt.Config{
		SigningKey: []byte(signedKey),
		NewClaimsFunc: func(c echo.Context) jwt.Claims {
			return new(my_jwt.JwtCustomClaims)
		},
	}))

	global.POST("/login", func(c echo.Context) error {
		username := c.FormValue("username")
		password := c.FormValue("password")
		rememberMe, err := strconv.Atoi(c.FormValue("remember"))
		if username == "" || password == "" || err != nil {
			return echo.ErrBadRequest
		}

		// find user in persistent layer.
		if username != "sharve" || password != "admin" {
			return echo.ErrUnauthorized
		}
		isAdmin := true
		var expiredDuration int
		if rememberMe == 1 {
			expiredDuration = 365
		} else {
			expiredDuration = 3
		}

		claims := &my_jwt.JwtCustomClaims{
			Name:  username,
			Admin: isAdmin,
			RegisteredClaims: jwt.RegisteredClaims{
				ExpiresAt: jwt.NewNumericDate(
					time.Now().Add(time.Hour * time.Duration(24*expiredDuration))),
			},
		}

		token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

		t, err := token.SignedString([]byte(signedKey))
		if err != nil {
			return err
		}

		return c.String(http.StatusOK, t)
	})

	// 如果有些赠品，打算不限制价格为负数了（？
	admin.POST("/product", func(c echo.Context) error {
		user := c.Get("user").(*jwt.Token)
		if !my_jwt.IsAdmin(user) {
			return echo.ErrUnauthorized
		}
		prod := new(product.ProductDto)
		err := echo.FormFieldBinder(c).
			MustString("name", &prod.Name).
			String("description", &prod.Description).
			String("photo", &prod.Photo).
			Float64("price", &prod.Price).BindError()
		if err != nil {
			return c.String(http.StatusBadRequest, "bad request")
		}

		id, err := productRepository.Save(prod)
		if err != nil {
			return c.String(http.StatusInternalServerError, err.Error())
		}

		return c.String(http.StatusCreated, strconv.Itoa(id))
	})

	global.GET("/product/:id", func(c echo.Context) error {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			log.Fatal("failed to convert string to int: ", err)
		}
		prod, err := productRepository.FindById(id)
		if err != nil {
			return c.String(http.StatusBadRequest, err.Error())
		}
		return c.JSON(http.StatusOK, prod)
	})

	// e.PUT("/product/:id", updateProduct)
	// e.DELETE("/product/:id", deleteProduct)

	global.GET("/stripe", stripe.CreateCheckoutSession)

	global.Logger.Fatal(global.Start(":1323"))
}
