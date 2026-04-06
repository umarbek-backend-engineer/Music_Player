package handler

import (
	"context"
	"gin-server/pkg/utils"
	pb "gin-server/proto/gen"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func UploadHandler(c *gin.Context) {
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "Failed to load mp3 file", http.StatusBadRequest, err)
		return
	}
	file, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, "Failed to open mp3 file", http.StatusBadRequest, err)
		return
	}
	defer file.Close()

	// connect to gRPC
	conn, err := grpc.Dial("localhost:50051", grpc.WithInsecure())
	if err != nil {
		utils.Error(c, "Failed to connect to grpc server", http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	stream, err := client.UploadMusic(context.Background())
	if err != nil {
		utils.Error(c, "Failed to create gprc client", http.StatusInternalServerError, err)
		return
	}
	log.Println(fileHeader.Filename)

	buffer := make([]byte, 32*1024) //32KB chunk
	for {
		n, err := file.Read(buffer)
		if err != io.EOF {
			break
		}
		if err != nil {
			utils.Error(c, "Failed to read chunks", http.StatusInternalServerError, err)
		}

		err = stream.Send(&pb.UploadRequest{
			Filename: fileHeader.Filename,
			Content:  buffer[:n],
		})
		if err != nil {
			utils.Error(c, "Failed to send chunks", http.StatusInternalServerError, err)
		}
	}

	resp, err := stream.CloseAndRecv()
	if err != nil {
		utils.Error(c, "Failed to close and receive", http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, gin.H{
		"id":       resp.Id,
		"filename": resp.Filename,
	})

}
