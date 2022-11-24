package middlewares

import (
	"log"
	"net/http"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

func ValidateJWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenSrv :=  tokenservice.NewTokenSrv("tokenString")

		const BEARER_HEADER = "Bearer "
		authHeader := c.GetHeader("Authorization")
		tokenString := authHeader[len(BEARER_HEADER):]

		token, err :=tokenSrv.ValidateToken(tokenString)

		if err != nil {
			c.Set("userId", token.Id)
			c.Next()
		}

		log.Println(err)
		c.AbortWithStatus(http.StatusUnauthorized)

	}
}
