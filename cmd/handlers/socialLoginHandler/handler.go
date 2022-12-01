package socialLoginHandler

import (
	"net/http"
	"test-va/internals/service/socialLoginService"
	// "test-va/internals/service/tokenservice"

	"github.com/gin-gonic/gin"
)

type socialLoginHandler struct {
	srv socialLoginService.LoginSrv
}

func NewCallHandler(srv socialLoginService.LoginSrv) *socialLoginHandler {
	return &socialLoginHandler{srv: srv}
}


func (t *socialLoginHandler) GoogleLogin(c *gin.Context) {
	c.Redirect(http.StatusFound, t.srv.GetAuthUrl())

	
}


func (t *socialLoginHandler) GoogleCallBack(c *gin.Context) {

	code := c.Query("code")


    user := t.srv.LoginResponse(code)
	c.Set("userId", user.UserId)
    // c.Redirect(http.StatusFound, "https://ticked.hng.tech/lo?accessToken="+user.Token)
	c.JSON(http.StatusAccepted, user)
}
