package grpc_init

import (
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"google.golang.org/grpc"
)

var AuthClient pb.AuthServiceClient

func InitauthGRPC() {
	// connection
	conn, err := grpc.Dial("", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatal("Failed to connect to lyrics-service:", err)
	}

	// assigning the variable authclient with new auth service client
	AuthClient = pb.NewAuthServiceClient(conn)
}
