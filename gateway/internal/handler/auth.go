package handler

import (
	"context"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/pkg/utils"
	"google.golang.org/grpc/metadata"
)

func Register(c *gin.Context) {

	// get the request context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// get the json request to register
	var request pb.RegisterRequest
	err := c.BindJSON(&request)
	if err != nil {
		utils.Error(c, "failed to get requst body", http.StatusBadRequest, err)
		return
	}

	// pass your request to the gRPC service
	resp, err := grpc_init.AuthClient.Register(ctx, &request)
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// give the reponse
	c.ShouldBind(resp)
}

func LogIn(c *gin.Context) {
	// get the request context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// get the nessessary date to pass to auth-service as Headers
	// in metadata
	md := metadata.New(map[string]string{
		"md-user-agent": c.Request.UserAgent(),
		"md-ip-address": c.ClientIP(),
	})

	// put the metadata inside ctx
	ctx = metadata.NewOutgoingContext(ctx, md)

	// get the json request body
	var logInRequest pb.LoginRequest
	err := c.BindJSON(&logInRequest)
	if err != nil {
		utils.Error(c, "failed to get requst body", http.StatusBadRequest, err)
		return
	}

	// pass your request to the gRPC service
	resp, err := grpc_init.AuthClient.Login(ctx, &logInRequest)
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// set cookies
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 1296000, "/", "", false, true)

	// pass the response to the user
	c.ShouldBind(gin.H{
		"message": "Log in successfully",
	})
}

func LogOut(c *gin.Context) {
	// get the user request
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	// get the refresh token from the cookies
	refresh_Token, err := c.Cookie("refresh_token")
	if err != nil {
		utils.Error(c, "Failed to get refresh token", http.StatusUnauthorized, err)
		return
	}

	// send the request to the auth-service
	resp, err := grpc_init.AuthClient.Logout(ctx, &pb.LogoutRequest{RefreshToken: refresh_Token})
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// delete the token from cookies
	c.SetCookie("access_token", "", 0, "/", "", false, true)
	c.SetCookie("refresh_token", "", 0, "/", "", false, true)

	// return the response to the user
	c.ShouldBind(gin.H{
		"log out": resp.Success,
	})
}

func DeleteAccount(c *gin.Context) {

	// get the request context with timeout of 10 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	// get the user id from parametr
	id := c.Param("id")

	// pass the request to the auth-service
	_, err := grpc_init.AuthClient.DeleteAccount(ctx, &pb.DeleteAccountRequest{Id: id})
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// pass the response to the user
	c.ShouldBind(gin.H{
		"message": "Success",
	})
}

func ResetPassword(c *gin.Context) {
	// get the request context with timeout of 10 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	// get the request body
	var RestPassword pb.ResetPasswordRequest
	err := c.BindJSON(&RestPassword)
	if err != nil {
		utils.Error(c, "failed to get requst body", http.StatusBadRequest, err)
		return
	}

	// send the request to the auth-service
	resp, err := grpc_init.AuthClient.ResetPassword(ctx, &RestPassword)
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// set cookies
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 1296000, "/", "", false, true)

	// pass the response to the user
	c.ShouldBind(gin.H{
		"message": "Success",
	})
}

func Refresh(c *gin.Context) {
	// get the request context with timeout of 10 seconds
	ctx, cancel := context.WithTimeout(c.Request.Context(), time.Second*10)
	defer cancel()

	// get the refresh token from the cookies
	refresh_Token, err := c.Cookie("refresh_token")
	if err != nil {
		utils.Error(c, "Failed to get refresh token", http.StatusUnauthorized, err)
		return
	}

	// send the request to the auth-service
	resp, err := grpc_init.AuthClient.Refresh(ctx, &pb.RefreshRequest{RefreshToken: refresh_Token})
	if err != nil {
		utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
		return
	}

	// set cookies
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", resp.AccessToken, 3600, "/", "", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 1296000, "/", "", false, true)

	// pass the response to the user
	c.ShouldBind(gin.H{
		"message": "Success",
	})
}


