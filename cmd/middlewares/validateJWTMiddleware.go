package middlewares

import (
	"log"
	"net/http"
	"strings"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenSrv := tokenservice.NewTokenSrv("tokenString")

		// const BEARER_HEADER = "Bearer "
		authHeader := c.GetHeader("Authorization")
		auth := strings.Split(authHeader, " ")

		token, err := tokenSrv.ValidateToken(auth[1])

		if err != nil {
			log.Println(err)
			c.AbortWithStatus(http.StatusUnauthorized)

		}

		c.Set("userId", token.Id)
		log.Println("middleware passed")
		c.Next()
	}
}
