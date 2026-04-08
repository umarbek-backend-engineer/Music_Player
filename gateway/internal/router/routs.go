package router

import (
	"gin-server/internal/handler"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	r := gin.Default()

	r.POST("/music/", handler.Upload)
	r.GET("/music/", handler.ListMusic)
	r.GET("/music/:id", handler.StreamMusic)

	return r
}
