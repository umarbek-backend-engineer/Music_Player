package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/handler"
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

	// music
	r.POST("/music", handler.Upload)
	r.GET("/music", handler.ListMusic)
	r.GET("/music/:id", handler.StreamMusic)

	// lyrics
	r.POST("/lyrics", handler.AddLyrics)
	r.GET("/lyrics/:id", handler.GetLyrics)

	// authentication
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.LogIn)
	r.POST("/auth/logout", handler.LogOut)
	r.POST("/auth/refresh")
	r.POST("/auth/validate")
	r.POST("/auth/resetpassword", handler.ResetPassword)
	r.POST("/auth/deleteaccount/:id", handler.DeleteAccount)

	return r
}
