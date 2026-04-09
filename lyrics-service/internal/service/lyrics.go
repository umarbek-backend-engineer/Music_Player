package service

import (
	"context"
	pb "lyrics-service/proto/gen"
)

func (s *Server) AddLyrics(ctx context.Context, req *pb.AddLyricsRequest) (*pb.Empty, error) {

	return &pb.Empty{}, nil
}

func (s *Server) GetLyrics(ctx context.Context, req *pb.GetLyricsRequest) (*pb.LyricsResponse, error) {

	return &pb.LyricsResponse{
		MusicId: "music_Id",
		Text:    "Music Lyrics",
	}, nil
}
