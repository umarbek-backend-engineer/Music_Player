package grpc_init

import (
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
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
