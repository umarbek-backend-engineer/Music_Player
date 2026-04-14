package service

import (
	"bytes"
	"context"
	"log"
	"lyrics-service/pkg/utils"
	pb "lyrics-service/proto/gen"

	"google.golang.org/grpc"
)

func (s *Server) AddLyrics(ctx context.Context, req *pb.AddLyricsRequest) (*pb.Empty, error) {

	conn, err := grpc.Dial("music-service:50051", grpc.WithInsecure())
	if err != nil {
		return nil, err
	}
	defer conn.Close()

	client := pb.NewMusicServiceClient(conn)

	stream, err := client.StreamMusic(ctx, &pb.StreamRequest{
		Id: req.MusicId,
	})

	if err != nil {
		return nil, utils.MapError(err)
	}
	var filename string
	var buffer bytes.Buffer
	for {
		res, err := stream.Recv()
		if err != nil {
			log.Println("Error recieving the music chunks: ", err)
			break
		}

		if res.Name != "" {
			filename = req.Text
		}

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
