package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/cmd/handlers/vaHandler"
	"test-va/internals/service/tokenService"
	"test-va/internals/service/vaService"
)

func VARoutes(v1 *gin.RouterGroup, service vaService.VAService, srv tokenservice.TokenSrv) {
	handler := vaHandler.NewVaHandler(srv, service)
	//mWare := vaMiddleware.NewVaMiddleWare(srv)

	va := v1.Group("/va")
	va.POST("/:va_id", handler.UpdateUser)
	va.POST("/login", handler.Login)

	//va.Use(mWare.MapMasterToReq)

	{
		//master middleware
		va.POST("/signup", handler.SignUp)
		va.POST("/delete/:va_id", handler.DeleteUser)
		va.POST("/change-password", handler.ChangePassword)
	}

}
