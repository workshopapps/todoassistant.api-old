package paymentHandler

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74/customer"
)

func HandleEvent(c *gin.Context) {
	event, err := getEvent(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(event.Type)

	if event.Type == "customer.subscription.created" {
		c, err := customer.Get(event.Data.Object["customer"].(string), nil)
		if err != nil {
			log.Fatal(err)
		}
		email := c.Metadata["FinalEmail"]
		log.Println("Subscription created by", email)
	}

}
