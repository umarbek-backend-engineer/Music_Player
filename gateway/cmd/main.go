package main

import (
	cgf "github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
	grp "github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"

	"log"

	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/router"
)

func main() {

	// connect to rabbit mq
	// rb, err := rabbitmq.Connect()
	// if err != nil {
	// 	log.Println("Error in connecting rabbit MQ message broker")
	// 	return
	// }

	grp.InitMusicGRPC()
	grp.InitLyricsGRPC()

	port := cgf.Load().Api_Port

	r := router.Route()

	log.Println("Gateway service is running on port: ", port)

	r.Run(":" + port)
}
