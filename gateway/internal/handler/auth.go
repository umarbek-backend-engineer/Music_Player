package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/pkg/utils"
)

func Register(c *gin.Context) {

	// get the request context
	ctx := c.Request.Context()

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
		utils.Error(c, "Failed to conenct the client", http.StatusInternalServerError, err)
		return
	}

	// give the reponse
	c.JSON(200, resp)
}

func LogIn(c *gin.Context) {
	// get the request context
	ctx := c.Request.Context()

	// get the json request body
	var logInRequest pb.LoginRequest
	err := c.BindJSON(&logInRequest)
	if err != nil {
		utils.Error(c, "failed to get requst body", http.StatusBadRequest, err)
		return
	}

	// pass your request to the gRPC service
	logInResponse, err := grpc_init.AuthClient.Login(ctx, &logInRequest)
	if err != nil {
		utils.Error(c, "Failed to pass the request", http.StatusInternalServerError, err)
		return
	}

	// pass the response to the user
	c.JSON(200, logInResponse)

}
