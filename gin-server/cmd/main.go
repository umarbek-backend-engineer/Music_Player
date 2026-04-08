package main

import (
	cgf "gin-server/internal/config"
	rabbitmq "gin-server/internal/rabbit-mq"
	"gin-server/internal/router"
	"log"
)

func main() {

	// connect to rabbit mq
	rb, err := rabbitmq.Connect()
	if err != nil {
		log.Println("Error in connecting rabbit MQ message broker")
		return
	}

	port := cgf.Load().Api_Port

	r := router.Route(rb)

	r.Run(":" + port)
}
