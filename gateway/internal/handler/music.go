package handler

import (
	"bytes"
	grpc_init "gin-server/internal/grpc"
	"gin-server/pkg/utils"
	pb "gin-server/proto/gen"
	"io"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// func UploadHandler(c *gin.Context, rb *modules.Rabbit) {

// 	// get the file
// 	fileHeader, err := c.FormFile("file")
// 	if err != nil {
// 		utils.Error(c, "Failed to recieve file", http.StatusBadRequest, err)
// 		return
// 	}
// 	// open the file
// 	file, err := fileHeader.Open()
// 	if err != nil {
// 		utils.Error(c, "Failed to open file", http.StatusInternalServerError, err)
// 		return
// 	}
// 	defer file.Close()

// 	// validate the file, check if the size of the file is less than 10mb and if the file is music or not
// 	err = utils.FileValidator(fileHeader)
// 	if err != nil {
// 		utils.Error(c, "Invalid File", http.StatusBadRequest, err)
// 		return
// 	}

// 	// defer rb.Conn.Close()
// 	// defer rb.Ch.Close()

// 	// make buffer for the passing chunks
// 	buffer := make([]byte, 1024*32)
// 	for {
// 		// read the file
// 		n, err := file.Read(buffer)
// 		if err != nil && err != io.EOF {
// 			utils.Error(c, "Failed to read the file", http.StatusInternalServerError, err)
// 			return
// 		}

// 		// detect the last chunk
// 		islast := err == io.EOF

// 		// create proto buffer chunk
// 		chunk := &musicpb.UploadMusicChunks{
// 			Filename: fileHeader.Filename,
// 			Data:     buffer[:n],
// 			IsLast:   islast,
// 		}

// 		// marshle the chunk
// 		body, err := proto.Marshal(chunk)
// 		if err != nil {
// 			utils.Error(c, "Failed to Marshal the proto chunk", http.StatusInternalServerError, err)
// 			return
// 		}
// 		// send the chunk to rabbit
// 		err = rb.Ch.Publish(
// 			"",
// 			rb.Q.Name,
// 			false,
// 			false,
// 			amqp.Publishing{
// 				ContentType: "application/x-protobuf",
// 				Body:        body,
// 			},
// 		)
// 		if err != nil {
// 			log.Println("Failed to publish message:", err)
// 			continue
// 		}
// 		if islast {
// 			break
// 		}
// 	}

// 	// send the json to the client
// 	c.JSON(200, gin.H{
// 		"Message": "Uploaded successfully",
// 	})
// }

func Upload(c *gin.Context) {

	ctx := c.Request.Context()

	// get the music from requst
	fileheader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "Failed to recieve file", http.StatusBadRequest, err)
		return
	}

	// open the file
	file, err := fileheader.Open()
	if err != nil {
		utils.Error(c, "Failed to open file", http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// connecting to grpc and  getting stream data
	stream, err := grpc_init.MusicClient.UploadMusic(ctx)
	if err != nil {
		utils.Error(c, "Failed to create grpc client for upload Music", http.StatusInternalServerError, err)
		return
	}

	// creating buffer size of 32kb to send data through
	buffer := make([]byte, 1024*32)

	firstChunk := true
	for {
		// reading the file chunk bu chunk of 32kb size
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Error(c, "Failed to read the file", http.StatusInternalServerError, err)
			return
		}

		chunk := &pb.UploadMusicChunks{
			Data: buffer[:n],
		}

		//send  the file only once
		if firstChunk {
			chunk.Filename = fileheader.Filename
			firstChunk = false
		}

		// send the read chunk
		err = stream.Send(chunk)
		if err != nil {
			utils.Error(c, "Failed to send the buffer", http.StatusInternalServerError, err)
			return
		}
	}
	// close the stream and recieve the  response
	res, err := stream.CloseAndRecv()
	if err != nil {
		utils.Error(c, "Failed to close and recieve the response", http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, res)
}

func ListMusic(c *gin.Context) {
	// context of the browser
	ctx := c.Request.Context()

	res, err := grpc_init.MusicClient.ListMusic(ctx, &pb.Empty{})
	if err != nil {
		utils.Error(c, "Failed to list music", http.StatusInternalServerError, err)
		return
	}

	c.JSON(200, res.Songs)

}

func StreamMusic(c *gin.Context) {
	ctx := c.Request.Context()
	var filename string

	id := c.Param("id")

	stream, err := grpc_init.MusicClient.StreamMusic(ctx, &pb.StreamRequest{Id: id})
	if err != nil {
		utils.Error(c, "Failed to get stream of music", http.StatusInternalServerError, err)
		return
	}

	// setting header so that browser know that the file that I am passing is music
	c.Header("Content-Type", "audio/mpeg")
	c.Header("Accept-Ranges", "bytes")
	c.Header("Content-Disposition", "inline")
	c.Status(http.StatusOK)

	var buffer bytes.Buffer

	for {
		chunk, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Error(c, "Failed to get chunk of music", http.StatusInternalServerError, err)
			return
		}

		filename = chunk.Name

		// // push the chunk of bytes to the browser
		// _, err = c.Writer.Write(chunk.Content)
		// if err != nil {
		// 	utils.Error(c, "Failed to write music in brower", http.StatusInternalServerError, err)
		// 	return
		// }

		// // Forces data to go immediately
		// c.Writer.Flush()

		// collects the bytes
		buffer.Write(chunk.Content)

	}

	// ServeContent handles Range requests automatically — lets the browser seek
	// to any position by returning only the requested bytes instead of the full file.
	http.ServeContent(c.Writer, c.Request, filename, time.Time{}, bytes.NewReader(buffer.Bytes()))

}
