package routes

import (
	"github.com/gin-gonic/gin"
	"test-va/internals/service/socialLoginService"
	"test-va/cmd/handlers/socialLoginHandler"
)

func SocialLoginRoute(v1 *gin.RouterGroup, srv socialLoginService.LoginSrv) {
	loginHandler := socialLoginHandler.NewLoginHandler(srv)

	v1.POST("/googlelogin", loginHandler.GoogleLogin)
	v1.POST("/facebooklogin", loginHandler.FacebookLogin)
	
}
