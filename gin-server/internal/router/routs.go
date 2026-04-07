package router

import (
	"gin-server/internal/handler"
	"gin-server/internal/modules"

	"github.com/gin-gonic/gin"
)

func Route(rabbit *modules.RabbitMQ) *gin.Engine {

	r := gin.Default()

	r.POST("/upload/", func(c *gin.Context) {
		handler.UploadHandler(c, rabbit)
	})

	return r
}
