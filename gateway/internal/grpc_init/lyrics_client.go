package grpc_init

import (
	"fmt"
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
	"google.golang.org/grpc"
)

var LyricsClient pb.LyricsServiceClient

func InitLyricsGRPC() {

	// load the config file
	cgf := config.Load()

	// make and address
	address := fmt.Sprintf("%s:%s", cgf.Grpc_lyrics_service_host, cgf.Grpc_lyrics_service_port)

	// connect to the service
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect to lyrics-service:", err)
	}

	LyricsClient = pb.NewLyricsServiceClient(conn)

}
