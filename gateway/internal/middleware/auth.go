package auth

import (
	"log"

	"github.com/gin-gonic/gin"
)

func authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		log.Println("authentication middleware start checking")
		c.Next()
	}
}
