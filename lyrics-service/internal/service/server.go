package service

import pb "github.com/umarbek-backend-engineer/Music_Player/lyrics-service/github.com/umarbek-backend-engineer/Music_Player/lyrics-service/proto/gen"

type Server struct {
	pb.UnimplementedLyricsServiceServer
}
