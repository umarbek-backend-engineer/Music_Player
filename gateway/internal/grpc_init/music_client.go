package grpc_init

import (
	"gin-server/internal/config"
	pb "gin-server/proto/gen"
	"log"

	"google.golang.org/grpc"
)

var MusicClient pb.MusicServiceClient

func InitMusicGRPC() {
	cgf := config.Load()
	conn, err := grpc.Dial(cgf.Grpc_host+":"+cgf.Grpc_music_service_port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	MusicClient = pb.NewMusicServiceClient(conn)
}
