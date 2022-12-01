package vaMiddleware

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	tokenservice "test-va/internals/service/tokenService"

	"github.com/gin-gonic/gin"
)

type vaMiddleWare struct {
	tokenSrv tokenservice.TokenSrv
}

func NewVaMiddleWare(tokenSrv tokenservice.TokenSrv) *vaMiddleWare {
	return &vaMiddleWare{tokenSrv: tokenSrv}
}

func (v *vaMiddleWare) MapVAToReq(c *gin.Context) {
	// const BEARER_HEADER = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid Token")
		return
	}
	auth := strings.Split(authHeader, " ")

	token, err := v.tokenSrv.ValidateToken(auth[1])
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("invalid Token: %v", err))

	}

	if token.Status != "VA" && token.Status != "MASTER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "You are not Authorized to access this resource")
		return
	}
	c.Set("id", token.Id)
	c.Set("status", token.Status)
	c.Set("email", token.Email)
	c.Next()
}

func (v *vaMiddleWare) MapMasterToReq(c *gin.Context) {
	log.Println("id", c.Request.Context().Value("id"))
	// const BEARER_HEADER = "Bearer "
	authHeader := c.GetHeader("Authorization")
	if authHeader == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "Invalid Token")
		return
	}
	auth := strings.Split(authHeader, " ")

	token, err := v.tokenSrv.ValidateToken(auth[1])
	if err != nil {
		log.Println(err)
		c.AbortWithStatusJSON(http.StatusUnauthorized, fmt.Sprintf("invalid Token: %v", err))

	}

	if token.Status != "MASTER" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, "You are not Authorized to access this resource")
		return
	}
	c.Set("id", token.Id)
	c.Set("status", token.Status)
	c.Set("email", token.Email)
	c.Next()
}
