package service

import (
	"io"
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

	for {
		req, err := stream.Recv()
		if err != nil {
			// Save metadata into DB (filename + path)
			if err == io.EOF {
				id, err = repository.UploadMusicDBHandler(
					stream.Context(),
					req.Filename,
					uploadDir+filename,
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

			return utils.MapErrors(err)
		}

		if outFile == nil {
			filename := req.Filename

			// create a file on desk
			var err error
			outFile, err = os.Create(uploadDir + filename)
			if err != nil {
				return utils.MapErrors(err)
			}
		}

		// write current chunk
		_, err = outFile.Write(req.Content)
		if err != nil {
			return utils.MapErrors(err)
		}
	}
}
