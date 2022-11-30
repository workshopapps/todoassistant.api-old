package routes

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

func SubscribeRoutes(v1 *gin.RouterGroup){

	// subscribe to newsletter route
	subscribe := v1.Group("/subscribe")
	{
		subscribe.POST("",func(c *gin.Context) {
		c.String(http.StatusOK, "subscribe to our newsletter")})
	}
}
