package subscribeHandler

import (
	"net/http"

	"github.com/gin-gonic/gin"
)


type subscribeHandler struct{

}

func NewSubscribeHandler() *subscribeHandler{
	return &subscribeHandler{}
}

func (t *subscribeHandler) AddSubscriber(c *gin.Context){
	// create a request of subsctibeEntity type

	// copy data from gin context to req

	// call function from service that saves email to DB

	c.JSON(http.StatusOK, "Subscribe to our newsletter")
}
