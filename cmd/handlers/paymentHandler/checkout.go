package paymentHandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
)

func checkout(itemPrice string) (*stripe.CheckoutSession, error) {

	var PriceId = itemPrice
	customerParams := stripe.CustomerParams{
		Description: stripe.String("Ticked Premium Customer"),
	}

	newCustomer, err := customer.New(&customerParams)

	if err != nil {
		return nil, err
	}

	params := &stripe.CheckoutSessionParams{
		Customer:   &newCustomer.ID,
		SuccessURL: stripe.String("https://ticked.hng.tech/success"),
		CancelURL:  stripe.String("https://ticked.hng.tech/cancel"),
		PaymentMethodTypes: stripe.StringSlice([]string{
			"card",
		}),
		Mode: stripe.String(string(stripe.CheckoutSessionModeSubscription)),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			&stripe.CheckoutSessionLineItemParams{
				Price:    stripe.String(PriceId),
				Quantity: stripe.Int64(1),
			},
		},
	}
	return session.New(params)
}

type PriceInput struct {
	Price string `json:"price"`
}

type SessionOutput struct {
	Id string `json:"id"`
}

func CheckoutCreator(c *gin.Context) {
	input := &PriceInput{}

	err := c.ShouldBindJSON(input)

	if err != nil {
		log.Fatal(err)
	}

	stripeSession, err := checkout(input.Price)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, &SessionOutput{Id: stripeSession.ID})

	if err != nil {
		log.Fatal(err)
	}

}
