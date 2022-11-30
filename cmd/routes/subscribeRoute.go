package routes

import (
	"test-va/cmd/handlers/subscribeHandler"

	"github.com/gin-gonic/gin"
)

func SubscribeRoutes(v1 *gin.RouterGroup){

	handler:= subscribeHandler.NewSubscribeHandler()
	// subscribe to newsletter route
	subscribe := v1.Group("/subscribe")
	{
		subscribe.POST("",handler.AddSubscriber)
	}
}
