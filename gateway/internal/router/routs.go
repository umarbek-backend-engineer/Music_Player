package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/handler"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/middleware"
)

func Route() *gin.Engine {

	r := gin.Default()

	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"http://localhost:3000"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))
	r.GET("/", func(c *gin.Context) {
		c.JSON(200, gin.H{"status": "gateway is running"})
	})

	// initialize ratelimiter
	rl := middleware.NewRateLimiter(50, time.Minute)

	// apply globally
	r.Use(rl.GinMiddleware())

	// use authentication verifier for the rest of the routes
	authGroup := r.Group("/")
	authGroup.Use((middleware.Authentication()))

	// public
	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.LogIn)

	// protected
	authGroup.POST("/auth/logout", handler.LogOut)
	authGroup.POST("/auth/resetpassword", handler.ResetPassword)
	authGroup.POST("/auth/deleteaccount", handler.DeleteAccount)
	authGroup.POST("/auth/refresh", handler.Refresh)

	// music
	authGroup.POST("/my_music", handler.Upload)
	authGroup.GET("/my_music", handler.ListMusic)
	authGroup.GET("/my_music/:id", handler.StreamMusic)

	// lyrics
	authGroup.POST("/lyrics", handler.AddLyrics)
	authGroup.GET("/lyrics/:id", handler.GetLyrics)

	return r
}
