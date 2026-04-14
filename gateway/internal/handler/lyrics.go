package handler

import (
	"gin-server/internal/grpc_init"
	"gin-server/pkg/utils"
	pb "gin-server/proto/gen"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddLyrics(c *gin.Context) {
	ctx := c.Request.Context()

	_, err := grpc_init.LyricsClient.AddLyrics(ctx, &pb.AddLyricsRequest{
		MusicId: "fa28f3c1-dd34-446b-be18-7c2563463406",
		Text:    "Amy Winehouse - Back To Black.mp3",
	})

	if err != nil {
		utils.Error(c, "failed to add lyrics", http.StatusBadGateway, err)
		return
	}

}
