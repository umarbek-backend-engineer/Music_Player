package main

import (
	"context"
	"fmt"
	"log"
	"music-service/internal/config"
	"music-service/internal/repository/db_connect.go"
	"music-service/internal/service"
	pb "music-service/proto/gen"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {

	conn, err := db_connect.Connect()
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: DataBase connection")
	}
	defer conn.Close(context.Background())

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

	reflection.Register(server)

	log.Println("Server is running on port :50051")
	err = server.Serve(lis)
	if err != nil {
		log.Println(err)
		log.Fatal("ERROR: Running the server")
	}
}
