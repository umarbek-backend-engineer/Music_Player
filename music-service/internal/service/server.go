package service

import pb "github.com/umarbek-backend-engineer/Music_Player/music-service/github.com/umarbek-backend-engineer/Music_Player/music-service/proto/gen"

type Server struct {
	pb.UnimplementedMusicServiceServer
}
