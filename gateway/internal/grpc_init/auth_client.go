package grpc_init

import (
	"fmt"
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
	"google.golang.org/grpc"
)

var AuthClient pb.AuthServiceClient

func InitauthGRPC() {

	// load the config file
	cgf := config.Load()

	address := fmt.Sprintf("%s:%s", cgf.Grpc_Auth_service_host, cgf.Grpc_Auth_service_port)

	// connection
	conn, err := grpc.Dial(address, grpc.WithInsecure())

	if err != nil {
		log.Fatal("Failed to connect to lyrics-service:", err)
	}

	// assigning the variable authclient with new auth service client
	AuthClient = pb.NewAuthServiceClient(conn)
}
