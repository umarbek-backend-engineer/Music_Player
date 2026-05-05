package router

import (
	"time"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
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

	// add swagger route
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// initialize ratelimiter
	rl := middleware.NewRateLimiter(50, time.Minute)

	// apply globally
	r.Use(rl.GinMiddleware())

	// use authentication verifier for the rest of the routes
	authGroup := r.Group("/auth/")
	authGroup.Use((middleware.Authentication()))

	// public
	// Register godoc
	// @Summary Register user
	// @Description Create a new user account
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param input body pb.RegisterRequest true "User registration data"
	// @Success 200 {object} map[string]interface{}
	// @Failure 400 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Router /register [post]
	r.POST("/register", handler.Register)
	// LogIn godoc
	// @Summary Login user
	// @Description Authenticates user and sets access & refresh tokens in cookies
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param input body pb.LoginRequest true "Login credentials"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Router /login [post]
	r.POST("/login", handler.LogIn)
	// Refresh godoc
	// @Summary Refresh access token
	// @Description Generates new access and refresh tokens using refresh token from cookie
	// @Tags auth
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Failure 401 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Security CookieAuth
	// @Router /refresh [post]
	r.POST("/refresh", handler.Refresh)

	// protected
	// LogOut godoc
	// @Summary Logout user
	// @Tags auth
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Failure 401 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Security CookieAuth
	// @Router /auth/logout [post]
	authGroup.POST("/logout", handler.LogOut) // logs out
	// ResetPassword godoc
	// @Summary Reset user password
	// @Description Reset password using current and new password
	// @Tags auth
	// @Accept json
	// @Produce json
	// @Param input body pb.ResetPasswordRequest true "Reset password data"
	// @Success 200 {object} map[string]string
	// @Failure 400 {object} ErrorResponse
	// @Failure 401 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Security CookieAuth
	// @Router /auth/resetpassword [post]
	authGroup.POST("/resetpassword", handler.ResetPassword) // resets users password
	// DeleteAccount godoc
	// @Summary Delete user account
	// @Description Deletes the authenticated user's account
	// @Tags auth
	// @Produce json
	// @Success 200 {object} map[string]string
	// @Failure 401 {object} ErrorResponse
	// @Failure 502 {object} ErrorResponse
	// @Security CookieAuth
	// @Router /auth/deleteaccount [delete]
	authGroup.DELETE("/deleteaccount", handler.DeleteAccount) // deletes  account

	// music
	authGroup.POST("/my_music", handler.Upload)   // upload the music
	authGroup.GET("/my_music", handler.ListMusic) // lists the users music
	// StreamMusic godoc
	// @Summary Stream music
	// @Tags music
	// @Param music_id path string true "Music ID"
	// @Success 200 {file} binary
	// @Router /auth/my_music/stream/{music_id} [get]
	authGroup.GET("/my_music/stream/:music_id", handler.StreamMusic)            // play the music
	authGroup.PATCH("/my_music/:music_id/visibility", handler.ChangeVisibilaty) // change the music visibility

	// lyrics
	authGroup.POST("/lyrics", handler.AddLyrics)          // generate lyrics
	authGroup.GET("/lyrics/:music_id", handler.GetLyrics) // get the generated music lyrics

	// social
	authGroup.GET("/users", handler.GetUsers)                     // get the existing user ID
	authGroup.GET("/user/:user_id/music", handler.GetPublicMusic) // get all public music of that user

	return r
}
