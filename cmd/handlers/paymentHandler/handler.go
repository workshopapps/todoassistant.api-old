package paymentHanlder

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
)

func Paymentsrv() gin.HandlerFunc {
	return func(c *gin.Context) {
		params := &stripe.CheckoutSessionParams{
			// all both URL are just placeholders remove when FE is ready
			SuccessURL: stripe.String("https://www.shutterstock.com/image-vector/green-tick-checkbox-vector-illustration-isolated-428282710"),
			CancelURL:  stripe.String("https://www.shutterstock.com/image-vector/vector-illustration-word-fail-red-ink-1077767732"),
			PaymentMethodTypes: stripe.StringSlice([]string{
				"card",
			}),
			Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),

			LineItems: []*stripe.CheckoutSessionLineItemParams{
				&stripe.CheckoutSessionLineItemParams{
					// when front end is set up set req params for price and quantity
					Price:    stripe.String("price_H5ggYwtDq4fbrJ"),
					Quantity: stripe.Int64(1),
				},
			},
		}
		s, _ := session.New(params)
		c.JSON(http.StatusOK, s)

	}
}
