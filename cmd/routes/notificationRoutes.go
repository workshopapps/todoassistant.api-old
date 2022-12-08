package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/cmd/middlewares"
	"test-va/cmd/handlers/notificationHandler"
	"test-va/internals/service/notificationService"
)

func NotificationRoutes(v1 *gin.RouterGroup, service notificationService.NotificationSrv) {
	notificationHandler := notificationHandler.NewNotificationHandler(service)

	v1.Use(middlewares.ValidateJWT()) 
	{
		//Create a Notification
		v1.POST("/notification", notificationHandler.RegisterForNotifications)

		//Register For Notifications
		v1.GET("/notification", notificationHandler.GetNotifications)
	}
}
