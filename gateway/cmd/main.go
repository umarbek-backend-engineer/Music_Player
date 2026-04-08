package main

import (
	cgf "gin-server/internal/config"
	grp "gin-server/internal/grpc"
	"gin-server/internal/router"
)

func main() {

	// connect to rabbit mq
	// rb, err := rabbitmq.Connect()
	// if err != nil {
	// 	log.Println("Error in connecting rabbit MQ message broker")
	// 	return
	// }

	grp.InitGRPC()

	port := cgf.Load().Api_Port

	r := router.Route()

	r.Run(":" + port)
}
