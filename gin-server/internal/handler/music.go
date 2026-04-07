package handler

import (
	"encoding/json"
	"gin-server/internal/modules"
	"gin-server/pkg/utils"
	"io"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/streadway/amqp"
)

func UploadHandler(c *gin.Context, rabbit *modules.RabbitMQ) {

	var filename string
	uploadid := uuid.New().String()

	// get the file
	fileHeader, err := c.FormFile("file")
	if err != nil {
		utils.Error(c, "Failed to recieve the file", http.StatusBadRequest, err)
		return
	}

	filename = fileHeader.Filename
	log.Println(filename)

	// open the file
	file, err := fileHeader.Open()
	if err != nil {
		utils.Error(c, "Failed to open the file", http.StatusInternalServerError, err)
		return
	}
	defer file.Close()

	// make buffer for the passing chunks
	buffer := make([]byte, 1024*32) //32KB chunk

	for {
		//read the file with the buffer
		n, err := file.Read(buffer)
		if err == io.EOF {
			break
		}
		if err != nil {
			utils.Error(c, "Failed to read file with buffer", http.StatusInternalServerError, err)
			return
		}

		msg := modules.MusicChunk{
			UploadID:  uploadid,
			FileName:  filename,
			ChunkData: buffer[:n],
			EOF:       false,
		}

		body, err := json.Marshal(msg)
		if err != nil {
			log.Println(err)
			return
		}
		err = rabbit.Ch.Publish(
			"",
			rabbit.Q.Name,
			false,
			false,
			amqp.Publishing{
				ContentType: "application/json",
				Body:        body,
			},
		)
		if err != nil {
			log.Println(err)
			return
		}

	}

	// send the json to the client
	msg := modules.MusicChunk{
		UploadID: uploadid,
		FileName: filename,
		EOF:      true,
	}

	body, err := json.Marshal(msg)
	if err != nil {
		log.Println(err)
		return
	}
	err = rabbit.Ch.Publish(
		"",
		rabbit.Q.Name,
		false,
		false,
		amqp.Publishing{
			ContentType: "application/json",
			Body:        body,
		},
	)

	if err != nil {
		log.Println(err)
		return
	}

	c.JSON(200, msg)

}
