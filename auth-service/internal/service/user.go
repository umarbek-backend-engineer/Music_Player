package service

import (
	"context"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
	"github.com/umarbek-backend-engineer/Music_Player/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {

	// sending req.Password(string) to utils.PasswordHash fucntion, and get encoded password  
	hash, err := utils.PasswordHash(req.Password)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	conn, err := postgres.Connect()
	if err != nil {
		return nil, utils.MapErrors(err)
	}
	defer conn.Close(ctx)

	conn.Exec(ctx, "insert into users (name, lastname, email, password, refreshtoken) values ($1,$2,$3,$4,$5)", req.Name, req.Lastname, req.Email, req.Password, )

	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *emptypb.Empty) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
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

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.AuthResponse, error) {
	return &pb.AuthResponse{
		Token: "Token",
	}, nil
}
