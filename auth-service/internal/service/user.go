package service

import (
	"context"
	"fmt"

	"github.com/golang-jwt/jwt/v5"
	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/config"
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

	if v := md.Get("md-user-agent"); len(v) > 0 {
		user_agent = v[0]
	}
	if v := md.Get("md-ip-address"); len(v) > 0 {
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

// in this method the gateway send id of the user and based on id it will delete the row  in database
func (s *Server) DeleteAccount(ctx context.Context, req *pb.DeleteAccountRequest) (*emptypb.Empty, error) {
	// delete accoutn crud operations, it will delete user row where id matrches
	err := repository.DeleAccountCrud(ctx, req.GetId())
	if err != nil {
		return nil, utils.MapErrors(err)
	}
	// return an error
	return &emptypb.Empty{}, nil
}

// The method refreshed the token expiration date in sessions table
func (s *Server) Refresh(ctx context.Context, req *pb.RefreshRequest) (*pb.AuthResponse, error) {

	// Get refresh and hash the token
	hashed_refresh_token := utils.HashToken(req.GetRefreshToken())

	// create db client
	id, role, err := repository.RefreshTokenCrud(ctx, hashed_refresh_token)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// generate new access token
	accessToken, err := utils.GenerateAccessJWT(id, role)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// generate new refresh token
	refreshToken, err := utils.GenerateRefreshTokne()

	// hash the new refersh token
	hashed_refresh_token = utils.HashToken(refreshToken)

	// update the old session and insert new one
	err = repository.UpdateRefreshToken(ctx, id, hashed_refresh_token)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  accessToken,
		RefreshToken: "Refresh_token",
	}, nil
}

func (s *Server) Validate(ctx context.Context, req *pb.ValidateRequest) (*pb.ValidateResponse, error) {

	// loading config ile
	cgf := config.Load()

	// get the access token from the request
	tokenStr := req.GetToken()

	// Parse and verify JWT
	token, err := jwt.Parse(tokenStr, func(token *jwt.Token) (interface{}, error) {
		// - check signature
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected Signed Method")
		}

		// - extract claims (user_id, role)
		return []byte(cgf.JWT_key), nil
	})

	// handler parsed token
	if err != nil {
		return nil, err
	}

	// If token is invalid or expired → return error
	if !token.Valid {
		return nil, utils.ErrInvalidToken
	}

	// 5. Return user info from token claims

	cliams, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return nil, utils.ErrInvalidToken
	}

	// extract user_id
	user_id, ok := cliams["user_id"].(string)
	if !ok {
		return nil, utils.ErrInvalidToken
	}

	// extract user_role
	role, ok := cliams["role"].(string)
	if !ok {
		return nil, utils.ErrInvalidToken
	}

	//returning the response
	return &pb.ValidateResponse{
		UserId: user_id,
		Role:   role,
	}, nil
}

func (s *Server) ResetPassword(ctx context.Context, req *pb.ResetPasswordRequest) (*pb.AuthResponse, error) {

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

	// get inputs
	user_id := req.GetId()
	currentPassword := req.GetCurrentPassword()
	newpassword := req.GetNewPassword()

	// update the password and delete all the old sessions
	role, err := repository.ResetPasswordCrud(ctx, user_id, currentPassword, newpassword)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// generate tokens
	token, err := utils.GenerateAccessJWT(user_id, role)
	if err != nil {
		return nil, utils.MapErrors(err)
	}
	ref_token, err := utils.GenerateRefreshTokne()
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// hash the ref_token
	hashed_refresh_token := utils.HashToken(ref_token)

	// store the refresh token
	err = repository.InsertRefreshToken(ctx, user_id, hashed_refresh_token, user_agent, ip_address)
	if err != nil {
		return nil, utils.MapErrors(err)
	}

	// returning the response
	return &pb.AuthResponse{
		AccessToken:  token,
		RefreshToken: hashed_refresh_token,
	}, nil
}
