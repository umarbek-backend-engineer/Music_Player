package main

import (
	"context"
	"log"
	"lyrics-service/internal/config"
	"lyrics-service/internal/repository/posgres"
	"lyrics-service/internal/service"
	lyricspb "lyrics-service/proto/gen"
	"net"

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

	lyricspb.RegisterLyricsServiceServer(gs, &service.Server{})

	// runnning the server
	log.Println("Server is running on port: ", cgf.Api_Port)
	err = gs.Serve(lis)
	if err != nil {
		log.Println("Failed to server: ", err)
	}
}
