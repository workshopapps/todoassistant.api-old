package routes

import (
	"test-va/cmd/handlers/subscribeHandler"
	"test-va/internals/service/subscribeService"

	"github.com/gin-gonic/gin"
)

func SubscribeRoutes(v1 *gin.RouterGroup, service subscribeService.SubscribeService) {

	handler := subscribeHandler.NewSubscribeHandler(service)
	// subscribe to newsletter route
	subscribe := v1.Group("/subscribe")
	{
		subscribe.POST("", handler.AddSubscriber)
	}

	contact := v1.Group("/contact-us")
	{
		contact.POST("", handler.ContactUs)
	}
}
