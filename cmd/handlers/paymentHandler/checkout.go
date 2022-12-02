package paymentHandler

import (
	"log"
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/stripe/stripe-go/v74"
	"github.com/stripe/stripe-go/v74/checkout/session"
	"github.com/stripe/stripe-go/v74/customer"
)

var PriceId = "price_1MAVOdFf5hgzULICDQBFEfDH"

func checkout(email string) (*stripe.CheckoutSession, error) {

	customerParams := stripe.CustomerParams{
		Email: stripe.String(email),
	}
	customerParams.AddMetadata("FinalEmail", email)
	newCustomer, err := customer.New(&customerParams)

	if err != nil {
		return nil, err
	}

	meta := map[string]string{
		"FinalEmail": email,
	}
	log.Println("Creating Meta for user: ", meta)

	params := &stripe.CheckoutSessionParams{
		Customer: &newCustomer.ID,
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
				Price:    stripe.String(PriceId),
				Quantity: stripe.Int64(1),
			},
		},
		SubscriptionData: &stripe.CheckoutSessionSubscriptionDataParams{
			TrialPeriodDays: stripe.Int64(7),
			Metadata:        meta,
		},
	}
	return session.New(params)
}

type EmailInput struct {
	Email string `json:"email"`
}

type SessionOutput struct {
	Id string `json:"id"`
}

func CheckoutCreator(c *gin.Context) {
	input := &EmailInput{}

	err := c.ShouldBindJSON(input)

	if err != nil {
		log.Fatal(err)
	}

	stripeSession, err := checkout(input.Email)
	if err != nil {
		log.Fatal(err)
	}

	c.JSON(http.StatusOK, &SessionOutput{Id: stripeSession.ID})

	if err != nil {
		log.Fatal(err)
	}

}
