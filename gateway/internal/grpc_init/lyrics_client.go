package grpc_init

import (
	pb "gin-server/proto/gen"
	"log"

	"google.golang.org/grpc"
)

var LyricsClient pb.LyricsServiceClient

func InitLyricsGRPC() {

	// cgf := config.Load()
	conn, err := grpc.Dial("lyrics-service:50052", grpc.WithInsecure(), grpc.WithBlock())

	if err != nil {
		log.Fatal("Failed to connect to lyrics-service:", err)
	}

	LyricsClient = pb.NewLyricsServiceClient(conn)

}
