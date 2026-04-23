package service

import (
	"context"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository"
	"github.com/umarbek-backend-engineer/Music_Player/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.AuthResponse, error) {

	// sending req.Password(string) to utils.PasswordHash fucntion, and get encoded password
	hash, err := utils.PasswordHash(req.Password)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// change the value of the req.Password so that when I am saving in database I am only passing
	req.Password = hash

	// the register crud function, it will save request data into database
	id, err := repository.RegisterDBCrud(ctx, req)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// generate access token
	token, err := utils.GenerateAccessJWT(id, req.Role)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  token,
		RefreshRokne: "Refresh_token",
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {
	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshRokne: "Refresh_token",
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
	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshRokne: "Refresh_token",
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshRokne: "Refresh_token",
	}, nil
}
