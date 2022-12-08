package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/cmd/handlers/notificationHandler"
	"test-va/cmd/middlewares"
	"test-va/internals/service/notificationService"
)

func NotificationRoutes(v1 *gin.RouterGroup, service notificationService.NotificationSrv) {
	notificationHandler := notificationHandler.NewNotificationHandler(service)

	not := v1.Group("/notification")

	not.Use(middlewares.ValidateJWT())
	{
		//Create a Notification
		not.POST("", notificationHandler.RegisterForNotifications)

		//Register For Notifications
		not.GET("", notificationHandler.GetNotifications)
	}
}
