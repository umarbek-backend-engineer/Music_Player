package service

import (
	"context"
	"errors"
	"fmt"
	"io"

	"os"
	"path/filepath"

	pb "github.com/umarbek-backend-engineer/Music_Player/music-service/github.com/umarbek-backend-engineer/Music_Player/music-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/config"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/repository"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/pkg/utils"
	"google.golang.org/grpc/metadata"
)

// func StartConsumer() error {
// 	cgf := config.Load()

// 	var filename string
// 	// cgf := config.Load()

// 	rabbit, err := rabbitmq.Connect()
// 	if err != nil {
// 		return utils.MapErrors(err)
// 	}
// 	// defer rabbit.Conn.Close()
// 	// defer rabbit.Ch.Close()

// 	msg, err := rabbit.Ch.Consume(
// 		rabbit.Q.Name,
// 		"",
// 		true,
// 		false,
// 		false,
// 		false,
// 		nil,
// 	)
// 	if err != nil {
// 		return utils.MapErrors(err)
// 	}

// 	forever := make(chan struct{})

// 	go func() {

// 		var files = make(map[string]*os.File)

// 		for d := range msg {

// 			// declare proto chunk
// 			var chunk pb.UploadMusicChunk

// 			err = proto.Unmarshal(d.Body, &chunk)
// 			if err != nil {
// 				utils.MapErrors(err)
// 				continue
// 			}

// 			filename = chunk.Filename

// 			filePath := fmt.Sprintf("./%s/%s", cgf.StoragePath, chunk.Filename)

// 			file := files[chunk.Filename]
// 			// 🟢 create file if not exists
// 			if file == nil {
// 				file, err = os.Create(filePath)
// 				if err != nil {
// 					utils.MapErrors(err)
// 					continue
// 				}
// 				files[chunk.Filename] = file
// 			}

// 			// 🟢 write chunk data
// 			_, err := file.Write(chunk.Data)
// 			if err != nil {
// 				utils.MapErrors(err)
// 				continue
// 			}

// 			// 🔴 if last chunk → close file
// 			if chunk.IsLast {
// 				log.Println("Finished file:", chunk.Filename)

// 				// saving the music metadata in database
// 				err = repository.UploadMusicDBHandler(context.Background(), filename, filePath)
// 				if err != nil {
// 					utils.MapErrors(err)
// 					continue
// 				}
// 				file.Close()
// 				delete(files, filename)
// 			}

// 		}
// 	}()
// 	<-forever
// 	return nil

// }

func (s *Server) UploadMusic(stream pb.MusicService_UploadMusicServer) error {
	cgf := config.Load()
	var file *os.File
	var title string
	var filePath string

	// extract the id form the incoming context (metadata)
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return utils.MapErrors(errors.New("Missing id in metadata"))
	}
	userID := md.Get("user_id")[0]
	if len(userID) == 0 {
		return utils.MapErrors(errors.New("Missing id in metadata"))
	}

	// the loop will run till the file chunks passed completely
	for {
		req, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			return utils.MapErrors(err)
		}
		if req.Title != "" {
			title = req.Title
		}
		// Create file on first chunk
		if file == nil {
			if title == "" {
				return utils.MapErrors(fmt.Errorf("Empty file name"))
			}
			err = os.MkdirAll(cgf.StoragePath, 0755)
			if err != nil {
				return utils.MapErrors(err)
			}
			filePath = filepath.Join(cgf.StoragePath, req.Title, ".", userID)

			// create the file inside the storage
			file, err = os.Create(filePath)
			if err != nil {
				return utils.MapErrors(err)
			}
		}
		//write the chunk iside the created file
		if len(req.Data) > 0 {
			_, err = file.Write(req.Data)
			if err != nil {
				return utils.MapErrors(err)
			}
		}

	}
	// safty check
	if file == nil {
		return utils.MapErrors(fmt.Errorf("no file received"))
	}
	// Close file explicitly
	err := file.Close()
	if err != nil {
		return utils.MapErrors(err)
	}
	// save metadata in database
	err = repository.UploadMusicDBHandler(stream.Context(), userID, title, filePath)
	if err != nil {
		return utils.MapErrors(err)
	}
	// send response
	return stream.SendAndClose(&pb.UploadMusicResponse{
		Status: title + " - File uploaded successfully",
	})
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

func (s *Server) StreamMusic(req *pb.StreamRequest, stream pb.MusicService_StreamMusicServer) error {

	ctx := stream.Context()

	// Build file path (you may map ID → filename later)
	music, err := repository.GetMusicIndoFromDB_on_ID(ctx, req.Id)
	if err != nil {
		return utils.MapErrors(err)
	}

	// open the file
	file, err := os.Open(music.FilePath)
	if err != nil {
		return utils.MapErrors(err)
	}

	// create a buffer, it is something like container with the size of 32kb using which I will pass data to gatewway
	buffer := make([]byte, 1024*32)

	for {
		n, err := file.Read(buffer)

		if err == io.EOF {
			break
		}

		if err != nil {
			return utils.MapErrors(err)
		}

		// send the  chunk of music
		err = stream.Send(&pb.MusicChunk{
			Title:   music.FileName,
			Content: buffer[:n],
		})

		if err != nil {
			return utils.MapErrors(err)
		}
	}
	return nil
}
