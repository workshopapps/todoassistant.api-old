package routes

import (
	"test-va/cmd/handlers/vaHandler"
	"test-va/cmd/middlewares/vaMiddleware"
	"test-va/internals/service/taskService"
	tokenservice "test-va/internals/service/tokenService"
	"test-va/internals/service/userService"
	"test-va/internals/service/vaService"

	"github.com/gin-gonic/gin"
)

func VARoutes(v1 *gin.RouterGroup, service vaService.VAService, srv tokenservice.TokenSrv, taskService taskService.TaskService, userService userService.UserSrv) {
	handler := vaHandler.NewVaHandler(srv, service, taskService, userService)
	mWare := vaMiddleware.NewVaMiddleWare(srv)

	va := v1.Group("/va")
	va.POST("/:va_id", handler.UpdateVA)
	va.GET("/:va_id", handler.GetVAByID)
	va.POST("/login", handler.Login)
	va.GET("/user/:va_id", handler.GetUserAssignedToVA)
	va.GET("/user/task/:user_id", handler.GetTaskByUser)
	va.GET("/user/profile/:user_id", handler.GetSingleUserProfile)

	va.Use(mWare.MapMasterToReq)

	{
		//master middleware
		va.POST("/signup", handler.SignUp)
		va.POST("/delete/:va_id", handler.DeleteVA)
		va.POST("/change-password", handler.ChangePassword)
	}

}
