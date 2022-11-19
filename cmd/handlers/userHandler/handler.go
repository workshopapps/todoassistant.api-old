package userHandler

import (
	"net/http"
	"test-va/internals/entity/ResponseEntity"
	"test-va/internals/entity/userEntity"
	"test-va/internals/service/userService"

	"github.com/gin-gonic/gin"
)

type userHandler struct {
	srv userService.UserSrv
}

func NewUserHandler(srv userService.UserSrv) *userHandler {
	return &userHandler{srv: srv}
}

func (u *userHandler) CreateUser(c *gin.Context) {
	var req userEntity.CreateUserReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error decoding into struct", err, nil))
		return
	}
	user, errorRes := u.srv.SaveUser(&req)
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "error saving into db", errorRes, nil))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) Login(c *gin.Context) {
	var req userEntity.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "invalid login credentials", err, nil))
		return
	}

	user, errorRes := u.srv.Login(&req)
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "invalid login credentials", errorRes, nil))
		return
	}

	c.JSON(http.StatusOK, user)
}
