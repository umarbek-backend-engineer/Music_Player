package service

import (
	"log"
	"music-service/internal/config"
	rabbitmq "music-service/internal/rabbit-mq"
	pb "music-service/proto/gen"
	"os"

	"google.golang.org/protobuf/proto"
)

func StartConsumer() error {
	cgf := config.Load()

	var filename string
	// cgf := config.Load()

	rabbit, err := rabbitmq.Connect()
	if err != nil {
		return err
	}
	// defer rabbit.Conn.Close()
	// defer rabbit.Ch.Close()

	msg, err := rabbit.Ch.Consume(
		rabbit.Q.Name,
		"",
		true,
		false,
		false,
		false,
		nil,
	)
	if err != nil {
		return err
	}

	forever := make(chan struct{})

	go func() {

		var files = make(map[string]*os.File)

		for d := range msg {

			// declare proto chunk
			var chunk pb.UploadMusicChunk

			err = proto.Unmarshal(d.Body, &chunk)
			if err != nil {
				log.Println("Failed to decode json:", err)
				continue
			}

			file := files[chunk.Filename]
			// 🟢 create file if not exists
			if file == nil {
				file, err = os.Create("./"+ cgf.StoragePath +"/" + chunk.Filename)
				if err != nil {
					log.Println("Failed to create file:", err)
					continue
				}
			}

			// 🟢 write chunk data
			_, err := file.Write(chunk.Data)
			if err != nil {
				log.Println("Failed to write chunk:", err)
				continue
			}

			// 🔴 if last chunk → close file
			if chunk.IsLast {
				log.Println("Finished file:", chunk.Filename)

				file.Close()
				file = nil
			}

		}
	}()

	log.Println("Consumer started...")
	<-forever
	log.Println("Recieved file ", filename)
	return nil

}
