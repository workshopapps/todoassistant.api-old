package socialLoginHandler

import (
	"log"
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/service/socialLoginService"
	"test-va/internals/entity/userEntity"

	"github.com/gin-gonic/gin"
)

type socialLoginHandler struct {
	srv socialLoginService.LoginSrv
}

func NewCallHandler(srv socialLoginService.LoginSrv) *socialLoginHandler {
	return &socialLoginHandler{srv: srv}
}



func (t *socialLoginHandler) GoogleLogin(c *gin.Context) {
	var req userEntity.GoogleLoginReq
	
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	user, errorRes := t.srv.LoginResponse(&req)

	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			ResponseEntity.BuildErrorResponse(http.StatusUnauthorized, "Authorization Error", errorRes, nil))
		return
	}

	log.Println("userid -", user.UserId)
	c.Set("userId", user.UserId)
	println(c.GetString("userId"))

	c.JSON(http.StatusOK,ResponseEntity.BuildSuccessResponse(http.StatusOK,"login successful",user,nil))
	
}



func (t *socialLoginHandler) FacebookLogin(c *gin.Context) {
	var req userEntity.FacebookLoginReq
	
	err := c.ShouldBindJSON(&req)

	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	user, errorRes := t.srv.FacebookLoginResponse(&req)

	if errorRes != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			ResponseEntity.BuildErrorResponse(http.StatusUnauthorized, "Authorization Error", errorRes, nil))
		return
	}

	log.Println("userid -", user.UserId)
	c.Set("userId", user.UserId)

    
	c.JSON(http.StatusOK,ResponseEntity.BuildSuccessResponse(http.StatusOK,"login successful",user,nil))
	
}
