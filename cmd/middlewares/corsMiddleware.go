package middlewares

import (
	"github.com/gin-gonic/gin"
	"log"
)

func CORS() gin.HandlerFunc {
	log.Println("inside here")
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "https://*")
		c.Writer.Header().Set("Access-Control-Allow-Origin", "http://*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		log.Println("cors done")
		c.Next()
	}
}
