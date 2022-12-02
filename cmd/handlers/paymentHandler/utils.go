package paymentHandler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v74"
)

func getEvent(c *gin.Context) (eventRes *stripe.Event, err error) {
	const MaxBodyBytes = int64(65536)
	var w http.ResponseWriter = c.Writer
	c.Request.Body = http.MaxBytesReader(w, c.Request.Body, MaxBodyBytes)
	payload, err := io.ReadAll(c.Request.Body)
	if err != nil {
		return nil, err
	}

	event := stripe.Event{}
	err = json.Unmarshal(payload, &event)
	if err != nil {
		return nil, err
	}
	return &event, nil
}
