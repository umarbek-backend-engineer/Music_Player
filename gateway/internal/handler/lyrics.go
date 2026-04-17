package handler

import (
	"gin-server/internal/grpc_init"
	"gin-server/internal/modules"
	"gin-server/pkg/utils"
	pb "gin-server/proto/gen"
	"net/http"

	"github.com/gin-gonic/gin"
)

func AddLyrics(c *gin.Context) {
	ctx := c.Request.Context()
	var req modules.AddLyricsPayload

	if err := c.ShouldBindJSON(&req); err != nil {
		utils.Error(c, "invalid request body", http.StatusBadRequest, err)
		return
	}

	_, err := grpc_init.LyricsClient.AddLyrics(ctx, &pb.AddLyricsRequest{
		MusicId: req.MusicID,
		Text:    req.Text,
	})

	if err != nil {
		utils.Error(c, "failed to add lyrics", http.StatusBadGateway, err)
		return
	}

	lyricsRes, err := grpc_init.LyricsClient.GetLyrics(ctx, &pb.GetLyricsRequest{MusicId: req.MusicID})
	if err != nil {
		utils.Error(c, "failed to fetch lyrics", http.StatusBadGateway, err)
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"music_id": lyricsRes.MusicId,
		"text":     lyricsRes.Text,
	})

}
