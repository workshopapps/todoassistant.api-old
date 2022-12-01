package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/internals/service/socialLoginService"
	"test-va/cmd/handlers/socialLoginHandler"
)

func SocialLoginRoute(v1 *gin.RouterGroup, srv socialLoginService.LoginSrv) {
	loginHandler := socialLoginHandler.NewCallHandler(srv)

	

	v1.GET("/googlelogin", loginHandler.GoogleLogin)
	v1.GET("/callback", loginHandler.GoogleCallBack)
	
}
