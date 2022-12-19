package routes

import (
	"test-va/cmd/handlers/socialLoginHandler"
	"test-va/internals/service/socialLoginService"

	"github.com/gin-gonic/gin"
)

func SocialLoginRoute(v1 *gin.RouterGroup, srv socialLoginService.LoginSrv) {
	loginHandler := socialLoginHandler.NewLoginHandler(srv)

	v1.POST("/googlelogin", loginHandler.GoogleLogin)
	v1.POST("/facebooklogin", loginHandler.FacebookLogin)

}
