package service

import (
	"io"
	"log"
	"music-service/internal/config"
	"music-service/internal/repository"
	"music-service/pkg/utils"
	pb "music-service/proto/gen"
	"os"

	"github.com/google/uuid"
)

func (s *Server) UploadMusic(stream pb.MusicService_UploadMusicServer) error {

	var filename string
	var outFile *os.File
	var id uuid.UUID
	uploadDir := config.Load().StoragePath

	defer func() {
		if outFile != nil {
			outFile.Close()
		}
	}()

	// ensure upload directory exists
	if _, err := os.Stat(uploadDir); os.IsNotExist(err) {
		os.MkdirAll(uploadDir, os.ModePerm)
	}

	for {
		req, err := stream.Recv()
		if err != nil {
			// Save metadata into DB (filename + path)
			if err == io.EOF {
				id, err = repository.UploadMusicDBHandler(
					stream.Context(),
					filename,
					uploadDir+filename+".mp3",
				)
				if err != nil {
					return utils.MapErrors(err)
				}

				// Send final response to client and close stream
				return stream.SendAndClose(&pb.UploadResponse{
					Id:       id.String(),
					Filename: filename,
				})
			}
			log.Println(err)
			return utils.MapErrors(err)
		}

		if outFile == nil || outFile.Fd() == 0 {
			log.Println("file is nil")
			filename = req.Filename
			log.Println("Received filename:", filename)
			if filename == "" {
				filename = uuid.New().String() // fallback unique name
				log.Println("Using fallback filename:", filename)
			}

			// create a file on desk
			var err error
			outFile, err = os.Create(uploadDir + filename + ".mp3")
			if err != nil {
				log.Println("Error creating file:", err)
				return utils.MapErrors(err)	
			}
			log.Println("Created file:", uploadDir+filename)
		}

		// write current chunk
		_, err = outFile.Write(req.Content)
		if err != nil {
			return utils.MapErrors(err)
		}
	}
}

// func (s *Server) ListMusic(ctx context.Context, req *pb.Empty) (*pb.ListResponse, error) {
// 	return &pb.ListResponse{
// 		Songs: []*pb.MusicItem{
// 			&pb.MusicItem{
// 				Id:       "id_123",
// 				Filename: "Hello",
// 			}},
// 	}, nil
// }

// func (s *Server) StreamMusic(ctx context.Context, req *pb.StreamRequest) (*pb.MusicService_StreamMusicServer, error) {
// 	return &pb.MusicService_StreamMusicServer{

// 	}, nil
// }
