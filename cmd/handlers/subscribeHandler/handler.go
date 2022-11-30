package subscribeHandler

import (
	"net/http"
	"test-va/internals/service/subscribeService"

	"github.com/gin-gonic/gin"
)


type subscribeHandler struct{
	srv subscribeService.SubscribeService
}

func NewSubscribeHandler(srv subscribeService.SubscribeService) *subscribeHandler{
	return &subscribeHandler{srv: srv}
}

func (t *subscribeHandler) AddSubscriber(c *gin.Context, ){
	// create a request of subsctibeEntity type

	// copy data from gin context to req

	// call function from service that saves email to DB

	c.JSON(http.StatusOK, "Subscribe to our newsletter")
}
