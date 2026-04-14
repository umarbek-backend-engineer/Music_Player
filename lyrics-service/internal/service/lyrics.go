package service

import (
	"bytes"
	"context"
	"log"
	"lyrics-service/internal/repository"
	"lyrics-service/pkg/utils"
	pb "lyrics-service/proto/gen"

	"google.golang.org/grpc"
)

// This is add lyrics rpc
// first it will ceck if the music exists by check the same names in data base. If yes the function will not create another lyrics, else it will.
// if there is not same music name in data base, it will connect to music-service and pull the stream music rpc
// the rpc will return music bytes and the bytes are send to wisper
// recieved lyrics with timestamp are saved in db

func (s *Server) AddLyrics(ctx context.Context, req *pb.AddLyricsRequest) (*pb.Empty, error) {

	exists, err := repository.Is_music_lyric_exists(ctx, req.Text)
	if err != nil {
		return nil, err
	}
	if exists {
		return &pb.Empty{}, nil
	}

	// connecting  to grpc server to get the music it self
	conn, err := grpc.Dial("music-service:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	// requesting for stream rpc
	stream, err := client.StreamMusic(ctx, &pb.StreamRequest{
		Id: req.MusicId,
	})

	if err != nil {
		return nil, utils.MapError(err)
	}
	var filename string
	var buffer bytes.Buffer
	for {
		//recieving the music chunks
		res, err := stream.Recv()
		if err != nil {
			log.Println("Error recieving the music chunks: ", err)
			break
		}

		if res.Name != "" {
			filename = req.Text
		}

		// writing in buffer
		_, err = buffer.Write(res.Content)
		if err != nil {
			log.Println("Error writing the file: ", err)
			break
		}
	}

	text, err := utils.SendToWisper(buffer.Bytes(), filename)
	if err != nil {
		return nil, utils.MapError(err)
	}

	log.Println("lyrics: ", text)

	return &pb.Empty{}, nil
}

func (s *Server) GetLyrics(ctx context.Context, req *pb.GetLyricsRequest) (*pb.LyricsResponse, error) {

	return &pb.LyricsResponse{
		MusicId: "music_Id",
		Text:    "Music Lyrics",
	}, nil
}
