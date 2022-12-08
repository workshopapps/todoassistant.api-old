package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/cmd/middlewares"
	"test-va/cmd/handlers/notificationHandler"
	"test-va/internals/service/notificationService"
)

func NotificationRoutes(v1 *gin.RouterGroup, service notificationService.NotificationSrv) {
	notificationHandler := notificationHandler.NewNotificationHandler(service)
	v1.POST("/notification", notificationHandler.RegisterForNotifications)

	v1.Use(middlewares.ValidateJWT()) 
	{
		v1.GET("/notification", notificationHandler.GetNotifications)
	}
}
