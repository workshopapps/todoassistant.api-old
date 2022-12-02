package userHandler

import (
	"log"
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
	//c.Set("userId", user.UserId)
	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(200, "Created user successfully", user, nil))
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
	log.Println("userid -", user.UserId)
	c.Set("userId", user.UserId)
	println(c.GetString("userId"))
	c.JSON(http.StatusOK, user)
}

func (u *userHandler) GetUsers(c *gin.Context) {
	page := c.Query("page")
	if page == "" {
		page = "1"
	}
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
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, nil, nil))
		return
	}

	c.JSON(http.StatusOK, users)
}

func (u *userHandler) GetUser(c *gin.Context) {
	user, err := u.srv.GetUser(userFromRequest(c))
	if err != nil {
		message := "No user with that ID in the system"
		c.AbortWithStatusJSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, message, nil, nil))
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

func (u *userHandler) ChangePassword(c *gin.Context) {
	var req userEntity.ChangePasswordReq
	userSession, _ := c.Get("userId")
	userURL := userFromRequest(c)
	if userSession != userURL {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Authorization error", nil, nil))
		return
	}

	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	errRes := u.srv.ChangePassword(&req, userURL)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Cannot Change Password", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "Password updated successfully", nil, nil))
}

func (u *userHandler) ResetPassword(c *gin.Context) {
	var req userEntity.ResetPasswordReq
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	token, err := u.srv.ResetPassword(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusOK, "Internal Server Error", err, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "Email sent, check your inbox!", token, nil))
}

func (u *userHandler) ResetPasswordWithToken(c *gin.Context) {
	var req userEntity.ResetPasswordWithTokenReq
	token := string(c.Query("token"))
	userId := string(c.Query("user_id"))
	err := c.ShouldBindJSON(&req)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Bad Request", err, nil))
		return
	}

	errRes := u.srv.ResetPasswordWithToken(&req, token, userId)
	if errRes != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.BuildErrorResponse(http.StatusForbidden, "Cannot Change Password", errRes, nil))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "Password changed successfully", nil, nil))
}

func (u *userHandler) DeleteUser(c *gin.Context) {
	err := u.srv.DeleteUser(userFromRequest(c))
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
		return
	}

	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "User deleted successfully!", nil, nil))
}

//The Id of the Virtual Assistant is Sent Along With this Request

func (u *userHandler) AssignVAToUser(c *gin.Context) {
	user_id := c.GetString("userId")
	va_id := c.Params.ByName("va_id")

	if user_id == "" || va_id == "" {
		c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "no user_id or va_id provided", nil, nil))
		return
	}
	err := u.srv.AssignVAToUser(user_id, va_id)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusInternalServerError, ResponseEntity.NewInternalServiceError(err))
		return
	}
	c.JSON(http.StatusOK, ResponseEntity.BuildSuccessResponse(http.StatusOK, "VA Assigned", nil, nil))
}

func userFromRequest(c *gin.Context) string {
	return c.Param("user_id")
}
