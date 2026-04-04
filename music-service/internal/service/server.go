package service

import pb "music-service/proto/gen"

type Server struct {
	pb.UnimplementedMusicServiceServer
}