package handler

import (
	"context"
	"gin-server/internal/config"
	"gin-server/pkg/utils"
	pb "gin-server/proto/gen"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
)

func UploadHandler(c *gin.Context) {

	var filename string
	cgf := config.Load()

	// get the file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "Failed to recieve the file", http.StatusBadRequest, err)
		return
	}

	filename = fileHeader.Filename

	// open the file
	file, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, "Failed to open the file", http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// connect to gRPC

	conn, err := grpc.Dial(cgf.Api_Host+":"+cgf.GRPC_PORT, grpc.WithInsecure())
	if err != nil {
		utils.Error(c, "Failed to connect gRPC server", http.StatusInternalServerError, err)
		return
	}
	defer conn.Close()
	// create service client

	client := pb.NewMusicServiceClient(conn)
	// get the stream from the client

	stream, err := client.UploadMusic(context.Background())
	if err != nil {
		utils.Error(c, "Failed to create stream", http.StatusInternalServerError, err)
		return
	}
	// make buffer for the passing chunks
	buffer := make([]byte, 1024*32) //32KB chunk

	for {
		//read the file with the buffer
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Error(c, "Failed to read chunk from the recieved file", http.StatusInternalServerError, err)
			return
		}

		// send the read buffer
		err = stream.Send(&pb.UploadRequest{
			Filename: filename,
			Content:  buffer[:n],
		})
		if err != nil {
			utils.Error(c, "Failed to read chunk from the recieved file", http.StatusInternalServerError, err)
			return
		}
	}
	// close the stream
	res, err := stream.CloseAndRecv()
	if err != nil {
		utils.Error(c, "Failed to Close and revieve response", http.StatusInternalServerError, err)
		return
	}
	// send the json to the client
	c.JSON(200, gin.H{
		"id":       res.Id,
		"filename": res.Filename,
	})
}
