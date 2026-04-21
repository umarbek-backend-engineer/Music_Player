package handler

import (
	"github.com/umarbek-backend-engineer/Music_Player/gateway/pkg/utils"
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/modules"
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

	c.JSON(http.StatusOK, gin.H{
		"success":  true,
		"music_id": req.MusicID,
		"message":  "Lyrics Added successfully",
	})
}

func GetLyrics(c *gin.Context) {

	ctx := c.Request.Context()

	idStr := c.Param("id")

	resp, err := grpc_init.LyricsClient.GetLyrics(ctx, &pb.GetLyricsRequest{
		MusicId: idStr,
	})

	if err != nil {
		utils.Error(c, "Failed to get lyrics from lyrics-service", http.StatusInternalServerError, err)
		return
	}

	segments := make([]modules.Segment, 0, len(resp.Lyrics))

	for _, i := range resp.Lyrics {
		segments = append(segments, modules.Segment{
			Start: float64(i.Start),
			End:   float64(i.End),
			Text:  i.Text,
		})
	}

	c.JSON(http.StatusOK, modules.Respond{
		Lyrics:   segments,
		Language: resp.Language,
	})

}
