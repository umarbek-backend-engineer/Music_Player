package grpc_init

import (
	"log"

	pb "github.com/umarbek-backend-engineer/Music_Player/gateway/github.com/umarbek-backend-engineer/Music_Player/gateway/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
	"google.golang.org/grpc"
)

var MusicClient pb.MusicServiceClient

func InitMusicGRPC() {

	// load the config file
	cgf := config.Load()
	// create a client of that service(connect)
	conn, err := grpc.Dial(cgf.Grpc_musci_service_host+":"+cgf.Grpc_music_service_port, grpc.WithInsecure())
	if err != nil {
		log.Fatal("Failed to connect: ", err)
	}

	// assign the client to ther variable
	MusicClient = pb.NewMusicServiceClient(conn)
}
