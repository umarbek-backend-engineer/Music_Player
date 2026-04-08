package handler

import (
	"gin-server/internal/modules"
	"gin-server/pkg/utils"
	musicpb "gin-server/proto/gen"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/streadway/amqp"
	"google.golang.org/protobuf/proto"
)

func UploadHandler(c *gin.Context, rb *modules.Rabbit) {

	// get the file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "Failed to recieve file", http.StatusBadRequest, err)
		return
	}
	// open the file
	file, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, "Failed to open file", http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// validate the file, check if the size of the file is less than 10mb and if the file is music or not
	err = utils.FileValidator(fileHeader)
	if err != nil {
		utils.Error(c, "Invalid File", http.StatusBadRequest, err)
		return
	}

	// defer rb.Conn.Close()
	// defer rb.Ch.Close()

	// make buffer for the passing chunks
	buffer := make([]byte, 1024*32)
	for {
		// read the file
		n, err := file.Read(buffer)
		if err != nil && err != io.EOF {
			utils.Error(c, "Failed to read the file", http.StatusInternalServerError, err)
			return
		}

		// detect the last chunk
		islast := err == io.EOF

		// create proto buffer chunk
		chunk := &musicpb.MusicChunk{
			Filename: fileHeader.Filename,
			Data:     buffer[:n],
			IsLast:   islast,
		}

		// marshle the chunk
		body, err := proto.Marshal(chunk)
		if err != nil {
			utils.Error(c, "Failed to Marshal the proto chunk", http.StatusInternalServerError, err)
			return
		}
		// send the chunk to rabbit
		err = rb.Ch.Publish(
			"",
			rb.Q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/x-protobuf",
				Body:        body,
			},
		)
		if err != nil {
			log.Println("Failed to publish message:", err)
			continue
		}
		if islast {
			break
		}
	}

	// send the json to the client
	c.JSON(200, gin.H{
		"Message": "Uploaded successfully",
	})
}
