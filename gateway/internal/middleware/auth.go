package middleware

import (
	"context"
	"errors"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/pkg/utils"
)

func Authentication() gin.HandlerFunc {
	return func(c *gin.Context) {
		// get the request context
		ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
		defer cancel()

		// get the access token from the cookies
		access_token, err := c.Cookie("access_token")
		if err != nil {
			utils.Error(c, "Missing access token", http.StatusUnauthorized, err)
			return
		}

		// validate refersh token
		if access_token == "" {
			utils.Error(c, "Missing access token", http.StatusUnauthorized, errors.New("Missing Access Token"))
			return
		}

		// call the grpc validator
		resp, err := grpc_init.AuthClient.Validate(ctx, &pb.ValidateRequest{Token: access_token})
		if err != nil {
			utils.Error(c, "Internal Error in auth-service", http.StatusBadGateway, err)
			return
		}

		// attach the user information to the context
		c.Set("user_id", resp.UserId)
		c.Set("role", resp.Role)

		c.Next()
	}
}
