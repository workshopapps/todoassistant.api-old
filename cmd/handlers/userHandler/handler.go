package userHandler

import (
	"net/http"
	"strconv"
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
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Input Data", err, nil))
		return
	}
	user, errorRes := u.srv.SaveUser(&req)
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusInternalServerError, "Failed To Save User", errorRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(200, "created user successfully", user, nil))
}

func (u *userHandler) Login(c *gin.Context) {
	var req userEntity.LoginReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest,
			ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	user, errorRes := u.srv.Login(&req)
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized,
			ResponseEntity.BuildErrorResponse(http.StatusUnauthorized, "Authorization Error", errorRes, nil))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) GetUsers(c *gin.Context) {
	page := c.Query("page")

	pageInt, err := strconv.Atoi(page)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
		return
	}

	users, err := u.srv.GetUsers(pageInt)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
		return
	}

	length := len(users)
	if length == 0 {
		message := "No users in the system"
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, nil))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *userHandler) GetUser(c *gin.Context) {
	user, err := u.srv.GetUser(userFromRequest(c))
	if err != nil {
		message := "No user with that ID in the system"
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, nil))
		return
	}
	// else if err != sql.ErrNoRows {
	// 	c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
	// 	return
	// }

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) UpdateUser(c *gin.Context) {
	var req userEntity.UpdateUserReq

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	user, errorRes := u.srv.UpdateUser(&req, userFromRequest(c))
	if errorRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Cannot Update!", err, nil))
		return
	}

	c.JSON(http.StatusOK, user)
}

func (u *userHandler) DeleteUser(c *gin.Context) {
	err := u.srv.DeleteUser(userFromRequest(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "User deleted successfully!", nil, nil))
}

func userFromRequest(c *gin.Context) string {
	return c.Param("user_id")
}
