package routes

import (
	"test-va/cmd/handlers/userHandler"
	"test-va/cmd/middlewares"
	"test-va/internals/service/userService"

	"github.com/gin-gonic/gin"
)

func UserRoutes(v1 *gin.RouterGroup, srv userService.UserSrv) {
	userHandler := userHandler.NewUserHandler(srv)

	// Register a user
	v1.POST("/user", userHandler.CreateUser)
	// Login into the user account
	v1.POST("/user/login", userHandler.Login)
	// Get a reset password token
	v1.POST("/user/reset-password", userHandler.ResetPassword)
	// Reset password with token id
	v1.PUT("/user/reset-password", userHandler.ResetPasswordWithToken)
	users := v1.Group("/user")

	users.Use(middlewares.ValidateJWT())
	{
		// Get all users
		users.GET("", userHandler.GetUsers)
		// Get a specific user
		users.GET("/:user_id", userHandler.GetUser)
		// Update a specific user
		users.PUT("/:user_id", userHandler.UpdateUser)
		// Change user password

		users.PUT("/:user_id/change-password", userHandler.ChangePassword)
		// Delete a user
		users.DELETE("/:user_id", userHandler.DeleteUser)
	}
}
