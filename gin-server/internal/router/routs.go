package router

import (
	"gin-server/internal/handler"

	"github.com/gin-gonic/gin"
)

func Route() *gin.Engine {

	r := gin.Default()

	r.POST("/upload/", handler.UploadHandler)

	return r
}
