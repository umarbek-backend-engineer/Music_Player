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
		utils.Error(c, "Failed to conenct the client", http.StatusBadGateway, err)
		return
	}

	// give the reponse
	c.JSON(200, resp)
}

func LogIn(c *gin.Context) {
	// get the request context
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	// get the nessessary date to pass to auth-service as Headers
	// in metadata
	md := metadata.New(map[string]string{
		"user-agent": c.Request.UserAgent(),
		"ip-address": c.ClientIP(),
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
		utils.Error(c, "Failed to pass the request", http.StatusBadGateway, err)
		return
	}

	// set cookies
	c.SetSameSite(http.SameSiteLaxMode)

	c.SetCookie("access_token", resp.AccessToken, 3600, "/", "localhost", false, true)
	c.SetCookie("refresh_token", resp.RefreshToken, 1296000, "/", "localhost", false, true)

	// pass the response to the user
	c.JSON(200, resp)
}
