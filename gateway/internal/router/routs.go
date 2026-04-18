package router

import (
	"gin-server/internal/handler"
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: false,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "gateway is running"})
	})

	r.POST("/music", handler.Upload)
	r.GET("/music", handler.ListMusic)
	r.GET("/music/:id", handler.StreamMusic)

	r.POST("/lyrics", handler.AddLyrics)
	r.GET("/lyrics/:id", handler.GetLyrics)

	return r
}
