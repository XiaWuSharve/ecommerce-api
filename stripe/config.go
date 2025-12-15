package stripe

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/stripe/stripe-go/v83"
	"github.com/stripe/stripe-go/v83/checkout/session"
)

type Config struct {
	Key string
}

func GetKey() string {
	return "<<secret key>>"
}

func GetPubKey() string {
	return ""
}

func GetPrivKey() string {
	return ""
}

func CreateCheckoutSession(c echo.Context) error {
	stripe.Key = GetPrivKey()

	// Set your domain
	domain := "http://localhost:1323"

	// Create line items for what customer is purchasing
	params := &stripe.CheckoutSessionParams{
		UIMode:    stripe.String("custom"),
		ReturnURL: stripe.String(domain + "/complete?session_id=1"),
		Mode:      stripe.String(string(stripe.CheckoutSessionModePayment)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("cny"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String("Product Name"),
					},
					UnitAmount: stripe.Int64(1000),
				},
				Quantity: stripe.Int64(1),
			},
		},
	}

	s, err := session.New(params)
	if err != nil {
		return err
	}

	// Return the client secret to your frontend
	return c.JSON(http.StatusOK, struct {
		ClientSecret string `json:"clientSecret"`
	}{
		ClientSecret: s.ClientSecret,
	})
}
