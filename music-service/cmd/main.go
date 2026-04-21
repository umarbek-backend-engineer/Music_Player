package main

import (
	"context"
	"fmt"
	"log"
	"net"

	pb "github.com/umarbek-backend-engineer/Music_Player/music-service/github.com/umarbek-backend-engineer/Music_Player/music-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/config"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/repository/db_connect"
	"github.com/umarbek-backend-engineer/Music_Player/music-service/internal/service"
	"google.golang.org/grpc"
)

func main() {

	conn, err := db_connect.Connect()
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: DataBase connection")
	}
	defer conn.Close(context.Background())

	// 🟢 START CONSUMER (add this)
	// go func() {
	// 	err := service.StartConsumer()
	// 	if err != nil {
	// 		log.Fatal("ERROR: Consumer failed:", err)
	// 	}
	// }()

	// loading port from config gile
	port := fmt.Sprintf(":%s", config.Load().API_Port)

	lis, err := net.Listen(config.Load().NetworkProtocol, port)
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: Initializing server listner")
	}
	defer lis.Close()

	server := grpc.NewServer()

	// registering the functionality of this service
	pb.RegisterMusicServiceServer(server, &service.Server{})

	log.Println("Server is running on port :50051")
	err = server.Serve(lis)
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: Running the server")
	}
}
