package main

import (
	"context"
	"log"
	"net"

	pb "github.com/umarbek-backend-engineer/Music_Player/lyrics-service/github.com/umarbek-backend-engineer/Music_Player/lyrics-service/proto/gen"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/config"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/repository/posgres"
	"github.com/umarbek-backend-engineer/Music_Player/lyrics-service/internal/service"
	"google.golang.org/grpc"
)

func main() {
	cgf := config.Load()

	// checking the connection of db
	conn, err := posgres.Connect()
	if err != nil {
		log.Println("Failed to connect to db: ", err)
		return
	}
	defer conn.Close(context.Background())

	// creating net http listener
	lis, err := net.Listen(cgf.NetWork_Protocol, ":"+cgf.Api_Port)
	if err != nil {
		log.Println("Failed to create http listener", err)
		return
	}
	defer lis.Close()

	// creating grpc server
	gs := grpc.NewServer()

	pb.RegisterLyricsServiceServer(gs, &service.Server{})

	// runnning the server
	log.Println("Server is running on port: ", cgf.Api_Port)
	err = gs.Serve(lis)
	if err != nil {
		log.Println("Failed to server: ", err)
	}
}
