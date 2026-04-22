package service

import (
	"context"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.Empty) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return &pb.ValidateResponse{
		UserId: "user_001",
		Role:   "Admin",
	}, nil
}

func (s *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*pb.Empty, error) {
	return &pb.Empty{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}
