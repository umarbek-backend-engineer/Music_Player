package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/handler"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/middleware"
	auth "github.com/umarbek-backend-engineer/Music_Player/gateway/internal/middleware"
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

	// initialize ratelimiter
	rl := auth.NewRateLimiter(50, time.Minute)

	// apply globally
	r.Use(rl.Middleware())

	// use authentication verifier for the rest of the routes
	authGroup := r.Group("/")
	authGroup.Use((middleware.Authentication()))

	// user - auth-service

	r.POST("/auth/register", handler.Register)
	r.POST("/auth/login", handler.LogIn)
	authGroup.POST("/auth/logout", handler.LogOut)
	authGroup.POST("/auth/resetpassword", handler.ResetPassword)
	authGroup.POST("/auth/deleteaccount/:id", handler.DeleteAccount)

	// music
	authGroup.POST("/music", handler.Upload)
	authGroup.GET("/music", handler.ListMusic)
	authGroup.GET("/music/:id", handler.StreamMusic)

	// lyrics
	authGroup.POST("/lyrics", handler.AddLyrics)
	authGroup.GET("/lyrics/:id", handler.GetLyrics)

	authGroup.POST("/auth/refresh", handler.Refresh)
	return r
}
