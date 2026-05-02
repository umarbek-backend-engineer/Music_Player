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
	"google.golang.org/protobuf/types/known/emptypb"
)

// the function will create and upload the recieved audio file and store the metadata inside the database
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

	userID := md.Get("user-id")
	// validate userID
	if len(userID) == 0 {
		return utils.MapErrors(errors.New("Missing id in metadata"))
	}

	user_id := userID[0]

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
			filePath = filepath.Join(cgf.StoragePath, user_id+"."+req.Title)

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
	err = repository.UploadMusicDBHandler(stream.Context(), user_id, title, filePath)
	if err != nil {
		return utils.MapErrors(err)
	}
	// send response
	return stream.SendAndClose(&pb.UploadMusicResponse{
		Status: title + " - File uploaded successfully",
	})
}

// this method will list every audio that is stored in db with user_id saved in db
func (s *Server) ListMusic(ctx context.Context, req *emptypb.Empty) (*pb.ListResponse, error) {

	// extract the id form the incoming context (metadata)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapErrors(errors.New("Missing id in metadata"))
	}
	userID := md.Get("user-id")
	// validate userID
	if len(userID) == 0 {
		return nil, utils.MapErrors(errors.New("Missing id in metadata"))
	}

	user_id := userID[0]

	// get the music from db
	musics, err := repository.ListMusicDB(ctx, user_id)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// send back the response
	return &pb.ListResponse{
		Songs: musics,
	}, nil

}

func (s *Server) StreamMusic(req *pb.StreamRequest, stream pb.MusicService_StreamMusicServer) error {

	ctx := stream.Context()

	// extract the id form the incoming context (metadata)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return utils.MapErrors(errors.New("Missing id in metadata"))
	}
	userID := md.Get("user-id")
	// validate userID
	if len(userID) == 0 {
		return utils.MapErrors(errors.New("Missing id in metadata"))
	}

	user_id := userID[0]

	// Build file path (you may map ID → filename later)
	music, err := repository.GetMusicIndoFromDB_on_ID(ctx, user_id, req.Id)
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
