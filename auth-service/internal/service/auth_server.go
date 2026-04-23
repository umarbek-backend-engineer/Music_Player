package service

import (
	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
)

type Server struct {
	pb.UnimplementedAuthServiceServer
}
