package router

import (
	"gin-server/internal/handler"
	"gin-server/internal/modules"

	"github.com/gin-gonic/gin"
)

func Route(rb *modules.Rabbit) *gin.Engine {

	r := gin.Default()

	r.POST("/upload/", func(cxt *gin.Context) {
		handler.UploadHandler(cxt, rb)
	})

	return r
}
