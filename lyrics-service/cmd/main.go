package main

import (
	"log"
	"lyrics-service/internal/config"
	"lyrics-service/internal/service"
	lyricspb "lyrics-service/proto/gen"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	cgf := config.Load()

	// creating net http listener
	lis, err := net.Listen(cgf.NetWork_Protocol, cgf.Api_Host+":"+cgf.Api_Port)
	if err != nil {
		log.Println("Failed to create http listener", err)
		return
	}
	defer lis.Close()

	// creating grpc server
	gs := grpc.NewServer()

	lyricspb.RegisterLyricsServiceServer(gs, &service.Server{})

	// only in production stage fro postman
	reflection.Register(gs)

	// runnning the server
	log.Println("Server is running on port: ", cgf.Api_Port)
	err = gs.Serve(lis)
	if err != nil {
		log.Println("Failed to server: ", err)
	}
}
