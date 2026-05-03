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

	r.SetTrustedProxies([]string{"127.0.0.1"})

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
	authGroup := r.Group("/auth/")
	authGroup.Use((middleware.Authentication()))

	// public
	r.POST("/register", handler.Register) // working
	r.POST("/login", handler.LogIn)       // working
	r.POST("/refresh", handler.Refresh)   // working

	// protected
	authGroup.POST("/logout", handler.LogOut)                 // working
	authGroup.POST("/resetpassword", handler.ResetPassword)   // working
	authGroup.DELETE("/deleteaccount", handler.DeleteAccount) // working

	// music
	authGroup.POST("/my_music", handler.Upload)                      // working
	authGroup.GET("/my_music", handler.ListMusic)                    // working
	authGroup.GET("/my_music/stream/:music_id", handler.StreamMusic) // working

	// lyrics
	authGroup.POST("/lyrics", handler.AddLyrics)          // working
	authGroup.GET("/lyrics/:music_id", handler.GetLyrics) // working

	// social
	authGroup.GET("/users", handler.GetUsers)                      // get the existing user ID
	authGroup.GET("/user/:user_id/music", handler.GetPublicMusic) // get all public music of that user

	return r
}
