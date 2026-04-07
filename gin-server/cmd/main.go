package main

import (
	cgf "gin-server/internal/config"
	rabbitmq "gin-server/internal/rabbit-mq"
	"gin-server/internal/router"
	"log"
)

func main() {

	// connect to rabbit mq
	rabbit, err := rabbitmq.Connect()
	if err != nil {
		log.Println("Error in connecting to rabbitMQ", err)
		return
	}
	defer rabbit.Conn.Close()
	defer rabbit.Ch.Close()

	port := cgf.Load().Api_Port

	r := router.Route(rabbit)

	r.Run(":" + port)
}
