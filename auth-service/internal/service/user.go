package service

import (
	"context"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
	"github.com/umarbek-backend-engineer/Music_Player/pkg/utils"
	"google.golang.org/protobuf/types/known/emptypb"
)

// A register function, user send request like name, email, role(default 'user'), password and the method will store this inside database
// this function will return access_token and refresh_token
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

	// generate refresh tokenm
	refreshToken, err := utils.GenerateRefreshTokne()
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// hash the token befor saving it inside database for more security
	hashed_refresh_token := utils.HashToken(token)

	// save refresh token in data base
	err = repository.InsertRefreshToken(ctx, id, hashed_refresh_token)

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshToken: "Refresh_token",
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	refresh_token := req.GetRefreshToken()

	// create the database client and close it after its usage
	conn, err := postgres.Connect()
	if err != nil {
		return nil, err
	}
	defer conn.Close(ctx)

	// check if it exists and valid
	tag, err := conn.Exec(ctx, "select * from sessions where refresh_token = $1", refresh_token)

	// return response
	return &pb.LogoutResponse{
		Success: true,
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {
	return &pb.ValidateResponse{
		UserId: "user_001",
		Role:   "Admin",
	}, nil
}

func (s *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	return &emptypb.Empty{}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.AuthResponse, error) {
	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshToken: "Refresh_token",
	}, nil
}

func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {
	// returning the response
	return &pb.AuthResponse{
		AccessToken:  "token",
		RefreshToken: "Refresh_token",
	}, nil
}
