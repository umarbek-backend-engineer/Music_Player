package main

import (
	"log"

	_ "github.com/umarbek-backend-engineer/Music_Player/gateway/docs"
	cgf "github.com/umarbek-backend-engineer/Music_Player/gateway/internal/config"
	grp "github.com/umarbek-backend-engineer/Music_Player/gateway/internal/grpc_init"

	"github.com/umarbek-backend-engineer/Music_Player/gateway/internal/router"
)

// @title SoundWave API
// @version 1.0
// @description SoundWave is a modern music social platform where users can upload, stream, discover, and share their music.
// @description
// @description Features:
// @description • Upload and stream audio tracks
// @description • Public / Private music visibility
// @description • User profiles and social feed
// @description • AI-powered lyrics generation (Whisper)
// @description • Like, search and discover new artists
// @host localhost:9090
// @BasePath /
// @schemes http
// @securityDefinitions.apikey Bearer
// @in header
// @name Authorization
// @description Type "Bearer" followed by a space and JWT token.
// @security Bearer
func main() {

	// connect to rabbit mq
	// rb, err := rabbitmq.Connect()
	// if err != nil {
	// 	log.Println("Error in connecting rabbit MQ message broker")
	// 	return
	// }

	grp.InitMusicGRPC()
	grp.InitLyricsGRPC()
	grp.InitauthGRPC()

	port := cgf.Load().Api_Port

	r := router.Route()

	log.Println("Gateway service is running on port: ", port)

	r.Run(":" + port)
}
