package middlewares

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"test-va/internals/entity/ResponseEntity"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenSrv := tokenservice.NewTokenSrv("tokenString")

		// const BEARER_HEADER = "Bearer "
		authHeader := c.GetHeader("Authorization")
		if authHeader == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid Token")
			return
		}
		auth := strings.Split(authHeader, " ")

		token, err := tokenSrv.ValidateToken(auth[1])

		if err != nil {
			log.Println(err)
			c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("invalid Token: %v", err))

		}

		c.Set("userId", token.Id)
		log.Println("middleware passed")
		c.Next()
	}
}

func CheckUserID() gin.HandlerFunc {
	return func(ctx *gin.Context) {

	}
}

func MapID() gin.HandlerFunc {
	return func(c *gin.Context) {
		userSession := c.GetString("userId")
		userURL := c.Param("user_id")

		if userSession == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Invalid UserID", nil, nil))
			return
		}

		if userURL == "" {
			c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "Invalid url", nil, nil))
			return
		}

		if userSession != userURL {
			c.AbortWithStatusJSON(http.StatusBadRequest, ResponseEntity.BuildErrorResponse(http.StatusBadRequest, "You are not allowed to modify resource", nil, nil))
			return
		}

		c.Next()
	}
}
