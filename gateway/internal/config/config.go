package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
)

// Configuration structure
type Config struct {
	Api_Port                 string
	Api_Host                 string
	Grpc_music_service_port  string
	Grpc_musci_service_host  string
	Grpc_lyrics_service_port string
	Grpc_lyrics_service_host string
	Grpc_Auth_service_port   string
	Grpc_Auth_service_host   string
}

// the function will laod the information from the .env file and assign them to the config  structure
func Load() *Config {

	// loading the .env file
	err := godotenv.Load(".env")
	if err != nil {
		log.Println(err)
		log.Fatal("Error in load configurations")
	}

	// assigning
	return &Config{
		Api_Port:                 os.Getenv("API_PORT"),
		Api_Host:                 os.Getenv("API_HOST"),
		Grpc_music_service_port:  os.Getenv("GRPC_MUSIC_SERVICE_PORT"),
		Grpc_musci_service_host:  os.Getenv("GRPC_MUSIC_SERVICE_HOST"),
		Grpc_lyrics_service_port: os.Getenv("GRPC_LYRICS_SERVICE_PORT"),
		Grpc_lyrics_service_host: os.Getenv("GRPC_LYRICS_SERVICE_HOST"),
		Grpc_Auth_service_port:   os.Getenv("GRPC_AUTH_SERVICE_PORT"),
		Grpc_Auth_service_host:   os.Getenv("GRPC_AUTH_SERVICE_HOST"),
	}
}
