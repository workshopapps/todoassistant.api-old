package routes

import (
	"test-va/cmd/handlers/notificationHandler"
	"test-va/cmd/middlewares"
	"test-va/internals/service/notificationService"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func NotificationRoutes(v1 *gin.RouterGroup, service notificationService.NotificationSrv, tokenSrv tokenservice.TokenSrv) {
	notificationHandler := notificationHandler.NewNotificationHandler(service)
	jwtMWare := middlewares.NewJWTMiddleWare(tokenSrv)

	not := v1.Group("/notification")

	not.Use(jwtMWare.ValidateJWT())
	{
		//Create a Notification
		not.POST("", notificationHandler.RegisterForNotifications)

		//Get Notifications
		not.GET("", notificationHandler.GetNotifications)

		//Delete Notifications
		not.DELETE("", notificationHandler.DeleteNotifications)

		not.PATCH("/:notification_id", notificationHandler.UpdateNotification)
	}
}
