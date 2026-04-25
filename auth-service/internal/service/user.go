package service

import (
	"context"
	"fmt"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository"
	"github.com/umarbek-backend-engineer/Music_Player/pkg/utils"
	"google.golang.org/grpc/metadata"
	"google.golang.org/protobuf/types/known/emptypb"
)

// A register function, user send request like name, email, role(default 'user'), password and the method will store this inside database
// this function will return access_token and refresh_token
func (s *Server) Register(ctx context.Context, req *pb.RegisterRequest) (*pb.RegisterResponse, error) {

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

	// returning the response
	return &pb.RegisterResponse{
		Id:       id,
		Name:     req.GetName(),
		Lastname: req.GetLastname(),
		Email:    req.GetEmail(),
		Role:     req.GetRole(),
	}, nil
}

// logIn method which will get email and password and based on that it will verify the user and save metadata wich is ip-address and user_aget
func (s *Server) Login(ctx context.Context, req *pb.LoginRequest) (*pb.AuthResponse, error) {

	// get the header information (user_agen and ip_address)
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok {
		return nil, utils.MapErrors(fmt.Errorf("Failed to get incoming request"))
	}

	// initialize variable
	user_agent := ""
	ip_address := ""

	if v := md.Get("user_agent"); len(v) > 0 {
		user_agent = v[0]
	}
	if v := md.Get("ip_address"); len(v) > 0 {
		ip_address = v[0]
	}

	// check if metadata has been sent. if no return error
	if user_agent == "" && ip_address == "" {
		return nil, utils.ErrMetaData
	}

	// get id, role, saved password where email matches from database
	id, role, dbpassword, err := repository.LogInCrud(ctx, req.Email)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// validate password
	err = utils.VerifyPassword(req.GetPassword(), dbpassword)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// generate access token
	token, err := utils.GenerateAccessJWT(id, role)
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
	err = repository.InsertRefreshToken(ctx, id, hashed_refresh_token, user_agent, ip_address)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  token,
		RefreshToken: refreshToken,
	}, nil
}

func (s *Server) Logout(ctx context.Context, req *pb.LogoutRequest) (*pb.LogoutResponse, error) {

	// get the refresh token from request and hash it so it is same as in database
	hashtoken := utils.HashToken(req.RefreshToken)

	// delete the row where the refresh token matches
	err := repository.LogoutCrud(ctx, hashtoken)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

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

	// delete accoutn crud operations, it will delete user row where id matrches
	err := repository.DeleAccountCrud(ctx, req.GetId())
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// return an error
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
