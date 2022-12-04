package paymentHandler

import (
	"log"

	"github.com/gin-gonic/gin"
)

func HandleEvent(c *gin.Context) {
	event, err := getEvent(c)

	if err != nil {
		log.Fatal(err)
	}

	log.Println(event.Type)

	if event.Type == "customer.subscription.created" {

		log.Println("Subscription created")
	}

}
