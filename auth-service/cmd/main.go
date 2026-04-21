package main

import (
	"context"
	"log"
	"net"

	pb "github.com/umarbek-backend-engineer/Music_Player/github.com/umarbek-backend-engineer/Music_Player/auth-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/internal/repository/postgres"
	"github.com/umarbek-backend-engineer/Music_Player/internal/service"
	"google.golang.org/grpc"
)

func main() {

	conn, err := postgres.Connect()
	if err != nil {
		log.Println(err)
		return
	}
	defer conn.Close(context.Background())

	// creating the listener
	lis, err := net.Listen("tcp", ":50053")
	if err != nil {
		log.Println(err)
		return
	}
	defer lis.Close()

	// creating grpc server
	server := grpc.NewServer()

	// register
	pb.RegisterAuthServiceServer(server, &service.Server{})

	// running the server
	log.Println("Server is running on port 50053")
	err = server.Serve(lis)
	if err != nil {
		log.Println(err)
		return
	}
}
