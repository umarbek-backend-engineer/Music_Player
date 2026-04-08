package service

import (
	"context"
	"fmt"
	"log"
	"music-service/internal/config"
	rabbitmq "music-service/internal/rabbit-mq"
	"music-service/internal/repository"
	"music-service/pkg/utils"
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
		return utils.MapErrors(err)
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
		return utils.MapErrors(err)
	}

	forever := make(chan struct{})

	go func() {

		var files = make(map[string]*os.File)

		for d := range msg {

			// declare proto chunk
			var chunk pb.UploadMusicChunk

			err = proto.Unmarshal(d.Body, &chunk)
			if err != nil {
				utils.MapErrors(err)
				continue
			}

			filename = chunk.Filename

			filePath := fmt.Sprintf("./%s/%s", cgf.StoragePath, chunk.Filename)

			file := files[chunk.Filename]
			// 🟢 create file if not exists
			if file == nil {
				file, err = os.Create(filePath)
				if err != nil {
					utils.MapErrors(err)
					continue
				}
				files[chunk.Filename] = file
			}

			// 🟢 write chunk data
			_, err := file.Write(chunk.Data)
			if err != nil {
				utils.MapErrors(err)
				continue
			}

			// 🔴 if last chunk → close file
			if chunk.IsLast {
				log.Println("Finished file:", chunk.Filename)

				// saving the music metadata in database
				err = repository.UploadMusicDBHandler(context.Background(), filename, filePath)
				if err != nil {
					utils.MapErrors(err)
					continue
				}
				file.Close()
				delete(files, filename)
			}

		}
	}()
	<-forever
	return nil

}

func (s *Server) ListMusic(ctx context.Context, req *pb.Empty) (*pb.ListResponse, error) {
	// get the music from db
	musics, err := repository.ListMusicDB(ctx)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	return &pb.ListResponse{
		Songs: musics,
	}, nil

}
