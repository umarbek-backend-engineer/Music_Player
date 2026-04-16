package main

import (
	cgf "gin-server/internal/config"
	grp "gin-server/internal/grpc_init"
	"gin-server/internal/router"
	"log"
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
